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

	server.GET("/locations", controller.GetLocations)
	server.GET("/images", controller.GetImages)
	server.GET("/types", controller.GetTypes)
	server.GET("/networks", controller.GetNetworks)
	server.GET("/firewalls", controller.GetFirewalls)
	server.GET("/keys", controller.GetSSHKeys)

	server.PUT("/templates", controller.CreateTemplate)
	server.PUT("/groups", controller.CreateGroup)

	server.GET("/targets", controller.GetTargets)

	server.Run()
}
