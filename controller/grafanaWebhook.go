package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReceiveGrafanaWebhook(g *gin.Context) {
	var alert Alert
	if err := g.ShouldBindJSON(&alert); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g.Status(http.StatusOK)
	jsonData, err := json.Marshal(alert)
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Println(string(jsonData))
}
