package db

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
	"database/sql"
	"github.com/tangtaoit/util"
)



type Account struct {

	Id uint64
	AppId string
	OpenId string
	//账户余额
	Amount int64
	//账户状态 1.正常 0.异常
	Status int

	//创建时间
	CreateTime time.Time

	//创建时间
	UpdateTime time.Time

}

type AccountRecord struct  {

	Id uint64
	//交易号
	TradeNo string
	//应用ID
	AppId string
	//用户OPENID
	OpenId string
	//账户号
	AccountId uint64
	//创建时间
	CreateTime time.Time
	//交易前的账户金额
	AmountBefore int64
	//交易后的账户金额
	AmountAfter int64
	//交易产生的账户变动金额
	ChangedAmount int64

}

func NewAccountRecord()  *AccountRecord{

	return &AccountRecord{}
}

const ACCOUNT_INSERT_SQL  =  "insert into accounts(app_id,open_id,amount,create_time,status) values(?,?,?,?,?)"

func NewAccount() *Account {

	return &Account{}
}

func (self *Account) QueryAccount(openId,appId string)  *Account {

	rows,err := GetDB().Query("select id,app_id,open_id,amount,create_time,status from accounts where open_id=? and app_id=?",openId,appId)
	defer rows.Close()

	if err!=nil {
		util.CheckErr(err)
		return nil;
	}

	if rows.Next() {

		var id uint64
		var appId *string
		var openId *string
		var amount int64
		var createTime *time.Time
		var status int

		err :=rows.Scan(&id,&appId,&openId,&amount,&createTime,&status)
		if err!=nil {
			util.CheckErr(err)
			return nil;
		}

		account :=NewAccount()
		account.Id=id;
		account.CreateTime=*createTime;
		account.Status=status;
		account.Amount=amount
		account.OpenId=*openId;
		account.AppId=*appId

		return account;
	}

	return nil;
}

func (account *Account) getInsertParam() [5]interface{}  {


	return [5]interface{}{account.AppId,account.OpenId,account.Amount,account.CreateTime,account.Status}
}

//添加账户
func (account *Account) Insert() bool {

	stmt,err := GetDB().Prepare(ACCOUNT_INSERT_SQL)

	defer stmt.Close()

	_,err =stmt.Exec(account.getInsertParam())
	if err!=nil {

		util.CheckErr(err)
		return false;
	}

	return true;
}

func (self *Account) InsertTx(tx  *sql.Tx)  {
	stmt,err := tx.Prepare(ACCOUNT_INSERT_SQL)
	util.CheckErr(err)
	_,err =stmt.Exec(self.getInsertParam())
	util.CheckErr(err)
}

