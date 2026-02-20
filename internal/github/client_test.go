package github

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"glens/internal/parser"
)

// TestGitHubIntegration tests the complete GitHub issue lifecycle
// This is an integration test that requires:
// - GITHUB_TOKEN environment variable
// - A test repository (e.g., aydabd/test-agent-ideas)
// - Internet connection
//
// Run with: go test -v -tags=integration ./pkg/github/...
// Or use: make test-integration
func TestGitHubIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Get credentials from environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("GITHUB_TOKEN not set, skipping integration test")
	}

	// Use test repository
	testRepo := os.Getenv("GITHUB_TEST_REPO")
	if testRepo == "" {
		testRepo = "aydabd/test-agent-ideas"
	}

	ctx := context.Background()

	// Create client
	client, err := NewClient(token)
	require.NoError(t, err, "Failed to create GitHub client")

	err = client.SetRepository(testRepo)
	require.NoError(t, err, "Failed to set repository")

	t.Run("CreateAndCleanupIssue", func(t *testing.T) {
		// Create a test endpoint
		endpoint := &parser.Endpoint{
			Method:      "GET",
			Path:        "/test/integration",
			OperationID: "testIntegration",
			Summary:     "Integration test endpoint",
			Description: "This is a test issue created by glens integration tests",
			Parameters: []parser.Parameter{
				{
					Name:        "id",
					In:          "path",
					Required:    true,
					Description: "Test ID",
					Schema: parser.Schema{
						Type: "string",
					},
				},
			},
			Responses: map[string]parser.Response{
				"200": {
					Description: "Success",
				},
				"404": {
					Description: "Not found",
				},
			},
		}

		aiModels := []string{"ollama"}

		// Step 1: Create issue
		t.Log("Creating test issue...")
		issueNumber, err := client.CreateEndpointIssue(ctx, endpoint, aiModels)
		require.NoError(t, err, "Failed to create issue")
		assert.Greater(t, issueNumber, 0, "Issue number should be positive")
		t.Logf("Created issue #%d", issueNumber)

		// Step 2: Update issue with test results
		t.Log("Updating issue with test results...")
		testResults := `
### Test Results
- ✅ Test 1: Passed
- ❌ Test 2: Failed
- ✅ Test 3: Passed

**Failure Details:**
Test 2 failed because of invalid response format.
`
		err = client.UpdateIssueWithResults(ctx, issueNumber, testResults)
		assert.NoError(t, err, "Failed to update issue with results")

		// Step 3: List issues by label
		t.Log("Listing issues by label...")
		issues, err := client.ListIssuesByLabel(ctx, []string{"ai-generated", "test-failure"})
		require.NoError(t, err, "Failed to list issues")
		assert.NotEmpty(t, issues, "Should have at least one issue")

		// Verify our issue is in the list
		found := false
		for _, issue := range issues {
			if issue.GetNumber() == issueNumber {
				found = true
				t.Logf("Found our issue in the list: #%d - %s", issue.GetNumber(), issue.GetTitle())
				break
			}
		}
		assert.True(t, found, "Created issue should be in the list")

		// Step 4: Close the test issue (cleanup)
		t.Log("Cleaning up test issue...")
		err = client.CloseIssue(ctx, issueNumber)
		assert.NoError(t, err, "Failed to close issue")
		t.Logf("Successfully closed issue #%d", issueNumber)
	})

	t.Run("CloseMultipleTestIssues", func(t *testing.T) {
		// This test demonstrates bulk cleanup
		t.Log("Testing bulk issue cleanup...")

		// Count open issues before
		issuesBefore, err := client.ListIssuesByLabel(ctx, []string{"ai-generated"})
		require.NoError(t, err, "Failed to list issues before cleanup")
		openBefore := 0
		for _, issue := range issuesBefore {
			if issue.GetState() == "open" {
				openBefore++
			}
		}
		t.Logf("Found %d open issues before cleanup", openBefore)

		// Close all test issues
		closedCount, err := client.CloseTestIssues(ctx, []string{"ai-generated"})
		require.NoError(t, err, "Failed to close test issues")
		t.Logf("Closed %d issues", closedCount)

		// Verify all are closed
		issuesAfter, err := client.ListIssuesByLabel(ctx, []string{"ai-generated"})
		require.NoError(t, err, "Failed to list issues after cleanup")
		openAfter := 0
		for _, issue := range issuesAfter {
			if issue.GetState() == "open" {
				openAfter++
			}
		}
		t.Logf("Found %d open issues after cleanup", openAfter)
		assert.Equal(t, 0, openAfter, "All test issues should be closed")
	})
}

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   "ghp_test123",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

// TestSetRepository tests repository setting
func TestSetRepository(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name      string
		repo      string
		wantErr   bool
		wantOwner string
		wantRepo  string
	}{
		{
			name:      "valid repository",
			repo:      "owner/repo",
			wantErr:   false,
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:    "invalid format - no slash",
			repo:    "owner-repo",
			wantErr: true,
		},
		{
			name:    "invalid format - too many parts",
			repo:    "owner/repo/extra",
			wantErr: true,
		},
		{
			name:    "empty string",
			repo:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.SetRepository(tt.repo)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantOwner, client.owner)
				assert.Equal(t, tt.wantRepo, client.repo)
			}
		})
	}
}
