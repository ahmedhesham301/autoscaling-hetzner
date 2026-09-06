package hetzner

import (
	"log/slog"
	"os"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var HClient *hcloud.Client

func SetupClient() {
	heyznerKey, exists := os.LookupEnv("HKEY")
	if !exists {
		slog.Error("env var is not set", "var", "HKEY")
		os.Exit(1)
	}
	HClient = hcloud.NewClient(
		hcloud.WithToken(heyznerKey),
	)
}
