# Step 6 — Map a Friendly Domain Name

Cloud Run services get auto-generated URLs like
`https://glens-api-<hash>-uc.a.run.app`.
This step maps `api.glens.example.com` (or any domain you own) to that service.

## 6.1 Choose your domain

Decide on a subdomain, for example `api.glens.example.com`.
You need control over the DNS records for `example.com`.

## 6.2 Verify domain ownership with Google

```bash
gcloud domains verify example.com
```

Follow the prompt — it asks you to add a `TXT` record to your DNS provider.
Wait a few minutes, then run again to confirm:

```bash
gcloud domains verify example.com --verify
```

> You only need to verify the root domain once; all subdomains are included.

## 6.3 Map the domain to the Cloud Run service

```bash
PROJECT=glens-prod     # use the project that hosts the prod service
REGION=us-central1

gcloud run domain-mappings create \
  --service=glens-api \
  --domain=api.glens.example.com \
  --region="$REGION" \
  --project="$PROJECT"
```

## 6.4 Configure DNS

After creating the mapping, GCP shows the required DNS records:

```bash
gcloud run domain-mappings describe \
  --domain=api.glens.example.com \
  --region="$REGION" --project="$PROJECT"
```

Add the records shown in the output to your DNS provider.
Typical records:

| Type | Name | Value |
|------|------|-------|
| `CNAME` | `api.glens` | `ghs.googlehosted.com.` |

> If your DNS provider does not support `CNAME` at the zone apex, use an `A` / `AAAA`
> record instead — GCP shows the IPs in the mapping description.

## 6.5 Wait for TLS provisioning

GCP automatically provisions a free managed TLS certificate.
This normally takes 5–15 minutes, but DNS propagation can take up to 24 hours.

Check the status:

```bash
gcloud run domain-mappings describe \
  --domain=api.glens.example.com \
  --region="$REGION" --project="$PROJECT" \
  --format="value(status.conditions)"
```

Look for `CertificateProvisioned: True`.

## 6.6 Verify end-to-end

Once DNS has propagated:

```bash
curl https://api.glens.example.com/healthz
# Expected: {"status":"ok","version":"..."}
```

## 6.7 Update the dev environment too (optional)

Repeat Steps 6.3–6.5 using:

- domain: `api-dev.glens.example.com`
- project: `glens-dev`

---

**Done!** Your Glens SaaS service is live at `https://api.glens.example.com`.

Return to the [overview](DEPLOY_GCP.md) for a summary of what was deployed.
