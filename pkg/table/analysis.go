package table

import (
	"fmt"
	"hash/crc32"
	"strconv"

	"gorm.io/gorm"
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

// e.g. DB.Scopes(DynamicTable(Crc32Hash, "test", "小圈圈")).Find(&tests)
// DynamicTable 设置动态表名scope params: f 分表计算函数 baseTable 基础表名 fieldValue 参与分表字段
func DynamicTable(f func(string) string, baseTable, fieldValue string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(fmt.Sprintf("%s_%s", baseTable, f(fieldValue)))
	}
}

func CreateSubTable(f func(string) string) {

}
