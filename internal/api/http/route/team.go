package route

import (
	"avito-test-assignment/internal/api/http/handler"

	"github.com/gin-gonic/gin"
)

func RegisterTeamRoutes(g *gin.RouterGroup, h *handler.TeamHandler) {
	g.POST("/add", h.AddTeam)
	g.GET("/get", h.GetTeam)
}
