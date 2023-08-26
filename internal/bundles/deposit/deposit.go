package deposit

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/pedrofaria/eventcard/internal/bundles/deposit/model"
	"github.com/pedrofaria/eventcard/internal/bundles/deposit/repository"
	mainModel "github.com/pedrofaria/eventcard/internal/model"
	"github.com/pedrofaria/eventcard/internal/service"
	"github.com/pedrofaria/eventcard/internal/utils"
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

func (s *Service) List(ctx context.Context, externalCardId uint32) ([]model.Deposit, error) {
	card, err := s.card.GetCardFull(ctx, int32(externalCardId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mainModel.ErrCardNotFound
		}

		return nil, err
	}

	dd, err := s.repo.GetListByCardId(ctx, card.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("card deposit list is empty")
		}

		return nil, err
	}

	models := make([]model.Deposit, len(dd))
	for i, d := range dd {
		models[i] = repository.DepositToModel(d, externalCardId)
	}

	return models, nil
}

func (s *Service) CreateDeposit(ctx context.Context, externalCardId uint32, amount float32, paid bool) (*model.Deposit, error) {
	if amount < 0 {
		return nil, mainModel.ErrAmountMustBePositive
	}

	card, err := s.card.GetCardFull(ctx, int32(externalCardId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mainModel.ErrCardNotFound
		}

		return nil, err
	}

	if !card.Enabled {
		return nil, mainModel.ErrCardDisabled
	}

	repoDeposit, err := s.repo.CreateDeposit(ctx, card.Id, amount, paid)
	if err != nil {
		return nil, err
	}

	model := repository.DepositToModel(repoDeposit, externalCardId)

	return &model, nil
}

func (s *Service) TogglePaidStatus(ctx context.Context, externalCardId uint32, externalId uuid.UUID) (*model.Deposit, error) {
	card, err := s.card.GetCardFull(ctx, int32(externalCardId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mainModel.ErrCardNotFound
		}

		return nil, err
	}

	d, err := s.repo.GetByCardAndExternalId(ctx, card.Id, externalId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("card deposit not found")
		}

		return nil, err
	}

	d, err = s.repo.UpdatePaid(ctx, d.ID, !d.Paid)
	if err != nil {
		return nil, err
	}

	model := repository.DepositToModel(d, externalCardId)

	return &model, nil
}

func (s *Service) Cancel(ctx context.Context, externalCardId uint32, externalId uuid.UUID) error {
	card, err := s.card.GetCardFull(ctx, int32(externalCardId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mainModel.ErrCardNotFound
		}

		return err
	}

	d, err := s.repo.GetByCardAndExternalId(ctx, card.Id, externalId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("card deposit not found")
		}

		return err
	}

	amount := utils.NumericToFloat32(d.Amount)

	if amount > card.Balance.Amount {
		return errors.New("not enough balance to cancel deposit")
	}

	return s.repo.CancelDeposit(ctx, d.ID)
}
