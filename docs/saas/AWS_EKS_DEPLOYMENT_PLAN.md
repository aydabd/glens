# AWS EKS Deployment Plan

> Plan for deploying Glens as a SaaS application on AWS using EKS.
> Mirrors the GCP SaaS architecture with AWS-native services.

## Goal

Deploy the same Glens SaaS application to AWS accounts so that it can
run on **Amazon EKS** (Elastic Kubernetes Service). The plan reuses all
existing Go backend code and OpenAPI contracts — only the infrastructure
layer and deployment tooling change.

## Guiding Principles

1. **Reuse** — same Go binaries, same OpenAPI spec, same Docker images.
2. **Kubernetes-native** — EKS with managed node groups or Fargate.
3. **IaC** — Terraform with AWS provider (same workflow as GCP).
4. **No vendor lock-in in code** — cloud services accessed via interfaces.
5. **Parity** — dev/prod differ only in scaling and log level.
6. **Security** — IRSA (IAM Roles for Service Accounts), no long-lived keys.

---

## GCP → AWS Service Mapping

| Concern | GCP Service | AWS Equivalent | Notes |
|---------|-------------|----------------|-------|
| Container hosting | Cloud Run | **EKS (Fargate)** | Kubernetes pods |
| Container registry | Artifact Registry | **ECR** | Same Docker images |
| Database | Firestore | **DynamoDB** | Document model, scales to zero |
| Object storage | Cloud Storage | **S3** | Static assets, reports |
| Secrets | Secret Manager | **Secrets Manager** | Credential storage |
| Analytics | BigQuery | **Athena + S3** | Query data in S3 |
| Events | Pub/Sub | **SNS + SQS** | Topic + queue pattern |
| Serverless functions | Cloud Functions | **Lambda** | Event consumers |
| Auth | Firebase Auth | **Cognito** | User pools + federation |
| Observability | Cloud Trace | **X-Ray** | Distributed tracing |
| Monitoring | Cloud Monitoring | **CloudWatch** | Metrics + alerts |
| Logging | Cloud Logging | **CloudWatch Logs** | Structured logs |
| DNS + TLS | Cloud DNS + managed cert | **Route 53 + ACM** | Domain + cert |
| Load balancer | Cloud LB | **ALB (ingress)** | Application LB |
| CI/CD auth | Workload Identity (OIDC) | **OIDC + IRSA** | GitHub → AWS |
| API gateway | API Gateway (GCP) | **API Gateway (AWS)** | Optional |
| IaC state | GCS bucket | **S3 + DynamoDB lock** | Terraform backend |

---

## Architecture

```text
GitHub Actions (push to main)
  │
  ├─ OIDC token ──► AWS IAM (GitHub OIDC provider) ──► IAM Role
  │
  ├─ docker build & push ──► ECR (private registry)
  │
  └─ terraform apply / kubectl apply
           │
           └─► EKS Cluster
                 ├─ Namespace: glens-dev
                 └─ Namespace: glens-prod
                      │
                      ├─ Deployment: glens-api (pods)
                      ├─ Service: ClusterIP
                      └─ Ingress: ALB → Route 53
                           │
                      DynamoDB · Secrets Manager · S3
                           │
                      SNS/SQS → Lambda (event consumers)
```

---

## Requirements

### R1 — EKS Cluster Setup

| ID | Requirement | Details |
|----|-------------|---------|
| R1.1 | EKS cluster with managed node group or Fargate | Cost-effective, auto-scaling |
| R1.2 | Namespaces for `dev` and `prod` | Environment isolation |
| R1.3 | IRSA for pod-level IAM | Pods assume IAM roles, no keys |
| R1.4 | ALB Ingress Controller | Route traffic to pods |
| R1.5 | Cluster autoscaler or Karpenter | Scale nodes to demand |

### R2 — Container Registry (ECR)

| ID | Requirement | Details |
|----|-------------|---------|
| R2.1 | Private ECR repository | `glens/api` image |
| R2.2 | Lifecycle policy | Expire untagged images after 30 days |
| R2.3 | Image scanning | Vulnerability scan on push |

### R3 — Database (DynamoDB)

| ID | Requirement | Details |
|----|-------------|---------|
| R3.1 | Tables for users, workspaces, runs | Same data model as Firestore |
| R3.2 | On-demand capacity | Scales to zero cost when idle |
| R3.3 | Point-in-time recovery | Data protection |
| R3.4 | DynamoDB Streams | Change data capture for events |

### R4 — Storage (S3)

