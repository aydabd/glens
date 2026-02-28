# Release Plan — Glens Monorepo

> Comprehensive release strategy for all modules, libraries, and services.
> Cross-references [Phase 16 — Releases & Versioning](saas/PHASE16_RELEASES_VERSIONING.md).

## Principles

1. **Isolated releases** — each module is versioned and released independently.
2. **Conventional commits** — version bumps are fully automated via commit messages.
3. **Signed artifacts** — binaries ship with SHA256 checksums and GPG signatures.
4. **Zero manual steps** — merging to `main` is the only human action required.

## Release Automation — Google Release Please

[Release Please](https://github.com/googleapis/release-please) manages version
bumps, changelogs, and release PRs for every module in the monorepo.

Configuration files:

| File | Purpose |
|------|---------|
| `release-please-config.json` | Per-module release type, tag format, changelog options |
| `.release-please-manifest.json` | Current version of each module (auto-updated) |
| `.github/workflows/release-please.yml` | GitHub Actions workflow triggered on `main` push |

### How it works

```text
Developer PR (conventional commits) → merge to main
  │
  └─ Release Please detects changes per module path
       │
       ├─ Opens (or updates) a Release PR per module
       │    • Bumps version in manifest
       │    • Generates CHANGELOG.md entries
       │
       └─ Merge Release PR → creates GitHub Release + git tag
            │
            ├─ Binary modules: cross-compile + sign + upload assets
            └─ Library modules: tag-only release (no binaries)
```

## Module Release Matrix

| Module | Path | Tag format | Binary | GPG signed | Release workflow |
|--------|------|------------|--------|------------|-----------------|
| **logging** | `pkg/logging` | `pkg/logging/v*` | No | N/A | `release-please.yml` / `release-pkg-logging.yml` |
| **glens** | `cmd/glens` | `cmd/glens/v*` | Yes | Yes | `release-please.yml` / `release-glens.yml` |
| **api** | `cmd/api` | `cmd/api/v*` | Yes | Yes | `release-please.yml` / `release-api.yml` |
| **demo** | `cmd/tools/demo` | `cmd/tools/demo/v*` | Yes | Yes | `release-please.yml` / `release-demo.yml` |
| **accuracy** | `cmd/tools/accuracy` | `cmd/tools/accuracy/v*` | Yes | Yes | `release-please.yml` / `release-accuracy.yml` |

### Binary build targets

All binary modules produce artifacts for:

| OS | Architecture |
|----|-------------|
| linux | amd64, arm64 |
| darwin | amd64, arm64 |
| windows | amd64, arm64 |

## Conventional Commits

All commits **must** follow the
[Conventional Commits](https://www.conventionalcommits.org/) specification.
This is enforced by the `conventional-pre-commit` hook in `.pre-commit-config.yaml`.

### Format

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Allowed types

| Type | Description | Version bump |
|------|-------------|-------------|
| `feat` | New feature | **minor** |
| `feat!` | Breaking feature | **major** |
| `fix` | Bug fix | **patch** |
| `docs` | Documentation only | no release |
| `style` | Formatting, no code change | no release |
| `refactor` | Code change, no feature/fix | **patch** |
| `perf` | Performance improvement | **patch** |
| `test` | Adding/correcting tests | no release |
| `build` | Build system or dependencies | **patch** |
| `ci` | CI configuration | no release |
| `chore` | Maintenance tasks | no release |
| `revert` | Revert a previous commit | **patch** |

A `BREAKING CHANGE:` footer in any commit type triggers a **major** bump.

### Scoping commits to modules

Use the module directory as scope to clearly attribute changes:

```text
feat(cmd/glens): add --output flag for JSON reports
fix(pkg/logging): handle nil logger gracefully
feat(cmd/api)!: redesign authentication endpoint
```

## Asset Signing

### SHA256 checksums

Every release with binary assets includes a `checksums.txt` file containing
SHA256 hashes for all artifacts:

```text
e3b0c44298fc1c149afbf4c...  glens-linux-amd64
a7ffc6f8bf1ed766...          glens-linux-arm64
...
```

### GPG signatures

When the repository secrets `GPG_PRIVATE_KEY` and `GPG_PASSPHRASE` are
configured, release workflows produce a `checksums.txt.asc` detached
signature file.

Users verify downloads with:

```bash
# Import the project's public key
gpg --import glens-release-key.pub

# Verify checksums file
gpg --verify checksums.txt.asc checksums.txt

# Verify individual binary
sha256sum -c checksums.txt --ignore-missing
```

### Setting up GPG signing

1. Generate a dedicated release signing key:
   ```bash
   gpg --batch --gen-key <<EOF
   Key-Type: RSA
   Key-Length: 4096
   Name-Real: Glens Release Signing
   Name-Email: release@glens.dev
   Expire-Date: 2y
   %no-protection
   EOF
   ```
2. Export and add to GitHub repository secrets:
   ```bash
   gpg --armor --export-secret-keys release@glens.dev  # → GPG_PRIVATE_KEY
   ```
3. Publish the public key in the repository or a keyserver.

## SaaS Service Releases

For the planned SaaS transformation (see [SAAS_PLAN.md](saas/SAAS_PLAN.md)),
service deployments follow the Phase 16 release flow:

```text
Release PR → merge → tag v1.2.3-rc.1
  ├─ Auto-deploy to dev environment
  ├─ Run regression + smoke tests
  └─ All pass → promote: tag v1.2.3
       ├─ Auto-deploy to prod
       └─ Post-deploy monitoring (error rate < 0.1%)
```

| Service | Deployment target | Pre-release tag | Stable tag |
|---------|------------------|-----------------|------------|
| API (`cmd/api`) | Cloud Run | `cmd/api/v*-rc.*` | `cmd/api/v*` |
| Frontend | Cloud Functions | `frontend/v*-rc.*` | `frontend/v*` |
| Infrastructure | Terraform | `infra/v*-rc.*` | `infra/v*` |

## Implementation Checklist

Cross-referenced with [Phase 16 success criteria](saas/PHASE16_RELEASES_VERSIONING.md):

- [x] Semver tagging per module (`pkg/logging/v*`, `cmd/glens/v*`, etc.)
- [x] Conventional commit-based version bumping
- [x] Cross-platform binary builds (linux/darwin/windows × amd64/arm64)
- [x] SHA256 checksums for release assets
- [x] Reusable release workflow (`release-module.yml`)
- [x] Individual release workflows for each module
- [x] Missing `release-api.yml` workflow added
- [x] Release Please configuration for monorepo
- [x] Release Please GitHub Actions workflow
- [x] GPG signing support in release workflows
- [x] Conventional commit pre-commit hook enforcement
- [x] AI agent instructions for conventional commits
- [ ] GPG release signing key generated and added to repository secrets
- [ ] Release Please initial run and first release PR
- [ ] RC tags auto-deploy to dev; stable tags to prod (Phase 16: RV-02, RV-03)
- [ ] Regression suite gates prod promotion (Phase 16: RV-04)
- [ ] Smoke tests against deployed environments (Phase 16: RV-05)
- [ ] Deployment e2e tests (Phase 16: RV-06)
- [ ] Rollback restores previous version in < 2 min
