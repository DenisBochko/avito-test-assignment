package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"avito-test-assignment/internal/model"
	"avito-test-assignment/internal/repository"
)

type PullRequestRepositoryForPR interface {
	Pool() *pgxpool.Pool

	InsertPullRequest(ctx context.Context, ext repository.RepoExtension, id, name, authorID string) (*model.PullRequest, error)
	SetReviewers(ctx context.Context, ext repository.RepoExtension, authorID, prID string) ([]string, error)
}

type UserRepositoryForPR interface {
	SelectUserByID(ctx context.Context, ext repository.RepoExtension, userID string) (*model.User, error)
}

type TeamRepositoryForPR interface {
	SelectTeamIDByUserID(ctx context.Context, ext repository.RepoExtension, userID string) (int, error)
}

type PullRequestService struct {
	pullRequestRepo PullRequestRepositoryForPR
	userRepo        UserRepositoryForPR
	teamRepo        TeamRepositoryForPR
}

func NewPullRequestService(
	pullRequestRepo PullRequestRepositoryForPR,
	userRepo UserRepositoryForPR,
	teamRepo TeamRepositoryForPR,
) *PullRequestService {
	return &PullRequestService{
		pullRequestRepo: pullRequestRepo,
		userRepo:        userRepo,
		teamRepo:        teamRepo,
	}
}

func (s *PullRequestService) Create(ctx context.Context, id, name, authorID string) (*model.PullRequestWithAssignedReviewers, error) {
	tx, err := s.pullRequestRepo.Pool().Begin(ctx)
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

	_, err = s.userRepo.SelectUserByID(ctx, tx, authorID)
	if err != nil {
		return nil, err
	}

	_, err = s.teamRepo.SelectTeamIDByUserID(ctx, tx, authorID)
	if err != nil {
		return nil, err
	}

	pr, err := s.pullRequestRepo.InsertPullRequest(ctx, tx, id, name, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert pull request: %w", err)
	}

	rIDs, err := s.pullRequestRepo.SetReviewers(ctx, tx, authorID, pr.PullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to set reviewers: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &model.PullRequestWithAssignedReviewers{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          pr.Status,
		Assigned:        rIDs,
	}, nil
}
