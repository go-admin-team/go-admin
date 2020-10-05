package tools

import (
	"strings"

	"github.com/gin-gonic/gin"
)

//获取URL中批量id并解析
func IdsStrToIdsIntGroup(key string, c *gin.Context) []int {
	return IdsStrToIdsIntGroupStr(c.Param(key))
}

type Ids struct {
	Ids []int
}

func IdsStrToIdsIntGroupStr(keys string) []int {
	IDS := make([]int, 0)
	ids := strings.Split(keys, ",")
	for i := 0; i < len(ids); i++ {
		ID, _ := StringToInt(ids[i])
		IDS = append(IDS, ID)
	}
	return IDS
}
