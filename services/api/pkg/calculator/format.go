package calculator

import (
	"math"
	"strconv"
)

func formatNum(n float64) string {
	if n == math.Trunc(n) && !math.IsInf(n, 0) {
		return strconv.FormatFloat(n, 'f', -1, 64)
	}
	return strconv.FormatFloat(n, 'g', -1, 64)
}

func formatBinary(a float64, op string, b float64) string {
	return formatNum(a) + " " + op + " " + formatNum(b)
}
