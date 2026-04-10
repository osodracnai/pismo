package service

import (
	"context"
	"testing"

	"github.com/iancardoso/pismo/internal/domain"
)

type accountsRepoStub struct {
	exists  bool
	account domain.Account
}

func (s accountsRepoStub) Create(_ context.Context, documentNumber string) (domain.Account, error) {
	return domain.Account{ID: 1, DocumentNumber: documentNumber}, nil
}

func (s accountsRepoStub) GetByID(_ context.Context, id int64) (domain.Account, error) {
	if id == s.account.ID {
		return s.account, nil
	}
	return domain.Account{}, ErrAccountNotFound
}

func (s accountsRepoStub) ExistsByDocumentNumber(_ context.Context, _ string) (bool, error) {
	return s.exists, nil
}

func TestAccountsServiceCreateRejectsBlankDocument(t *testing.T) {
	service := NewAccountsService(accountsRepoStub{})

	_, err := service.Create(context.Background(), "   ")
	if err != ErrInvalidDocumentNumber {
		t.Fatalf("expected ErrInvalidDocumentNumber, got %v", err)
	}
}

func TestAccountsServiceCreateRejectsDuplicateDocument(t *testing.T) {
	service := NewAccountsService(accountsRepoStub{exists: true})

	_, err := service.Create(context.Background(), "12345678900")
	if err != ErrDuplicateDocument {
		t.Fatalf("expected ErrDuplicateDocument, got %v", err)
	}
}

func TestAccountsServiceGetByIDRejectsInvalidID(t *testing.T) {
	service := NewAccountsService(accountsRepoStub{})

	_, err := service.GetByID(context.Background(), 0)
	if err != ErrInvalidAccountID {
		t.Fatalf("expected ErrInvalidAccountID, got %v", err)
	}
}
