package grafana

import (
	"log"
	"net/url"
	"os"

	"github.com/go-openapi/strfmt"
	gapi "github.com/grafana/grafana-openapi-client-go/client"
)

var GClient *gapi.GrafanaHTTPAPI
var PrometheusUid string
var FolderUid string

func InitGrafana() {
	GClient = gapi.NewHTTPClientWithConfig(strfmt.Default, &gapi.TransportConfig{
		Host:      os.Getenv("GRAFANA_HOST"),
		BasePath:  "/api",
		BasicAuth: url.UserPassword("admin", "admin"),
		Schemes:   []string{"http"},
	})
	resp, err := GClient.Datasources.GetDataSources()
	if err != nil {
		log.Panicf("%v", err)
	}
	PrometheusUid = resp.GetPayload()[0].UID
	// resp, err := GClient.Folders.GetFolders(folders.NewGetFoldersParams())
	// if err != nil {
	// 	log.Panic(err)
	// }
	// res, err := GClient.Folders.CreateFolder(&models.CreateFolderCommand{Title: "test"})
	// if err != nil {
	// 	log.Panicf("%v", err)
	// }
	re, err := GClient.Folders.GetFolderByID(1)
	if err != nil {
		log.Panicf("%v", err)
	}

	FolderUid = re.Payload.UID
	println(re.Payload.UID)

}
