package middlewares

import (
	"log/slog"
	"net/http"

	"api-server/data"

	"github.com/gin-gonic/gin"
)

func ValidateParams() gin.HandlerFunc {
	return func(g *gin.Context) {
		var params data.CreateServiceParams
		if err := g.ShouldBindJSON(&params); err != nil {
			slog.Error("Failed to bind body to createServiceParams struct ")
			g.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		err := params.Validate()
		if err != nil {
			slog.Error("error", "validation failed", err.Error())
			g.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		g.Set("params", params)

		// Pre-handler phase
		g.Next()
	}
}
