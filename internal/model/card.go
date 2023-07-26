package model

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	Id         uuid.UUID
	ExternalId int32
	Name       string
	Enabled    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Balance    *Balance
}
