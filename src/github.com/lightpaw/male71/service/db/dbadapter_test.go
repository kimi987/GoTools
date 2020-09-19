package db

import (
	//"cloud.google.com/go/datastore"
	"context"
	"database/sql"
	"fmt"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	//"github.com/lightpaw/male7/service/db/nsds"
	"github.com/lightpaw/male7/service/db/sqldb"
	"github.com/lightpaw/male7/service/timeservice"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/must"
	. "github.com/onsi/gomega"
	//"google.golang.org/api/iterator"
	"math"
	"path"
	"testing"
	"time"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/util/timeutil"
	"strings"
)

//func newDatastoreAdapter(datas *config.ConfigDatas) DbServiceAdapter {
//
//	client, err := nsds.NewNamespaceClient("127.0.0.1:8431", "male7test", "test")
//	if err != nil {
//		fmt.Println("创建nsds.Client失败，跳过测试", err.Error())
//		return nil
//	}
//
//	truncateDatastore(client)
//
//	db, err := datastoredb.NewDatastoreDbServiceWithNsdsClient(datas, timeservice.NewDefaultTimeService(), client)
//	if err != nil {
//		fmt.Println("创建datastoredb.Client失败，跳过测试", err.Error())
//		return nil
//	}
//	return db
//}

func TestTruncate(t *testing.T) {
	RegisterTestingT(t)

	//truncateDB(5)
	//truncateDB(6)
	//truncateDB(40)
	//truncateDB(99)

}

//func truncateDB(sid int) {
//	client, err := nsds.NewNamespaceClient("192.168.1.5:8432", "male7", fmt.Sprintf("1-%d", sid))
//	if err != nil {
//		fmt.Println("创建nsds.Client失败，跳过测试", err.Error())
//		return
//	}
//
//	//truncateGuild(client)
//
//	truncateTable(client, "mail")
//}
//
//func truncateGuild(client *nsds.NamespaceClient) {
//
//	truncateTable(client, "guild")
//
//	// 遍历所有英雄数据，把guild字段设置为0
//	q := datastore.NewQuery("hero")
//
//	for t := client.Run(context.TODO(), q); ; {
//		var x hero_entity
//		key, err := t.Next(&x)
//		if err != nil {
//			if err == iterator.Done {
//				break
//			}
//
//			fmt.Println("db错误", err.Error())
//			return
//		}
//
//		heroProto := &server_proto.HeroServerProto{}
//		heroProto.Unmarshal(x.HeroData)
//		if heroProto.GuildId != 0 {
//			fmt.Println("清除玩家联盟", heroProto.Id, heroProto.Name)
//			heroProto.GuildId = 0
//			x.HeroData, _ = heroProto.Marshal()
//
//			client.Put(context.TODO(), key, &x)
//		}
//	}
//}

type hero_entity struct {
	Name       string `datastore:",noindex"`
	HeroData   []byte `datastore:",noindex"`
	BaseRegion int64
}

//func truncateTable(client *nsds.NamespaceClient, kind string) {
//	keys, err := client.GetAll(context.TODO(), datastore.NewQuery(kind).KeysOnly(), nil)
//	Ω(err).Should(Succeed())
//
//	if len(keys) > 0 {
//		Ω(client.DeleteMulti(context.TODO(), keys)).Should(Succeed())
//
//		Eventually(func() int {
//			count, err := client.Count(context.TODO(), datastore.NewQuery(kind).KeysOnly())
//			Ω(err).Should(Succeed())
//
//			return count
//		}, "10s").Should(Equal(0))
//	}
//}
//
//func truncateDatastore(client *nsds.NamespaceClient) {
//
//	truncate := func(kind string) {
//		truncateTable(client, kind)
//	}
//
//	truncate("kv")
//	truncate("user")
//	truncate("heroname")
//	truncate("hero")
//	truncate("guild")
//	truncate("mail")
//	truncate("baizhan")
//	truncate("farm")
//}

func newMysqlAdapter(datas *config.ConfigDatas, dataSourceName string, truncate bool) DbServiceAdapter {

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return nil
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return nil
	}

	dbName := path.Base(dataSourceName)
	if idx := strings.Index(dbName, "?"); idx >= 0 {
		dbName = dbName[:idx]
	}

	ds, err := sqldb.NewSqlDbService(db, dbName, datas, timeservice.NewDefaultTimeService(), ifacemock.MetricsRegister)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return nil
	}

	if truncate {
		truncateSqlDb(db)
	}

	return ds
}

func truncateSqlDb(db *sql.DB) {

	truncate := func(tableName string) {
		_, err := db.Exec("truncate table " + tableName)
		Ω(err).Should(Succeed())
	}

	//truncate("kv")
	truncate("user")
	truncate("hero")
	truncate("guild")
	truncate("guild_logs")
	truncate("mail")
	truncate("bai_zhan")
	truncate("chat_msg")
	truncate("chat_window")
	truncate("farm")
}

//func TestDatastoreAdapter(t *testing.T) {
//	RegisterTestingT(t)
//
//	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
//	Ω(err).Should(Succeed())
//
//	if db := newDatastoreAdapter(datas); db != nil {
//		testAdapter(db, datas)
//	}
//}

func TestMysqlAdapter(t *testing.T) {
	RegisterTestingT(t)

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	Ω(err).Should(Succeed())

	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/male7test"
	if db := newMysqlAdapter(datas, dataSourceName, true); db != nil {
		testAdapter(db, datas)
	}
}

