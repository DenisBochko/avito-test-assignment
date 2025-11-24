package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"avito-test-assignment/internal/model"
)

const (
	StatusNotAvailable = "NOT_AVAILABLE"
)

type ResponseWithError struct {
	Error ResponseError `json:"error"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ResponseWithMessage struct {
	Status  string `json:"status"`
	Message string `son:"message"`
}

type ResponseWithUser struct {
	User *model.UserResponseWithTeamName `json:"user"`
}

type ResponseWithPR struct {
	PR any `json:"pr"`
}

func NoMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, ResponseWithError{
		Error: ResponseError{
			Code:    StatusNotAvailable,
			Message: "method not allowed on this endpoint",
		},
	})
}

func NoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, ResponseWithError{
		Error: ResponseError{
			Code:    StatusNotAvailable,
			Message: "page not found",
		},
	})
}
