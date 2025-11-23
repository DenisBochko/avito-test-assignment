package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"avito-test-assignment/internal/apperrors"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Pool() *pgxpool.Pool {
	return r.db
}

func (r *TeamRepository) InsertTeam(ctx context.Context, ext RepoExtension, teamName string) (int, error) {
	if ext == nil {
		ext = r.db
	}

	const query = `
		INSERT INTO teams (team_name)
		VALUES ($1)
		RETURNING id;
	`

	var id int

	if err := ext.QueryRow(ctx, query, teamName).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, apperrors.ErrTeamAlreadyExists
		}

		return 0, err
	}

	return id, nil
}

func (r *TeamRepository) SelectTeamIDByName(ctx context.Context, ext RepoExtension, teamName string) (int, error) {
	if ext == nil {
		ext = r.db
	}

	const query = `
		SELECT id FROM teams 
		WHERE team_name = $1;
	`

	var id int

	if err := ext.QueryRow(ctx, query, teamName).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, apperrors.ErrTeamNotExist
		}

		return 0, err
	}

	return id, nil
}

func (r *TeamRepository) InsertTeamLinkWithUser(ctx context.Context, ext RepoExtension, teamID int, userID string) error {
	if ext == nil {
		ext = r.db
	}

	const query = `
		INSERT INTO team_lnk (user_id, team_id)
		VALUES ($1, $2)
	`

	_, err := ext.Exec(ctx, query, userID, teamID)
	if err != nil {
		return err
	}

	return nil
}
