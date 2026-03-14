package controller

import (
	"autoscaling-hetzner/model"
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

			// jsonData, err := json.Marshal(alert)
			// if err != nil {
			// 	log.Fatal("Error:", err)
			// }
			// fmt.Println(string(jsonData))

			groupId, err := strconv.Atoi(alert.Labels["groupId"])
			if err != nil {
				slog.Error("Failed to convert groupId label to int", "error", err)
				return
			}

			var group model.Group
			err = group.GetById(groupId)
			if err != nil {
				slog.Error("Failed to get group by id", "error", err)
				return
			}

			if group.ScalingAlgorithm == "simple" {
				if alert.Values["B0"] <= float64(*group.ScaleDownThreshold) {
					err = services.ScaleOut(services.ScaleOps{Group: &group})
					if err != nil {
						slog.Error("Failed to scale out", "groupID", group.Id, "error", err)
					}
				} else {
					err = services.ScaleUp(services.ScaleOps{Group: &group}, 1, "alert")
					if err != nil {
						slog.Error("Failed to scale up", "groupId", groupId, "error", err)
					} else {
						slog.Info("group scaled up", "groupID", group.Id)
					}
				}
			}
		}
	}
}
