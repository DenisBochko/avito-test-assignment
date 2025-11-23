package route

import (
	"avito-test-assignment/internal/api/http/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUsersRoutes(g *gin.RouterGroup, h *handler.UserHandler) {
	g.POST("/setIsActive", h.SetIsActive)
}
