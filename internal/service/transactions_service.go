package service

import (
	"context"
	"time"

	"github.com/iancardoso/pismo/internal/domain"
	"github.com/iancardoso/pismo/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type TransactionsService struct {
	accountsRepo     repository.AccountsRepository
	transactionsRepo repository.TransactionsRepository
	now              func() time.Time
}

func NewTransactionsService(accountsRepo repository.AccountsRepository, transactionsRepo repository.TransactionsRepository) TransactionsService {
	return TransactionsService{
		accountsRepo:     accountsRepo,
		transactionsRepo: transactionsRepo,
		now:              time.Now,
	}
}

func (s TransactionsService) Create(ctx context.Context, accountID, operationTypeID int64, amount float64) (domain.Transaction, error) {
	_, span := otel.Tracer("github.com/iancardoso/pismo/internal/service").Start(ctx, "transactions-service.create")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("app.account.id", accountID),
		attribute.Int64("app.transaction.operation_type_id", operationTypeID),
	)
	if accountID <= 0 {
		span.SetStatus(codes.Error, ErrInvalidAccountID.Error())
		return domain.Transaction{}, ErrInvalidAccountID
	}
	if !domain.IsValidOperationType(operationTypeID) {
		span.SetStatus(codes.Error, ErrInvalidOperationType.Error())
		return domain.Transaction{}, ErrInvalidOperationType
	}
	if amount <= 0 {
		span.SetStatus(codes.Error, ErrInvalidAmount.Error())
		return domain.Transaction{}, ErrInvalidAmount
	}

	if _, err := s.accountsRepo.GetByID(ctx, accountID); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return domain.Transaction{}, err
	}

	transaction := domain.Transaction{
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          domain.NormalizeAmount(operationTypeID, amount),
		EventDate:       s.now().UTC(),
	}
	span.SetAttributes(attribute.Float64("app.transaction.amount", transaction.Amount))

	createdTransaction, err := s.transactionsRepo.Create(ctx, transaction)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return domain.Transaction{}, err
	}
	span.SetAttributes(attribute.Int64("app.transaction.id", createdTransaction.ID))
	return createdTransaction, nil
}
