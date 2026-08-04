package calculator

import (
	"math"
	"strconv"
)

func formatNum(n float64) string {
	if n == math.Trunc(n) && !math.IsInf(n, 0) {
		return strconv.FormatFloat(n, 'f', -1, 64)
	}
	// Expressions are display-only: round to 12 significant digits so
	// float artifacts (0.1+0.2 → 0.30000000000000004) don't surface when a
	// prior result is chained into the next operation.
	rounded, err := strconv.ParseFloat(strconv.FormatFloat(n, 'g', 12, 64), 64)
	if err != nil {
		return strconv.FormatFloat(n, 'g', -1, 64)
	}
	return strconv.FormatFloat(rounded, 'g', -1, 64)
}

func formatBinary(a float64, op string, b float64) string {
	return formatNum(a) + " " + op + " " + formatNum(b)
}
