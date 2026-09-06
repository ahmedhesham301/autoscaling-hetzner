package temporal

import (
	"log"
	"log/slog"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
)

var TemporalClient client.Client

func SetupClient() {
	var err error
	// Loads the "default" profile from the standard location and environment variables.
	_, exists := os.LookupEnv("TEMPORAL_ADDRESS")
	if !exists {
		slog.Error("env var is not set", "var", "TEMPORAL_ADDRESS")
		os.Exit(1)
	}
	TemporalClient, err = client.Dial(envconfig.MustLoadDefaultClientOptions())
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
}
