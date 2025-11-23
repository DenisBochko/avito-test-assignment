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

type PullRequestService interface {
	Create(ctx context.Context, id, name, authorID string) (*model.PullRequestWithAssignedReviewers, error)
}

type PullRequestHandler struct {
	l   *zap.Logger
	svc PullRequestService
}

func NewPullRequestHandler(logger *zap.Logger, svc PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{l: logger, svc: svc}
}

func (s *PullRequestHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var req model.PullRequestCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ResponseWithError{
			Error: ResponseError{
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			},
		})
	}

	pr, err := s.svc.Create(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotExist) || errors.Is(err, apperrors.ErrTeamNotExist) {
			c.JSON(http.StatusNotFound, ResponseWithError{
				Error: ResponseError{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})

			return
		}

		if errors.Is(err, apperrors.ErrPullRequestAlreadyExists) {
			c.JSON(http.StatusConflict, ResponseWithError{
				Error: ResponseError{
					Code:    "PR_EXISTS",
					Message: "PR id already exists",
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

	c.JSON(http.StatusOK, ResponseWithPR{
		PR: pr,
	})
}
