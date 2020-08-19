package tools

import (
	"math"
	"strconv"
)

func Float64ToString(e float64) string {
	return strconv.FormatFloat(e, 'E', -1, 64)
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n // TODO +0.5 是为了四舍五入，如果不希望这样去掉这个
}
