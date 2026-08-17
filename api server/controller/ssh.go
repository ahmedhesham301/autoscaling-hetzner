package controller

import (
	"github.com/ahmedhesham301/autoscaling-hetzner/modules/services"
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"

	"github.com/gin-gonic/gin"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func GetAllSSHKeys(g *gin.Context) {
	resp, err := hetzner.HClient.SSHKey.All(context.TODO())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sshKeySchemas := services.MapStructToSchema(resp, hcloud.SchemaFromSSHKey)
	g.JSON(http.StatusOK, sshKeySchemas)
}

func CreateSSHKey(g *gin.Context) {
	var keySchema schema.SSHKey

	if err := g.ShouldBindJSON(&keySchema); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to bind keySchema json", "error", err)
		return
	}

	key := hcloud.SSHKeyFromSchema(keySchema)

	opts := hcloud.SSHKeyCreateOpts{
		Name:      key.Name,
		PublicKey: key.PublicKey,
		Labels:    key.Labels,
	}
	if err := opts.Validate(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to validate SSHKeyCreateOpts", "error", err)
		return
	}

	key, _, err := hetzner.HClient.SSHKey.Create(context.TODO(), opts)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to create SSH key", "error", err)
		return
	}

	g.JSON(http.StatusCreated, hcloud.SchemaFromSSHKey(key))
}

func DeleteSSHKey(g *gin.Context) {
	id, err := strconv.ParseInt(g.Param("id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	_, err = hetzner.HClient.SSHKey.Delete(context.TODO(), &hcloud.SSHKey{ID: id})
	if err != nil {
		if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
			g.Status(http.StatusNotFound)
			return
		}
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to delete SSH key", "error", err)
		return
	}
	g.Status(http.StatusOK)
}

func GetSSHKeyByID(g *gin.Context) {
	id, err := strconv.ParseInt(g.Param("id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	key, _, err := hetzner.HClient.SSHKey.GetByID(context.TODO(), id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to get SSH key", "error", err)
		return
	}
	if key == nil {
		g.Status(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, hcloud.SchemaFromSSHKey(key))
}
