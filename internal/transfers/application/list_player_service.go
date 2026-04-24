package application

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/sergiojaa/soccer-manager-api/internal/transfers/infrastructure"
)

var (
	ErrTransferPlayerNotFoundOrNotOwned = errors.New("player not found or not owned by user")
	ErrInvalidAskingPrice               = errors.New("asking price must be greater than zero")
	ErrPlayerAlreadyListed              = errors.New("player is already listed for transfer")
)

type ListPlayerService struct {
	transferRepo *infrastructure.TransferRepository
}

func NewListPlayerService(db *sql.DB) *ListPlayerService {
	return &ListPlayerService{
		transferRepo: infrastructure.NewTransferRepository(db),
	}
}

func (s *ListPlayerService) Execute(
	ctx context.Context,
	userID int64,
	playerID int64,
	askingPrice int64,
) error {
	if askingPrice <= 0 {
		return ErrInvalidAskingPrice
	}

	err := s.transferRepo.ListOwnedPlayer(ctx, userID, playerID, askingPrice)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTransferPlayerNotFoundOrNotOwned
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrPlayerAlreadyListed
		}

		return err
	}

	return nil
}
