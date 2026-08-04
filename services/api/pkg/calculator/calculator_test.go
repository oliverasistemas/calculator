package calculator

import (
	"errors"
	"testing"
)

func num(s string) *Number {
	n, err := NewNumber(s)
	if err != nil {
		panic(err)
	}
	return n
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		want float64
	}{
		{"add positive", Request{Operation: Add, A: num("2"), B: num("3")}, 5},
		{"add negative", Request{Operation: Add, A: num("-2"), B: num("-3")}, -5},
		{"add mixed", Request{Operation: Add, A: num("-2"), B: num("3")}, 1},
		{"add zero", Request{Operation: Add, A: num("0"), B: num("0")}, 0},
		{"add decimals", Request{Operation: Add, A: num("0.1"), B: num("0.2")}, 0.3},
		{"add cancellation", Request{Operation: Add, A: num("1.0000001"), B: num("-1")}, 1e-7},
		{"subtract", Request{Operation: Subtract, A: num("10"), B: num("4")}, 6},
		{"subtract cancellation", Request{Operation: Subtract, A: num("1.0000001"), B: num("1")}, 1e-7},
		{"subtract deep cancellation", Request{Operation: Subtract, A: num("1.0000000001"), B: num("1")}, 1e-10},
		{"subtract decimals", Request{Operation: Subtract, A: num("0.3"), B: num("0.1")}, 0.2},
		{"multiply", Request{Operation: Multiply, A: num("3"), B: num("7")}, 21},
		{"multiply large", Request{Operation: Multiply, A: num("1e15"), B: num("1e15")}, 1e30},
		{"multiply exact integer beyond 12 digits", Request{Operation: Multiply, A: num("9999999"), B: num("9999999")}, 99999980000001},
		{"multiply beyond 12 significant digits", Request{Operation: Multiply, A: num("1.000000000001"), B: num("1.000000000001")}, 1.000000000002},
		{"multiply small decimals", Request{Operation: Multiply, A: num("0.1"), B: num("0.2")}, 0.02},
		{"divide", Request{Operation: Divide, A: num("10"), B: num("4")}, 2.5},
		{"power", Request{Operation: Power, A: num("2"), B: num("8")}, 256},
		{"power zero exponent", Request{Operation: Power, A: num("5"), B: num("0")}, 1},
		{"power zero to the zero", Request{Operation: Power, A: num("0"), B: num("0")}, 1},
		{"power negative exponent", Request{Operation: Power, A: num("2"), B: num("-2")}, 0.25},
		{"sqrt", Request{Operation: Sqrt, A: num("16")}, 4},
		{"sqrt zero", Request{Operation: Sqrt, A: num("0")}, 0},
		{"sqrt irrational", Request{Operation: Sqrt, A: num("2")}, 1.41421356237},
		{"percentage", Request{Operation: Percentage, A: num("25"), B: num("200")}, 50},
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

// TestResultText covers the digits float64 cannot: exact integer powers
// beyond 2^53, operands beyond 2^53, and honest rounding when an exact
// answer would not fit the digit budget.
func TestResultText(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		want string
	}{
		{"power beyond 2^53 is exact", Request{Operation: Power, A: num("3"), B: num("35")}, "50031545098999707"},
		{"power 2^100 is exact", Request{Operation: Power, A: num("2"), B: num("100")}, "1267650600228229401496703205376"},
		{"operand beyond 2^53 survives", Request{Operation: Add, A: num("9007199254740993"), B: num("0")}, "9007199254740993"},
		{"add decimals", Request{Operation: Add, A: num("0.1"), B: num("0.2")}, "0.3"},
		{"divide terminating", Request{Operation: Divide, A: num("10"), B: num("8")}, "1.25"},
		{"divide non-terminating rounds to 12 digits", Request{Operation: Divide, A: num("1"), B: num("3")}, "0.333333333333"},
		{"power negative exponent", Request{Operation: Power, A: num("2"), B: num("-2")}, "0.25"},
		{"power non-integer exponent falls back to float", Request{Operation: Power, A: num("2"), B: num("0.5")}, "1.41421356237"},
		{"oversized exact result rounds honestly", Request{Operation: Power, A: num("10"), B: num("60")}, "1e+60"},
		{"sqrt irrational", Request{Operation: Sqrt, A: num("2")}, "1.41421356237"},
		{"percentage", Request{Operation: Percentage, A: num("25"), B: num("200")}, "50"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Calculate(tt.req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.ResultText != tt.want {
				t.Errorf("got %q, want %q", r.ResultText, tt.want)
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
		{"divide by zero", Request{Operation: Divide, A: num("10"), B: num("0")}, ErrDivisionByZero},
		{"zero to a negative power", Request{Operation: Power, A: num("0"), B: num("-2")}, ErrDivisionByZero},
		{"sqrt negative", Request{Operation: Sqrt, A: num("-4")}, ErrNegativeSqrt},
		{"unknown operation", Request{Operation: "modulo", A: num("10"), B: num("3")}, ErrUnknownOperation},
		{"missing operand a", Request{Operation: Add, B: num("3")}, ErrMissingOperand},
		{"missing operand b", Request{Operation: Add, A: num("3")}, ErrMissingOperand},
		{"power producing NaN", Request{Operation: Power, A: num("-1"), B: num("0.5")}, ErrNonFiniteResult},
		{"power overflowing to Inf", Request{Operation: Power, A: num("10"), B: num("1000")}, ErrNonFiniteResult},
		{"multiply overflowing to Inf", Request{Operation: Multiply, A: num("1e308"), B: num("10")}, ErrNonFiniteResult},
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

func TestNewNumberRejectsGarbage(t *testing.T) {
	for _, s := range []string{"", "abc", "1.2.3", "1e", "0x10", "Inf", "NaN", "1/3", "--5"} {
		if _, err := NewNumber(s); err == nil {
			t.Errorf("NewNumber(%q) succeeded, want error", s)
		}
	}
}

func TestExpressionFormatting(t *testing.T) {
	tests := []struct {
		req  Request
		want string
	}{
		{Request{Operation: Add, A: num("5"), B: num("3")}, "5 + 3"},
		{Request{Operation: Subtract, A: num("10"), B: num("2.5")}, "10 - 2.5"},
		{Request{Operation: Multiply, A: num("3"), B: num("7")}, "3 * 7"},
		{Request{Operation: Divide, A: num("10"), B: num("3")}, "10 / 3"},
		{Request{Operation: Divide, A: num("-2.5"), B: num("0.5")}, "-2.5 / 0.5"},
		{Request{Operation: Power, A: num("2"), B: num("8")}, "2 ^ 8"},
		{Request{Operation: Sqrt, A: num("16")}, "√16"},
		{Request{Operation: Percentage, A: num("25"), B: num("200")}, "25% of 200"},
		{Request{Operation: Add, A: num("1e20"), B: num("1")}, "100000000000000000000 + 1"},
		{Request{Operation: Add, A: num("0.30000000000000004"), B: num("5")}, "0.30000000000000004 + 5"},
		{Request{Operation: Multiply, A: num("0.1"), B: num("0.2")}, "0.1 * 0.2"},
		{Request{Operation: Multiply, A: num("1.000000000001"), B: num("1.000000000001")}, "1.000000000001 * 1.000000000001"},
		{Request{Operation: Add, A: num("5."), B: num(".5")}, "5 + 0.5"},
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
