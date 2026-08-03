package router

import (
	"calculator/services/api/pkg/controller/calculator"
	"log/slog"

	"github.com/go-chi/chi/v5"
)

func SetupAPIRoutes(r chi.Router, logger *slog.Logger) {
	calcHandler := calculator.NewHandler(logger)

	r.Post("/calculate", calcHandler.Calculate)
}
