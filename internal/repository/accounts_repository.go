package repository

import (
	"context"

	"github.com/iancardoso/pismo/internal/domain"
)

type AccountsRepository interface {
	Create(ctx context.Context, documentNumber string) (domain.Account, error)
	GetByID(ctx context.Context, id int64) (domain.Account, error)
	ExistsByDocumentNumber(ctx context.Context, documentNumber string) (bool, error)
}
