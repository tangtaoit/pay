package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/tangtaoit/util"
)

var environments = map[string]string{
	"production":    "config/prod.json",
	"preproduction": "config/pre.json",
	"tests":         "config/tests.json",
}

type Settings struct {
	PublicKeyPath      string
	JWTExpirationDelta int
	MysqlHost	   string
	MysqlPassword 	   string
	MysqlUser	   string
	MysqlDB  	   string
	//token失效时间
        TokenExpire        float32

	RedisAddress string

}


var settings *Settings
var env = "preproduction"

func Init() {
	env = os.Getenv("GO_ENV")
	pwd, _ := os.Getwd()
	fmt.Println(pwd)
	if env == "" {
		fmt.Println("Warning: Setting preproduction environment due to lack of GO_ENV value")
		env = "preproduction"
	}
	LoadSettingsByEnv(env)
}

func LoadSettingsByEnv(env string) {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)

		util.CheckErr(err)
	}
	settings = &Settings{}
	jsonErr := json.Unmarshal(content, &settings)
	util.CheckErr(jsonErr)
}


func GetSetting() *Settings {
	if settings == nil {
		fmt.Println("------")
		Init()
	}
	return settings
}

func IsTestEnvironment() bool {
	return env == "tests"
}
