// Command accuracy evaluates OpenAPI spec parsing accuracy.
// It is a cross-platform replacement for scripts/test_accuracy.sh.
package main

import (
	"flag"
	"fmt"
	"os"

	"glens/tools/accuracy/internal/analyze"
	"glens/tools/accuracy/internal/report"
)

// version is set at build time via -ldflags="-X main.version=<tag>".
var version = "0.1.0"

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

	results := analyze.Specs(specs)
	output := report.Build(results)

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(output), 0o600); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing report: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Report written to %s\n", outputFile)
	} else {
		fmt.Print(output)
	}

	for _, r := range results {
		if r.Err != nil {
			os.Exit(1)
		}
	}
}
