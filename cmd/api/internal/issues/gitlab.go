package issues

import (
	"context"
	"fmt"
)

func init() {
	Register("gitlab", NewGitLabProvider)
}

// GitLabProvider is a stub that returns "not implemented" for all operations.
type GitLabProvider struct{}

// NewGitLabProvider creates a GitLabProvider.
func NewGitLabProvider(_ ProviderConfig) (IssueProvider, error) {
	return &GitLabProvider{}, nil
}

// CreateIssue is not yet implemented.
func (g *GitLabProvider) CreateIssue(_ context.Context, _ CreateIssueRequest) (IssueResult, error) {
	return IssueResult{}, fmt.Errorf("gitlab: CreateIssue not implemented")
}

// UpdateIssue is not yet implemented.
func (g *GitLabProvider) UpdateIssue(_ context.Context, _ string, _ UpdateIssueRequest) error {
	return fmt.Errorf("gitlab: UpdateIssue not implemented")
}

// CloseIssue is not yet implemented.
func (g *GitLabProvider) CloseIssue(_ context.Context, _ string) error {
	return fmt.Errorf("gitlab: CloseIssue not implemented")
}

// ListIssues is not yet implemented.
func (g *GitLabProvider) ListIssues(_ context.Context, _ IssueFilter) ([]IssueResult, error) {
	return nil, fmt.Errorf("gitlab: ListIssues not implemented")
}
