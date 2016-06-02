package db

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
	"github.com/gocraft/dbr"
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

	Password string

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


func NewAccount() *Account {

	return &Account{}
}

func (self *Account) QueryAccount(openId,appId string)  (*Account,error) {

	session := NewSession()
	var account *Account
	_,err := session.Select("*").From("accounts").Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&account)
	return account,err;
}



//添加账户
func (account *Account) Insert() error {

	sess := NewSession()
	_,err := sess.InsertInto("accounts").Columns("app_id","open_id","amount","create_time","update_time","status","password").Record(account).Exec()

	return err;
}

func (self *Account) InsertTx(tx  *dbr.Tx) (sql.Result,error) {

	return tx.InsertInto("accounts").Columns("app_id","open_id","amount","create_time","status","password").Record(self).Exec()
}

func (self *Account) UpdatePwd(openId string,password string) error{

	session :=NewSession()

	_,err :=session.Update("accounts").Set("password",password).Set("update_time",time.Now()).Where("open_id=?",openId).Exec()

	return err

}

func (self *AccountRecord) InsertTx(tx *dbr.Tx) error  {

	_,err := tx.InsertInto("accounts_record").Columns("trade_no","app_id","open_id","account_id","amount_before","amount_after","changed_amount").Record(self).Exec()

	return err
}

//账户金额改变
func AccountAmountChange(changeAmount int64,tradeNo string,openId string,appId string,tx *dbr.Tx) error {

	//修改账户余额
	//如果用户没有创建账户,那么就创建一个新的账户
	//添加账户变动记录
	var account *Account
	err:=tx.Select("*").From("accounts").Where("open_id=?",openId).LoadStruct(&account)
	if err!=nil{
		return err
	}

	var amountBefrore int64
	var amountAfter int64
	if account==nil{
		amountBefrore=0
		amountAfter =changeAmount
		account = NewAccount()
		account.OpenId=openId
		account.AppId=appId
		account.CreateTime=time.Now()
		account.UpdateTime=time.Now()
		account.Status=1
		account.Amount=changeAmount
		result,err := account.InsertTx(tx)
		util.CheckErr(err)
		lastId,_ := result.LastInsertId()
		account.Id=uint64(lastId)
	}else{
		amountBefrore=account.Amount
		amountAfter = account.Amount+changeAmount

		_,err :=tx.Update("accounts").Where("id=?",account.Id).Set("amount",account.Amount+changeAmount).Exec()
		util.CheckErr(err)

	}
	accRecod :=NewAccountRecord()
	accRecod.TradeNo=tradeNo
	accRecod.OpenId=account.OpenId
	accRecod.AccountId=account.Id
	accRecod.AmountBefore=amountBefrore
	accRecod.AmountAfter=amountAfter
	accRecod.ChangedAmount=changeAmount
	accRecod.AppId=account.AppId
	accRecod.CreateTime=time.Now()
	err = accRecod.InsertTx(tx)
	if err!=nil{
		return err
	}

	return nil
}


