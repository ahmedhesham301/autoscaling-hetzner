 plans:
- [ ] implement frontend
- [ ] implement end points to delete modify templates/groups 
- [ ] implement scale out

## Architecture overview
<img width="971" height="681" alt="hetzner auto scaling (1)" src="https://github.com/user-attachments/assets/5a46de66-61dd-46b7-b12b-4ea1943126f4" />

app: the app stores and gets groups/templates/servers in the database, calls the hetzner api, setups alerts in Grafana using grafana api, provides servers ip to grafana alloy.
Grafana: evaluates alert rules and send a webhook to the app.
PostgreSQL: stores the templates,groups,servers data.
Prometheus: stores metrics.
Alloy: scrapes metics annd pushes it to Promethues (why im using it? Because it has http targets discovery)

## What this project is

### How it works



## Installation
> make sure that docker and docker compose are installed
 
 clone this repo
```bash
git clone https://github.com/ahmedhesham301/autoscaling-hetzner.git
```

create a file `.env.compose` inside the repo with the following values
```
HKEY=Your hetzner api key
DATABASE_HOST=db
GRAFANA_HOST=grafana:3000
ENV=prod
```
> if you are running this outside of Hetzner (not recommended) set ENV=dev and enable public ip in your template. so it uses public internet to scrape metrics

start docker compose 
```bash
docker compse up -d --build
```



