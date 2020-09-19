package sqldb

import (
	. "github.com/onsi/gomega"
	"context"
	"github.com/lightpaw/male7/pb/shared_proto"
	"time"
	"github.com/lightpaw/male7/util/i64"
	"database/sql"
	"fmt"
	"path"
	"github.com/lightpaw/male7/service/timeservice"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/confpath"
	"testing"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/entity"
	"sort"
)

//func TestQuery(t *testing.T) {
//	RegisterTestingT(t)
//
//	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
//	Ω(err).Should(Succeed())
//
//	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/xianyou"
//	if db := newMysqlAdapter(datas, dataSourceName, false); db != nil {
//
//		rows, err := db.db.Query("select id, misc from user where id >= 11201 and id < 11701")
//		Ω(err).Should(Succeed())
//
//		totalCount := 0
//		completeCount := 0
//		newUserMap := make(map[uint64]uint64)
//		for rows.Next() {
//			var id int64
//			var data []byte
//			err := rows.Scan(&id, &data)
//			Ω(err).Should(Succeed())
//
//			proto := server_proto.UserMiscProto{}
//			err = proto.Unmarshal(data)
//			Ω(err).Should(Succeed())
//
//			totalCount++
//			if proto.IsTutorialComplete {
//				completeCount++
//			}
//
//			newUserMap[uint64(proto.TutorialProgress)]++
//		}
//
//		heros, err := db.LoadAllHeroData(context.TODO())
//		Ω(err).Should(Succeed())
//
//		totalCreate := 0
//		create18 := 0
//		create18Login19 := 0
//		create19 := 0
//		time18, err := time.ParseInLocation("2006-01-02", "2018-04-18", timeutil.GameZone)
//		Ω(err).Should(Succeed())
//		time19, err := time.ParseInLocation("2006-01-02", "2018-04-19", timeutil.GameZone)
//		Ω(err).Should(Succeed())
//		fmt.Println(time18, time19)
//
//		baseLevelMap := make(map[uint64]uint64)
//		heroLevelMap := make(map[uint64]uint64)
//		guanfuLevelMap := make(map[uint64]uint64)
//		completedBayeStageMap := make(map[uint64]uint64)
//		completedBayeTaskMap := make(map[uint64]uint64)
//		collectedBayeTaskMap := make(map[uint64]uint64)
//		completedAchiveTaskMap := make(map[uint64]uint64)
//		collectedAchiveTaskMap := make(map[uint64]uint64)
//		completedBwzkTaskMap := make(map[uint64]uint64)
//		collectedBwzkTaskMap := make(map[uint64]uint64)
//		for _, hero := range heros {
//			if hero.Id() <= 11201 || hero.Id() >= 11701 {
//				continue
//			}
//
//			if time18.Before(hero.CreateTime()) {
//				if time19.After(hero.CreateTime()) {
//					create18++
//					if hero.LastOfflineTime().After(time19) {
//						create18Login19++
//					}
//				} else {
//					create19++
//				}
//			}
//			totalCreate++
//
//			baseLevelMap[hero.HomeHistoryMaxLevel()]++
//			heroLevelMap[hero.Level()]++
//			guanfuLevelMap[hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU).Level]++
//
//			completedBayeStageMap[ hero.TaskList().GetCompletedBaYeStage()]++
//
//			for _, v := range datas.GetBaYeStageDataArray() {
//				for _, t := range v.Tasks {
//					if hero.TaskList().GetCompletedBaYeStage() >= v.Stage {
//						completedBayeTaskMap[t.Id]++
//						collectedBayeTaskMap[t.Id]++
//					} else {
//						if v.Stage == hero.TaskList().BaYeStage().Data().Stage {
//							task := hero.TaskList().BaYeStage().GetTask(t.Id)
//							if task != nil {
//								if task.Progress().IsCompleted() {
//									completedBayeTaskMap[t.Id]++
//								}
//								if task.IsCollectPrize() {
//									collectedBayeTaskMap[t.Id]++
//								}
//							}
//						}
//					}
//				}
//			}
//
//			for _, v := range datas.GetAchieveTaskDataArray() {
//				task := hero.TaskList().AchieveTaskList().GetTaskByAchieveType(v.AchieveType)
//				if task != nil {
//
//					if task.Data().TotalStar > v.TotalStar {
//						completedAchiveTaskMap[v.Id]++
//						collectedAchiveTaskMap[v.Id]++
//					} else if task.Data().TotalStar == v.TotalStar && task.Progress().IsCompleted() {
//						completedAchiveTaskMap[v.Id]++
//
//						if task.IsCollectPrize() {
//							collectedAchiveTaskMap[v.Id]++
//						}
//					}
//				}
//			}
//
//			for _, v := range datas.GetBwzlTaskDataArray() {
//				task := hero.TaskList().BwzlTaskList().Task(v.Id)
//				if task != nil {
//					if task.Progress().IsCompleted() {
//						completedBwzkTaskMap[v.Id]++
//					}
//
//					if task.IsCollectPrize() {
//						collectedBwzkTaskMap[v.Id]++
//					}
//				}
//			}
//
//		}
//
//		fmt.Println("创建角色总数", totalCreate)
//		fmt.Println("18号创建角色", create18)
//		fmt.Println("18号创角次日留存", create18Login19)
//		fmt.Println("19号创建角色", create19)
//
//		f := func(m map[uint64]uint64) string {
//
//			var ks []uint64
//			for k := range m {
//				ks = append(ks, k)
//			}
//			sortkeys.Uint64s(ks)
//
//			b := bytes.Buffer{}
//			b.WriteString("\n")
//			for _, k := range ks {
//				b.WriteString(fmt.Sprintf("%d	%d\n", k, m[k]))
//			}
//
//			return b.String()
//		}
//
//		fmt.Println("主城等级分布（历史最高）", f(baseLevelMap))
//		fmt.Println("君主等级分布", f(heroLevelMap))
//		fmt.Println("官府等级分布", f(guanfuLevelMap))
//		fmt.Println("完成霸业阶段分布", f(completedBayeStageMap))
//		fmt.Println("完成霸业任务分布", f(completedBayeTaskMap))
//		fmt.Println("领取霸业任务分布", f(collectedBayeTaskMap))
//		fmt.Println("完成成就任务分布", f(completedAchiveTaskMap))
//		fmt.Println("领取成就任务分布", f(collectedAchiveTaskMap))
//		fmt.Println("完成霸王之路分布", f(completedBwzkTaskMap))
//		fmt.Println("领取霸王之路分布", f(collectedBwzkTaskMap))
//
//		fmt.Println("新手总数", totalCount)
//		fmt.Println("完成新手引导", completeCount)
//		fmt.Println("新手引导分布", f(newUserMap))
//	}
//}

