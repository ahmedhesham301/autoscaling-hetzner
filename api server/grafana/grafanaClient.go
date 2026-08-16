package grafana

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag/conv"
	gapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/provisioning"
	"github.com/grafana/grafana-openapi-client-go/models"
)

var GClient *gapi.GrafanaHTTPAPI
var PrometheusUid string
var FolderUid string
var ContactPointUid string

func InitGrafana() {
	grafanaHost, exists := os.LookupEnv("GRAFANA_HOST")
	if !exists {
		slog.Error("env var GRAFANA_HOST is not set")
		os.Exit(1)
	}

	controllerHost, exists := os.LookupEnv("CONTROLLER_HOST")
	if !exists {
		slog.Error("env var CONTROLLER_HOST is not set")
		os.Exit(1)
	}

	GClient = gapi.NewHTTPClientWithConfig(strfmt.Default, &gapi.TransportConfig{
		Host:      grafanaHost,
		BasePath:  "/api",
		BasicAuth: url.UserPassword("admin", "admin"),
		Schemes:   []string{"http"},
	})

	// get data source uid
	resp, err := GClient.Datasources.GetDataSources()
	if err != nil {
		slog.Error("Failed to get datasource uid", "error", err)
		os.Exit(1)
	}
	PrometheusUid = resp.GetPayload()[0].UID

	// Ensure the alerts folder exists in Grafana.
	// If it already exists, use its UID.
	// Otherwise, create it and store the newly created folder UID.
	re, err := GClient.Folders.GetFolderByID(1)
	if err != nil {
		resp, err := GClient.Folders.CreateFolder(&models.CreateFolderCommand{Title: "alerts"})
		if err != nil {
			slog.Error("could not create an alert folder in grafana", "error", err)
			os.Exit(1)
		}
		FolderUid = resp.Payload.UID

	} else {
		FolderUid = re.Payload.UID
	}

	// Ensure a "server" webhook contact point exists in Grafana.
	// Reuse its UID if found, otherwise create it and store the new UID.
	resp2, err := GClient.Provisioning.GetContactpoints(&provisioning.GetContactpointsParams{Name: conv.Pointer("server")})
	if err != nil {
		slog.Error("Failed to get Grafana contact points", "error", err)
		os.Exit(1)
	}
	if len(resp2.Payload) == 0 {
		body := models.EmbeddedContactPoint{
			Name: "server",
			Type: conv.Pointer("webhook"),
			Settings: map[string]any{
				"url": fmt.Sprintf("http://%s:8080/webhooks/grafana/alerts", controllerHost),
			},
		}
		_, err := GClient.Provisioning.PostContactpoints(provisioning.NewPostContactpointsParams().WithBody(&body))
		if err != nil {
			slog.Error("Failed to create Grafana contact points", "error", err)
			os.Exit(1)
		}
	}
}
