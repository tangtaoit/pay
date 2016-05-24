package pay

import (
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"
)

var alienvironments = map[string]string{
	"production":    "config/alipay.json",
	"preproduction": "config/alipay.json",
	"tests":         "../../alipay.json",
}

type AlipaySettings struct {

	//合作伙伴ID
	Partner   string `json:"partner"`
	//商家私钥
	PrivateKey string `json:"private_key"`
	//支付宝公钥
	AliPublickey string `json:"ali_public_key"`

	//商家账号
	SellerID string `json:"seller_id"`

	//app_id
	AppId string `json:"app_id"`

	//商家支付宝私钥
	AlipayPrivateKeyPath string `json:"alipay_privateKey_Path"`


}

var alisettings *AlipaySettings
var alienv = "preproduction"

func InitAlipay() {
	alienv = os.Getenv("GO_ENV")

	pwd, _ := os.Getwd()
	fmt.Println(pwd)
	if alienv == "" {
		fmt.Println("Warning: Setting preproduction environment due to lack of GO_ENV value")
		alienv = "preproduction"
	}
	LoadAlipaySettingsByEnv(alienv)
}

func LoadAlipaySettingsByEnv(env string) {
	content, err := ioutil.ReadFile(alienvironments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)
	}
	alisettings = &AlipaySettings{}
	jsonErr := json.Unmarshal(content, &alisettings)
	if jsonErr != nil {
		fmt.Println("Error while parsing config file", jsonErr)
	}
}

func GetAlipaySetting() *AlipaySettings {
	if alisettings == nil {
		fmt.Println("------")
		InitAlipay()
	}
	return alisettings
}