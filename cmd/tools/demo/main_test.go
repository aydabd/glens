package main

import (
	"testing"

	"glens/tools/demo/internal/render"
)

func TestRunDemo_missingFile(t *testing.T) {
	err := runDemo("/nonexistent/spec.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRenderBanner(_ *testing.T) {
	render.Banner() // must not panic
}

func TestRenderModelComparison(_ *testing.T) {
	render.ModelComparison() // must not panic
}

func TestRenderSampleTest(_ *testing.T) {
	render.SampleTest() // must not panic
}
