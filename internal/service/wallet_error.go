package service

import "errors"

var (
	ErrNegativeAmount      = errors.New("amount cannot be negative")
	ErrZeroAmount          = errors.New("amount cannot be zero")
	ErrNonExistentCurrency = errors.New("currency does not exist")
)
