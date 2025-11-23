package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"avito-test-assignment/internal/apperrors"
	"avito-test-assignment/internal/model"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*model.UserResponseWithTeamName, error)
}

type UserHandler struct {
	l   *zap.Logger
	svc UserService
}

func NewUserHandler(l *zap.Logger, svc UserService) *UserHandler {
	return &UserHandler{
		l:   l,
		svc: svc,
	}
}

func (h *UserHandler) SetIsActive(c *gin.Context) {
	ctx := c.Request.Context()

	var req model.UserIsActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ResponseWithError{
			Error: ResponseError{
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			},
		})
	}

	user, err := h.svc.SetIsActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotExist) {
			c.JSON(http.StatusNotFound, ResponseWithError{
				Error: ResponseError{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})

			return
		}

		c.JSON(http.StatusInternalServerError, ResponseWithError{
			Error: ResponseError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})

		return
	}

	c.JSON(http.StatusOK, ResponseWithUser{
		User: user,
	})
}
