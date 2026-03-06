package controller

import (
	"autoscaling-hetzner/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	for _, alert := range alert.Alerts {
		if alert.Labels["alertname"] != "DatasourceNoData" && alert.Status != "resolved"{
			groupId, err := strconv.Atoi(alert.Labels["groupId"])
			if err != nil {
				log.Fatal("Error:", err)
			}	
			services.ScaleUp(services.ScaleUpOps{GroupId: groupId}, 1)
		}
	}
}
