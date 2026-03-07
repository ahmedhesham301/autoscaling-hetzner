package controller

import (
	"autoscaling-hetzner/model"
	"autoscaling-hetzner/services"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGroup(g *gin.Context) {
	var group model.Group
	if err := g.ShouldBindJSON(&group); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to bind group json", "error", err)
		return
	}
	if err := group.Save(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to save group", "error", err)
		return
	}
	err := services.ScaleUp(services.ScaleUpOps{Group: &group}, group.DesiredSize, "init")
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to initiate group scale up", "error", err)
		return
	}
	g.Status(http.StatusOK)
	slog.Info("a Group has been created", "groupId", group.Id, "templateId", group.TemplateId)
}
