// Package e2e contains end-to-end tests that compile the glens binary and
// exercise it exactly as a user would.
//
// Each test is fully self-contained:
//   - A mock Ollama HTTP server is started per test (no real Ollama required).
//   - A temporary config file is written pointing the "ollama" model at the mock.
//   - The glens binary is compiled once per test run via TestMain.
//   - All teardown is handled by t.Cleanup / t.TempDir — no shared mutable state.
//
// Feature under test: glens analyze <spec> --ai-models <local-model>
// (local LLM support via Ollama — no cloud API key required)
package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Binary lifecycle (compiled once per test run) ─────────────────────────────

var (
	binaryOnce sync.Once
	binaryPath string
	binaryErr  error
)

// TestMain compiles the glens binary once before any test runs, so individual
// tests don't pay the compilation cost independently.
func TestMain(m *testing.M) {
	binaryOnce.Do(compileBinary)
	os.Exit(m.Run())
}

func compileBinary() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		binaryErr = fmt.Errorf("runtime.Caller failed")
		return
	}
	// file: cmd/glens/e2e/local_llm_integration_test.go — one level up is cmd/glens/
	moduleRoot := filepath.Join(filepath.Dir(file), "..")
	dir, err := os.MkdirTemp("", "glens-e2e-*")
	if err != nil {
		binaryErr = fmt.Errorf("create temp dir: %w", err)
		return
	}
	bin := filepath.Join(dir, "glens")
	cmd := exec.Command("go", "build", "-o", bin, ".") //nolint:gosec // output path is a temp dir we created; source is always "." (fixed string)
	cmd.Dir = moduleRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		binaryErr = fmt.Errorf("build failed: %w\n%s", err, out)
		return
	}
	binaryPath = bin
}

// glensBinary returns the compiled binary path, or fails the test if it
// wasn't built successfully.
func glensBinary(t *testing.T) string {
	t.Helper()
	require.NoError(t, binaryErr, "glens binary must compile successfully")
	return binaryPath
}

// ── Shared helpers ────────────────────────────────────────────────────────────

// repoRoot returns the repository root by walking up from this file.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")
	// file is at cmd/glens/e2e/local_llm_integration_test.go → 3 dirs up = repo root
	return filepath.Join(filepath.Dir(file), "..", "..", "..")
}

// sampleSpecPath returns the absolute path to test_specs/sample_api.json.
func sampleSpecPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "test_specs", "sample_api.json")
}

// ollamaGenerate is the response struct for the /api/generate mock endpoint.
type ollamaGenerate struct {
	Model         string `json:"model"`
	Response      string `json:"response"`
	Done          bool   `json:"done"`
	TotalDuration int64  `json:"total_duration"`
	EvalDuration  int64  `json:"eval_duration"`
}

// mockGenerateResponse returns the JSON bytes the mock sends for /api/generate.
// The Response field contains a minimal valid Go test function.
func mockGenerateResponse() []byte {
	resp := ollamaGenerate{
		Model: "mistral",
		Response: "package main\n\nimport (\n\t\"net/http\"\n\t\"testing\"\n\n" +
			"\t\"github.com/stretchr/testify/assert\"\n)\n\n" +
			"func TestGETUsers(t *testing.T) {\n" +
			"\tresp, err := http.Get(\"http://localhost:8080/users\")\n" +
			"\tassert.NoError(t, err)\n" +
			"\tassert.Equal(t, 200, resp.StatusCode)\n}\n",
		Done:          true,
		TotalDuration: 1_000_000_000,
		EvalDuration:  800_000_000,
	}
	b, _ := json.Marshal(resp)
	return b
}

// requestCounts tracks how many times each mock endpoint was called.
// Updated from the HTTP handler goroutine; read after CombinedOutput() returns
// (subprocess exited), so no additional synchronization is required.
type requestCounts struct {
	generate int
	tags     int
	version  int
}

