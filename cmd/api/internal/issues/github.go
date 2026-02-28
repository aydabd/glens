package issues

import (
	"context"
	"fmt"
	"strings"
)

func init() {
	Register("github", NewGitHubProvider)
}

// GitHubProvider is a stub implementation of IssueProvider for GitHub.
type GitHubProvider struct {
	owner string
	repo  string
}

// NewGitHubProvider creates a GitHubProvider after validating the config.
func NewGitHubProvider(cfg ProviderConfig) (IssueProvider, error) {
	parts := strings.Split(cfg.Repository, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("github: repository must be in owner/repo format, got %q", cfg.Repository)
	}
	return &GitHubProvider{owner: parts[0], repo: parts[1]}, nil
}

func (g *GitHubProvider) CreateIssue(_ context.Context, req CreateIssueRequest) (IssueResult, error) {
	return IssueResult{
		ID:     fmt.Sprintf("%s/%s/1", g.owner, g.repo),
		URL:    fmt.Sprintf("https://github.com/%s/%s/issues/1", g.owner, g.repo),
		Number: 1,
	}, nil
}

func (g *GitHubProvider) UpdateIssue(_ context.Context, _ string, _ UpdateIssueRequest) error {
	return fmt.Errorf("github: UpdateIssue not implemented")
}

func (g *GitHubProvider) CloseIssue(_ context.Context, _ string) error {
	return fmt.Errorf("github: CloseIssue not implemented")
}

func (g *GitHubProvider) ListIssues(_ context.Context, _ IssueFilter) ([]IssueResult, error) {
	return []IssueResult{}, nil
}
