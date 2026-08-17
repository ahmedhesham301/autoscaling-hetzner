package main

import (
	"context"
	"os"
	"github.com/ahmedhesham301/autoscaling-hetzner/control-plane/controller"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/database"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	database.InitDB()
	hetzner.SetupClient(ctx, os.Getenv("HKEY"))

	server := gin.Default()

	server.POST("/webhooks/grafana/alerts", controller.ReceiveGrafanaWebhook)

	server.GET("/targets", controller.GetTargets)
	
	server.Run("0.0.0.0:8085")
}
