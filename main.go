package main

import (
	"autoscaling-hetzner/controller"
	"autoscaling-hetzner/database"
	"autoscaling-hetzner/hetzner"
	"context"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	database.InitDB()
	hetzner.SetupClient(ctx, os.Getenv("HKEY"))

	server := gin.Default()

	server.GET("/locations", controller.GetLocations)
	server.GET("/images", controller.GetImages)
	server.GET("/types", controller.GetTypes)
	server.PUT("/templates", controller.CreateTemplate)
	server.PUT("/groups", controller.CreateGroup)
	server.GET("/networks", controller.GetNetworks)
	server.Run()
}
