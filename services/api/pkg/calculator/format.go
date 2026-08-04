package calculator

import (
	"math"
	"strconv"
)

func formatNum(n float64) string {
	if n == math.Trunc(n) && !math.IsInf(n, 0) {
		return strconv.FormatFloat(n, 'f', -1, 64)
	}
	return strconv.FormatFloat(roundTo12(n), 'g', -1, 64)
}

// roundTo12 rounds n to 12 significant digits, absorbing binary float
// artifacts (0.1+0.2 → 0.30000000000000004) while leaving typical
// calculator results intact. Integers are exact in float64 up to 2^53 and
// may need more than 12 digits, so they pass through untouched.
func roundTo12(n float64) float64 {
	if n == math.Trunc(n) && !math.IsInf(n, 0) {
		return n
	}
	rounded, err := strconv.ParseFloat(strconv.FormatFloat(n, 'g', 12, 64), 64)
	if err != nil {
		return n
	}
	return rounded
}

func formatBinary(a float64, op string, b float64) string {
	return formatNum(a) + " " + op + " " + formatNum(b)
}
