package db

import "pay/comm"

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

}

func NewAPP() *APP  {

	return &APP{}
}

func (self *APP)  Insert() error{

	stmt,err :=db.Prepare("insert into app(app_id,open_id,app_key,app_name,app_desc,status) values(?,?,?,?,?,?)")
	comm.CheckErr(err)
	_,err =stmt.Exec(self.AppId,self.OpenId,self.AppKey,self.AppName,self.AppDesc,self.Status)
	comm.CheckErr(err)
	return err
}

//查询可用的APP
func (self *APP) QueryCanUseApp(appId string) *APP {

	stmt,err := db.Prepare("select id,app_id,open_id,app_key,app_name,app_desc,status from app where app_id=? and status=1")
	comm.CheckErr(err)

	rows,err := stmt.Query(appId);

	defer rows.Close()
	comm.CheckErr(err)
	if rows.Next()  {
		app :=NewAPP()
		rows.Scan(&app.Id,&app.AppId,&app.OpenId,&app.AppKey,&app.AppName,&app.AppDesc,&app.Status)

		return app
	}

	return nil;
}
