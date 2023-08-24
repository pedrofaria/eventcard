package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pedrofaria/eventcard/internal/bundles/deposit/model"
	"github.com/pedrofaria/eventcard/internal/bundles/deposit/repository/sqlc"
	"github.com/pedrofaria/eventcard/internal/repository"
	"github.com/pedrofaria/eventcard/internal/repository/ledger"
	"github.com/pedrofaria/eventcard/internal/utils"
)

type Deposit struct {
	deposit *sqlc.Queries
	ledger  *ledger.Queries
	db      *sql.DB
}

func NewDeposit(db *sql.DB) *Deposit {
	return &Deposit{
		deposit: sqlc.New(db),
		ledger:  ledger.New(),
		db:      db,
	}
}

func (c *Deposit) RunAtomic(ctx context.Context, fn func(ledger.DBTX, *sqlc.Queries, *ledger.Queries) error) error {
	tx, _ := c.db.BeginTx(ctx, nil)

	if err := fn(tx, c.deposit.WithTx(tx), c.ledger); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c *Deposit) CreateDeposit(ctx context.Context, cardId uuid.UUID, amount float32, paid bool) (*sqlc.Deposit, error) {
	id := uuid.New()
	externalId := uuid.New()
	now := time.Now().UTC()
	amountStr := fmt.Sprintf("%f", amount)

	err := c.RunAtomic(ctx, func(dbtx ledger.DBTX, repo *sqlc.Queries, l *ledger.Queries) error {
		if err := repo.CreateDeposit(ctx, sqlc.CreateDepositParams{
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

		return repository.RegisterLedger(ctx, dbtx, l, cardId, ledger.ReferenceDeposit, id, amountStr, now)
	})

	if err != nil {
		return nil, err
	}

	return &sqlc.Deposit{
		ID:         id,
		ExternalID: externalId,
		CardID:     cardId,
		Amount:     amountStr,
		Paid:       paid,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (c *Deposit) CancelDeposit(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()

	return c.RunAtomic(ctx, func(dbtx ledger.DBTX, repo *sqlc.Queries, l *ledger.Queries) error {
		deposit, err := repo.CancelDepositById(ctx, id)
		if err != nil {
			return err
		}

		return repository.RegisterLedger(ctx, dbtx, l, deposit.CardID, ledger.ReferenceDeposit, deposit.ID, "-"+deposit.Amount, now)
	})
}

func DepositToModel(dep *sqlc.Deposit) model.Deposit {
	return model.Deposit{
		Id:         dep.ID,
		ExternalId: dep.ExternalID,
		CardId:     dep.CardID,
		Amount:     utils.NumericToFloat32(dep.Amount),
		Paid:       dep.Paid,
		CreatedAt:  dep.CreatedAt,
		UpdatedAt:  dep.UpdatedAt,
	}
}
