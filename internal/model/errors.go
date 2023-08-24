package model

import "errors"

var (
	ErrCardNotFound         = errors.New("card not found")
	ErrAmountMustBePositive = errors.New("amount must be positive")
	ErrCardDisabled         = errors.New("card is disabled")
)
