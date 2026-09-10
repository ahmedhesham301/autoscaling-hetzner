package controller

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ahmedhesham301/autoscaling-hetzner/modules/services"

	"github.com/ahmedhesham301/autoscaling-hetzner/modules/hetzner"

	"github.com/gin-gonic/gin"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func GetAllImages(g *gin.Context) {
	resp, err := hetzner.HClient.Image.All(context.Background())
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	imageSchemas := services.MapStructToSchema(resp, hcloud.SchemaFromImage)
	g.JSON(http.StatusOK, imageSchemas)
}

func DeleteImage(g *gin.Context) {
	id, err := strconv.ParseInt(g.Param("id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	_, err = hetzner.HClient.Image.Delete(context.TODO(), &hcloud.Image{ID: id})
	if err != nil {
		if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
			g.Status(http.StatusNotFound)
			return
		}
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to delete image", "error", err)
		return
	}
	g.Status(http.StatusOK)
}

func GetImageByID(g *gin.Context) {
	id, err := strconv.ParseInt(g.Param("id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	image, _, err := hetzner.HClient.Image.GetByID(context.TODO(), id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("failed to get image", "error", err)
		return
	}
	if image == nil {
		g.Status(http.StatusNotFound)
		return
	}
	g.JSON(http.StatusOK, hcloud.SchemaFromImage(image))
}
