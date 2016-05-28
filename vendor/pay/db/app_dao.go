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

	stmt,err :=GetDB().Prepare("insert into app(app_id,open_id,app_key,app_name,app_desc,status,create_time,update_time) values(?,?,?,?,?,?,?,?)")
	util.CheckErr(err)
	_,err =stmt.Exec(self.AppId,self.OpenId,self.AppKey,self.AppName,self.AppDesc,self.Status,self.CreateTime,self.UpdateTime)
	util.CheckErr(err)
	return err
}

//查询可用的APP
func (self *APP) QueryCanUseApp(appId string) *APP {

	stmt,err := GetDB().Prepare("select id,app_id,open_id,app_key,app_name,app_desc,status from app where app_id=? and status=1")
	util.CheckErr(err)

	rows,err := stmt.Query(appId);

	defer rows.Close()
	util.CheckErr(err)
	if rows.Next()  {
		app :=NewAPP()
		rows.Scan(&app.Id,&app.AppId,&app.OpenId,&app.AppKey,&app.AppName,&app.AppDesc,&app.Status)

		return app
	}

	return nil;
}
