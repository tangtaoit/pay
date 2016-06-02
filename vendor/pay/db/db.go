package db

import (
	"pay/config"
	"fmt"
	"time"
	"github.com/tangtaoit/util"
	"github.com/gocraft/dbr"
)
var conn *dbr.Connection


func InitMysql() {

	fmt.Println("init mysql...");
	loc,_ := time.LoadLocation("Local")

	setting :=config.GetSetting()
	connInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=%s&parseTime=true",setting.MysqlUser,setting.MysqlPassword,setting.MysqlHost,setting.MysqlDB,loc.String())
	fmt.Println(connInfo);
	var err error;
	conn,err = dbr.Open("mysql",connInfo,nil)
	util.CheckErr(err)
	conn.SetMaxOpenConns(2000)
	conn.SetMaxIdleConns(1000)
	conn.Ping()

	fmt.Println("mysql inital is success");
}

func NewSession() *dbr.Session {
	if conn==nil {
		InitMysql()
	}


	return conn.NewSession(nil)
}
