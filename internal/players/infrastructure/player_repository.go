package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sergiojaa/soccer-manager-api/internal/players/domain"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) CreateBatch(
	ctx context.Context,
	tx *sql.Tx,
	players []domain.PlayerInput,
) error {

	if len(players) == 0 {
		return nil
	}

	var (
		queryBuilder strings.Builder
		args         []interface{}
	)

	queryBuilder.WriteString(`
		INSERT INTO players (team_id, first_name, last_name, country, age, position)
		VALUES
	`)

	for i, p := range players {
		start := i * 6

		queryBuilder.WriteString(fmt.Sprintf(
			"($%d,$%d,$%d,$%d,$%d,$%d)",
			start+1,
			start+2,
			start+3,
			start+4,
			start+5,
			start+6,
		))

		if i < len(players)-1 {
			queryBuilder.WriteString(",")
		}

		args = append(args,
			p.TeamID,
			p.FirstName,
			p.LastName,
			p.Country,
			p.Age,
			p.Position,
		)
	}

	_, err := tx.ExecContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *PlayerRepository) UpdateOwnedPlayer(
	ctx context.Context,
	userID int64,
	playerID int64,
	firstName string,
	lastName string,
	country string,
) error {
	query := `
		UPDATE players AS p
		SET first_name = $1,
		    last_name = $2,
		    country = $3,
		    updated_at = NOW()
		FROM teams AS t
		WHERE p.id = $4
		  AND p.team_id = t.id
		  AND t.user_id = $5
	`

	result, err := r.db.ExecContext(ctx, query, firstName, lastName, country, playerID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
