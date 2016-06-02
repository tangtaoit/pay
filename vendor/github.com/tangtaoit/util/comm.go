package util

import (
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"bytes"
	"github.com/sumory/idgen"
	"hash"
	"sort"
	"encoding/hex"
	"crypto/md5"
	"bufio"
	"fmt"
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

	w.WriteHeader(statusCode)
	WriteJson(w,err)
}

func ResponseSuccess(w http.ResponseWriter)  {

	err := NewResultError(0,"OK")
	WriteJson(w,err)

}

func WriteJson(w io.Writer,obj interface{})  {

	jsonData,_:= json.Marshal(obj);

	io.WriteString(w,string(jsonData))
}

func ReadJsonByByte(body []byte,obj interface{}) error {

	mdz:=json.NewDecoder(bytes.NewBuffer(body))

	mdz.UseNumber()
	err := mdz.Decode(obj)

	if  err != nil {
		return err;
	}
	return nil;
}

func ReadJson( r io.ReadCloser,obj interface{})  error {

	body, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	if err := r.Close(); err != nil {
		panic(err)
	}


	return ReadJsonByByte(body,obj);

	
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

func SignWithBaseSign(params map[string]interface{}, apiKey string,basesign string, fn func() hash.Hash) string {
	if fn == nil {
		fn = md5.New
	}
	h := fn()
	bufw := bufio.NewWriterSize(h, 128)

	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := params[k]
		if v == "" {
			continue
		}
		vs := ObjToStr(v)
		bufw.WriteString(k)
		bufw.WriteByte('=')
		bufw.WriteString(vs)
		bufw.WriteByte('&')
	}
	bufw.WriteString("key=")
	bufw.WriteString(apiKey)

	if basesign!=""{
		bufw.WriteString("&")
		bufw.WriteString("basesign=")
		bufw.WriteString(basesign)
	}


	bufw.Flush()

	signature := make([]byte, hex.EncodedLen(h.Size()))
	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}

// Sign 支付签名.
//  params: 待签名的参数集合
//  apiKey: api密钥
//   basesign 基础sign
//  fn:     func() hash.Hash, 如果为 nil 则默认用 md5.New
func Sign(params map[string]string, apiKey string, fn func() hash.Hash) string {

	objparams :=make(map[string]interface{})

	for key,_ :=range params {

		objparams[key] = params[key]
	}

	return SignWithBaseSign(objparams,apiKey,"",fn)
}

func ObjToStr(v interface{}) string {
	var strV string
	switch v.(type) {

		case int:
			strV= fmt.Sprintf("%d",v)
			break
		case uint:
			strV= fmt.Sprintf("%d",v)
			break
		case int64:
			strV= fmt.Sprintf("%d",v)
			break
		case uint64:
			strV= fmt.Sprintf("%d",v)
			break
		case int8:
			strV= fmt.Sprintf("%d",v)
			break
		case uint8:
			strV= fmt.Sprintf("%d",v)
			break
		case int16:
			strV= fmt.Sprintf("%d",v)
			break
		case uint16:
			strV= fmt.Sprintf("%d",v)
			break
		case int32:
			strV= fmt.Sprintf("%d",v)
			break
		case uint32:
			strV= fmt.Sprintf("%s",v)
			break
		case string:
			strV= fmt.Sprintf("%s",v)
			break
		case float32:
			strV= fmt.Sprintf("%s",v)
			break
		case float64:
			strV= fmt.Sprintf("%s",v)
			break
		default:
			strV= fmt.Sprintf("%s",v)

	}
	return strV
}
