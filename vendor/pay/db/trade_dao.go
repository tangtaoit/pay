package db

import (
	"time"
	"database/sql"
	"github.com/tangtaoit/util"
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
	//IN 收入 OUT 支出
	InOut string
	//交易标题
	Title string
	//交易备注
	Remark string
	//状态 1.交易成功 0.待交易
	Status int

}

type TradePay struct {
	Id uint64
	//交易号
	TradeNo string

	//支付类型
	PayType int
	//支付金额
	PayAmount int64

	//实际交易金额(单位:分)
	ActualAmount int64

	//状态 1.交易成功 0.待交易
	Status int
	
}

func NewTradePay() *TradePay {

	return &TradePay{}
}

const (

	TRADE_UPDATE_STATUS_AND_ACTUALAMOUNT ="update trades set status=?,actual_amount=? where trade_no=?"

	TRADE_PAY_UPDATE_STATUS_AND_ACTUALAMOUNT ="update trades_pay set status=?,actual_amount=? where trade_no=? and pay_type=?"

	TRADE_QUERY_BY_TRADENO_AND_STATUS_SQL ="select id,trade_no,trade_type,out_trade_no,out_trade_type,app_id,open_id,create_time,update_time,changed_amount,actual_amount,in_out,title,remark,status from trades where trade_no=? and status=?"
	TRADE_INSERT_SQL string ="insert into trades(trade_no,trade_type,out_trade_no,out_trade_type,app_id,open_id,create_time,update_time,changed_amount,actual_amount,in_out,title,remark,status) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	TRADE_PAY_INSERT_SQL string="insert into trades_pay(trade_no,pay_type,pay_amount,actual_amount,status)values(?,?,?,?,?)"
)

func NewTrade()  *Trade{

	return &Trade{}
}



func (self *Trade) UpdateStatusAndActualAmountTx(tradeNo string,status int,actualAmount int64,tx *sql.Tx)  {

	Exec(tx,TRADE_UPDATE_STATUS_AND_ACTUALAMOUNT,status,actualAmount,tradeNo)
}

func (self TradePay) UpdateStatusAndActualAmountTx(tradeNo string,status int,actualAmount int64,payType int,tx *sql.Tx) {

	Exec(tx,TRADE_PAY_UPDATE_STATUS_AND_ACTUALAMOUNT,status,actualAmount,tradeNo,payType)
}

func (self *Trade) QueryByTradeNo(tradeNo string,status int) *Trade  {

	stmt,err:=GetDB().Prepare(TRADE_QUERY_BY_TRADENO_AND_STATUS_SQL)
	util.CheckErr(err)
	defer stmt.Close()

	rows,err :=stmt.Query(tradeNo,status)
	util.CheckErr(err)
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&self.Id,&self.TradeNo,&self.TradeType,&self.OutTradeNo,&self.OutTradeType,&self.AppId,&self.OpenId,&self.CreateTime,&self.UpdateTime,&self.ChangedAmount,&self.ActualAmount,&self.InOut,&self.Title,&self.Remark,&self.Status)
		util.CheckErr(err)

		return self;
	}
	return nil;

}

func (self *Trade) Insert() bool {
	stmt,err:=GetDB().Prepare(TRADE_INSERT_SQL)
	defer stmt.Close()

	util.CheckErr(err)
	_,err =stmt.Exec(self.TradeNo,self.TradeType,self.OutTradeNo,self.OutTradeType,self.AppId,self.OpenId,self.CreateTime,self.UpdateTime,self.ChangedAmount,self.ActualAmount,self.InOut,self.Title,self.Remark,self.Status)
	util.CheckErr(err)
	return true
}

func (self *Trade) InsertTx(tx *sql.Tx)   {

	stmt, err := tx.Prepare(TRADE_INSERT_SQL)
	defer stmt.Close()
	util.CheckErr(err)
	_,err =stmt.Exec(self.TradeNo,self.TradeType,self.OutTradeNo,self.OutTradeType,self.AppId,self.OpenId,self.CreateTime,self.UpdateTime,self.ChangedAmount,self.ActualAmount,self.InOut,self.Title,self.Remark,self.Status)
	util.CheckErr(err)

}

func (self *TradePay) Insert() bool {

	stmt,err:=GetDB().Prepare(TRADE_PAY_INSERT_SQL)
	defer stmt.Close()
	util.CheckErr(err)
	_,err =stmt.Exec(self.TradeNo,self.PayType,self.PayAmount,self.ActualAmount,self.Status)
	util.CheckErr(err)
	return true
}
func (self *TradePay) InsertTx(tx *sql.Tx) bool {

	stmt,err:=tx.Prepare(TRADE_PAY_INSERT_SQL)
	defer stmt.Close()
	util.CheckErr(err)
	_,err =stmt.Exec(self.TradeNo,self.PayType,self.PayAmount,self.ActualAmount,self.Status)
	util.CheckErr(err)
	return true
}