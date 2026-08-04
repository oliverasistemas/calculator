package calculator

import (
	"errors"
	"math"
	"testing"
)

func ptr(f float64) *float64 { return &f }

func TestCalculate(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		want float64
	}{
		{"add positive", Request{Operation: Add, A: ptr(2), B: ptr(3)}, 5},
		{"add negative", Request{Operation: Add, A: ptr(-2), B: ptr(-3)}, -5},
		{"add mixed", Request{Operation: Add, A: ptr(-2), B: ptr(3)}, 1},
		{"add zero", Request{Operation: Add, A: ptr(0), B: ptr(0)}, 0},
		{"add decimals", Request{Operation: Add, A: ptr(0.1), B: ptr(0.2)}, 0.30000000000000004},
		{"subtract", Request{Operation: Subtract, A: ptr(10), B: ptr(4)}, 6},
		{"multiply", Request{Operation: Multiply, A: ptr(3), B: ptr(7)}, 21},
		{"multiply large", Request{Operation: Multiply, A: ptr(1e15), B: ptr(1e15)}, 1e30},
		{"divide", Request{Operation: Divide, A: ptr(10), B: ptr(4)}, 2.5},
		{"power", Request{Operation: Power, A: ptr(2), B: ptr(8)}, 256},
		{"power zero exponent", Request{Operation: Power, A: ptr(5), B: ptr(0)}, 1},
		{"sqrt", Request{Operation: Sqrt, A: ptr(16)}, 4},
		{"sqrt zero", Request{Operation: Sqrt, A: ptr(0)}, 0},
		{"sqrt irrational", Request{Operation: Sqrt, A: ptr(2)}, math.Sqrt(2)},
		{"percentage", Request{Operation: Percentage, A: ptr(25), B: ptr(200)}, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Calculate(tt.req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.Result != tt.want {
				t.Errorf("got %v, want %v", r.Result, tt.want)
			}
		})
	}
}

func TestCalculateErrors(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		want error
	}{
		{"divide by zero", Request{Operation: Divide, A: ptr(10), B: ptr(0)}, ErrDivisionByZero},
		{"sqrt negative", Request{Operation: Sqrt, A: ptr(-4)}, ErrNegativeSqrt},
		{"unknown operation", Request{Operation: "modulo", A: ptr(10), B: ptr(3)}, ErrUnknownOperation},
		{"missing operand a", Request{Operation: Add, B: ptr(3)}, ErrMissingOperand},
		{"missing operand b", Request{Operation: Add, A: ptr(3)}, ErrMissingOperand},
		{"power producing NaN", Request{Operation: Power, A: ptr(-1), B: ptr(0.5)}, ErrNonFiniteResult},
		{"power overflowing to Inf", Request{Operation: Power, A: ptr(10), B: ptr(1000)}, ErrNonFiniteResult},
		{"multiply overflowing to Inf", Request{Operation: Multiply, A: ptr(1e308), B: ptr(10)}, ErrNonFiniteResult},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Calculate(tt.req)
			if !errors.Is(err, tt.want) {
				t.Errorf("got %v, want %v", err, tt.want)
			}
		})
	}
}

func TestExpressionFormatting(t *testing.T) {
	tests := []struct {
		req  Request
		want string
	}{
		{Request{Operation: Add, A: ptr(5), B: ptr(3)}, "5 + 3"},
		{Request{Operation: Subtract, A: ptr(10), B: ptr(2.5)}, "10 - 2.5"},
		{Request{Operation: Multiply, A: ptr(3), B: ptr(7)}, "3 * 7"},
		{Request{Operation: Divide, A: ptr(10), B: ptr(3)}, "10 / 3"},
		{Request{Operation: Divide, A: ptr(-2.5), B: ptr(0.5)}, "-2.5 / 0.5"},
		{Request{Operation: Power, A: ptr(2), B: ptr(8)}, "2 ^ 8"},
		{Request{Operation: Sqrt, A: ptr(16)}, "√16"},
		{Request{Operation: Percentage, A: ptr(25), B: ptr(200)}, "25% of 200"},
		{Request{Operation: Add, A: ptr(1e20), B: ptr(1)}, "100000000000000000000 + 1"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			r, err := Calculate(tt.req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.Expression != tt.want {
				t.Errorf("got %q, want %q", r.Expression, tt.want)
			}
		})
	}
}
