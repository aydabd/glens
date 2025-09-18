package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"

	"glens/pkg/parser"
)

// Client wraps GitHub API operations
type Client struct {
	client *github.Client
	owner  string
	repo   string
}

// NewClient creates a new GitHub client
func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("GitHub token is required")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &Client{
		client: github.NewClient(tc),
	}, nil
}

// SetRepository sets the target repository
func (c *Client) SetRepository(repository string) error {
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return fmt.Errorf("repository must be in format 'owner/repo'")
	}

	c.owner = parts[0]
	c.repo = parts[1]

	log.Debug().
		Str("owner", c.owner).
		Str("repo", c.repo).
		Msg("Repository set")

	return nil
}

// CreateEndpointIssue creates a GitHub issue for an endpoint with AI model subtasks
func (c *Client) CreateEndpointIssue(ctx context.Context, endpoint *parser.Endpoint, aiModels []string) (int, error) {
	title := fmt.Sprintf("Integration Test: %s %s", endpoint.Method, endpoint.Path)

	body := c.generateIssueBody(endpoint, aiModels)

	issue := &github.IssueRequest{
		Title: &title,
		Body:  &body,
		Labels: &[]string{
			"integration-test",
			"ai-generated",
			"openapi",
			strings.ToLower(endpoint.Method),
		},
	}

	createdIssue, _, err := c.client.Issues.Create(ctx, c.owner, c.repo, issue)
	if err != nil {
		return 0, fmt.Errorf("failed to create issue: %w", err)
	}

	issueNumber := createdIssue.GetNumber()

	log.Info().
		Int("issue_number", issueNumber).
		Str("endpoint", fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)).
		Msg("GitHub issue created")

	// Create subtasks for each AI model
	for _, aiModel := range aiModels {
		if err := c.createSubtask(ctx, issueNumber, endpoint, aiModel); err != nil {
			log.Error().
				Err(err).
				Str("ai_model", aiModel).
				Int("parent_issue", issueNumber).
				Msg("Failed to create subtask")
		}
	}

	return issueNumber, nil
}

// generateIssueBody creates the markdown body for the main issue
func (c *Client) generateIssueBody(endpoint *parser.Endpoint, aiModels []string) string {
	var body strings.Builder

	body.WriteString("## 🎯 Endpoint Integration Test\n\n")
	body.WriteString(fmt.Sprintf("**Method:** `%s`\n", endpoint.Method))
	body.WriteString(fmt.Sprintf("**Path:** `%s`\n", endpoint.Path))

	if endpoint.OperationID != "" {
		body.WriteString(fmt.Sprintf("**Operation ID:** `%s`\n", endpoint.OperationID))
	}

	if endpoint.Summary != "" {
		body.WriteString(fmt.Sprintf("**Summary:** %s\n", endpoint.Summary))
	}

	if endpoint.Description != "" {
		body.WriteString(fmt.Sprintf("\n**Description:**\n%s\n", endpoint.Description))
	}

	// Parameters section
	if len(endpoint.Parameters) > 0 {
		body.WriteString("\n### 📋 Parameters\n\n")
		body.WriteString("| Name | Type | In | Required | Description |\n")
		body.WriteString("|------|------|----|---------|--------------|\n")

		for i := range endpoint.Parameters {
			param := &endpoint.Parameters[i]
			required := "No"
			if param.Required {
				required = "Yes"
			}
			body.WriteString(fmt.Sprintf("| `%s` | `%s` | `%s` | %s | %s |\n",
				param.Name, param.Schema.Type, param.In, required, param.Description))
		}
	}

	// Request body section
	if endpoint.RequestBody != nil {
		body.WriteString("\n### 📤 Request Body\n\n")
		if endpoint.RequestBody.Description != "" {
			body.WriteString(fmt.Sprintf("**Description:** %s\n\n", endpoint.RequestBody.Description))
		}
		body.WriteString("**Content Types:**\n")
		for contentType := range endpoint.RequestBody.Content {
			body.WriteString(fmt.Sprintf("- `%s`\n", contentType))
		}
	}

	// Responses section
	if len(endpoint.Responses) > 0 {
		body.WriteString("\n### 📥 Responses\n\n")
		body.WriteString("| Status Code | Description |\n")
		body.WriteString("|-------------|-------------|\n")

		for code, response := range endpoint.Responses {
			body.WriteString(fmt.Sprintf("| `%s` | %s |\n", code, response.Description))
		}
	}

	// AI Models section
	body.WriteString("\n### 🤖 AI Models for Test Generation\n\n")
	body.WriteString("This endpoint will be tested using the following AI models:\n\n")

	for _, model := range aiModels {
		body.WriteString(fmt.Sprintf("- [ ] **%s** - Integration test generation and execution\n", model))
	}

	body.WriteString("\n### 🧪 Test Categories\n\n")
	body.WriteString("Each AI model will generate tests covering:\n\n")
	body.WriteString("- [ ] **Happy Path Testing** - Valid requests with expected responses\n")
	body.WriteString("- [ ] **Error Handling** - Invalid inputs and error responses\n")
	body.WriteString("- [ ] **Boundary Testing** - Edge cases and limits\n")
	body.WriteString("- [ ] **Security Testing** - Authentication and authorization\n")
	body.WriteString("- [ ] **Performance Testing** - Response times and load handling\n")

	body.WriteString("\n### 📊 Success Criteria\n\n")
	body.WriteString("- All generated tests execute successfully\n")
	body.WriteString("- Tests cover all documented response codes\n")
	body.WriteString("- Security requirements are validated\n")
	body.WriteString("- Performance benchmarks are met\n")
	body.WriteString("- AI model comparison report is generated\n")

	body.WriteString("\n---\n")
	body.WriteString("*This issue was automatically generated by Glens*")

	return body.String()
}

