package model

import (
	"time"
)

type UserRequest struct {
	UserID   string `binding:"required" json:"user_id"`
	Username string `binding:"required" json:"username"`
	IsActive bool   `binding:"required" json:"is_active"`
}

type User struct {
	ID        string
	Username  string
	TeamID    int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
