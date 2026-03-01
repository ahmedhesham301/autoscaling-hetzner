# autoscaling-hetzner

## Overview
`autoscaling-hetzner` is a Go service that:

- exposes a Gin HTTP API for Hetzner metadata and provisioning,
- stores templates/groups/servers in PostgreSQL,
- provisions Hetzner servers immediately when a group is created,
- exposes Prometheus scrape targets from DB state,
- initializes Grafana client state used by alert provisioning helpers.

This README documents the repository exactly as it is currently implemented.

## Current State Snapshot
- API routes are defined in `main.go` and listen on Gin default `:8080`.
- `docker-compose.yaml` currently starts: `db`, `pg_admin`, `prometheus`, `alloy`, `grafana`.
- `server` and `mimir` services are present but commented out in compose.
- Alert creation helper exists (`services/alerts.go`) but is not called from request flow.
- No `_test.go` files exist.

## Architecture and Service Topology

### Compose Services (as committed)
| Service | Status in `docker-compose.yaml` | Port(s) | Purpose |
| --- | --- | --- | --- |
| `db` | enabled | `5432:5432` | PostgreSQL with init schema from `configs/schema.sql` |
| `pg_admin` | enabled | `81:80` | pgAdmin UI |
| `prometheus` | enabled | `9090:9090` | Prometheus with remote-write receiver flag |
| `alloy` | enabled | `12345:12345` | HTTP discovery + scrape + remote_write forwarder |
| `grafana` | enabled | `3000:3000` | Grafana with provisioned Prometheus datasource |
| `server` | commented out | `8080:8080` | Go API container (disabled unless manually uncommented) |
| `mimir` | commented out | none mapped | Example metrics backend config (disabled) |

### Runtime Flow
1. API starts, initializes DB pool, Hetzner client, and Grafana client.
2. `POST /templates` stores provisioning template data.
3. `POST /groups` stores group row, then immediately calls `ScaleUp` for `desiredSize`.
4. Each created Hetzner server is inserted into `servers`.
5. `GET /targets` returns `IP:9100` targets from `servers`.
6. Alloy discovery polls `http://server:8080/targets` and forwards scraped metrics to Prometheus remote write (`/api/v1/write`).
7. Grafana reads Prometheus via provisioned datasource.

## Prerequisites
- Go `1.25.6+` (per `go.mod`).
- Docker Engine + Docker Compose.
- Hetzner Cloud API token with permissions for reading metadata and creating servers.
- Reachable Grafana instance for API startup (because Grafana initialization runs at process start).

## Security Warning and Remediation
Do not commit real secrets.

Current security-sensitive facts in this repo:
- Local env files in this workspace currently contain a real Hetzner token (`.env`, `.env.compose`).
- DB credentials are hardcoded in code as `postgres:1234` (`database/db.go`).
- Grafana basic auth is hardcoded as `admin/admin` (`grafana/grafanaClient.go`).

Recommended immediate remediation:
1. Revoke/rotate any exposed Hetzner token from Hetzner Cloud Console.
2. Replace local env files with non-secret placeholders.
3. Move credentials to secure secret storage for real deployments.
4. Replace hardcoded DB/Grafana credentials in code before production use.

## Configuration

### Environment Variables
| Variable | Required | Used by | Description |
| --- | --- | --- | --- |
| `HKEY` | yes | Hetzner client init | Hetzner API token |
| `DATABASE_HOST` | yes | DB init | Postgres hostname only (no scheme/port) |
| `GRAFANA_HOST` | yes | Grafana init | Grafana host:port (no scheme), example `localhost:3000` |

Safe example:

```bash
export HKEY=REPLACE_WITH_HETZNER_TOKEN
export DATABASE_HOST=localhost
export GRAFANA_HOST=localhost:3000
```

### Static Configuration Files
- SQL schema: `configs/schema.sql`
- Alloy pipeline: `configs/alloy/main.alloy`
- Grafana datasource provisioning: `configs/grafana/datasource.yaml`
- Optional example Mimir config: `configs/mimir.yaml`
- Example cloud-init text: `cloud-config.yaml` (stored as escaped newline string)

