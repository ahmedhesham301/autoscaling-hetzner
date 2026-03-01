# autoscaling-hetzner

## What This Project Does
This service provisions Hetzner Cloud servers from saved templates and groups, stores infrastructure metadata in PostgreSQL, exposes Prometheus scrape targets from the database, and initializes Grafana API access for alert-related operations. The runtime is a Gin HTTP API plus a monitoring stack (Alloy, Prometheus, Grafana) wired through Docker Compose.

## Architecture
Services and default ports from `docker-compose.yaml`:

- `server:8080` (this Go API)
- `db:5432` (PostgreSQL)
- `pg_admin:81` (pgAdmin UI)
- `prometheus:9090`
- `alloy:12345`
- `grafana:3000`

High-level flow:

1. `PUT /groups` stores group configuration and triggers server creation on Hetzner.
2. Created servers are persisted in the `servers` table.
3. `GET /targets` returns `IP:9100` endpoints from DB rows.
4. Alloy `discovery.http` pulls `http://server:8080/targets`.
5. Alloy scrapes targets and forwards metrics to Prometheus remote write (`/api/v1/write`).
6. Grafana reads Prometheus data via provisioned datasource.

## Prerequisites
- Docker Engine with Docker Compose.
- A Hetzner Cloud API token with permissions to create servers and read locations/networks/images.
- A reachable Grafana host for backend calls via `GRAFANA_HOST` (host:port, no scheme).
- Awareness of hardcoded assumptions in current code:
- Postgres connection string uses user `postgres` and password `1234`.
- Grafana API client uses basic auth `admin/admin`.

## Required Environment Variables
Do not commit real secrets.

| Variable | Purpose | Example |
| --- | --- | --- |
| `HKEY` | Hetzner API token used by the server client | `HKEY=hc_xxxxxxxxx` |
| `DATABASE_HOST` | PostgreSQL hostname reachable by the Go service | `DATABASE_HOST=db` (Compose) or `DATABASE_HOST=localhost` (host run) |
| `GRAFANA_HOST` | Grafana host and port (without `http://`) | `GRAFANA_HOST=grafana:3000` |

## Quick Start (Docker Compose)
1. Create/update `.env.compose`:

```bash
cat > .env.compose <<'EOF'
HKEY=REPLACE_WITH_HETZNER_TOKEN
DATABASE_HOST=db
GRAFANA_HOST=grafana:3000
EOF
```

2. Start the stack:

```bash
docker compose up --build
```

3. Verify first checks:

```bash
curl -i http://localhost:8080/locations
curl -i http://localhost:8080/networks
```

4. Open dashboards and tools:
- Grafana: `http://localhost:3000` (current default auth in code is `admin/admin`)
- Prometheus: `http://localhost:9090`
- pgAdmin: `http://localhost:81`

Startup note: `server` waits for a healthy `db` container. Other services start independently.

## API
Base URL: `http://localhost:8080`

### GET /locations
Purpose: Return Hetzner locations grouped by network zone.

Example:

```bash
curl -s http://localhost:8080/locations
```

Representative response:

```json
{
  "eu-central": {
    "fsn1": 12345,
    "nbg1": 12346
  }
}
```

### GET /images
Purpose: Return available image versions grouped by OS flavor.

Example:

```bash
curl -s http://localhost:8080/images
```

Representative response:

```json
{
  "ubuntu": ["22.04", "24.04"],
  "debian": ["12"]
}
```

### GET /types
Purpose: Intended to inspect available server types.

Example:

```bash
curl -i http://localhost:8080/types
```

Current behavior caveat: handler logs type details to server stdout and does not return a structured JSON payload.

### GET /networks
Purpose: Return Hetzner networks as `name -> id`.

Example:

```bash
curl -s http://localhost:8080/networks
```

Representative response:

```json
{
  "private-net": 11952339
}
```

### PUT /templates
Purpose: Save a template used when creating servers.

Example:

```bash
curl -i -X PUT http://localhost:8080/templates \
  -H 'Content-Type: application/json' \
  -d '{
    "OSFlavor": "ubuntu",
    "OSVersion": "24.04",
    "cloudConfig": "#cloud-config\npackage_update: true\npackages:\n  - prometheus-node-exporter"
  }'
```

Success response: `200 OK` with empty body.

### PUT /groups
Purpose: Save a group definition, then create `desiredSize` servers immediately.

Example:

```bash
curl -i -X PUT http://localhost:8080/groups \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "web",
    "templateId": 1,
    "zone": "eu-central",
    "locations": ["nbg1", "fsn1"],
    "serverTypes": ["cx22", "cx32"],
    "minSize": 1,
    "desiredSize": 2,
    "maxSize": 5,
    "networks": ["private-net"],
    "monitoringType": "cpu",
    "target": 80
  }'
```

Success response: `200 OK` with empty body.  
Validation/service errors return `400` with `{"error":"..."}`.

### GET /targets
Purpose: Return Prometheus target list built from DB `servers`.

Example:

```bash
curl -s http://localhost:8080/targets
```

Representative response:

```json
[
  {
    "targets": ["10.0.0.12:9100"],
    "labels": {
      "groupId": "1",
      "name": "web-a1B2c"
    }
  }
]
```

## Database Schema Overview
- `templates`: stores OS flavor/version and optional cloud-init payload.
- `groups`: stores autoscaling configuration (`zone`, `locations`, `server_types`, bounds, monitoring fields).
- `servers`: stores provisioned server metadata and private IP used for scraping.

Important schema details:

- `monitoring_types` enum values: `cpu`, `memory`.
- `groups.rule_uid` column exists in schema.

## Monitoring and Alerting Flow
- Scrape targets are generated from `servers` via `GET /targets`.
- Alloy config (`configs/alloy/main.alloy`) uses `discovery.http` against `http://server:8080/targets`.
- Alloy forwards scraped metrics to Prometheus remote write endpoint (`http://prometheus:9090/api/v1/write`).
- On startup, Grafana client initialization fetches datasource UID and folder UID for alert operations.
- The alert creation helper exists (`services.SetupAlert`) but is not invoked by the current request flow.

## Known Limitations
- `GET /types` does not return a structured API response.
- Group scaling currently uses only the first `serverTypes` entry.
- Group scaling currently uses a fixed network ID in code (`11952339`) instead of the `networks` payload values.
- Grafana defaults are hardcoded in code (for example basic auth `admin/admin`, first datasource selection, folder ID assumptions).
- `cloud-config.yaml` currently contains escaped newline sequences (`\n`) instead of multiline YAML formatting.
- No automated tests (`*_test.go`) are present in this repository.

## Local Development (Optional)
Run the binary locally while keeping dependencies available.

1. Start dependency services (example):

```bash
docker compose up -d db grafana prometheus alloy
```

2. Export environment variables for host execution:

```bash
export HKEY=REPLACE_WITH_HETZNER_TOKEN
export DATABASE_HOST=localhost
export GRAFANA_HOST=localhost:3000
```

3. Run the service:

```bash
go run .
```

This mode still requires reachable Postgres, Grafana, and Hetzner APIs.

## Troubleshooting
- Missing env vars:
- Symptom: startup panic like `DATABASE_HOST is not set`.
- Fix: set `HKEY`, `DATABASE_HOST`, `GRAFANA_HOST` before running.

- Database host mismatch:
- Symptom: DB connection/ping failures at startup.
- Fix: use `DATABASE_HOST=db` inside Compose, or `DATABASE_HOST=localhost` for host-run with mapped port.

- Grafana host/auth mismatch:
- Symptom: startup panic during Grafana init.
- Fix: verify `GRAFANA_HOST` resolves correctly and that Grafana credentials still match current hardcoded values.

- Hetzner token/permission issues:
- Symptom: `/images`, `/networks`, `/locations`, or group creation returns provider errors.
- Fix: validate token scope and project access in Hetzner Cloud.

- Empty `/targets` response:
- Symptom: `[]` from `/targets`.
- Fix: ensure at least one successful `PUT /groups` created servers and rows were inserted into `servers`.
