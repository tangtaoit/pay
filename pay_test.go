package main

import (
	"testing"
	"github.com/tangtaoit/util"
	"os"
	"pay/config"
	"log"
	"pay/db"
	"time"
)

func InitTest()  {
	os.Setenv("GO_ENV", "tests")
}

func TestNewAccountRecord(t *testing.T) {
	os.Setenv("GO_ENV", "tests")
	config.GetSetting()

	err := AccountRecharge("1160527760470002",Pay_Type_Account,int64(10))
	util.CheckErr(err)
}

func TestAccount_Insert(t *testing.T) {
	os.Setenv("GO_ENV", "tests")
	sess := db.NewSession()

	account := db.NewAccount()
	account.OpenId="test_openid"
	account.AppId="test_app_id"
	account.CreateTime=time.Now()
	account.UpdateTime=time.Now()
	account.Status=1
	result,err :=sess.InsertInto("accounts").Columns("open_id","app_id","status").Record(account).Exec()
	util.CheckErr(err)
	log.Println(result.LastInsertId())
}

func TestAddTradeInfo(t *testing.T) {
	os.Setenv("GO_ENV", "tests")
	//{"open_id":"20B7828A-0B26-45BF-A833-9B40A0B5CCF9","amount":1,"trade_type":1,"pay_type":2}
	tradeModel := &TradeModel{
		TradeNo:"23435",
		TradeType:Trade_Type_Recharge,
		AppId:"194669081868111872",
		OpenId:"20B7828A-0B26-45BF-A833-9B40A0B5CCF9",
		Amount:1,
		PayType:Pay_Type_WXPAY,
		Title:"充值啦",
		Description:"微信充值",
		NotifyUrl:"http://baidu.com",
	}
	addTradeInfo(tradeModel)
}

func TestAccountRecharge(t *testing.T) {
	InitTest()
	err := AccountRecharge("1160602407870000",Pay_Type_WXPAY,1)
	util.CheckErr(err)
}

func TestDB(t *testing.T)  {

	ss := db.NewSession()
	tx,_ :=ss.Begin()

	app := db.NewAPP()
	app.AppName="sdsdsd"
	tx.InsertInto("app").Columns("app_name").Record(app).Exec()

	tx.Commit()

}


