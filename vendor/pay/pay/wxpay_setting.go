package pay

import (
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"
)

var wxenvironments = map[string]string{
	"production":    "pay/wxpay.json",
	"preproduction": "pay/wxpay.json",
	"tests":         "pay/wxpay.json",
}

type WXpaySettings struct {

	//微信开发平台应用id
	AppId string `json:"app_id"`

	//应用对应的凭证
	AppKey string `json:"app_key"`

	//财付通商户号
	Partner  string `json:"partner"`

	//下单接口
	PlaceOrderUrl string `json:"place_order_url"`

	//查询订单接口
	QueryOrderUrl string `json:"query_order_url"`


}

var wxsettings *WXpaySettings
var wxenv = "preproduction"

func InitWXPay() {
	wxenv = os.Getenv("GO_ENV")

	pwd, _ := os.Getwd()
	fmt.Println(pwd)
	if wxenv == "" {
		fmt.Println("Warning: Setting preproduction environment due to lack of GO_ENV value")
		wxenv = "preproduction"
	}
	LoadWXpaySettingsByEnv(wxenv)
}

func LoadWXpaySettingsByEnv(env string) {
	content, err := ioutil.ReadFile(wxenvironments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)
	}
	wxsettings = &WXpaySettings{}
	jsonErr := json.Unmarshal(content, &wxsettings)
	if jsonErr != nil {
		fmt.Println("Error while parsing config file", jsonErr)
	}
}

func GetWXpaySetting() *WXpaySettings {
	if wxsettings == nil {
		fmt.Println("------")
		InitWXPay()
	}
	return wxsettings
}