package calculator

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateEndpoint(t *testing.T) {
	h := NewHandler(slog.New(slog.NewTextHandler(io.Discard, nil)))

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantResult *float64
		wantExpr   string
		wantError  string
	}{
		{
			name:       "add",
			body:       `{"operation":"add","a":5,"b":3}`,
			wantStatus: http.StatusOK,
			wantResult: ptr(8),
			wantExpr:   "5 + 3",
		},
		{
			name:       "divide by zero",
			body:       `{"operation":"divide","a":10,"b":0}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "division by zero",
		},
		{
			name:       "invalid json",
			body:       `not json`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid JSON body",
		},
		{
			name:       "unknown operation",
			body:       `{"operation":"modulo","a":10,"b":3}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "unknown operation",
		},
		{
			name:       "missing operand",
			body:       `{"operation":"add","a":5}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "missing required operand",
		},
		{
			name:       "sqrt negative",
			body:       `{"operation":"sqrt","a":-4}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "square root of negative number",
		},
		{
			name:       "sqrt success",
			body:       `{"operation":"sqrt","a":25}`,
			wantStatus: http.StatusOK,
			wantResult: ptr(5),
			wantExpr:   "√25",
		},
		{
			name:       "percentage",
			body:       `{"operation":"percentage","a":10,"b":200}`,
			wantStatus: http.StatusOK,
			wantResult: ptr(20),
			wantExpr:   "10% of 200",
		},
		{
			name:       "power producing NaN",
			body:       `{"operation":"power","a":-1,"b":0.5}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "result is not a finite number",
		},
		{
			name:       "multiply overflowing to Inf",
			body:       `{"operation":"multiply","a":1e308,"b":10}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "result is not a finite number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.Calculate(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantResult != nil {
				var resp struct {
					Result     float64 `json:"result"`
					Expression string  `json:"expression"`
				}
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.Result != *tt.wantResult {
					t.Errorf("result = %v, want %v", resp.Result, *tt.wantResult)
				}
				if resp.Expression != tt.wantExpr {
					t.Errorf("expression = %q, want %q", resp.Expression, tt.wantExpr)
				}
			}

			if tt.wantError != "" {
				var resp errorResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode error response: %v", err)
				}
				if resp.Error != tt.wantError {
					t.Errorf("error = %q, want %q", resp.Error, tt.wantError)
				}
			}
		})
	}
}

func ptr(f float64) *float64 { return &f }