func TestMysqlTicker(t *testing.T) {
	RegisterTestingT(t)

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	Ω(err).Should(Succeed())

	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/male7test"
	if db := newMysqlAdapter(datas, dataSourceName, true); db != nil {
		//testTickDeleteExpireChat(db)
		//testTickDeleteExpireGuildLog(db)
		//testTickDeleteExpireMail(db)
		testHeroField(db)
	}
}

func TestMysqlTicker1(t *testing.T) {
	RegisterTestingT(t)

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	Ω(err).Should(Succeed())

	createTime := "2018-08-28"
	totalCount := 0
	loginCountMap := make(map[string]uint64)

	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	if db := newMysqlAdapter(datas, dataSourceName, false); db != nil {
		heros, err := db.LoadAllHeroData(context.TODO())
		Ω(err).Should(Succeed())

		for _, hero := range heros {
			heroCreateTime := hero.CreateTime().Format(timeutil.DayLayout)

			if heroCreateTime == createTime {
				totalCount++

				lastTime := timeutil.Max(hero.LastOfflineTime(), hero.LastOnlineTime()).Format(timeutil.DayLayout)
				loginCountMap[lastTime]++
			}
		}
	}

	fmt.Println(totalCount)

	var keys []string
	for k := range loginCountMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	var count uint64
	for i := 0; i < len(keys); i++ {
		k := keys[len(keys)-1-i]
		count += loginCountMap[k]
		fmt.Println(k, count, loginCountMap[k])
	}
}

