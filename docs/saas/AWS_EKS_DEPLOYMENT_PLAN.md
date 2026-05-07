# AWS EKS Deployment Plan

> Plan for deploying Glens as a SaaS application on AWS using EKS
> with PostgreSQL/MongoDB and RabbitMQ/Kafka.

## Goal

Deploy the same Glens SaaS application to AWS accounts using **Amazon
EKS** with industry-standard, portable data and messaging services.
The stack uses **PostgreSQL** (relational) or **MongoDB** (document) for
persistence and **RabbitMQ** or **Kafka** for event streaming. These run
inside the cluster or as managed AWS services — the application code
stays cloud-agnostic.

## Guiding Principles

1. **Reuse** — same Go binaries, same OpenAPI spec, same Docker images.
2. **Kubernetes-native** — all services run as pods or Helm charts in EKS.
3. **Portable stack** — PostgreSQL, MongoDB, RabbitMQ, and Kafka work
   on any cloud or on-premises Kubernetes cluster.
4. **IaC** — Terraform with AWS provider (same workflow as GCP).
5. **No vendor lock-in in code** — data and messaging accessed via
   Go interfaces; swap implementations without changing handlers.
6. **Parity** — dev/prod differ only in scaling and log level.
7. **Security** — IRSA (IAM Roles for Service Accounts), no long-lived keys.

---

## GCP → AWS Service Mapping

| Concern | GCP Service | AWS + EKS Stack | Notes |
|---------|-------------|-----------------|-------|
| Container hosting | Cloud Run | **EKS** | Kubernetes pods |
| Container registry | Artifact Registry | **ECR** | Same Docker images |
| Relational database | Firestore | **PostgreSQL (RDS or in-cluster)** | Structured data |
| Document database | Firestore | **MongoDB (Atlas or in-cluster)** | Flexible schema |
| Object storage | Cloud Storage | **S3** | Reports, static assets |
| Secrets | Secret Manager | **Secrets Manager** | Credential storage |
| Analytics | BigQuery | **Athena + S3** | Query data in S3 |
| Messaging | Pub/Sub | **RabbitMQ (in-cluster or Amazon MQ)** | Simple queue/topic |
| Streaming | Pub/Sub | **Kafka (in-cluster or Amazon MSK)** | High-throughput |
| Event consumers | Cloud Functions | **Worker pods in EKS** | Same Go binary |
| Auth | Firebase Auth | **Cognito** | User pools + federation |
| Observability | Cloud Trace | **X-Ray** | Distributed tracing |
| Monitoring | Cloud Monitoring | **CloudWatch** | Metrics + alerts |
| Logging | Cloud Logging | **CloudWatch Logs** | Structured logs |
| DNS + TLS | Cloud DNS | **Route 53 + ACM** | Domain + cert |
| Load balancer | Cloud LB | **ALB (ingress)** | Application LB |
| CI/CD auth | Workload Identity | **OIDC + IRSA** | GitHub → AWS |
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
                      ├─ Deployment: glens-worker (event consumers)
                      ├─ StatefulSet: postgresql (or RDS external)
                      ├─ StatefulSet: rabbitmq (or Amazon MQ external)
                      ├─ Service: ClusterIP (api, db, mq)
                      └─ Ingress: ALB → Route 53
                           │
                      S3 · Secrets Manager · CloudWatch
