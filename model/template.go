package model

import (
	"autoscaling-hetzner/database"
	"context"
)

type Template struct {
	Id          int     `json:"id"`
	ImageId     int64   `json:"image_id"     binding:"required"`
	Networks    []int64 `json:"networks"`
	SSHKeys     []int64 `json:"SSH_keys"`
	CloudConfig string  `json:"cloudConfig"`
}

func (t *Template) Save() error {
	query := `INSERT INTO templates (image_id, networks, SSH_keys, cloud_config)
	VALUES ($1, $2, $3, $4) RETURNING id;`
	err := database.Pool.QueryRow(
		context.Background(),
		query,
		t.ImageId, t.Networks, t.SSHKeys, t.CloudConfig,
	).Scan(&t.Id)
	return err
}

func (t *Template) GetById(id int) error {
	t.Id = id
	query := `SELECT image_id, networks, SSH_keys, cloud_config
	FROM templates WHERE id = $1`
	err := database.Pool.QueryRow(context.Background(), query, id).
		Scan(&t.ImageId, &t.Networks, &t.SSHKeys, &t.CloudConfig)
	return err
}
