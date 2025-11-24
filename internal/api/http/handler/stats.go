package handler

import (
	"avito-test-assignment/internal/model"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StatsService interface {
	GetStats(ctx context.Context) (response *model.StatsResponse, err error)
}

type StatsHandler struct {
	l   *zap.Logger
	svc StatsService
}

func NewStatsHandler(logger *zap.Logger, svc StatsService) *StatsHandler {
	return &StatsHandler{
		l:   logger,
		svc: svc,
	}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.svc.GetStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseWithError{
			Error: ResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
	}

	c.JSON(http.StatusOK, stats)
}