// createSubtask creates a subtask issue for a specific AI model
func (c *Client) createSubtask(ctx context.Context, parentIssue int, endpoint *parser.Endpoint, aiModel string) error {
	title := fmt.Sprintf("[%s] Generate tests for %s %s", aiModel, endpoint.Method, endpoint.Path)

	body := c.generateSubtaskBody(parentIssue, endpoint, aiModel)

	issue := &github.IssueRequest{
		Title: &title,
		Body:  &body,
		Labels: &[]string{
			"integration-test",
			"ai-generated",
			"subtask",
			strings.ToLower(aiModel),
			strings.ToLower(endpoint.Method),
		},
	}

	createdIssue, _, err := c.client.Issues.Create(ctx, c.owner, c.repo, issue)
	if err != nil {
		return fmt.Errorf("failed to create subtask: %w", err)
	}

	// Add comment to parent issue linking to subtask
	comment := fmt.Sprintf("🤖 **%s Subtask Created:** #%d", aiModel, createdIssue.GetNumber())
	_, _, err = c.client.Issues.CreateComment(ctx, c.owner, c.repo, parentIssue, &github.IssueComment{
		Body: &comment,
	})

	if err != nil {
		log.Error().
			Err(err).
			Int("parent_issue", parentIssue).
			Int("subtask_issue", createdIssue.GetNumber()).
			Msg("Failed to link subtask to parent issue")
	}

	log.Debug().
		Int("subtask_issue", createdIssue.GetNumber()).
		Int("parent_issue", parentIssue).
		Str("ai_model", aiModel).
		Msg("Subtask created")

	return nil
}