| ID | Requirement | Details |
|----|-------------|---------|
| R4.1 | Bucket for reports and static assets | Versioned, encrypted |
| R4.2 | Lifecycle rules | Archive old reports after 90 days |
| R4.3 | CloudFront distribution | CDN for static assets (optional) |

### R5 — Secrets (Secrets Manager)

| ID | Requirement | Details |
|----|-------------|---------|
| R5.1 | Store API credentials | Same ref-based pattern as GCP |
| R5.2 | IAM policy scoped to pods | Only API pods can read secrets |
| R5.3 | Automatic rotation | Optional for managed credentials |

### R6 — Events (SNS + SQS + Lambda)

| ID | Requirement | Details |
|----|-------------|---------|
| R6.1 | SNS topics for domain events | Same event catalogue as GCP |
| R6.2 | SQS queues per consumer | Report gen, issue creator, etc. |
| R6.3 | Lambda functions for consumers | Same Go code, different trigger |
| R6.4 | Dead-letter queue | Failed event retry |

### R7 — Authentication (Cognito)

| ID | Requirement | Details |
|----|-------------|---------|
| R7.1 | Cognito User Pool | Email + social login |
| R7.2 | JWT verification in API | Same middleware pattern |
| R7.3 | OIDC federation | Enterprise SSO support |
| R7.4 | API key management | Application-level keys |

### R8 — Observability

| ID | Requirement | Details |
|----|-------------|---------|
| R8.1 | OpenTelemetry → X-Ray | OTel SDK with AWS exporter |
| R8.2 | CloudWatch metrics | Custom metrics from API |
| R8.3 | CloudWatch Logs | Structured JSON logs |
| R8.4 | CloudWatch alarms | Error rate, latency alerts |

### R9 — CI/CD Pipeline

| ID | Requirement | Details |
|----|-------------|---------|
| R9.1 | GitHub OIDC → AWS IAM | No long-lived access keys |
| R9.2 | ECR push workflow | Build, tag, push on merge |
| R9.3 | EKS deploy workflow | `kubectl apply` or Helm |
| R9.4 | Terraform workflow | Plan on PR, apply on merge |
| R9.5 | Namespace-based environments | `dev` on RC, `prod` on stable |

### R10 — Networking & Security

| ID | Requirement | Details |
|----|-------------|---------|
| R10.1 | VPC with public/private subnets | EKS nodes in private subnets |
| R10.2 | Security groups | Restrict pod-to-service traffic |
| R10.3 | Network policies | Kubernetes-level isolation |
| R10.4 | Route 53 + ACM | Custom domain with TLS |
| R10.5 | WAF (optional) | Web Application Firewall on ALB |

---

## Terraform Module Layout

```text
infra/
├── aws/
│   ├── main.tf                # AWS provider, backend (S3 + DynamoDB lock)
│   ├── variables.tf           # region, cluster_name, env vars
│   ├── outputs.tf             # cluster endpoint, ECR URL, ALB DNS
│   ├── modules/
│   │   ├── vpc/               # VPC, subnets, NAT gateway
│   │   ├── eks/               # EKS cluster, node groups, IRSA
│   │   ├── ecr/               # Container registry
│   │   ├── dynamodb/          # Tables, indexes, streams
│   │   ├── s3/                # Buckets, lifecycle, CDN
│   │   ├── secrets/           # Secrets Manager
│   │   ├── events/            # SNS topics, SQS queues, Lambda
│   │   ├── cognito/           # User pool, app clients
│   │   ├── observability/     # X-Ray, CloudWatch, alarms
│   │   └── dns/               # Route 53, ACM certs
│   └── environments/
│       ├── dev.tfvars          # dev config (small nodes, debug)
│       └── prod.tfvars         # prod config (larger nodes, info)
```

---

## Kubernetes Manifests

```text
k8s/
├── base/
│   ├── namespace.yaml
│   ├── deployment.yaml         # glens-api pods
│   ├── service.yaml            # ClusterIP
│   ├── ingress.yaml            # ALB ingress
│   ├── hpa.yaml                # Horizontal Pod Autoscaler
│   └── serviceaccount.yaml     # IRSA-annotated SA
├── overlays/
│   ├── dev/
│   │   └── kustomization.yaml  # dev replica count, log level
│   └── prod/
│       └── kustomization.yaml  # prod replica count, log level
```

---

## Code Changes Required

The existing Go code is cloud-agnostic in its core logic. Only the
infrastructure integration points need AWS alternatives:

