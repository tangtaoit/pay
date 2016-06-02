package main

import (
	"fmt"
	"time"
	"strconv"
	"sync"
	"errors"
	"net/http"
	"crypto/md5"
	"github.com/tangtaoit/util"
	"strings"
	"pay/db"
)

var ordermx sync.Mutex

//如果是分布式的话 此处要修改,不能声明一个变量 需要所有分布式程序都统一读取同一个原子数据
var orderInc int

func GetAuthToken(r *http.Request) string {

	authorization :=r.FormValue("auth_token");

	if authorization=="" {

		return r.Header.Get("auth_token")
	}
	return ""
}

func NewOrderNo(orderType int) (string,error) {

	ordermx.Lock()

	defer ordermx.Unlock()

	var orderPrefix string

	if orderType==Trade_Type_Recharge{ //充值
		orderPrefix=strconv.Itoa(Trade_Type_Recharge)
	}else if orderType==Trade_Type_Buy{ //购买
		orderPrefix=strconv.Itoa(Trade_Type_Buy)
	}else{
		return "",errors.New("没有声明的交易类型["+strconv.Itoa(orderType)+"]")
	}

	nw := time.Now()
	yy := util.Substr(strconv.Itoa(nw.Year()),2,2)
	var mm string
	if nw.Month()<=9 {

		mm=fmt.Sprintf("0%d",nw.Month())
	}else{
		mm=fmt.Sprintf("%d",nw.Month())
	}

	var dd string
	if nw.Day()<=9 {
		dd=fmt.Sprintf("0%d",nw.Day())
	}else{
		dd=fmt.Sprintf("%d",nw.Day())
	}
	sedStr := strconv.Itoa(nw.Hour()*60*60+nw.Minute()*60+nw.Second())

	for len(sedStr)<5 {
		sedStr="0"+sedStr
	}
	if orderInc>=9999 {
		orderInc=0
	}

	orderStr := strconv.Itoa(orderInc)

	for len(orderStr)<4{
		orderStr="0"+orderStr;
	}

	orderInc++

	return fmt.Sprintf("%s%s%s%s%s%s",orderPrefix,yy,mm,dd,sedStr,orderStr),nil
}




func SignApi(sign string,data []byte,appId,appKey,baseSign string) error {

	if sign=="" {
		return errors.New("非法请求!")
	}

	signs :=strings.Split(sign,".")

	var gotSign string
	if len(signs)>=2 {
		gotSign =signs[1];
	}else{
		gotSign =signs[0];
	}


	var signMap  map[string]interface{}
	err :=util.ReadJsonByByte(data,&signMap)
	if err!=nil{
		return err
	}

	wantSign := util.SignWithBaseSign(signMap,appKey,baseSign,nil)

	if wantSign!=gotSign {

		return errors.New("签名不匹配!")
	}
	return err
}

func AppIsOk(w http.ResponseWriter,r *http.Request) (appId string,appKey string,basesign string,isOk bool) {
	app_id := r.Header.Get("app_id");
	if app_id=="" {
		util.ResponseError(w,http.StatusBadRequest,"app_id不能为空!");
		return "","","",false;
	}

	app := db.NewAPP()
	app,_ = app.QueryCanUseApp(app_id)
	if app==nil {
		util.ResponseError(w,http.StatusBadRequest,"系统中没有此应用信息!");
		return app_id,"","",false;
	}
	sign :=r.Header.Get("sign")
	if sign =="" {
		util.ResponseError(w,http.StatusBadRequest,"签名信息(sign)不能为空!");
		return app_id,app.AppKey,"",false;
	}
	signs := strings.Split(sign,".")
	gotSign := signs[0]

	noncestr :=r.Header.Get("noncestr")
	timestamp :=r.Header.Get("timestamp")

	if noncestr=="" {
		util.ResponseError(w,http.StatusBadRequest,"随机码不能为空!");
		return app_id,app.AppKey,"",false;
	}

	if timestamp=="" {
		util.ResponseError(w,http.StatusBadRequest,"时间戳不能为空!");
		return app_id,app.AppKey,"",false;
	}


	timestam64,_ := strconv.ParseInt(timestamp,10,64)
	timeBtw := time.Now().Unix()-int64(timestam64)
	if timeBtw > 5*60*1000 {
		util.ResponseError(w,http.StatusBadRequest,"签名已失效!");
		return app_id,app.AppKey,"",false;
	}

	signStr:= fmt.Sprintf("%s%s%s",app.AppKey,noncestr,timestamp)
	wantSign :=fmt.Sprintf("%X",md5.Sum([]byte(signStr)))

	if gotSign!=wantSign {
		fmt.Println("wantSign: ",wantSign,"gotSign: ",gotSign)
		util.ResponseError(w,http.StatusBadRequest,"请求不合法!");
		return app_id,app.AppKey,"",false;
	}

	if app==nil{
		util.ResponseError(w,http.StatusUnauthorized,"应用信息未找到!请检查APPID是否正确");
		return app_id,app.AppKey,"",false;
	}

	return app_id,app.AppKey,gotSign,true;
}