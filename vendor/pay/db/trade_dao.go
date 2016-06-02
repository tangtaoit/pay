package db

import (
	"time"
	"github.com/gocraft/dbr"
)

type Trade struct {

	Id uint64
	//交易号
	TradeNo string
	//交易类型 1.充值 2.普通支出
	TradeType int
	//第三方系统中的交易号
	OutTradeNo string
	//第三方系统中的交易类型
	OutTradeType int
	//通知地址
	NotifyUrl string
	//通知状态
	NotifyStatus int
	//应用ID
	AppId string
	//用户openID
	OpenId string
	//创建时间
	CreateTime time.Time
	//更新时间
	UpdateTime time.Time
	//变动金额(单位:分)
	ChangedAmount int64
	//实际交易金额(单位:分)
	ActualAmount int64
	//交易标题
	Title string
	//交易备注
	Remark string
	//状态 1.交易成功 0.待交易
	Status int
	//是否必须一次付清
	NoOnce int

}

type TradePay struct {
	Id uint64
	//交易号
	TradeNo string
	//付款人ID
	OpenId string
	//支付类型
	PayType int
	//创建时间
	CreateTime time.Time
	//修改时间
	UpdateTime time.Time
	//付款金额
	PayAmount int64
	
}

func NewTradePay() *TradePay {

	return &TradePay{}
}


func NewTrade()  *Trade{

	return &Trade{}
}



func (self *Trade) QueryByTradeNo(tradeNo string) *Trade  {

	session := NewSession()
	var trade *Trade
	session.Select("*").From("trades").Where("trade_no=?",tradeNo).LoadStructs(&trade)

	return trade;

}

func (self *Trade) InsertTx(tx *dbr.Tx) error {

	result,err :=tx.InsertInto("trades").Columns("trade_no","trade_type","out_trade_no","out_trade_type","app_id","open_id","create_time","update_time","changed_amount","actual_amount","title","remark","notify_url","notify_status","status","no_once").Record(self).Exec()

	tradeId,er := result.LastInsertId()
	if er!=nil{
		return er
	}
	self.Id=uint64(tradeId)
	return err
}

func (self *TradePay) InsertTx(tx *dbr.Tx) error  {

	result,err := tx.InsertInto("trades_pay").Columns("trade_no","open_id","pay_type","pay_amount","create_time","update_time").Record(self).Exec()

	if err==nil{
		id,er := result.LastInsertId()
		if er!=nil{
			return er
		}
		self.Id=uint64(id)
	}

	return err
}


//查询已交易支付信息
func TradePayList(tradeNo string) ([]*TradePay,error) {
	sess := NewSession()
	var tradePays []*TradePay
	_,err :=sess.Select("*").From("trades_pay").Where("trade_no=?",tradeNo).LoadStructs(&tradePays)
	if err!=nil{
		return nil,err
	}

	return tradePays,nil
}

