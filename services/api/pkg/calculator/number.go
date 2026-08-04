package calculator

import (
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

// numberPattern accepts plain decimal numbers: JSON number syntax plus the
// forms a calculator display can produce ("5.", ".5", a leading "+").
var numberPattern = regexp.MustCompile(`^[+-]?([0-9]+\.?[0-9]*|\.[0-9]+)([eE][+-]?[0-9]+)?$`)

// maxNumberLen bounds operand literals so hostile requests cannot feed
// arbitrarily large rationals into the exact-arithmetic paths.
const maxNumberLen = 512

// Number is a request operand. It keeps the literal decimal text from the
// JSON payload (a number or a string) alongside its exact rational value,
// because decoding into float64 silently alters integers beyond 2^53 —
// {"a":9007199254740993} would become 9007199254740992 before any math ran.
type Number struct {
	text string   // literal as received, echoed in expressions
	rat  *big.Rat // exact decimal value
	f    float64  // nearest float64, for irrational fallbacks
}

func NewNumber(s string) (*Number, error) {
	if len(s) > maxNumberLen || !numberPattern.MatchString(s) {
		return nil, fmt.Errorf("invalid number %q", s)
	}
	norm := strings.TrimPrefix(s, "+")
	// big.Rat.SetString needs a digit on both sides of the decimal point.
	norm = strings.Replace(norm, ".e", "e", 1)
	norm = strings.Replace(norm, ".E", "E", 1)
	norm = strings.TrimSuffix(norm, ".")
	norm = strings.Replace(norm, "-.", "-0.", 1)
	if strings.HasPrefix(norm, ".") {
		norm = "0" + norm
	}
	r, ok := new(big.Rat).SetString(norm)
	if !ok {
		return nil, fmt.Errorf("invalid number %q", s)
	}
	f, _ := r.Float64()
	return &Number{text: s, rat: r, f: f}, nil
}

// UnmarshalJSON accepts both a JSON number and a JSON string. Numbers are
// captured from their literal digits, never through float64.
func (n *Number) UnmarshalJSON(data []byte) error {
	s := string(data)
	if strings.HasPrefix(s, `"`) {
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
	}
	num, err := NewNumber(s)
	if err != nil {
		return err
	}
	*n = *num
	return nil
}

// String renders the operand for expression strings: the canonical decimal
// form of its exact value, or the literal input when the expansion is too
// long to spell out.
func (n *Number) String() string {
	if s, ok := exactDecimal(n.rat); ok {
		return s
	}
	return n.text
}
