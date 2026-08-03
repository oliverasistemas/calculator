package calculator

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"

	calc "calculator/services/api/pkg/calculator"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	var req calc.Request
	h.logger.Debug("calculate request received", "method", r.Method, "path", r.URL.Path)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	h.logger.Debug("performing calculation", "operation", req.Operation, "a", req.A, "b", req.B)

	result, err := calc.Calculate(req)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case errors.Is(err, calc.ErrDivisionByZero):
			status = http.StatusUnprocessableEntity
		case errors.Is(err, calc.ErrNegativeSqrt):
			status = http.StatusUnprocessableEntity
		case errors.Is(err, calc.ErrUnknownOperation):
			status = http.StatusBadRequest
		case errors.Is(err, calc.ErrMissingOperand):
			status = http.StatusBadRequest
		}
		h.logger.Warn("calculation failed", "operation", req.Operation, "status", status, "error", err)
		writeJSON(w, status, errorResponse{Error: err.Error()})
		return
	}

	h.logger.Debug("calculation succeeded", "operation", req.Operation, "result", result.Result)
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
