package controller

import (
	"autoscaling-hetzner/model"
	"autoscaling-hetzner/services"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func CreateGroup(g *gin.Context) {
	var group model.Group
	if err := g.ShouldBindJSON(&group); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to bind group json", "error", err)
		return
	}

	if err := group.Validate(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to validate group", "error", err)
		return
	}

	if err := group.Save(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to save group", "error", err)
		return
	}

	err := services.ScaleUp(services.ScaleOps{Group: &group}, group.DesiredSize, "init")
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("failed to initiate group scale up", "error", err)
		return
	}
	g.Status(http.StatusOK)
	slog.Info("a Group has been created", "groupId", group.Id, "templateId", group.TemplateId)
}

func GetAllGroups(c *gin.Context) {
	groups, err := model.GetAllGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("Failed to get all groups", "error", err)
		return
	}
	if len(groups) == 0 {
		c.JSON(http.StatusOK, []any{})
		return
	}
	c.JSON(http.StatusOK, groups)
}

func GetGroupByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	var group model.Group
	if err := group.GetById(id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("Failed to get group by id", "error", err)
		return
	}

	c.JSON(http.StatusOK, group)
}
