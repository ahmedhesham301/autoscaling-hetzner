package controller

import (
	"autoscaling-hetzner/hetzner"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLocations(g *gin.Context) {
	resp, err := hetzner.HClient.Location.All(context.TODO())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, resp)
}

func GetImages(g *gin.Context) {
	resp, err := hetzner.HClient.Image.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, resp)
}

func GetTypes(g *gin.Context) {
	resp, err := hetzner.HClient.ServerType.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, resp)
}

func GetNetworks(g *gin.Context) {
	resp, err := hetzner.HClient.Network.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, resp)
}

func GetFirewalls(g *gin.Context) {
	resp, err := hetzner.HClient.Firewall.All(context.TODO())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, resp)
}

func GetSSHKeys(g *gin.Context) {
	resp, err := hetzner.HClient.SSHKey.All(context.TODO())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, resp)
}
