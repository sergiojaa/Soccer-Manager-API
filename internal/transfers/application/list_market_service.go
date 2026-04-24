package application

import (
	"context"
	"database/sql"

	"github.com/sergiojaa/soccer-manager-api/internal/transfers/infrastructure"
)

type ListMarketService struct {
	transferRepo *infrastructure.TransferRepository
}

func NewListMarketService(db *sql.DB) *ListMarketService {
	return &ListMarketService{
		transferRepo: infrastructure.NewTransferRepository(db),
	}
}

func (s *ListMarketService) Execute(
	ctx context.Context,
) ([]infrastructure.MarketListingView, error) {
	return s.transferRepo.FindActiveListings(ctx)
}
