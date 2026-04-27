# Analytics (Umami) — Integration Guide

Umami is deployed as a **separate stack** outside this repository. The portfolio
site is a *consumer*: it sends pageviews to Umami via a browser script and
reads aggregated stats server-side for the admin dashboard.

One Umami instance can track multiple sites — each gets its own `website` entry
in the Umami UI with a unique UUID.

## Required environment variables

Add to the prod `.env` of this repo once Umami is reachable.

### Backend (server-side, never sent to browser)

| Variable | Example | Purpose |
|---|---|---|
| `UMAMI_API_URL` | `https://umami.lavier.tech` | Base URL the Go backend calls |
| `UMAMI_API_KEY` | `<long token>` | `Authorization: Bearer` for Umami API |
| `UMAMI_WEBSITE_ID` | `58828d09-...` | UUID of the Website entry in Umami UI |

### Frontend (public, shipped to browser)

| Variable | Example | Purpose |
|---|---|---|
| `NEXT_PUBLIC_UMAMI_SRC` | `https://umami.lavier.tech/script.js` | Tracker script URL |
| `NEXT_PUBLIC_UMAMI_WEBSITE_ID` | `58828d09-...` | Same UUID, as `data-website-id` |
| `NEXT_PUBLIC_UMAMI_ALLOWED_HOSTS` | `umami.lavier.tech` | Comma-separated host allowlist for the tracker script |

`localhost` is always in the allowlist implicitly. Override with a comma list to
add your production host, e.g. `umami.lavier.tech,analytics.example.com`.

## SSRF validator (backend)

The Go backend validates `UMAMI_API_URL` at startup. Rejected:
- non-`http(s)` schemes
- cloud metadata endpoints (`169.254.169.254`, `metadata.google.internal`, etc.)
- link-local IPs (`169.254.0.0/16`)
- loopback hosts (`localhost`, `127.0.0.1`, `::1`, `0.0.0.0`) unless `APP_ENV=dev`

Any other host is trusted — the value is operator-configured and only parsed at
boot, so an attacker cannot substitute a hostile URL at runtime.

## Linking a new site in the Umami UI

1. Log into your Umami admin (e.g. `https://umami.lavier.tech`)
2. *Websites → Add website* — name, domain
3. Copy the UUID from the website detail page
4. Set `UMAMI_WEBSITE_ID` and `NEXT_PUBLIC_UMAMI_WEBSITE_ID` in this project's
   prod `.env` to that UUID
5. Set `NEXT_PUBLIC_UMAMI_SRC` to `https://<your-umami-host>/script.js`
6. *Settings → API Keys → Create key* — copy the token to `UMAMI_API_KEY`
7. Restart backend and frontend:
   `docker compose -f docker-compose.prod.yml up -d --force-recreate backend frontend`
8. Backend log should show `analytics service: enabled url=...` at INFO level
9. Visit the site in a browser; check Umami UI for the pageview within 1 minute

## Local development

`docker-compose.dev.yml` ships an **optional** Umami pair under the
`analytics` Compose profile for offline work:

```
docker compose -f docker-compose.yml -f docker-compose.dev.yml \
  --profile analytics up -d umami umami-db
```

This starts Umami on `http://localhost:3001` with dev-only defaults for
`HASH_SALT` / `APP_SECRET`. Default admin login is `admin` / `umami` — change it
in *Settings → Profile* before use.

When the profile is not active, dev backend starts in *disabled* analytics
mode and the admin widget renders the "not configured" state.

## Disabling analytics

Leave `UMAMI_API_URL`, `UMAMI_API_KEY`, `UMAMI_WEBSITE_ID` empty. Backend logs
`analytics service: disabled (missing envs)` and all `/api/analytics/*` endpoints
return 503 `analytics_not_configured`. The admin widget shows an empty-state
card. The browser-side tracker is also skipped when `NEXT_PUBLIC_UMAMI_SRC` or
`NEXT_PUBLIC_UMAMI_WEBSITE_ID` is empty.

## Rotating the API key

1. In Umami UI: *Settings → API Keys → Revoke*, then create a new key
2. Replace `UMAMI_API_KEY` in the prod `.env`
3. `docker compose -f docker-compose.prod.yml up -d --force-recreate backend`
4. No downtime for the tracker script — that path uses the public UUID only.
