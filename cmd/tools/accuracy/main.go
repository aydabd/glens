// Command accuracy evaluates OpenAPI spec parsing accuracy.
// It is a cross-platform replacement for scripts/test_accuracy.sh.
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

// minimalSpec holds only the fields needed for accuracy metrics.
type minimalSpec struct {
	Info struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`
	Paths map[string]map[string]interface{} `json:"paths"`
}

// result holds the outcome for a single spec analysis.
type result struct {
	name      string
	specPath  string
	title     string
	endpoints int
	elapsed   time.Duration
	err       error
}

func main() {
	var outputFile string
	var showVersion bool

	flag.StringVar(&outputFile, "output", "", "write markdown report to file (default: stdout)")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: accuracy [flags] <spec> [spec...]\n\n")
		fmt.Fprintf(os.Stderr, "Evaluates OpenAPI spec parsing accuracy and generates a report.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  accuracy test_specs/sample_api.json\n")
		fmt.Fprintf(os.Stderr, "  accuracy --output report.md spec1.json spec2.json\n")
	}
	flag.Parse()

	if showVersion {
		fmt.Printf("glens-accuracy version %s\n", version)
		return
	}

	specs := flag.Args()
	if len(specs) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	results := analyzeSpecs(specs)
	report := buildReport(results)

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(report), 0o600); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing report: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Report written to %s\n", outputFile)
	} else {
		fmt.Print(report)
	}

	for _, r := range results {
		if r.err != nil {
			os.Exit(1)
		}
	}
}

func analyzeSpecs(specs []string) []result {
	results := make([]result, 0, len(specs))
	for _, specPath := range specs {
		start := time.Now()
		spec, err := loadSpec(specPath)
		elapsed := time.Since(start)

		r := result{
			name:    specName(specPath),
			specPath: specPath,
			elapsed: elapsed,
			err:     err,
		}
		if err == nil {
			r.title = spec.Info.Title
			r.endpoints = countEndpoints(spec)
		}
		results = append(results, r)
	}
	return results
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
	if jsonErr := json.Unmarshal(data, &spec); jsonErr != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", jsonErr)
	}
	return &spec, nil
}

func countEndpoints(spec *minimalSpec) int {
	count := 0
	for _, methods := range spec.Paths {
		count += len(methods)
	}
	return count
}

func buildReport(results []result) string {
	var sb strings.Builder
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05 UTC")

	sb.WriteString("# Glens Accuracy Report\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", timestamp))

	total := len(results)
	passed := 0
	totalEndpoints := 0
	for _, r := range results {
		if r.err == nil {
			passed++
			totalEndpoints += r.endpoints
		}
	}

	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Metric | Value |\n")
	sb.WriteString("|--------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Specs Tested | %d |\n", total))
	sb.WriteString(fmt.Sprintf("| Successful | %d |\n", passed))
	sb.WriteString(fmt.Sprintf("| Failed | %d |\n", total-passed))
	sb.WriteString(fmt.Sprintf("| Total Endpoints | %d |\n", totalEndpoints))
	if total > 0 {
		sb.WriteString(fmt.Sprintf("| Success Rate | %d%% |\n", passed*100/total))
	}
	sb.WriteString("\n")

	sb.WriteString("## Results\n\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("### %s\n\n", r.name))
		sb.WriteString(fmt.Sprintf("**Spec:** `%s`\n\n", r.specPath))
		sb.WriteString(fmt.Sprintf("**Duration:** %s\n\n", r.elapsed.Round(time.Millisecond)))
		if r.err != nil {
			sb.WriteString("**Status:** ❌ Failed\n\n")
			sb.WriteString(fmt.Sprintf("**Error:**\n```\n%v\n```\n\n", r.err))
		} else {
			sb.WriteString("**Status:** ✅ Success\n\n")
			if r.title != "" {
				sb.WriteString(fmt.Sprintf("**Title:** %s\n\n", r.title))
			}
			sb.WriteString(fmt.Sprintf("**Endpoints Found:** %d\n\n", r.endpoints))
		}
		sb.WriteString("---\n\n")
	}
	return sb.String()
}

// specName returns a short display name for a spec path.
func specName(path string) string {
	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]
	if idx := strings.LastIndex(name, "."); idx > 0 {
		name = name[:idx]
	}
	return name
}