```

---

## Database Choice: PostgreSQL vs MongoDB

Both are supported. Pick one based on data access patterns:

| Concern | PostgreSQL | MongoDB |
|---------|-----------|---------|
| Data model | Relational + JSONB | Document (BSON) |
| Schema | Strict migrations | Flexible, schema-less |
| Query | SQL + full-text search | MQL, aggregation pipeline |
| Transactions | Full ACID | Multi-document ACID |
| Managed AWS | **RDS for PostgreSQL** | **DocumentDB** or Atlas |
| In-cluster | Helm: `bitnami/postgresql` | Helm: `bitnami/mongodb` |
| Go driver | `pgx` (jackc/pgx) | `mongo-go-driver` |
| Best for | Structured data, joins, analytics | Flexible documents, rapid iteration |
| GCP equivalent | Cloud SQL / AlloyDB | Firestore / MongoDB Atlas |

### Recommendation

**PostgreSQL** as the primary database — structured data with JSONB
for semi-structured fields. Run results, user profiles, and workspace
configs fit naturally in relational tables. JSONB columns handle
variable payloads (test results, endpoint metadata).

**MongoDB** as an alternative when the document model from Firestore
is preferred, or when data shapes vary significantly between features.

Both options use the same Go interface — swap at startup via config.

---

## Messaging Choice: RabbitMQ vs Kafka

Both are supported. Pick one based on messaging patterns:

| Concern | RabbitMQ | Kafka |
|---------|----------|-------|
| Pattern | Message queue (push) | Event log (pull) |
| Delivery | At-most/at-least-once | Exactly-once (with config) |
| Ordering | Per-queue FIFO | Per-partition ordered |
| Replay | No (consumed = gone) | Yes (configurable retention) |
| Throughput | ~50K msg/s | ~1M msg/s |
| Managed AWS | **Amazon MQ** | **Amazon MSK** |
| In-cluster | Helm: `bitnami/rabbitmq` | Helm: `bitnami/kafka` |
| Go client | `rabbitmq/amqp091-go` | `segmentio/kafka-go` |
| Best for | Task queues, simple pub/sub | Event sourcing, high volume |

### Recommendation

**RabbitMQ** for the initial deployment — simpler to operate, lower
resource footprint, and sufficient for the current event catalogue
(report generation, issue creation, notifications).

**Kafka** when the platform scales to high-throughput event streaming,
event replay, or multi-consumer fan-out with ordering guarantees.

Both options use the same Go `EventPublisher` interface.

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

### R3 — Database (PostgreSQL)

| ID | Requirement | Details |
|----|-------------|---------|
| R3.1 | PostgreSQL 16+ instance | RDS or in-cluster StatefulSet |
| R3.2 | Tables: users, workspaces, runs, results | Relational schema with JSONB |
| R3.3 | Connection pooling (PgBouncer) | Efficient pod-to-DB connections |
| R3.4 | Automated backups | Point-in-time recovery |
| R3.5 | Schema migrations | `golang-migrate` or `goose` |
| R3.6 | Read replica (prod) | Offload read queries |

### R3-ALT — Database (MongoDB Alternative)

| ID | Requirement | Details |
|----|-------------|---------|
| R3-ALT.1 | MongoDB 7+ instance | DocumentDB, Atlas, or in-cluster |
| R3-ALT.2 | Collections: users, workspaces, runs | Same document model as Firestore |
| R3-ALT.3 | Indexes on workspace_id + created_at | Query performance |
| R3-ALT.4 | Replica set | High availability |
| R3-ALT.5 | Automated backups | Point-in-time recovery |

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
| R5.4 | DB credentials in Secrets Manager | RDS password rotation |

### R6 — Messaging (RabbitMQ)

| ID | Requirement | Details |
|----|-------------|---------|
| R6.1 | RabbitMQ instance | Amazon MQ or in-cluster Helm chart |
| R6.2 | Exchanges for domain events | Same event catalogue as GCP |
| R6.3 | Queues per consumer | Report gen, issue creator, etc. |
| R6.4 | Dead-letter exchange | Failed message retry |
| R6.5 | Management UI | Monitor queues and consumers |
| R6.6 | Worker deployment in EKS | Go consumer pods (not Lambda) |

### R6-ALT — Messaging (Kafka Alternative)

| ID | Requirement | Details |
|----|-------------|---------|
| R6-ALT.1 | Kafka cluster | Amazon MSK or in-cluster Helm chart |
| R6-ALT.2 | Topics for domain events | Same event catalogue as GCP |
| R6-ALT.3 | Consumer groups per service | Report gen, issue creator, etc. |
| R6-ALT.4 | Dead-letter topic | Failed message handling |
| R6-ALT.5 | Schema registry | Avro/JSON schema validation |
| R6-ALT.6 | Configurable retention | Event replay capability |

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
| R9.6 | DB migration step in pipeline | Run migrations before deploy |

### R10 — Networking & Security

| ID | Requirement | Details |
|----|-------------|---------|
| R10.1 | VPC with public/private subnets | EKS nodes in private subnets |
| R10.2 | Security groups | Restrict pod-to-DB, pod-to-MQ traffic |
| R10.3 | Network policies | Kubernetes-level isolation |
| R10.4 | Route 53 + ACM | Custom domain with TLS |
| R10.5 | WAF (optional) | Web Application Firewall on ALB |
| R10.6 | DB in private subnet only | No public access to PostgreSQL/MongoDB |
| R10.7 | MQ in private subnet only | No public access to RabbitMQ/Kafka |

---

## Terraform Module Layout

```text
infra/
├── aws/
│   ├── main.tf                # AWS provider, backend (S3 + DynamoDB lock)
│   ├── variables.tf           # region, cluster_name, db_engine, mq_engine
│   ├── outputs.tf             # cluster endpoint, ECR URL, ALB DNS, DB host
│   ├── modules/
│   │   ├── vpc/               # VPC, subnets, NAT gateway
│   │   ├── eks/               # EKS cluster, node groups, IRSA
│   │   ├── ecr/               # Container registry
│   │   ├── rds/               # RDS PostgreSQL (managed option)
│   │   ├── mongodb/           # DocumentDB or Atlas (managed option)
│   │   ├── s3/                # Buckets, lifecycle, CDN
│   │   ├── secrets/           # Secrets Manager (API creds + DB password)
│   │   ├── rabbitmq/          # Amazon MQ for RabbitMQ (managed option)
│   │   ├── kafka/             # Amazon MSK (managed option)
│   │   ├── cognito/           # User pool, app clients
│   │   ├── observability/     # X-Ray, CloudWatch, alarms
│   │   └── dns/               # Route 53, ACM certs
│   └── environments/
│       ├── dev.tfvars          # dev: small RDS, single-node RabbitMQ
│       └── prod.tfvars         # prod: multi-AZ RDS, RabbitMQ cluster
```

---

## Kubernetes Manifests

```text
k8s/
├── base/
│   ├── namespace.yaml
│   ├── api-deployment.yaml     # glens-api pods
│   ├── worker-deployment.yaml  # glens-worker pods (event consumers)
│   ├── service.yaml            # ClusterIP for api
│   ├── ingress.yaml            # ALB ingress
│   ├── hpa.yaml                # Horizontal Pod Autoscaler (api + worker)
│   ├── serviceaccount.yaml     # IRSA-annotated SA
│   ├── configmap.yaml          # DB_HOST, MQ_HOST, LOG_LEVEL
│   ├── db-secret.yaml          # ExternalSecret → Secrets Manager
│   └── mq-secret.yaml          # ExternalSecret → Secrets Manager
├── overlays/
│   ├── dev/
│   │   └── kustomization.yaml  # dev replica count, in-cluster DB/MQ
│   └── prod/
│       └── kustomization.yaml  # prod replicas, managed RDS/Amazon MQ
├── in-cluster/                  # Optional: run DB/MQ inside EKS
│   ├── postgresql.yaml          # bitnami/postgresql Helm values
│   ├── mongodb.yaml             # bitnami/mongodb Helm values
│   ├── rabbitmq.yaml            # bitnami/rabbitmq Helm values
│   └── kafka.yaml               # bitnami/kafka Helm values
```

---

## Database Schema (PostgreSQL)

```sql
-- users
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       TEXT UNIQUE NOT NULL,
    plan        TEXT NOT NULL DEFAULT 'free',
    settings    JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ DEFAULT now()
);

