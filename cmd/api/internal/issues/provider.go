package issues

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

// CreateIssueRequest contains fields for creating a new issue.
type CreateIssueRequest struct {
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Labels   []string `json:"labels,omitempty"`
	ParentID string   `json:"parent_id,omitempty"`
}

// UpdateIssueRequest contains fields for updating an existing issue.
type UpdateIssueRequest struct {
	Title  *string  `json:"title,omitempty"`
	Body   *string  `json:"body,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

// IssueFilter specifies criteria when listing issues.
type IssueFilter struct {
	State  string   `json:"state,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

// IssueResult represents the outcome of an issue operation.
type IssueResult struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Number int    `json:"number"`
}

// IssueProvider is the interface every issue-tracker backend must implement.
type IssueProvider interface {
	CreateIssue(ctx context.Context, req CreateIssueRequest) (IssueResult, error)
	UpdateIssue(ctx context.Context, id string, req UpdateIssueRequest) error
	CloseIssue(ctx context.Context, id string) error
	ListIssues(ctx context.Context, filter IssueFilter) ([]IssueResult, error)
}

// ProviderConfig holds connection details for a provider instance.
type ProviderConfig struct {
	Name          string `json:"name"`
	Repository    string `json:"repository"`
	CredentialRef string `json:"credential_ref"`
}

// ProviderFactory creates an IssueProvider from a config.
type ProviderFactory func(cfg ProviderConfig) (IssueProvider, error)

var (
	mu        sync.RWMutex
	providers = make(map[string]ProviderFactory)
)

// Register adds a provider factory under the given name.
func Register(name string, factory ProviderFactory) {
	mu.Lock()
	defer mu.Unlock()
	providers[name] = factory
}

// NewProvider creates a provider by name using the supplied config.
func NewProvider(name string, cfg ProviderConfig) (IssueProvider, error) {
	mu.RLock()
	factory, ok := providers[name]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown issue provider: %s", name)
	}
	return factory(cfg)
}

// ListProviders returns the sorted names of all registered providers.
func ListProviders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(providers))
	for n := range providers {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
