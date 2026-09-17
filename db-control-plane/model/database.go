package model

import (
	"context"
	"net"
	"strconv"

	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/data"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/database"
)

var metricsPorts = map[string]int{
	"postgresql": 9187,
}
var healthPaths = map[string]string{
	"postgresql": ":8008/patroni",
}

func GetOSTargets(ctx context.Context) (*[]data.Target, error) {
	targets := []data.Target{}

	query := "SELECT private_ip, server_name, server_id, type, engine FROM services WHERE node_exporter=true;"
	rows, err := database.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ip net.IP
		var serverName, serviceType, serviceEngine string
		var serverID int64
		if err := rows.Scan(&ip, &serverName, &serverID, &serviceType, &serviceEngine); err != nil {
			return nil, err
		}
		targets = append(targets, data.Target{
			Targets: []string{ip.String() + ":9100"},
			Labels: map[string]string{
				"server_name": serverName,
				"server_id":   strconv.FormatInt(serverID, 10),
				"type":        serviceType,
				"engine":      serviceEngine,
			},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &targets, nil
}
