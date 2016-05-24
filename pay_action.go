package main

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"pay/config"
	"github.com/gorilla/mux"
	"pay/comm"
	"pay/db"
	"time"
	"fmt"
	"pay/pay"
)

const (
	//支付宝
	AliPay = iota

	//微信支付
	WXPAY
)

type PayToken struct  {

	PayToken string `json:"pay_token"`

}


type RechargeModel struct {

	//充值金额(单位分)
	amount int64 `json:"amount"`

	//支付类型 1.支付宝  2.微信
	payType int `json:"pay_type"`
}



type AlipayOrderModel struct  {

	//合作商家ID
	partner string
	//商家账号
	sellerID string
	//订单ID（由商家自行制定）
	outTradeNO string
	//商品标题
	Subject string
	////商品描述
	Body string
	//商品价格
	TotalFee string
	//回调URL
	NotifyURL string
	service string
	paymentType string
	inputCharset string
	itBPay string
	showURL string
	AppId string
}





func (self *AlipayOrderModel) ToString()  string{



	return fmt.Sprintf("partner=%s&seller_id=%s&out_trade_no=%s&subject=%s&body=%s&total_fee=%s&notify_url=%s&service=%s&payment_type=%s&_input_charset=%s&it_b_pay=%s&show_url=%s&app_id=%s",
		self.partner,self.sellerID,self.outTradeNO,self.Subject,self.Body,self.TotalFee,self.NotifyURL,self.service,self.paymentType,self.inputCharset,self.itBPay,self.showURL,self.AppId)
}


func NewPayToken(token string)  *PayToken{

	return &PayToken{PayToken:token}
}


var PAY_TOKEN_PREFIX ="PAY_TOKEN:";


//获取支付TOKEN
func GetPayToken(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if isOk:=appIsOk(w,r);!isOk{
		return;
	}

	authorization :=r.FormValue("auth_token");
	token :=authPayToken(w,authorization);
	if token==nil {return}

	paytoken := comm.GenerUUId();
	 sub,_:=token.Claims["sub"].(string)
	key :=PAY_TOKEN_PREFIX+paytoken+sub;
	SetAndExpire(key,"1",config.GetSetting().TokenExpire)
	comm.WriteJson(w,NewPayToken(paytoken))
}

//绑定支付信息
func BindPayInfo(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if isOk:=appIsOk(w,r);!isOk{
		return;
	}

	auth_token:=r.Header.Get("auth_token")
	if auth_token==""{
		comm.ResponseError(w,http.StatusBadRequest,"认证信息不能为空!");
		return;
	}
	token :=authPayToken(w,auth_token);
	if token==nil {return}

	openId :=token.Claims["sub"].(string)

	account := db.QueryAccount(openId)

	if account==nil{
		account =db.NewAccount()
		account.OpenId=openId
		account.Amount = 0;
		account.CreateTime=time.Now()
		account.Status=1

		isSuccess := db.InsertAccount(account);
		if !isSuccess {
			comm.ResponseError(w,http.StatusBadRequest,"添加账户失败!")
		}
	}else{
		comm.ResponseError(w,http.StatusBadRequest,"账户已存在!")
	}
}

//充值
func Recharge(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	token :=authPayToken(w,r.Header.Get("auth_token"));
	if token==nil {return}

	var rechargeModel *RechargeModel
	comm.ReadJson(r.Body,&rechargeModel)

	result,err := WXPrepay(r)

	if err!=nil{

		comm.ResponseError(w,http.StatusBadRequest,err.Error())
	}else{
		comm.WriteJson(w,result)
	}



}

//微信预支付
func WXPrepay(r *http.Request) (pay.PaymentRequest,error) {

	notifyUrl :="http://"+":"+r.Host+"/pay/wxpay_callback"
	fmt.Println(notifyUrl);
	wxPrepay := pay.NewWXPrepay()
	wxPrepay.OrderId=pay.NewNonceString()
	wxPrepay.Amount="1"
	wxPrepay.ClientIp="192.168.0.2"
	wxPrepay.Desc="ceshixia"
	wxPrepay.NotifyUrl=notifyUrl

	return wxPrepay.Prepay();

}

//func AlipayPrepay(w http.ResponseWriter)  {
//	order :=NewAlipayOrderModel()
//	order.outTradeNO="12234445"
//	order.TotalFee = "12";
//	order.Subject="充值"
//	order.Body="用户充值"
//
//	var hash crypto.Hash =crypto.SHA256
//
//	hasher :=hash.New()
//	hasher.Write([]byte(order.ToString()))
//
//	if sigBytes, err := rsa.SignPKCS1v15(rand.Reader, keys.GetAlipayPrivateKey(), hash, hasher.Sum(nil)); err == nil {
//		signString := strings.TrimRight(base64.URLEncoding.EncodeToString(sigBytes), "=")
//		comm.WriteJson(w,order.ToString()+"&sign="+signString+"&sign_type=RSA")
//	} else {
//
//		log.Printf("",err)
//	}
//}



func authPayToken(w http.ResponseWriter,userToken string) *jwt.Token {

	if userToken==""{
		comm.ResponseError(w,http.StatusBadRequest,"认证信息不能为空!");
		return nil;
	}

	parts := strings.Split(userToken, ".")
	var token *jwt.Token;
	if len(parts)==3 {
		var err error;
		token, err = InitJWTAuthenticationBackend().FetchToken(userToken);
		if err!=nil{
			comm.ResponseError(w,http.StatusUnauthorized,"用户认证不通过!");
			return nil;
		}
	}else{
		comm.ResponseError(w,http.StatusBadRequest,"用户认证信息格式不正确!");
		return nil;
	}

	if  token.Valid {

		return token;

	} else {
		comm.ResponseError(w,http.StatusUnauthorized,"token已失效!");
		return nil;
	}
}

func appIsOk(w http.ResponseWriter,r *http.Request) bool {
	vars := mux.Vars(r);
	app_id := vars["app_id"];
	app_key := vars["app_key"];

	if err:=comm.AuthApp(app_id,app_key);err!=nil{
		comm.ResponseError(w,http.StatusUnauthorized,"appid和appkey不合法!")
		return false;
	}

	return true;
}