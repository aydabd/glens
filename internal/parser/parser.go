package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// ParseOpenAPISpec parses an OpenAPI specification from a URL or file path
func ParseOpenAPISpec(source string) (*OpenAPISpec, error) {
	log.Debug().Str("source", source).Msg("Parsing OpenAPI specification")

	var data []byte
	var err error

	if isURL(source) {
		data, err = fetchFromURL(source)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch from URL: %w", err)
		}
	} else {
		data, err = os.ReadFile(source)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
	}

	// Determine format based on content or extension
	var rawSpec map[string]interface{}
	if isYAML(source, data) {
		if err := yaml.Unmarshal(data, &rawSpec); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &rawSpec); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	}

	spec, err := convertToSpec(rawSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to internal format: %w", err)
	}

	spec.ParsedAt = time.Now()

	log.Info().
		Int("endpoints_count", len(spec.Endpoints)).
		Str("version", spec.Version).
		Str("title", spec.Info.Title).
		Msg("OpenAPI specification parsed successfully")

	return spec, nil
}

// isURL checks if the source is a URL
func isURL(source string) bool {
	u, err := url.Parse(source)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

// fetchFromURL fetches content from a URL
func fetchFromURL(urlStr string) ([]byte, error) {
	// Validate URL to mitigate G107 security warning
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Only allow HTTP and HTTPS schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Debug().Err(closeErr).Msg("failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// isYAML determines if the content is YAML based on file extension or content
func isYAML(source string, data []byte) bool {
	// Check file extension
	if strings.HasSuffix(strings.ToLower(source), ".yaml") ||
		strings.HasSuffix(strings.ToLower(source), ".yml") {
		return true
	}

	// Check content - YAML typically starts with "openapi:" or "swagger:"
	content := strings.TrimSpace(string(data))
	return strings.HasPrefix(content, "openapi:") ||
		strings.HasPrefix(content, "swagger:")
}

// convertToSpec converts the raw specification to our internal format
func convertToSpec(rawSpec map[string]interface{}) (*OpenAPISpec, error) {
	spec := &OpenAPISpec{
		Endpoints: []Endpoint{},
	}

	// Extract version
	if openapi, ok := rawSpec["openapi"].(string); ok {
		spec.Version = openapi
	} else if swagger, ok := rawSpec["swagger"].(string); ok {
		spec.Version = swagger
	}

	// Extract info
	if infoRaw, ok := rawSpec["info"].(map[string]interface{}); ok {
		spec.Info = extractInfo(infoRaw)
	}

	// Extract servers
	if serversRaw, ok := rawSpec["servers"].([]interface{}); ok {
		spec.Servers = extractServers(serversRaw)
	}

	// Extract paths and convert to endpoints
	if pathsRaw, ok := rawSpec["paths"].(map[string]interface{}); ok {
		endpoints, err := extractEndpoints(pathsRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to extract endpoints: %w", err)
		}
		spec.Endpoints = endpoints
	}

	return spec, nil
}

// extractInfo extracts API info
func extractInfo(infoRaw map[string]interface{}) Info {
	info := Info{}

	if title, ok := infoRaw["title"].(string); ok {
		info.Title = title
	}
	if description, ok := infoRaw["description"].(string); ok {
		info.Description = description
	}
	if version, ok := infoRaw["version"].(string); ok {
		info.Version = version
	}

	if contactRaw, ok := infoRaw["contact"].(map[string]interface{}); ok {
		contact := &Contact{}
		if name, ok := contactRaw["name"].(string); ok {
			contact.Name = name
		}
		if url, ok := contactRaw["url"].(string); ok {
			contact.URL = url
		}
		if email, ok := contactRaw["email"].(string); ok {
			contact.Email = email
		}
		info.Contact = contact
	}

	return info
}

// extractServers extracts server information
func extractServers(serversRaw []interface{}) []Server {
	var servers []Server

	for _, serverRaw := range serversRaw {
		if serverMap, ok := serverRaw.(map[string]interface{}); ok {
			server := Server{}
			if url, ok := serverMap["url"].(string); ok {
				server.URL = url
			}
			if description, ok := serverMap["description"].(string); ok {
				server.Description = description
			}
			servers = append(servers, server)
		}
	}

	return servers
}

// extractEndpoints extracts endpoints from paths
func extractEndpoints(pathsRaw map[string]interface{}) ([]Endpoint, error) {
	var endpoints []Endpoint

	for path, pathItemRaw := range pathsRaw {
		if pathItem, ok := pathItemRaw.(map[string]interface{}); ok {
			for method, operationRaw := range pathItem {
				if method == "parameters" || method == "servers" {
					continue // Skip path-level parameters and servers
				}

				if operation, ok := operationRaw.(map[string]interface{}); ok {
					endpoint := Endpoint{
						ID:        fmt.Sprintf("%s_%s", strings.ToUpper(method), strings.ReplaceAll(path, "/", "_")),
						Method:    strings.ToUpper(method),
						Path:      path,
						Responses: make(map[string]Response),
					}

					// Extract operation details
					if operationID, ok := operation["operationId"].(string); ok {
						endpoint.OperationID = operationID
					}
					if summary, ok := operation["summary"].(string); ok {
						endpoint.Summary = summary
					}
					if description, ok := operation["description"].(string); ok {
						endpoint.Description = description
					}

					// Extract tags
					if tagsRaw, ok := operation["tags"].([]interface{}); ok {
						for _, tagRaw := range tagsRaw {
							if tag, ok := tagRaw.(string); ok {
								endpoint.Tags = append(endpoint.Tags, tag)
							}
						}
					}

					// Extract parameters
					if parametersRaw, ok := operation["parameters"].([]interface{}); ok {
						endpoint.Parameters = extractParameters(parametersRaw)
					}

					// Extract request body
					if requestBodyRaw, ok := operation["requestBody"].(map[string]interface{}); ok {
						endpoint.RequestBody = extractRequestBody(requestBodyRaw)
					}

					// Extract responses
					if responsesRaw, ok := operation["responses"].(map[string]interface{}); ok {
						endpoint.Responses = extractResponses(responsesRaw)
					}

					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	return endpoints, nil
}

// extractParameters extracts parameters from operation
func extractParameters(parametersRaw []interface{}) []Parameter {
	var parameters []Parameter

	for _, paramRaw := range parametersRaw {
		param, ok := paramRaw.(map[string]interface{})
		if !ok {
			continue
		}

		parameter := Parameter{}

		if name, ok := param["name"].(string); ok {
			parameter.Name = name
		}
		if in, ok := param["in"].(string); ok {
			parameter.In = in
		}
		if description, ok := param["description"].(string); ok {
			parameter.Description = description
		}
		if required, ok := param["required"].(bool); ok {
			parameter.Required = required
		}
		if schemaRaw, ok := param["schema"].(map[string]interface{}); ok {
			parameter.Schema = extractSchema(schemaRaw)
		}
		if example := param["example"]; example != nil {
			parameter.Example = example
		}

		parameters = append(parameters, parameter)
	}

	return parameters
}

// extractRequestBody extracts request body information
func extractRequestBody(requestBodyRaw map[string]interface{}) *RequestBody {
	requestBody := &RequestBody{
		Content: make(map[string]MediaType),
	}

	if description, ok := requestBodyRaw["description"].(string); ok {
		requestBody.Description = description
	}
	if required, ok := requestBodyRaw["required"].(bool); ok {
		requestBody.Required = required
	}
	if contentRaw, ok := requestBodyRaw["content"].(map[string]interface{}); ok {
		requestBody.Content = extractContent(contentRaw)
	}

	return requestBody
}

// extractResponses extracts response information
func extractResponses(responsesRaw map[string]interface{}) map[string]Response {
	responses := make(map[string]Response)

	for code, responseRaw := range responsesRaw {
		if response, ok := responseRaw.(map[string]interface{}); ok {
			resp := Response{}

			if description, ok := response["description"].(string); ok {
				resp.Description = description
			}
			if contentRaw, ok := response["content"].(map[string]interface{}); ok {
				resp.Content = extractContent(contentRaw)
			}

			responses[code] = resp
		}
	}

	return responses
}

// extractContent extracts media type content
func extractContent(contentRaw map[string]interface{}) map[string]MediaType {
	content := make(map[string]MediaType)

	for mediaType, mediaTypeRaw := range contentRaw {
		if mediaTypeData, ok := mediaTypeRaw.(map[string]interface{}); ok {
			mt := MediaType{}

			if schemaRaw, ok := mediaTypeData["schema"].(map[string]interface{}); ok {
				mt.Schema = extractSchema(schemaRaw)
			}
			if example := mediaTypeData["example"]; example != nil {
				mt.Example = example
			}

			content[mediaType] = mt
		}
	}

	return content
}

// extractSchema extracts schema information
func extractSchema(schemaRaw map[string]interface{}) Schema {
	schema := Schema{}

	if schemaType, ok := schemaRaw["type"].(string); ok {
		schema.Type = schemaType
	}
	if format, ok := schemaRaw["format"].(string); ok {
		schema.Format = format
	}
	if description, ok := schemaRaw["description"].(string); ok {
		schema.Description = description
	}
	if ref, ok := schemaRaw["$ref"].(string); ok {
		schema.Ref = ref
	}

	// Extract properties for object types
	if propertiesRaw, ok := schemaRaw["properties"].(map[string]interface{}); ok {
		schema.Properties = make(map[string]Schema)
		for propName, propSchemaRaw := range propertiesRaw {
			if propSchema, ok := propSchemaRaw.(map[string]interface{}); ok {
				schema.Properties[propName] = extractSchema(propSchema)
			}
		}
	}

	// Extract required fields
	if requiredRaw, ok := schemaRaw["required"].([]interface{}); ok {
		for _, reqRaw := range requiredRaw {
			if req, ok := reqRaw.(string); ok {
				schema.Required = append(schema.Required, req)
			}
		}
	}

	return schema
}
