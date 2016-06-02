package main

import (
	"testing"
	"github.com/codegangsta/negroni"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"fmt"
	"encoding/json"
	"github.com/tangtaoit/util"
	"time"
	"crypto/md5"
	"bytes"
	"github.com/tangtaoit/queue"
)

var server *negroni.Negroni

var basesign string
var noncestr string
var timestamp string
var apikey string
var appid string
var auth_token string
var openId string

func initSetting()  {

	os.Setenv("GO_ENV", "tests")

	appid="196124939491741696"
	apikey ="4537C07A563C4899B5A592DA3CC84A10"
	noncestr ="23435"
	openId="ECE50557AA3047E4BB00D60B141C7066"
	timestamp =fmt.Sprintf("%d",time.Now().Unix())

	auth_token="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBfaWQiOiIxOTQyNTEyNzc5ODE0NTQzMzYiLCJleHAiOjE0NjQ5NjAxODQsImlhdCI6MTQ2NDcwMDk4NCwic3ViIjoiRUNFNTA1NTdBQTMwNDdFNEJCMDBENjBCMTQxQzcwNjYifQ.dprwXU-v_Uvpw-6l1rQuFnmRxauLzCMHCVE3rmLkbtWyg4-ujQ4AU9Q98lGEPbHnFA4v5zPvdo0mX8K7wgvVlDqMKa08Vabh7bd718myq6z-gkluNnwm0ofGSliOOLf1Eld7u7O-aadqq9xK7I9i4muT6MgfPrHCYsLZSXv3UnE"

	signStr := apikey+noncestr+timestamp
	bytes  := md5.Sum([]byte(signStr))
	basesign =fmt.Sprintf("%X",bytes)

	router :=GetRouters()
	server = negroni.Classic()
	server.UseHandler(router)
}

func TestSubmitApp(t *testing.T) {
	initSetting()
	resource := "/pay/app"
	params :=map[string]interface{}{
		"app_name":"测试应用",
		"app_desc":"应用描述",
	}

	sign := util.SignWithBaseSign(params,apikey,basesign,nil)
	paramsBytes,err := json.Marshal(params)
	util.CheckErr(err)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", resource, bytes.NewReader(paramsBytes))
	request.Header.Set("sign",fmt.Sprintf("%s.%s",basesign,sign))
	request.Header.Set("noncestr",noncestr)
	request.Header.Set("timestamp",timestamp)

	//request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	server.ServeHTTP(response, request)

	var result *util.ResultError
	err = json.Unmarshal(response.Body.Bytes(),&result)
	util.CheckErr(err)

	fmt.Println("result =%s",result)

	assert.Equal(t, response.Code, http.StatusOK)
}

//下单测试
func TestMakePrePay(t *testing.T) {

	initSetting()
	resource := "/pay/makeprepay"
	params :=map[string]interface{}{
		"pay_type":2,
		"open_id":openId,
		"trade_type":1,
		"amount":100,
		"notify_url":"http://127.0.0.1:8080/pay/test",
		"no_once":1,
	}

	sign := util.SignWithBaseSign(params,apikey,basesign,nil)
	paramsBytes,err := json.Marshal(params)
	util.CheckErr(err)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", resource, bytes.NewReader(paramsBytes))
	setAuthHeader(request)
	request.Header.Set("sign",fmt.Sprintf("%s.%s",basesign,sign))

	//request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	server.ServeHTTP(response, request)

	var result map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(),&result)
	util.CheckErr(err)

	fmt.Println("result =%s",result)

	assert.Equal(t, response.Code, http.StatusOK)
}

func TestSetAccountPassword(t *testing.T) {
	initSetting()
	resource := "/pay/password"
	params :=map[string]interface{}{
		"newpwd":"123456",
	}

	sign := util.SignWithBaseSign(params,apikey,basesign,nil)
	paramsBytes,err := json.Marshal(params)
	util.CheckErr(err)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", resource, bytes.NewReader(paramsBytes))
	setAuthHeader(request)
	request.Header.Set("sign",fmt.Sprintf("%s.%s",basesign,sign))
	request.Header.Set("auth_token",auth_token)
	server.ServeHTTP(response, request)

	var result *util.ResultError
	err = json.Unmarshal(response.Body.Bytes(),&result)
	util.CheckErr(err)

	fmt.Println("result=",result)

	assert.Equal(t, response.Code, http.StatusOK)
}

func TestNotifyThirdServer(t *testing.T) {

	initSetting()
	queue.SetupAMQP("")

	NotifyThirdServer("1160531590220000")
}

//测试获取付款token
func TestGetPayToken(t *testing.T) {
	initSetting()
	resource := "/pay/token?password=123456"

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	setAuthHeader(request)
	request.Header.Set("sign",fmt.Sprintf("%s.%s",basesign,""))
	server.ServeHTTP(response, request)

	var result map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(),&result)
	util.CheckErr(err)

	fmt.Println("result=",result)

	assert.Equal(t, response.Code, http.StatusOK)
}

//付款测试
func TestMakePay(t *testing.T) {
	initSetting()
	resource := "/pay/makepay"
	params :=map[string]interface{}{
		"pay_token":"672E518475C7417785456CA8A267F1B0",
		"open_id":openId,
		"amount":91,
		"pay_no":"1160602449370000",
	}
	sign := util.SignWithBaseSign(params,apikey,basesign,nil)
	paramsBytes,err := json.Marshal(params)
	util.CheckErr(err)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", resource, bytes.NewReader(paramsBytes))
	setAuthHeader(request)
	request.Header.Set("sign",fmt.Sprintf("%s.%s",basesign,sign))
	server.ServeHTTP(response, request)

	var result map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(),&result)
	util.CheckErr(err)

	fmt.Println("result=",result)

	assert.Equal(t, response.Code, http.StatusOK)
}

func setAuthHeader(request *http.Request)  {
	request.Header.Set("app_id",appid)
	request.Header.Set("noncestr",noncestr)
	request.Header.Set("timestamp",timestamp)
	request.Header.Set("auth_token",auth_token)
}

func TestTime(t *testing.T) {

	fmt.Println(time.Now())
}
