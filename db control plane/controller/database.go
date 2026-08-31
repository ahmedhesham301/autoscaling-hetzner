package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/data"
	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/temporal"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

func CreateDatabase(g *gin.Context) {
	var params data.CreateServiceParams
	params.AppName = g.Param("serviceName")
	if err := g.ShouldBindJSON(&params); err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
		slog.Error("Failed to bind body to createServiceParams struct ")
		return
	}

	options := client.StartWorkflowOptions{
		ID:        "create-service-workflow",
		TaskQueue: "task-queue",
	}

	_, err := temporal.TemporalClient.ExecuteWorkflow(context.TODO(), options, temporal.CreateServiceWorkflow, params)
	if err != nil {
		g.Status(http.StatusInternalServerError)
		slog.Error("Failed to bind body to go struct", "error", err)
		return
	}
}

func ListMangedServices(g *gin.Context) {
	templatesPath, exists := os.LookupEnv("PACKER_TEMPLATES_PATH")
	if !exists {
		slog.Error("env var PACKER_TEMPLATES_PATH is not set")
		g.Status(http.StatusInternalServerError)
		return
	}

	files, err := os.ReadDir(templatesPath)
	if err != nil {
		slog.Error("error", "failed to read so dir templates path", err)
		g.Status(http.StatusInternalServerError)
		return
	}
	services := []string{}
	for _, file := range files {
		if file.Name()[0] != '.' {
			services = append(services, file.Name())
		}
	}
	g.JSON(http.StatusOK, services)
}

func GetMangedServiceCreateOps(g *gin.Context) {
	serviceName := g.Param("serviceName")

	templatesPath, exists := os.LookupEnv("PACKER_TEMPLATES_PATH")
	if !exists {
		slog.Error("env var PACKER_TEMPLATES_PATH is not set")
		os.Exit(1)
	}

	content, err := os.ReadFile(templatesPath + "/" + serviceName + "/info.json")
	if err != nil {
		slog.Error("Error when opening file", "err", err)
		g.JSON(http.StatusBadRequest, gin.H{"error": "no template called " + serviceName})
		return
	}

	var serviceInfo data.ServiceInfo
	err = json.Unmarshal(content, &serviceInfo)
	if err != nil {
		slog.Error("Error during Unmarshal", "err", err)
		return
	}
	g.JSON(http.StatusOK, serviceInfo)

}
