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

func (r *TransferRepository) FindActiveListings(
	ctx context.Context,
) ([]MarketListingView, error) {
	query := `
		SELECT
			tl.id,
			p.id,
			p.first_name,
			p.last_name,
			p.country,
			p.age,
			p.position,
			p.market_value,
			tl.asking_price,
			t.id,
			t.name
		FROM transfer_listings tl
		JOIN players p ON p.id = tl.player_id
		JOIN teams t ON t.id = tl.seller_team_id
		WHERE tl.status = 'ACTIVE'
		ORDER BY tl.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	listings := make([]MarketListingView, 0)

	for rows.Next() {
		var listing MarketListingView

		if err := rows.Scan(
			&listing.ID,
			&listing.PlayerID,
			&listing.FirstName,
			&listing.LastName,
			&listing.Country,
			&listing.Age,
			&listing.Position,
			&listing.MarketValue,
			&listing.AskingPrice,
			&listing.SellerTeamID,
			&listing.SellerTeam,
		); err != nil {
			return nil, err
		}

		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listings, nil
}
