package tools

import "strconv"

func IntToString(e int) string {
	return strconv.Itoa(e)
}

func IntArrayToString(e []int) string {
	str := ""
	for i := 0; i < len(e); i++ {
		str+=strconv.Itoa(e[i])
	}
	return str
}
