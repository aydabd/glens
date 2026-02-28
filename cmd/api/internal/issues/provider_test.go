package issues

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProviders_ReturnsAllRegistered(t *testing.T) {
	names := ListProviders()
	assert.Contains(t, names, "github")
	assert.Contains(t, names, "gitlab")
	assert.Contains(t, names, "jira")
}

func TestListProviders_ReturnsSorted(t *testing.T) {
	names := ListProviders()
	for i := 1; i < len(names); i++ {
		assert.True(t, names[i-1] <= names[i], "providers should be sorted")
	}
}

func TestNewProvider_UnknownProvider_ReturnsError(t *testing.T) {
	_, err := NewProvider("nonexistent", ProviderConfig{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown issue provider")
}

func TestNewProvider_GitHub_ValidConfig(t *testing.T) {
	p, err := NewProvider("github", ProviderConfig{
		Repository: "owner/repo",
	})
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewProvider_GitHub_InvalidConfig(t *testing.T) {
	tests := []struct {
		name string
		repo string
	}{
		{"empty string", ""},
		{"missing slash", "noslash"},
		{"trailing slash", "owner/"},
		{"leading slash", "/repo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProvider("github", ProviderConfig{
				Repository: tt.repo,
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "owner/repo")
		})
	}
}

func TestGitHubProvider_CreateIssue_ReturnsStubResult(t *testing.T) {
	p, err := NewProvider("github", ProviderConfig{Repository: "myorg/myrepo"})
	require.NoError(t, err)

	result, err := p.CreateIssue(context.Background(), CreateIssueRequest{
		Title: "test issue",
		Body:  "body",
	})
	require.NoError(t, err)
	assert.Equal(t, 1, result.Number)
	assert.Contains(t, result.URL, "myorg/myrepo")
}

func TestGitHubProvider_ListIssues_ReturnsEmptySlice(t *testing.T) {
	p, err := NewProvider("github", ProviderConfig{Repository: "o/r"})
	require.NoError(t, err)

	results, err := p.ListIssues(context.Background(), IssueFilter{})
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestStubProviders_ReturnNotImplemented(t *testing.T) {
	stubProviders := []string{"gitlab", "jira"}
	for _, name := range stubProviders {
		t.Run(name, func(t *testing.T) {
			p, err := NewProvider(name, ProviderConfig{Repository: "o/r"})
			require.NoError(t, err)

			ctx := context.Background()

			_, err = p.CreateIssue(ctx, CreateIssueRequest{Title: "t"})
			assert.ErrorContains(t, err, "not implemented")

			err = p.UpdateIssue(ctx, "1", UpdateIssueRequest{})
			assert.ErrorContains(t, err, "not implemented")

			err = p.CloseIssue(ctx, "1")
			assert.ErrorContains(t, err, "not implemented")

			_, err = p.ListIssues(ctx, IssueFilter{})
			assert.ErrorContains(t, err, "not implemented")
		})
	}
}