func testAdapter(db DbServiceAdapter, datas *config.ConfigDatas) {

	testKV(db)

	testUser(db)

	testHero(db, datas)

	testGuild(db, datas)

	testGuildLogs(db)

	testMail(db)

	testBaizhan(db)

	testChat(db)

	testFarm(db)

	textFarmLog(db)

	testFarmSteal(db)
}

func testKV(db DbServiceAdapter) {

	key := server_proto.Key_UNIT_TEST

	// 保存一个数据
	toSet := []byte("test")
	err := db.SaveKey(context.TODO(), key, toSet)
	Ω(err).Should(Succeed())

	newData, err := db.LoadKey(context.TODO(), key)
	Ω(err).Should(Succeed())

	Ω(newData).Should(Equal(toSet))

}

const heroId int64 = 9999
const heroName string = "hero9999"

func testUser(db DbServiceAdapter) {

	proto := &server_proto.UserMiscProto{
		TutorialProgress:   100,
		IsTutorialComplete: true,
	}

	err := db.UpdateUserMisc(context.TODO(), heroId, proto)
	Ω(err).Should(Succeed())

	newProto, err := db.LoadUserMisc(context.TODO(), heroId)
	Ω(err).Should(Succeed())

	Ω(*newProto).Should(Equal(*proto))
}

func testHero(db DbServiceAdapter, datas *config.ConfigDatas) {

	newHeroId := heroId

	// 创建之前
	heroId, err := db.HeroId(context.TODO(), heroName)
	Ω(err).Should(Succeed())
	Ω(heroId).Should(Equal(int64(0)))

	exist, err := db.HeroNameExist(context.TODO(), heroName)
	Ω(err).Should(Succeed())
	Ω(exist).Should(BeFalse())

	hero, err := db.LoadHero(context.TODO(), newHeroId)
	Ω(err).Should(Succeed())
	Ω(hero).Should(BeNil())

	hero = entity.NewHero(newHeroId, heroName, datas.HeroInitData(), time.Now())

	// 第一次创建，成功
	result, err := db.CreateHero(context.TODO(), hero)
	Ω(err).Should(Succeed())
	Ω(result).Should(BeTrue())

	// 第二次创建，失败
	result, err = db.CreateHero(context.TODO(), hero)
	Ω(err).Should(HaveOccurred())
	Ω(result).Should(BeFalse())

	// 换个id，名字一样，失败
	hero = entity.NewHero(newHeroId+1, heroName, datas.HeroInitData(), time.Now())
	result, err = db.CreateHero(context.TODO(), hero)
	Ω(err).Should(HaveOccurred())
	Ω(result).Should(BeFalse())

	// 创建之之后
	heroId, err = db.HeroId(context.TODO(), heroName)
	Ω(err).Should(Succeed())
	Ω(heroId).Should(Equal(newHeroId))

	exist, err = db.HeroNameExist(context.TODO(), heroName)
	Ω(err).Should(Succeed())
	Ω(exist).Should(BeTrue())

	hero, err = db.LoadHero(context.TODO(), newHeroId)
	Ω(err).Should(Succeed())
	Ω(hero).ShouldNot(BeNil())
	Ω(hero.Id()).Should(Equal(newHeroId))
	Ω(hero.Name()).Should(Equal(heroName))

	if _, ok := db.(*sqldb.SqlDbService); ok {
		// 模糊查询
		text := heroName[:len(heroName)-2]
		heros, err := db.LoadNoGuildHeroListByName(context.TODO(), text, 0, 10)
		Ω(err).Should(Succeed())
		Ω(len(heros)).Should(Equal(1))
		Ω(heros[0].GuildId()).Should(Equal(int64(0)))
		Ω(heros[0].Name()).Should(Equal(heroName))
	}

	// 改名
	newName := "SB250"
	Ω(db.UpdateHeroName(context.TODO(), heroId, heroName, newName)).Should(BeTrue())

	heroId, err = db.HeroId(context.TODO(), heroName)
	Ω(err).Should(Succeed())
	Ω(heroId).Should(Equal(int64(0)))

	exist, err = db.HeroNameExist(context.TODO(), heroName)
	Ω(err).Should(Succeed())
	Ω(exist).Should(BeFalse())

	heroId, err = db.HeroId(context.TODO(), newName)
	Ω(err).Should(Succeed())
	Ω(heroId).Should(Equal(newHeroId))

	exist, err = db.HeroNameExist(context.TODO(), newName)
	Ω(err).Should(Succeed())
	Ω(exist).Should(BeTrue())

	// 没有主城
	regionHeros, err := db.LoadAllRegionHero(context.TODO())
	Ω(err).Should(Succeed())
	Ω(regionHeros).Should(BeEmpty())

	Eventually(func() int {
		heros, err := db.LoadAllHeroData(context.TODO())
		Ω(err).Should(Succeed())
		return len(heros)
	}, "10s").Should(Equal(1))

	// 添加主城，保存
	hero.SetBaseXY(1, 1, 1)
	db.SaveHero(context.TODO(), hero)

	Eventually(func() int {
		regionHeros, err := db.LoadAllRegionHero(context.TODO())
		Ω(err).Should(Succeed())
		return len(regionHeros)
	}, "10s").Should(Equal(1))

}

func newGuild(id int64, datas *config.ConfigDatas) *sharedguilddata.Guild {
	return sharedguilddata.NewGuild(id,
		fmt.Sprintf("g%d", id),
		fmt.Sprintf("f%d", id),
		datas,
		time.Now(),
	)
}

