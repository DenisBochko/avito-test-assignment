package model

type TeamResponse struct {
	TeamName string         `json:"team_name"`
	Members  []UserResponse `json:"members"`
}

type AddTeamRequest struct {
	TeamName string        `binding:"required" json:"team_name"`
	Members  []UserRequest `binding:"required" json:"members"`
}

type TeamNameQueryParam struct {
	TeamName string `binding:"required" form:"team_name"`
}
