package comm

import (
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"bytes"
	"github.com/sumory/idgen"
)


//认证APP是否合法
func AuthApp(appId string,appKey string)  error{

	return nil;
}


func CheckErr(err error)  {
	if err != nil {
		panic(err)
	}
}

func ResponseError(w http.ResponseWriter, statusCode int,msg string)  {
	err := ResultError{statusCode, msg}
	if jsonData,er := json.Marshal(err);er==nil{
		http.Error(w,string(jsonData),err.ErrCode)
		return;
	}
	http.Error(w,"未知错误",500);
}

func ResponseSuccess(w http.ResponseWriter)  {

	err := NewResultError(0,"OK")
	WriteJson(w,err)

}

func WriteJson(w io.Writer,obj interface{})  {

	jsonData,_:= json.Marshal(obj);

	io.WriteString(w,string(jsonData))
}

func ReadJson( r io.ReadCloser,obj interface{})  error {

	body, err := ioutil.ReadAll(io.LimitReader(r, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Close(); err != nil {
		panic(err)
	}
	mdz:=json.NewDecoder(bytes.NewBuffer(body))

	mdz.UseNumber()
	err = mdz.Decode(obj)

	if  err != nil {
		return err;
	}

	return nil;

	
}

func GenerUUId()  string{

	out, _ := exec.Command("uuidgen").Output()


	return strings.Replace(strings.TrimSpace(string(out)),"-","",-1)
}

//生成APPID
func GenerAppId() int64  {
	err, idWorker := idgen.NewIdWorker(1)
	CheckErr(err)
	err,appid := idWorker.NextId()
	CheckErr(err)
	return appid;
}

type ResultError struct {

	ErrCode int `json:"err_code"`
	ErrMsg string `json:"err_msg"`

}


func NewResultError(errCode int,errMsg string) *ResultError  {

	resultError := &ResultError{}
	resultError.ErrCode=errCode;
	resultError.ErrMsg=errMsg

	return  resultError
}