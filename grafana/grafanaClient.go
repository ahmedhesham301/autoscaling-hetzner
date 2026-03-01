package grafana

import (
	"log"
	"net/url"
	"os"

	"github.com/go-openapi/strfmt"
	gapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
)

var GClient *gapi.GrafanaHTTPAPI
var PrometheusUid string
var FolderUid string

func InitGrafana() {
	grafanaHost, exists := os.LookupEnv("GRAFANA_HOST")
	if !exists {
		log.Panic("GRAFANA_HOST is not set; expected the grafana hostname")
	}

	GClient = gapi.NewHTTPClientWithConfig(strfmt.Default, &gapi.TransportConfig{
		Host:      grafanaHost,
		BasePath:  "/api",
		BasicAuth: url.UserPassword("admin", "admin"),
		Schemes:   []string{"http"},
	})
	
	resp, err := GClient.Datasources.GetDataSources()
	if err != nil {
		log.Panicf("%v", err)
	}
	PrometheusUid = resp.GetPayload()[0].UID

	re, err := GClient.Folders.GetFolderByID(1)
	if err != nil {
		resp, err := GClient.Folders.CreateFolder(&models.CreateFolderCommand{Title: "alerts"})
		if err != nil {
			log.Panicf("could not create an alert folder in grafana: %v", err)
		}
		FolderUid = resp.Payload.UID
		return
	}

	FolderUid = re.Payload.UID
}
