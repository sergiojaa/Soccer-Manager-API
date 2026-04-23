package infrastructure

import (
	"context"
	"database/sql"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	name string,
	country string,
) (int64, error) {
	query := `
		INSERT INTO teams (user_id, name, country)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(ctx, query, userID, name, country).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