## Run Modes

### Mode A: Compose dependencies + API on host (as-is friendly)
1. Start dependency services:

```bash
docker compose up -d db pg_admin prometheus alloy grafana
```

2. Run API on host:

```bash
export HKEY=REPLACE_WITH_HETZNER_TOKEN
export DATABASE_HOST=localhost
export GRAFANA_HOST=localhost:3000
go run .
```

3. Verify API:

```bash
curl -i http://localhost:8080/locations
```

Important caveat for this mode:
- Alloy is configured with `http://server:8080/targets`, which resolves only when API runs as compose service named `server`. If API runs on host, Alloy target discovery will not reach it unless config/networking is adjusted.

### Mode B: Full compose with API container (requires manual file edit)
`server` is currently commented out in `docker-compose.yaml`. Uncomment it to run API inside compose network so Alloy can reach `http://server:8080/targets`.

## API Reference
Base URL (host-run): `http://localhost:8080`

### GET /locations
Returns Hetzner locations (`[]Location` from hcloud client).

```bash
curl -s http://localhost:8080/locations
```

Representative shape:

```json
[
  {
    "id": 1,
    "name": "fsn1",
    "network_zone": "eu-central"
  }
]
```

### GET /images
Returns Hetzner images (`[]Image`).

```bash
curl -s http://localhost:8080/images
```

### GET /types
Returns Hetzner server types (`[]ServerType`).

```bash
curl -s http://localhost:8080/types
```

### GET /networks
Returns Hetzner networks (`[]Network`).

```bash
curl -s http://localhost:8080/networks
```

### GET /firewalls
Returns Hetzner firewalls (`[]Firewall`).

```bash
curl -s http://localhost:8080/firewalls
```

### GET /keys
Returns Hetzner SSH keys (`[]SSHKey`).

```bash
curl -s http://localhost:8080/keys
```

### POST /templates
Creates and persists a template row.

```bash
curl -i -X POST http://localhost:8080/templates \
  -H 'Content-Type: application/json' \
  -d '{
    "image_id": 45557056,
    "networks": [11952339],
    "SSH_keys": [987654],
    "publicIPv4": true,
    "publicIPv6": true,
    "firewalls": [123456],
    "cloudConfig": "#cloud-config\npackage_update: true\npackages:\n  - prometheus-node-exporter"
  }'
```

JSON fields (from `model.Template`):
- `image_id` (`int64`, required)
- `networks` (`[]int64`, required)
- `SSH_keys` (`[]int64`, optional)
- `publicIPv4` (`bool`, required)
- `publicIPv6` (`bool`, required)
- `firewalls` (`[]int64`, optional)
- `cloudConfig` (`string`, optional)

Notes:
- `publicIPv4` and `publicIPv6` are pointer-backed booleans with `binding:"required"`, so fields must be present in JSON.
- Success returns `200` with inserted template object (including `id`).

### POST /groups
Creates and persists a group row, then immediately provisions `desiredSize` servers.

```bash
curl -i -X POST http://localhost:8080/groups \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "web",
    "templateId": 1,
    "zone": "eu-central",
    "locations": [1, 2],
    "serverType": "cx22",
    "minSize": 1,
    "desiredSize": 2,
    "maxSize": 5,
    "monitoringType": "cpu",
    "target": 80
  }'
```

JSON fields (from `model.Group`):
- `name` (`string`, required)
- `templateId` (`int`, required)
- `zone` (`string`, required)
- `locations` (`[]int64`, required)
- `serverType` (`string`, required)
- `minSize` (`int`, required)
- `desiredSize` (`int`, required)
- `maxSize` (`int`, required)
- `monitoringType` (`string`, required, DB enum `cpu|memory`)
- `target` (`int16`, required, DB check `1..100`)

Behavior:
- Group row is inserted before provisioning starts.
- If provisioning fails, request returns `400` but inserted group row is not rolled back.

