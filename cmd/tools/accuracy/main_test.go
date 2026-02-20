package main

import (
	"fmt"
	"strings"
	"testing"

	"glens/tools/accuracy/internal/analyze"
	"glens/tools/accuracy/internal/report"
)

func TestReport_empty(t *testing.T) {
	out := report.Build(nil)
	if !strings.Contains(out, "# Glens Accuracy Report") {
		t.Error("report missing title")
	}
}

func TestReport_withResults(t *testing.T) {
	results := []analyze.Result{
		{Name: "sample_api", SpecPath: "test.json", Endpoints: 3},
		{Name: "bad_spec", SpecPath: "bad.json", Err: fmt.Errorf("parse error")}, //nolint:err113
	}
	out := report.Build(results)
	if !strings.Contains(out, "sample_api") {
		t.Error("report missing spec name")
	}
	if !strings.Contains(out, "✅ Success") {
		t.Error("report missing success marker")
	}
	if !strings.Contains(out, "❌ Failed") {
		t.Error("report missing failure marker")
	}
}
