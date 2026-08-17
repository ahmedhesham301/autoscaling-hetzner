package models

import (
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/database"
	"context"
	"errors"
	"log/slog"
)

type Group struct {
	Id                 int     `json:"id"`
	Name               string  `json:"name"                binding:"required"`
	TemplateId         int     `json:"templateId"          binding:"required"`
	Zone               string  `json:"zone"                binding:"required"`
	Locations          []int64 `json:"locations"           binding:"required"`
	ServerType         string  `json:"serverType"          binding:"required"`
	MinSize            int     `json:"minSize"             binding:"required"`
	DesiredSize        int     `json:"desiredSize"         binding:"required"`
	MaxSize            int     `json:"maxSize"             binding:"required"`
	MonitoringType     string  `json:"monitoringType"      binding:"required"`
	ScalingAlgorithm   string  `json:"scalingAlgorithm"    binding:"required"`
	TargetThreshold    *int16  `json:"targetThreshold"`
	ScaleUpThreshold   *int16  `json:"scaleUpThreshold"`
	ScaleDownThreshold *int16  `json:"scaleDownThreshold"`
}

func (g *Group) Save() error {
	query := `INSERT INTO groups (name, template_id, zone, locations, server_type, min_size, desired_size, max_size, monitoring_type, scaling_algorithm, target_threshold, scale_up_threshold,scale_down_threshold)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING id;`
	err := database.Pool.QueryRow(
		context.Background(),
		query,
		g.Name, g.TemplateId, g.Zone, g.Locations, g.ServerType, g.MinSize, g.DesiredSize, g.MaxSize, g.MonitoringType, g.ScalingAlgorithm, g.TargetThreshold, g.ScaleUpThreshold, g.ScaleDownThreshold,
	).Scan(&g.Id)
	return err
}

func (g *Group) GetById(id int) error {
	g.Id = id
	query := `SELECT name, template_id, zone, locations, server_type, min_size, desired_size, max_size, monitoring_type, scaling_algorithm, target_threshold, scale_up_threshold,scale_down_threshold
	FROM groups WHERE id = $1;`
	err := database.Pool.QueryRow(context.Background(), query, id).
		Scan(&g.Name, &g.TemplateId, &g.Zone, &g.Locations, &g.ServerType, &g.MinSize, &g.DesiredSize, &g.MaxSize, &g.MonitoringType, &g.ScalingAlgorithm, &g.TargetThreshold, &g.ScaleUpThreshold, &g.ScaleDownThreshold)
	return err
}

func (g *Group) UpdateDesiredSize(s int) error {
	query := `UPDATE groups SET desired_size=$1 WHERE id=$2;`
	g.DesiredSize += s
	_, err := database.Pool.Exec(context.TODO(), query, g.DesiredSize, g.Id)
	if err != nil {
		return err
	}
	slog.Info("desired state updated", "groupID", g.Id, "before", g.DesiredSize-s, "after", g.DesiredSize)
	return nil
}

func (g *Group) Validate() error {
	if g.MinSize > g.MaxSize {
		return errors.New("max size cant be smaller than min")
	}
	return nil
}

func GetAllGroups() ([]Group, error) {
	query := `SELECT id, name, template_id, zone, locations, server_type, min_size, desired_size, max_size, monitoring_type, 
	scaling_algorithm, target_threshold, scale_up_threshold, scale_down_threshold
	FROM groups;`
	rows, err := database.Pool.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group

	for rows.Next() {
		var g Group
		if err := rows.Scan(
			&g.Id, &g.Name, &g.TemplateId, &g.Zone, &g.Locations, &g.ServerType, &g.MinSize, &g.DesiredSize, &g.MaxSize,
			&g.MonitoringType, &g.ScalingAlgorithm, &g.TargetThreshold, &g.ScaleUpThreshold, &g.ScaleDownThreshold,
		); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func DeleteGroupByID(groupID int) error {
	query := `DELETE FROM public.groups WHERE id=$1;`
	_, err := database.Pool.Exec(context.TODO(), query, groupID)
	return err
}
