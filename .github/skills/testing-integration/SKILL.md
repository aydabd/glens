---
name: testing-integration
description: >
  Integration and end-to-end testing design for this project. Use this when writing
  integration tests (*_integration_test.go), E2E binary tests, or when deciding
  what belongs in integration vs unit tests.
---

# Skill: Expert Test Automation — Integration & E2E Testing

## Integration Tests in This Codebase

Integration tests use `test_specs/sample_api.json` — **no network dependency**.
File name convention: `*_integration_test.go`.

```go
// Arrange: load real spec files from test_specs/ or testdata/
// Act:     run the full pipeline (parse → analyze → report) with real components
// Assert:  verify combined output and side-effects, not internal state
```

## What Belongs in Integration Tests

- Two or more **real** components working together (no mocks for the glue between them).
- The full parse → analyze → report pipeline with real files.
- Error propagation across component boundaries.
- File I/O, CLI flag parsing, or configuration loading.

**Do NOT add integration tests for:**

- Single-function logic — use a unit test instead.
- Flaky external network calls — stub or mock them.

## End-to-End Test Design

An E2E test exercises the compiled binary exactly as a user would:

```go
func TestCLI_AnalyzeSpec_PrintsReport(t *testing.T) {
    binary := buildBinary(t)          // compile once; reuse across table rows
    specPath := sampleSpecPath(t)

    out, err := exec.Command(binary, "analyze", "--spec", specPath).CombinedOutput()

    assert.NoError(t, err)
    assert.Contains(t, string(out), "Sample API")
}
```

## Corner Cases for Integration & E2E

Before declaring a feature integration-tested, verify:

- **Empty spec** — zero endpoints, zero schemas.
- **Malformed JSON** — truncated file, wrong value types.
- **Duplicate paths** — same path + method defined twice.
- **Large spec** — hundreds of endpoints with deep nesting.
- **Permission error** — unreadable or missing file path.
- **Concurrent execution** — two goroutines/workers on the same spec.

## Locating Spec Files in Tests

Use `runtime.Caller` to locate `test_specs/` without hardcoding absolute paths:

```go
func sampleSpecPath(t *testing.T) string {
    t.Helper()
    _, file, _, ok := runtime.Caller(0)
    if !ok {
        t.Fatal("runtime.Caller failed")
    }
    // Navigate from this file up to the repository root, then into test_specs/.
    root := filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "..")
    return filepath.Join(root, "test_specs", "sample_api.json")
}
```

Adjust the number of `".."` segments to match the depth of the test file.

## Compiling the Binary for E2E Tests

Compile the binary once per test run and cache the path:

```go
func buildBinary(t *testing.T) string {
    t.Helper()
    bin := filepath.Join(t.TempDir(), "tool")
    cmd := exec.Command("go", "build", "-o", bin, ".")
    out, err := cmd.CombinedOutput()
    require.NoError(t, err, "build failed: %s", out)
    return bin
}
```
