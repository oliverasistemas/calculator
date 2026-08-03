package calculator

import (
	"fmt"
	"math"
	"strconv"
)

func formatNum(n float64) string {
	if n == math.Trunc(n) && !math.IsInf(n, 0) {
		return strconv.FormatInt(int64(n), 10)
	}
	return fmt.Sprintf("%g", n)
}

func formatBinary(a float64, op string, b float64) string {
	return formatNum(a) + " " + op + " " + formatNum(b)
}
