package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"avito-test-assignment/internal/model"
	"avito-test-assignment/internal/repository"
)

type TeamRepositoryForUser interface {
	Pool() *pgxpool.Pool

	SelectTeamNameByUserID(ctx context.Context, ext repository.RepoExtension, userID string) (string, error)
}

type UserRepositoryForUser interface {
	UpdateUserActive(ctx context.Context, ext repository.RepoExtension, userID string, isActive bool) error
	SelectUserByID(ctx context.Context, ext repository.RepoExtension, userID string) (*model.User, error)
}

type PullRequestRepositoryForUser interface {
	Pool() *pgxpool.Pool

	SelectPullRequestsByUserID(ctx context.Context, ext repository.RepoExtension, userID string) ([]*model.PullRequest, error)
}

type UserService struct {
	teamRepo        TeamRepositoryForUser
	userRepo        UserRepositoryForUser
	pullRequestRepo PullRequestRepositoryForUser
}

func NewUserService(teamRepo TeamRepositoryForUser, userRepo UserRepositoryForUser, pullRequestRepo PullRequestRepositoryForUser) *UserService {
	return &UserService{
		teamRepo:        teamRepo,
		userRepo:        userRepo,
		pullRequestRepo: pullRequestRepo,
	}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (user *model.UserResponseWithTeamName, err error) {
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

	if err := s.userRepo.UpdateUserActive(ctx, tx, userID, isActive); err != nil {
		return nil, fmt.Errorf("failed to update user active: %w", err)
	}

	teamName, err := s.teamRepo.SelectTeamNameByUserID(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select team name: %w", err)
	}

	userFull, err := s.userRepo.SelectUserByID(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &model.UserResponseWithTeamName{
		TeamName: teamName,
		UserID:   userFull.ID,
		Username: userFull.Username,
		IsActive: userFull.IsActive,
	}, nil
}

func (s *UserService) GetReview(ctx context.Context, userID string) (*model.GetReviewResponse, error) {
	_, err := s.userRepo.SelectUserByID(ctx, nil, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select user: %w", err)
	}

	prs, err := s.pullRequestRepo.SelectPullRequestsByUserID(ctx, nil, userID)
	if err != nil {
		return nil, err
	}

	prsResponse := make([]model.PullRequestResponse, 0, len(prs))

	for _, pr := range prs {
		prsResponse = append(prsResponse, model.PullRequestResponse{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		})
	}

	return &model.GetReviewResponse{
		UserID:       userID,
		PullRequests: prsResponse,
	}, nil
}
