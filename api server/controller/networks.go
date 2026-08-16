package controller

import (
	"autoscaling-hetzner/hetzner"
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func GetAllNetworks(g *gin.Context) {
	resp, err := hetzner.HClient.Network.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, resp)
}

func CreateNetwork(g *gin.Context) {
	var sNetwork schema.Network
	if err := g.ShouldBindJSON(&sNetwork); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to bind sNetwork json", "error", err)
		return
	}
	network := hcloud.NetworkFromSchema(sNetwork)

	opts := hcloud.NetworkCreateOpts{
		Name:                  network.Name,
		IPRange:               network.IPRange,
		Subnets:               network.Subnets,
		Routes:                network.Routes,
		Labels:                network.Labels,
		ExposeRoutesToVSwitch: network.ExposeRoutesToVSwitch,
	}

	if err := opts.Validate(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to validate network", "error", err)
		return
	}

	network, _, err := hetzner.HClient.Network.Create(context.TODO(), opts)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to create network", "error", err)
		return
	}
	g.JSON(http.StatusOK, hcloud.SchemaFromNetwork(network))
}

func DeleteNetwork(g *gin.Context) {
	id, err := strconv.ParseInt(g.Param("id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	_, err = hetzner.HClient.Network.Delete(context.TODO(), &hcloud.Network{ID: id})
	if err != nil {
		if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
			g.Status(http.StatusNotFound)
			return
		}
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to delete network", "error", err)
		return
	}
	g.Status(http.StatusOK)
}

func GetNetworkByID(g *gin.Context) {
	id, err := strconv.ParseInt(g.Param("id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	network, _, err := hetzner.HClient.Network.GetByID(context.TODO(), id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to get network", "error", err)
		return
	}
	if network == nil {
		g.Status(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, hcloud.SchemaFromNetwork(network))
}
