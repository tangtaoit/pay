package main

import (
	"net/http"
	. "pay/route"
)

func main() {







	router := NewRouter([]Route{

		Route{  //应用申请
			"SubmitApp",
			"POST",
			"/pay/app",
			SubmitApp,
		},
		Route{ //获取支付token
			"GetPayToken",
			"GET",
			"/pay/info/{app_id}/{app_key}",
			GetPayToken,
		},
		Route{ //绑定支付信息
			"BindPayInfo",
			"POST",
			"/pay/info/{app_id}/{app_key}",
			BindPayInfo,
		},
		Route{  //充值
			"MakePrePayOrder",
			"POST",
			"/pay/makeorder",
			MakePrePayOrder,
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
	})

	http.ListenAndServe(":8080", router)
}
