# Phase 17 — Issue Tracker Provider Abstraction

> Pluggable issue creation: GitHub first, then GitLab, Jira.

## Requirements Covered

IP-01 (provider interface), IP-02 (GitHub provider), IP-03 (GitLab
future), IP-04 (Jira future), IP-05 (provider config via UI).

## Problem

Current `internal/github/client.go` is hard-coded to GitHub API.
SaaS users may use GitLab, Jira, or other trackers. We need a
provider pattern so implementations are swappable.

## Interface Design

```go
type IssueProvider interface {
    CreateIssue(ctx context.Context, req CreateIssueRequest) (IssueResult, error)
    UpdateIssue(ctx context.Context, id string, req UpdateIssueRequest) error
    CloseIssue(ctx context.Context, id string) error
    ListIssues(ctx context.Context, filter IssueFilter) ([]IssueResult, error)
}

type CreateIssueRequest struct {
    Title, Body string; Labels []string; ParentID string
}
type IssueResult struct {
    ID, URL string; Number int
}
```

## Provider Registry

```go
var providers = map[string]func(ProviderConfig) (IssueProvider, error){
    "github": NewGitHubProvider,
    "gitlab": NewGitLabProvider, // future
    "jira":   NewJiraProvider,   // future
}
func NewProvider(name string, cfg ProviderConfig) (IssueProvider, error) {
    f, ok := providers[name]
    if !ok { return nil, fmt.Errorf("unknown provider: %s", name) }
    return f(cfg)
}
```

## GitHub Provider (IP-02)

```go
func (g *GitHubProvider) CreateIssue(ctx context.Context,
    req CreateIssueRequest) (IssueResult, error) {
    issue, _, err := g.client.Issues.Create(ctx, g.owner, g.repo,
        &github.IssueRequest{Title: &req.Title, Body: &req.Body,
            Labels: &req.Labels})
    if err != nil { return IssueResult{}, fmt.Errorf("github: %w", err) }
    return IssueResult{ID: fmt.Sprintf("%d", issue.GetNumber()),
        URL: issue.GetHTMLURL(), Number: issue.GetNumber()}, nil
}
```

## Future Providers (IP-03, IP-04)

Future: GitLab (REST v4, project token), Jira (REST v3, OAuth2).
Not implemented now — interface ensures they slot in cleanly.

## Workspace Config

```json
{ "issue_provider": "github",
  "provider_config": {
    "repository": "owner/repo",
    "credential_ref": "projects/x/secrets/github-token/versions/1" } }
```

Credential is a Secret Manager ref (never raw — see Phase 7).
Frontend: provider dropdown (GitHub | GitLab | Jira coming soon);
credential stored via `POST /api/v1/secrets`.

## Package: `cmd/api/internal/issues/`

`provider.go` (interface + registry), `github.go`, `gitlab.go` (stub),
`jira.go` (stub), `provider_test.go` (compliance tests).

## Steps

1. Define `IssueProvider` interface + provider registry
2. Refactor `internal/github/` → `GitHubProvider`
3. Update event subscriber to use provider; add workspace config

## Success Criteria

- [ ] `IssueProvider` interface defined with 4 methods
- [ ] GitHub provider passes all interface compliance tests
- [ ] Existing issue-creation flow unchanged (backward compatible)
- [ ] Provider selected per workspace via config
