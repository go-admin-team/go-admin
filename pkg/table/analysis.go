package table

import (
	"hash/crc32"
	"strconv"
)

// Crc32Hash 用于32张分表
func Crc32Hash(src string) string {
	sum := crc32.ChecksumIEEE([]byte(src)) % 32
	return strconv.Itoa(int(sum))
}

// Crc16Hash 用于16张分表
func Crc16Hash(src string) string {
	sum := crc32.ChecksumIEEE([]byte(src)) % 16
	return strconv.Itoa(int(sum))
}

// Crc8Hash 用于8张分表
func Crc8Hash(src string) string {
	sum := crc32.ChecksumIEEE([]byte(src)) % 8
	return strconv.Itoa(int(sum))
}