func testGuild(db DbServiceAdapter, datas *config.ConfigDatas) {

	idGen, err := db.MaxGuildId(context.TODO())
	Ω(err).Should(Succeed())
	Ω(idGen).Should(Equal(int64(0)))

	// load 数据不存在
	g, err := db.LoadGuild(context.TODO(), 1)
	Ω(err).Should(Succeed())
	Ω(g).Should(BeNil())

	// 插入数据
	err = db.CreateGuild(context.TODO(), 1, must.Marshal(newGuild(1, datas).EncodeServer()))
	Ω(err).Should(Succeed())

	err = db.CreateGuild(context.TODO(), 2, must.Marshal(newGuild(2, datas).EncodeServer()))
	Ω(err).Should(Succeed())

	err = db.CreateGuild(context.TODO(), 3, must.Marshal(newGuild(3, datas).EncodeServer()))
	Ω(err).Should(Succeed())

	// load 存在的数据
	g, err = db.LoadGuild(context.TODO(), 2)
	Ω(err).Should(Succeed())
	Ω(g).ShouldNot(BeNil())
	Ω(g.Id()).Should(Equal(int64(2)))

	// load全部数据
	gs, err := db.LoadAllGuild(context.TODO())
	Ω(err).Should(Succeed())
	Ω(gs).Should(HaveLen(3))

	var ids []int64
	for _, g := range gs {
		ids = append(ids, g.Id())

		Ω(g.GetText()).Should(Equal(""))

		// 改数据，并保存
		g.SetText(g.Name())
		err := db.SaveGuild(context.TODO(), g.Id(), must.Marshal(g.EncodeServer()))
		Ω(err).Should(Succeed())

		// 再load一次，看看数据对不对
		g, err := db.LoadGuild(context.TODO(), g.Id())
		Ω(err).Should(Succeed())
		Ω(g.GetText()).Should(Equal(g.Name()))
	}

	Ω(ids).Should(ConsistOf(int64(1), int64(2), int64(3)))

	idGen, err = db.MaxGuildId(context.TODO())
	Ω(err).Should(Succeed())
	Ω(idGen).Should(Equal(int64(3)))

	for _, id := range ids {
		err := db.DeleteGuild(context.TODO(), id)
		Ω(err).Should(Succeed())
	}
}

func testGuildLogs(db DbServiceAdapter) {

	// load 数据不存在
	logs, err := db.LoadGuildLogs(context.TODO(), 1, shared_proto.GuildLogType_GLTFight, 0, 10)
	Ω(err).Should(Succeed())
	Ω(logs).Should(BeEmpty())

	proto := &shared_proto.GuildLogProto{
		Type:   shared_proto.GuildLogType_GLTFight,
		FightX: 10,
		FightY: 20,
	}

	oldId := proto.Id

	// 插入数据
	err = db.InsertGuildLog(context.TODO(), 1, proto)
	Ω(err).Should(Succeed())
	Ω(proto.Id > oldId).Should(BeTrue())
	oldId = proto.Id

	err = db.InsertGuildLog(context.TODO(), 1, proto)
	Ω(err).Should(Succeed())
	Ω(proto.Id > oldId).Should(BeTrue())
	oldId = proto.Id

	proto.Type = shared_proto.GuildLogType_GLT_Memorabilia
	err = db.InsertGuildLog(context.TODO(), 1, proto)
	Ω(err).Should(Succeed())
	Ω(proto.Id > oldId).Should(BeTrue())
	oldId = proto.Id

	proto.Type = shared_proto.GuildLogType_GLTFight
	err = db.InsertGuildLog(context.TODO(), 2, proto)
	Ω(err).Should(Succeed())
	Ω(proto.Id > oldId).Should(BeTrue())
	oldId = proto.Id

	logs, err = db.LoadGuildLogs(context.TODO(), 1, shared_proto.GuildLogType_GLTFight, 0, 10)
	Ω(err).Should(Succeed())
	Ω(logs).Should(HaveLen(2))

	for _, v := range logs {
		Ω(v.FightX).Should(Equal(int32(10)))
		Ω(v.FightY).Should(Equal(int32(20)))
	}

	logs, err = db.LoadGuildLogs(context.TODO(), 1, shared_proto.GuildLogType_GLT_Memorabilia, 0, 10)
	Ω(err).Should(Succeed())
	Ω(logs).Should(HaveLen(1))

	for _, v := range logs {
		Ω(v.FightX).Should(Equal(int32(10)))
		Ω(v.FightY).Should(Equal(int32(20)))
	}

	logs, err = db.LoadGuildLogs(context.TODO(), 2, shared_proto.GuildLogType_GLTFight, 0, 10)
	Ω(err).Should(Succeed())
	Ω(logs).Should(HaveLen(1))

	for _, v := range logs {
		Ω(v.FightX).Should(Equal(int32(10)))
		Ω(v.FightY).Should(Equal(int32(20)))
	}

	err = db.DeleteGuild(context.TODO(), 1)
	Ω(err).Should(Succeed())

	logs, err = db.LoadGuildLogs(context.TODO(), 1, shared_proto.GuildLogType_GLTFight, 0, 10)
	Ω(err).Should(Succeed())
	Ω(logs).Should(BeEmpty())

	err = db.DeleteGuild(context.TODO(), 2)
	Ω(err).Should(Succeed())

	logs, err = db.LoadGuildLogs(context.TODO(), 2, shared_proto.GuildLogType_GLTFight, 0, 10)
	Ω(err).Should(Succeed())
	Ω(logs).Should(BeEmpty())
}

