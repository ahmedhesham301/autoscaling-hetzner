package controller

import (
	"autoscaling-hetzner/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGroup(g *gin.Context) {
	var group model.Group
	if err := g.ShouldBindJSON(&group); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := group.Save(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g.Status(http.StatusOK)
}