// generateSubtaskBody creates the markdown body for AI model subtasks
func (c *Client) generateSubtaskBody(parentIssue int, endpoint *parser.Endpoint, aiModel string) string {
	var body strings.Builder

	body.WriteString(fmt.Sprintf("## 🤖 %s Integration Test Generation\n\n", aiModel))
	body.WriteString(fmt.Sprintf("**Parent Issue:** #%d\n", parentIssue))
	body.WriteString(fmt.Sprintf("**Endpoint:** `%s %s`\n", endpoint.Method, endpoint.Path))
	body.WriteString(fmt.Sprintf("**AI Model:** %s\n\n", aiModel))

	body.WriteString("### 🎯 Objective\n\n")
	body.WriteString(fmt.Sprintf("Generate comprehensive integration tests for the `%s %s` endpoint using the %s AI model.\n\n",
		endpoint.Method, endpoint.Path, aiModel))

	body.WriteString("### 📋 Tasks\n\n")
	body.WriteString("- [ ] **Analyze Endpoint Specification**\n")
	body.WriteString("  - Review parameters, request body, and response schemas\n")
	body.WriteString("  - Identify security requirements\n")
	body.WriteString("  - Understand business logic constraints\n\n")

	body.WriteString("- [ ] **Generate Test Cases**\n")
	body.WriteString("  - Happy path scenarios\n")
	body.WriteString("  - Error handling cases\n")
	body.WriteString("  - Boundary value testing\n")
	body.WriteString("  - Security validation\n\n")

	body.WriteString("- [ ] **Create Test Code**\n")
	body.WriteString("  - Generate executable test code\n")
	body.WriteString("  - Include proper assertions\n")
	body.WriteString("  - Add test data generation\n")
	body.WriteString("  - Implement cleanup procedures\n\n")

	body.WriteString("- [ ] **Execute Tests**\n")
	body.WriteString("  - Run generated test suite\n")
	body.WriteString("  - Capture execution results\n")
	body.WriteString("  - Document any failures\n")
	body.WriteString("  - Generate performance metrics\n\n")

	body.WriteString("### 🔍 Test Focus Areas\n\n")

	if len(endpoint.Parameters) > 0 {
		body.WriteString("**Parameters to Test:**\n")
		for i := range endpoint.Parameters {
			param := &endpoint.Parameters[i]
			required := "optional"
			if param.Required {
				required = "required"
			}
			body.WriteString(fmt.Sprintf("- `%s` (%s, %s): %s\n",
				param.Name, param.In, required, param.Description))
		}
		body.WriteString("\n")
	}

	if endpoint.RequestBody != nil {
		body.WriteString("**Request Body Testing:**\n")
		body.WriteString("- Valid payload structures\n")
		body.WriteString("- Invalid/malformed data\n")
		body.WriteString("- Missing required fields\n")
		body.WriteString("- Content-type validation\n\n")
	}

	if len(endpoint.Responses) > 0 {
		body.WriteString("**Response Validation:**\n")
		for code := range endpoint.Responses {
			body.WriteString(fmt.Sprintf("- HTTP %s response handling\n", code))
		}
		body.WriteString("\n")
	}

	body.WriteString("### 🛠 Technical Requirements\n\n")
	body.WriteString("- **Framework:** Go with testify\n")
	body.WriteString("- **HTTP Client:** Standard library or custom\n")
	body.WriteString("- **Assertions:** Comprehensive validation\n")
	body.WriteString("- **Documentation:** Clear test descriptions\n")
	body.WriteString("- **Maintainability:** Readable and modular code\n\n")

	body.WriteString("### 📊 Success Criteria\n\n")
	body.WriteString("- [ ] All test cases execute without compilation errors\n")
	body.WriteString("- [ ] Tests demonstrate endpoint functionality\n")
	body.WriteString("- [ ] Error scenarios are properly handled\n")
	body.WriteString("- [ ] Performance metrics are captured\n")
	body.WriteString("- [ ] Test results are documented\n\n")

	body.WriteString("### 📈 Deliverables\n\n")
	body.WriteString("1. **Generated Test Code** - Complete test suite\n")
	body.WriteString("2. **Execution Report** - Test run results\n")
	body.WriteString("3. **Performance Metrics** - Response time analysis\n")
	body.WriteString("4. **Issue Report** - Any discovered problems\n")
	body.WriteString("5. **AI Prompt Details** - Prompt used for generation\n\n")

	body.WriteString("---\n")
	body.WriteString(fmt.Sprintf("*Generated by Glens for %s*", aiModel))

	return body.String()
}

// UpdateIssueWithResults updates an issue with test execution results
func (c *Client) UpdateIssueWithResults(ctx context.Context, issueNumber int, results string) error {
	comment := fmt.Sprintf("## 📊 Test Execution Results\n\n%s", results)

	_, _, err := c.client.Issues.CreateComment(ctx, c.owner, c.repo, issueNumber, &github.IssueComment{
		Body: &comment,
	})

	if err != nil {
		return fmt.Errorf("failed to update issue with results: %w", err)
	}

	return nil
}

// CloseIssue closes an issue when testing is complete
func (c *Client) CloseIssue(ctx context.Context, issueNumber int) error {
	state := "closed"
	_, _, err := c.client.Issues.Edit(ctx, c.owner, c.repo, issueNumber, &github.IssueRequest{
		State: &state,
	})

	if err != nil {
		return fmt.Errorf("failed to close issue: %w", err)
	}

	return nil
}
