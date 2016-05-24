package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"pay/config"
	"pay/comm"
	"time"
	"net/url"
)

var db *sql.DB

func init() {

	setting :=config.GetSetting()
	connInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=%s&parseTime=true",setting.MysqlUser,setting.MysqlPassword,setting.MysqlHost,setting.MysqlDB,url.QueryEscape("Asia/shanghai"))
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

type Account struct {

	id int64
	OpenId string
	//账户余额
	Amount float64
	//账户状态 1.正常 0.异常
	Status int

	//创建时间
	CreateTime time.Time

}

func NewAccount() *Account {

	return &Account{}
}

func QueryAccount(openId string)  *Account {

	rows,err := db.Query("select id,open_id,amount,create_time,status from accounts where open_id=?",openId)
	defer rows.Close()

	if err!=nil {
		comm.CheckErr(err)
		return nil;
	}

	if rows.Next() {

		var id int64
		var openId *string
		var amount float64
		var createTime *time.Time
		var status int

		err :=rows.Scan(&id,&openId,&amount,&createTime,&status)
		if err!=nil {
			comm.CheckErr(err)
			return nil;
		}

		account :=NewAccount()
		account.id=id;
		account.CreateTime=*createTime;
		account.Status=status;
		account.Amount=amount
		account.OpenId=*openId;

		return account;
	}

	return nil;
}

//添加账户
func InsertAccount(account *Account) bool {

	_,err := db.Exec("insert into accounts(open_id,amount,create_time,status) values(?,?,?,?)",account.OpenId,account.Amount,account.CreateTime,account.Status)

	if err!=nil {

		comm.CheckErr(err)
		return false;
	}

	return true;
}

