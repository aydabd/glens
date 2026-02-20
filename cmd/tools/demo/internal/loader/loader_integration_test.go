package loader_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"glens/tools/demo/internal/loader"
)

// sampleSpecPath returns the absolute path to test_specs/sample_api.json
// by navigating from this test file up to the repository root.
func sampleSpecPath(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// file is at: cmd/tools/demo/internal/loader/loader_integration_test.go
	// repo root is 5 directories up
	root := filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "..")
	return filepath.Join(root, "test_specs", "sample_api.json")
}

func TestLoad_sampleAPI(t *testing.T) {
	specPath := sampleSpecPath(t)

	spec, err := loader.Load(specPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if spec.Info.Title != "Sample API" {
		t.Errorf("title = %q, want %q", spec.Info.Title, "Sample API")
	}
	if spec.Info.Version != "1.0.0" {
		t.Errorf("version = %q, want %q", spec.Info.Version, "1.0.0")
	}
	if len(spec.Servers) == 0 {
		t.Error("expected at least one server")
	} else if spec.Servers[0].URL != "https://api.example.com/v1" {
		t.Errorf("server URL = %q, want %q", spec.Servers[0].URL, "https://api.example.com/v1")
	}
	if len(spec.Paths) != 3 {
		t.Errorf("paths count = %d, want 3", len(spec.Paths))
	}
}

func TestLoad_missingFile(t *testing.T) {
	_, err := loader.Load("/nonexistent/path/spec.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_expectedPaths(t *testing.T) {
	specPath := sampleSpecPath(t)

	spec, err := loader.Load(specPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantPaths := map[string]string{
		"/users":      "get",
		"/users/{id}": "get",
		"/posts":      "post",
	}
	for path, method := range wantPaths {
		methods, ok := spec.Paths[path]
		if !ok {
			t.Errorf("missing expected path %q", path)
			continue
		}
		if _, hasMethod := methods[method]; !hasMethod {
			t.Errorf("path %q missing expected method %q", path, method)
		}
	}
}
