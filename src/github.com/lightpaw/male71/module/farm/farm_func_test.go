package farm

import (
	"database/sql"
	"fmt"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/db"
	"github.com/lightpaw/male7/service/db/sqldb"
	"github.com/lightpaw/male7/service/timeservice"
	. "github.com/onsi/gomega"
	"path"
	"strings"
	"testing"
	"time"
)

func TestFarmFunc(t *testing.T) {
	RegisterTestingT(t)

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	Ω(err).Should(Succeed())

	farmFunc := &FarmFunc{}
	farmFunc.datas = datas
	farmFunc.miscConf = datas.FarmMiscConfig()
	farmFunc.worldService = ifacemock.WorldService
	farmFunc.timeService = ifacemock.TimeService
	farmFunc.db = newMysqlAdapter(datas)
	ifacemock.TimeService.Mock(ifacemock.TimeService.CurrentTime, func() time.Time {
		return time.Now()
	})

	if farmFunc.db != nil {
		getFarmTest(farmFunc)
		plantTest(farmFunc)
		harvestTest(farmFunc)
		//oneKeyPlantTest(farmFunc)
	}
}

func getFarmTest(farmFunc *FarmFunc) {
	allCubes := make([]cb.Cube, 0)
	for i := 1; i <= 10; i++ {
		allCubes = append(allCubes, cb.XYCube(i, i))
	}

	errMsg, farmCubes := farmFunc.GetFarm(1, allCubes)
	Ω(errMsg).Should(BeNil())
	Ω(len(farmCubes)).Should(Equal(10))
}

func plantTest(farmFunc *FarmFunc) {
	heroId := int64(0)
	resId := uint64(1)
	cube := cb.XYCube(-10, 20)
	resConf := farmFunc.datas.GetFarmResConfig(resId)
	errMsg := farmFunc.Plant(ifacemock.HeroController, cube, resConf)
	Ω(errMsg).Should(BeNil())
	succ, fc := farmFunc.loadFarmCube(heroId, cube)
	Ω(succ).Should(BeTrue())
	Ω(fc.ResId).Should(Equal(resId))
}

func oneKeyPlantTest(farmFunc *FarmFunc) {
	heroId := int64(0)
	resId := uint64(1)
	cube := cb.XYCube(-10, 20)
	allCubes := []cb.Cube{cb.XYCube(-10, 20), cb.XYCube(10, -20), cb.XYCube(-20, 10), cb.XYCube(20, -10)}
	okConf := farmFunc.datas.GetFarmOneKeyConfig(10)
	errMsg := farmFunc.OneKeyPlant(ifacemock.HeroController, okConf, []cb.Cube{}, allCubes, farmFunc.datas.FarmMiscConfig().OneKeyResConfig[shared_proto.ResType_GOLD], farmFunc.datas.FarmMiscConfig().OneKeyResConfig[shared_proto.ResType_STONE], 10, 10)
	Ω(errMsg).Should(BeNil())
	succ, fc := farmFunc.loadFarmCube(heroId, cube)
	Ω(succ).Should(BeTrue())
	Ω(fc.ResId).Should(Equal(resId))
}

func harvestTest(farmFunc *FarmFunc) {
	resId := uint64(1)
	cube := cb.XYCube(-30, 50)
	resConf := farmFunc.datas.GetFarmResConfig(resId)
	errMsg := farmFunc.Plant(ifacemock.HeroController, cube, resConf)
	Ω(errMsg).Should(BeNil())

	errMsg, fc := farmFunc.Harvest(ifacemock.HeroController, cube)
	Ω(errMsg).Should(BeNil())
	Ω(fc.StealTimes).Should(Equal(uint64(0)))
}

func newMysqlAdapter(datas *config.ConfigDatas) iface.DbService {

	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/male7test"

	db0, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return nil
	}

	err = db0.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return nil
	}

	dbName := path.Base(dataSourceName)
	if idx := strings.Index(dbName, "?"); idx >= 0 {
		dbName = dbName[:idx]
	}

	ds, err := sqldb.NewSqlDbService(db0, dbName, datas, timeservice.NewDefaultTimeService(), ifacemock.MetricsRegister)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return nil
	}

	truncateSqlDb(db0)

	return db.NewMysqlDbService(ds)
}

func truncateSqlDb(db *sql.DB) {

	truncate := func(tableName string) {
		_, err := db.Exec("truncate table " + tableName)
		Ω(err).Should(Succeed())
	}
	truncate("farm")
}
