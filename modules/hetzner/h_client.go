package hetzner

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var HClient *hcloud.Client

func SetupClient(ctx context.Context, token string) {
	HClient = hcloud.NewClient(
		hcloud.WithToken(token),
	)
 }
