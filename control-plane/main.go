package main

import (
	"github.com/ahmedhesham301/hetzner-control-plane/control-plane/controller"
	"github.com/ahmedhesham301/hetzner-control-plane/modules/database"
	"github.com/ahmedhesham301/hetzner-control-plane/modules/hetzner"

	"github.com/gin-gonic/gin"
)

func main() {

	database.InitDB()
	hetzner.SetupClient()

	server := gin.Default()

	server.POST("/webhooks/grafana/alerts", controller.ReceiveGrafanaWebhook)

	server.GET("/targets", controller.GetTargets)

	server.Run("0.0.0.0:8085")
}
