# Step 3 — GitHub Environments, Secrets & Variables

GitHub **Environments** gate deployments behind manual approval (prod) and
store per-environment secrets that are only exposed to the matching workflow job.

## 3.1 Create environments in GitHub

1. Open your repository on GitHub.
2. Go to **Settings → Environments → New environment**.
3. Create two environments: **`dev`** and **`prod`**.
4. For **`prod`**: enable **Required reviewers** and add yourself (or your team).
   This prevents accidental production deploys.

## 3.2 Secrets — add to each environment

For **both** `dev` and `prod` environments
(**Settings → Environments → \<env\> → Environment secrets → Add secret**):

| Secret name | Value | Description |
|-------------|-------|-------------|
| `WIF_PROVIDER` | `projects/…/providers/github-provider` | Full resource name from Step 2.6 |
| `SA_EMAIL` | `glens-deploy@glens-dev.iam.gserviceaccount.com` | Deployment service account (change project for prod) |

Use the **dev** values for the `dev` environment and **prod** values for the `prod` environment.

## 3.3 Variables — add to each environment

(**Settings → Environments → \<env\> → Environment variables → Add variable**)

| Variable name | dev value | prod value | Description |
|---------------|-----------|-----------|-------------|
| `GCP_PROJECT` | `glens-dev` | `glens-prod` | GCP project ID |
| `GCP_REGION` | `us-central1` | `us-central1` | Deployment region |

## 3.4 Repository-level secrets (optional)

These are shared across all environments and used by optional features:

| Secret name | Description |
|-------------|-------------|
| `GPG_PRIVATE_KEY` | ASCII-armoured GPG key for signing release checksums |
| `GPG_PASSPHRASE` | Passphrase for the GPG key |

Add at **Settings → Secrets and variables → Actions → New repository secret**.
You can skip these if you do not need signed releases.

## 3.5 How secrets flow into workflows

The existing workflow files already reference these names:

```yaml
# .github/workflows/infra.yml (existing)
environment: ${{ matrix.workspace }}   # 'dev' or 'prod'
# secrets.WIF_PROVIDER and secrets.SA_EMAIL are injected automatically
```

No workflow changes are needed — the environment name selects the correct secret set.

## 3.6 Verify

After adding secrets, navigate to each environment in GitHub Settings and confirm
all four items appear in the **Secrets** and **Variables** lists.

---

**Next:** [Step 4 — Provision GCP Infrastructure with Terraform](STEP4_TERRAFORM.md)
