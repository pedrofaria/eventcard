package model

import "github.com/google/uuid"

type Balance struct {
	CardId uuid.UUID
	Amount float32
}