func newMail(id uint64, i int) *shared_proto.MailProto {

	keep := i&(1<<0) != 0
	readed := i&(1<<1) != 0
	hasReport := i&(1<<2) != 0
	hasPrize := i&(1<<3) != 0
	collected := i&(1<<4) != 0

	var report *shared_proto.FightReportProto
	if hasReport {
		report = &shared_proto.FightReportProto{
			ShowPrize: &shared_proto.PrizeProto{
				Gold: 100,
			},
		}
	}

	var prize *shared_proto.PrizeProto
	if hasPrize {
		prize = &shared_proto.PrizeProto{
			Gold: 100,
		}
	}

	proto := &shared_proto.MailProto{}
	proto.Id = i64.ToBytesU64(id)
	proto.Prize = prize
	proto.Report = report

	proto.Keep = keep
	proto.Read = readed
	proto.HasReport = hasReport
	proto.HasPrize = hasPrize
	proto.Collected = collected

	return proto
}

func testMail(db DbServiceAdapter) {

	var reportIds []uint64
	var systemIds []uint64

	n := 1 << 5
	for i := 0; i < n; i++ {

		id := uint64(i + 1)
		proto := newMail(id, i)

		if proto.HasReport {
			reportIds = append(reportIds, id)
		} else {
			systemIds = append(systemIds, id)
		}

		err := db.CreateMail(context.TODO(), id, heroId, must.Marshal(proto), false, proto.HasReport, proto.HasPrize, 0, time.Now().Unix())
		Ω(err).Should(Succeed())

		if proto.Keep {
			err = db.UpdateMailKeep(context.TODO(), id, heroId, false)
			Ω(err).Should(Succeed())
		} else {
			err = db.UpdateMailKeep(context.TODO(), id, heroId, true)
			Ω(err).Should(Succeed())
		}

		if proto.Read {
			err = db.UpdateMailRead(context.TODO(), id, heroId, false)
			Ω(err).Should(Succeed())
		} else {
			err = db.UpdateMailRead(context.TODO(), id, heroId, true)
			Ω(err).Should(Succeed())
		}

		if proto.Collected {
			err = db.UpdateMailCollected(context.TODO(), id, heroId, false)
			Ω(err).Should(Succeed())
		} else {
			err = db.UpdateMailCollected(context.TODO(), id, heroId, true)
			Ω(err).Should(Succeed())
		}
	}

	idGen, err := db.MaxMailId(context.TODO())
	Ω(err).Should(Succeed())
	Ω(idGen).Should(Equal(uint64(n)))

	for i := int32(0); i < 3; i++ {
		keep := i
		for i := int32(0); i < 3; i++ {
			readed := i
			for i := int32(0); i < 3; i++ {
				hasReport := i
				for i := int32(0); i < 3; i++ {
					hasPrize := i
					for i := int32(0); i < 3; i++ {
						collected := i

						checkFunc := func(v *shared_proto.MailProto) {
							if keep > 0 {
								Ω(v.Keep).Should(Equal(keep == 2))
							}

							if readed > 0 {
								Ω(v.Read).Should(Equal(readed == 2))
							}

							if hasReport > 0 {
								Ω(v.HasReport).Should(Equal(hasReport == 2))
							}

							if hasPrize > 0 {
								Ω(v.HasPrize).Should(Equal(hasPrize == 2))
							}

							if collected > 0 {
								Ω(v.Collected).Should(Equal(collected == 2))
							}

							if hasPrize > 0 && collected > 0 {
								id, _ := i64.FromBytesU64(v.Id)
								p, err := db.LoadCollectMailPrize(context.TODO(), id, heroId)
								Ω(err).Should(Succeed())

								if hasPrize == 2 && collected == 1 {
									Ω(p).ShouldNot(BeNil())
								} else {
									Ω(p).Should(BeNil())
								}
							}
						}

						// 首先不限制id
						id := uint64(0)
						list, err := db.LoadHeroMailList(context.TODO(), heroId, id, keep, readed, hasReport, hasPrize, collected, 0, 50)
						Ω(err).Should(Succeed())
						Ω(list).ShouldNot(BeEmpty())

						for _, v := range list {
							checkFunc(v)
						}

						// 取1个id，限制这个值
						if n := len(list); n > 0 {
							x := list[n/2]
							id, _ = i64.FromBytesU64(x.Id)
							id++

							list, err := db.LoadHeroMailList(context.TODO(), heroId, id, keep, readed, hasReport, hasPrize, collected, 0, 50)
							Ω(err).Should(Succeed())
							Ω(list).ShouldNot(BeEmpty())

							for _, v := range list {
								vid, _ := i64.FromBytesU64(v.Id)
								Ω(vid < id).Should(BeTrue())

								checkFunc(v)
							}
						}
					}
				}
			}
		}
	}

	noreport := int32(1)
	report := int32(2)

	// 先来一波删除已读
	// 删除已读
	f := func(isReport bool) {
		ids := systemIds
		intReport := noreport
		if isReport {
			ids = reportIds
			intReport = report
		}

		err := db.DeleteMultiMail(context.TODO(), heroId, ids, isReport)
		Ω(err).Should(Succeed())

		// 留下来的肯定是未读，并且有奖励可以领取的的邮件
		//list, err := db.LoadHeroMailList(context.TODO(), heroId, 0, 0, 0, intReport, 0, 0, 50)
		//Ω(err).Should(Succeed())
		//Ω(list).ShouldNot(BeEmpty())
		//
		//for _, v := range list {
		//	Ω(v.HasReport).Should(Equal(isReport))
		//
		//	if v.Keep || !v.Read {
		//		// 收藏邮件，或者未读邮件，不会被删除
		//		continue
		//	}
		//
		//	// 否则肯定是有奖励可以领取
		//	Ω(v.HasPrize && !v.Collected).Should(BeTrue())
		//}

		Eventually(func() bool {
			// 留下来的肯定是收藏的邮件
			list, err := db.LoadHeroMailList(context.TODO(), heroId, 0, 0, 0, intReport, 0, 0, 0, 50)
			Ω(err).Should(Succeed())
			Ω(list).ShouldNot(BeEmpty())

			for _, v := range list {
				if v.HasReport != isReport {
					return false
				}

				if v.Keep || !v.Read {
					// 收藏邮件，或者未读邮件，不会被删除
					continue
				}

				// 否则肯定是有奖励可以领取
				if !v.HasPrize || v.Collected {
					return false
				}
			}
			return true
		}, "10s").Should(BeTrue())
	}

	f(true)
	f(false)

	// 一键已读
	f = func(isReport bool) {

		ids := systemIds
		intReport := noreport
		if isReport {
			ids = reportIds
			intReport = report
		}

		// 一键已读
		prize, err := db.ReadMultiMail(context.TODO(), heroId, ids, isReport)
		Ω(err).Should(Succeed())

		if isReport {
			Ω(prize).Should(BeNil())
		} else {
			Ω(prize).ShouldNot(BeNil())
		}

		// 没有未读邮件
		//list, err := db.LoadHeroMailList(context.TODO(), heroId, 0, 0, 0, intReport, 0, 0, 50)
		//Ω(err).Should(Succeed())
		//Ω(list).ShouldNot(BeEmpty())
		//
		//for _, v := range list {
		//	Ω(v.HasReport).Should(Equal(isReport))
		//	Ω(v.Read).Should(BeTrue())
		//
		//	if !isReport && v.HasPrize {
		//		Ω(v.Collected).Should(BeTrue())
		//	}
		//}

		Eventually(func() bool {
			// 留下来的肯定是收藏的邮件
			list, err := db.LoadHeroMailList(context.TODO(), heroId, 0, 0, 0, intReport, 0, 0, 0, 50)
			Ω(err).Should(Succeed())
			Ω(list).ShouldNot(BeEmpty())

			for _, v := range list {
				if v.HasReport != isReport {
					return false
				}

				if !v.Read {
					return false
				}

				if !isReport && v.HasPrize {
					if !v.Collected {
						return false
					}
				}
			}
			return true
		}, "10s").Should(BeTrue())
	}

	f(true)
	f(false)

	// 删除已读
	f = func(isReport bool) {
		ids := systemIds
		intReport := noreport
		if isReport {
			ids = reportIds
			intReport = report
		}

		err := db.DeleteMultiMail(context.TODO(), heroId, ids, isReport)
		Ω(err).Should(Succeed())

		//// 留下来的肯定是收藏的邮件
		//list, err := db.LoadHeroMailList(context.TODO(), heroId, 0, 0, 0, intReport, 0, 0, 50)
		//Ω(err).Should(Succeed())
		//Ω(list).ShouldNot(BeEmpty())
		//
		//for _, v := range list {
		//	Ω(v.HasReport).Should(Equal(isReport))
		//	Ω(v.Keep).Should(BeTrue())
		//
		//	if !isReport && v.HasPrize {
		//		Ω(v.Collected).Should(BeTrue())
		//	}
		//}

		Eventually(func() bool {
			// 留下来的肯定是收藏的邮件
			list, err := db.LoadHeroMailList(context.TODO(), heroId, 0, 0, 0, intReport, 0, 0, 0, 50)
			Ω(err).Should(Succeed())
			Ω(list).ShouldNot(BeEmpty())

			for _, v := range list {
				if v.HasReport != isReport {
					return false
				}

				if !v.Keep {
					return false
				}

				if !isReport && v.HasPrize {
					if !v.Collected {
						return false
					}
				}
			}
			return true
		}, "10s").Should(BeTrue())
	}

	f(true)
	f(false)

	for i := 0; i < n; i++ {
		id := uint64(i + 1)
		err := db.DeleteMail(context.TODO(), id, heroId)
		Ω(err).Should(Succeed())
	}

	// 可能索引没更新，导致失败
	//idGen, err = db.MaxMailId()
	//Ω(err).Should(Succeed())
	//Ω(idGen).Should(Equal(uint64(0)))
}

