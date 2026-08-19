package controller

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ahmedhesham301/autoscaling-hetzner/modules/helper"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/services"

	"github.com/gin-gonic/gin"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func GetAllServers(g *gin.Context) {
	resp, err := hetzner.HClient.Server.All(context.TODO())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	serverSchemas := services.MapStructToSchema(resp, hcloud.SchemaFromServer)
	g.JSON(http.StatusOK, serverSchemas)
}

func CreateServer(g *gin.Context) {
	var serverSchema schema.ServerCreateRequest

	if err := g.ShouldBindJSON(&serverSchema); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to bind sServer json", "error", err)
		return
	}

	err, opts := helper.ToServerCreateOpts(serverSchema)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to create validate server create ops", "error", err)
		return
	}

	result, _, err := hetzner.HClient.Server.Create(context.TODO(), *opts)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to create server", "error", err)
		return
	}

	g.JSON(http.StatusCreated, hcloud.SchemaFromServer(result.Server))
}

func GetServerByID(g *gin.Context) {
	id, err := strconv.ParseInt(g.Param("id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	server, _, err := hetzner.HClient.Server.GetByID(context.TODO(), id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to get firewall", "error", err)
		return
	}
	if server == nil {
		g.Status(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, hcloud.SchemaFromServer(server))
}

func DeleteServer(g *gin.Context) {
	id, err := strconv.ParseInt(g.Param("id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	_, _, err = hetzner.HClient.Server.DeleteWithResult(context.TODO(), &hcloud.Server{ID: id})
	if err != nil {
		if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
			g.Status(http.StatusNotFound)
			return
		}
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to delete server", "error", err)
		return
	}
	g.Status(http.StatusOK)
}
