package model

import (
	"autoscaling-hetzner/database"
	"context"
)

type Group struct {
	Id          int      `json:"id"`
	TemplateId  int      `json:"templateId"  binding:"required"`
	Zone        string   `json:"zone"        binding:"required"`
	Locations   []string `json:"locations"   binding:"required"`
	ServerTypes []string `json:"serverTypes" binding:"required"`
	MinSize     int      `json:"minSize"     binding:"required"`
	DesiredSize int      `json:"desiredSize" binding:"required"`
	MaxSize     int      `json:"maxSize"     binding:"required"`
	Networks    []string `json:"networks"    binding:"required"`
}

func (g *Group) Save() error {
	query := `INSERT INTO groups (template_id, zone, locations, server_types, min_size, desired_size, max_size, networks)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id;`
	err := database.Pool.QueryRow(
		context.Background(),
		query,
		g.TemplateId, g.Zone, g.Locations, g.ServerTypes, g.MinSize, g.DesiredSize, g.MaxSize, g.Networks,
	).Scan(&g.Id)
	return err
}

func (g *Group) GetById(id int) error {
	g.Id = id
	query := `SELECT template_id, zone, locations, server_types, min_size, desired_size, max_size, networks
	FROM groups WHERE id = $1`
	err := database.Pool.QueryRow(context.Background(), query, id).
		Scan(&g.TemplateId, &g.Zone, &g.Locations, &g.ServerTypes, &g.MinSize, &g.DesiredSize, &g.MaxSize, &g.Networks)
	return err
}
