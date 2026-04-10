package service

import "errors"

var (
	ErrInvalidDocumentNumber = errors.New("document_number is required")
	ErrDuplicateDocument     = errors.New("document_number already exists")
	ErrInvalidAccountID      = errors.New("invalid account_id")
	ErrAccountNotFound       = errors.New("account not found")
	ErrInvalidOperationType  = errors.New("invalid operation_type_id")
	ErrInvalidAmount         = errors.New("amount must be greater than zero")
	ErrInsufficientFunds     = errors.New("insufficient funds")
)
