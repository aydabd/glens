package safety

import "strings"

// Risk represents the risk level of an endpoint.
type Risk string

// Risk level constants.
const (
	RiskSafe   Risk = "safe"
	RiskMedium Risk = "medium"
	RiskHigh   Risk = "high"
)

// Category represents the operational category of an endpoint.
type Category string

// Category constants for endpoint operations.
const (
	CategoryRead    Category = "read"
	CategoryWrite   Category = "write"
	CategoryMutate  Category = "mutate"
	CategoryDestroy Category = "destroy"
)

// EndpointCategory holds the categorisation result for a single endpoint.
type EndpointCategory struct {
	Path     string   `json:"path"`
	Method   string   `json:"method"`
	Category Category `json:"category"`
	Risk     Risk     `json:"risk"`
}

// EndpointInput is the input for batch categorisation.
type EndpointInput struct {
	Method string
	Path   string
	XSafe  bool
}

// safePostSuffixes are path segments that indicate a POST is read-only.
var safePostSuffixes = []string{
	"/search", "/query", "/list", "/find", "/check", "/validate", "/verify",
}

// Categorise returns the category and risk for a single endpoint.
func Categorise(method, path string, xSafe bool) EndpointCategory {
	ec := EndpointCategory{
		Path:   path,
		Method: strings.ToUpper(method),
	}

	if xSafe {
		ec.Category = CategoryRead
		ec.Risk = RiskSafe
		return ec
	}

	switch ec.Method {
	case "GET", "HEAD", "OPTIONS":
		ec.Category = CategoryRead
		ec.Risk = RiskSafe
	case "POST":
		if isSafePost(path) {
			ec.Category = CategoryRead
			ec.Risk = RiskSafe
		} else {
			ec.Category = CategoryWrite
			ec.Risk = RiskMedium
		}
	case "PUT", "PATCH":
		ec.Category = CategoryMutate
		ec.Risk = RiskMedium
	case "DELETE":
		ec.Category = CategoryDestroy
		ec.Risk = RiskHigh
	default:
		ec.Category = CategoryWrite
		ec.Risk = RiskMedium
	}

	return ec
}

// CategoriseAll categorises a batch of endpoints.
func CategoriseAll(endpoints []EndpointInput) []EndpointCategory {
	results := make([]EndpointCategory, len(endpoints))
	for i, ep := range endpoints {
		results[i] = Categorise(ep.Method, ep.Path, ep.XSafe)
	}
	return results
}

// Warnings returns human-readable warnings for medium and high risk endpoints.
func Warnings(categories []EndpointCategory) []string {
	var warnings []string
	for _, c := range categories {
		switch c.Risk {
		case RiskMedium:
			warnings = append(warnings, c.Method+" "+c.Path+" is "+string(c.Category)+" (medium risk)")
		case RiskHigh:
			warnings = append(warnings, c.Method+" "+c.Path+" is "+string(c.Category)+" (high risk)")
		}
	}
	return warnings
}

func isSafePost(path string) bool {
	lower := strings.ToLower(path)
	for _, suffix := range safePostSuffixes {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}
