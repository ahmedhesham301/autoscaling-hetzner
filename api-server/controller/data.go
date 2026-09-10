package controller

import (
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/services"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func GetLocations(g *gin.Context) {
	resp, err := hetzner.HClient.Location.All(context.TODO())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	locationSchemas := services.MapStructToSchema(resp, hcloud.SchemaFromLocation)
	g.JSON(http.StatusOK, locationSchemas)
}

func GetServerTypes(g *gin.Context) {
	resp, err := hetzner.HClient.ServerType.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	typeSchemas := services.MapStructToSchema(resp, hcloud.SchemaFromServerType)

	g.JSON(http.StatusOK, typeSchemas)
}
