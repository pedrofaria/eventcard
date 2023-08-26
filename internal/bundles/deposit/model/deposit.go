package model

import (
	"time"

	"github.com/google/uuid"
)

type Deposit struct {
	Id             uuid.UUID
	ExternalId     uuid.UUID
	CardId         uuid.UUID
	ExternalCardId uint32
	Amount         float32
	Paid           bool
	Cancelled      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
