package controller

import (
	"autoscaling-hetzner/hetzner"
	"autoscaling-hetzner/vars"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLocations(g *gin.Context) {
	g.JSON(http.StatusOK, vars.Zones)
}

func GetImages(g *gin.Context) {
	resp, err := hetzner.HClient.Image.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	images := make(map[string][]string)
	for _, image := range resp {
		images[image.OSFlavor] = append(images[image.OSFlavor], image.OSVersion)
	}
	g.JSON(http.StatusOK, images)
}

func GetTypes(g *gin.Context) {
	resp, err := hetzner.HClient.ServerType.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// types := make()
	for _, t := range resp {
		fmt.Println(t.Category)
		fmt.Println(t.CPUType)
		fmt.Println(t.Architecture)
		fmt.Println(t.Cores)
		fmt.Println(t.Memory)
		fmt.Println(t.Disk)
		fmt.Printf("%+v", t.Pricings)
	}
}

func GetNetworks(g *gin.Context) {
	resp, err := hetzner.HClient.Network.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	networks := make(map[string]int64)
	for _, network := range resp {
		networks[network.Name] = network.ID
	}
	g.JSON(http.StatusOK, networks)
}
