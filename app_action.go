package main

import (
	"net/http"
	"pay/comm"
	"pay/db"
	"fmt"
	"strings"
	"crypto/md5"
	"strconv"
	"time"
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
	}else{
		comm.ResponseSuccess(w)
	}

}


func AppIsOk(w http.ResponseWriter,r *http.Request) (appId string,appKey string,basesign string,isOk bool) {
	app_id := r.Header.Get("app_id");
	if app_id=="" {
		comm.ResponseError(w,http.StatusBadRequest,"app_id不能为空!");
		return "","","",false;
	}

	app := db.NewAPP()
	app = app.QueryCanUseApp(app_id)

	if app==nil {
		comm.ResponseError(w,http.StatusBadRequest,"系统中没有此应用信息!");
		return app_id,"","",false;
	}
	sign :=r.Header.Get("sign")
	if sign =="" {
		comm.ResponseError(w,http.StatusBadRequest,"签名信息(sign)不能为空!");
		return app_id,app.AppKey,"",false;
	}
	signs := strings.Split(sign,".")
	gotSign := signs[0]

	noncestr :=r.Header.Get("noncestr")
	timestamp :=r.Header.Get("timestamp")

	if noncestr=="" {
		comm.ResponseError(w,http.StatusBadRequest,"随机码不能为空!");
		return app_id,app.AppKey,"",false;
	}

	if timestamp=="" {
		comm.ResponseError(w,http.StatusBadRequest,"时间戳不能为空!");
		return app_id,app.AppKey,"",false;
	}


	timestam64,_ := strconv.ParseInt(timestamp,10,64)
	timeBtw := time.Now().Unix()-int64(timestam64)
	if timeBtw > 5*60 {
		comm.ResponseError(w,http.StatusBadRequest,"签名已失效!");
		return app_id,app.AppKey,"",false;
	}

	signStr:= fmt.Sprintf("%s%s%s",app.AppKey,noncestr,timestamp)
	wantSign :=fmt.Sprintf("%X",md5.Sum([]byte(signStr)))

	if gotSign!=wantSign {
		fmt.Println("wantSign: ",wantSign,"gotSign: ",gotSign)
		comm.ResponseError(w,http.StatusBadRequest,"请求不合法!");
		return app_id,app.AppKey,"",false;
	}

	if app==nil{
		comm.ResponseError(w,http.StatusUnauthorized,"应用信息未找到!请检查APPID是否正确");
		return app_id,app.AppKey,"",false;
	}

	return app_id,app.AppKey,gotSign,true;
}