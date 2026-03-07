# autoscaling-hetzner

Autoscaling control plane for Hetzner Cloud instances.

## What This Project Does

- Stores templates, groups, and provisioned servers in PostgreSQL.
- Creates Hetzner servers from saved templates.
- Exposes dynamic scrape targets for Grafana Alloy at `GET /targets`.
- Creates Grafana alert rules per group.
- Receives Grafana webhook alerts and performs scale-up.

## Architecture

<img width="971" height="681" alt="hetzner auto scaling (1)" src="https://github.com/user-attachments/assets/5a46de66-61dd-46b7-b12b-4ea1943126f4" />

Components:

- `app` (Go): API + orchestration logic.
- `PostgreSQL`: storage for templates, groups, and servers.
- `Prometheus`: metric storage.
- `Alloy`: target discovery + scraping + remote_write to Prometheus.
- `Grafana`: alert evaluation + webhook notifications.

## Current Status

- Implemented: create template, create group, scale-up, webhook handling, dynamic targets.
- Pending: scale-down/scale-out, update/delete endpoints for templates/groups, frontend.

## Prerequisites

- Docker + Docker Compose
- Hetzner Cloud API token

Important before running:

- `configs/alloy/main.alloy` currently points to `http://192.168.1.40:8080/targets`.
- `grafana/grafanaClient.go` currently sets webhook URL to `http://192.168.1.40:8080/webhooks/grafana/alerts`.

Update both to your app's reachable address.

## Installation

1. Clone the repo:

```bash
git clone https://github.com/ahmedhesham301/autoscaling-hetzner.git
cd autoscaling-hetzner
```

2. create a `.env.compose` file with the following values in it

```bash
HKEY=<your_hetzner_api_token>
DATABASE_HOST=db
GRAFANA_HOST=grafana:3000
ENV=prod
```

3. Start the the containers:

```bash
docker compose up -d --build
```

By default the API listens on `:8080`.

## Default Local Ports

- App: `8080`
- Grafana: `3000` (default credentials in this setup: `admin` / `admin`)
- Prometheus: `9090`
- PostgreSQL: `5432`
- Alloy HTTP: `12345`

## API Endpoints

Discovery and Hetzner metadata:

- `GET /locations`
- `GET /images`
- `GET /types`
- `GET /networks`
- `GET /firewalls`
- `GET /keys`

Core operations:

- `POST /templates`
- `POST /groups`
- `GET /targets`
- `POST /webhooks/grafana/alerts`

## Notes and Limitations

- `POST /groups` creates servers immediately via Hetzner API.
- Alert webhook currently triggers scale-up by `1` when firing and not resolved.
- No authentication layer is implemented for API endpoints yet.
- Several config values are currently hardcoded and should be parameterized for production.

