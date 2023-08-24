package deposit

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pedrofaria/eventcard/internal/bundles/deposit/model"
	"github.com/pedrofaria/eventcard/internal/bundles/deposit/repository"
	mainModel "github.com/pedrofaria/eventcard/internal/model"
	"github.com/pedrofaria/eventcard/internal/service"
)

type Service struct {
	repo *repository.Deposit
	card *service.Card
}

func Init(db *sql.DB, card *service.Card) *Service {
	return NewService(repository.NewDeposit(db), card)
}

func NewService(repo *repository.Deposit, card *service.Card) *Service {
	return &Service{
		repo: repo,
		card: card,
	}
}

func (c *Service) CreateDeposit(ctx context.Context, externalCardId uint32, amount float32, paid bool) (*model.Deposit, error) {
	if amount < 0 {
		return nil, mainModel.ErrAmountMustBePositive
	}

	card, err := c.card.GetCardFull(ctx, int32(externalCardId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mainModel.ErrCardNotFound
		}

		return nil, err
	}

	if !card.Enabled {
		return nil, mainModel.ErrCardDisabled
	}

	repoDeposit, err := c.repo.CreateDeposit(ctx, card.Id, amount, paid)
	if err != nil {
		return nil, err
	}

	model := repository.DepositToModel(repoDeposit)

	return &model, nil
}
