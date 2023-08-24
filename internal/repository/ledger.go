package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pedrofaria/eventcard/internal/repository/ledger"
)

func RegisterLedger(
	ctx context.Context,
	dbtx ledger.DBTX,
	repo *ledger.Queries,
	cardId uuid.UUID,
	reference ledger.Reference,
	referenceId uuid.UUID,
	amount string,
	createdAt time.Time,
) error {
	if err := repo.CreateLedger(ctx, dbtx, ledger.CreateLedgerParams{
		ID:          uuid.New(),
		CardID:      cardId,
		Reference:   ledger.ReferenceDeposit,
		ReferenceID: referenceId,
		Amount:      amount,
		CreatedAt:   createdAt,
	}); err != nil {
		return err
	}

	return repo.IncreaseCardBalance(ctx, dbtx, amount, cardId)
}
