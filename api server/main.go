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

	server.POST("/webhooks/grafana/alerts", controller.ReceiveGrafanaWebhook)

	server.GET("/targets", controller.GetTargets)

	server.Run()
}
