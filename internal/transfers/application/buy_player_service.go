package application

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/sergiojaa/soccer-manager-api/internal/transfers/infrastructure"
)

var (
	ErrListingNotFound    = errors.New("transfer listing not found")
	ErrCannotBuyOwnPlayer = errors.New("cannot buy your own player")
	ErrInsufficientBudget = errors.New("insufficient budget")
)

type BuyPlayerService struct {
	db           *sql.DB
	transferRepo TransferBuyerRepository
}

type TransferBuyerRepository interface {
	FindActiveListingForUpdate(
		ctx context.Context,
		tx *sql.Tx,
		listingID int64,
	) (*infrastructure.ActiveListingForUpdate, error)
	FindTeamByUserIDForUpdate(
		ctx context.Context,
		tx *sql.Tx,
		userID int64,
	) (*infrastructure.TeamForUpdate, error)
	DecreaseTeamBudget(
		ctx context.Context,
		tx *sql.Tx,
		teamID int64,
		amount int64,
	) error
	IncreaseTeamBudget(
		ctx context.Context,
		tx *sql.Tx,
		teamID int64,
		amount int64,
	) error
	MovePlayerToTeam(
		ctx context.Context,
		tx *sql.Tx,
		playerID int64,
		newTeamID int64,
		newMarketValue int64,
	) error
	MarkListingSold(
		ctx context.Context,
		tx *sql.Tx,
		listingID int64,
	) error
	CreateTransferHistory(
		ctx context.Context,
		tx *sql.Tx,
		listingID int64,
		playerID int64,
		sellerTeamID int64,
		buyerTeamID int64,
		salePrice int64,
		oldMarketValue int64,
		newMarketValue int64,
	) error
}

func NewBuyPlayerService(db *sql.DB) *BuyPlayerService {
	return NewBuyPlayerServiceWithRepository(db, infrastructure.NewTransferRepository(db))
}

func NewBuyPlayerServiceWithRepository(db *sql.DB, transferRepo TransferBuyerRepository) *BuyPlayerService {
	return &BuyPlayerService{
		db:           db,
		transferRepo: transferRepo,
	}
}

func (s *BuyPlayerService) Execute(
	ctx context.Context,
	buyerUserID int64,
	listingID int64,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	listing, err := s.transferRepo.FindActiveListingForUpdate(ctx, tx, listingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrListingNotFound
		}
		return err
	}

	buyerTeam, err := s.transferRepo.FindTeamByUserIDForUpdate(ctx, tx, buyerUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrListingNotFound
		}
		return err
	}

	if buyerTeam.ID == listing.SellerTeamID {
		return ErrCannotBuyOwnPlayer
	}

	if buyerTeam.Budget < listing.AskingPrice {
		return ErrInsufficientBudget
	}

	increasePercent := randomIncreasePercent()
	newMarketValue := listing.MarketValue + (listing.MarketValue * increasePercent / 100)

	if err := s.transferRepo.DecreaseTeamBudget(ctx, tx, buyerTeam.ID, listing.AskingPrice); err != nil {
		return err
	}

	if err := s.transferRepo.IncreaseTeamBudget(ctx, tx, listing.SellerTeamID, listing.AskingPrice); err != nil {
		return err
	}

	if err := s.transferRepo.MovePlayerToTeam(ctx, tx, listing.PlayerID, buyerTeam.ID, newMarketValue); err != nil {
		return err
	}

	if err := s.transferRepo.MarkListingSold(ctx, tx, listing.ID); err != nil {
		return err
	}

	if err := s.transferRepo.CreateTransferHistory(
		ctx,
		tx,
		listing.ID,
		listing.PlayerID,
		listing.SellerTeamID,
		buyerTeam.ID,
		listing.AskingPrice,
		listing.MarketValue,
		newMarketValue,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func randomIncreasePercent() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(91) + 10) // 10–100
}
