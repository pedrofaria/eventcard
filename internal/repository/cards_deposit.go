package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pedrofaria/eventcard/internal/repository/cards"
)

func (c *Cards) CreateDeposit(ctx context.Context, cardId uuid.UUID, amount float32, paid bool) (*cards.Deposit, error) {
	id := uuid.New()
	externalId := uuid.New()
	now := time.Now().UTC()
	amountStr := fmt.Sprintf("%f", amount)

	err := c.RunAtomic(ctx, func(repo *cards.Queries) error {
		if err := repo.CreateDeposit(ctx, cards.CreateDepositParams{
			ID:         id,
			ExternalID: externalId,
			CardID:     cardId,
			Amount:     amountStr,
			Paid:       paid,
			CreatedAt:  now,
			UpdatedAt:  now,
		}); err != nil {
			return err
		}

		return registerLedger(ctx, repo, cardId, cards.ReferenceDeposit, id, amountStr, now)
	})

	if err != nil {
		return nil, err
	}

	return &cards.Deposit{
		ID:         id,
		ExternalID: externalId,
		CardID:     cardId,
		Amount:     amountStr,
		Paid:       paid,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (c *Cards) CancelDeposit(ctx context.Context, externalId uuid.UUID) error {
	now := time.Now().UTC()

	return c.RunAtomic(ctx, func(repo *cards.Queries) error {
		deposit, err := repo.CancelDepositByExternalId(ctx, externalId)
		if err != nil {
			return err
		}

		return registerLedger(ctx, repo, deposit.CardID, cards.ReferenceDeposit, deposit.ID, "-"+deposit.Amount, now)
	})
}
