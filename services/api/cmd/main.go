package main

import (
	"calculator/services/api/pkg/logger"
	"calculator/services/api/pkg/router"
	"calculator/services/api/pkg/types"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	httputil "calculator/services/api/pkg/utils/http"
)

func main() {
	cfg, err := types.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg)
	slog.SetDefault(log)

	r := chi.NewRouter()

	r.Use(httputil.RequestLogger(log))

	httputil.ApplyCORSOptions(r)

	router.SetupAPIRoutes(r, log)

	log.Info("calculator-api listening", "addr", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, r); err != nil {
		log.Error("server failed", "error", err)
		os.Exit(1)
	}
}
