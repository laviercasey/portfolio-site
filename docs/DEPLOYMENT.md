# portfolio-site — deployment runbook

Production stack: **Caddy (shared reverse-proxy) → Next.js (frontend) + Go (backend) + Postgres**, deployed to `/opt/apps/portfolio/` on VPS, built & pushed from GitHub Actions to GHCR.

This runbook mirrors the pattern used by `TicketHub` and `med-reminder-bot`. See `deploy-plan.md` at the repo-root of the parent workspace for the shared Caddy architecture.

---

## 1. One-time VPS setup

Assumes Docker + Docker Compose v2 are already installed (done when Caddy / TicketHub / Med-Reminder were set up).

```bash
# 1. Create app + backup directories
sudo mkdir -p /opt/apps/portfolio
sudo mkdir -p /opt/backups/portfolio
sudo chown -R $USER:$USER /opt/apps/portfolio /opt/backups/portfolio

# 2. Ensure shared network exists (idempotent)
docker network create caddy-net 2>/dev/null || true

# 3. Create .env (use .env.example as template)
nano /opt/apps/portfolio/.env
# fill DB_PASSWORD, JWT_SECRET, ADMIN_PASSWORD_HASH, NEXT_PUBLIC_SITE_URL, IMAGE_TAG=latest

# 4. GHCR login (one-time, for manual pulls / rescue)
echo "$GHCR_PAT" | docker login ghcr.io -u laviercasey --password-stdin
# PAT needs scope: read:packages
```

Required GitHub secrets (repo → Settings → Secrets and variables → Actions):

| Secret | Purpose |
|---|---|
| `DEPLOY_HOST` | VPS IP or hostname |
| `DEPLOY_USER` | SSH user (not root) |
| `DEPLOY_SSH_KEY` | Private key matching `~/.ssh/authorized_keys` on VPS |

Required GitHub variable (same page, Variables tab):

| Variable | Value |
|---|---|
| `PRODUCTION_URL` | `https://lavier.tech` (shown in the Environment tile) |

Create the `production` environment (Settings → Environments → New → `production`). Optionally set:
- Required reviewers (human approval before every deploy)
- Deployment branches: only `main`

---

## 2. Caddy block to add

Append to `/opt/apps/caddy/Caddyfile`, then `docker compose exec caddy caddy reload --config /etc/caddy/Caddyfile`.

Apex + www, with the `/uploads/` path served as static files from the shared `portfolio_uploads` volume for performance:

```caddyfile
lavier.tech, www.lavier.tech {
    import security-headers
    encode gzip zstd

    log {
        output file /data/logs/portfolio.log {
            roll_size 10mb
            roll_keep 5
        }
    }

    @uploads path /uploads/*
    handle @uploads {
        root * /srv/portfolio
        file_server {
            precompressed gzip
        }
        header Cache-Control "public, max-age=2592000, immutable"
        header Content-Disposition "attachment"
        header X-Content-Type-Options "nosniff"
    }

    handle /healthz {
        reverse_proxy portfolio-backend:8080
    }

    handle {
        reverse_proxy portfolio-frontend:3000 {
            header_up X-Real-IP {remote_host}
            header_up X-Forwarded-Proto {scheme}
            health_uri /
            health_interval 30s
            health_timeout 5s
        }
    }
}
```

For Caddy to read the uploads volume, add to `/opt/apps/caddy/docker-compose.yml` → `caddy.volumes`:

```yaml
      - portfolio_uploads:/srv/portfolio/uploads:ro
```

And declare it as external at the bottom of the same file:

```yaml
volumes:
  caddy_data:
  caddy_config:
  caddy_logs:
  portfolio_uploads:
    external: true
    name: portfolio_uploads
```

The volume is created when portfolio-site first comes up (`docker-compose.prod.yml` declares `name: portfolio_uploads`).

---

## 3. First deploy

DNS: `lavier.tech` and `www.lavier.tech` A-records → VPS IP.

```bash
# On the VPS (after .env is filled):
cd /opt/apps/portfolio

# Pull compose file from the first GitHub Actions run (it scp-s it),
# OR bootstrap manually this first time:
curl -L -o docker-compose.yml https://raw.githubusercontent.com/laviercasey/portfolio-site/main/docker-compose.prod.yml

# Trigger the first CI/CD run by pushing to main:
#   - build & push images to GHCR
#   - scp docker-compose.prod.yml
#   - pull images, run migrations, start stack, health-gate
```

After the first successful deploy, `.last-good-sha` is created — **this is what unlocks auto-rollback for every subsequent deploy**.

---

## 4. Day-to-day operations

