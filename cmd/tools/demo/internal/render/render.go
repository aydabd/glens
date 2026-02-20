// Package render prints demo output for the demo tool.
package render

import (
	"fmt"
	"strings"

	"glens/tools/demo/internal/loader"
)

// Banner prints the glens demo banner.
func Banner() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════════════╗
║                                                                       ║
║              🚀 GLENS MODERN AI TEST GENERATION 🚀                   ║
║                                                                       ║
╚═══════════════════════════════════════════════════════════════════════╝`)
	fmt.Println()
}

// SpecInfo prints API metadata from a parsed spec.
func SpecInfo(spec *loader.Spec) {
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

// Endpoints prints all endpoints found in the spec.
func Endpoints(spec *loader.Spec) {
	count := len(spec.Paths)
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

// ModelComparison prints a table of available AI models.
func ModelComparison() {
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

// SampleTest prints a sample generated test snippet.
func SampleTest() {
	fmt.Println("─── Sample Generated Test ────────────────────────────────────")
	fmt.Print(`
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
