package cmd

import (
	"context"
	"fmt"
	"strings"

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
2. Generates integration tests using AI models (defaults to GPT-4 only)
3. Executes tests against the implementation
4. Creates GitHub issues ONLY for endpoints where tests fail
5. Generates comparison reports

GitHub issues are created only when tests fail, indicating a mismatch
between the OpenAPI specification and the actual implementation.`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringSlice("ai-models", []string{"gpt4"}, "AI models to use for test generation (gpt4, ollama, ollama:model-name, etc.)")
	analyzeCmd.Flags().String("github-repo", "", "GitHub repository in owner/repo format (can also use GITHUB_REPOSITORY env var)")
	analyzeCmd.Flags().String("test-framework", "testify", "Test framework to use (testify, ginkgo)")
	analyzeCmd.Flags().Bool("create-issues", true, "Create GitHub issues when tests fail (requires github-repo and GITHUB_TOKEN)")
	analyzeCmd.Flags().Bool("run-tests", true, "Execute generated tests")
	analyzeCmd.Flags().String("output", "reports/report.md", "Output file for the final report")

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

func runAnalyze(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	openapiURL := args[0]

	// Handle github repository with proper precedence: CLI flag > env var > config file
	// If CLI flag is explicitly set, it should override config file values
	if cmd.Flags().Changed("github-repo") {
		flagValue, _ := cmd.Flags().GetString("github-repo")
		viper.Set("github.repository", flagValue)
	}

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

		// Set the target repository
		repo := viper.GetString("github.repository")
		if repo == "" {
			return fmt.Errorf("github repository is required when create-issues is enabled (use --github-repo flag or GITHUB_REPOSITORY env var)")
		}
		if err := githubClient.SetRepository(repo); err != nil {
			return fmt.Errorf("failed to set github repository: %w", err)
		}

		log.Info().
			Str("repository", repo).
			Msg("GitHub client configured")
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

		// Track if we should create an issue (only if tests fail)
		hasFailedTests := false
		failedModels := []string{}

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
					// Check if this is a real test failure, not just connection/setup issues
					if isRealTestFailure(err, execResult) {
						hasFailedTests = true
						failedModels = append(failedModels, modelName)
					}
				} else {
					testResult.ExecutionResult = execResult
					log.Info().
						Str("ai_model", modelName).
						Bool("passed", execResult.Passed).
						Dur("duration", execResult.Duration).
						Msg("Test execution completed")

					// Check if tests failed (not passed and has actual test failures)
					if execResult.Failed && (execResult.FailureCount > 0 || execResult.ErrorCount > 0) {
						hasFailedTests = true
						failedModels = append(failedModels, modelName)
					}
				}
			}

			result.Tests[modelName] = testResult
		}

		// Create GitHub issue ONLY if tests failed
		if githubClient != nil && hasFailedTests {
			log.Info().
				Str("endpoint", fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)).
				Strs("failed_models", failedModels).
				Msg("Creating GitHub issue for failed tests")

			issueNumber, err := githubClient.CreateEndpointIssue(ctx, endpoint, failedModels)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create GitHub issue")
			} else {
				result.IssueNumber = issueNumber
				log.Info().
					Int("issue_number", issueNumber).
					Msg("GitHub issue created for test failures")

				// Update issue with test results
				resultsComment := formatTestFailureResults(result, failedModels)
				if err := githubClient.UpdateIssueWithResults(ctx, issueNumber, resultsComment); err != nil {
					log.Error().Err(err).Msg("Failed to update issue with results")
				}
			}
		} else if githubClient != nil && !hasFailedTests {
			log.Info().
				Str("endpoint", fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)).
				Msg("All tests passed - no issue created")
		}

		results = append(results, result)
	}

	// Generate final report
	log.Info().Msg("Generating final report")
	report := reporter.GenerateReport(spec, results)

	outputFile := viper.GetString("output")

	// Ensure the reports directory exists
	if err := reporter.EnsureReportDirectory(outputFile); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}

	if err := reporter.WriteReport(report, outputFile); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	log.Info().
		Str("output_file", outputFile).
		Int("endpoints_processed", len(results)).
		Msg("Analysis completed successfully")

	return nil
}

// isRealTestFailure determines if an error represents a real test failure
// against the OpenAPI spec, not just connection or setup issues
func isRealTestFailure(err error, result *generator.ExecutionResult) bool {
	if err == nil {
		return false
	}

	// If we have execution results with actual test failures, it's a real failure
	if result != nil && (result.FailureCount > 0 || result.ErrorCount > 0) {
		// Check if errors are not just compilation or setup errors
		for _, testErr := range result.Errors {
			// These are real test failures, not setup issues
			if testErr.Type != "error" && testErr.TestName != "compilation" {
				return true
			}
		}
	}

	errMsg := err.Error()

	// These are setup/infrastructure issues, not real test failures:
	setupIssues := []string{
		"failed to create temp directory",
		"failed to write test file",
		"failed to create test module",
		"go mod tidy failed",
		"compilation",
	}

	for _, issue := range setupIssues {
		if strings.Contains(strings.ToLower(errMsg), issue) {
			return false
		}
	}

	// If we can't determine, consider it a real failure to be safe
	return true
}

// formatTestFailureResults formats test failure information for GitHub issue
func formatTestFailureResults(result reporter.EndpointResult, failedModels []string) string {
	var sb strings.Builder

	sb.WriteString("## Test Execution Results\n\n")
	fmt.Fprintf(&sb, "**Endpoint:** `%s %s`\n\n", result.Endpoint.Method, result.Endpoint.Path)

	for _, modelName := range failedModels {
		if testResult, ok := result.Tests[modelName]; ok {
			fmt.Fprintf(&sb, "### ❌ %s - Tests Failed\n\n", modelName)

			if testResult.ExecutionResult != nil {
				execResult := testResult.ExecutionResult
				fmt.Fprintf(&sb, "- **Test Count:** %d\n", execResult.TestCount)
				fmt.Fprintf(&sb, "- **Failures:** %d\n", execResult.FailureCount)
				fmt.Fprintf(&sb, "- **Errors:** %d\n", execResult.ErrorCount)
				fmt.Fprintf(&sb, "- **Duration:** %s\n\n", execResult.Duration)

				if len(execResult.Errors) > 0 {
					sb.WriteString("#### Failed Tests:\n\n")
					for _, testErr := range execResult.Errors {
						fmt.Fprintf(&sb, "**%s** (%s):\n", testErr.TestName, testErr.Type)
						fmt.Fprintf(&sb, "```\n%s\n```\n\n", testErr.Message)
					}
				}

				if execResult.Output != "" {
					sb.WriteString("<details>\n<summary>Full Test Output</summary>\n\n")
					fmt.Fprintf(&sb, "```\n%s\n```\n", execResult.Output)
					sb.WriteString("</details>\n\n")
				}
			} else if testResult.ExecutionError != "" {
				fmt.Fprintf(&sb, "**Execution Error:**\n```\n%s\n```\n\n", testResult.ExecutionError)
			}

			sb.WriteString("---\n\n")
		}
	}

	return sb.String()
}
