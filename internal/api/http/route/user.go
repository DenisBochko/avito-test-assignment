package route

import (
	"github.com/gin-gonic/gin"

	"avito-test-assignment/internal/api/http/handler"
)

func RegisterUsersRoutes(g *gin.RouterGroup, h *handler.UserHandler) {
	g.POST("/setIsActive", h.SetIsActive)
	g.GET("/getReview", h.GetReview)
}