func testBaizhan(db DbServiceAdapter) {

	ctime := int32(time.Now().Unix())

	attackerId := int64(9999)
	defenserId := int64(6666)
	replay := &shared_proto.BaiZhanReplayProto{
		Time: ctime,
	}

	err := db.InsertBaiZhanReplay(context.TODO(), attackerId, defenserId, replay, true, int64(replay.Time))
	Ω(err).Should(Succeed())

	// attacker
	Eventually(func() int32 {
		datas, err := db.LoadBaiZhanRecord(context.TODO(), attackerId, 1)
		Ω(err).Should(Succeed())

		result := &shared_proto.BaiZhanReplayProto{}
		if len(datas) > 0 {
			Ω(result.Unmarshal(datas[0])).Should(Succeed())
		}

		return result.Time
	}, "10s").Should(Equal(ctime))

	// defenser
	datas, err := db.LoadBaiZhanRecord(context.TODO(), defenserId, 1)
	Ω(err).Should(Succeed())

	if len(datas) > 0 {
		result := &shared_proto.BaiZhanReplayProto{}
		Ω(result.Unmarshal(datas[0])).Should(Succeed())
		Ω(result.Time).ShouldNot(Equal(ctime))
	}

	// defenser not a npc
	ctime++
	replay.Time = ctime
	err = db.InsertBaiZhanReplay(context.TODO(), attackerId, defenserId, replay, false, int64(replay.Time))
	Ω(err).Should(Succeed())

	// attacker
	Eventually(func() int32 {
		datas, err := db.LoadBaiZhanRecord(context.TODO(), attackerId, 1)
		Ω(err).Should(Succeed())

		result := &shared_proto.BaiZhanReplayProto{}
		if len(datas) > 0 {
			Ω(result.Unmarshal(datas[0])).Should(Succeed())
		}

		return result.Time
	}, "10s").Should(Equal(ctime))

	// defenser
	Eventually(func() int32 {
		datas, err := db.LoadBaiZhanRecord(context.TODO(), defenserId, 1)
		Ω(err).Should(Succeed())

		result := &shared_proto.BaiZhanReplayProto{}
		if len(datas) > 0 {
			Ω(result.Unmarshal(datas[0])).Should(Succeed())
		}

		return result.Time
	}, "10s").Should(Equal(ctime))

}

