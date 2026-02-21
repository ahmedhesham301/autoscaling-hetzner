package model

import (
	"autoscaling-hetzner/database"
	"context"
	"net"
	"time"
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
		context.Background(), query,
		s.Name, s.GroupId, s.Type, s.Location, s.PrivateIp.To4(),
	)
	return err
}
