package apperrors

import (
	"errors"
)

var (
	ErrTeamNotExist      = errors.New("team does not exist")
	ErrTeamAlreadyExists = errors.New("team already exists")

	ErrUserNotExist = errors.New("user does not exist")

	ErrPullRequestAlreadyExists     = errors.New("pull request already exists")
	ErrPullRequestNotExist          = errors.New("pull request does not exist")
	ErrPullRequestAlreadyMerged     = errors.New("pull request already merged")
	ErrNoActiveReplacementCandidate = errors.New("no active replacement candidate in team")
	ErrUserIsNotAssignedAsReviewer  = errors.New("user is not assigned as reviewer on pr")
)
