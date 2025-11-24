package model

import (
	"time"
)

type PullRequest struct {
	PullRequestID   string     `json:"pull_request_id"`
	PullRequestName string     `json:"pull_request_name"`
	AuthorID        string     `json:"author_id"`
	Status          string     `json:"status"`
	Assigned        []string   `json:"assigned_reviewers"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	MergedAt        *time.Time `json:"mergedAt,omitempty"`
}

type PullRequestWithAssignedReviewers struct {
	PullRequestID   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorID        string   `json:"author_id"`
	Status          string   `json:"status"`
	Assigned        []string `json:"assigned_reviewers"`
}

type PullRequestCreateRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type PullRequestResponse struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type GetReviewResponse struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestResponse `json:"pull_requests"`
}

type GetReviewRequestUserIDParam struct {
	UserID string `binding:"required" form:"user_id"`
}

type MergedResponse struct {
	PullRequestWithAssignedReviewers

	MergedAt time.Time `json:"mergedAt"`
}

type MergedRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type ReassignResponse struct {
	PR         PullRequestWithAssignedReviewers `json:"pr"`
	ReplacedBy string                           `json:"replaced_by"`
}
