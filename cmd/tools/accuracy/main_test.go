package main

import (
	"strings"
	"testing"
)

func TestSpecName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"test_specs/sample_api.json", "sample_api"},
		{"api.yaml", "api"},
		{"dir/sub/my_spec.json", "my_spec"},
		{"noext", "noext"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := specName(tt.input)
			if got != tt.want {
				t.Errorf("specName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuildReport_empty(t *testing.T) {
	report := buildReport(nil)
	if !strings.Contains(report, "# Glens Accuracy Report") {
		t.Error("report missing title")
	}
}

func TestBuildReport_withResults(t *testing.T) {
	results := []result{
		{name: "sample_api", specPath: "test.json", endpoints: 3},
		{name: "bad_spec", specPath: "bad.json", err: errMock("parse error")},
	}
	report := buildReport(results)
	if !strings.Contains(report, "sample_api") {
		t.Error("report missing spec name")
	}
	if !strings.Contains(report, "✅ Success") {
		t.Error("report missing success marker")
	}
	if !strings.Contains(report, "❌ Failed") {
		t.Error("report missing failure marker")
	}
}

type errMock string

func (e errMock) Error() string { return string(e) }
