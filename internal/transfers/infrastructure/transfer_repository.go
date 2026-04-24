package infrastructure

import (
	"context"
	"database/sql"
)

type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) ListOwnedPlayer(
	ctx context.Context,
	userID int64,
	playerID int64,
	askingPrice int64,
) error {
	query := `
		INSERT INTO transfer_listings (player_id, seller_team_id, asking_price)
		SELECT p.id, t.id, $1
		FROM players p
		JOIN teams t ON t.id = p.team_id
		WHERE p.id = $2
		  AND t.user_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, askingPrice, playerID, userID)
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
