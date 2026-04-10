package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/iancardoso/pismo/docs"
	apihttp "github.com/iancardoso/pismo/internal/http"
	"github.com/iancardoso/pismo/internal/repository/sqlite"
	"github.com/iancardoso/pismo/internal/service"
	"github.com/iancardoso/pismo/internal/telemetry"
)

// @title Pismo Assessment API
// @version 1.0
// @description REST API for the Pismo code assessment.
// @BasePath /
func main() {
	ctx := context.Background()
	shutdownTelemetry, err := telemetry.Setup(ctx)
	if err != nil {
		log.Fatalf("telemetry setup failed: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTelemetry(shutdownCtx); err != nil {
			log.Printf("telemetry shutdown failed: %v", err)
		}
	}()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "file:pismo.db?_pragma=foreign_keys(1)"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := sqlite.Open(databaseURL)
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	defer db.Close()

	accountsRepository := sqlite.NewAccountsRepository(db)
	transactionsRepository := sqlite.NewTransactionsRepository(db)

	accountsService := service.NewAccountsService(accountsRepository)
	transactionsService := service.NewTransactionsService(accountsRepository, transactionsRepository)

	server := &http.Server{
		Addr: ":" + port,
		Handler: apihttp.NewRouter(
			apihttp.NewAccountsHandler(accountsService),
			apihttp.NewTransactionsHandler(transactionsService),
		),
	}

	go func() {
		log.Printf("listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}
