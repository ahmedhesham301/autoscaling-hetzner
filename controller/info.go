package controller

import (
	"autoscaling-hezner/hetzner"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLocations(g *gin.Context) {
	resp, err := hetzner.HClient.Location.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
	}

	locations := make(map[string][]string)
	for _, location := range resp {
		locations[string(location.NetworkZone)] = append(locations[string(location.NetworkZone)], location.City)
	}
	g.JSON(http.StatusOK, locations)
}

func GetImages(g *gin.Context) {
	resp, err := hetzner.HClient.Image.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
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


