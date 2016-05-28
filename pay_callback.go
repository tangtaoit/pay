package main

import (
	"net/http"
	"encoding/xml"
	"io/ioutil"
	"log"
	"github.com/tangtaoit/util"
	"bytes"
	"errors"
	"fmt"
	"github.com/tangtaoit/security"
	"pay/pay"
	"strconv"
	"pay/db"
	"math"
	"time"
	"github.com/gocraft/dbr"
)

var (
	ErrNotFoundReturnCode = errors.New("not found return_code parameter")
	ErrNotFoundResultCode = errors.New("not found result_code parameter")
	ErrNotFoundSign       = errors.New("not found sign parameter")
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

func NewServer(appId, mchId, apiKey string, handler Handler, errorHandler ErrorHandler) *Server {
	if apiKey == "" {
		panic("empty apiKey")
	}
	if handler == nil {
		panic("nil Handler")
	}
	if errorHandler == nil {
		errorHandler = DefaultErrorHandler
	}

	return &Server{
		appId:        appId,
		mchId:        mchId,
		apiKey:       apiKey,
		handler:      handler,
		errorHandler: errorHandler,
	}
}

func (srv *Server) AppId() string {
	return srv.appId
}
func (srv *Server) MchId() string {
	return srv.mchId
}
func (srv *Server) ApiKey() string {
	return srv.apiKey
}

type WXNotifyReturnModel struct  {
	XMLName  xml.Name `xml:"xml"`
	ReturnCode string `xml:"return_code"`

	ReturnMsg string `xml:"return_msg"`
	
}

type Server struct {
	appId  string
	mchId  string
	apiKey string

	handler      Handler
	errorHandler ErrorHandler
}

func NewWXNotifyReturnModel(returnCode,returnMsg string) *WXNotifyReturnModel  {

	returnModel := &WXNotifyReturnModel{}
	returnModel.ReturnCode=returnCode
	returnModel.ReturnMsg=returnMsg

	return returnModel
}

//支付宝回调
func AlipayCallback(w http.ResponseWriter, r *http.Request)  {

	log.Println("支付宝回调了...");
}

func (srv *Server) Callback(w http.ResponseWriter, r *http.Request)  {

	errorHandler := srv.errorHandler

	switch r.Method {
	case "POST":
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errorHandler.ServeError(w, r, err)
			return
		}
		log.Println(string(requestBody))

		msg, err := util.DecodeXMLToMap(bytes.NewReader(requestBody))
		if err != nil {
			errorHandler.ServeError(w, r, err)
			return
		}


		returnCode, ok := msg["return_code"]
		if returnCode == ReturnCodeSuccess || !ok {
			haveAppId := msg["appid"]
			wantAppId := srv.appId
			if haveAppId != "" && wantAppId != "" && !security.SecureCompareString(haveAppId, wantAppId) {
				err = fmt.Errorf("appid mismatch, have: %s, want: %s", haveAppId, wantAppId)
				errorHandler.ServeError(w, r, err)
				return
			}

			haveMchId := msg["mch_id"]
			wantMchId := srv.mchId
			if haveMchId != "" && wantMchId != "" && !security.SecureCompareString(haveMchId, wantMchId) {
				err = fmt.Errorf("mch_id mismatch, have: %s, want: %s", haveMchId, wantMchId)
				errorHandler.ServeError(w, r, err)
				return
			}

			// 认证签名
			haveSignature, ok := msg["sign"]
			if !ok {
				err = ErrNotFoundSign
				errorHandler.ServeError(w, r, err)
				return
			}
			wantSignature := util.Sign(msg, srv.apiKey, nil)
			if !security.SecureCompareString(haveSignature, wantSignature) {
				err = fmt.Errorf("sign mismatch,\nhave: %s,\nwant: %s", haveSignature, wantSignature)
				errorHandler.ServeError(w, r, err)
				return
			}
		}

		ctx := &Context{
			Server: srv,

			ResponseWriter: w,
			Request:        r,
			RequestBody: requestBody,
			Msg:         msg,

			handlerIndex: initHandlerIndex,
		}
		srv.handler.ServeMsg(ctx)
	default:
		errorHandler.ServeError(w, r, errors.New("Unexpected HTTP Method: "+r.Method))
	}

}

func GetServerCallback()  http.HandlerFunc {

	hander := &HandlerChain{}
	hander.AppendHandler(HandlerFunc(func(cxt *Context) {

		log.Println("dealing trade...");
		out_trade_no,_ := cxt.Msg["out_trade_no"]

		if out_trade_no!=""&&len(out_trade_no)>=1 {

			orderType,_ :=strconv.Atoi(util.Substr(out_trade_no,0,1))
			totalFeeStr := cxt.Msg["total_fee"]
			totalFee,_  :=strconv.Atoi(totalFeeStr)
			if orderType==Trade_Type_Recharge { //充值
				err := AccountRecharge(out_trade_no,int64(totalFee))
				if err!=nil {
					cxt.Server.errorHandler.ServeError(cxt.ResponseWriter,cxt.Request,err)
				}
			}else{
				err :=errors.New("trade type ["+strconv.Itoa(orderType)+"] is error]")
				cxt.Server.errorHandler.ServeError(cxt.ResponseWriter,cxt.Request,err)
			}

		}else{
			err :=errors.New("out_trade_no is nil")
			cxt.Server.errorHandler.ServeError(cxt.ResponseWriter,cxt.Request,err)
		}
	}))

	callbackServer := NewServer(pay.GetWXpaySetting().AppId,"",pay.GetWXpaySetting().AppKey,hander,nil)

	return callbackServer.Callback
}




