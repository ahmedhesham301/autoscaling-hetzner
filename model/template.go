package model

import (
	"autoscaling-hetzner/database"
	"context"
)

type Template struct {
	Id          int64  `json:"id"`
	Image       string `json:"image"      binding:"required"`
	CloudConfig string `json:"cloudConfig"`
}

func (t *Template) Save() error {
	query := `INSERT INTO templates (image, cloud_config)
	VALUES ($1, $2);`
	_, err := database.Pool.Exec(
		context.Background(),
		query,
		t.Image, t.CloudConfig,
	)
	return err
}

// func GetTemplate(id int64) error {

// 	return err
// }
