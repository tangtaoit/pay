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

	if orderType==Trade_Type_Recharge{

		orderPrefix=strconv.Itoa(Trade_Type_Recharge)
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





func SignStr(data string)  string {


	return fmt.Sprintf("%X",md5.Sum([]byte(data)))
}