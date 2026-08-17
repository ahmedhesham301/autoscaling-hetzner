package controller

import (
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/services"
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func GetAllFirewalls(g *gin.Context) {
	resp, err := hetzner.HClient.Firewall.All(context.TODO())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	firewallSchemas := services.MapStructToSchema(resp, hcloud.SchemaFromFirewall)
	g.JSON(http.StatusOK, firewallSchemas)
}

func CreateFirewall(g *gin.Context) {
	var sFirewall schema.Firewall

	if err := g.ShouldBindJSON(&sFirewall); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to bind sFirewall json", "error", err)
		return
	}

	firewall := hcloud.FirewallFromSchema(sFirewall)

	opts := hcloud.FirewallCreateOpts{
		Name:    firewall.Name,
		Labels:  firewall.Labels,
		Rules:   firewall.Rules,
		ApplyTo: firewall.AppliedTo,
	}
	if err := opts.Validate(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to validate firewall", "error", err)
		return
	}

	result, _, err := hetzner.HClient.Firewall.Create(context.TODO(), opts)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to create firewall", "error", err)
		return
	}

	g.JSON(http.StatusCreated, hcloud.SchemaFromFirewall(result.Firewall))
}

func GetFirewallByID(g *gin.Context) {
	id, err := strconv.ParseInt(g.Param("id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	firewall, _, err := hetzner.HClient.Firewall.GetByID(context.TODO(), id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to get firewall", "error", err)
		return
	}
	if firewall == nil {
		g.Status(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, hcloud.SchemaFromFirewall(firewall))
}

func DeleteFirewall(g *gin.Context) {
	id, err := strconv.ParseInt(g.Param("id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	_, err = hetzner.HClient.Firewall.Delete(context.TODO(), &hcloud.Firewall{ID: id})
	if err != nil {
		if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
			g.Status(http.StatusNotFound)
			return
		}
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to delete firewall", "error", err)
		return
	}
	g.Status(http.StatusOK)
}
