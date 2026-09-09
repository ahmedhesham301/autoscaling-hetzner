package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/data"
	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/model"
	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/temporal"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

func CreateService(g *gin.Context) {
	var params data.CreateDBParams
	if err := g.ShouldBindJSON(&params); err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
		slog.Error("Failed to bind body to createServiceParams struct ")
		return
	}

	DB_ID, err := model.CreateDBRecord(g, params)
	if err != nil {
		g.Status(http.StatusInternalServerError)
		slog.Error("Failed to create database record", "error", err)
		return
	}

	options := client.StartWorkflowOptions{
		ID:        "create-database-workflow" + strconv.Itoa(DB_ID),
		TaskQueue: "task-queue",
	}

	env := os.Getenv("ENV")
	templatesPath := os.Getenv("PACKER_TEMPLATES_PATH")

	var networkID *int64
	networkvar, exists := os.LookupEnv("networkID")
	if exists {
		id, err := strconv.ParseInt(networkvar, 10, 64)
		if err != nil {
			g.Status(http.StatusInternalServerError)
			slog.Error("Failed to Parse network id", "error", err)
			return
		}
		networkID = &id
	}
	_, err = temporal.TemporalClient.ExecuteWorkflow(context.TODO(), options, temporal.CreateServiceWorkflow, params, DB_ID, env, templatesPath, networkID)
	if err != nil {
		g.Status(http.StatusInternalServerError)
		slog.Error("Failed to Execute workflow", "error", err)
		return
	}
	g.JSON(http.StatusAccepted, gin.H{"id": DB_ID})
}

func ListMangedServices(g *gin.Context) {
	templatesPath := os.Getenv("PACKER_TEMPLATES_PATH")

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
	serviceName := g.Param("kind")

	templatesPath := os.Getenv("PACKER_TEMPLATES_PATH")

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

func GetOSTargets(g *gin.Context) {

	targets, err := model.GetOSTargets(g)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("Failed to expose targets", "error", err)
		return
	}

	g.JSON(http.StatusOK, targets)
}
