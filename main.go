package main

import (
	"net/http"
	. "pay/route"
	"github.com/tangtaoit/queue"
	"github.com/gorilla/mux"
	"github.com/tangtaoit/util"
	"log"
	"io/ioutil"
	"pay/config"
)

func GetRouters()  *mux.Router{

	return  NewRouter([]Route{

		Route{  //应用申请
			"SubmitApp",
			"POST",
			"/pay/app",
			SubmitApp,
		},
		Route{ //获取支付token
			"GetPayToken",
			"GET",
			"/pay/token",
			GetPayToken,
		},
		Route{ //绑定支付信息
			"BindPayInfo",
			"POST",
			"/pay/info/{app_id}/{app_key}",
			BindPayInfo,
		},
		Route{  //预备支付
			"MakePrePayOrder",
			"POST",
			"/pay/makeprepay",
			MakePrePay,
		},
		Route{  //支付
			"MakePay",
			"POST",
			"/pay/makepay",
			MakePay,
		},
		Route{  //支付宝回调
			"AlipayCallback",
			"POST",
			"/pay/alipay_callback",
			AlipayCallback,
		},
		Route{  //微信支付回调
			"AlipayCallback",
			"POST",
			"/pay/wxpay_callback",
			GetServerCallback(),
		},
		Route{  //test
			"SetAccountPassword",
			"POST",
			"/pay/password",
			SetAccountPassword,
		},

	})
}

func TestTimeout(w http.ResponseWriter, r *http.Request)  {

	byts,err := ioutil.ReadAll(r.Body)
	util.CheckErr(err)
	log.Println(string(byts))
	util.ResponseSuccess(w)
}

func main() {

	log.Println("amqpUrl="+config.GetSetting().AmqpUrl)

	queue.SetupAMQP(config.GetSetting().AmqpUrl)

	
	http.ListenAndServe(":8080", GetRouters())
}
