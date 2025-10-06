# Claude Instructions for Glens Project

## Project Philosophy

**Simplicity First**: Always choose the simplest solution that works.
**Long-term Support**: Design for minimal maintenance and easy understanding.
**Quality Assurance**: Every change must be testable and maintainable.
**Go Best Practices**: Follow idiomatic Go patterns and conventions.

## Core Principles

### 1. Code Changes

- Write idiomatic Go code (gofmt, golint compliant)
- Keep functions small and focused (max 50 lines)
- Use meaningful variable names, avoid abbreviations
- Add comments only when code isn't self-explanatory
- Handle errors explicitly, never ignore them
- Write tests for new functionality

### 2. Documentation

- Keep docs minimal and focused on "why", not "how"
- Update only affected docs when code changes
- Avoid duplication - single source of truth
- Use README.md for user-facing info
- Use DEVELOPMENT.md for contributor info
- Keep inline code comments for complex logic only
- Visual diagrams in docs/diagrams/ - update only if architecture changes

### 3. Simplicity Rules

- No premature optimization
- No complex abstractions until needed
- Prefer standard library over third-party deps
- Keep file structure flat and obvious
- One concept per file when possible

### 4. GitHub Issue Creation Logic

**Critical**: Issues are created ONLY when tests fail against OpenAPI spec

- NOT for connection errors
- NOT for compilation errors
- NOT for infrastructure issues
- ONLY for real API specification violations

When modifying this logic, preserve the distinction between:

- Real test failures (spec violations)
- Setup/infrastructure failures

### 5. Testing

- Write table-driven tests
- Use testify for assertions
- Mock external dependencies
- Test error cases, not just happy paths
- Keep test code as simple as production code

### 6. File Organization

```txt
cmd/        - CLI commands (one per file)
pkg/        - Reusable packages (focused, small)
docs/       - Only QUICKSTART.md and DEVELOPMENT.md
configs/    - Example configurations only
```

### 7. Before Any Change

Ask yourself:

1. Is this the simplest solution?
2. Will this be easy to maintain in 2 years?
3. Can someone new understand this in 5 minutes?
4. Does this follow Go idioms?
5. Is this change necessary?

If any answer is "no", reconsider the approach.

### 8. Documentation Updates

Only update docs when:

- Public API changes
- New CLI flags added
- Configuration format changes
- Setup process changes

Do NOT update docs for:

- Internal refactoring
- Bug fixes (unless behavior changes)
- Code organization changes
- Performance improvements

### 9. Common Patterns to Follow

```go
// Error handling
if err != nil {
    return fmt.Errorf("context: %w", err)
}

// Configuration
type Config struct {
    // Required fields first
    Token string `yaml:"token"`
    // Optional fields after
    Debug bool `yaml:"debug,omitempty"`
}

// Logging
log.Info().
    Str("key", value).
    Msg("What happened")
```

### 10. What NOT to Do

- ❌ Don't add features "just in case"
- ❌ Don't create complex inheritance hierarchies
- ❌ Don't use reflection unless absolutely necessary
- ❌ Don't write documentation novels
- ❌ Don't create multiple docs saying the same thing
- ❌ Don't add dependencies without strong justification

## Current Project State

This is an active development project. Code and features may change.
Focus on working code over extensive documentation.
