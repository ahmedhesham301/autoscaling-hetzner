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



## Prerequisites
- Docker + Docker Compose 
- Hetzner Cloud API token

## Installation
### Locally
(not recommended only for testing)
> must enable public ip in the template 

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
ENV=dev
```

3. Start the the containers:

```bash
docker compose up -d --build
```

By default the API listens on `:8080`.
### on Hetzner (recommended)
> must allow outgoing traffic to public internet (to be able to access Hetzner api )
>  must be in the network where the server are going to be
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

## Default Ports

- App: `8080`
- Grafana: `3000` (default credentials in this setup: `admin` / `admin`)
- Prometheus: `9090`
- PostgreSQL: `5432` (default credentials: postgres / 1234)
- Alloy: `12345`

## usage

currently the project does not have a frontend so u will have to use the api
### create a templates
make a post request to the endpoint with the following body

| param       | type         | required | where to get   |
| ----------- | ------------ | -------- | -------------- |
| image_id    | int          | yes      | GET /images    |
| SSH_Keys    | list of ints | no       | GET /keys      |
| firewalls   | list of ints | no       | GET /firewalls |
| networks    | list of ints | no       | GET /networks  |
| publicIPv4  | bool         | yes      |                |
| publicIPv6  | bool         | yes      |                |
| cloudConfig | string       | no       |                |
> prometheus node exporter must be installed on the server
> it is recommended to make your own snapshot with it installed
> but it can also be installed at startup using the cloud config below(tested on ubuntu and debain) 

example 
```json
{
	"image_id": 310554929,
	"Networks": [11952339],
	"SSH_Keys": [107916411],
	"publicIPv4": true,
	"publicIPv6": true,
	"cloudConfig":"#cloud-config\npackage_update: true\npackage_upgrade: true\npackages:\n - prometheus-node-exporter\n - stress"
}
```


