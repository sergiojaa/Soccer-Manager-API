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

func (r *TeamRepository) FindByUserID(
	ctx context.Context,
	userID int64,
) (*TeamView, error) {
	teamQuery := `
		SELECT
			t.id,
			t.user_id,
			t.name,
			t.country,
			t.budget,
			COALESCE(SUM(p.market_value), 0) AS total_team_value
		FROM teams t
		LEFT JOIN players p ON p.team_id = t.id
		WHERE t.user_id = $1
		GROUP BY t.id
	`

	var team TeamView

	err := r.db.QueryRowContext(ctx, teamQuery, userID).Scan(
		&team.ID,
		&team.UserID,
		&team.Name,
		&team.Country,
		&team.Budget,
		&team.TotalTeamValue,
	)
	if err != nil {
		return nil, err
	}

	playersQuery := `
		SELECT id, first_name, last_name, country, age, position, market_value
		FROM players
		WHERE team_id = $1
		ORDER BY position, id
	`

	rows, err := r.db.QueryContext(ctx, playersQuery, team.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var player PlayerView

		if err := rows.Scan(
			&player.ID,
			&player.FirstName,
			&player.LastName,
			&player.Country,
			&player.Age,
			&player.Position,
			&player.MarketValue,
		); err != nil {
			return nil, err
		}

		team.Players = append(team.Players, player)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &team, nil
}