func testChat(db DbServiceAdapter) {
	//if _, ok := db.(*datastoredb.DatastoreDbService); ok {
	//	fmt.Println("datastore 不支持Chat")
	//	return
	//}

	testChatMsg(db)
	testChatWindow(db)
}

func testChatMsg(db DbServiceAdapter) {

	roomIdBytes := util.SafeMarshal(&shared_proto.ChatRoomId{
		T:         shared_proto.ChatType_ChatPrivate,
		MemberIds: [][]byte{idbytes.ToBytes(math.MaxInt64), idbytes.ToBytes(math.MaxInt64)},
	})

	msg1 := &shared_proto.ChatMsgProto{
		SendTime: 1,
	}
	_, err := db.AddChatMsg(context.TODO(), 1, roomIdBytes, msg1)
	Ω(err).Should(Succeed())
	Ω(msg1.ChatId).ShouldNot(BeEmpty())

	msg2 := &shared_proto.ChatMsgProto{
		SendTime: 2,
	}
	_, err = db.AddChatMsg(context.TODO(), 2, roomIdBytes, msg2)
	Ω(err).Should(Succeed())
	Ω(msg2.ChatId).ShouldNot(BeEmpty())

	// 查询所有
	msgs, err := db.ListHeroChatMsg(context.TODO(), roomIdBytes, 0)
	Ω(err).Should(Succeed())
	Ω(msgs).Should(Equal([]*shared_proto.ChatMsgProto{msg2, msg1}))

	// 查询msg2之前的
	minChatId, ok := i64.FromBytes(msg2.ChatId)
	Ω(ok).Should(BeTrue())
	msgs, err = db.ListHeroChatMsg(context.TODO(), roomIdBytes, uint64(minChatId))
	Ω(err).Should(Succeed())
	Ω(msgs).Should(Equal([]*shared_proto.ChatMsgProto{msg1}))

	// 查询msg1之前的
	minChatId, ok = i64.FromBytes(msg1.ChatId)
	Ω(ok).Should(BeTrue())
	msgs, err = db.ListHeroChatMsg(context.TODO(), roomIdBytes, uint64(minChatId))
	Ω(err).Should(Succeed())
	Ω(msgs).Should(BeEmpty())

	err = db.RemoveChatMsg(context.TODO(), 1)
	Ω(err).Should(Succeed())
	err = db.RemoveChatMsg(context.TODO(), 2)
	Ω(err).Should(Succeed())

	msgs, err = db.ListHeroChatMsg(context.TODO(), roomIdBytes, 0)
	Ω(err).Should(Succeed())
	Ω(msgs).Should(BeEmpty())

	roomIdBytes1 := make([]byte, 255)
	roomIdBytes2 := make([]byte, 256)
	roomIdBytes3 := make([]byte, 257)
	roomIdBytes4 := make([]byte, 257)

	for i := 0; i < 255; i++ {
		roomIdBytes1[i] = byte(i)
		roomIdBytes2[i] = byte(i)
		roomIdBytes3[i] = byte(i)
		roomIdBytes4[i] = byte(i)
	}
	copy(roomIdBytes1[255:], roomIdBytes1)
	copy(roomIdBytes2[255:], roomIdBytes2)
	copy(roomIdBytes3[255:], roomIdBytes3)
	copy(roomIdBytes4[255:], roomIdBytes4)
	roomIdBytes4[len(roomIdBytes4)-1] = 255

	ids := [][]byte{roomIdBytes1, roomIdBytes2, roomIdBytes3, roomIdBytes4}
	msgs = make([]*shared_proto.ChatMsgProto, len(ids))
	for i, roomIdBytes := range ids {
		//fmt.Println(roomIdBytes)
		msgs[i] = &shared_proto.ChatMsgProto{
			SendTime: int32(i),
		}
		_, err = db.AddChatMsg(context.TODO(), 2, roomIdBytes, msgs[i])
		Ω(err).Should(Succeed())
		Ω(msgs[i].ChatId).ShouldNot(BeEmpty())
	}

	for i, roomIdBytes := range ids {
		msgs0, err := db.ListHeroChatMsg(context.TODO(), roomIdBytes, 0)
		Ω(err).Should(Succeed())
		Ω(msgs0).Should(Equal(msgs[i : i+1]))
	}

	err = db.RemoveChatMsg(context.TODO(), 2)
	Ω(err).Should(Succeed())

	for _, roomIdBytes := range ids {
		msgs0, err := db.ListHeroChatMsg(context.TODO(), roomIdBytes, 0)
		Ω(err).Should(Succeed())
		Ω(msgs0).Should(BeEmpty())
	}
}

