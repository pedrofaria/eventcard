package model

import (
	"time"

	"github.com/google/uuid"
)

type Deposit struct {
	Id         uuid.UUID
	ExternalId uuid.UUID
	CardId     uuid.UUID
	Amount     float32
	Paid       bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