| Component | Current (GCP) | AWS Change | Effort |
|-----------|---------------|------------|--------|
| HTTP server | `net/http` | No change | — |
| Business logic | `internal/` | No change | — |
| Docker image | Dockerfile | No change | — |
| OpenAPI spec | `openapi.yaml` | No change | — |
| Firestore client | `cloud.google.com/go/firestore` | DynamoDB SDK | Medium |
| Secret Manager | `cloud.google.com/go/secretmanager` | AWS SDK | Small |
| Pub/Sub publisher | `cloud.google.com/go/pubsub` | SNS SDK | Small |
| OTel exporter | GCP Trace exporter | X-Ray exporter | Small |
| Auth (Firebase) | Firebase Admin SDK | Cognito JWT | Medium |

### Recommended Approach — Provider Interface

Abstract cloud services behind Go interfaces so both GCP and AWS
implementations can coexist:

```go
// internal/cloud/storage.go
type DocumentStore interface {
    Get(ctx context.Context, collection, id string) (map[string]any, error)
    Put(ctx context.Context, collection, id string, doc map[string]any) error
    Delete(ctx context.Context, collection, id string) error
    Query(ctx context.Context, q Query) ([]map[string]any, error)
}

// internal/cloud/gcp/firestore.go  — GCP implementation
// internal/cloud/aws/dynamodb.go   — AWS implementation
```

Selected at startup via config:

```go
var store cloud.DocumentStore
switch cfg.CloudProvider {
case "gcp":
    store = gcp.NewFirestoreStore(ctx, cfg.GCPProject)
case "aws":
    store = aws.NewDynamoDBStore(ctx, cfg.AWSRegion, cfg.DynamoTable)
}
```

---

## Implementation Phases

| Phase | Description | Depends on |
|-------|-------------|------------|
| A1 | VPC + EKS cluster (Terraform) | — |
| A2 | ECR + Docker push workflow | A1 |
| A3 | Kubernetes manifests + deploy | A1, A2 |
| A4 | DynamoDB tables + Go interface | A1 |
| A5 | Secrets Manager + Go interface | A1 |
| A6 | SNS/SQS + Lambda consumers | A1 |
| A7 | Cognito user pool + auth middleware | A1 |
| A8 | Observability (X-Ray + CloudWatch) | A3 |
| A9 | Route 53 + ACM + ALB ingress | A3 |
| A10 | CI/CD pipeline (GitHub OIDC → AWS) | A1, A2 |

---

## Cost Estimate (AWS Free Tier)

| Service | Free Tier | Expected usage |
|---------|-----------|----------------|
| EKS | $0.10/hr cluster fee | ~$73/mo (cluster always on) |
| Fargate | 750 hrs/mo (first 12 mo) | Pods scale to zero |
| ECR | 500 MB storage | < 500 MB |
| DynamoDB | 25 GB + 25 RCU/WCU | < 1 GB |
| S3 | 5 GB | < 1 GB |
| Secrets Manager | $0.40/secret/mo | ~$2/mo |
| Lambda | 1M requests/mo | < 100K |
| CloudWatch | 10 custom metrics | Sufficient |
| Route 53 | $0.50/hosted zone | $0.50/mo |

**Note**: EKS has a fixed cluster cost ($0.10/hr ≈ $73/mo) unlike
Cloud Run which scales to zero. For cost parity, consider ECS Fargate
as an alternative to EKS for small deployments.

### ECS Fargate Alternative

If the EKS cluster cost is too high for early-stage deployment, ECS
Fargate provides a simpler, lower-cost option:

| Concern | EKS | ECS Fargate |
|---------|-----|-------------|
| Cluster cost | $73/mo fixed | $0 (no cluster fee) |
| Complexity | High (K8s) | Low (task definitions) |
| Scaling | HPA + Karpenter | Auto-scaling built in |
| Portability | Standard K8s | AWS-specific |
| Best for | Large, multi-service | Single-service, cost-sensitive |

The Terraform modules and CI/CD workflows in this plan work with both
EKS and ECS Fargate. The choice can be made per environment.

---

## Success Criteria

- [ ] EKS cluster running with `dev` and `prod` namespaces
- [ ] Glens API pods healthy, reachable via ALB
- [ ] DynamoDB stores and retrieves analysis results
- [ ] Secrets Manager resolves credential references
- [ ] SNS/SQS events trigger Lambda consumers
- [ ] GitHub Actions deploys to AWS via OIDC (no access keys)
- [ ] Route 53 + ACM provides custom domain with TLS
- [ ] CloudWatch shows metrics, logs, and X-Ray traces
- [ ] Same Docker image runs on both GCP and AWS
