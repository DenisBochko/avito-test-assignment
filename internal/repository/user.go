package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"avito-test-assignment/internal/model"
)

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
	teamID int,
) error {
	if ext == nil {
		ext = r.db
	}

	const query = `
        INSERT INTO users (id, username, is_active, team_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE
        SET username = EXCLUDED.username,
            is_active = EXCLUDED.is_active,
            team_id = EXCLUDED.team_id,
            updated_at = NOW();
    `

	_, err := ext.Exec(ctx, query, userID, username, isActive, teamID)
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
		SELECT id, username, team_id, is_active, created_at, updated_at
		FROM users
		WHERE team_id = $1;
	`

	rows, err := ext.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User

		if err := rows.Scan(&user.ID, &user.Username, &user.TeamID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
