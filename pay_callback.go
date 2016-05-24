package main

import (
	"net/http"
	"log"
)


//支付宝回调
func AlipayCallback(w http.ResponseWriter, r *http.Request)  {

	log.Println("支付宝回调了...");
}

func WXpayCallback(w http.ResponseWriter, r *http.Request)  {

	log.Println("微信回调了...");
}