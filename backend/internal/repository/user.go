package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrDisa/Raspy/backend/internal/model"
)



type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByTelegramID(telegramID int64) (*model.User, error) {
	const query = `
		SELECT id, telegram_id, group_id, subgroup,
		       notifications_enabled, created_at, updated_at
		FROM users
		WHERE telegram_id = $1
	`

	var user model.User

	err := r.db.QueryRowContext(
		context.Background(),
		query,
		telegramID,
	).Scan(
		&user.ID,
		&user.TelegramID,
		&user.GroupID,
		&user.Subgroup,
		&user.NotificationsEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	const query = `
		INSERT INTO users (
			telegram_id,
			group_id,
			subgroup,
			notifications_enabled
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, telegram_id, group_id, subgroup,
		          notifications_enabled, created_at, updated_at
	`

	var createdUser model.User

	err := r.db.QueryRowContext(
		context.Background(),
		query,
		user.TelegramID,
		user.GroupID,
		user.Subgroup,
		user.NotificationsEnabled,
	).Scan(
		&createdUser.ID,
		&createdUser.TelegramID,
		&createdUser.GroupID,
		&createdUser.Subgroup,
		&createdUser.NotificationsEnabled,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &createdUser, nil
}