package repository

import (
	"avito-test-assignment/internal/apperrors"
	"avito-test-assignment/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepository struct {
	db *pgxpool.Pool
}

func NewPullRequestRepository(db *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) Pool() *pgxpool.Pool {
	return r.db
}

func (r *PullRequestRepository) InsertPullRequest(ctx context.Context, ext RepoExtension, id, name, authorID string) (*model.PullRequest, error) {
	if ext == nil {
		ext = r.db
	}

	const query = `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id)
		VALUES ($1, $2, $3)
		RETURNING pull_request_id, pull_request_name, author_id, status, created_at, merged_at;
	`

	var pr model.PullRequest

	err := ext.QueryRow(ctx, query, id, name, authorID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, apperrors.ErrPullRequestAlreadyExists
		}

		return nil, err
	}

	return &pr, nil
}

func (r *PullRequestRepository) SelectPullRequestByID(ctx context.Context, ext RepoExtension, id string) (*model.PullRequest, error) {
	if ext == nil {
		ext = r.db
	}

	const query = `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1;
	`

	var pr model.PullRequest

	err := ext.QueryRow(ctx, query, id).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrPullRequestNotExist
		}

		return nil, err
	}

	return &pr, nil
}

func (r *PullRequestRepository) SetReviewers(ctx context.Context, ext RepoExtension, authorID, prID string) ([]string, error) {
	if ext == nil {
		ext = r.db
	}

	const query = `
		WITH team_members AS (
		    SELECT u.id AS user_id,
		           COUNT(pr.pull_request_id) AS open_prs
		    FROM team_lnk tl
		    JOIN users u ON u.id = tl.user_id
		    LEFT JOIN pr_reviewers prr ON prr.reviewer_id = u.id
		    LEFT JOIN pull_requests pr ON pr.pull_request_id = prr.pull_request_id AND pr.status = 'OPEN'
		    WHERE tl.team_id IN (
		        SELECT team_id
		        FROM team_lnk
		        WHERE user_id = $1
		    )
		    AND u.id <> $1
			AND u.is_active = true
		    GROUP BY u.id
		    ORDER BY open_prs ASC
		    LIMIT 2
		)
		INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
		SELECT $2, user_id
		FROM team_members
		RETURNING reviewer_id;
	`

	rows, err := ext.Query(ctx, query, authorID, prID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reviewers []string

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		reviewers = append(reviewers, id)
	}

	return reviewers, nil
}

func (r *PullRequestRepository) SelectPullRequestsByUserID(ctx context.Context, ext RepoExtension, userID string) ([]*model.PullRequest, error) {
	if ext == nil {
		ext = r.db
	}

	const query = `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at
		FROM pull_requests pr
		JOIN pr_reviewers r ON r.pull_request_id = pr.pull_request_id
		WHERE r.reviewer_id = $1;
	`

	rows, err := ext.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var prs []*model.PullRequest

	for rows.Next() {
		var pr model.PullRequest

		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
			&pr.CreatedAt,
			&pr.MergedAt,
		); err != nil {
			return nil, err
		}

		prs = append(prs, &pr)
	}

	return prs, nil
}
