package calculator

import (
	"errors"
	"math"
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
	res.Result = roundTo12(res.Result)
	return res, nil
}

func binary(op Operation, a, b float64) (*Result, error) {
	var result float64
	var expression string

	switch op {
	case Add:
		result = a + b
		expression = formatBinary(a, "+", b)
	case Subtract:
		result = a - b
		expression = formatBinary(a, "-", b)
	case Multiply:
		result = a * b
		expression = formatBinary(a, "*", b)
	case Divide:
		if b == 0 {
			return nil, ErrDivisionByZero
		}
		result = a / b
		expression = formatBinary(a, "/", b)
	case Power:
		result = math.Pow(a, b)
		expression = formatBinary(a, "^", b)
	case Percentage:
		result = (a / 100) * b
		expression = formatNum(a) + "% of " + formatNum(b)
	}

	return &Result{Result: result, Expression: expression}, nil
}

func sqrt(a float64) (*Result, error) {
	if a < 0 {
		return nil, ErrNegativeSqrt
	}
	return &Result{
		Result:     math.Sqrt(a),
		Expression: "√" + formatNum(a),
	}, nil
}
