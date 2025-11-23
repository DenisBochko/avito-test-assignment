package apperrors

import (
	"errors"
)

var (
	ErrTeamNotExist      = errors.New("team does not exist")
	ErrTeamAlreadyExists = errors.New("team already exists")

	ErrUserNotExist = errors.New("user does not exist")
)
