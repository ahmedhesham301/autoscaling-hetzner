package controller

import (
	"autoscaling-hetzner/model"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTemplate(g *gin.Context) {
	var template model.Template
	if err := g.ShouldBindJSON(&template); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("Failed to bind template json template", "error", err)
		return
	}
	if err := template.Save(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("Failed to save template", "error", err)
		return
	}
	g.JSON(http.StatusOK, template)
	slog.Info("a template has been created", "templateId", template.Id)

}
