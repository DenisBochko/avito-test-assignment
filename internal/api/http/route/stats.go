package route

import (
	"github.com/gin-gonic/gin"

	"avito-test-assignment/internal/api/http/handler"
)

func RegisterStatsRoutes(g *gin.RouterGroup, h *handler.StatsHandler) {
	g.GET("", h.GetStats)
}
