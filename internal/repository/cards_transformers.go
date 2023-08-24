package repository

import (
	"github.com/pedrofaria/eventcard/internal/model"
	"github.com/pedrofaria/eventcard/internal/repository/card"
	"github.com/pedrofaria/eventcard/internal/utils"
)

func GetCardFullRowToModel(repoCard *card.GetCardFullRow) model.Card {
	return model.Card{
		Id:         repoCard.ID,
		ExternalId: repoCard.ExternalID,
		Name:       repoCard.Name,
		Enabled:    repoCard.Enabled,
		CreatedAt:  repoCard.CreatedAt,
		UpdatedAt:  repoCard.UpdatedAt,
		Balance: &model.Balance{
			CardId: repoCard.ID,
			Amount: utils.NumericToFloat32(repoCard.Balance),
		},
	}
}

func GetCardFullByExternalIdRowToModel(repoCard *card.GetCardFullByExternalIdRow) model.Card {
	return model.Card{
		Id:         repoCard.ID,
		ExternalId: repoCard.ExternalID,
		Name:       repoCard.Name,
		Enabled:    repoCard.Enabled,
		CreatedAt:  repoCard.CreatedAt,
		UpdatedAt:  repoCard.UpdatedAt,
		Balance: &model.Balance{
			CardId: repoCard.ID,
			Amount: utils.NumericToFloat32(repoCard.Balance),
		},
	}
}
