package calculator

import (
	"math"
	"math/big"
	"strconv"
)

// maxExactDigits caps how many digits of an exact decimal result are served
// verbatim; longer results are rounded to 12 significant digits instead.
const maxExactDigits = 50

// exactDecimal renders r as an exact decimal string, or reports false when
// the expansion does not terminate within maxExactDigits.
func exactDecimal(r *big.Rat) (string, bool) {
	places, ok := terminatingPlaces(r)
	if !ok {
		return "", false
	}
	s := r.FloatString(places)
	digits := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			digits++
		}
	}
	if digits > maxExactDigits {
		return "", false
	}
	return s, true
}

// terminatingPlaces reports how many decimal places r needs to print
// exactly: its decimal expansion terminates iff the reduced denominator is
// 2^i * 5^j, at max(i, j) places. big.Rat keeps denominators reduced, so
// that last digit is never a trailing zero.
func terminatingPlaces(r *big.Rat) (int, bool) {
	d := new(big.Int).Set(r.Denom())
	twos := 0
	for d.Bit(0) == 0 {
		if twos++; twos > maxExactDigits {
			return 0, false
		}
		d.Rsh(d, 1)
	}
	five := big.NewInt(5)
	q, m := new(big.Int), new(big.Int)
	fives := 0
	for {
		q.QuoRem(d, five, m)
		if m.Sign() != 0 {
			break
		}
		if fives++; fives > maxExactDigits {
			return 0, false
		}
		d.Set(q)
	}
	if d.Cmp(big.NewInt(1)) != 0 {
		return 0, false
	}
	return max(twos, fives), true
}

// formatRat renders r exactly when its decimal expansion fits, otherwise
// rounded to 12 significant digits via the nearest float64.
func formatRat(r *big.Rat) string {
	if s, ok := exactDecimal(r); ok {
		return s
	}
	f, _ := r.Float64()
	return formatFloat(roundTo12(f))
}

// formatFloat renders a float result as display text. Integer-valued floats
// beyond 2^53 are deliberately not printed digit-for-digit — those digits
// would be rounding artifacts.
func formatFloat(n float64) string {
	if n == math.Trunc(n) && math.Abs(n) <= 1<<53 {
		return strconv.FormatFloat(n, 'f', -1, 64)
	}
	return strconv.FormatFloat(n, 'g', -1, 64)
}

// roundTo12 rounds n to 12 significant digits, absorbing binary float
// artifacts (0.1+0.2 → 0.30000000000000004) while leaving typical
// calculator results intact. Integers are exact in float64 only up to 2^53;
// above that the trailing digits are rounding artifacts, so large integers
// are rounded too rather than served with false precision.
func roundTo12(n float64) float64 {
	if n == math.Trunc(n) && math.Abs(n) <= 1<<53 {
		return n
	}
	rounded, err := strconv.ParseFloat(strconv.FormatFloat(n, 'g', 12, 64), 64)
	if err != nil {
		return n
	}
	return rounded
}

func formatBinary(a *Number, op string, b *Number) string {
	return a.String() + " " + op + " " + b.String()
}
