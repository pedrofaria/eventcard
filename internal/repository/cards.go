package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pedrofaria/eventcard/internal/repository/cards"
)

type Cards struct {
	cards.Queries

	db *sql.DB
}

func NewCards(db *sql.DB) *Cards {
	cQ := cards.New(db)

	return &Cards{
		Queries: *cQ,
		db:      db,
	}
}

func (c *Cards) RunAtomic(ctx context.Context, fn func(repo *cards.Queries) error) error {
	tx, _ := c.db.BeginTx(ctx, nil)

	if err := fn(c.WithTx(tx)); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *Cards) GetCardFullFromExternalId(ctx context.Context, externalId int32) (*cards.GetCardFullByExternalIdRow, error) {
	card, err := c.GetCardFullByExternalId(ctx, externalId)
	if err != nil {
		return nil, err
	}

	return &card, err
}

func (c *Cards) CreateCard(ctx context.Context, externalId int32, name string, enabled bool) (*cards.GetCardFullRow, error) {
	cardId := uuid.New()
	now := time.Now().UTC()
	balance := "0.00"

	err := c.RunAtomic(ctx, func(repo *cards.Queries) error {
		cardParam := cards.CreateCardParams{
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

		balanceParam := cards.CreateBalanceParams{
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

	return &cards.GetCardFullRow{
		ID:         cardId,
		ExternalID: externalId,
		Name:       name,
		Enabled:    enabled,
		CreatedAt:  now,
		UpdatedAt:  now,
		Balance:    balance,
	}, nil
}

func (c *Cards) UpdateEnableCard(ctx context.Context, externalId int32, enabled bool) error {
	return c.UpdateEnabledCardByExternalId(
		ctx,
		cards.UpdateEnabledCardByExternalIdParams{ExternalID: externalId, Enabled: enabled})
}

func registerLedger(
	ctx context.Context,
	repo *cards.Queries,
	cardId uuid.UUID,
	reference cards.Reference,
	referenceId uuid.UUID,
	amount string,
	createdAt time.Time,
) error {
	if err := repo.CreateLedger(ctx, cards.CreateLedgerParams{
		ID:          uuid.New(),
		CardID:      cardId,
		Reference:   cards.ReferenceDeposit,
		ReferenceID: referenceId,
		Amount:      amount,
		CreatedAt:   createdAt,
	}); err != nil {
		return err
	}

	return repo.IncreaseCardBalance(ctx, cards.IncreaseCardBalanceParams{
		Amount: amount,
		CardID: cardId,
	})
}
