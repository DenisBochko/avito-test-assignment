package service

import (
	"avito-test-assignment/internal/model"
	"avito-test-assignment/internal/repository"
	"context"
	"fmt"
)

type PullRequestRepositoryForStats interface {
	GetReviewerStats(ctx context.Context, ext repository.RepoExtension) ([]model.ReviewerStats, error)
	GetPRStats(ctx context.Context, ext repository.RepoExtension) ([]model.PRStats, error)
}

type StatsService struct {
	pullRequestRepo PullRequestRepositoryForStats
}

func NewStatsService(pullRequestRepo PullRequestRepositoryForStats) *StatsService {
	return &StatsService{
		pullRequestRepo: pullRequestRepo,
	}
}

func (s *StatsService) GetStats(ctx context.Context) (response *model.StatsResponse, err error) {
	reviewer, err := s.pullRequestRepo.GetReviewerStats(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewer stats: %w", err)
	}

	pr, err := s.pullRequestRepo.GetPRStats(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR stats: %w", err)
	}

	return &model.StatsResponse{
		ReviewerStats: reviewer,
		PRStats:       pr,
	}, nil
}
