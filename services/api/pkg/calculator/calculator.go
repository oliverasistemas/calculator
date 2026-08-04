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
	A         *float64  `json:"a"`
	B         *float64  `json:"b,omitempty"`
}

type Result struct {
	Result     float64 `json:"result"`
	Expression string  `json:"expression"`
}

func Calculate(req Request) (*Result, error) {
	if req.A == nil {
		return nil, ErrMissingOperand
	}

	a := *req.A

	var res *Result
	var err error
	switch req.Operation {
	case Sqrt:
		res, err = sqrt(a)
	case Add, Subtract, Multiply, Divide, Power, Percentage:
		if req.B == nil {
			return nil, ErrMissingOperand
		}
		res, err = binary(req.Operation, a, *req.B)
	default:
		return nil, ErrUnknownOperation
	}
	if err != nil {
		return nil, err
	}
	// NaN and ±Inf are not representable in JSON; reject them here rather
	// than letting response encoding fail after the status is written.
	if math.IsNaN(res.Result) || math.IsInf(res.Result, 0) {
		return nil, ErrNonFiniteResult
	}
	return res, nil
}

func binary(op Operation, a, b float64) (*Result, error) {
	var result float64
	var expression string

	switch op {
	case Add:
		result = decimalExact(op, a, b)
		expression = formatBinary(a, "+", b)
	case Subtract:
		result = decimalExact(op, a, b)
		expression = formatBinary(a, "-", b)
	case Multiply:
		result = decimalExact(op, a, b)
		expression = formatBinary(a, "*", b)
	case Divide:
		if b == 0 {
			return nil, ErrDivisionByZero
		}
		result = roundTo12(a / b)
		expression = formatBinary(a, "/", b)
	case Power:
		result = roundTo12(math.Pow(a, b))
		expression = formatBinary(a, "^", b)
	case Percentage:
		result = roundTo12((a / 100) * b)
		expression = formatNum(a) + "% of " + formatNum(b)
	}

	return &Result{Result: result, Expression: expression}, nil
}

// decimalExact computes a+b, a-b, or a*b in exact decimal arithmetic and
// returns the nearest float64. Doing these in float64 leaks representation
// error into the significant digits when operands cancel (1.0000001 - 1)
// or when the true result carries more digits than a fixed rounding keeps
// (1.000000000001²), so no after-the-fact rounding can recover the exact
// answer. The operands' decimal values (their shortest decimal form) are
// exact rationals, and add/subtract/multiply are closed over them.
func decimalExact(op Operation, a, b float64) float64 {
	ra, okA := new(big.Rat).SetString(strconv.FormatFloat(a, 'g', -1, 64))
	rb, okB := new(big.Rat).SetString(strconv.FormatFloat(b, 'g', -1, 64))
	if !okA || !okB {
		switch op {
		case Subtract:
			return a - b
		case Multiply:
			return a * b
		default:
			return a + b
		}
	}

	r := new(big.Rat)
	switch op {
	case Subtract:
		r.Sub(ra, rb)
	case Multiply:
		r.Mul(ra, rb)
	default:
		r.Add(ra, rb)
	}
	f, _ := r.Float64()
	return f
}

func sqrt(a float64) (*Result, error) {
	if a < 0 {
		return nil, ErrNegativeSqrt
	}
	return &Result{
		Result:     roundTo12(math.Sqrt(a)),
		Expression: "√" + formatNum(a),
	}, nil
}
