# Autoscaling Hetzner

Autoscaling control plane for Hetzner Cloud instances.

## Overview

This project is a control plane for autoscaling Hetzner Cloud servers based on monitoring signals.

## Architecture

<img width="971" height="681" alt="hetzner auto scaling (1)" src="https://github.com/user-attachments/assets/5a46de66-61dd-46b7-b12b-4ea1943126f4" />

Components:

- `server` (Go/Gin): API + orchestration logic
- `PostgreSQL`: stores templates, groups, and servers
- `Alloy`: discovers scrape targets from `/targets` and scrapes node exporter
- `Prometheus`: stores metrics (remote_write receiver enabled)
- `Grafana`: evaluates alerts and sends webhook events back to the API

## Prerequisites

- Docker + Docker Compose
- Hetzner Cloud API token
- Outbound internet access from the control plane (to call Hetzner API)
- Node Exporter available on each provisioned server (`prometheus-node-exporter`)

## Quick Start (Docker Compose)

1. Clone the repository:

```bash
git clone https://github.com/ahmedhesham301/autoscaling-hetzner.git
cd autoscaling-hetzner
```

2. Create `.env.compose`:

```bash
HKEY=<your_hetzner_api_token>
DATABASE_HOST=db
GRAFANA_HOST=grafana:3000
ENV=prod
```

3. Start services:

```bash
docker compose up -d --build
```

### `ENV` mode

- `ENV=prod` (recommended on Hetzner): uses private IPs for scraping.
- `ENV=dev` (local testing): uses public IPv4 for scraping.

For production deployment, the control plane should be in the same Hetzner network as managed servers.

## Default Ports

| Service | Port | Notes |
| --- | --- | --- |
| API server | `8080` | Main API |
| Grafana | `3000` | Default login: `admin` / `admin` |
| Prometheus | `9090` | Metrics storage |
| PostgreSQL | `5432` | Default: `postgres` / `1234` |
| Alloy | `12345` | Alloy HTTP endpoint |
| pgAdmin | `81` | Optional DB UI (`docker-compose`) |

## API Usage

No frontend is included yet; interact via HTTP API.

Base URL (local): `http://localhost:8080`

### 1) Discovery endpoints

Use these to fetch IDs/names before creating resources:

- `GET /locations`
- `GET /images`
- `GET /types`
- `GET /networks`
- `GET /firewalls`
- `GET /keys`

### 2) Create a template

`POST /templates`

| Field         | Type     | Required | Notes                 |
| ------------- | -------- | -------- | --------------------- |
| `image_id`    | `int`    | yes      | from `GET /images`    |
| `networks`    | `[]int`  | yes      | from `GET /networks`  |
| `SSH_keys`    | `[]int`  | no       | from `GET /keys`      |
| `firewalls`   | `[]int`  | no       | from `GET /firewalls` |
| `publicIPv4`  | `bool`   | yes      | must be true in dev   |
| `publicIPv6`  | `bool`   | yes      |                       |
| `cloudConfig` | `string` | no       |                       |

Example:

```json
{
  "image_id": 310554929,
  "networks": [11952339],
  "SSH_keys": [107916411],
  "publicIPv4": true,
  "publicIPv6": true,
  "cloudConfig": "#cloud-config\npackage_update: true\npackage_upgrade: true\npackages:\n  - prometheus-node-exporter\n  - stress"
}
```

### 3) Create a group

`POST /groups`

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `templateId` | `int` | yes | from `GET /templates` |
| `name` | `string` | yes | prefix for created server names |
| `zone` | `string` | yes | e.g. `eu-central` |
| `locations` | `[]int` | yes | from `GET /locations` |
| `serverType` | `string` | yes | from `GET /types` |
| `minSize` | `int` | yes | |
| `desiredSize` | `int` | yes | initial size |
| `maxSize` | `int` | yes | must be `>= minSize` |
| `monitoringType` | `string` | yes | `cpu` or `memory` |
| `scalingAlgorithm` | `string` | yes | currently `simple` |
| `scaleUpThreshold` | `int` | yes* | required for `simple` |
| `scaleDownThreshold` | `int` | yes* | required for `simple` |

Example:

```json
{
  "templateId": 1,
  "name": "testgroup",
  "zone": "eu-central",
  "locations": [2, 3],
  "serverType": "cx23",
  "minSize": 1,
  "desiredSize": 1,
  "maxSize": 5,
  "monitoringType": "cpu",
  "scalingAlgorithm": "simple",
  "scaleUpThreshold": 70,
  "scaleDownThreshold": 40
}
```

### 4) Management endpoints

- `GET /templates`
- `GET /templates/:id`
- `GET /groups`
- `GET /groups/:id`
- `DELETE /groups/:id`
- `GET /targets` (used by Alloy)

## Node Exporter Requirement

Servers must expose node exporter on port `9100` for scaling signals to work.

You can install it through `cloudConfig`, for example:

```yaml
#cloud-config
package_update: true
package_upgrade: true
packages:
  - prometheus-node-exporter
```

For production, it is better to install Node Exporter once, create a snapshot, and use that snapshot as your image.

## Notes

- Grafana datasource and contact point are auto-created at startup.
- A Grafana alert rule is created when a group is created.
- Change default credentials (`Grafana`, `PostgreSQL`, `pgAdmin`) before production use.
