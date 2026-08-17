package models

import (
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/database"
	"context"
	"net"
	"time"

	"github.com/jackc/pgx/v5"
)

type Server struct {
	ID        int64
	Name      string
	GroupId   int
	Type      string
	Location  int64
	PrivateIp net.IP
	CreatedAt time.Time
}

func (s *Server) Save() error {
	query := `INSERT INTO servers (id, name, group_id, type, location, private_ip)
	VALUES ($1 ,$2, $3, $4, $5, $6);`
	_, err := database.Pool.Exec(
		context.TODO(), query,
		s.ID, s.Name, s.GroupId, s.Type, s.Location, s.PrivateIp.String(),
	)
	return err
}

func (s *Server) DeleteServer() error {
	query := "DELETE FROM servers WHERE name=$1;"
	_, err := database.Pool.Exec(context.TODO(), query, s.Name)
	return err
}

func GetAllServers() (pgx.Rows, error) {
	query := "SELECT name, group_id, private_ip FROM servers;"
	resp, err := database.Pool.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAllServersInGroup(groupID int) (map[int64][]Server, error) {
	query := "SELECT id, name, location, created_at FROM servers WHERE group_id=$1;"
	resp, err := database.Pool.Query(context.TODO(), query, groupID)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	servers := make(map[int64][]Server)

	for resp.Next() {
		var server Server
		if err := resp.Scan(&server.ID, &server.Name, &server.Location, &server.CreatedAt); err != nil {
			return nil, err
		}

		servers[server.Location] = append(servers[server.Location], server)
	}
	return servers, nil
}
