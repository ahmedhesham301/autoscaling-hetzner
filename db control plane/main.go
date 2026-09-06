package main

import (
	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/controller"
	"github.com/ahmedhesham301/autoscaling-hetzner/db-control-plane/temporal"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/config"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/database"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ValidateEnvVars([]string{"TARGET_ENV", "ENV", "PACKER_TEMPLATES_PATH"})
	temporal.SetupClient()
	go temporal.StartWorker()
	database.InitDB()
	hetzner.SetupClient()

	server := gin.Default()
	server.POST("/services/:kind", controller.CreateService)
	server.GET("/services", controller.ListMangedServices)
	server.GET("/services/:kind", controller.GetMangedServiceCreateOps)

	server.GET("/services/monitoring/os/targets", controller.GetOSTargets)
	// server.GET("/services/monitoring/:kind/targets", controller.GetDatabaseTargets)

	// server.GET("/targets", controller.GetTargets)

	server.Run("0.0.0.0:8090")
}
