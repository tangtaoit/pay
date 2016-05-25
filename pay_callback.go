package main

import (
	"net/http"
	"log"
	"encoding/xml"
	"pay_BAK/comm"
	"io"
	"io/ioutil"
)

type WXNotifyResult  struct{
	XMLName  xml.Name `xml:"xml"`
	//应用ID
	AppId string `xml:"app_id"`
	//商户号
	MchId string `xml:"mch_id"`
	//设备号(微信支付分配的终端设备号，)
	DeviceInfo string `xml:"device_info"`
	//随机字符串，不长于32位
	NonceStr string `xml:"nonce_str"`
	//签名，详见签名算法
	Sign string `xml:"sign"`
	//业务结果 SUCCESS/FAIL
	ResultCode string `xml:"result_code"`
	//错误返回的信息描述
	ErrCode string `xml:"err_code"`
	//错误代码描述
	ErrCodeDes string `xml:"err_code_des"`
	//用户在商户appid下的唯一标识
	OpenId string `xml:"openid"`
	//用户是否关注公众账号，Y-关注，N-未关注，仅在公众账号类型支付有效
	IsSubscribe string `xml:"is_subscribe"`
	//交易类型 APP
	TradeType string `xml:"trade_type"`
	//付款银行 银行类型，采用字符串类型的银行标识，银行类型见银行列表
	BankType string `xml:"bank_type"`
	//订单总金额，单位为分
	TotalFee int64 `xml:"total_fee"`
	//货币类型，符合ISO4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	FeeType string `xml:"fee_type"`
	//现金支付金额订单现金支付金额，详见支付金额
	CashFee string `xml:"cash_fee"`
	//现金支付货币类型 货币类型，符合ISO4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashFeeType string `xml:"cash_fee_type"`
	//微信支付订单号
	TransactionId string `xml:"transaction_id"`
	//商户系统的订单号，与请求一致
	OutTradeNo string `xml:"out_trade_no"`
	//商家数据包，原样返回
	Attach  string `xml:"attach"`
	//支付完成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	TimeEnd string `xml:"time_end"`

}

//支付宝回调
func AlipayCallback(w http.ResponseWriter, r *http.Request)  {

	log.Println("支付宝回调了...");
}

func WXpayCallback(w http.ResponseWriter, r *http.Request)  {

	log.Println("微信回调了...");

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	comm.CheckErr(err)
	var wxNotifyResult WXNotifyResult;
	comm.CheckErr(xml.Unmarshal(body,&wxNotifyResult))

	log.Println("wxNotifyResult==",wxNotifyResult.ResultCode)

}