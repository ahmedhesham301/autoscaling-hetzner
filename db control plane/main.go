package main

import (
	"context"
	"os"

	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/controller"
	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/temporal"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/database"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	temporal.SetupClient()
	go temporal.StartWorker()
	database.InitDB()
	hetzner.SetupClient(ctx, os.Getenv("HKEY"))

	server := gin.Default()
	server.POST("/services/:serviceName", controller.CreateDatabase)
	server.GET("/services", controller.ListMangedServices)
	server.GET("/services/:serviceName", controller.GetMangedServiceCreateOps)

	server.GET("/services/os/targets", controller.GetOSTargets)

	// server.GET("/targets", controller.GetTargets)

	server.Run("0.0.0.0:8090")
}
