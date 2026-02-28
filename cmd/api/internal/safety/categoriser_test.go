package safety

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategorise(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		xSafe    bool
		wantCat  Category
		wantRisk Risk
	}{
		// Read-safe methods
		{"GET is read/safe", "GET", "/users", false, CategoryRead, RiskSafe},
		{"HEAD is read/safe", "HEAD", "/users", false, CategoryRead, RiskSafe},
		{"OPTIONS is read/safe", "OPTIONS", "/users", false, CategoryRead, RiskSafe},
		{"GET lowercase normalised", "get", "/items", false, CategoryRead, RiskSafe},

		// POST — default write/medium
		{"POST default is write/medium", "POST", "/users", false, CategoryWrite, RiskMedium},

		// POST — safe paths
		{"POST /search is read/safe", "POST", "/users/search", false, CategoryRead, RiskSafe},
		{"POST /query is read/safe", "POST", "/data/query", false, CategoryRead, RiskSafe},
		{"POST /list is read/safe", "POST", "/items/list", false, CategoryRead, RiskSafe},
		{"POST /find is read/safe", "POST", "/records/find", false, CategoryRead, RiskSafe},
		{"POST /check is read/safe", "POST", "/health/check", false, CategoryRead, RiskSafe},
		{"POST /validate is read/safe", "POST", "/schema/validate", false, CategoryRead, RiskSafe},
		{"POST /verify is read/safe", "POST", "/token/verify", false, CategoryRead, RiskSafe},

		// PUT, PATCH
		{"PUT is mutate/medium", "PUT", "/users/1", false, CategoryMutate, RiskMedium},
		{"PATCH is mutate/medium", "PATCH", "/users/1", false, CategoryMutate, RiskMedium},

		// DELETE
		{"DELETE is destroy/high", "DELETE", "/users/1", false, CategoryDestroy, RiskHigh},

		// x-safe override
		{"x-safe overrides DELETE to read/safe", "DELETE", "/users/1", true, CategoryRead, RiskSafe},
		{"x-safe overrides POST to read/safe", "POST", "/users", true, CategoryRead, RiskSafe},

		// Unknown method
		{"unknown method is write/medium", "TRACE", "/debug", false, CategoryWrite, RiskMedium},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Categorise(tt.method, tt.path, tt.xSafe)
			assert.Equal(t, tt.wantCat, got.Category)
			assert.Equal(t, tt.wantRisk, got.Risk)
			assert.Equal(t, tt.path, got.Path)
		})
	}
}

func TestCategoriseAll(t *testing.T) {
	inputs := []EndpointInput{
		{Method: "GET", Path: "/users"},
		{Method: "POST", Path: "/users"},
		{Method: "DELETE", Path: "/users/1"},
		{Method: "POST", Path: "/users/search", XSafe: false},
		{Method: "PUT", Path: "/users/1", XSafe: true},
	}

	results := CategoriseAll(inputs)

	assert.Len(t, results, 5)
	assert.Equal(t, CategoryRead, results[0].Category)
	assert.Equal(t, RiskSafe, results[0].Risk)

	assert.Equal(t, CategoryWrite, results[1].Category)
	assert.Equal(t, RiskMedium, results[1].Risk)

	assert.Equal(t, CategoryDestroy, results[2].Category)
	assert.Equal(t, RiskHigh, results[2].Risk)

	assert.Equal(t, CategoryRead, results[3].Category)
	assert.Equal(t, RiskSafe, results[3].Risk)

	// x-safe override on PUT
	assert.Equal(t, CategoryRead, results[4].Category)
	assert.Equal(t, RiskSafe, results[4].Risk)
}

func TestWarnings(t *testing.T) {
	categories := []EndpointCategory{
		{Path: "/users", Method: "GET", Category: CategoryRead, Risk: RiskSafe},
		{Path: "/users", Method: "POST", Category: CategoryWrite, Risk: RiskMedium},
		{Path: "/users/1", Method: "DELETE", Category: CategoryDestroy, Risk: RiskHigh},
		{Path: "/users/1", Method: "PUT", Category: CategoryMutate, Risk: RiskMedium},
	}

	warnings := Warnings(categories)

	assert.Len(t, warnings, 3)
	assert.Contains(t, warnings[0], "POST /users")
	assert.Contains(t, warnings[0], "medium risk")
	assert.Contains(t, warnings[1], "DELETE /users/1")
	assert.Contains(t, warnings[1], "high risk")
	assert.Contains(t, warnings[2], "PUT /users/1")
	assert.Contains(t, warnings[2], "medium risk")
}

func TestWarnings_empty(t *testing.T) {
	categories := []EndpointCategory{
		{Path: "/users", Method: "GET", Category: CategoryRead, Risk: RiskSafe},
	}

	warnings := Warnings(categories)

	assert.Empty(t, warnings)
}
