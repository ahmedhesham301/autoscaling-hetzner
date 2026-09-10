package model

import (
	"context"
	"net"
	"os"
	"strconv"

	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/data"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/database"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var metricsPorts = map[string]int{
	"postgresql": 9187,
}
var healthPaths = map[string]string{
	"postgresql": ":8008/patroni",
}

func SaveDB(ctx context.Context, server hcloud.Server, params data.CreateDBParams, DB_ID int) error {
	env := os.Getenv("ENV")
	var ip string
	if env == "dev" {
		ip = server.PublicNet.IPv4.IP.String()
	} else {
		ip = server.PrivateNet[0].IP.String()
	}
	query := `UPDATE databases SET server_id=$1, server_name=$2, location=$3, private_ip=$4, node_exporter=$5, service_exporter=$6
	WHERE id = $7;`
	_, err := database.Pool.Exec(
		ctx, query,
		server.ID, server.Name, server.Location.Name, ip, params.NodeExporter, params.ServiceExporter, DB_ID,
	)
	return err
}

func CreateDBRecord(ctx context.Context, params data.CreateDBParams) (int, error) {
	var id int
	query := `INSERT INTO databases (kind)
	VALUES ($1)
	RETURNING id;`
	row := database.Pool.QueryRow(
		ctx, query,
		params.AppName,
	)
	err := row.Scan(&id)
	return id, err
}

func GetOSTargets(ctx context.Context) (*[]data.Target, error) {
	targets := []data.Target{}

	query := "SELECT private_ip, server_name, server_id, kind FROM databases WHERE node_exporter=true;"
	rows, err := database.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ip net.IP
		var serverName, kind string
		var serverID int64
		if err := rows.Scan(&ip, &serverName, &serverID, &kind); err != nil {
			return nil, err
		}
		targets = append(targets, data.Target{
			Targets: []string{ip.String() + ":9100"},
			Labels: map[string]string{
				"server_name": serverName,
				"server_id":   strconv.FormatInt(serverID, 10),
				"kind":        kind,
			},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &targets, nil
}
