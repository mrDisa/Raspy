package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrDisa/Raspy/backend/internal/model"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

func (r *GroupRepository) FindByID(ctx context.Context, id int) (*model.Group, error) {
	const query = `
		SELECT id, external_id, name
		FROM groups
		WHERE id = $1
	`

	var group model.Group

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&group.ID,
		&group.ExternalID,
		&group.Name,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find group: %w", err)
	}

	return &group, nil
}

func (r *GroupRepository) Upsert(ctx context.Context, externalID string, name string) (*model.Group, error) {
	const query = `
		INSERT INTO groups (external_id, name)
		VALUES ($1, $2)
		ON CONFLICT (external_id)
		DO UPDATE SET name = EXCLUDED.name
		RETURNING id, external_id, name
	`

	var group model.Group

	err := r.db.QueryRowContext(
		ctx,
		query,
		externalID,
		name,
	).Scan(
		&group.ID,
		&group.ExternalID,
		&group.Name,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to upsert group: %w", err)
	}

	return &group, nil
}