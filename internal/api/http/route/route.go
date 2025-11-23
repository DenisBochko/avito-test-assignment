package route

import (
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"avito-test-assignment/internal/api/http/handler"
	"avito-test-assignment/internal/api/http/middleware"
	"avito-test-assignment/internal/config"
)

func SetupRouter(
	l *zap.Logger,
	cfg *config.Config,
	teamHdl *handler.TeamHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	router := gin.Default()

	router.Use(middleware.Logger(l))
	router.Use(middleware.RequestTimeout(cfg.Timeout.Request))

	router.HandleMethodNotAllowed = true
	router.NoMethod(handler.NoMethod)
	router.NoRoute(handler.NoRoute)

	basePath := router.Group(cfg.BasePath)

	teamGroup := basePath.Group("/team")
	RegisterTeamRoutes(teamGroup, teamHdl)

	return router
}
