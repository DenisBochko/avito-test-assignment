package route

import (
	"github.com/gin-gonic/gin"

	"avito-test-assignment/internal/api/http/handler"
)

func RegisterTeamRoutes(g *gin.RouterGroup, h *handler.TeamHandler) {
	g.POST("/add", h.AddTeam)
	g.GET("/get", h.GetTeam)
}
