// Command demo demonstrates glens OpenAPI parsing capabilities.
// It is a cross-platform replacement for scripts/demo_modern.sh.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const version = "0.1.0"

// minimalSpec holds only what we need for the demo output.
type minimalSpec struct {
	Info struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Version     string `json:"version"`
	} `json:"info"`
	Servers []struct {
		URL string `json:"url"`
	} `json:"servers"`
	Paths map[string]map[string]struct {
		Summary string   `json:"summary"`
		Tags    []string `json:"tags"`
	} `json:"paths"`
}

func main() {
	var specPath string
	var showVersion bool

	flag.StringVar(&specPath, "spec", "", "path to OpenAPI spec file or URL")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: demo [flags] [spec-path]\n\n")
		fmt.Fprintf(os.Stderr, "Demonstrates glens OpenAPI parsing capabilities.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  demo --spec test_specs/sample_api.json\n")
		fmt.Fprintf(os.Stderr, "  demo https://petstore3.swagger.io/api/v3/openapi.json\n")
	}
	flag.Parse()

	if showVersion {
		fmt.Printf("glens-demo version %s\n", version)
		return
	}

	if specPath == "" && flag.NArg() > 0 {
		specPath = flag.Arg(0)
	}

	if specPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := runDemo(specPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runDemo(specPath string) error {
	printBanner()
	fmt.Printf("Parsing OpenAPI spec: %s\n\n", specPath)

	spec, err := loadSpec(specPath)
	if err != nil {
		return fmt.Errorf("failed to load spec: %w", err)
	}

	printSpecInfo(spec)
	printEndpoints(spec)
	printModelComparison()
	printSampleTest()

	return nil
}

func loadSpec(source string) (*minimalSpec, error) {
	var data []byte
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		client := &http.Client{Timeout: 30 * time.Second}
		resp, httpErr := client.Get(source) //nolint:gosec
		if httpErr != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", httpErr)
		}
		defer resp.Body.Close()
		data, err = io.ReadAll(resp.Body)
	} else {
		data, err = os.ReadFile(source) //nolint:gosec
	}
	if err != nil {
		return nil, err
	}

	var spec minimalSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &spec, nil
}

func printBanner() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════════════╗
║                                                                       ║
║              🚀 GLENS MODERN AI TEST GENERATION 🚀                   ║
║                                                                       ║
╚═══════════════════════════════════════════════════════════════════════╝`)
	fmt.Println()
}

func printSpecInfo(spec *minimalSpec) {
	fmt.Println("─── API Information ──────────────────────────────────────────")
	fmt.Printf("  Title:   %s\n", spec.Info.Title)
	fmt.Printf("  Version: %s\n", spec.Info.Version)
	if spec.Info.Description != "" {
		desc := spec.Info.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		fmt.Printf("  Desc:    %s\n", desc)
	}
	if len(spec.Servers) > 0 {
		fmt.Printf("  Server:  %s\n", spec.Servers[0].URL)
	}
	fmt.Println()
}

func printEndpoints(spec *minimalSpec) {
	count := 0
	for range spec.Paths {
		count++
	}
	fmt.Printf("─── Endpoints (%d paths) ─────────────────────────────────────\n", count)
	i := 1
	for path, methods := range spec.Paths {
		for method, op := range methods {
			tags := ""
			if len(op.Tags) > 0 {
				tags = fmt.Sprintf(" [%s]", strings.Join(op.Tags, ", "))
			}
			summary := op.Summary
			if len(summary) > 40 {
				summary = summary[:37] + "..."
			}
			fmt.Printf("  %2d. %-6s %-35s %s%s\n", i, strings.ToUpper(method), path, summary, tags)
			i++
		}
	}
	fmt.Println()
}

func printModelComparison() {
	fmt.Println("─── Available AI Models ──────────────────────────────────────")
	fmt.Println()
	fmt.Println("  Provider    Model                    Cost/1M tokens  Speed")
	fmt.Println("  ─────────── ──────────────────────── ──────────────  ─────")
	fmt.Println("  OpenAI      gpt-4o                   $5.00           Fast")
	fmt.Println("  OpenAI      gpt-4o-mini              $0.15           Fast")
	fmt.Println("  Anthropic   claude-3.5-sonnet         $3.00           Fast")
	fmt.Println("  Google      gemini-2.0-flash          $0.00 (free)    Very Fast")
	fmt.Println("  Google      gemini-2.0-pro            $1.25           Fast")
	fmt.Println("  Local       enhanced-mock             $0.00 (free)    Very Fast")
	fmt.Println("  Local       ollama:*                  $0.00 (free)    Depends")
	fmt.Println()
}

func printSampleTest() {
	fmt.Println("─── Sample Generated Test ────────────────────────────────────")
	fmt.Println(`
  func TestGETEndpoint(t *testing.T) {
      client := &http.Client{Timeout: 10 * time.Second}
      resp, err := client.Get(baseURL + "/endpoint")
      require.NoError(t, err)
      defer resp.Body.Close()
      assert.Equal(t, http.StatusOK, resp.StatusCode)
  }
`)
	fmt.Println("─── Quick Start ──────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("  # Offline demo (no API key needed):")
	fmt.Println("  glens analyze <spec> --ai-models=enhanced-mock")
	fmt.Println()
	fmt.Println("  # With OpenAI:")
	fmt.Println("  export OPENAI_API_KEY=sk-...")
	fmt.Println("  glens analyze <spec> --ai-models=gpt-4o")
	fmt.Println()
}
