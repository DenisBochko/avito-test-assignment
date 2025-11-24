package route

import (
	"avito-test-assignment/internal/api/http/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStatsRoutes(g *gin.RouterGroup, h *handler.StatsHandler) {
	g.GET("", h.GetStats)
}