//账户充值
func AccountRecharge(tradeNo string,amount int64) error {

	//查询等待交易信息
	trade := db.NewTrade()
	trade =trade.QueryByTradeNo(tradeNo,Trade_Status_Wait)
	if trade==nil {
		return errors.New(fmt.Sprintf("trade[%s] not found or status is error!",tradeNo))
	}
	//查询除了当前支付方式的其他支付方式是否付款完成
	tradesPays := tradePaysStatusSuccess(tradeNo)

	sess := db.NewSession()
	//开启事务
	tx,err := sess.Begin()
	util.CheckErr(err)

	var tradeStatus int
	if tradesPays==nil{ //修改交易信息收到的实际金额和交易支付信息的对应支付方式的实际金额和交易成功状态

		//应付款跟实际收到的金额不一致
		if math.Abs(float64(trade.ChangedAmount))!=float64(amount) {

			tradeStatus=Trade_Status_Fail
			fmt.Println(fmt.Sprintf("weixin total_fee is %f  but ChangedAmount is %f",float64(amount),math.Abs(float64(trade.ChangedAmount))))

		}else{
			tradeStatus=Trade_Status_Success
		}
	}else{//如果其他支付方式付款完成并且加当前支付的金额等于需要支付的金额,那么交易成功完成

		var actualedAmount int64
		for _,tradepay:= range tradesPays {
			actualedAmount =actualedAmount+tradepay.ActualAmount
		}

		actualedAmount=actualedAmount+amount
		if math.Abs(float64(trade.ChangedAmount))!=float64(actualedAmount) {
			//交易失败,并记录收到的金额
			tradeStatus=Trade_Status_Fail
		}else{
			tradeStatus=Trade_Status_Success
		}
	}

	tradeChangeAndRecordActualAmount(amount,tradeStatus,trade,Pay_Type_WXPAY,tx)

	if tradeStatus==Trade_Status_Success {
		accountChange(amount,trade,tx)
	}
	tx.Commit()

	defer func() {
		if err:=recover();err!=nil{
			log.Println(err)
			tx.Rollback()

			sess.Close()

		}
	}();


		//如果其他支付方式付款未完成


	return nil;
}

//交易失败,并记录收到的金额
func tradeChangeAndRecordActualAmount(actualAmount int64,status int,trade *db.Trade,payType int,tx *dbr.Tx)  {

	//更新实际交易金额和交易状态
	_,err :=tx.Update("trades").Set("status",status).Set("actual_amount",trade.ActualAmount+actualAmount).Where("trade_no=?",trade.TradeNo).Exec()
	util.CheckErr(err)
	var tradePay *db.TradePay
	err =tx.Select("*").From("trades_pay").Where("pay_type=?",payType).Where("trade_no=?",trade.TradeNo).LoadStruct(&tradePay)
	util.CheckErr(err)
	if tradePay!=nil{
		_,err :=tx.Update("trades_pay").Where("id=?",trade.Id).Set("status",status).Set("actual_amount",tradePay.ActualAmount+actualAmount).Exec()
		util.CheckErr(err)
	}
}

func accountChange(changeAmount int64,trade *db.Trade,tx *dbr.Tx)  {

	//修改账户余额
	//如果用户没有创建账户,那么就创建一个新的账户
	//添加账户变动记录
	//将交易信息和交易支付信息状态置为成功并修改实际收到的交易金额

	var account *db.Account
	err:=tx.Select("*").From("accounts").Where("open_id=?",trade.OpenId).LoadStruct(&account)
	util.CheckErr(err)

	var amountBefrore int64
	var amountAfter int64
	if account==nil{
		amountBefrore=0
		amountAfter =changeAmount
		account = db.NewAccount()
		account.OpenId=trade.OpenId
		account.AppId=trade.AppId
		account.CreateTime=time.Now()
		account.UpdateTime=time.Now()
		account.Status=1
		account.Amount=changeAmount
		result,err :=tx.InsertInto("accounts").Columns("open_id","app_id","status","create_time","update_time","amount").Record(account).Exec()
		util.CheckErr(err)
		lastId,_ := result.LastInsertId()
		account.Id=uint64(lastId)
	}else{
		amountBefrore=account.Amount
		amountAfter = account.Amount+changeAmount

		_,err :=tx.Update("accounts").Where("id=?",account.Id).Set("amount",account.Amount+changeAmount).Exec()
		util.CheckErr(err)

	}
	accRecod :=db.NewAccountRecord()
	accRecod.TradeNo=trade.TradeNo
	accRecod.OpenId=account.OpenId
	accRecod.AccountId=account.Id
	accRecod.AmountBefore=amountBefrore
	accRecod.AmountAfter=amountAfter
	accRecod.ChangedAmount=trade.ChangedAmount
	accRecod.AppId=account.AppId
	accRecod.CreateTime=time.Now()
	_,err =tx.InsertInto("accounts_record").Columns("trade_no","app_id","open_id","account_id","create_time","amount_before","amount_after","changed_amount").Record(accRecod).Exec()
	util.CheckErr(err)

}





func tradePaysStatusSuccess(tradeNo string) []*db.TradePay {
	sess := db.NewSession()
	var tradePays []*db.TradePay
	_,err :=sess.Select("*").From("trades_pay").Where("trade_no=?",tradeNo).Where("status=?",Trade_Status_Success).LoadStructs(&tradePays)
	util.CheckErr(err)

	return tradePays
}