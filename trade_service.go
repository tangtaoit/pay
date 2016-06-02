package main

import (
	"pay/db"
	"github.com/gocraft/dbr"
	"math"
	"fmt"
	"time"
	"errors"
)

func CalTradeStatus(amount,changeAmount int64,noOnce int,tradeNo string) (status int,er error) {

	tradesPays,err := db.TradePayList(tradeNo)
	if err!=nil{
		return 0,err
	}
	var tradeStatus int
	if tradesPays==nil{ //如果没有其他付款

		//应付款跟实际收到的金额不一致
		if math.Abs(float64(changeAmount))>float64(amount) {

			if noOnce==1{
				tradeStatus=Trade_Status_NOFULL
			}else{
				tradeStatus=Trade_Status_Fail
				fmt.Println(fmt.Sprintf("收到的金额为 %f  当时应收为 %f",float64(amount),math.Abs(float64(changeAmount))))
			}

		}else{
			tradeStatus=Trade_Status_Success
		}
	}else{//如果有其他付款 计算是否支付完成

		var actualedAmount int64
		for _,tradepay:= range tradesPays {
			actualedAmount =actualedAmount+tradepay.PayAmount
		}
		actualedAmount=actualedAmount+amount
		if math.Abs(float64(changeAmount))>float64(actualedAmount) {
			if noOnce==1{
				tradeStatus=Trade_Status_NOFULL
			}else{
				tradeStatus=Trade_Status_Fail
				fmt.Println(fmt.Sprintf("收到的金额为 %f  当时应收为 %f",float64(amount),math.Abs(float64(changeAmount))))
			}
		}else{
			tradeStatus=Trade_Status_Success
		}
	}

	return tradeStatus,nil
}

//交易发送改变
func TradeChange(amount int64,payType int,tradeNo string,tx *dbr.Tx) (*db.Trade,error)  {

	//查询等待交易信息
	trade := db.NewTrade()
	trade =trade.QueryByTradeNo(tradeNo)
	if trade==nil {
		return trade,errors.New(fmt.Sprintf("交易信息[%s]没有找到或交易信息状态有误!",tradeNo))
	}
	if trade.Status!=Trade_Status_Wait&&trade.Status!=Trade_Status_NOFULL{

		return trade,errors.New(fmt.Sprintf("交易信息[%s]不是合法的状态!",tradeNo))
	}


	status,err := CalTradeStatus(amount,trade.ChangedAmount,trade.NoOnce,tradeNo)
	if err!=nil{
		return trade,err
	}

	//更新实际交易金额和交易状态
	_,err =tx.Update("trades").Set("status",status).Set("update_time",time.Now()).Set("actual_amount",trade.ActualAmount+amount).Where("trade_no=?",trade.TradeNo).Exec()
	if err!=nil{
		return trade,err
	}

	//插入交易支付信息
	tradePay :=db.NewTradePay()
	tradePay.PayAmount=amount
	tradePay.CreateTime=time.Now()
	tradePay.UpdateTime=time.Now()
	tradePay.OpenId=trade.OpenId
	tradePay.PayType=payType
	tradePay.TradeNo=trade.TradeNo
	err = tradePay.InsertTx(tx)
	if err!=nil{
		return trade,err
	}

	return trade,nil;
}