func testTickDeleteExpireChat(db *SqlDbService) {
	ctime := time.Now()

	room0 := []byte{0}

	proto := &shared_proto.ChatMsgProto{}
	proto.ChatType = shared_proto.ChatType_ChatWorld
	proto.SendTime = int32(ctime.Add(-time.Hour).Unix())
	_, err := db.AddChatMsg(context.TODO(), 1, room0, proto)
	Ω(err).Should(Succeed())

	proto.SendTime = int32(ctime.Add(-time.Minute).Unix())
	_, err = db.AddChatMsg(context.TODO(), 2, room0, proto)
	Ω(err).Should(Succeed())

	proto.ChatType = shared_proto.ChatType_ChatGuild
	proto.SendTime = int32(ctime.Add(-time.Minute).Unix())
	room1 := []byte{1}
	_, err = db.AddChatMsg(context.TODO(), 1, room1, proto)
	Ω(err).Should(Succeed())

	// 2个都还在
	db.tickDeleteExpireChat(shared_proto.ChatType_ChatWorld, ctime, time.Hour+time.Second)

	chats, err := db.ListHeroChatMsg(context.TODO(), room0, 0)
	Ω(err).Should(Succeed())
	Ω(chats).Should(HaveLen(2))

	// 只剩1个
	db.tickDeleteExpireChat(shared_proto.ChatType_ChatWorld, ctime, time.Hour)

	chats, err = db.ListHeroChatMsg(context.TODO(), room0, 0)
	Ω(err).Should(Succeed())
	Ω(chats).Should(HaveLen(1))

	db.tickDeleteExpireChat(shared_proto.ChatType_ChatWorld, ctime, time.Minute)

	chats, err = db.ListHeroChatMsg(context.TODO(), room0, 0)
	Ω(err).Should(Succeed())
	Ω(chats).Should(BeEmpty())

	// 联盟聊天还在
	chats, err = db.ListHeroChatMsg(context.TODO(), room1, 0)
	Ω(err).Should(Succeed())
	Ω(chats).Should(HaveLen(1))

	db.tickDeleteExpireChat(shared_proto.ChatType_ChatGuild, ctime, time.Minute)

	chats, err = db.ListHeroChatMsg(context.TODO(), room1, 0)
	Ω(err).Should(Succeed())
	Ω(chats).Should(HaveLen(0))
}

func testTickDeleteExpireGuildLog(db *SqlDbService) {

	ts := []shared_proto.GuildLogType{
		shared_proto.GuildLogType_GLTDaily,
		shared_proto.GuildLogType_GLTFight,
		shared_proto.GuildLogType_GLTPrestige,
		shared_proto.GuildLogType_GLT_Memorabilia,
	}

	guildId := int64(50)

	for _, t := range ts {
		proto := &shared_proto.GuildLogProto{}
		proto.Type = t

		for i := 0; i < 10; i++ {
			err := db.InsertGuildLog(context.TODO(), guildId, proto)
			Ω(err).Should(Succeed())
		}

		logs, err := db.LoadGuildLogs(context.TODO(), guildId, t, 0, 10)
		Ω(err).Should(Succeed())
		Ω(logs).Should(HaveLen(10))
	}

	limit := uint64(5)
	db.tickDeleteExpireGuildLog(limit)
	for _, t := range ts {
		if t != shared_proto.GuildLogType_GLT_Memorabilia {
			logs, err := db.LoadGuildLogs(context.TODO(), guildId, t, 0, 10)
			Ω(err).Should(Succeed())
			Ω(logs).Should(HaveLen(int(limit)))
		} else {
			logs, err := db.LoadGuildLogs(context.TODO(), guildId, t, 0, 10)
			Ω(err).Should(Succeed())
			Ω(logs).Should(HaveLen(10))
		}
	}

	// 没变化，
	db.tickDeleteExpireGuildLog(8)
	for _, t := range ts {
		if t != shared_proto.GuildLogType_GLT_Memorabilia {
			logs, err := db.LoadGuildLogs(context.TODO(), guildId, t, 0, 10)
			Ω(err).Should(Succeed())
			Ω(logs).Should(HaveLen(int(limit)))
		} else {
			logs, err := db.LoadGuildLogs(context.TODO(), guildId, t, 0, 10)
			Ω(err).Should(Succeed())
			Ω(logs).Should(HaveLen(10))
		}
	}

	limit = uint64(3)
	db.tickDeleteExpireGuildLog(limit)
	for _, t := range ts {
		if t != shared_proto.GuildLogType_GLT_Memorabilia {
			logs, err := db.LoadGuildLogs(context.TODO(), guildId, t, 0, 10)
			Ω(err).Should(Succeed())
			Ω(logs).Should(HaveLen(int(limit)))
		} else {
			logs, err := db.LoadGuildLogs(context.TODO(), guildId, t, 0, 10)
			Ω(err).Should(Succeed())
			Ω(logs).Should(HaveLen(10))
		}
	}
}