func testChatWindow(db DbServiceAdapter) {

	//roomIdBytes1 := util.SafeMarshal(&shared_proto.ChatRoomId{
	//	T:         shared_proto.ChatType_ChatPrivate,
	//	MemberIds: [][]byte{idbytes.ToBytes(math.MaxInt64), idbytes.ToBytes(math.MaxInt64),},
	//})
	//
	//roomIdBytes2 := util.SafeMarshal(&shared_proto.ChatRoomId{
	//	T:         shared_proto.ChatType_ChatPrivate,
	//	MemberIds: [][]byte{idbytes.ToBytes(1), idbytes.ToBytes(9850),},
	//})

	roomIdBytes1 := make([]byte, 256)
	roomIdBytes2 := make([]byte, 256)
	for i := 0; i < 255; i++ {
		roomIdBytes1[i] = byte(i)
		roomIdBytes2[i] = byte(i + 1)
	}

	target := &shared_proto.ChatSenderProto{
		Id: idbytes.ToBytes(996),
	}

	err := db.UpdateChatWindow(context.TODO(), 1, roomIdBytes1, util.SafeMarshal(target), false, 2, false)
	Ω(err).Should(Succeed())

	unread, target0, err := db.ListHeroChatWindow(context.TODO(), 1)
	Ω(err).Should(Succeed())
	Ω(unread).Should(Equal([]uint64{0}))
	Ω([][]byte(target0)).Should(Equal([][]byte{util.SafeMarshal(target)}))

	// 再保存一次
	target.Name = "haha"
	err = db.UpdateChatWindow(context.TODO(), 1, roomIdBytes1, util.SafeMarshal(target), false, 3, false)
	Ω(err).Should(Succeed())

	unread, target0, err = db.ListHeroChatWindow(context.TODO(), 1)
	Ω(err).Should(Succeed())
	Ω(unread).Should(Equal([]uint64{0}))
	Ω([][]byte(target0)).Should(Equal([][]byte{util.SafeMarshal(target)}))

	// roomIdBytes2
	target2 := &shared_proto.ChatSenderProto{
		Id: idbytes.ToBytes(9850),
	}

	err = db.UpdateChatWindow(context.TODO(), 1, roomIdBytes2, util.SafeMarshal(target2), true, 1, true)
	Ω(err).Should(Succeed())

	unread, target0, err = db.ListHeroChatWindow(context.TODO(), 1)
	Ω(err).Should(Succeed())
	Ω(unread).Should(Equal([]uint64{0, 1}))
	Ω([][]byte(target0)).Should(Equal([][]byte{util.SafeMarshal(target), util.SafeMarshal(target2)}))

	// 再保存一次
	target2.Name = "haha2"
	err = db.UpdateChatWindow(context.TODO(), 1, roomIdBytes2, util.SafeMarshal(target2), true, 3, true)
	Ω(err).Should(Succeed())

	unread, target0, err = db.ListHeroChatWindow(context.TODO(), 1)
	Ω(err).Should(Succeed())
	Ω(unread).Should(Equal([]uint64{2, 0}))
	Ω([][]byte(target0)).Should(Equal([][]byte{util.SafeMarshal(target2), util.SafeMarshal(target)}))

	// 读消息
	err = db.ReadChat(context.TODO(), 1, roomIdBytes2)
	Ω(err).Should(Succeed())

	unread, target0, err = db.ListHeroChatWindow(context.TODO(), 1)
	Ω(err).Should(Succeed())
	Ω(unread).Should(Equal([]uint64{0, 0}))
	Ω([][]byte(target0)).Should(Equal([][]byte{util.SafeMarshal(target2), util.SafeMarshal(target)}))

	// 删除
	err = db.DeleteChatWindow(context.TODO(), 1, roomIdBytes1)
	Ω(err).Should(Succeed())
	err = db.DeleteChatWindow(context.TODO(), 1, roomIdBytes2)
	Ω(err).Should(Succeed())

	unread, target0, err = db.ListHeroChatWindow(context.TODO(), 1)
	Ω(err).Should(Succeed())
	Ω(unread).Should(BeEmpty())
	Ω(target0).Should(BeEmpty())
}

