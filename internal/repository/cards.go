package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pedrofaria/eventcard/internal/repository/card"
	"github.com/pedrofaria/eventcard/internal/repository/ledger"
)

type Cards struct {
	*card.Queries
	db *sql.DB
}

func NewCards(db *sql.DB) *Cards {
	return &Cards{
		Queries: card.New(db),
		db:      db,
	}
}

func (c *Cards) RunAtomic(ctx context.Context, fn func(repo *card.Queries) error) error {
	tx, _ := c.db.BeginTx(ctx, nil)

	if err := fn(c.WithTx(tx)); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *Cards) CreateCard(ctx context.Context, externalId int32, name string, enabled bool) (*card.GetCardFullRow, error) {
	cardId := uuid.New()
	now := time.Now().UTC()
	balance := "0.00"

	err := c.RunAtomic(ctx, func(repo *card.Queries) error {
		cardParam := card.CreateCardParams{
			ID:         cardId,
			ExternalID: externalId,
			Name:       name,
			Enabled:    enabled,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		if err := repo.CreateCard(ctx, cardParam); err != nil {
			return err
		}

		balanceParam := card.CreateBalanceParams{
			ID:        uuid.New(),
			CardID:    cardId,
			Amount:    balance,
			CreatedAt: now,
			UpdatedAt: now,
		}
		return repo.CreateBalance(ctx, balanceParam)
	})

	if err != nil {
		return nil, err
	}

	return &card.GetCardFullRow{
		ID:         cardId,
		ExternalID: externalId,
		Name:       name,
		Enabled:    enabled,
		CreatedAt:  now,
		UpdatedAt:  now,
		Balance:    balance,
	}, nil
}

func registerLedger(
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