func testTickDeleteExpireMail(db *SqlDbService) {

	heroId := int64(10110)
	proto := &shared_proto.MailProto{}

	proto.Id = i64.ToBytes(1)
	data, _ := proto.Marshal()
	err := db.CreateMail(context.TODO(), 1, heroId, data,
		true, false, false, 0, 1)
	Ω(err).Should(Succeed())

	proto.Id = i64.ToBytes(2)
	data, _ = proto.Marshal()
	err = db.CreateMail(context.TODO(), 2, heroId, data,
		false, false, false, 0, 1)
	Ω(err).Should(Succeed())

	proto.Id = i64.ToBytes(3)
	data, _ = proto.Marshal()
	err = db.CreateMail(context.TODO(), 3, heroId, data,
		true, false, true, 0, 1)
	Ω(err).Should(Succeed())

	proto.Id = i64.ToBytes(4)
	data, _ = proto.Marshal()
	err = db.CreateMail(context.TODO(), 4, heroId, data,
		false, false, true, 0, 1)
	Ω(err).Should(Succeed())
	err = db.UpdateMailCollected(context.TODO(), 4, heroId, true)
	Ω(err).Should(Succeed())

	proto.Id = i64.ToBytes(5)
	data, _ = proto.Marshal()
	err = db.CreateMail(context.TODO(), 5, heroId, data,
		false, false, true, 0, 1)
	Ω(err).Should(Succeed())

	f := func(ids []uint64) {
		list, err := db.LoadHeroMailList(context.TODO(), heroId, 0, 0, 0, 0, 0, 0, 0, 10)
		Ω(err).Should(Succeed())

		var listIds []uint64
		for _, v := range list {
			id, _ := i64.FromBytesU64(v.Id)
			listIds = append(listIds, id)
		}
		Ω(listIds).Should(Equal(ids))
	}

	db.tickDeleteExpireMail(5)
	f([]uint64{5, 4, 3, 2, 1})

	db.tickDeleteExpireMail(4)
	f([]uint64{5, 4, 3, 1})

	db.tickDeleteExpireMail(3)
	f([]uint64{5, 3, 1})

	db.tickDeleteExpireMail(2)
	f([]uint64{5, 3})

	db.tickDeleteExpireMail(1)
	f([]uint64{5})
}

