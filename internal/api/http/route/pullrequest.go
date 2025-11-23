package route

import (
	"github.com/gin-gonic/gin"

	"avito-test-assignment/internal/api/http/handler"
)

func RegisterPRRoutes(g *gin.RouterGroup, h *handler.PullRequestHandler) {
	g.POST("/create", h.Create)
}
