// Package analyze loads and analyses OpenAPI specs for the accuracy tool.
package analyze

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Result holds the outcome of analysing a single spec.
type Result struct {
	Name      string
	SpecPath  string
	Title     string
	Endpoints int
	Elapsed   time.Duration
	Err       error
}

// minimalSpec holds only the fields needed for accuracy metrics.
type minimalSpec struct {
	Info struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`
	Paths map[string]map[string]interface{} `json:"paths"`
}

// Specs analyses each spec and returns a Result per spec.
func Specs(paths []string) []Result {
	results := make([]Result, 0, len(paths))
	for _, p := range paths {
		start := time.Now()
		spec, err := loadSpec(p)
		elapsed := time.Since(start)

		r := Result{
			Name:     specName(p),
			SpecPath: p,
			Elapsed:  elapsed,
			Err:      err,
		}
		if err == nil {
			r.Title = spec.Info.Title
			r.Endpoints = countEndpoints(spec)
		}
		results = append(results, r)
	}
	return results
}

func loadSpec(source string) (*minimalSpec, error) {
	data, err := fetch(source)
	if err != nil {
		return nil, err
	}
	var spec minimalSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &spec, nil
}

func fetch(source string) ([]byte, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Get(source) //nolint:gosec
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", err)
		}
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}
	return os.ReadFile(source) //nolint:gosec
}

func countEndpoints(spec *minimalSpec) int {
	count := 0
	for _, methods := range spec.Paths {
		count += len(methods)
	}
	return count
}

// specName derives a short display name from the file path.
func specName(path string) string {
	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]
	if idx := strings.LastIndex(name, "."); idx > 0 {
		name = name[:idx]
	}
	return name
}