func testHeroField(db *SqlDbService) {

	for i := 0; i < 10; i++ {
		heroId := int64(i)
		heroName := fmt.Sprintf("君主%d", i)

		hero := entity.NewHero(heroId, heroName, db.datas.HeroInitData(), time.Now())
		result, err := db.CreateHero(context.TODO(), hero)
		Ω(err).Should(Succeed())
		Ω(result).Should(BeTrue())
	}

	field, err := db.loadHeroField(context.TODO(), offlineBoolFieldName, 1)
	Ω(err).Should(Succeed())
	Ω(field).Should(Equal(int64(0)))

	ids, err := db.loadHeroIdsByFieldValue(context.TODO(), offlineBoolFieldName, []int64{1, 2, 3, 4}, 1)
	Ω(err).Should(Succeed())
	Ω(ids).Should(BeEmpty())

	ids, err = db.loadHeroIdsByFieldValue(context.TODO(), offlineBoolFieldName, []int64{1, 2, 3, 4}, 0)
	Ω(err).Should(Succeed())
	Ω(ids).Should(ConsistOf([]int64{1, 2, 3, 4}))

	err = db.updateHeroField(context.TODO(), offlineBoolFieldName, 1, 1)
	Ω(err).Should(Succeed())

	suc, err := db.updateHeroFieldIfExpected(context.TODO(), offlineBoolFieldName, 2, 1, 2)
	Ω(err).Should(Succeed())
	Ω(suc).Should(BeFalse())

	suc, err = db.updateHeroFieldIfExpected(context.TODO(), offlineBoolFieldName, 3, 0, 1)
	Ω(err).Should(Succeed())
	Ω(suc).Should(BeTrue())

	ids, err = db.loadHeroIdsByFieldValue(context.TODO(), offlineBoolFieldName, []int64{1, 2, 3, 4}, 1)
	Ω(err).Should(Succeed())
	Ω(ids).Should(ConsistOf([]int64{1, 3}))

	b, err := db.loadHeroBoolField(context.TODO(), offlineBoolFieldName, 5, 0)
	Ω(err).Should(Succeed())
	Ω(b).Should(BeFalse())

	ids, err = db.loadHeroIdsByBoolFieldValue(context.TODO(), offlineBoolFieldName, []int64{5, 6, 7, 8}, 1, true)
	Ω(err).Should(Succeed())
	Ω(ids).Should(BeEmpty())

	ids, err = db.loadHeroIdsByBoolFieldValue(context.TODO(), offlineBoolFieldName, []int64{5, 6, 7, 8}, 1, false)
	Ω(err).Should(Succeed())
	Ω(ids).Should(ConsistOf([]int64{5, 6, 7, 8}))

	err = db.updateHeroBoolField(context.TODO(), offlineBoolFieldName, 5, 1, true)
	Ω(err).Should(Succeed())

	b, err = db.updateHeroBoolFieldIfExpected(context.TODO(), offlineBoolFieldName, 6, 1, true, false)
	Ω(err).Should(Succeed())
	Ω(b).Should(BeFalse())

	b, err = db.updateHeroBoolFieldIfExpected(context.TODO(), offlineBoolFieldName, 7, 1, false, true)
	Ω(err).Should(Succeed())
	Ω(b).Should(BeTrue())

	ids, err = db.loadHeroIdsByBoolFieldValue(context.TODO(), offlineBoolFieldName, []int64{5, 6, 7, 8}, 1, true)
	Ω(err).Should(Succeed())
	Ω(ids).Should(ConsistOf([]int64{5, 7}))

	b, err = db.updateHeroBoolFieldIfExpected(context.TODO(), offlineBoolFieldName, 7, 1, true, false)
	Ω(err).Should(Succeed())
	Ω(b).Should(BeTrue())

	ids, err = db.loadHeroIdsByBoolFieldValue(context.TODO(), offlineBoolFieldName, []int64{5, 6, 7, 8}, 1, true)
	Ω(err).Should(Succeed())
	Ω(ids).Should(ConsistOf([]int64{5}))
}

func newMysqlAdapter(datas *config.ConfigDatas, dataSourceName string, truncate bool) *SqlDbService {

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

	ds, err := NewSqlDbService(db, dbName, datas, timeservice.NewDefaultTimeService(), ifacemock.MetricsRegister)
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

func BenchmarkUpdate(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	stmt, err := db.Prepare("update test set b=? where id=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	ids := []int{1, 2, 3}
	bb := true
	n := len(ids)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%n]
		stmt.Exec(&bb, &id)
	}
}

func BenchmarkGetUpdate(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	stmt, err := db.Prepare("update test set b=? where id=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	queryStmt, err := db.Prepare("select b from test where id=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	ids := []int{1, 2, 3}
	bb := true
	n := len(ids)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%n]

		var bbb bool
		queryStmt.QueryRow(&id).Scan(&bbb)

		if !bbb {
			stmt.Exec(&bb, &id)
		}
	}
}

func BenchmarkBool1Match(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	//db.Exec("delete from hero_bool where id = 1")
	//db.Exec("insert into hero_bool(id,type,bool)values(1,1,true)")
	db.Exec("truncate table hero_bool")
	n := 1000
	for i := 0; i < n; i++ {
		db.Exec(fmt.Sprintf("insert into hero_bool(id,type,bool)values(%d,1,true)", i))
	}

	stmt, err := db.Prepare("update hero_bool set bool=bool&? where id=? and type=? and bool&?<>0")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	orStmt, _ := db.Prepare("update hero_bool set bool=bool|? where id=? and type=? and bool&?=0")

	t := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := i % n
		b := i % 1
		if b == 0 {
			stmt.Exec(&b, &id, &t, &b)
		} else {
			c := ^b
			orStmt.Exec(&c, &id, &t, &b)
		}
	}
}

