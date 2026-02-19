package model

import (
	"autoscaling-hetzner/database"
	"context"
	"net"
)

type Server struct {
	Id        int
	GroupId   int
	Type      string
	Location  string
	PrivateIp net.IP
}

func (s *Server) Save() error {
	query := `INSERT INTO servers (group_id, type, location, private_ip)
	VALUES ($1 ,$2, $3, $4);`
	_, err := database.Pool.Exec(
		context.Background(), query,
		s.GroupId, s.Type, s.Location, s.PrivateIp.To4(),
	)
	return err
}
