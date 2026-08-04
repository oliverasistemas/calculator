package calculator

import (
	"errors"
	"math"
	"math/big"
	"strconv"
)

var (
	ErrDivisionByZero   = errors.New("division by zero")
	ErrNegativeSqrt     = errors.New("square root of negative number")
	ErrUnknownOperation = errors.New("unknown operation")
	ErrMissingOperand   = errors.New("missing required operand")
	ErrNonFiniteResult  = errors.New("result is not a finite number")
)

type Operation string

const (
	Add        Operation = "add"
	Subtract   Operation = "subtract"
	Multiply   Operation = "multiply"
	Divide     Operation = "divide"
	Power      Operation = "power"
	Sqrt       Operation = "sqrt"
	Percentage Operation = "percentage"
)

type Request struct {
	Operation Operation `json:"operation"`
	A         *Number   `json:"a"`
	B         *Number   `json:"b,omitempty"`
}

type Result struct {
	// ResultText carries the authoritative digits: a float64 cannot hold
	// integers beyond 2^53 (3^35 = 50031545098999707 has no float64
	// representation), so Result is only the nearest float64 to ResultText,
	// kept for compatibility. Clients must display and chain from
	// ResultText.
	Result     float64 `json:"result"`
	ResultText string  `json:"resultText"`
	Expression string  `json:"expression"`
}

func Calculate(req Request) (*Result, error) {
	if req.A == nil {
		return nil, ErrMissingOperand
	}
	switch req.Operation {
	case Sqrt:
		return sqrt(req.A)
	case Add, Subtract, Multiply, Divide, Power, Percentage:
		if req.B == nil {
			return nil, ErrMissingOperand
		}
		return binary(req.Operation, req.A, req.B)
	default:
		return nil, ErrUnknownOperation
	}
}

// Add, subtract, multiply, divide and percentage are closed over the
// rationals, so they are computed exactly; float64 only enters when the
// result is formatted.
func binary(op Operation, a, b *Number) (*Result, error) {
	r := new(big.Rat)
	switch op {
	case Add:
		r.Add(a.rat, b.rat)
		return ratResult(r, formatBinary(a, "+", b))
	case Subtract:
		r.Sub(a.rat, b.rat)
		return ratResult(r, formatBinary(a, "-", b))
	case Multiply:
		r.Mul(a.rat, b.rat)
		return ratResult(r, formatBinary(a, "*", b))
	case Divide:
		if b.rat.Sign() == 0 {
			return nil, ErrDivisionByZero
		}
		r.Quo(a.rat, b.rat)
		return ratResult(r, formatBinary(a, "/", b))
	case Power:
		return power(a, b)
	case Percentage:
		r.Quo(a.rat, big.NewRat(100, 1)).Mul(r, b.rat)
		return ratResult(r, a.String()+"% of "+b.String())
	}
	return nil, ErrUnknownOperation
}

// Bounds on the exact power path: exponent magnitude and base size (numerator
// plus denominator bits) both capped so hostile requests cannot force huge
// big.Int allocations. Anything outside falls back to math.Pow.
const (
	maxExactExponent = 4096
	maxExactBaseBits = 256
)

// power computes integer exponents exactly with big.Int. math.Pow alone can
// never be correct here: at the magnitude of 3^35 consecutive float64 values
// are 8 apart, so the true answer is simply not representable.
func power(a, b *Number) (*Result, error) {
	expr := formatBinary(a, "^", b)
	if e, ok := intExponent(b); ok && a.rat.Num().BitLen()+a.rat.Denom().BitLen() <= maxExactBaseBits {
		if a.rat.Sign() == 0 && e < 0 {
			return nil, ErrDivisionByZero
		}
		return ratResult(ratPow(a.rat, e), expr)
	}
	return floatResult(math.Pow(a.f, b.f), expr)
}

func intExponent(b *Number) (int64, bool) {
	if !b.rat.IsInt() || !b.rat.Num().IsInt64() {
		return 0, false
	}
	e := b.rat.Num().Int64()
	if e < -maxExactExponent || e > maxExactExponent {
		return 0, false
	}
	return e, true
}

func ratPow(a *big.Rat, e int64) *big.Rat {
	neg := e < 0
	if neg {
		e = -e
	}
	exp := big.NewInt(e)
	num := new(big.Int).Exp(a.Num(), exp, nil)
	den := new(big.Int).Exp(a.Denom(), exp, nil)
	if neg {
		num, den = den, num
	}
	return new(big.Rat).SetFrac(num, den)
}

func sqrt(a *Number) (*Result, error) {
	if a.rat.Sign() < 0 {
		return nil, ErrNegativeSqrt
	}
	return floatResult(math.Sqrt(a.f), "√"+a.String())
}

// ratResult formats the exact value first, then derives the compatibility
// float64 from that text so the two fields can never disagree.
func ratResult(r *big.Rat, expr string) (*Result, error) {
	text := formatRat(r)
	f, err := strconv.ParseFloat(text, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		// NaN and ±Inf are not representable in JSON; reject here rather
		// than letting response encoding fail after the status is written.
		return nil, ErrNonFiniteResult
	}
	return &Result{Result: f, ResultText: text, Expression: expr}, nil
}

func floatResult(f float64, expr string) (*Result, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nil, ErrNonFiniteResult
	}
	f = roundTo12(f)
	return &Result{Result: f, ResultText: formatFloat(f), Expression: expr}, nil
}
