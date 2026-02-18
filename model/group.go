package model

import (
	"autoscaling-hetzner/database"
	"context"
)

type Group struct {
	Id          int64    `json:"id"`
	TemplateId  int64    `json:"templateId"  binding:"required"`
	Zone        string   `json:"zone"        binding:"required"`
	Locations   []string `json:"locations"   binding:"required"`
	ServerTypes []string `json:"serverTypes" binding:"required"`
	MinSize     int      `json:"minSize"     binding:"required"`
	DesiredSize int      `json:"desiredSize" binding:"required"`
	MaxSize     int      `json:"maxSize"     binding:"required"`
}

func (g *Group) Save() error {
	query := `INSERT INTO groups (template_id, zone, locations, server_types, min_size, desired_size, max_size)
	VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err := database.Pool.Exec(
		context.Background(),
		query,
		g.TemplateId, g.Zone, g.Locations, g.ServerTypes, g.MinSize, g.DesiredSize, g.MaxSize,
	)
	return err
}
