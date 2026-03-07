package model

import (
	"autoscaling-hetzner/database"
	"context"
)

type Group struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"           binding:"required"`
	TemplateId     int     `json:"templateId"     binding:"required"`
	Zone           string  `json:"zone"           binding:"required"`
	Locations      []int64 `json:"locations"      binding:"required"`
	ServerType     string  `json:"serverType"     binding:"required"`
	MinSize        int     `json:"minSize"        binding:"required"`
	DesiredSize    int     `json:"desiredSize"    binding:"required"`
	MaxSize        int     `json:"maxSize"        binding:"required"`
	MonitoringType string  `json:"monitoringType" binding:"required"`
	Target         int16   `json:"target"         binding:"required"`
}

func (g *Group) Save() error {
	query := `INSERT INTO groups (name, template_id, zone, locations, server_type, min_size, desired_size, max_size, monitoring_type, target)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id;`
	err := database.Pool.QueryRow(
		context.Background(),
		query,
		g.Name, g.TemplateId, g.Zone, g.Locations, g.ServerType, g.MinSize, g.DesiredSize, g.MaxSize, g.MonitoringType, g.Target,
	).Scan(&g.Id)
	return err
}

func (g *Group) GetById(id int) error {
	g.Id = id
	query := `SELECT name, template_id, zone, locations, server_type, min_size, desired_size, max_size, monitoring_type, target
	FROM groups WHERE id = $1;`
	err := database.Pool.QueryRow(context.Background(), query, id).
		Scan(&g.Name, &g.TemplateId, &g.Zone, &g.Locations, &g.ServerType, &g.MinSize, &g.DesiredSize, &g.MaxSize, &g.MonitoringType, &g.Target)
	return err
}

func (g *Group) UpdateDesiredSize(s int) error {
	query := `UPDATE groups SET desired_size=$1 WHERE id=$2;`
	g.DesiredSize += s
	_, err := database.Pool.Exec(context.TODO(), query, g.DesiredSize, g.Id)
	return err
}
