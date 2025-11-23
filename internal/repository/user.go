package repository

import (
	"avito-test-assignment/internal/apperrors"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"avito-test-assignment/internal/model"
)

const userListDefaultCap = 10

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UpsertUser(
	ctx context.Context,
	ext RepoExtension,
	userID string,
	username string,
	isActive bool,
) error {
	if ext == nil {
		ext = r.db
	}

	const query = `
        INSERT INTO users (id, username, is_active)
        VALUES ($1, $2, $3)
        ON CONFLICT (id) DO UPDATE
        SET username = EXCLUDED.username,
            is_active = EXCLUDED.is_active,
            updated_at = NOW();
    `

	_, err := ext.Exec(ctx, query, userID, username, isActive)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) SelectUsersByTeamID(ctx context.Context, ext RepoExtension, teamID int) ([]model.User, error) {
	if ext == nil {
		ext = r.db
	}

	const query = `
		SELECT u.id, u.username, u.is_active, u.created_at, u.updated_at
		FROM users u 
		JOIN team_lnk l ON u.id = l.user_id
		WHERE l.team_id = $1;
	`

	rows, err := ext.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]model.User, 0, userListDefaultCap)

	for rows.Next() {
		var user model.User

		if err := rows.Scan(&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) SelectUserByID(ctx context.Context, ext RepoExtension, userID string) (*model.User, error) {
	if ext == nil {
		ext = r.db
	}

	const query = `
		SELECT u.id, u.username, u.is_active, u.created_at, u.updated_at 
		FROM users u
		WHERE u.id = $1;
	`

	var user model.User

	if err := ext.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Username, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotExist
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserActive(ctx context.Context, ext RepoExtension, userID string, isActive bool) error {
	if ext == nil {
		ext = r.db
	}

	const query = `
		Update users 
		SET is_active = $1
		WHERE id = $2;
	`

	cmd, err := ext.Exec(ctx, query, isActive, userID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return apperrors.ErrUserNotExist
	}

	return nil
}
