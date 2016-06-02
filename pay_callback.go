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
	"time"
	"github.com/tangtaoit/queue"
)

var (
	ErrNotFoundReturnCode = errors.New("not found return_code parameter")
	ErrNotFoundResultCode = errors.New("not found result_code parameter")
	ErrNotFoundSign       = errors.New("not found sign parameter")
)

type NotifyTradeMsgModel  struct{
	//交易号
	TradeNo string
	//交易类型 1.充值 2.普通支出
	TradeType int
	//第三方系统中的交易号
	OutTradeNo string
	//第三方系统中的交易类型
	OutTradeType int
	//应用ID
	AppId string
	//用户openID
	OpenId string
	//交易时间
	TradeTime time.Time
	//交易金额
	 Amount int64
	//交易标题
	Title string
	//交易备注
	Remark string
	//交易通知地址
	NotifyUrl string

}

func NewNotifyTradeMsgModel() *NotifyTradeMsgModel {

	return &NotifyTradeMsgModel{}
}

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

				err := AccountRecharge(out_trade_no,Pay_Type_WXPAY,int64(totalFee))
				if err!=nil {
					cxt.Server.errorHandler.ServeError(cxt.ResponseWriter,cxt.Request,err)
				}else{
					//通知第三方服务器
					NotifyThirdServer(out_trade_no)
					//返回成功
					cxt.Response(map[string]string{
						"return_code": "SUCCESS",
						"return_msg": "OK",
					})
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

//通知第三方服务支付完成
func NotifyThirdServer(tradeNo string) error  {

	session := db.NewSession()
	var trade *db.Trade
	session.Select("*").From("trades").Where("trade_no=?",tradeNo).Where("status=?",Trade_Status_Success).Where("notify_status=?",0).LoadStruct(&trade)

	if trade==nil{
		fmt.Println("null")
		log.Println("warning: 交易信息没有找到 或者交易通知状态不是待通知状态",tradeNo)
		return errors.New(fmt.Sprintf("warning: 交易信息没有找到 交易号为[%s] 或者交易通知状态不是待通知状态",tradeNo))
	}
	fmt.Println("push")
	msgModel := queue.NewTradeMsg()
	msgModel.AppId=trade.AppId
	msgModel.Amount = trade.ChangedAmount
	msgModel.OpenId=trade.OpenId
	msgModel.OutTradeNo=trade.OutTradeNo
	msgModel.OutTradeType=trade.OutTradeType
	msgModel.Remark=trade.Remark
	msgModel.Title=trade.Title
	msgModel.TradeTime=trade.UpdateTime
	msgModel.TradeType=trade.TradeType
	msgModel.NotifyUrl=trade.NotifyUrl
	msgModel.TradeNo=tradeNo
	msgModel.NoOnce=trade.NoOnce
	//发送订单消息到队列
	err := queue.PublishTradeMsg(msgModel)
	
	if err==nil {
		//修改通知状态为已完成
		updateNotifyStatus(Notify_Status_Finish,tradeNo)
	}else{
		fmt.Println(err)
	}
	return err
}

func updateNotifyStatus(notifyStatus int,tradeNo string)  {
	session := db.NewSession()
	session.Update("trades").Set("notify_status",notifyStatus).Where("trade_no=?",tradeNo).Exec()
}

//账户充值
func AccountRecharge(tradeNo string,payType int,amount int64) error {

	session := db.NewSession()
	tx,er := session.Begin()
	if er!= nil {
		return er
	}
	defer func() {
		if err:=recover();err!=nil{
			log.Println(err)
			tx.Rollback()
		}
	}();

	trade,err :=TradeChange(amount,payType,tradeNo,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}

	tradeStatus,err :=CalTradeStatus(amount,trade.ChangedAmount,trade.NoOnce,tradeNo)
	if err!=nil{
		tx.Rollback()
		return err
	}

	if tradeStatus==Trade_Status_Success {
		err = db.AccountAmountChange(amount,tradeNo,trade.OpenId,trade.AppId,tx)
		if err!=nil{
			tx.Rollback()
			return err
		}
	}else{
		log.Println("充值出现错误!")
	}
	tx.Commit()

	log.Println("充值完成!")
	return nil;
}







