package pay

import (
	"fmt"
	"crypto/tls"
	"net/http"
	"bytes"
	"io/ioutil"
)



type PaymentRequest struct {
	AppId     string `json:"app_id"`
	PartnerId string `json:"partner_id"`
	PrepayId  string `json:"prepay_id"`
	Package   string `json:"package"`
	NonceStr  string `json:"noncestr"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

//微信预下单model
type WXPrepay struct  {
	//订单号
	OrderId string
	//订单金额
	Amount string
	//描述
	Desc string
	//客户端IP
	ClientIp string

	//通知地址
	NotifyUrl string
}

type WXPrepayResult struct  {

	Appid string `json:"appid"`

	Partnerid string `json:"partnerid"`

	Package  string `json:"package"`

	Noncestr string `json:"noncestr"`

	Timestamp int64 `json:"timestamp"`

	Prepayid string `json:"prepayid"`

	Sign string `json:"sign"`

}

func NewWXPrepay() *WXPrepay {

	return &WXPrepay{}
}

func NewWXPrepayResult() *WXPrepayResult {

	return &WXPrepayResult{}
}

//预支付
func ( self *WXPrepay) Prepay() (PaymentRequest,error) {

	var paymentRequest PaymentRequest;

	order := self.newOrderRequest()
	odrInXml := self.signedOrderRequestXmlString(order)
	fmt.Println("print:",odrInXml);
	resp, err := doHttpPost(GetWXpaySetting().PlaceOrderUrl, []byte(odrInXml))
	if err != nil {
		return paymentRequest, err
	}
	placeOrderResult, err := ParsePlaceOrderResult(resp)
	if err != nil {
		return paymentRequest, err
	}

	//Verify the sign of response
	resultInMap := placeOrderResult.ToMap()
	wantSign := Sign(resultInMap, GetWXpaySetting().AppKey)
	gotSign := resultInMap["sign"]
	if wantSign != gotSign {
		return paymentRequest, fmt.Errorf("sign not match, want:%s, got:%s", wantSign, gotSign)
	}

	if placeOrderResult.ReturnCode != "SUCCESS" {
		return paymentRequest, fmt.Errorf("%s[%s]", placeOrderResult.ReturnMsg,placeOrderResult.ReturnCode)
	}

	if placeOrderResult.ResultCode != "SUCCESS" {

		return paymentRequest, fmt.Errorf("%s[%s]", placeOrderResult.ErrCodeDesc,placeOrderResult.ErrCode)
	}



	return self.NewPaymentRequest(placeOrderResult.PrepayId), nil

}

func (self *WXPrepay) NewPaymentRequest(prepayId string) PaymentRequest {

	conf :=GetWXpaySetting()

	nonceStr :=NewNonceString()
	timestampStr :=NewTimestampString()

	param := make(map[string]string)
	param["appid"] = conf.AppId
	param["partnerid"] = conf.Partner
	param["prepayid"] = prepayId
	param["package"] = "Sign=WXPay"
	param["noncestr"] = nonceStr
	param["timestamp"] = timestampStr

	sign := Sign(param, conf.AppKey)

	payRequest := PaymentRequest{
		AppId:     conf.AppId,
		PartnerId: conf.Partner,
		PrepayId:  prepayId,
		Package:   "Sign=WXPay",
		NonceStr:  nonceStr,
		Timestamp: timestampStr,
		Sign:      sign,
	}

	return payRequest
}

func (self *WXPrepay) newOrderRequest() map[string]string {

	conf :=GetWXpaySetting()
	param := make(map[string]string)
	param["appid"] = conf.AppId
	param["body"] = self.Desc
	param["mch_id"] = conf.Partner
	param["nonce_str"] = NewNonceString()
	param["notify_url"] = self.NotifyUrl
	param["out_trade_no"] = self.OrderId
	param["spbill_create_ip"] = self.ClientIp
	param["total_fee"] = self.Amount
	param["trade_type"] = "APP"

	return param
}

func (self *WXPrepay) signedOrderRequestXmlString(orderMap map[string]string) string {

	sign := Sign(orderMap, GetWXpaySetting().AppKey)
	 fmt.Println(sign)

	orderMap["sign"] = sign

	return ToXmlString(orderMap)
}

// doRequest post the order in xml format with a sign
func doHttpPost(targetUrl string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return []byte(""), err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return respData, nil
}