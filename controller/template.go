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

func GetAllTemplates(g *gin.Context) {
	templates, err := model.GetAllTemplates()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("Failed to get all templates", "error", err)
		return
	}
	if len(templates) == 0 {
		g.JSON(http.StatusOK, []any{})
		return
	}
	g.JSON(http.StatusOK, templates)
}
