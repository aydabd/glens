---
name: testing-patterns
description: >
  Go-specific test patterns for this project. Use this when writing Go tests —
  table-driven tests, interface mocking, test helpers, testify assertions,
  and parallel sub-tests.
---

# Skill: Expert Test Automation — Go Patterns

## Table-Driven Tests

Always prefer table-driven tests when testing the same function across multiple scenarios:

```go
func TestValidateEndpoint(t *testing.T) {
    tests := []struct {
        name    string
        input   Endpoint
        wantErr bool
        errMsg  string
    }{
        {
            name:  "valid endpoint",
            input: Endpoint{Path: "/users", Method: "GET"},
        },
        {
            name:    "missing path",
            input:   Endpoint{Method: "GET"},
            wantErr: true,
            errMsg:  "path is required",
        },
        {
            name:    "unsupported HTTP method",
            input:   Endpoint{Path: "/users", Method: "PATCH"},
            wantErr: true,
            errMsg:  "unsupported method",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEndpoint(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
                return
            }
            assert.NoError(t, err)
        })
    }
}
```

## Mocking External Dependencies

Mock at the **interface** boundary, not the implementation:

```go
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

type mockHTTPClient struct {
    response *http.Response
    err      error
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
    return m.response, m.err
}
```

Keep mocks in the same package as the test that uses them.

## Test Helpers

Extract repetitive setup into helpers tagged with `t.Helper()`:

```go
func newTestEndpoint(t *testing.T, path, method string) *Endpoint {
    t.Helper()
    return &Endpoint{Path: path, Method: method}
}
```

## Assertions

Use `testify/assert` for readable failure messages:

- `assert.Equal(t, want, got)` — value equality
- `assert.NoError(t, err)` — clean success path
- `assert.Error(t, err)` — expected failure
- `assert.Contains(t, s, substr)` — partial string match
- `require.NoError(t, err)` — stop the test immediately on failure

Use `require` (not `assert`) when subsequent assertions are meaningless after a failure.

## Parallel Tests

Mark independent sub-tests parallel to shorten the test run with `t.Parallel()`.
Do **not** use it when tests share mutable state or write to the same file.