```bash
# Status
docker compose -f /opt/apps/portfolio/docker-compose.yml ps

# Logs
docker compose -f /opt/apps/portfolio/docker-compose.yml logs -f --tail 100 backend frontend

# Restart a single service
docker compose -f /opt/apps/portfolio/docker-compose.yml restart frontend

# Ad-hoc migration
cd /opt/apps/portfolio && docker compose --profile migrate run --rm migrate

# Health
curl -fsS https://lavier.tech/healthz
```

---

## 5. Rollback

### 5.1 Automatic (during deploy)

The `deploy` job in `.github/workflows/deploy.yml`:

1. Pulls new `backend:<sha>` + `frontend:<sha>`, writes `IMAGE_TAG=<sha>` into `.env`.
2. Runs migrations.
3. Waits up to 120s polling `/healthz` and `/` inside the containers.
4. On success → writes `<sha>` to `.last-good-sha` and prunes dangling images.
5. On failure → reads previous `.last-good-sha`, pulls that tag, `compose up -d`, verifies containers are running.

### 5.2 Manual via GitHub UI

Actions → **Rollback** → Run workflow → optionally paste a 7-char SHA, else empty = use `.last-good-sha`.

### 5.3 Manual via SSH

```bash
cd /opt/apps/portfolio

# Roll back to last known good
PREV=$(cat .last-good-sha)
sed -i "s/^IMAGE_TAG=.*/IMAGE_TAG=$PREV/" .env
docker compose up -d --remove-orphans

# Or roll back to an arbitrary older tag
sed -i 's/^IMAGE_TAG=.*/IMAGE_TAG=a1b2c3d/' .env
docker pull ghcr.io/laviercasey/portfolio-site/backend:a1b2c3d
docker pull ghcr.io/laviercasey/portfolio-site/frontend:a1b2c3d
docker compose up -d --remove-orphans
```

### 5.4 Database rollback

If a migration broke prod, the pre-deploy backup in `/opt/backups/portfolio/` is the source of truth:

```bash
# Find most recent pre-deploy dump
ls -lt /opt/backups/portfolio/ | head

# Restore (destructive — current DB is wiped by --clean in the dump)
cd /opt/apps/portfolio
gunzip -c /opt/backups/portfolio/pre-deploy-YYYYMMDDTHHMMSSZ-SHA.sql.gz \
  | docker compose exec -T postgres psql -U portfolio -d portfolio
```

---

## 6. Backups

The deploy workflow auto-dumps Postgres **before** every deploy to `/opt/backups/portfolio/pre-deploy-*.sql.gz` and prunes files older than 30 days.

For nightly backups of DB **and** uploads, install the helper script as a cron:

```bash
# Copy scripts/prod-backup.sh from repo to /opt/apps/portfolio/scripts/
chmod +x /opt/apps/portfolio/scripts/prod-backup.sh

# Root crontab — nightly at 03:17 UTC
sudo crontab -e
# Add:
17 3 * * * /opt/apps/portfolio/scripts/prod-backup.sh >> /var/log/portfolio-backup.log 2>&1
```

---

## 7. Monitoring & alerts (minimal)

```bash
# Uptime via cron (root crontab)
*/5 * * * * curl -fsS https://lavier.tech/healthz > /dev/null || \
  echo "portfolio-site healthz failed at $(date -u)" | tee -a /var/log/portfolio-alerts.log
```

For richer monitoring, reuse the Uptime Kuma / Grafana setup recommended in `deploy-plan.md` §8.3.

---

## 8. Branch protection (repo settings)

Settings → Branches → `main`:

- Require pull request before merging
- Require status checks to pass: `CI / Backend — lint & vet`, `CI / Backend — test & race`, `CI / Frontend — typecheck & build`, `CI / Docker — build verify`
- Require branches to be up to date before merging
- Do not allow force pushes

---

## 9. First-time checklist

- [ ] Repo pushed to `github.com/laviercasey/portfolio-site`
- [ ] GitHub secrets: `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_SSH_KEY`
- [ ] GitHub variable: `PRODUCTION_URL`
- [ ] GitHub Environment `production` created
- [ ] DNS A-records for `lavier.tech` + `www.lavier.tech`
- [ ] `/opt/apps/portfolio/.env` filled on VPS (passwords + secrets)
- [ ] `/opt/apps/caddy/Caddyfile` updated + reloaded
- [ ] `portfolio_uploads` volume mounted into Caddy container
- [ ] First push to `main` triggers green deploy
- [ ] `.last-good-sha` file exists on VPS after first deploy
- [ ] Nightly backup cron installed
- [ ] Branch protection rules on `main`
