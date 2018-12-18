package alipay

import (
	"crypto"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
	"zcm_tools/crypt"
	"zcm_tools/http"
)

//芝麻信用
var (
	zhimaPublicKey  []byte //阿里公钥
	zhimaPrivateKey []byte //商户私钥
	zhimaApp_id     string //支付宝分配给开发者的应用ID
	zhimaUrl        string
)

func InitZhiMa(app_id, request_url string, privateKey, publicKey []byte) {
	zhimaApp_id = app_id
	zhimaPrivateKey = privateKey
	zhimaPublicKey = publicKey
	zhimaUrl = request_url
}

func RequestZhiMaCredit(method, biz_content string) ([]byte, error) {
	var param = url.Values{}
	param.Add("biz_content", biz_content)
	return doRequest(zhimaApp_id, method, biz_content, param, zhimaPrivateKey)
}

func doRequest(app_id, method, biz_content string, param url.Values, privateKey []byte) ([]byte, error) {
	//添加公共参数
	param.Add("app_id", app_id)
	param.Add("method", method)
	param.Add("format", "JSON")
	param.Add("charset", "utf-8")
	param.Add("sign_type", "RSA2")
	param.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	param.Add("version", "1.0")
	fmt.Println(url.QueryUnescape(param.Encode()))
	s, err := url.QueryUnescape(param.Encode())

	signData := []byte(s) //GetSignData(param)
	sign, err := crypt.SignPKCS1v15(signData, privateKey, crypto.SHA256)
	if err != nil {
		return nil, err
	}
	param.Add("sign", sign)
	if strings.HasPrefix(zhimaUrl, "https://") {
		return http.HttpsPost(zhimaUrl, param.Encode())
	} else {
		return http.HttpPost(zhimaUrl, param.Encode())
	}
}

//获取需要签名的数据
func GetSignData(param url.Values) (s string) {
	var keys = make([]string, 0, 0)
	for key, _ := range param {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value) //url.QueryEscape(value)
		}
	}
	return strings.Join(pList, "&")
}
