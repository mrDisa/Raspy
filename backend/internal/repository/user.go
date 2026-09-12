package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrDisa/Raspy/backend/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

type UserRepository interface {
    FindByTelegramID(telegramID int64) (model.User, error)
	Create(telegramID int64, groupID *int, subgroup model.Subgroup) (model.User, error)
}

func (r *PostgresUserRepository) FindByTelegramID(telegramID int64) (model.User, error) {
    var user model.User

    query := `SELECT id, telegram_id, group_id, subgroup, notifications_enabled, created_at, updated_at
              FROM users WHERE telegram_id = $1`

    err := r.db.QueryRow(query, telegramID).Scan(
        &user.ID,
        &user.TelegramID,
		&user.GroupID,
		&user.Subgroup,
		&user.NotificationsEnabled,
		&user.CreatedAt,
		&user.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("telegram_id: %d %w",telegramID, ErrUserNotFound)
		}
        return model.User{}, err
    }

    return user, nil
}

func (r *PostgresUserRepository) Create(telegramID int64, groupID *int, subgroup model.Subgroup) (model.User, error) {
    query := `INSERT INTO users (telegram_id, group_id, subgroup)
              VALUES ($1, $2, $3)
              RETURNING id, notifications_enabled, created_at, updated_at`

    user := model.User{
        TelegramID: telegramID,
        GroupID:    groupID,
        Subgroup:   subgroup,
    }

    err := r.db.QueryRow(query, telegramID, groupID, subgroup).Scan(
        &user.ID,
        &user.NotificationsEnabled,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return model.User{}, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}