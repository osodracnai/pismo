package repository

import (
	"context"

	"github.com/iancardoso/pismo/internal/domain"
)

type TransactionsRepository interface {
	Create(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error)
}
