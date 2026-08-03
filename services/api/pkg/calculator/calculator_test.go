package calculator

import (
	"math"
	"testing"
)

func ptr(f float64) *float64 { return &f }

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive", 2, 3, 5},
		{"negative", -2, -3, -5},
		{"mixed", -2, 3, 1},
		{"zero", 0, 0, 0},
		{"decimals", 0.1, 0.2, 0.30000000000000004},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Calculate(Request{Operation: Add, A: ptr(tt.a), B: ptr(tt.b)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.Result != tt.want {
				t.Errorf("got %v, want %v", r.Result, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	r, err := Calculate(Request{Operation: Subtract, A: ptr(10), B: ptr(4)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 6 {
		t.Errorf("got %v, want 6", r.Result)
	}
}

func TestMultiply(t *testing.T) {
	r, err := Calculate(Request{Operation: Multiply, A: ptr(3), B: ptr(7)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 21 {
		t.Errorf("got %v, want 21", r.Result)
	}
}

func TestDivide(t *testing.T) {
	r, err := Calculate(Request{Operation: Divide, A: ptr(10), B: ptr(4)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 2.5 {
		t.Errorf("got %v, want 2.5", r.Result)
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Calculate(Request{Operation: Divide, A: ptr(10), B: ptr(0)})
	if err != ErrDivisionByZero {
		t.Errorf("got %v, want ErrDivisionByZero", err)
	}
}

func TestPower(t *testing.T) {
	r, err := Calculate(Request{Operation: Power, A: ptr(2), B: ptr(8)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 256 {
		t.Errorf("got %v, want 256", r.Result)
	}
}

func TestSqrt(t *testing.T) {
	r, err := Calculate(Request{Operation: Sqrt, A: ptr(16)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 4 {
		t.Errorf("got %v, want 4", r.Result)
	}
}

func TestSqrtNegative(t *testing.T) {
	_, err := Calculate(Request{Operation: Sqrt, A: ptr(-4)})
	if err != ErrNegativeSqrt {
		t.Errorf("got %v, want ErrNegativeSqrt", err)
	}
}

func TestPercentage(t *testing.T) {
	r, err := Calculate(Request{Operation: Percentage, A: ptr(25), B: ptr(200)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 50 {
		t.Errorf("got %v, want 50", r.Result)
	}
}

func TestUnknownOperation(t *testing.T) {
	_, err := Calculate(Request{Operation: "modulo", A: ptr(10), B: ptr(3)})
	if err != ErrUnknownOperation {
		t.Errorf("got %v, want ErrUnknownOperation", err)
	}
}

func TestMissingOperandA(t *testing.T) {
	_, err := Calculate(Request{Operation: Add, B: ptr(3)})
	if err != ErrMissingOperand {
		t.Errorf("got %v, want ErrMissingOperand", err)
	}
}

func TestMissingOperandB(t *testing.T) {
	_, err := Calculate(Request{Operation: Add, A: ptr(3)})
	if err != ErrMissingOperand {
		t.Errorf("got %v, want ErrMissingOperand", err)
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
		{Request{Operation: Power, A: ptr(2), B: ptr(8)}, "2 ^ 8"},
		{Request{Operation: Sqrt, A: ptr(16)}, "√16"},
		{Request{Operation: Percentage, A: ptr(25), B: ptr(200)}, "25% of 200"},
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

func TestSqrtZero(t *testing.T) {
	r, err := Calculate(Request{Operation: Sqrt, A: ptr(0)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 0 {
		t.Errorf("got %v, want 0", r.Result)
	}
}

func TestPowerZeroExponent(t *testing.T) {
	r, err := Calculate(Request{Operation: Power, A: ptr(5), B: ptr(0)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 1 {
		t.Errorf("got %v, want 1", r.Result)
	}
}

func TestLargeNumbers(t *testing.T) {
	r, err := Calculate(Request{Operation: Multiply, A: ptr(1e15), B: ptr(1e15)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Result != 1e30 {
		t.Errorf("got %v, want 1e30", r.Result)
	}
}

func TestSqrtIrrational(t *testing.T) {
	r, err := Calculate(Request{Operation: Sqrt, A: ptr(2)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(r.Result-math.Sqrt(2)) > 1e-15 {
		t.Errorf("got %v, want %v", r.Result, math.Sqrt(2))
	}
}