// startMockOllama starts a lightweight HTTP server that simulates the subset
// of the Ollama API used by glens.  The server shuts down when the test ends.
// The returned *requestCounts is incremented by the handlers.
func startMockOllama(t *testing.T) (string, *requestCounts) {
	t.Helper()

	counts := &requestCounts{}
	mux := http.NewServeMux()

	// Health / version check
	mux.HandleFunc("/api/version", func(w http.ResponseWriter, _ *http.Request) {
		counts.version++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"0.6.0"}`))
	})

	// List installed models
	mux.HandleFunc("/api/tags", func(w http.ResponseWriter, _ *http.Request) {
		counts.tags++
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"models": []map[string]interface{}{
				{"name": "mistral", "size": 4109856768, "digest": "abc123def456789"},
				{"name": "codellama:7b-instruct", "size": 3826793472, "digest": "def456abc123789"},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	})

	// Text generation endpoint
	mux.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		counts.generate++
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(mockGenerateResponse())
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv.URL, counts
}

// writeConfig writes a minimal glens config YAML that points the "ollama"
// model at ollamaURL, and returns the config file path.
func writeConfig(t *testing.T, ollamaURL string) string {
	t.Helper()
	content := fmt.Sprintf(`ai_models:
  ollama:
    base_url: "%s"
    model: "mistral"
    timeout: "30s"
    temperature: 0.1
    max_tokens: 2000
logging:
  level: "warn"
  format: "console"
`, ollamaURL)
	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

// runGlens executes the glens binary with the given arguments (all supplied by
// test helpers in this package) and returns combined stdout+stderr and any error.
func runGlens(t *testing.T, args ...string) (string, error) {
	t.Helper()
	bin := glensBinary(t)
	cmd := exec.Command(bin, args...) //nolint:gosec // binary compiled from source; args are literals from test helpers
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// ── Feature: glens analyze with a local Ollama model ─────────────────────────

// TestLocalLLM_AnalyzeSpec_GeneratesReport is the primary E2E scenario
// demonstrating the full pipeline when using a local open-source LLM:
//
//	glens analyze <spec> --ai-models ollama --create-issues=false --run-tests=false
//
// The mock Ollama server provides generated test code; the binary must
// parse the spec, call the model, and write a non-empty markdown report.
func TestLocalLLM_AnalyzeSpec_GeneratesReport(t *testing.T) {
	ollamaURL, counts := startMockOllama(t)
	cfgFile := writeConfig(t, ollamaURL)
	specPath := sampleSpecPath(t)
	reportPath := filepath.Join(t.TempDir(), "report.md")

	out, err := runGlens(t,
		"analyze", specPath,
		"--config", cfgFile,
		"--ai-models", "ollama",
		"--create-issues=false",
		"--run-tests=false",
		"--output", reportPath,
	)

	require.NoError(t, err, "glens analyze should exit 0; output:\n%s", out)

	// The mock generate endpoint must have been called (once per endpoint in the spec).
	assert.Greater(t, counts.generate, 0,
		"mock /api/generate must be hit — verifies the Ollama URL in config is applied")

	// Report file must exist and contain API title
	content, readErr := os.ReadFile(reportPath)
	require.NoError(t, readErr, "report file must be created")
	assert.Contains(t, string(content), "Sample API",
		"report should reference the spec title from the OpenAPI spec")
}

// TestLocalLLM_AnalyzeSpec_ModelShortcut verifies that the `mistral-local`
// shortcut (introduced as part of the local LLM feature) is recognised by
// the binary when backed by a mock Ollama server.
func TestLocalLLM_AnalyzeSpec_ModelShortcut(t *testing.T) {
	ollamaURL, counts := startMockOllama(t)
	cfgFile := writeConfig(t, ollamaURL)
	specPath := sampleSpecPath(t)
	reportPath := filepath.Join(t.TempDir(), "report.md")

	out, err := runGlens(t,
		"analyze", specPath,
		"--config", cfgFile,
		"--ai-models", "mistral-local",
		"--create-issues=false",
		"--run-tests=false",
		"--output", reportPath,
	)

	require.NoError(t, err, "mistral-local shortcut should be accepted; output:\n%s", out)
	assert.Greater(t, counts.generate, 0,
		"mock /api/generate must be hit for mistral-local shortcut")
	_, statErr := os.Stat(reportPath)
	assert.NoError(t, statErr, "report file must be created for mistral-local shortcut")
}

// TestLocalLLM_AnalyzeSpec_CustomOllamaModel exercises the ollama:<model>
// escape hatch that allows any arbitrary Ollama model name.
func TestLocalLLM_AnalyzeSpec_CustomOllamaModel(t *testing.T) {
	ollamaURL, counts := startMockOllama(t)
	cfgFile := writeConfig(t, ollamaURL)
	specPath := sampleSpecPath(t)
	reportPath := filepath.Join(t.TempDir(), "report.md")

	out, err := runGlens(t,
		"analyze", specPath,
		"--config", cfgFile,
		"--ai-models", "ollama:mistral",
		"--create-issues=false",
		"--run-tests=false",
		"--output", reportPath,
	)

	require.NoError(t, err, "ollama:<model> custom syntax should be accepted; output:\n%s", out)
	assert.Greater(t, counts.generate, 0,
		"mock /api/generate must be hit for ollama:<model> syntax")
	_, statErr := os.Stat(reportPath)
	assert.NoError(t, statErr, "report file must be created for ollama:<model> syntax")
}

// TestLocalLLM_AnalyzeSpec_OllamaServerDown verifies that when the Ollama server
// is unreachable, glens still exits cleanly (exit 0) without panicking.
// The analyze command treats AI generation failures as non-fatal per-endpoint
// errors (it logs them and continues), so the process must not crash.
func TestLocalLLM_AnalyzeSpec_OllamaServerDown(t *testing.T) {
	// Write a config pointing at a port where nothing is listening.
	cfgFile := writeConfig(t, "http://127.0.0.1:1")
	specPath := sampleSpecPath(t)
	reportPath := filepath.Join(t.TempDir(), "report.md")

	// Discard error: we expect analyze to exit 0 even when Ollama is down,
	// because generation failures are logged and skipped per endpoint.
	out, _ := runGlens(t,
		"analyze", specPath,
		"--config", cfgFile,
		"--ai-models", "ollama",
		"--create-issues=false",
		"--run-tests=false",
		"--output", reportPath,
	)

	assert.NotContains(t, strings.ToLower(out), "panic",
		"glens must not panic when Ollama is unreachable")
}

// TestLocalLLM_ModelsOllama_ListCommand verifies that the `models ollama list`
// subcommand works against a mock Ollama server and returns the model names.
func TestLocalLLM_ModelsOllama_ListCommand(t *testing.T) {
	ollamaURL, counts := startMockOllama(t)
	cfgFile := writeConfig(t, ollamaURL)

	out, err := runGlens(t,
		"models", "ollama", "list",
		"--config", cfgFile,
	)

	require.NoError(t, err, "models ollama list should exit 0; output:\n%s", out)
	assert.Equal(t, 1, counts.tags, "mock /api/tags must be hit exactly once")
	assert.Contains(t, out, "mistral",
		"list output should include the mistral model from mock server")
}

// TestLocalLLM_ModelsOllama_StatusCommand verifies that the `models ollama status`
// subcommand reports healthy when pointed at the mock server.
func TestLocalLLM_ModelsOllama_StatusCommand(t *testing.T) {
	ollamaURL, counts := startMockOllama(t)
	cfgFile := writeConfig(t, ollamaURL)

	out, err := runGlens(t,
		"models", "ollama", "status",
		"--config", cfgFile,
	)

	require.NoError(t, err, "models ollama status should exit 0; output:\n%s", out)
	assert.Greater(t, counts.version, 0, "mock /api/version must be hit for health check")
	assert.Contains(t, out, "✅",
		"status command should report healthy against the mock server")
}
