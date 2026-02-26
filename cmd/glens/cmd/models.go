package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"glens/tools/glens/internal/ai"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Manage AI models (list, download, status)",
	Long: `Manage AI models including local LLMs via Ollama and cloud providers.
Supports listing available models, checking status, and downloading models.`,
}

var modelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available AI models",
	Long:  `List all available AI models including local Ollama models and cloud providers.`,
	RunE:  runModelsList,
}

var modelsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check status of AI model providers",
	Long:  `Check the health and availability of AI model providers including Ollama and cloud APIs.`,
	RunE:  runModelsStatus,
}

var modelsOllamaCmd = &cobra.Command{
	Use:   "ollama",
	Short: "Ollama-specific commands",
	Long:  `Commands for managing Ollama local LLM models.`,
}

var modelsOllamaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed Ollama models",
	Long:  `List all models currently installed in Ollama.`,
	RunE:  runOllamaList,
}

var modelsOllamaPullCmd = &cobra.Command{
	Use:   "pull [model]",
	Short: "Pull an Ollama model from the registry",
	Long: `Pull (download) a model from the Ollama registry for fully local use.
No cloud account or API key is required.

Examples:
  glens models ollama pull mistral          # Mistral 7B
  glens models ollama pull mistral-nemo     # Mistral Nemo 12B
  glens models ollama pull llama3           # Meta Llama 3
  glens models ollama pull phi4             # Microsoft Phi-4
  glens models ollama pull gemma2           # Google Gemma 2
  glens models ollama pull codellama:7b-instruct`,
	Args: cobra.ExactArgs(1),
	RunE: runOllamaPull,
}

var modelsOllamaStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Ollama server status",
	Long:  `Check if Ollama server is running and accessible.`,
	RunE:  runOllamaStatus,
}

func init() {
	rootCmd.AddCommand(modelsCmd)

	// Add subcommands
	modelsCmd.AddCommand(modelsListCmd)
	modelsCmd.AddCommand(modelsStatusCmd)
	modelsCmd.AddCommand(modelsOllamaCmd)

	// Add Ollama subcommands
	modelsOllamaCmd.AddCommand(modelsOllamaListCmd)
	modelsOllamaCmd.AddCommand(modelsOllamaStatusCmd)
	modelsOllamaCmd.AddCommand(modelsOllamaPullCmd)
}

