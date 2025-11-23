package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"avito-test-assignment/internal/model"
	"avito-test-assignment/internal/repository"
)

type TeamRepository interface {
	Pool() *pgxpool.Pool

	InsertTeam(ctx context.Context, ext repository.RepoExtension, teamName string) (int, error)
	SelectTeamIDByName(ctx context.Context, ext repository.RepoExtension, teamName string) (int, error)
}

type UserRepository interface {
	UpsertUser(ctx context.Context, ext repository.RepoExtension, userID string, username string, isActive bool, teamID int) error
	SelectUsersByTeamID(ctx context.Context, ext repository.RepoExtension, teamID int) ([]model.User, error)
}

type TeamService struct {
	teamRepo TeamRepository
	userRepo UserRepository
}

func NewTeamService(teamRepo TeamRepository, userRepo UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s TeamService) AddTeam(ctx context.Context, teamName string, members []model.UserRequest) (err error) {
	tx, err := s.teamRepo.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rErr := tx.Rollback(ctx); rErr != nil {
				err = fmt.Errorf("%w, failed to rollback: %w", err, rErr)
			}
		}
	}()

	teamID, err := s.teamRepo.InsertTeam(ctx, tx, teamName)
	if err != nil {
		return fmt.Errorf("failed to insert team: %w", err)
	}

	for _, user := range members {
		if err = s.userRepo.UpsertUser(ctx, tx, user.UserID, user.Username, user.IsActive, teamID); err != nil {
			return fmt.Errorf("failed to upsert user: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s TeamService) GetTeam(ctx context.Context, teamName string) (team *model.TeamResponse, err error) {
	tx, err := s.teamRepo.Pool().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rErr := tx.Rollback(ctx); rErr != nil {
				err = fmt.Errorf("%w, failed to rollback: %w", err, rErr)
			}
		}
	}()

	teamID, err := s.teamRepo.SelectTeamIDByName(ctx, tx, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to select team ID: %w", err)
	}

	users, err := s.userRepo.SelectUsersByTeamID(ctx, tx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	usersResponse := make([]model.UserResponse, 0, len(users))

	for _, user := range users {
		usersResponse = append(usersResponse, model.UserResponse{
			UserID:   user.ID,
			Username: user.Username,
			IsActive: user.IsActive,
		})
	}

	return &model.TeamResponse{
		TeamName: teamName,
		Members:  usersResponse,
	}, nil
}
