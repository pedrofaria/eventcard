package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/pedrofaria/eventcard/internal/model"
	"github.com/pedrofaria/eventcard/internal/repository"
)

var (
	ErrCardNotFound         = errors.New("card not found")
	ErrAmountMustBePositive = errors.New("amount must be positive")
	ErrCardDisabled         = errors.New("card is disabled")
)

type Card struct {
	repo *repository.Cards
}

func NewCard(repo *repository.Cards) *Card {
	return &Card{repo: repo}
}

func (c *Card) GetCardFull(ctx context.Context, externalId int32) (*model.Card, error) {
	repoCard, err := c.repo.GetCardFullByExternalId(ctx, externalId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCardNotFound
		}

		return nil, err
	}

	model := repository.GetCardFullByExternalIdRowToModel(repoCard)

	return &model, nil
}

func (c *Card) GetCardBalance(ctx context.Context, externalId int32) (*model.Balance, error) {
	cardId, err := c.repo.GetCardIdByExternalId(ctx, externalId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCardNotFound
		}

		return nil, err
	}

	repoAmount, err := c.repo.GetCardBalance(ctx, cardId)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseFloat(repoAmount, 32)
	if err != nil {
		return nil, err
	}

	return &model.Balance{
		CardId: cardId,
		Amount: float32(amount),
	}, nil
}

func (c *Card) CreateCard(ctx context.Context, externalId uint32, name string, enabled bool) (*model.Card, error) {
	repoCard, err := c.repo.CreateCard(ctx, int32(externalId), name, enabled)
	if err != nil {
		return nil, err
	}

	model := repository.GetCardFullRowToModel(repoCard)

	return &model, nil
}
