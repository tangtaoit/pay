package db

import (
	"pay/config"
	"fmt"
	"net/url"
	"database/sql"
	"pay/comm"
	"time"
)
var db *sql.DB
func init() {

	loc,_ := time.LoadLocation("Local")

	setting :=config.GetSetting()
	connInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=%s&parseTime=true",setting.MysqlUser,setting.MysqlPassword,setting.MysqlHost,setting.MysqlDB,loc.String())
	fmt.Println(connInfo);
	var err error;
	db, err = sql.Open("mysql",connInfo)
	if err!=nil{
		comm.CheckErr(err)
	}

	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}
