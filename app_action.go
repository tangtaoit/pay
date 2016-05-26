package main

import (
	"net/http"
	"pay/comm"
	"pay/db"
	"fmt"
)

type AppDto struct  {
	AppId string `json:"app_id"`
	AppKey string `json:"app_key"`
	AppName string `json:"app_name"`
	AppDesc string `json:"app_desc"`
	Status int `json:"status"`
}

//提交应用申请
func SubmitApp(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authorization :=GetAuthToken(r);
	token :=authPayToken(w,authorization);
	if token==nil {return}

	var appDto *AppDto
	comm.CheckErr(comm.ReadJson(r.Body,&appDto))

	app := db.NewAPP()
	app.AppId = fmt.Sprintf("%d",comm.GenerAppId())
	app.AppName = appDto.AppName
	app.OpenId = token.Claims["sub"].(string)
	app.AppDesc = appDto.AppDesc
	app.Status=0
	app.AppKey = comm.GenerUUId()

	if app.Insert()!=nil {
		comm.ResponseError(w,http.StatusBadRequest,"添加APP失败!")
		return;
	}

}