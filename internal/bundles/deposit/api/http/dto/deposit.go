package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/pedrofaria/eventcard/internal/bundles/deposit/model"
)

type Deposit struct {
	ExternalId     uuid.UUID `json:"externalId"`
	ExternalCardId uint32    `json:"externalCardId"`
	Amount         float32   `json:"amount"`
	Paid           bool      `json:"paid"`
	Cancelled      bool      `json:"cancelled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func FromDepositModel(m *model.Deposit) Deposit {
	return Deposit{
		ExternalId:     m.ExternalId,
		ExternalCardId: m.ExternalCardId,
		Amount:         m.Amount,
		Paid:           m.Paid,
		Cancelled:      m.Cancelled,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