func runModelsList(_ *cobra.Command, _ []string) error {
	fmt.Println("📋 Available AI Models")
	fmt.Println("=====================")

	// Cloud providers
	fmt.Println("\n🌐 Cloud Providers (require API keys):")
	fmt.Println("  • gpt4         - OpenAI GPT-4 Turbo")
	fmt.Println("  • sonnet4      - Anthropic Claude 3.5 Sonnet")
	fmt.Println("  • flash-pro    - Google Gemini 1.5 Flash Pro")
	fmt.Println("  • mistral      - Mistral AI (cloud)")

	// Local open-source model shortcuts
	fmt.Println("\n🔓 Local Open-Source Models (no cloud/API-key required):")
	fmt.Println("  Mistral:")
	fmt.Println("    • mistral-local, mistral7b    → mistral (7B)")
	fmt.Println("    • mistral-nemo-local          → mistral-nemo (12B)")
	fmt.Println("    • mistral-small-local         → mistral-small")
	fmt.Println("  Meta Llama:")
	fmt.Println("    • llama3, llama3-local        → llama3")
	fmt.Println("    • llama3.1, llama3.1-local    → llama3.1")
	fmt.Println("    • llama3.2, llama3.2-local    → llama3.2")
	fmt.Println("  Microsoft Phi:")
	fmt.Println("    • phi3, phi3-local            → phi3")
	fmt.Println("    • phi4, phi4-local            → phi4")
	fmt.Println("  Google Gemma (open-weights):")
	fmt.Println("    • gemma2, gemma2-local        → gemma2")
	fmt.Println("    • gemma3, gemma3-local        → gemma3")
	fmt.Println("  Custom: ollama:<model>          e.g. ollama:mistral:7b-instruct")
	fmt.Println("\n💡 Pull a model first:  glens models ollama pull <model-name>")

	// Check Ollama models
	fmt.Println("\n🏠 Installed Ollama Models:")

	ollamaClient, err := ai.NewOllamaClient("")
	if err != nil {
		fmt.Printf("  ❌ Ollama not configured: %v\n", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	models, err := ollamaClient.ListModels(ctx)
	if err != nil {
		fmt.Printf("  ❌ Ollama not accessible: %v\n", err)
		fmt.Println("\n💡 Install Ollama: https://ollama.ai")
		fmt.Println("💡 Then pull a model: glens models ollama pull mistral")
		return nil
	}

	if len(models) == 0 {
		fmt.Println("  📭 No models installed")
		fmt.Println("\n💡 Recommended models (no API key needed):")
		fmt.Println("     glens models ollama pull mistral")
		fmt.Println("     glens models ollama pull llama3")
		fmt.Println("     glens models ollama pull codellama:7b-instruct")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "  Model\tSize\tModified"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := fmt.Fprintln(w, "  -----\t----\t--------"); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	for _, model := range models {
		size := formatSize(model.Size)
		modified := model.ModifiedAt.Format("2006-01-02 15:04")
		if _, err := fmt.Fprintf(w, "  %s\t%s\t%s\n", model.Name, size, modified); err != nil {
			return fmt.Errorf("failed to write model data: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}

	fmt.Println("\n💡 Usage examples:")
	fmt.Printf("  glens analyze spec.json --ai-models mistral-local\n")
	fmt.Printf("  glens analyze spec.json --ai-models ollama:%s\n", models[0].Name)
	fmt.Printf("  glens analyze spec.json --ai-models llama3-local,mistral-local\n")

	return nil
}

func runModelsStatus(_ *cobra.Command, _ []string) error {
	fmt.Println("🔍 AI Model Provider Status")
	fmt.Println("===========================")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check Ollama
	fmt.Print("\n🏠 Ollama: ")
	ollamaClient, err := ai.NewOllamaClient("")
	if err != nil {
		fmt.Printf("❌ Not configured (%v)\n", err)
	} else {
		if err := ollamaClient.HealthCheck(ctx); err != nil {
			fmt.Printf("❌ Unhealthy (%v)\n", err)
		} else {
			fmt.Println("✅ Available")

			// Show available models
			if models, err := ollamaClient.ListModels(ctx); err == nil {
				fmt.Printf("   📦 %d models installed\n", len(models))
			}
		}
	}

	// Check cloud providers (would need API keys to test)
	fmt.Println("\n🌐 Cloud Providers:")
	fmt.Println("   🤖 OpenAI: API key required")
	fmt.Println("   🧠 Anthropic: API key required")
	fmt.Println("   🌟 Google: Credentials required")

	fmt.Println("\n💡 To test cloud providers, set environment variables:")
	fmt.Println("   export OPENAI_API_KEY=your_key")
	fmt.Println("   export ANTHROPIC_API_KEY=your_key")
	fmt.Println("   export GOOGLE_APPLICATION_CREDENTIALS=path/to/credentials.json")

	return nil
}

func runOllamaList(_ *cobra.Command, _ []string) error {
	ollamaClient, err := ai.NewOllamaClient("")
	if err != nil {
		return fmt.Errorf("failed to create Ollama client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	models, err := ollamaClient.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("failed to list Ollama models: %w", err)
	}

	if len(models) == 0 {
		fmt.Println("📭 No Ollama models installed")
		fmt.Println("\n💡 Install recommended coding models:")
		fmt.Println("   ollama pull codellama:7b-instruct")
		fmt.Println("   ollama pull deepseek-coder:6.7b-instruct")
		fmt.Println("   ollama pull qwen2.5-coder:7b-instruct")
		return nil
	}

	fmt.Printf("📦 Found %d Ollama models:\n\n", len(models))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintln(w, "Model\tSize\tModified\tDigest"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := fmt.Fprintln(w, "-----\t----\t--------\t------"); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	for _, model := range models {
		size := formatSize(model.Size)
		modified := model.ModifiedAt.Format("2006-01-02 15:04")
		digest := model.Digest[:12] + "..."
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", model.Name, size, modified, digest); err != nil {
			return fmt.Errorf("failed to write model data: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}

	return nil
}

func runOllamaStatus(_ *cobra.Command, _ []string) error {
	ollamaClient, err := ai.NewOllamaClient("")
	if err != nil {
		return fmt.Errorf("failed to create Ollama client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Print("🔍 Checking Ollama status... ")

	if err := ollamaClient.HealthCheck(ctx); err != nil {
		fmt.Printf("❌ Failed\n\nError: %v\n", err)
		fmt.Println("\n💡 Make sure Ollama is installed and running:")
		fmt.Println("   • Install: https://ollama.ai")
		fmt.Println("   • Start: ollama serve")
		fmt.Println("   • Pull a model: ollama pull codellama:7b-instruct")
		return nil
	}

	fmt.Println("✅ Healthy")

	// Get model count
	if models, err := ollamaClient.ListModels(ctx); err == nil {
		fmt.Printf("📦 %d models available\n", len(models))

		if len(models) > 0 {
			fmt.Println("\n🎯 Recommended for code generation:")
			for _, model := range models {
				if isCodeModel(model.Name) {
					fmt.Printf("   ✅ %s\n", model.Name)
				}
			}
		}
	}

	return nil
}

// formatSize converts bytes to human readable format
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// isCodeModel checks if a model is recommended for code generation
func isCodeModel(name string) bool {
	codeModels := []string{
		"codellama", "deepseek-coder", "qwen2.5-coder", "codeqwen",
		"starcoder", "wizard-coder", "phind-codellama",
	}

	for _, model := range codeModels {
		if len(name) >= len(model) && name[:len(model)] == model {
			return true
		}
	}
	return false
}

func runOllamaPull(_ *cobra.Command, args []string) error {
	modelName := args[0]

	ollamaClient, err := ai.NewOllamaClient("")
	if err != nil {
		return fmt.Errorf("failed to create Ollama client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	fmt.Printf("⬇️  Pulling model %q from Ollama registry...\n", modelName)
	fmt.Println("   (This may take several minutes for large models)")

	if err := ollamaClient.PullModel(ctx, modelName, os.Stdout); err != nil {
		fmt.Printf("\n❌ Failed to pull model: %v\n", err)
		fmt.Println("\n💡 Make sure Ollama is running: ollama serve")
		return nil
	}

	fmt.Printf("\n✅ Model %q pulled successfully\n", modelName)
	fmt.Printf("\n💡 Use it with:\n")
	fmt.Printf("   glens analyze spec.json --ai-models ollama:%s\n", modelName)
	return nil
}
