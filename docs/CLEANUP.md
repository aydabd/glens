# Issue Cleanup

Clean up test issues created during integration testing.

## Quick Commands

```bash
# Preview (safe - shows what would be closed)
make cleanup-test-issues

# Execute (actually closes issues)
make cleanup-test-issues-confirm

# Test repository: aydabd/test-agent-ideas
```

## CLI Usage

```bash
# Dry-run (preview only)
glens cleanup --github-repo aydabd/test-agent-ideas --dry-run

# Close issues
glens cleanup --github-repo aydabd/test-agent-ideas

# Custom labels
glens cleanup --github-repo owner/repo --labels test-failure,integration-test
```

**Default label:** `ai-generated` (catches all AI-generated issues)

## Test Workflow

```bash
# 1. Run tests and create issues
make run-ollama-issues OP_ID=getPetById

# 2. Preview cleanup
make cleanup-test-issues

# 3. Clean up
make cleanup-test-issues-confirm
```

## Integration Tests

```bash
# Run integration tests (requires GITHUB_TOKEN)
make test-integration

# Run unit tests only
make test-short
```

## Programmatic Usage

```go
import (
    "context"
    "glens/pkg/github"
)

client, _ := github.NewClient(token)
client.SetRepository("aydabd/test-agent-ideas")

// List issues by label
issues, _ := client.ListIssuesByLabel(ctx, []string{"ai-generated"})

// Close specific issue
client.CloseIssue(ctx, issueNumber)

// Close all test issues
closedCount, _ := client.CloseTestIssues(ctx, []string{"ai-generated"})
```

## Issue Labels

- `ai-generated` - All AI-generated issues (default)
- `test-failure` - Main test failure issues
- `integration-test` - Integration test issues
- `subtask` - AI model subtasks
- `openapi` - OpenAPI-related
- HTTP methods: `get`, `post`, `put`, `delete`

## Environment Setup

```bash
export GITHUB_TOKEN=$(gh auth token)
export GITHUB_TEST_REPO=aydabd/test-agent-ideas  # optional
```

## Notes

- GitHub API **cannot delete** issues, only close them
- Always dry-run first (`--dry-run` flag)
- Test repository is safe for cleanup: `aydabd/test-agent-ideas`
- For permanent deletion, use GitHub web UI (admin only)

## Troubleshooting

```bash
# Token not set?
gh auth login
export GITHUB_TOKEN=$(gh auth token)

# Show help
glens cleanup --help
make help
```
