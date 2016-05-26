package pay

import (

	"fmt"
	"crypto/md5"
	"time"
	"sort"
	"strings"
	"reflect"
	"strconv"
)

func NewNonceString() string {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 36)

	nonceStr :=fmt.Sprintf("%x", md5.Sum([]byte(nonce)))
	fmt.Println("111notice=",nonceStr)
	return nonceStr
}

// SortAndConcat sort the map by key in ASCII order,
// and concat it in form of "k1=v1&k2=2"
func SortAndConcat(param map[string]string) string {
	var keys []string
	for k := range param {
		keys = append(keys, k)
	}

	var sortedParam []string
	sort.Strings(keys)
	for _, k := range keys {
		// fmt.Println(k, "=", param[k])
		sortedParam = append(sortedParam, k+"="+param[k])
	}

	return strings.Join(sortedParam, "&")
}

// Sign the parameter in form of map[string]string with app key.
// Empty string and "sign" key is excluded before sign.
// Please refer to http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=4_3
func Sign(param map[string]string, key string) string {
	newMap := make(map[string]string)
	// fmt.Printf("%#v\n", param)
	for k, v := range param {
		if k == "sign" {
			continue
		}
		if v == "" {
			continue
		}
		newMap[k] = v
	}
	// fmt.Printf("%#v\n\n", newMap)

	preSignStr := SortAndConcat(newMap)
	preSignWithKey := preSignStr + "&key=" + key
	//fmt.Println(preSignWithKey)
	//fmt.Println("preSignStr==",fmt.Sprintf("%X", md5.Sum([]byte(preSignWithKey))))
	return fmt.Sprintf("%X", md5.Sum([]byte(preSignWithKey)))
}

const ChinaTimeZoneOffset = 8 * 60 * 60 //Beijing(UTC+8:00)

// NewTimestampString return
func NewTimestampString() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}



// ToXmlString convert the map[string]string to xml string
func ToXmlString(param map[string]string) string {
	xml := "<xml>"
	for k, v := range param {
		xml = xml + fmt.Sprintf("<%s>%s</%s>", k, v, k)
	}
	xml = xml + "</xml>"

	return xml
}

// ToMap convert the xml struct to map[string]string
func ToMap(in interface{}) (map[string]string, error) {
	out := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get("xml"); tagv != "" && tagv != "xml" {
			// set key of map to value in struct field
			out[tagv] = v.Field(i).String()
		}
	}
	return out, nil
}

// ToMap convert the json struct to map[string]string
func ToMapOfJson(in interface{}) (map[string]string, error) {
	out := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get("json"); tagv != "" && tagv != "json" {
			// set key of map to value in struct field
			out[tagv] = v.Field(i).String()
		}
	}
	return out, nil
}

