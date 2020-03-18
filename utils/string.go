package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func StringToInt64(e string) (int64, error) {
	return strconv.ParseInt(e, 10, 64)
}

func IntToString(e int) string {
	return strconv.Itoa(e)
}

func Float64ToString(e float64) string {
	return strconv.FormatFloat(e, 'E', -1, 64)
}

func Int64ToString(e int64) string {
	return strconv.FormatInt(e, 10)
}

//获取URL中批量id并解析
func IdsStrToIdsInt64Group(key string, c *gin.Context) []int64 {
	IDS := make([]int64, 0)
	ids := strings.Split(c.Param(key), ",")
	for i := 0; i < len(ids); i++ {
		ID, _ := strconv.ParseInt(ids[i], 10, 64)
		IDS = append(IDS, ID)
	}
	return IDS
}

func GetCurrntTime() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func GetLocation(ip string) string {
	if ip == "127.0.0.1" || ip == "localhost" {
		return "内部IP"
	}
	resp, err := http.Get("https://restapi.amap.com/v3/ip?ip=" + ip + "&key=3fabc36c20379fbb9300c79b19d5d05e")
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(s))

	m := make(map[string]string)

	err = json.Unmarshal(s, &m)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
	}
	if m["province"] == "" {
		return "未知位置"
	}
	return m["province"] + "-" + m["city"]
}

func StructToJsonStr(e interface{}) (string, error) {
	if b, err := json.Marshal(e); err == nil {
		return string(b), err
	} else {
		return "", err
	}
}

func GetBodyString(c *gin.Context) (string, error) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return string(body), nil
	} else {
		return "", err
	}
}

func JsonStrToMap(e string) (map[string]interface{}, error) {
	var dict map[string]interface{}
	if err := json.Unmarshal([]byte(e), &dict); err == nil {
		return dict, err
	} else {
		return nil, err
	}
}
