package model

import (
	"autoscaling-hetzner/database"
	"context"
)

type Template struct {
	Id          int     `json:"id"`
	ImageId     int64   `json:"image_id"     binding:"required"`
	Networks    []int64 `json:"networks"     binding:"required"`
	SSHKeys     []int64 `json:"SSH_keys"`
	PublicIPv4  *bool   `json:"publicIPv4"   binding:"required"`
	PublicIPv6  *bool   `json:"publicIPv6"   binding:"required"`
	CloudConfig string  `json:"cloudConfig"`
}

func (t *Template) Save() error {
	query := `INSERT INTO templates (image_id, networks, SSH_keys, public_ipv4, public_ipv6,cloud_config)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
	err := database.Pool.QueryRow(
		context.Background(),
		query,
		t.ImageId, t.Networks, t.SSHKeys, t.PublicIPv4, t.PublicIPv6, t.CloudConfig,
	).Scan(&t.Id)
	return err
}

func (t *Template) GetById(id int) error {
	t.Id = id
	query := `SELECT image_id, networks, SSH_keys, public_ipv4, public_ipv6, cloud_config
	FROM templates WHERE id = $1`
	err := database.Pool.QueryRow(context.Background(), query, id).
		Scan(&t.ImageId, &t.Networks, &t.SSHKeys, &t.PublicIPv4, &t.PublicIPv6, &t.CloudConfig)
	return err
}
