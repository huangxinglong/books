package zcm_tools

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

//截取小数点后几位
func SubFloatToString(f float64, m int) string {
	n := strconv.FormatFloat(f, 'f', -1, 64)
	if n == "" {
		return ""
	}
	if m >= len(n) {
		return n
	}
	newn := strings.Split(n, ".")
	if m == 0 {
		return newn[0]
	}
	if len(newn) < 2 || m >= len(newn[1]) {
		return n
	}
	return newn[0] + "." + newn[1][:m]
}

//截取小数点后几位
func SubFloatToFloat(f float64, m int) float64 {
	newn := SubFloatToString(f, m)
	newf, _ := strconv.ParseFloat(newn, 64)
	return newf
}

//序列化
func ToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// 获取数字随机字符
func GetRandDigit(n int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(n)+"d", rnd.Intn(int(math.Pow10(n))))
}

func UrlEncode(s string) string {
	return url.QueryEscape(s)
}
