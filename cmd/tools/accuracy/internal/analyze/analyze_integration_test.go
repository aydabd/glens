package analyze_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"glens/tools/accuracy/internal/analyze"
)

// sampleSpecPath returns the absolute path to test_specs/sample_api.json
// by navigating from this test file up to the repository root.
func sampleSpecPath(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// file is at: cmd/tools/accuracy/internal/analyze/analyze_integration_test.go
	// repo root is 5 directories up
	root := filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "..")
	return filepath.Join(root, "test_specs", "sample_api.json")
}

func TestSpecs_sampleAPI(t *testing.T) {
	specPath := sampleSpecPath(t)

	results := analyze.Specs([]string{specPath})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]

	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if r.Title != "Sample API" {
		t.Errorf("title = %q, want %q", r.Title, "Sample API")
	}
	if r.Endpoints != 3 {
		t.Errorf("endpoints = %d, want 3", r.Endpoints)
	}
	if r.Name != "sample_api" {
		t.Errorf("name = %q, want %q", r.Name, "sample_api")
	}
	if r.Elapsed <= 0 {
		t.Error("elapsed duration should be positive")
	}
}

func TestSpecs_missingFile(t *testing.T) {
	results := analyze.Specs([]string{"/nonexistent/path/spec.json"})

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestSpecs_multipleSpecs(t *testing.T) {
	specPath := sampleSpecPath(t)

	// Run the same spec twice to verify multi-spec handling
	results := analyze.Specs([]string{specPath, specPath})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for i, r := range results {
		if r.Err != nil {
			t.Errorf("result[%d] unexpected error: %v", i, r.Err)
		}
		if r.Endpoints != 3 {
			t.Errorf("result[%d] endpoints = %d, want 3", i, r.Endpoints)
		}
	}
}
