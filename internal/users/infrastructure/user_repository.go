package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	email string,
	passwordHash string,
) (int64, error) {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(ctx, query, email, passwordHash).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, ErrDuplicateEmail
		}
		return 0, err
	}

	return id, nil
}
