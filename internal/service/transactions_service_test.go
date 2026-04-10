package service

import (
	"context"
	"testing"
	"time"

	"github.com/iancardoso/pismo/internal/domain"
)

type transactionsRepoStub struct {
	created domain.Transaction
}

func (s *transactionsRepoStub) Create(_ context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	s.created = transaction
	transaction.ID = 1
	return transaction, nil
}

func TestTransactionsServiceCreateNormalizesSigns(t *testing.T) {
	accountsRepo := accountsRepoStub{account: domain.Account{ID: 1, DocumentNumber: "123"}}
	transactionsRepo := &transactionsRepoStub{}
	service := NewTransactionsService(accountsRepo, transactionsRepo)
	service.now = func() time.Time { return time.Date(2026, 4, 6, 12, 0, 0, 0, time.UTC) }

	tests := []struct {
		name            string
		operationTypeID int64
		amount          float64
		expected        float64
	}{
		{name: "purchase", operationTypeID: domain.OperationTypePurchase, amount: 50, expected: -50},
		{name: "installment purchase", operationTypeID: domain.OperationTypeInstallmentPurchase, amount: 50, expected: -50},
		{name: "withdrawal", operationTypeID: domain.OperationTypeWithdrawal, amount: 50, expected: -50},
		{name: "payment", operationTypeID: domain.OperationTypePayment, amount: 50, expected: 50},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transaction, err := service.Create(context.Background(), 1, test.operationTypeID, test.amount)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if transaction.Amount != test.expected {
				t.Fatalf("expected amount %v, got %v", test.expected, transaction.Amount)
			}
		})
	}
}

func TestTransactionsServiceCreateRejectsInvalidOperationType(t *testing.T) {
	accountsRepo := accountsRepoStub{account: domain.Account{ID: 1, DocumentNumber: "123"}}
	transactionsRepo := &transactionsRepoStub{}
	service := NewTransactionsService(accountsRepo, transactionsRepo)

	_, err := service.Create(context.Background(), 1, 9, 10)
	if err != ErrInvalidOperationType {
		t.Fatalf("expected ErrInvalidOperationType, got %v", err)
	}
}

func TestTransactionsServiceCreateRejectsMissingAccount(t *testing.T) {
	transactionsRepo := &transactionsRepoStub{}
	service := NewTransactionsService(accountsRepoStub{}, transactionsRepo)

	_, err := service.Create(context.Background(), 1, domain.OperationTypePayment, 10)
	if err != ErrAccountNotFound {
		t.Fatalf("expected ErrAccountNotFound, got %v", err)
	}
}
