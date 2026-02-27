# Phase 15 — Environment Parity (Dev / Prod)

> Identical infra with only log level and feature flags differing.

## Requirements Covered

EP-01 (dev account), EP-02 (prod account), EP-03 (parity guarantee),
EP-04 (debug in dev only), EP-05 (no env-specific code paths).

## Principle

Dev and prod use **the same Terraform modules**, **the same Docker
images**, and **the same Cloud Function source**. The only differences
are input variables.

## GCP Project Layout

| Project | Purpose | Log level |
|---------|---------|-----------|
| `glens-dev` | Development + staging | `debug` |
| `glens-prod` | Production | `info` |

Both projects have identical IAM roles, APIs enabled, and VPC config.

## Terraform Variables (only differences)

```hcl
# environments/dev.tfvars
project       = "glens-dev"
log_level     = "debug"
min_instances = 0
max_instances = 1
alert_channel = "dev-slack"

# environments/prod.tfvars
project       = "glens-prod"
log_level     = "info"
min_instances = 1
max_instances = 10
alert_channel = "prod-pagerduty"
```

All resource definitions are shared via modules. No `if env == prod`.

## GCP Account Setup

1. Create projects `glens-dev` + `glens-prod`; enable same APIs
2. Enable: Cloud Run, Functions, Firestore, Pub/Sub, Secret Manager,
   IAM, Artifact Registry, Cloud Trace
3. WIF pool + GitHub OIDC in each; SA: `deployer@glens-{dev,prod}.iam`
4. Roles: Run/Functions/Firestore/SecretMgr/Storage/PubSub Admin
5. Artifact Registry repo + Terraform state GCS bucket in each

## CI Deploy (matrix selects env)

```yaml
deploy:
  strategy: { matrix: { env: [dev, prod] } }
  steps:
    - uses: google-github-actions/auth@v2
      with:
        workload_identity_provider: ${{ secrets[format('WIF_{0}', matrix.env)] }}
        service_account: ${{ secrets[format('SA_{0}', matrix.env)] }}
    - run: cd infra && terraform apply -var-file=environments/${{ matrix.env }}.tfvars -auto-approve
```

## Parity Validation (EP-03)

CI diffs `terraform plan` between dev and prod; only allowed vars
(project, log_level, instances, alert) may differ. Any other diff
fails the pipeline.

## Runtime Config

```go
type Config struct {
    LogLevel string `env:"LOG_LEVEL" default:"info"`
    Project  string `env:"GCP_PROJECT" required:"true"`
}
```

No `if os.Getenv("ENV") == "production"` — behaviour controlled by
config values, not environment names.

## Steps

1. Create GCP projects with identical API enablement
2. Set up WIF + service accounts; create tfvars files
3. Add parity-check CI job; verify plan diff is expected only

## Success Criteria

- [ ] `terraform apply` succeeds in both dev and prod
- [ ] Only log level, instance count, and alert channel differ
- [ ] No env-specific code paths in application source
- [ ] Parity-check CI job blocks unexpected divergence
