package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iancardoso/pismo/internal/repository/sqlite"
	"github.com/iancardoso/pismo/internal/service"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()

	db, err := sqlite.Open("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	accountsRepository := sqlite.NewAccountsRepository(db)
	transactionsRepository := sqlite.NewTransactionsRepository(db)
	accountsService := service.NewAccountsService(accountsRepository)
	transactionsService := service.NewTransactionsService(accountsRepository, transactionsRepository)

	return NewRouter(NewAccountsHandler(accountsService), NewTransactionsHandler(transactionsService))
}

func TestAccountsLifecycle(t *testing.T) {
	router := newTestRouter(t)

	createRequest := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(`{"document_number":"12345678900"}`))
	createRecorder := httptest.NewRecorder()
	router.ServeHTTP(createRecorder, createRequest)

	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createRecorder.Code)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRecorder.Code)
	}
}

func TestCreateTransactionReturnsCreated(t *testing.T) {
	router := newTestRouter(t)

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(`{"document_number":"12345678900"}`)))

	transactionRequest := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(`{"account_id":1,"operation_type_id":4,"amount":123.45}`))
	transactionRecorder := httptest.NewRecorder()
	router.ServeHTTP(transactionRecorder, transactionRequest)

	if transactionRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, transactionRecorder.Code)
	}
}

func TestCreateTransactionRejectsMissingAccount(t *testing.T) {
	router := newTestRouter(t)

	transactionRequest := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(`{"account_id":99,"operation_type_id":4,"amount":123.45}`))
	transactionRecorder := httptest.NewRecorder()
	router.ServeHTTP(transactionRecorder, transactionRequest)

	if transactionRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, transactionRecorder.Code)
	}
}

func TestCreateTransactionRejectsInsufficientFunds(t *testing.T) {
	router := newTestRouter(t)

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(`{"document_number":"12345678900"}`)))
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(`{"account_id":1,"operation_type_id":4,"amount":100}`)))

	transactionRequest := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(`{"account_id":1,"operation_type_id":1,"amount":150}`))
	transactionRecorder := httptest.NewRecorder()
	router.ServeHTTP(transactionRecorder, transactionRequest)

	if transactionRecorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, transactionRecorder.Code)
	}
}

func TestOpenAPIRouteReturnsSpec(t *testing.T) {
	router := newTestRouter(t)

	request := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if body := recorder.Body.String(); body == "" {
		t.Fatal("expected non-empty openapi body")
	}
}

func TestSwaggerRouteReturnsDocs(t *testing.T) {
	router := newTestRouter(t)

	request := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if body := recorder.Body.String(); body == "" {
		t.Fatal("expected non-empty swagger body")
	}
}
