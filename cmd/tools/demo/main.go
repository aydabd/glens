// Command demo demonstrates glens OpenAPI parsing capabilities.
// It is a cross-platform replacement for scripts/demo_modern.sh.
package main

import (
	"flag"
	"fmt"
	"os"

	"glens/tools/demo/internal/loader"
	"glens/tools/demo/internal/render"
)

// version is set at build time via -ldflags="-X main.version=<tag>".
var version = "0.1.0"

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
	render.Banner()
	fmt.Printf("Parsing OpenAPI spec: %s\n\n", specPath)

	spec, err := loader.Load(specPath)
	if err != nil {
		return fmt.Errorf("failed to load spec: %w", err)
	}

	render.SpecInfo(spec)
	render.Endpoints(spec)
	render.ModelComparison()
	render.SampleTest()

	return nil
}
