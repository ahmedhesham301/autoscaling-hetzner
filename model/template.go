package model

import (
	"autoscaling-hetzner/database"
	"context"
)

type Template struct {
	Id          int    `json:"id"`
	OSFlavor    string `json:"OSFlavor"     binding:"required"`
	OSVersion   string `json:"OSVersion"    binding:"required"`
	CloudConfig string `json:"cloudConfig"`
}

func (t *Template) Save() error {
	query := `INSERT INTO templates (os_flavor, os_version, cloud_config)
	VALUES ($1, $2, $3);`
	_, err := database.Pool.Exec(
		context.Background(),
		query,
		t.OSFlavor, t.OSVersion, t.CloudConfig,
	)
	return err
}

func (t *Template) GetById(id int) error {
	t.Id = id
	query := `SELECT os_flavor, os_version, cloud_config
	FROM templates WHERE id = $1`
	err := database.Pool.QueryRow(context.Background(), query, id).
		Scan(&t.OSFlavor, &t.OSVersion, &t.CloudConfig)
	return err
}
