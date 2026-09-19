package controller

import (
	"log/slog"
	"net/http"

	"github.com/ahmedhesham301/autoscaling-hetzner/services-control-plane/model"
	"github.com/gin-gonic/gin"
)

func GetOSTargets(g *gin.Context) {

	targets, err := model.GetOSTargets(g)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("Failed to expose targets", "error", err)
		return
	}

	g.JSON(http.StatusOK, targets)
}
