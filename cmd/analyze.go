package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"glens/pkg/ai"
	"glens/pkg/generator"
	"glens/pkg/github"
	"glens/pkg/parser"
	"glens/pkg/reporter"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [openapi-url]",
	Short: "Analyze OpenAPI specification and generate integration tests",
	Long: `Analyzes an OpenAPI specification from a URL or file path and:
1. Parses the OpenAPI spec to extract endpoints
2. Creates GitHub issues for each endpoint
3. Generates integration tests using AI models (defaults to GPT-4 only)
4. Executes tests and creates comparison reports`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringSlice("ai-models", []string{"gpt4"}, "AI models to use for test generation (gpt4, ollama, ollama:model-name, etc.)")
	analyzeCmd.Flags().String("github-repo", "", "GitHub repository (owner/repo)")
	analyzeCmd.Flags().String("test-framework", "testify", "Test framework to use (testify, ginkgo)")
	analyzeCmd.Flags().Bool("create-issues", true, "Create GitHub issues for endpoints")
	analyzeCmd.Flags().Bool("run-tests", true, "Execute generated tests")
	analyzeCmd.Flags().String("output", "report.md", "Output file for the final report")

	// Endpoint filtering options
	analyzeCmd.Flags().String("op-id", "", "Target specific endpoint by operation ID (e.g., getPetById, addPet)")

	_ = viper.BindPFlag("ai_models", analyzeCmd.Flags().Lookup("ai-models"))
	_ = viper.BindPFlag("github.repository", analyzeCmd.Flags().Lookup("github-repo"))
	_ = viper.BindPFlag("test_framework", analyzeCmd.Flags().Lookup("test-framework"))
	_ = viper.BindPFlag("create_issues", analyzeCmd.Flags().Lookup("create-issues"))
	_ = viper.BindPFlag("run_tests", analyzeCmd.Flags().Lookup("run-tests"))
	_ = viper.BindPFlag("output", analyzeCmd.Flags().Lookup("output"))
	_ = viper.BindPFlag("op_id", analyzeCmd.Flags().Lookup("op-id"))
}

func runAnalyze(_ *cobra.Command, args []string) error {
	ctx := context.Background()
	openapiURL := args[0]

	log.Info().
		Str("openapi_url", openapiURL).
		Strs("ai_models", viper.GetStringSlice("ai_models")).
		Str("github_repo", viper.GetString("github.repository")).
		Msg("Starting OpenAPI analysis")

	// Parse OpenAPI specification
	log.Info().Msg("Parsing OpenAPI specification")
	spec, err := parser.ParseOpenAPISpec(openapiURL)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	log.Info().
		Int("endpoints_count", len(spec.Endpoints)).
		Msg("OpenAPI specification parsed successfully")

	// Initialize GitHub client
	var githubClient *github.Client
	if viper.GetBool("create_issues") {
		log.Info().Msg("Initializing GitHub client")
		githubClient, err = github.NewClient(viper.GetString("github.token"))
		if err != nil {
			return fmt.Errorf("failed to initialize GitHub client: %w", err)
		}
	}

	// Initialize AI clients
	log.Info().Msg("Initializing AI model clients")
	aiManager, err := ai.NewManager(viper.GetStringSlice("ai_models"))
	if err != nil {
		return fmt.Errorf("failed to initialize AI clients: %w", err)
	}

	// Initialize test generator
	testGen := generator.NewTestGenerator(viper.GetString("test_framework"))

	// Filter endpoints if operation ID is specified
	var endpointsToProcess []parser.Endpoint
	opID := viper.GetString("op_id")

	if opID != "" {
		log.Info().
			Str("operation_id", opID).
			Msg("Filtering endpoints by operation ID")

		// Find endpoint with matching operation ID
		found := false
		for i := range spec.Endpoints {
			endpoint := &spec.Endpoints[i]
			if endpoint.OperationID == opID {
				endpointsToProcess = append(endpointsToProcess, *endpoint)
				found = true
				break
			}
		}

		if !found {
			// List available operation IDs to help user
			var availableOps []string
			for i := range spec.Endpoints {
				endpoint := &spec.Endpoints[i]
				if endpoint.OperationID != "" {
					availableOps = append(availableOps, endpoint.OperationID)
				}
			}
			return fmt.Errorf("operation ID '%s' not found. Available operation IDs: %v", opID, availableOps)
		}

		log.Info().
			Str("operation_id", opID).
			Int("matching_endpoints", len(endpointsToProcess)).
			Msg("Found matching endpoint")
	} else {
		// Process all endpoints
		endpointsToProcess = spec.Endpoints
	}

	// Process each endpoint
	var results []reporter.EndpointResult

	for i := range endpointsToProcess {
		endpoint := &endpointsToProcess[i]
		log.Info().
			Str("method", endpoint.Method).
			Str("path", endpoint.Path).
			Msg("Processing endpoint")

		result := reporter.EndpointResult{
			Endpoint: *endpoint,
			Tests:    make(map[string]reporter.TestResult),
		}

		// Create GitHub issue if enabled
		if githubClient != nil {
			issueNumber, err := githubClient.CreateEndpointIssue(ctx, endpoint, viper.GetStringSlice("ai_models"))
			if err != nil {
				log.Error().Err(err).Msg("Failed to create GitHub issue")
			} else {
				result.IssueNumber = issueNumber
				log.Info().Int("issue_number", issueNumber).Msg("GitHub issue created")
			}
		}

		// Generate and run tests for each AI model
		for _, modelName := range viper.GetStringSlice("ai_models") {
			log.Info().
				Str("ai_model", modelName).
				Str("endpoint", fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)).
				Msg("Generating tests with AI model")

			testCode, prompt, err := aiManager.GenerateTest(ctx, modelName, endpoint)
			if err != nil {
				log.Error().
					Err(err).
					Str("ai_model", modelName).
					Msg("Failed to generate test")
				continue
			}

			testResult := reporter.TestResult{
				AIModel:   modelName,
				Prompt:    prompt,
				TestCode:  testCode,
				Framework: viper.GetString("test_framework"),
			}

			// Execute test if enabled
			if viper.GetBool("run_tests") {
				log.Info().
					Str("ai_model", modelName).
					Msg("Executing generated test")

				execResult, err := testGen.ExecuteTest(ctx, testCode, endpoint)
				if err != nil {
					log.Error().
						Err(err).
						Str("ai_model", modelName).
						Msg("Test execution failed")
					testResult.ExecutionError = err.Error()
				} else {
					testResult.ExecutionResult = execResult
					log.Info().
						Str("ai_model", modelName).
						Bool("passed", execResult.Passed).
						Dur("duration", execResult.Duration).
						Msg("Test execution completed")
				}
			}

			result.Tests[modelName] = testResult
		}

		results = append(results, result)
	}

	// Generate final report
	log.Info().Msg("Generating final report")
	report := reporter.GenerateReport(spec, results)

	outputFile := viper.GetString("output")
	if err := reporter.WriteReport(report, outputFile); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	log.Info().
		Str("output_file", outputFile).
		Int("endpoints_processed", len(results)).
		Msg("Analysis completed successfully")

	return nil
}
