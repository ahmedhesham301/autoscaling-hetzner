package controller

import (
	"autoscaling-hezner/model"
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)
func CreateTemplate(g *gin.Context) {
	var template model.Template
	if err := g.ShouldBindJSON(&template); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	template.Id = rand.Int64()
	if err := template.Save(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g.Status(http.StatusOK)
}

