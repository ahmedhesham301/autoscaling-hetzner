package model

import (
	"autoscaling-hetzner/database"
	"context"
	"net"
	"time"

	"github.com/jackc/pgx/v5"
)

type Server struct {
	Id        int
	Name      string
	GroupId   int
	Type      string
	Location  string
	PrivateIp net.IP
	Create_at time.Time
}

func (s *Server) Save() error {
	query := `INSERT INTO servers (name, group_id, type, location, private_ip)
	VALUES ($1 ,$2, $3, $4, $5);`
	_, err := database.Pool.Exec(
		context.TODO(), query,
		s.Name, s.GroupId, s.Type, s.Location, s.PrivateIp.String(),
	)
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
