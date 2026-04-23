package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrDuplicateEmail = errors.New("duplicate email")
var ErrUserNotFound = errors.New("user not found")

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

type UserAuthRecord struct {
	ID           int64
	Email        string
	PasswordHash string
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*UserAuthRecord, error) {
	query := `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`

	var user UserAuthRecord

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
