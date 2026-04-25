package application

import (
	"context"
	"database/sql"

	"github.com/sergiojaa/soccer-manager-api/internal/transfers/infrastructure"
)

type ListMarketService struct {
	transferRepo MarketReaderRepository
}

type MarketReaderRepository interface {
	FindActiveListings(ctx context.Context) ([]infrastructure.MarketListingView, error)
}

func NewListMarketService(db *sql.DB) *ListMarketService {
	return NewListMarketServiceWithRepository(infrastructure.NewTransferRepository(db))
}

func NewListMarketServiceWithRepository(transferRepo MarketReaderRepository) *ListMarketService {
	return &ListMarketService{
		transferRepo: transferRepo,
	}
}

func (s *ListMarketService) Execute(
	ctx context.Context,
) ([]infrastructure.MarketListingView, error) {
	return s.transferRepo.FindActiveListings(ctx)
}
