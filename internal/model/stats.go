package model

type ReviewerStats struct {
	ReviewerID    string `db:"reviewer_id"    json:"reviewer_id"`
	AssignedCount int    `db:"assigned_count" json:"assigned_count"`
}

type PRStats struct {
	PullRequestID string `db:"pull_request_id" json:"pull_request_id"`
	ReviewerCount int    `db:"reviewer_count"  json:"reviewer_count"`
}

type StatsResponse struct {
	ReviewerStats []ReviewerStats `json:"reviewer_stats"`
	PRStats       []PRStats       `json:"pr_stats"`
}
