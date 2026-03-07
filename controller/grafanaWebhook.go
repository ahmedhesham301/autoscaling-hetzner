package controller

import (
	"autoscaling-hetzner/services"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ReceiveGrafanaWebhook(g *gin.Context) {
	var alert Alert
	if err := g.ShouldBindJSON(&alert); err != nil {
		g.Status(http.StatusInternalServerError)
		slog.Error("Failed to bind Grafana webhook json")
		return
	}
	g.Status(http.StatusOK)

	for _, alert := range alert.Alerts {
		if alert.Labels["alertname"] != "DatasourceNoData" && alert.Status != "resolved" {
			groupId, err := strconv.Atoi(alert.Labels["groupId"])
			if err != nil {
				slog.Error("Failed to convert groupId label to int", "error", err)
				return
			}
			err = services.ScaleUp(services.ScaleUpOps{GroupId: groupId}, 1, "alert")
			if err != nil {
				slog.Error("Failed to scale up", "groupId", groupId, "error", err)
				return
			}
		}
	}
}
