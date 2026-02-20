package parser

import (
	"time"
)

// OpenAPISpec represents the parsed OpenAPI specification
type OpenAPISpec struct {
	Info      Info       `json:"info"`
	Servers   []Server   `json:"servers"`
	Endpoints []Endpoint `json:"endpoints"`
	Version   string     `json:"version"`
	ParsedAt  time.Time  `json:"parsed_at"`
}

// Info contains API metadata
type Info struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Contact     *Contact `json:"contact,omitempty"`
}

// Contact information for the API
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// Server represents an API server
type Server struct {
	URL         string            `json:"url"`
	Description string            `json:"description,omitempty"`
	Variables   map[string]string `json:"variables,omitempty"`
}

// Endpoint represents a single API endpoint
type Endpoint struct {
	ID          string                `json:"id"`
	Method      string                `json:"method"`
	Path        string                `json:"path"`
	OperationID string                `json:"operation_id,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"request_body,omitempty"`
	Responses   map[string]Response   `json:"responses,omitempty"`
	Security    []SecurityRequirement `json:"security,omitempty"`
}

// Parameter represents an endpoint parameter
type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"` // query, header, path, cookie
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required"`
	Schema      Schema      `json:"schema"`
	Example     interface{} `json:"example,omitempty"`
}

// RequestBody represents the request body
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required"`
	Content     map[string]MediaType `json:"content"`
}

// Response represents an API response
type Response struct {
	Description string               `json:"description"`
	Headers     map[string]Header    `json:"headers,omitempty"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// MediaType represents a media type specification
type MediaType struct {
	Schema   Schema             `json:"schema,omitempty"`
	Example  interface{}        `json:"example,omitempty"`
	Examples map[string]Example `json:"examples,omitempty"`
}

// Schema represents a JSON schema
type Schema struct {
	Type        string            `json:"type,omitempty"`
	Format      string            `json:"format,omitempty"`
	Description string            `json:"description,omitempty"`
	Properties  map[string]Schema `json:"properties,omitempty"`
	Items       *Schema           `json:"items,omitempty"`
	Required    []string          `json:"required,omitempty"`
	Enum        []interface{}     `json:"enum,omitempty"`
	Example     interface{}       `json:"example,omitempty"`
	Minimum     *float64          `json:"minimum,omitempty"`
	Maximum     *float64          `json:"maximum,omitempty"`
	MinLength   *int              `json:"min_length,omitempty"`
	MaxLength   *int              `json:"max_length,omitempty"`
	Pattern     string            `json:"pattern,omitempty"`
	Ref         string            `json:"$ref,omitempty"`
}

// Header represents a response header
type Header struct {
	Description string      `json:"description,omitempty"`
	Schema      Schema      `json:"schema,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}

// Example represents an example value
type Example struct {
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Value       interface{} `json:"value,omitempty"`
}

// SecurityRequirement represents security requirements
type SecurityRequirement map[string][]string
