package http

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	_ "github.com/iancardoso/pismo/docs"
	docs "github.com/iancardoso/pismo/docs"
)

func NewRouter(accountsHandler AccountsHandler, transactionsHandler TransactionsHandler) stdhttp.Handler {
	router := chi.NewRouter()

	router.Use(otelhttp.NewMiddleware("http-server"))

	router.Get("/health", func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	router.Get("/openapi.yaml", func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write(docs.SwaggerYAML)
	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	router.Post("/accounts", otelhttp.WithRouteTag("/accounts", stdhttp.HandlerFunc(accountsHandler.Create)).ServeHTTP)
	router.Get("/accounts/{accountId}", otelhttp.WithRouteTag("/accounts/{accountId}", stdhttp.HandlerFunc(accountsHandler.GetByID)).ServeHTTP)
	router.Post("/transactions", otelhttp.WithRouteTag("/transactions", stdhttp.HandlerFunc(transactionsHandler.Create)).ServeHTTP)

	return router
}
