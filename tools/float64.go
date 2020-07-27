package tools

import "strconv"

func Float64ToString(e float64) string {
	return strconv.FormatFloat(e, 'E', -1, 64)
}
