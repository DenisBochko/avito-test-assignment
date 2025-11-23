package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"avito-test-assignment/internal/apperrors"
	"avito-test-assignment/internal/model"
	"avito-test-assignment/internal/repository"
)

const (
	prStatusMerged = "MERGED"
)

type PullRequestRepositoryForPR interface {
	Pool() *pgxpool.Pool

	InsertPullRequest(ctx context.Context, ext repository.RepoExtension, id, name, authorID string) (*model.PullRequest, error)
	SetReviewers(ctx context.Context, ext repository.RepoExtension, authorID, prID string) ([]string, error)
	MergePullRequest(ctx context.Context, ext repository.RepoExtension, prID string) error
	GetAssignedReviewers(ctx context.Context, ext repository.RepoExtension, prID string) ([]string, error)
	SelectPullRequestByID(ctx context.Context, ext repository.RepoExtension, id string) (*model.PullRequest, error)
	GetPRStatus(ctx context.Context, ext repository.RepoExtension, prID string) (string, error)
	GetReviewerCandidates(ctx context.Context, ext repository.RepoExtension, oldReviewerID, prID string) ([]string, error)
	RemoveReviewer(ctx context.Context, ext repository.RepoExtension, prID, reviewerID string) error
	AddReviewer(ctx context.Context, ext repository.RepoExtension, prID, reviewerID string) error
	IsReviewerAssigned(ctx context.Context, ext repository.RepoExtension, prID, reviewerID string) (bool, error)
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

func (s *PullRequestService) Merge(ctx context.Context, pullRequestID string) (*model.MergedResponse, error) {
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

	if err := s.pullRequestRepo.MergePullRequest(ctx, tx, pullRequestID); err != nil {
		return nil, fmt.Errorf("failed to merge pull request: %w", err)
	}

	pr, err := s.pullRequestRepo.SelectPullRequestByID(ctx, tx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to select pull request by ID: %w", err)
	}

	reviewers, err := s.pullRequestRepo.GetAssignedReviewers(ctx, tx, pr.PullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to select assigned reviewers: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &model.MergedResponse{
		PullRequestWithAssignedReviewers: model.PullRequestWithAssignedReviewers{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
			Assigned:        reviewers,
		},
		MergedAt: *pr.MergedAt,
	}, nil
}

func (s *PullRequestService) Reassign(ctx context.Context, pullRequestID, oldReviewerID string) (*model.ReassignResponse, error) {
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

	pr, err := s.pullRequestRepo.SelectPullRequestByID(ctx, tx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to select pull request by ID: %w", err)
	}

	if pr.Status == prStatusMerged {
		return nil, apperrors.ErrPullRequestAlreadyMerged
	}

	isAssigned, err := s.pullRequestRepo.IsReviewerAssigned(ctx, tx, pullRequestID, oldReviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if reviewer is assigned: %w", err)
	}

	if !isAssigned {
		return nil, apperrors.ErrUserIsNotAssignedAsReviewer
	}

	candidates, err := s.pullRequestRepo.GetReviewerCandidates(ctx, tx, oldReviewerID, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers candidates: %w", err)
	}

	if len(candidates) == 0 {
		return nil, apperrors.ErrNoActiveReplacementCandidate
	}

	newReviewer := candidates[0]

	if err := s.pullRequestRepo.RemoveReviewer(ctx, tx, pullRequestID, oldReviewerID); err != nil {
		return nil, fmt.Errorf("failed to remove reviewer: %w", err)
	}

	if err := s.pullRequestRepo.AddReviewer(ctx, tx, pullRequestID, newReviewer); err != nil {
		return nil, fmt.Errorf("failed to add reviewer: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	pr, err = s.pullRequestRepo.SelectPullRequestByID(ctx, nil, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to select pull request by ID: %w", err)
	}

	assigned, err := s.pullRequestRepo.GetAssignedReviewers(ctx, nil, pr.PullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned reviewers: %w", err)
	}

	return &model.ReassignResponse{
		PR: model.PullRequestWithAssignedReviewers{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
			Assigned:        assigned,
		},
		ReplacedBy: newReviewer,
	}, nil
}
