package temporal

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/envconfig"
)

var TemporalClient client.Client

func SetupClient() {
	var err error
	// Loads the "default" profile from the standard location and environment variables.
	TemporalClient, err = client.Dial(envconfig.MustLoadDefaultClientOptions())
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
}
