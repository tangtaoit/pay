package main

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"pay/config"
	"pay/db"
	"time"
	"fmt"
	"pay/pay"
	"strconv"
	"github.com/tangtaoit/util"
	"log"
	"io/ioutil"

	"encoding/json"
)



type PayToken struct  {

	PayToken string `json:"pay_token"`

}

type SignModel struct  {

	Sign string `json:"sign"`
	Noncestr string `json:"noncestr"`
	Timestamp string `json:"timestamp"`
}

//账户预支付body
type AccountPrepayBody struct {

	TradeNo string `json:"TradeNo"`

}

func NewAccountPrepayBody(tradeNo string) *AccountPrepayBody {

	return &AccountPrepayBody{TradeNo:tradeNo}
}

type PrepayWrap  struct {

	PayType int `json:"pay_type"`

	PayNo string `json:"pay_no"`

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

	TradeNo string `json:"trade_no,omitempty"`

	//交易类型
	TradeType  int  `json:"trade_type"`

	//支付类型 1.支付宝  2.微信
	PayType int `json:"pay_type"`

	//交易标题
	Title string `json:"title"`
	//描述
	Description string `json:"description"`
	ClientIp string `json:"client_ip"`
	//通知地址
	NotifyUrl string `json:"notify_url"`
	//0.一次支付 1.分批支付
	NoOnce int `json:"no_once"`
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

//绑定支付信息
type BindPayInfoModel struct{

	Password string `json:"password"`

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
	appId,_,_,isOk:=AppIsOk(w,r);
	if !isOk{
		return;
	}
	authorization :=GetAuthToken(r);
	token :=authUserToken(w,authorization);
	if token==nil {return}

	openId := token.Claims["sub"].(string)
	//支付密码
	password := r.FormValue("password")
	if password=="" {
		util.ResponseError(w,http.StatusBadRequest,"请输入支付密码!")
		return;
	}

	account := db.NewAccount()
	account,err :=account.QueryAccount(openId,appId)
	log.Println(err)
	if account==nil {
		util.ResponseError(w,http.StatusBadRequest,"请先设置账户密码!")
		return;
	}else{
		if account.Password!=password {
			util.ResponseError(w,http.StatusBadRequest,"密码错误!")
			return;
		}
	}
	paytoken := util.GenerUUId();
	key :=PAY_TOKEN_PREFIX+paytoken+":"+openId;
	SetAndExpire(key,"1",config.GetSetting().TokenExpire)
	util.WriteJson(w,NewPayToken(paytoken))
}



// PUT pay/password
func SetAccountPassword(w http.ResponseWriter, r *http.Request)  {

	fmt.Println("SetAccountPassword...")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	appId,_,_,isOk:=AppIsOk(w,r)
	if !isOk{
		return;
	}

	var bodyParams map[string]interface{}
	util.ReadJson(r.Body,&bodyParams)

	 newpwd,ok := bodyParams["newpwd"].(string)
	if !ok{
		newpwd=""
	}
	 oldpwd,ok :=bodyParams["oldpwd"].(string)
	if !ok{
		oldpwd=""
	}

	openId := bodyParams["open_id"].(string)
	if !ok{
		openId=""
	}


	if len(newpwd)!= 6 {
		util.ResponseError(w,http.StatusBadRequest,"密码长度必须为6位!")
		return;
	}

	if _,er := strconv.Atoi(newpwd);er!=nil{
		util.ResponseError(w,http.StatusBadRequest,"密码必须为纯数字!")
		return;
	}

	if openId==""{
		util.ResponseError(w,http.StatusBadRequest,"用户ID不能为空!")
		return;
	}


	account := db.NewAccount()
	var err error
	account,err = account.QueryAccount(openId,appId)
	if err!=nil{
		log.Println("查询账户信息错误",err,account)
	}
	if account==nil{
		account = db.NewAccount()
		account.OpenId=openId
		account.Amount=0
		account.CreateTime=time.Now()
		account.UpdateTime=time.Now()
		account.AppId=appId
		account.Password=newpwd
		account.Status=Account_Status_Enable
		account.Insert()
	}else{

		fmt.Println("account=",account.Password)

		if account.Password!=""&&oldpwd==""{ //系统中存在支付密码的话需要用户输入旧密码进行修改
			util.ResponseError(w,http.StatusBadRequest,"请输入旧密码!.")
			return
		}else if newpwd!=""&&oldpwd!=account.Password { //旧密码校验失败
			util.ResponseError(w,http.StatusBadRequest,"旧密码不正确!.")
			return
		}
	}

	//修改新密码
	err =account.UpdatePwd(openId,newpwd)
	if err!=nil{
		util.ResponseError(w,http.StatusBadRequest,"修改密码失败!.")
		return
	}
	util.ResponseSuccess(w)
	return;
}



//绑定支付信息
func BindPayInfo(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	appId,appKey,baseSign,isOk:=AppIsOk(w,r)
	if !isOk{
		return;
	}
	sign := r.Header.Get("sign")
	data,err := ioutil.ReadAll(r.Body)
	err = SignApi(sign,data,appId,appKey,baseSign)
	if err!=nil{
		util.ResponseError(w,http.StatusBadRequest,err.Error())
		return
	}
	var resultMap map[string]interface{}
	util.CheckErr(util.ReadJsonByByte(data,&resultMap))

	openId := resultMap["open_id"].(string)

	account :=db.NewAccount()
	account,_ = account.QueryAccount(openId,appId)
	if account==nil{
		account =db.NewAccount()
		account.OpenId=openId
		account.Amount = 0;
		account.CreateTime=time.Now()
		account.Status=1

		err := account.Insert();
		if err!=nil {
			util.ResponseError(w,http.StatusBadRequest,"添加账户失败!")
			return
		}
	}else{
		util.ResponseError(w,http.StatusBadRequest,"账户已存在!")
	}
}


//预备支付接口(付款之前调用)
func MakePrePay(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	appId,appKey,basesign,isOk:=AppIsOk(w,r)
	if !isOk{
		return;
	}

	sign := r.Header.Get("sign")
	data,err := ioutil.ReadAll(r.Body)
	err = SignApi(sign,data,appId,appKey,basesign)
	if err!=nil{
		util.ResponseError(w,http.StatusBadRequest,err.Error())
		return
	}


	util.CheckErr(err)

	var tradeModel *TradeModel
	util.CheckErr(util.ReadJsonByByte(data,&tradeModel))
	//util.CheckErr(util.ReadJson(r.Body,&tradeModel))

	if tradeModel.OpenId=="" {
		util.ResponseError(w,http.StatusBadRequest,"用户open_id不能为空!")
		return
	}

	if tradeModel.NotifyUrl=="" {
		util.ResponseError(w,http.StatusBadRequest,"通知地址[notify_url]不能为空!")
		return
	}
	if !(strings.HasPrefix(tradeModel.NotifyUrl,"http:")||strings.HasPrefix(tradeModel.NotifyUrl,"https:")){
		util.ResponseError(w,http.StatusBadRequest,"通知地址[notify_url]格式不正确,没有前缀http或者https!")
		return
	}


	if tradeModel.TradeType ==Trade_Type_Recharge {
		tradeModel.Title="账户充值"
		tradeModel.Description="账户充值"
	}

	if tradeModel.Amount<=0 {
		util.ResponseError(w,http.StatusBadRequest,"交易金额不能小于或等于0");
		return;
	}
	orderNo,err :=NewOrderNo(tradeModel.TradeType)
	if err!=nil{
		util.ResponseError(w,http.StatusBadRequest,err.Error());
		return;
	}

	account := db.NewAccount()
	account,err = account.QueryAccount(tradeModel.OpenId,appId)
	if err!=nil{
		log.Println("查询账户信息错误",err)
	}
	if account==nil {
		util.ResponseError(w,http.StatusBadRequest,"此用户还未开通支付功能!")
		return
	}

	var result interface{}
	if tradeModel.PayType==Pay_Type_WXPAY { //微信支付
		result,err = WXPrepay(r,tradeModel,orderNo)
		util.CheckErr(err)
	}else if tradeModel.PayType==Pay_Type_Account { //账户余额支付
		if tradeModel.TradeType==Trade_Type_Recharge {
			util.ResponseError(w,http.StatusBadRequest,"不支持账户余额对账户进行充值!")
			return
		}
		result = NewAccountPrepayBody(orderNo)
	}else {
		util.ResponseError(w,http.StatusBadRequest,"不支持的支付方式["+strconv.Itoa(tradeModel.PayType)+"]")
		return
	}

	if err!=nil{
		util.ResponseError(w,http.StatusBadRequest,err.Error())
	}else{
		tradeModel.TradeNo = orderNo
		tradeModel.AppId=appId
		err := addTradeInfo(tradeModel)
		if err!=nil{
			log.Println("订单信息插入失败[%s]",err)
			util.ResponseError(w,http.StatusBadRequest,"订单信息插入失败!")
		}else{
			prepayWrap := NewPrepayWrap(tradeModel.PayType,result)
			prepayWrap.PayNo=orderNo
			util.WriteJson(w,prepayWrap)
		}
	}
}

//付款(只限账户)
func MakePay(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	appId,appKey,baseSign,isOk:=AppIsOk(w,r)
	if !isOk{
		return;
	}
	sign := r.Header.Get("sign")
	data,err := ioutil.ReadAll(r.Body)
	err =SignApi(sign,data,appId,appKey,baseSign)
	if err!=nil{
		util.ResponseError(w,http.StatusBadRequest,err.Error())
		return
	}

	var resultMap map[string]interface{}
	util.CheckErr(util.ReadJsonByByte(data,&resultMap))

	if _,ok := resultMap["pay_token"];!ok{
		util.ResponseError(w,http.StatusBadRequest,"支付token不能为空!")
		return
	}
	if _,ok := resultMap["open_id"];!ok{
		util.ResponseError(w,http.StatusBadRequest,"用户ID不能为空!")
		return
	}

	if _,ok := resultMap["amount"];!ok{
		util.ResponseError(w,http.StatusBadRequest,"付款金额不能为空!")
		return
	}
	if _,ok := resultMap["pay_no"];!ok{
		util.ResponseError(w,http.StatusBadRequest,"付费编号不能为空!")
		return
	}

	payToken := resultMap["pay_token"].(string)
	openId := resultMap["open_id"].(string)
	amount,_:=resultMap["amount"].(json.Number).Int64();
	payNo := resultMap["pay_no"].(string)


	key :=PAY_TOKEN_PREFIX+payToken+":"+openId;
	if amount<=0 {
		util.ResponseError(w,http.StatusBadRequest,"付款金额不能等于或小于0!")
		return
	}

	isVail := GetString(key)
	if isVail!="1" {
		util.ResponseError(w,http.StatusBadRequest,"支付token不存在或已失效!")
		return
	}

	account := db.NewAccount()
	account,_ = account.QueryAccount(openId,appId)
	if account==nil{
		util.ResponseError(w,http.StatusBadRequest,"此账户还没有开通支付功能!")
		return
	}

	if account.Amount<int64(amount) {
		util.ResponseError(w,http.StatusBadRequest,"余额不足以支付!")
		return
	}

	trade := db.NewTrade()
	trade =trade.QueryByTradeNo(payNo)
	if trade==nil{
		util.ResponseError(w,http.StatusBadRequest,"没有找到对应的付款信息!")
		return
	}
	if trade.Status!=Trade_Status_Wait&&trade.Status!=Trade_Status_NOFULL{
		util.ResponseError(w,http.StatusBadRequest,"交易信息不是待支付状态!")
		return
	}

	sess :=db.NewSession()
	tx,err := sess.Begin()
	defer func() {
		if er := recover();er!=nil{
			fmt.Println(er)
			util.ResponseError(w,500,"系统错误!")
			tx.Rollback()
		}
	}()

	trade,err=TradeChange(amount,Pay_Type_Account,trade.TradeNo,tx)
	util.CheckErr(err)

	//计算当前交易状态
	tradeStatus,err :=CalTradeStatus(amount,trade.ChangedAmount,trade.NoOnce,trade.TradeNo)
	util.CheckErr(err)

	//如果交易状态为成功
	if tradeStatus!=Trade_Status_Fail {
		err = db.AccountAmountChange(-amount,trade.TradeNo,trade.OpenId,trade.AppId,tx)
		util.CheckErr(err)
	}
	tx.Commit()

	util.ResponseSuccess(w)
}




func addTradeInfo(tradeModel *TradeModel) error {

	session := db.NewSession()
	tx,err := session.Begin()
	util.CheckErr(err)

	trade := db.NewTrade()
	trade.TradeNo=tradeModel.TradeNo
	trade.TradeType=tradeModel.TradeType
	trade.AppId =tradeModel.AppId
	trade.OpenId = tradeModel.OpenId
	trade.CreateTime = time.Now()
	trade.UpdateTime =time.Now()
	trade.ChangedAmount = tradeModel.Amount

	trade.Title = tradeModel.Title
	trade.Remark = tradeModel.Description
	trade.NotifyUrl=tradeModel.NotifyUrl
	trade.NotifyStatus=Notify_Status_Wait
	trade.Status = Trade_Status_Wait
	trade.NoOnce=tradeModel.NoOnce
	err = trade.InsertTx(tx)
	util.CheckErr(err)
	err = tx.Commit()
	util.CheckErr(err)

	defer func() {
		if err:=recover();err!=nil{
			log.Println(err)
			tx.Rollback()
		}
	}()
	return nil
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



func authUserToken(w http.ResponseWriter,userToken string) *jwt.Token {

	if userToken==""{
		util.ResponseError(w,http.StatusBadRequest,"认证信息不能为空!");
		return nil;
	}

	parts := strings.Split(userToken, ".")
	var token *jwt.Token;
	if len(parts)==3 {
		var err error;
		token, err = InitJWTAuthenticationBackend().FetchToken(userToken);
		if err!=nil{
			util.ResponseError(w,http.StatusUnauthorized,"用户认证不通过!");
			return nil;
		}
	}else{
		util.ResponseError(w,http.StatusBadRequest,"用户认证信息格式不正确!");
		return nil;
	}

	if  token.Valid {

		return token;

	} else {
		util.ResponseError(w,http.StatusUnauthorized,"token已失效!");
		return nil;
	}
}
