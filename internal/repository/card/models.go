// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package card

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID         uuid.UUID
	ExternalID int32
	Name       string
	Enabled    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
