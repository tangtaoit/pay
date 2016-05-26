package main

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"pay/config"
	"pay/comm"
	"pay/db"
	"time"
	"fmt"
	"pay/pay"
	"strconv"
)

const (
	//等待交易
	Trade_Status_Wait =0

	//交易成功
	Trade_Status_Success =1

	//交易失败
	Trade_Status_Fail =2
)

//支付类型
const (
	//支付宝
	Pay_Type_AliPay = 1

	//微信支付
	Pay_Type_WXPAY =2
)


//订单类型
const  (

	//充值订单
	Order_Type_Recharge = 1
	//普通订单
	Order_Type_CommOrder =2
)

type PayToken struct  {

	PayToken string `json:"pay_token"`

}

type SignModel struct  {

	Sign string `json:"sign"`
	Noncestr string `json:"noncestr"`
	Timestamp string `json:"timestamp"`
}

type PrepayWrap  struct {

	PayType int `json:"pay_type"`

	Body interface{} `json:"body"`

}

func NewPrepayWrap(payType int,body interface{}) *PrepayWrap  {

	wrap :=&PrepayWrap{}
	wrap.PayType=payType
	wrap.Body = body

	return wrap
}

//交易model
type TradeModel struct {

	SignModel

	AppId string `json:"app_id"`

	//用户OPENID
	OpenId string `json:"open_id"`

	//交易金额(单位分)
	Amount int64 `json:"amount"`

	//交易类型
	TradeType  int  `json:"trade_type"`

	//支付类型 1.支付宝  2.微信
	PayType int `json:"pay_type"`

	//交易标题
	Title string `json:"title"`
	//描述
	Description string `json:"description"`

	ClientIp string `json:"client_ip"`
}

func (self *TradeModel) toMap() (map[string]string,error)  {

	return  pay.ToMapOfJson(self)
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

	if _,_,_,isOk:=AppIsOk(w,r);!isOk{
		return;
	}
	authorization :=GetAuthToken(r);
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
	appid,_,_,isOk:=AppIsOk(w,r)
	if !isOk{
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
	account := db.QueryAccount(openId,appid)
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

//统一下单接口
func MakePrePayOrder(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	appid,appkey,basesign,isOk:=AppIsOk(w,r)
	if !isOk{

		return;
	}
	sign := r.Header.Get("sign")
	signs :=strings.Split(sign,".")
	if len(signs)!=2 {
		comm.ResponseError(w,http.StatusBadRequest,"非法请求!")
		return
	}

	var tradeModel *TradeModel
	comm.CheckErr(comm.ReadJson(r.Body,&tradeModel))

	if tradeModel.OpenId=="" {
		comm.ResponseError(w,http.StatusBadRequest,"用户open_id不能为空!")
		return
	}

	signStr :=fmt.Sprintf("%s%s%s%s%d%d%d",basesign,appid,appkey,tradeModel.OpenId,tradeModel.Amount,tradeModel.PayType,tradeModel.TradeType)

	wantSign := SignStr(signStr)
	gotSign :=signs[1];
	if wantSign!=gotSign {
		fmt.Println("wantSign: ",wantSign,"gotSign: ",gotSign)
		comm.ResponseError(w,http.StatusBadRequest,"签名不匹配!")
		return
	}

	//pay.Sign(tradeModel.toMap(),)

	if tradeModel.TradeType ==Order_Type_Recharge {
		tradeModel.Title="账户充值"
		tradeModel.Description="账户充值"
	}

	if tradeModel.Amount<=0 {
		comm.ResponseError(w,http.StatusBadRequest,"充值金额不能小于或等于0");
		return;
	}
	orderNo,err :=NewOrderNo(tradeModel.TradeType)
	if err!=nil{
		comm.ResponseError(w,http.StatusBadRequest,err.Error());
		return;
	}

	var result interface{}
	if tradeModel.PayType==Pay_Type_WXPAY {
		result,err = WXPrepay(r,tradeModel,orderNo)
		comm.CheckErr(err)
	}else {
		comm.ResponseError(w,http.StatusBadRequest,"不支持的支付方式["+strconv.Itoa(tradeModel.PayType)+"]")
		return
	}

	if err!=nil{
		comm.ResponseError(w,http.StatusBadRequest,err.Error())
	}else{
		trade := db.NewTrade()
		trade.TradeNo=orderNo
		trade.TradeType=Order_Type_Recharge
		trade.AppId =appid
		trade.OpenId = tradeModel.OpenId
		trade.CreateTime = time.Now()
		trade.UpdateTime =time.Now()
		trade.ChangedAmount = tradeModel.Amount
		trade.InOut = "IN"
		trade.Title = tradeModel.Title
		trade.Remark = tradeModel.Description
		trade.Status = Trade_Status_Wait

		if trade.Insert() {
			comm.WriteJson(w,NewPrepayWrap(tradeModel.PayType,result))
		}else{
			comm.ResponseError(w,http.StatusBadRequest,"订单入库出错!")
		}
	}
}

//微信预支付
func WXPrepay(r *http.Request,tradeModel *TradeModel,orderNo string) (pay.PaymentRequest,error) {


	notifyUrl :="http://"+r.Host+"/pay/wxpay_callback"
	wxPrepay := pay.NewWXPrepay()
	wxPrepay.OrderId=orderNo
	wxPrepay.Amount= fmt.Sprintf("%d",tradeModel.Amount)
	wxPrepay.ClientIp=tradeModel.ClientIp
	wxPrepay.Desc=tradeModel.Description
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
