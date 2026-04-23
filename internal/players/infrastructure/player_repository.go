package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	playersApp "github.com/sergiojaa/soccer-manager-api/internal/players/application"
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
	players []playersApp.PlayerInput,
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