func BenchmarkBool1Match2(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("truncate table hero_bool")
	for i := 0; i < 100; i++ {
		db.Exec(fmt.Sprintf("insert into hero_bool(id,type,bool)values(%d,1,true)", i))
	}

	stmt, err := db.Prepare("update hero_bool set bool=? where id in (?) and type=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	//id := 1
	t := 1

	ids := "1,3,4,6,7,8,9"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b := i%1 == 0
		stmt.Exec(&b, &ids, &t)
	}
}

func BenchmarkBool1Delete(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("truncate table hero_bool")
	for i := 0; i < 100; i++ {
		db.Exec(fmt.Sprintf("insert into hero_bool(id,type,bool)values(%d,1,true)", i))
	}

	stmt, err := db.Prepare("delete from hero_bool where id=? and type=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	//id := 1
	t := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stmt.Exec(&i, &t)
	}
}

func BenchmarkBool1Empty(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("delete from hero_bool where id = 1")
	//db.Exec("insert into hero_bool(id,type,bool)values(1,1,true)")

	stmt, err := db.Prepare("update hero_bool set bool=? where id=? and type=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	id := 1
	t := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b := i%1 == 0
		stmt.Exec(&b, &id, &t)
	}
}

func BenchmarkBool1Empty2(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("delete from hero_bool where id = 1")
	db.Exec("insert into hero_bool(id,type,bool)values(1,1,true)")

	stmt, err := db.Prepare("update hero_bool set bool=? where id=? and type=? and bool=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	id := 1
	t := 1
	bl := false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b := i%2 == 0
		stmt.Exec(&b, &id, &t, &bl)
	}
}

func BenchmarkBool1Get(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("truncate table hero_bool")
	for i := 0; i < 100; i++ {
		db.Exec(fmt.Sprintf("insert into hero_bool(id,type,bool)values(%d,1,true)", i))
	}
	stmt, err := db.Prepare("select bool from hero_bool where id in (?) and type=? and bool&?<>0")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	id := "1,2,3,4,5,6,7,8,9"
	t := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var b bool
		stmt.QueryRow(&id, &t, &t).Scan(&b)
	}
}

func BenchmarkBool1GetEmpty(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("delete from hero_bool where id = 1")
	//db.Exec("insert into hero_bool(id,type,bool)values(1,1,true)")

	stmt, err := db.Prepare("select bool from hero_bool where id=? and type=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	id := 1
	t := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var b bool
		stmt.QueryRow(&id, &t).Scan(&b)
	}
}

func BenchmarkBool2Match(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("delete from hero_bool2 where id = 1")
	db.Exec("insert into hero_bool2(id,bools)values(1,0)")

	stmt, err := db.Prepare("update hero_bool2 set bools=? where id=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	id := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stmt.Exec(&i, &id)
	}
}

func BenchmarkBool2Empty(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("delete from hero_bool2 where id = 1")
	//db.Exec("insert into hero_bool2(id,bools)values(1,0)")

	stmt, err := db.Prepare("update hero_bool2 set bools=? where id=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	id := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stmt.Exec(&i, &id)
	}
}

func BenchmarkBool2Empty2(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("delete from hero_bool2 where id = 1")
	db.Exec("insert into hero_bool2(id,bools)values(1,0)")

	stmt, err := db.Prepare("update hero_bool2 set bools=? where id=? and bools=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	id := 1
	bl := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stmt.Exec(&i, &id, &bl)
	}
}

func BenchmarkBool2Get(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("delete from hero_bool2 where id = 1")
	db.Exec("insert into hero_bool2(id,bools)values(1,0)")

	stmt, err := db.Prepare("select bools from hero_bool2 where id=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	id := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var b int
		stmt.QueryRow(&id).Scan(&b)
	}
}

func BenchmarkBool2GetEmpty(b *testing.B) {
	dataSourceName := "root:my-secret-pw@tcp(127.0.0.1:3306)/liwei1"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	db.Exec("delete from hero_bool2 where id = 1")
	db.Exec("insert into hero_bool2(id,bools)values(1,0)")

	stmt, err := db.Prepare("update hero_bool2 set bools=? where id=?")
	if err != nil {
		fmt.Println("创建mysql.Client失败，跳过测试", err.Error())
		return
	}

	id := 1
	t := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var b bool
		stmt.QueryRow(&id, &t).Scan(&b)
	}
}
