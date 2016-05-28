package db

import (
	"pay/config"
	"fmt"
	"database/sql"
	"time"
	"github.com/tangtaoit/util"
	"github.com/gocraft/dbr"
)
var conn *dbr.Connection
var sqldb  *sql.DB

func init()  {

}

func GetDB()  *sql.DB {

	if sqldb==nil{
		InitMysql()
	}
	return sqldb;
}

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

	sqldb =conn.DB

	fmt.Println("mysql inital is success");
}

func NewSession() *dbr.Session {
	if conn==nil {
		InitMysql()
	}


	return conn.NewSession(nil)
}

func Begin() *sql.Tx {

	tx,err :=sqldb.Begin()

	util.CheckErr(err)

	return tx
}

func Exec(dbTx interface{},sqlStr string,param ...interface{})  {



	var stmt *sql.Stmt
	var err error
	if tx,isOk := dbTx.(*sql.Tx);isOk {


		stmt,err = tx.Prepare(sqlStr)

	}else if db,isOk :=dbTx.(*sql.DB);isOk {
		stmt,err = db.Prepare(sqlStr)
	}
	util.CheckErr(err)
	defer stmt.Close()
	_,err = stmt.Exec(param)
	util.CheckErr(err)
}