package controller

import (
	"autoscaling-hezner/model"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"

	"github.com/aws/smithy-go"
	"github.com/gin-gonic/gin"
)

func CreateGroup(g *gin.Context) {
	var group model.Group
	if err := g.ShouldBindJSON(&group); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	group.Id = rand.Int64()
	err := model.GetTemplate(group.TemplateId)
	if err != nil {
		var opErr smithy.APIError
		if errors.As(err, &opErr) {
			if opErr.ErrorCode() == "NoSuchKey" {
				g.JSON(http.StatusBadRequest, gin.H{"message": "template does not exist"})
				return
			} else {
				g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

		} else {
			g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	}

	if err := group.Save(); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("%+v", group)
	g.Status(http.StatusOK)
}
