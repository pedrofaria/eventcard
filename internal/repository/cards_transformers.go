package repository

import (
	"strconv"

	"github.com/pedrofaria/eventcard/internal/model"
	"github.com/pedrofaria/eventcard/internal/repository/cards"
)

func NumericToFloat32(v string) float32 {
	f, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return 0
	}

	return float32(f)
}

func GetCardFullRowToModel(repoCard *cards.GetCardFullRow) model.Card {
	return model.Card{
		Id:         repoCard.ID,
		ExternalId: repoCard.ExternalID,
		Name:       repoCard.Name,
		Enabled:    repoCard.Enabled,
		CreatedAt:  repoCard.CreatedAt,
		UpdatedAt:  repoCard.UpdatedAt,
		Balance: &model.Balance{
			CardId: repoCard.ID,
			Amount: NumericToFloat32(repoCard.Balance),
		},
	}
}

func GetCardFullByExternalIdRowToModel(repoCard *cards.GetCardFullByExternalIdRow) model.Card {
	return model.Card{
		Id:         repoCard.ID,
		ExternalId: repoCard.ExternalID,
		Name:       repoCard.Name,
		Enabled:    repoCard.Enabled,
		CreatedAt:  repoCard.CreatedAt,
		UpdatedAt:  repoCard.UpdatedAt,
		Balance: &model.Balance{
			CardId: repoCard.ID,
			Amount: NumericToFloat32(repoCard.Balance),
		},
	}
}

func DepositToModel(deposit *cards.Deposit) model.Deposit {
	return model.Deposit{
		Id:         deposit.ID,
		ExternalId: deposit.ExternalID,
		CardId:     deposit.CardID,
		Amount:     NumericToFloat32(deposit.Amount),
		Paid:       deposit.Paid,
		CreatedAt:  deposit.CreatedAt,
		UpdatedAt:  deposit.UpdatedAt,
	}
}
