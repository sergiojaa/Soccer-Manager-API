package infrastructure

import (
	"context"
	"database/sql"
)

type TransferRepository struct {
	db *sql.DB
}
type ActiveListingForUpdate struct {
	ID           int64
	PlayerID     int64
	SellerTeamID int64
	AskingPrice  int64
	MarketValue  int64
}

type TeamForUpdate struct {
	ID     int64
	Budget int64
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

func (r *TransferRepository) FindActiveListingForUpdate(
	ctx context.Context,
	tx *sql.Tx,
	listingID int64,
) (*ActiveListingForUpdate, error) {
	query := `
		SELECT
			tl.id,
			tl.player_id,
			tl.seller_team_id,
			tl.asking_price,
			p.market_value
		FROM transfer_listings tl
		JOIN players p ON p.id = tl.player_id
		WHERE tl.id = $1
		  AND tl.status = 'ACTIVE'
		FOR UPDATE
	`

	var listing ActiveListingForUpdate

	err := tx.QueryRowContext(ctx, query, listingID).Scan(
		&listing.ID,
		&listing.PlayerID,
		&listing.SellerTeamID,
		&listing.AskingPrice,
		&listing.MarketValue,
	)
	if err != nil {
		return nil, err
	}

	return &listing, nil
}

func (r *TransferRepository) FindTeamByUserIDForUpdate(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
) (*TeamForUpdate, error) {
	query := `
		SELECT id, budget
		FROM teams
		WHERE user_id = $1
		FOR UPDATE
	`

	var team TeamForUpdate

	err := tx.QueryRowContext(ctx, query, userID).Scan(
		&team.ID,
		&team.Budget,
	)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (r *TransferRepository) DecreaseTeamBudget(
	ctx context.Context,
	tx *sql.Tx,
	teamID int64,
	amount int64,
) error {
	query := `
		UPDATE teams
		SET budget = budget - $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, amount, teamID)
	return err
}

func (r *TransferRepository) IncreaseTeamBudget(
	ctx context.Context,
	tx *sql.Tx,
	teamID int64,
	amount int64,
) error {
	query := `
		UPDATE teams
		SET budget = budget + $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, amount, teamID)
	return err
}

func (r *TransferRepository) MovePlayerToTeam(
	ctx context.Context,
	tx *sql.Tx,
	playerID int64,
	newTeamID int64,
	newMarketValue int64,
) error {
	query := `
		UPDATE players
		SET team_id = $1,
		    market_value = $2,
		    updated_at = NOW()
		WHERE id = $3
	`

	_, err := tx.ExecContext(ctx, query, newTeamID, newMarketValue, playerID)
	return err
}

func (r *TransferRepository) MarkListingSold(
	ctx context.Context,
	tx *sql.Tx,
	listingID int64,
) error {
	query := `
		UPDATE transfer_listings
		SET status = 'SOLD',
		    closed_at = NOW()
		WHERE id = $1
	`

	_, err := tx.ExecContext(ctx, query, listingID)
	return err
}

func (r *TransferRepository) CreateTransferHistory(
	ctx context.Context,
	tx *sql.Tx,
	listingID int64,
	playerID int64,
	sellerTeamID int64,
	buyerTeamID int64,
	salePrice int64,
	oldMarketValue int64,
	newMarketValue int64,
) error {
	query := `
		INSERT INTO transfers (
			listing_id,
			player_id,
			seller_team_id,
			buyer_team_id,
			sale_price,
			market_value_before,
			market_value_after
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		listingID,
		playerID,
		sellerTeamID,
		buyerTeamID,
		salePrice,
		oldMarketValue,
		newMarketValue,
	)

	return err
}