-- workspaces
CREATE TABLE workspaces (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    config      JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ DEFAULT now()
);

-- runs
CREATE TABLE runs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    spec_url    TEXT NOT NULL,
    models      TEXT[] NOT NULL,
    status      TEXT NOT NULL DEFAULT 'pending',
    summary     JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ DEFAULT now(),
    completed_at TIMESTAMPTZ
);
CREATE INDEX idx_runs_workspace ON runs(workspace_id, created_at DESC);

-- results (per endpoint per model)
CREATE TABLE results (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id      UUID REFERENCES runs(id) ON DELETE CASCADE,
    endpoint    TEXT NOT NULL,
    method      TEXT NOT NULL,
    model       TEXT NOT NULL,
    passed      BOOLEAN NOT NULL,
    duration_ms INT NOT NULL,
    output      JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_results_run ON results(run_id);
```

---

## Event Consumer Architecture

Instead of serverless functions (Lambda/Cloud Functions), event consumers
run as **long-lived worker pods** in EKS. This simplifies the stack —
one binary, one deployment model.

```text
glens-api (pod)
  │
  └── Publish event ──► RabbitMQ exchange (or Kafka topic)
                             │
            ┌────────────────┼────────────────┐
            ▼                ▼                ▼
      Queue: reports   Queue: issues   Queue: notifications
            │                │                │
            ▼                ▼                ▼
      glens-worker     glens-worker     glens-worker
      (report gen)     (issue create)   (notify)
```

Workers are the **same Go binary** started with a `--worker` flag
and a `--queue` argument:

```bash
glens --worker --queue=reports
glens --worker --queue=issues
glens --worker --queue=notifications
```

Each worker deployment in Kubernetes scales independently via HPA.

---

## Code Changes Required

The existing Go code is cloud-agnostic in its core logic. Only the
infrastructure integration points need alternatives:

| Component | Current (GCP) | AWS + EKS Change | Effort |
|-----------|---------------|------------------|--------|
| HTTP server | `net/http` | No change | — |
| Business logic | `internal/` | No change | — |
| Docker image | Dockerfile | No change | — |
| OpenAPI spec | `openapi.yaml` | No change | — |
| Firestore client | `cloud.google.com/go/firestore` | `pgx` or `mongo-go-driver` | Medium |
| Secret Manager | `cloud.google.com/go/secretmanager` | AWS SDK | Small |
| Pub/Sub publisher | `cloud.google.com/go/pubsub` | `streadway/amqp` or `kafka-go` | Small |
| OTel exporter | GCP Trace exporter | X-Ray exporter | Small |
| Auth (Firebase) | Firebase Admin SDK | Cognito JWT | Medium |

### Provider Interfaces

Abstract data and messaging behind Go interfaces so GCP, AWS, and
self-hosted implementations coexist:

```go
// internal/platform/store.go — database abstraction
type Store interface {
    GetUser(ctx context.Context, id string) (*User, error)
    CreateRun(ctx context.Context, run *Run) error
    ListRuns(ctx context.Context, workspaceID string) ([]*Run, error)
    SaveResult(ctx context.Context, result *Result) error
}

// internal/platform/gcp/firestore.go   — Firestore implementation
// internal/platform/aws/postgres.go    — PostgreSQL implementation
// internal/platform/aws/mongodb.go     — MongoDB implementation
```

```go
// internal/platform/events.go — messaging abstraction
type EventPublisher interface {
    Publish(ctx context.Context, topic string, event Event) error
}

type EventConsumer interface {
    Subscribe(ctx context.Context, queue string, handler EventHandler) error
}

// internal/platform/gcp/pubsub.go      — Pub/Sub implementation
// internal/platform/aws/rabbitmq.go    — RabbitMQ implementation
// internal/platform/aws/kafka.go       — Kafka implementation
```

Selected at startup via config:

```go
switch cfg.Platform {
case "gcp":
    store = gcp.NewFirestoreStore(ctx, cfg.GCPProject)
    publisher = gcp.NewPubSubPublisher(ctx, cfg.GCPProject)
case "aws-pg":
    store = aws.NewPostgresStore(ctx, cfg.DatabaseURL)
    publisher = aws.NewRabbitMQPublisher(cfg.AMQPURL)
case "aws-mongo":
    store = aws.NewMongoStore(ctx, cfg.MongoURL)
    publisher = aws.NewKafkaPublisher(cfg.KafkaBrokers)
}
```

---

## Implementation Phases

| Phase | Description | Depends on |
|-------|-------------|------------|
| A1 | VPC + EKS cluster (Terraform) | — |
| A2 | ECR + Docker push workflow | A1 |
| A3 | Kubernetes manifests (api + worker) | A1, A2 |
| A4 | PostgreSQL (RDS or in-cluster) + Go `Store` interface | A1 |
| A5 | Schema migrations + seed data | A4 |
| A6 | RabbitMQ (Amazon MQ or in-cluster) + Go publisher/consumer | A1 |
| A7 | Worker deployments (report gen, issue creator) | A3, A6 |
| A8 | Secrets Manager + Go interface | A1 |
| A9 | Cognito user pool + auth middleware | A1 |
| A10 | Observability (X-Ray + CloudWatch) | A3 |
| A11 | Route 53 + ACM + ALB ingress | A3 |
| A12 | CI/CD pipeline (GitHub OIDC → AWS) | A1, A2 |

### MongoDB + Kafka Alternative Path

| Phase | Description | Depends on |
|-------|-------------|------------|
| A4-ALT | MongoDB (DocumentDB or in-cluster) + Go `Store` interface | A1 |
| A6-ALT | Kafka (MSK or in-cluster) + Go publisher/consumer | A1 |

---

## Cost Estimate

### Managed Services (RDS + Amazon MQ)

| Service | Estimated monthly cost | Notes |
|---------|----------------------|-------|
| EKS | ~$73 | Cluster fee ($0.10/hr) |
| EC2 (node group) | ~$30–70 | t3.medium (dev), t3.large (prod) |
| RDS PostgreSQL | ~$15–30 | db.t4g.micro (dev), db.t4g.small (prod) |
| Amazon MQ (RabbitMQ) | ~$22 | mq.t3.micro single-instance |
| ECR | < $1 | Image storage |
| S3 | < $1 | Reports storage |
| Secrets Manager | ~$2 | $0.40/secret/mo |
| CloudWatch | < $5 | Metrics + logs |
| Route 53 | ~$1 | Hosted zone |
| **Total (dev)** | **~$150** | Managed services |
| **Total (prod)** | **~$250** | Multi-AZ, larger instances |

### In-Cluster (Helm Charts — Lower Cost)

| Service | Estimated monthly cost | Notes |
|---------|----------------------|-------|
| EKS | ~$73 | Cluster fee |
| EC2 (node group) | ~$50–100 | Larger nodes to run DB + MQ |
| PostgreSQL pod | $0 | Runs on existing nodes |
| RabbitMQ pod | $0 | Runs on existing nodes |
| **Total (dev)** | **~$125** | Everything in-cluster |

### Managed Kafka Alternative

| Service | Estimated monthly cost | Notes |
|---------|----------------------|-------|
| Amazon MSK | ~$75–150 | kafka.t3.small (min 2 brokers) |
| Amazon DocumentDB | ~$30 | db.t4g.medium |

**Note**: Kafka (MSK) is significantly more expensive than RabbitMQ
(Amazon MQ) for small workloads. Start with RabbitMQ; migrate to Kafka
when throughput exceeds ~10K events/second.

---

## Local Development

```yaml
# docker-compose.aws-local.yml
services:
  postgres:
    image: postgres:16-alpine
    ports: ["5432:5432"]
    environment:
      POSTGRES_DB: glens
      POSTGRES_USER: glens
      POSTGRES_PASSWORD: localdev
    volumes: ["pgdata:/var/lib/postgresql/data"]

  rabbitmq:
    image: rabbitmq:4-management-alpine
    ports: ["5672:5672", "15672:15672"]  # AMQP + management UI
    environment:
      RABBITMQ_DEFAULT_USER: glens
      RABBITMQ_DEFAULT_PASS: localdev

  # Optional: MongoDB alternative
  mongodb:
    image: mongo:7
    ports: ["27017:27017"]
    volumes: ["mongodata:/data/db"]

  # Optional: Kafka alternative
  kafka:
    image: bitnami/kafka:3.7
    ports: ["9092:9092"]
    environment:
      KAFKA_CFG_NODE_ID: 0
      KAFKA_CFG_PROCESS_ROLES: controller,broker
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 0@kafka:9093
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER

volumes:
  pgdata:
  mongodata:
```

Run locally:

```bash
docker compose -f docker-compose.aws-local.yml up -d
cd cmd/api && DATABASE_URL="postgres://glens:localdev@localhost:5432/glens" \
              AMQP_URL="amqp://glens:localdev@localhost:5672/" \
              go run .
```

---

## Success Criteria

- [ ] EKS cluster running with `dev` and `prod` namespaces
- [ ] Glens API pods healthy, reachable via ALB
- [ ] PostgreSQL stores and retrieves analysis results
- [ ] Schema migrations run automatically in CI/CD pipeline
- [ ] RabbitMQ routes domain events to worker pods
- [ ] Worker pods process events (reports, issues, notifications)
- [ ] Secrets Manager resolves credential references
- [ ] GitHub Actions deploys to AWS via OIDC (no access keys)
- [ ] Route 53 + ACM provides custom domain with TLS
- [ ] CloudWatch shows metrics, logs, and X-Ray traces
- [ ] Same API binary runs as both HTTP server and event worker
- [ ] `docker-compose.aws-local.yml` starts full local dev stack
