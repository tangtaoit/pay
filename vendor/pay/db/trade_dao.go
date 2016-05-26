package db

import (
	"time"
	"pay/comm"
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
	//IN 收入 OUT 支出
	InOut string
	//交易标题
	Title string
	//交易备注
	Remark string
	//状态 1.交易成功 0.待交易
	Status int

}

func NewTrade()  *Trade{

	return &Trade{}
}

func (self *Trade) Insert() bool {
	stmt,err:=db.Prepare("insert into trades(trade_no,trade_type,out_trade_no,out_trade_type,app_id,open_id,create_time,update_time,changed_amount,in_out,title,remark,status) values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
	comm.CheckErr(err)
	_,err =stmt.Exec(self.TradeNo,self.TradeType,self.OutTradeNo,self.OutTradeType,self.AppId,self.OpenId,self.CreateTime,self.UpdateTime,self.ChangedAmount,self.InOut,self.Title,self.Remark,self.Status)
	comm.CheckErr(err)
	return true

}