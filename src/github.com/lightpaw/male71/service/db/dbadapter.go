package db

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/service/db/sqldb"
	"github.com/lightpaw/male7/service/timeservice"
)

var _ DbServiceAdapter = (*DbService)(nil)
var _ DbServiceAdapter = (*sqldb.SqlDbService)(nil)
//var _ DbServiceAdapter = (*datastoredb.DatastoreDbService)(nil)

func NewDbService(datas *config.ConfigDatas, serverConfig *kv.IndividualServerConfig, timeService *timeservice.TimeService, register iface.MetricsRegister) *DbService {

	//var adapter DbServiceAdapter
	//var err error
	//if !serverConfig.EnableMysql {
	//	adapter, err = datastoredb.NewDatastoreDbService(datas, serverConfig, timeService)
	//	if err != nil {
	//		logrus.WithField("datastore_host", serverConfig.DatastoreHost).
	//			WithField("projectID", serverConfig.ProjectID).
	//			WithField("serverID", serverConfig.ServerID).
	//			WithError(err).Panic("初始化datastores service失败")
	//	}
	//} else {
	//	adapter, err = sqldb.NewMysqlDbService(serverConfig.DBSN, datas, timeService, register)
	//	if err != nil {
	//		logrus.WithField("dbsn", serverConfig.DBSN).WithError(err).Panic("初始化db service失败")
	//	}
	//}

	adapter, err := sqldb.NewMysqlDbService(serverConfig, datas, timeService, register)
	if err != nil {
		logrus.WithField("dbsn", serverConfig.DBSN).WithError(err).Panic("初始化db service失败")
	}

	return NewMysqlDbService(adapter)
}

func NewMysqlDbService(adapter *sqldb.SqlDbService) *DbService {
	return &DbService{
		adapter: adapter,
	}
}