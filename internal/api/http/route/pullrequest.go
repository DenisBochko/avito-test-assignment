package route

import (
	"avito-test-assignment/internal/api/http/handler"

	"github.com/gin-gonic/gin"
)

func RegisterPRRoutes(g *gin.RouterGroup, h *handler.PullRequestHandler) {
	g.POST("/create", h.Create)
}
