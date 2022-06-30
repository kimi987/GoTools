package db

import (
	"dgServer/conf"
	"dgServer/log"
	"errors"
	"runtime/debug"

	"dgServer/mysqldb"
)

//Module 数据库模块
type Module struct {
	*mysqldb.Dbs
}

//OnInit 数据库模块初始化
func (module *Module) OnInit() {
	module.Dbs = new(mysqldb.Dbs)
	module.Dbs.DbDetails = make(map[string]*mysqldb.Db_detail)
	for n, _ := range conf.Config.DBName {
		db := new(mysqldb.Db_detail)
		db.Name = conf.Config.DBName[n]
		db.ConnAddr = conf.Config.ConnAddr[n]
		db.MaxOpenConn = conf.Config.MaxOpenConn
		db.MaxIdleConn = conf.Config.MaxIdleConn
		module.Dbs.DbDetails[db.Name] = db
	}

}

//RPCDB 远程模块
type RPCDB struct {
	M *Module
}

//模块
var (
	DBModule = new(Module)
	DBRpc    = &RPCDB{DBModule}
)

//RequestDB 数据库请求结构
type RequestDB struct {
	DBName    string
	Query     string
	QueryType int //执行类型
	Args      []interface{}
}

//ResponseDB 数据库返回结构
type ResponseDB struct {
	Rows map[int]map[string]string
}

//CreateRequestDB 创建请求
func CreateRequestDB(Query string, Args []interface{}) *RequestDB {
	return &RequestDB{DBName: conf.Config.DBName[0], Query: Query, Args: Args}
}

func CreateRequestDBCRM(Query string, Args []interface{}) *RequestDB {
	return &RequestDB{DBName: conf.Config.DBName[1], Query: Query, Args: Args}
}

//CreateResponseDB 创建返回
func CreateResponseDB() *ResponseDB {
	return &ResponseDB{
		Rows: make(map[int]map[string]string),
	}
}

//OnRun 启动
func OnRun(closeSign chan bool) {
	DBModule.OnInit()
	DBModule.Run(closeSign)
	<-closeSign
	DBModule.OnDestroy()
}

//错误恢复
func doPanic() {
	if err := recover(); err != nil {
		debug.PrintStack()
		log.Error("%s", string(debug.Stack()))
	}
}

//QuerySelect 查询调用
func QuerySelect(req *RequestDB, res *ResponseDB) error {
	defer doPanic()
	dbDetail := DBRpc.M.Dbs.DbDetails[req.DBName]
	if dbDetail != nil {
		result, err := dbDetail.QuerySelect(req.Query, req.Args...)
		if err != nil {
			log.Error("query selece error = %v", err)
			return errors.New("query select error ")
		} else {
			res.Rows = result

			//blob测试代码
			//if strings.Contains(req.Query, "t_mail") {
			//	for row, value := range res.Rows {
			//		for k, v := range value {
			//			utils.Message("blob测试2:row[%d],col[%s]=%s", row, k, v)
			//			utils.Message("blob测试2:row[%d],col[%s].len=%d", row, k, len(v))
			//		}
			//	}
			//}
		}
	} else {
		log.Error("query select %v,  can not find DB %v", req.Query, req.DBName)
		return errors.New("Can not find db")
	}

	return nil
}

//QueryInsert 插入调用
func QueryInsert(req *RequestDB, res *int) error {
	defer doPanic()
	dbDetail := DBRpc.M.Dbs.DbDetails[req.DBName]
	if dbDetail != nil {
		result, err := dbDetail.QueryInsert(req.Query, req.Args...)

		if err != nil {
			return err
		}

		*res = result
	} else {
		log.Error("query insert %v,  can not find DB %v", req.Query, req.DBName)
		return errors.New("Can not find DB")
	}

	return nil
}

//QueryUpdate 更新调用
func QueryUpdate(req *RequestDB, res *int) error {
	defer doPanic()
	dbDetail := DBRpc.M.Dbs.DbDetails[req.DBName]
	if dbDetail != nil {
		result, err := dbDetail.QueryUpdate(req.Query, req.Args...)
		if err != nil {
			return err
		}

		*res = result
	} else {
		log.Error("query update %v,  can not find DB %v", req.Query, req.DBName)
		return errors.New("Can not find DB")
	}
	return nil
}

//QueryRemove 删除调用
func QueryRemove(req *RequestDB, res *int) error {
	defer doPanic()
	dbDetail := DBRpc.M.Dbs.DbDetails[req.DBName]
	if dbDetail != nil {
		result, err := dbDetail.QueryUpdate(req.Query, req.Args...)

		if err != nil {
			return err
		}
		*res = result
	} else {
		log.Error("query delete %v,  can not find DB %v", req.Query, req.DBName)
		return errors.New("Can not find DB")
	}

	return nil
}
