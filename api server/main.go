package main

import (
	"autoscaling-hetzner/controller"
	"context"
	"os"

	"github.com/ahmedhesham301/autoscaling-hetzner/modules/grafana"

	"github.com/ahmedhesham301/autoscaling-hetzner/modules/database"

	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	database.InitDB()
	hetzner.SetupClient(ctx, os.Getenv("HKEY"))
	grafana.InitGrafana()

	server := gin.Default()

	server.GET("/locations", controller.GetLocations)
	server.GET("/server_types", controller.GetServerTypes)

	server.GET("/networks", controller.GetAllNetworks)
	server.GET("/networks/:id", controller.GetNetworkByID)
	server.POST("/networks", controller.CreateNetwork)
	server.DELETE("/networks/:id", controller.DeleteNetwork)

	server.GET("/firewalls", controller.GetAllFirewalls)
	server.GET("/firewalls/:id", controller.GetFirewallByID)
	server.POST("/firewalls", controller.CreateFirewall)
	server.DELETE("/firewalls/:id", controller.DeleteFirewall)

	server.GET("/ssh_keys", controller.GetAllSSHKeys)
	server.GET("/ssh_keys/:id", controller.GetSSHKeyByID)
	server.POST("/ssh_keys", controller.CreateSSHKey)
	server.DELETE("/ssh_keys/:id", controller.DeleteSSHKey)

	server.POST("/templates", controller.CreateTemplate)
	server.GET("/templates", controller.GetAllTemplates)
	server.GET("/templates/:id", controller.GetTemplateByID)

	server.POST("/groups", controller.CreateGroup)
	server.GET("/groups", controller.GetAllGroups)
	server.GET("/groups/:id", controller.GetGroupByID)
	server.DELETE("/groups/:id", controller.DeleteGroupByID)

	server.GET("/servers", controller.GetAllServers)
	server.GET("/servers/:id", controller.GetServerByID)
	server.POST("/servers", controller.CreateServer)
	server.DELETE("/servers/:id", controller.DeleteServer)

	server.GET("/images", controller.GetAllImages)
	server.GET("/images/:id", controller.GetImageByID)
	server.DELETE("/images/:id", controller.DeleteImage)
	server.Run()
}
