package controller

import (
	"autoscaling-hetzner/services"
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


	for _, alert := range alert.Alerts {
		if alert.Labels["alertname"] != "DatasourceNoData" && alert.Status != "resolved" {
			groupId, err := strconv.Atoi(alert.Labels["groupId"])
			if err != nil {
				log.Printf("%v", err)
				return
			}
			err = services.ScaleUp(services.ScaleUpOps{GroupId: groupId}, 1, "alert")
			if err != nil {
				log.Printf("%v", err)
				return
			}
		}
	}
}
