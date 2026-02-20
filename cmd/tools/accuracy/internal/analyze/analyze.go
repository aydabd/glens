// Package analyze loads and analyses OpenAPI specs for the accuracy tool.
package analyze

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			bodySnippet, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			return nil, fmt.Errorf("HTTP %d %s: %s", resp.StatusCode, resp.Status, strings.TrimSpace(string(bodySnippet)))
		}
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

// specName derives a short display name from the file path, cross-platform.
func specName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