### GET /targets
Returns Prometheus target groups from stored servers.

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
      "name": "web-Ab12C"
    }
  }
]
```

## Database Schema Mapping

### `templates`
- `id SERIAL PRIMARY KEY`
- `image_id BIGINT NOT NULL` <- `image_id`
- `networks BIGINT[] NOT NULL` <- `networks`
- `SSH_keys BIGINT[]` <- `SSH_keys`
- `public_ipv4 BOOL NOT NULL` <- `publicIPv4`
- `public_ipv6 BOOL NOT NULL` <- `publicIPv6`
- `firewalls BIGINT[]` <- `firewalls`
- `cloud_config VARCHAR` <- `cloudConfig`

### `groups`
- `id SERIAL PRIMARY KEY`
- `name VARCHAR NOT NULL`
- `template_id INTEGER NOT NULL REFERENCES templates(id)`
- `zone VARCHAR NOT NULL`
- `locations INTEGER[] NOT NULL`
- `server_type VARCHAR NOT NULL`
- `min_size SMALLINT NOT NULL`
- `desired_size SMALLINT NOT NULL`
- `max_size SMALLINT NOT NULL`
- `monitoring_type monitoring_types NOT NULL` (`cpu`, `memory`)
- `target SMALLINT NOT NULL CHECK (target BETWEEN 1 AND 100)`
- `rule_uid VARCHAR`

### `servers`
- `id SERIAL PRIMARY KEY`
- `name VARCHAR NOT NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- `group_id INTEGER NOT NULL REFERENCES groups(id)`
- `type VARCHAR NOT NULL`
- `location INTEGER NOT NULL`
- `private_ip INET NOT NULL`
- index on `group_id`

## Monitoring and Alerting (Implemented vs Wired)

Implemented:
- `GET /targets` provides scrape targets from DB.
- Alloy config scrapes discovered targets and forwards to Prometheus remote write.
- Grafana datasource provisioning points to `http://prometheus:9090`.
- `services.SetupAlert` can create Grafana alert rules using Prometheus expressions.

Not wired into current request flow:
- `services.SetupAlert` is never invoked by `POST /groups` or other active handlers.
- `groups.rule_uid` exists in schema but is not set by current API flow.

## Known Limitations and Behavioral Caveats
- Compose file does not run API by default (`server` is commented).
- Autoscaling loop is not implemented; creation path performs one-time `ScaleUp`.
- `minSize`/`maxSize` are stored but not enforced by runtime scaler logic.
- `zone` is stored in DB but not used in `ScaleUp`.
- `ScaleUp` assumes private network exists and reads `res.Server.PrivateNet[0].IP`.
- `ScaleUp` indexes locations with modulo; empty `locations` can break runtime assumptions.
- Grafana init assumes datasource list is non-empty and uses the first datasource UID.
- DB and Grafana credentials are hardcoded in code.
- API startup depends on DB and Grafana availability; failures panic at startup.

## Troubleshooting

### Panic: `DATABASE_HOST is not set`
Set env var before starting API:

```bash
export DATABASE_HOST=localhost
```

### Panic: `GRAFANA_HOST is not set`
Set env var before starting API:

```bash
export GRAFANA_HOST=localhost:3000
```

### API starts but `/groups` fails with Hetzner errors
- Validate `HKEY` token and project permissions.
- Ensure referenced IDs (`image_id`, `networks`, `SSH_keys`, `firewalls`, `locations`) exist in your Hetzner project.

### `/targets` returns empty array
No rows in `servers` table yet. Create at least one successful group provisioning first.

### Alloy target discovery does not work in host-run mode
Alloy uses static URL `http://server:8080/targets` from `configs/alloy/main.alloy`. That host resolves only when API is a compose service named `server`.

### Schema changes do not apply
`configs/schema.sql` is applied only on first Postgres initialization. If DB volume already exists, recreate container/volume for schema re-init.

## Development Notes
- Build binary:

```bash
go build .
```

- Run compiled binary:

```bash
./autoscaling-hetzner
```

- Container image build:

```bash
docker build -t autoscaling-hetzner:local .
```

- Security workflow file exists at `.github/workflows/snyk-security.yml`.
