package main

import (
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/database"
	"github.com/ahmedhesham301/autoscaling-hetzner/services-control-plane/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	server := gin.Default()

	server.GET("/services/monitoring/os/targets", controller.GetOSTargets)


	server.Run("0.0.0.0:8090")
}
