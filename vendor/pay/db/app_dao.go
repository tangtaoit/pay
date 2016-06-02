package db

import (
	"time"
	"github.com/tangtaoit/util"
)

type APP struct  {
	Id uint64
	//应用ID
	AppId string
	//应用KEY
	AppKey string
	//应用名称
	AppName string
	//应用描述
	AppDesc string
	//应用状态 0.待审核 1.已审核
	Status int
	//openID
	OpenId string

	//创建时间
	CreateTime time.Time

	//修改时间
	UpdateTime time.Time

}

func NewAPP() *APP  {

	return &APP{}
}

func (self *APP)  Insert() error{

	//nw :=time.Now()
	self.CreateTime = time.Now()
	self.UpdateTime= time.Now()

	session := NewSession()
	_,err :=session.InsertInto("app").Columns("app_id","open_id","app_key","app_name","app_desc","status","create_time","update_time").Record(self).Exec()
	util.CheckErr(err)

	return err
}

//查询可用的APP
func (self *APP) QueryCanUseApp(appId string) (*APP,error) {

	sess := NewSession()
	var app *APP
	_,err :=sess.Select("*").From("app").Where("app_id=?",appId).Where("status=?",1).LoadStructs(&app)


	return app,err;
}
