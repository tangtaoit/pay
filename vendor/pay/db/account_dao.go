package db

import (
	_ "github.com/go-sql-driver/mysql"
	"pay/comm"
	"time"
)



type Account struct {

	id uint64
	AppId string
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

func QueryAccount(openId string,appId string)  *Account {

	rows,err := db.Query("select id,app_id,open_id,amount,create_time,status from accounts where open_id=? and app_id=?",openId,appId)
	defer rows.Close()

	if err!=nil {
		comm.CheckErr(err)
		return nil;
	}

	if rows.Next() {

		var id uint64
		var appId *string
		var openId *string
		var amount float64
		var createTime *time.Time
		var status int

		err :=rows.Scan(&id,&appId,&openId,&amount,&createTime,&status)
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
		account.AppId=*appId

		return account;
	}

	return nil;
}

//添加账户
func InsertAccount(account *Account) bool {

	stmt,err := db.Prepare("insert into accounts(app_id,open_id,amount,create_time,status) values(?,?,?,?,?)")

	_,err =stmt.Exec(account.AppId,account.OpenId,account.Amount,account.CreateTime,account.Status)
	if err!=nil {

		comm.CheckErr(err)
		return false;
	}

	return true;
}


