package main

import (
	"autoscaling-hetzner/controller"
	"autoscaling-hetzner/database"
	"autoscaling-hetzner/grafana"
	"autoscaling-hetzner/hetzner"
	"context"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	database.InitDB()
	hetzner.SetupClient(ctx, os.Getenv("HKEY"))
	grafana.InitGrafana()

	server := gin.Default()

	server.POST("/webhooks/grafana/alerts", controller.ReceiveGrafanaWebhook)

	server.GET("/targets", controller.GetTargets)

	server.Run()
}
