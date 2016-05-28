package db

import (
	"testing"
	"fmt"
	"time"
	"github.com/tangtaoit/util"
	"log"
)



func TestTrade_Update(t *testing.T) {
	sess := NewSession()
	sess.Update("trades").Where("trade_no=?","1160527553280001").Set("status",1).Exec()
}

func TestTrade_UpdateTX(t *testing.T) {
	sess := NewSession()
	tx, _ :=sess.Begin()
	tx.Update("trades").Where("trade_no=?","1160527553280001").Set("status",1).Exec()
	tx.Commit()
}

func TestTradePay_Query(t *testing.T)  {
	sess := NewSession()
	var tradePay *TradePay
	sess.Select("*").From("trades_pay").Where("pay_type=?",2).Where("trade_no=?","1160527553280001").LoadStruct(&tradePay)

	if tradePay!=nil {
		fmt.Println(tradePay.ActualAmount)
	}

}

func TestAccount_Insert(t *testing.T) {

	sess := NewSession()

	account := NewAccount()
	account.OpenId="test_openid"
	account.AppId="test_app_id"
	account.CreateTime=time.Now()
	account.UpdateTime=time.Now()
	account.Status=1
	result,err :=sess.InsertInto("accounts").Columns("open_id","app_id","status").Record(account).Exec()
	util.CheckErr(err)
	log.Println(result.LastInsertId())
}

