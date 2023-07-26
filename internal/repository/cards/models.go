// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package cards

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Reference string

const (
	ReferenceDeposit  Reference = "deposit"
	ReferenceTransfer Reference = "transfer"
	ReferencePurchase Reference = "purchase"
)

func (e *Reference) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Reference(s)
	case string:
		*e = Reference(s)
	default:
		return fmt.Errorf("unsupported scan type for Reference: %T", src)
	}
	return nil
}

type NullReference struct {
	Reference Reference
	Valid     bool // Valid is true if Reference is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullReference) Scan(value interface{}) error {
	if value == nil {
		ns.Reference, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Reference.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullReference) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Reference), nil
}

type Card struct {
	ID         uuid.UUID
	ExternalID int32
	Name       string
	Enabled    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Deposit struct {
	ID         uuid.UUID
	ExternalID uuid.UUID
	CardID     uuid.UUID
	Amount     string
	Paid       bool
	Cancelled  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