func testFarm(db DbServiceAdapter) {
	//if _, ok := db.(*datastoredb.DatastoreDbService); ok {
	//	fmt.Println("datastore 不支持Farm")
	//	return
	//}

	ctime := time.Now()
	ripeDuration, _ := time.ParseDuration("1h")
	ripeTime := ctime.Add(ripeDuration)
	d, _ := time.ParseDuration("24h")
	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())

	heroId := int64(10001)
	cubes := make(map[int]*entity.FarmCube, 0)
	for i := 0; i < 10; i++ {
		cube := entity.NewFarmCube(heroId, cb.XYCube(i, i), ctime, ripeTime, datas.GetFarmResConfig(1), datas.FarmMiscConfig())
		cube.ConflictTime = ctime.Add(d)
		cubes[i] = cube
	}

	// remove
	for _, cube := range cubes {
		err := db.RemoveFarmCube(context.TODO(), heroId, cube.Cube)
		Ω(err).Should(Succeed())
	}

	// create
	for _, cube := range cubes {
		err := db.CreateFarmCube(context.TODO(), cube)
		Ω(err).Should(Succeed())
	}

	// load
	result, err := db.LoadFarmCube(context.TODO(), heroId, cb.XYCube(1, 1))
	Ω(err).Should(Succeed())
	Ω(timeutil.Marshal64(result.ConflictTime)).Should(Equal(timeutil.Marshal64(ctime.Add(d))))

	// save
	saveCubes := make([]*entity.FarmCube, 0)
	for _, i := range []int{2, 3} {
		saveCube := cubes[i]
		saveCube.ConflictTime = time.Time{}
		saveCubes = append(saveCubes, saveCube)
	}
	for _, i := range []int{1, 3} {
		saveCube := cubes[i]
		saveCube.RemoveTime = ctime
		saveCubes = append(saveCubes, saveCube)
	}
	for _, saveCube := range saveCubes {
		err = db.SaveFarmCube(context.TODO(), saveCube)
		Ω(err).Should(Succeed())
	}

	err = db.PlantFarmCube(context.TODO(), heroId, cubes[1].Cube, timeutil.Marshal64(ctime), timeutil.Marshal64(ripeTime), 1)
	Ω(err).Should(Succeed())
	err = db.UpdateFarmCubeState(context.TODO(), heroId, cubes[1].Cube, timeutil.Marshal64(ctime), timeutil.Marshal64(ctime))
	Ω(err).Should(Succeed())

	// load harvest
	results, err := db.LoadFarmHarvestCubes(context.TODO(), heroId)
	Ω(err).Should(Succeed())
	Ω(len(results)).Should(Equal(2))

	// load all
	results, err = db.LoadFarmCubes(context.TODO(), heroId)
	Ω(err).Should(Succeed())
	Ω(len(results)).Should(Equal(10))

	// load steal
	results, err = db.LoadFarmStealCubes(context.TODO(), heroId, timeutil.Marshal64(ripeTime), 3)
	Ω(err).Should(Succeed())
	Ω(len(results)).Should(Equal(0))

	// update steal time
	cbs := []cb.Cube{cb.XYCube(1, 1), cb.XYCube(3, 3), cb.XYCube(5, 5)}
	for i := 0; i < 3; i++ {
		err = db.UpdateFarmStealTimes(context.TODO(), heroId, cbs)
	}
	Ω(err).Should(Succeed())

	// load
	for _, cb := range cbs {
		result, err = db.LoadFarmCube(context.TODO(), heroId, cb)
		Ω(err).Should(Succeed())
		Ω(result.StealTimes).Should(Equal(uint64(3)))
	}
	result, err = db.LoadFarmCube(context.TODO(), heroId, cb.XYCube(2, 2))
	Ω(err).Should(Succeed())
	Ω(result.StealTimes).Should(Equal(uint64(0)))

	canStealCube, err := db.LoadCanStealCube(context.TODO(), heroId, 10002, timeutil.Marshal64(ripeTime), 3)
	Ω(err).Should(Succeed())
	Ω(len(canStealCube)).Should(Equal(0))

	// remove
	for _, cube := range cubes {
		err := db.RemoveFarmCube(context.TODO(), heroId, cube.Cube)
		Ω(err).Should(Succeed())
	}

}

func textFarmLog(db DbServiceAdapter) {
	//if _, ok := db.(*datastoredb.DatastoreDbService); ok {
	//	fmt.Println("datastore 不支持Farm")
	//	return
	//}

	ctime := time.Now()
	heroId := int64(100001)
	// create
	for i := 0; i < 10; i++ {
		log := &shared_proto.FarmStealLogProto{}
		log.HeroId = idbytes.ToBytes(heroId)
		log.ThiefName = "a"
		log.ThiefGuildFlag = "b"
		log.GoldOutput = 3
		log.StoneOutput = 4
		d, _ := time.ParseDuration(fmt.Sprintf("-%vh", i*24))
		log.LogTime = timeutil.Marshal32(ctime.Add(d))
		err := db.CreateFarmLog(context.TODO(), log)
		Ω(err).Should(Succeed())
	}

	// load
	d, _ := time.ParseDuration(fmt.Sprintf("-%vh", 24*2))
	fls, err := db.LoadFarmLog(context.TODO(), heroId, 3)
	Ω(err).Should(Succeed())
	Ω(len(fls)).Should(Equal(3))

	// remove
	d, _ = time.ParseDuration("1h")
	err = db.RemoveFarmLog(context.TODO(), timeutil.Marshal32(ctime.Add(d)))
	Ω(err).Should(Succeed())
}

func testFarmSteal(db DbServiceAdapter) {
	//if _, ok := db.(*datastoredb.DatastoreDbService); ok {
	//	fmt.Println("datastore 不支持Farm")
	//	return
	//}

	heroId := int64(10001)
	thiefId := int64(10002)

	cubes := make([]cb.Cube, 0)
	for i := 0; i < 100; i++ {
		cube := cb.XYCube(i, i)
		err := db.AddFarmSteal(context.TODO(), heroId, thiefId, cube)
		Ω(err).Should(Succeed())
		cubes = append(cubes, cube)
	}

	count, err := db.LoadFarmStealCount(context.TODO(), heroId, thiefId, cubes[1])
	Ω(err).Should(Succeed())
	Ω(count).Should(Equal(uint64(1)))

	ctime := time.Now()
	ripeTimeInt := timeutil.Marshal64(ctime)
	canCount, err := db.LoadCanStealCount(context.TODO(), heroId, thiefId, ripeTimeInt, 3)
	Ω(err).Should(Succeed())
	Ω(canCount).Should(Equal(uint64(0)))

	err = db.RemoveFarmSteal(context.TODO(), heroId, cubes)
	Ω(err).Should(Succeed())
}
