package handler

import (
	"avito-test-assignment/internal/apperrors"
	"avito-test-assignment/internal/model"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TeamService interface {
	AddTeam(ctx context.Context, teamName string, members []model.UserRequest) (err error)
	GetTeam(ctx context.Context, teamName string) (team *model.TeamResponse, err error)
}

type TeamHandler struct {
	l   *zap.Logger
	svc TeamService
}

func NewTeamHandler(l *zap.Logger, svc TeamService) *TeamHandler {
	return &TeamHandler{l, svc}
}

func (h *TeamHandler) AddTeam(c *gin.Context) {
	ctx := c.Request.Context()

	var req model.AddTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ResponseWithError{
			Error: ResponseError{
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			},
		})
	}

	if err := h.svc.AddTeam(ctx, req.TeamName, req.Members); err != nil {
		if errors.Is(err, apperrors.ErrTeamAlreadyExists) {
			c.JSON(http.StatusBadRequest, ResponseWithError{
				Error: ResponseError{
					Code:    "TEAM_EXISTS",
					Message: err.Error(),
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

	c.JSON(http.StatusCreated, req)
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	ctx := c.Request.Context()

	var qp model.TeamNameQueryParam
	if err := c.ShouldBindQuery(&qp); err != nil {
		c.JSON(http.StatusBadRequest, ResponseWithError{
			Error: ResponseError{
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			},
		})
	}

	team, err := h.svc.GetTeam(ctx, qp.TeamName)
	if err != nil {
		if errors.Is(err, apperrors.ErrTeamNotExist) {
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

	c.JSON(http.StatusOK, team)
}
