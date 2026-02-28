package issues

import (
	"context"
	"fmt"
)

func init() {
	Register("jira", NewJiraProvider)
}

// JiraProvider is a stub that returns "not implemented" for all operations.
type JiraProvider struct{}

// NewJiraProvider creates a JiraProvider.
func NewJiraProvider(_ ProviderConfig) (IssueProvider, error) {
	return &JiraProvider{}, nil
}

// CreateIssue is not yet implemented.
func (j *JiraProvider) CreateIssue(_ context.Context, _ CreateIssueRequest) (IssueResult, error) {
	return IssueResult{}, fmt.Errorf("jira: CreateIssue not implemented")
}

// UpdateIssue is not yet implemented.
func (j *JiraProvider) UpdateIssue(_ context.Context, _ string, _ UpdateIssueRequest) error {
	return fmt.Errorf("jira: UpdateIssue not implemented")
}

// CloseIssue is not yet implemented.
func (j *JiraProvider) CloseIssue(_ context.Context, _ string) error {
	return fmt.Errorf("jira: CloseIssue not implemented")
}

// ListIssues is not yet implemented.
func (j *JiraProvider) ListIssues(_ context.Context, _ IssueFilter) ([]IssueResult, error) {
	return nil, fmt.Errorf("jira: ListIssues not implemented")
}
