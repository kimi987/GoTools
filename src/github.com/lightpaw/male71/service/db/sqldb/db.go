package sqldb

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/db/isql"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/service/timeservice"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/pkg/errors"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"
)

func NewMysqlDbService(serverConfig *kv.IndividualServerConfig, datas *config.ConfigDatas, timeService *timeservice.TimeService, register iface.MetricsRegister) (*SqlDbService, error) {
	dataSourceName := serverConfig.DBSN
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, errors.Wrapf(err, "连接数据库失败, %s", dataSourceName)
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "连接数据库失败(Ping), %s", dataSourceName)
	}

	db.SetMaxOpenConns(serverConfig.DBMaxOpenConns)
	db.SetMaxIdleConns(serverConfig.DBMaxIdleConns)
	db.SetConnMaxLifetime(serverConfig.DBConnMaxLifetime)

	dbName := path.Base(dataSourceName)
	if idx := strings.Index(dbName, "?"); idx >= 0 {
		dbName = dbName[:idx]
	}

	return NewSqlDbService(db, dbName, datas, timeService, register)
}

func NewSqlDbService(db *sql.DB, dbName string, datas *config.ConfigDatas, timeService *timeservice.TimeService, register iface.MetricsRegister) (*SqlDbService, error) {

	//db.SetMaxOpenConns(500)

	// 数据库升级
	upgrade(db, dbName)

	// 处理kv表的初始化
	if err := initKv(db); err != nil {
		return nil, err
	}

	var idb isql.DB

	if register != nil && register.EnableDBMetrics() {
		mdb := metrics.NewMetricsDb(db)
		for _, c := range mdb.Collectors() {
			register.Register(c)
		}

		idb = mdb
	} else {
		idb = isql.NewDB(db)
	}

	dbService := &SqlDbService{
		db:            idb,
		datas:         datas,
		timeService:   timeService,
		closeNotifier: make(chan struct{}),
		loopNotifier:  make(chan struct{}),
	}

	go call.CatchLoopPanic(dbService.loop, "SqlDbService.loop()")

	return dbService, nil
}

func initKv(db *sql.DB) error {
	empty := make([]byte, 0)
	for k := range server_proto.Key_value {
		if _, err := db.Exec("insert ignore into kv(k,v)values(?,?)", &k, &empty); err != nil {
			return errors.Wrapf(err, "初始化kv表出错")
		}
	}

	return nil
}

//gogen:iface face DbServiceAdapter DbService service/db/adapter.go CallingTimes
type SqlDbService struct {
	//db *sql.DB
	db          isql.DB
	datas       *config.ConfigDatas
	timeService *timeservice.TimeService

	closeNotifier chan struct{}
	loopNotifier  chan struct{}
}

//func NewDbService(datas *config.ConfigDatas, serverConfig *kv.IndividualServerConfig, timeService *timeservice.TimeService, register iface.MetricsRegister) *SqlDbService {
//	dbsn := serverConfig.DBSN
//	db, err := newMysqlDbService(dbsn, datas, timeService, register)
//	if err != nil {
//		logrus.WithField("dbsn", dbsn).WithError(err).Panic("初始化db service失败")
//	}
//	return db
//}

func (d *SqlDbService) Close() error {
	close(d.closeNotifier)
	<-d.loopNotifier

	if err := d.db.Close(); err != nil {
		logrus.WithError(err).Error("SqlDbService.Close() fail")
	}

	return nil
}

func (d *SqlDbService) loop() {
	defer close(d.loopNotifier)

	ctime := d.timeService.CurrentTime()
	// 凌晨1-3点
	afd := timeutil.DailyTime.NextTime(ctime).Sub(ctime) + time.Hour +
		time.Duration(rand.Int63n(int64(time.Hour)))

	select {
	case <-time.After(afd):
		d.deleteExpireData()
	case <-d.closeNotifier:
		return
	}

	periodTicker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-periodTicker.C:
			d.deleteExpireData()
		case <-d.closeNotifier:
			return
		}
	}
}

func (d *SqlDbService) deleteExpireData() {

	ctime := d.timeService.CurrentTime()

	// 删除聊天信息
	// 世界聊天
	d.tickDeleteExpireChat(shared_proto.ChatType_ChatWorld, ctime, d.datas.MiscConfig().DbWorldChatExpireDuration)

	// 联盟聊天
	d.tickDeleteExpireChat(shared_proto.ChatType_ChatGuild, ctime, d.datas.MiscConfig().DbGuildChatExpireDuration)

	// 私人聊天
	d.tickDeleteExpireChat(shared_proto.ChatType_ChatPrivate, ctime, d.datas.MiscConfig().DbPrivateChatExpireDuration)

	// 删除联盟日志
	d.tickDeleteExpireGuildLog(d.datas.MiscConfig().DbGuildLogCountLimit)

	// 删除邮件
	d.tickDeleteExpireMail(d.datas.MiscConfig().DbMailCountLimit)
}

func (d *SqlDbService) tickDeleteExpireChat(chatType shared_proto.ChatType, ctime time.Time, expireDuration time.Duration) {
	if expireDuration <= 0 {
		return
	}

	expireTime := ctime.Add(-expireDuration).Unix()
	if result, err := d.db.Exec("delete from chat_msg where chat_type = ? and time <= ?", &chatType, &expireTime); err != nil {
		logrus.WithField("chat_type", chatType).WithError(err).Error("定时删除聊天失败")
	} else {
		row, _ := result.RowsAffected()
		logrus.WithField("chat_type", chatType).WithField("row", row).Info("定时删除聊天")
	}
}

func (d *SqlDbService) tickDeleteExpireGuildLog(limit uint64) {
	if limit <= 0 {
		return
	}

	t := shared_proto.GuildLogType_GLT_Memorabilia
	if rows, err := d.db.Query("select guild, type, count(id) n from guild_logs where type<? group by guild, type having n>?", &t, &limit); err != nil {
		logrus.WithError(err).Error("定时删除联盟日志失败")
	} else {
		var guildIds []int64
		var logTypes []int32
		var counts []uint64
		defer rows.Close()
		for rows.Next() {
			var guild int64
			var logType int32
			var count uint64
			err = rows.Scan(&guild, &logType, &count)
			if err != nil {
				err = errors.Wrapf(err, "定时删除联盟日志失败")
				break
			}

			guildIds = append(guildIds, guild)
			logTypes = append(logTypes, logType)
			counts = append(counts, count)
		}

		var totalRow int64
		for i, id := range guildIds {
			logType := logTypes[i]
			count := counts[i]
			toReduceRow := u64.Sub(count, limit)

			if toReduceRow > 0 {
				if result, err := d.db.Exec("delete from guild_logs where guild=? and type=? order by id limit ?", &id, &logType, &toReduceRow); err != nil {
					logrus.WithError(err).Error("定时删除联盟日志失败")
				} else {
					row, _ := result.RowsAffected()
					totalRow += row
				}
			}
		}

		logrus.WithField("row", totalRow).Info("定时删除联盟日志")
	}
}

func (d *SqlDbService) tickDeleteExpireMail(limit uint64) {
	if limit <= 0 {
		return
	}

	// select receiver, count(id) n from mail group by receiver having n>?
	// stmtGetDeleteMailMsg
	if rows, err := d.db.Query("select receiver, count(id) n from mail group by receiver having n>?", &limit); err != nil {
		logrus.WithError(err).Error("定时删除邮件失败")
	} else {
		var heroIds []int64
		var counts []uint64
		defer rows.Close()
		for rows.Next() {
			var heroId int64
			var count uint64
			err = rows.Scan(&heroId, &count)
			if err != nil {
				err = errors.Wrapf(err, "定时删除邮件失败")
				break
			}

			heroIds = append(heroIds, heroId)
			counts = append(counts, count)
		}

		var totalRow, prizeRow int64
		for i, heroId := range heroIds {
			count := counts[i]
			toReduceRow := u64.Sub(count, limit)

			if toReduceRow > 0 {
				// 什么都不管，优先删除一波战报和普通邮件
				// delete from mail where receiver=? and keep=false and (has_prize=false or collected=true) order by id limit ?
				if result, err := d.db.Exec("delete from mail where receiver=? and keep=false and (has_prize=false or collected=true) order by id limit ?", &heroId, &toReduceRow); err != nil {
					logrus.WithError(err).Error("定时删除邮件失败（未收藏）")
				} else {
					row, _ := result.RowsAffected()
					totalRow += row
					toReduceRow = u64.Sub(toReduceRow, uint64(row))
				}
			}

			if toReduceRow > 0 {
				// 还超出上限，开始删除keep的邮件
				// delete from mail where receiver=? and (has_prize=false or collected=true) order by id limit ?
				if result, err := d.db.Exec("delete from mail where receiver=? and (has_prize=false or collected=true) order by id limit ?", &heroId, &toReduceRow); err != nil {
					logrus.WithError(err).Error("定时删除邮件失败（含收藏）")
				} else {
					row, _ := result.RowsAffected()
					totalRow += row
					toReduceRow = u64.Sub(toReduceRow, uint64(row))
				}
			}

			if toReduceRow > 0 {
				// 还超出上限（不管了，一通删）
				// delete from mail where receiver=? order by id limit ?
				if result, err := d.db.Exec("delete from mail where receiver=? order by id limit ?", &heroId, &toReduceRow); err != nil {
					logrus.WithError(err).Error("定时删除邮件失败（含奖励）")
				} else {
					row, _ := result.RowsAffected()
					totalRow += row
					prizeRow += row
					toReduceRow = u64.Sub(toReduceRow, uint64(row))
				}
			}
		}

		logrus.WithField("row", totalRow).WithField("prize_row", prizeRow).Info("定时删除联盟日志")
	}
}

func (d *SqlDbService) LoadKey(ctx context.Context, key server_proto.Key) ([]byte, error) {

	keyName, ok := server_proto.Key_name[int32(key)]
	if !ok {
		return nil, errors.Errorf("db.LoadKey(%v) key不存在", key)
	}

	var v []byte
	err := d.db.QueryRowContext(ctx, "select v from kv where k=?", &keyName).Scan(&v)
	if err != nil {
		return nil, errors.Wrapf(err, "db.LoadKey(%s) 失败", keyName)
	}

	return v, nil
}

func (d *SqlDbService) SaveKey(ctx context.Context, key server_proto.Key, data []byte) error {

	keyName, ok := server_proto.Key_name[int32(key)]
	if !ok {
		return errors.Errorf("db.SaveKey(%v) key不存在", key)
	}

	_, err := d.db.ExecContext(ctx, "update kv set v=? where k=?", &data, &keyName)
	if err != nil {
		return errors.Wrapf(err, "db.SaveKey() fail")
	}

	return nil
}

func (d *SqlDbService) UpdateUserMisc(ctx context.Context, id int64, proto *server_proto.UserMiscProto) error {
	data, err := proto.Marshal()
	if err != nil {
		return err
	}

	_, err = d.db.ExecContext(ctx, "replace into user(id,misc)values(?,?)", &id, &data)
	if err != nil {
		return errors.Wrapf(err, "db.UpdateUserMisc() 失败: %d, %#v", id, proto)
	}

	return nil
}

func (d *SqlDbService) LoadUserMisc(ctx context.Context, id int64) (*server_proto.UserMiscProto, error) {
	var data []byte
	err := d.db.QueryRowContext(ctx, "select misc from user where id=?", &id).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return &server_proto.UserMiscProto{}, nil
		}

		return nil, errors.Wrapf(err, "SqlDbService.LoadUserMisc(%d)", id)
	}

	proto := &server_proto.UserMiscProto{}
	if err := proto.Unmarshal(data); err != nil {
		return nil, err
	}

	return proto, nil
}

func (d *SqlDbService) CreateHero(ctx context.Context, hero *entity.Hero) (bool, error) {

	heroData, err := hero.EncodeServer().Marshal()
	if err != nil {
		return false, errors.Wrapf(err, "CreateHero ServerProto Marshal heroData fail")
	}

	heroId := hero.Id()
	heroName := hero.Name()
	baseRegion := hero.BaseRegion()
	_, err = d.db.ExecContext(ctx, "insert into hero(id,name,hero_data,base_region)values(?,?,?,?)", &heroId, &heroName, &heroData, &baseRegion)
	if err != nil {
		return false, errors.Wrapf(err, "db.CreateHero() 失败")
	}

	return true, nil
}

const settingFieldName = "settings"

func (d *SqlDbService) UpdateSettings(ctx context.Context, id int64, settings uint64) error {
	return d.updateHeroField(ctx, settingFieldName, id, int64(settings))
}

func (d *SqlDbService) FindSettingsOpen(ctx context.Context, settingType shared_proto.SettingType, ids []int64) (result []int64, err error) {
	if len(ids) <= 0 {
		return ids, nil
	}

	return d.loadHeroIdsByBoolFieldValue(ctx, settingFieldName, ids, uint32(settingType), true)
}

func (d *SqlDbService) SaveHero(ctx context.Context, hero *entity.Hero) error {

	heroData, err := hero.EncodeServer().Marshal()
	if err != nil {
		return errors.Wrapf(err, "SaveHero Proto Marshal fail")
	}

	heroId := hero.Id()
	baseRegion := hero.BaseRegion()
	settings := hero.Settings().EncodeToUint64()
	guildId := hero.GuildId()
	countryId := hero.CountryId()
	location := hero.Location()
	lastOnlineTime := timeutil.Marshal64(hero.LastOnlineTime())
	level := hero.Level()

	sql := "update hero set hero_data=?,base_region=?,settings=?, guild_id=?, country_id=?, location=?, last_online_time=?, level=? where id=?"
	_, err = d.db.ExecContext(ctx, sql, &heroData, &baseRegion, &settings, &guildId, &countryId, &location, &lastOnlineTime, &level, &heroId)
	if err != nil {
		return errors.Wrapf(err, "db.SaveHero() 失败")
	}

	return nil
}

func (d *SqlDbService) LoadHero(ctx context.Context, id int64) (*entity.Hero, error) {
	var heroName string

	var heroData []byte
	err := d.db.QueryRowContext(ctx, "select name,hero_data from hero where id=?", &id).Scan(&heroName, &heroData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "SqlDbService.LoadHero(%d)", id)
	}

	return d.parseHero(id, heroName, heroData)
}

func (d *SqlDbService) parseHero(id int64, name string, heroData []byte) (*entity.Hero, error) {
	var heroProto *server_proto.HeroServerProto
	if len(heroData) > 0 {
		heroProto = &server_proto.HeroServerProto{}
		err := heroProto.Unmarshal(heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "DbSerbive.parseHero(%d), Unmarshal HeroServerProto fail", id)
		}
	}

	hero := entity.UnmarshalHero(id, name, d.datas.HeroInitData(), heroProto, d.datas, d.timeService.CurrentTime())

	return hero, nil
}

func (d *SqlDbService) LoadHeroCount(ctx context.Context) (uint64, error) {
	var heroCount uint64
	err := d.db.QueryRowContext(ctx, "select count(id) from hero").Scan(&heroCount)
	if err != nil {
		return 0, errors.Wrapf(err, "db.LoadHeroCount() 失败")
	}

	return heroCount, nil
}

func (d *SqlDbService) LoadAllHeroData(ctx context.Context) ([]*entity.Hero, error) {

	rows, err := d.db.QueryContext(ctx, "select id, name, hero_data from hero")
	if err != nil {
		return nil, errors.Wrapf(err, "db.LoadAllSharedHeroData() 失败")
	}

	array := make([]*entity.Hero, 0)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var heroData []byte
		err := rows.Scan(&id, &name, &heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadAllSharedHeroData() 失败")
		}

		hero, err := d.parseHero(id, name, heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadAllSharedHeroData() 失败")
		}

		array = append(array, hero)
	}

	return array, nil
}

func (d *SqlDbService) LoadAllRegionHero(ctx context.Context) ([]*entity.Hero, error) {

	rows, err := d.db.QueryContext(ctx, "select id, name, hero_data from hero where base_region != 0")
	if err != nil {
		return nil, errors.Wrapf(err, "db.LoadAllSharedHeroData() 失败")
	}

	array := make([]*entity.Hero, 0)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var heroData []byte
		err := rows.Scan(&id, &name, &heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadAllSharedHeroData() 失败")
		}

		hero, err := d.parseHero(id, name, heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadAllSharedHeroData() 失败")
		}

		array = append(array, hero)
	}

	return array, nil
}

func (d *SqlDbService) HeroNameExist(ctx context.Context, name string) (bool, error) {
	heroId, err := d.HeroId(ctx, name)
	if err != nil {
		return false, errors.Wrapf(err, "DB查询英雄名字是否存在失败")
	}

	return heroId != 0, nil

}

func (d *SqlDbService) HeroIds(ctx context.Context) (ids []int64,  err error) {
	var rows isql.Rows
	rows, err = d.db.QueryContext(ctx, "select id from hero")
	if err != nil {
		return nil, errors.Wrapf(err, "DB查询所有玩家id失败")
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, errors.Wrapf(err, "db.HeroIds() Scan 失败")
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (d *SqlDbService) HeroId(ctx context.Context, name string) (heroId int64, err error) {
	err = d.db.QueryRowContext(ctx, "select id from hero where name=?", name).Scan(&heroId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, errors.Wrapf(err, "DB查询英雄名字对应的玩家id失败")
	}

	return heroId, nil
}

func (d *SqlDbService) HeroIdExist(ctx context.Context, heroId int64) (exist bool, err error) {
	count := 0
	err = d.db.QueryRowContext(ctx, "select count(id) from hero where id=?", &heroId).Scan(&count)
	if err != nil {
		err = errors.Wrapf(err, "SqlDbService.HeroIdExist(heroId:%v)", heroId)
		return
	}
	return count > 0, nil
}

func (d *SqlDbService) LoadNoGuildHeroListByName(ctx context.Context, text string, index, size uint64) (heros []*entity.Hero, err error) {
	text = text + "%"
	rows, err := d.db.QueryContext(ctx, "select id, name, hero_data from hero where guild_id = 0 and name like ? limit ?, ?", &text, &index, &size)
	if err != nil {
		return nil, errors.Wrapf(err, "DB LoadNoGuildHeroListByName 失败")
	}

	array := make([]*entity.Hero, 0)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var heroData []byte
		err := rows.Scan(&id, &name, &heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadNoGuildHeroListByName() 失败")
		}

		hero, err := d.parseHero(id, name, heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadNoGuildHeroListByName() 失败")
		}

		array = append(array, hero)
	}

	return array, nil
}

func (d *SqlDbService) UpdateHeroGuildId(ctx context.Context, id, guildId int64) (err error) {
	_, err = d.db.ExecContext(ctx, "update hero set guild_id=? where id=?", &guildId, &id)
	if err != nil {
		logrus.WithError(err).Debugf("英雄更新联盟，DB操作失败")
		return
	}

	return
}

func (d *SqlDbService) LoadHeroListByNameAndCountry(ctx context.Context, text string, countryId, index, size uint64) (heros []*entity.Hero, err error) {
	text = text + "%"
	rows, err := d.db.QueryContext(ctx, "select id, name, hero_data from hero where country_id =? and name like ? limit ?, ?", &countryId, &text, &index, &size)
	if err != nil {
		return nil, errors.Wrapf(err, "DB LoadNoGuildHeroListByName 失败")
	}

	array := make([]*entity.Hero, 0)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var heroData []byte
		err := rows.Scan(&id, &name, &heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadHeroListByNameAndCountry() 失败")
		}

		hero, err := d.parseHero(id, name, heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadHeroListByNameAndCountry() 失败")
		}

		array = append(array, hero)
	}

	return array, nil
}

func (d *SqlDbService) LoadHeroListByCountry(ctx context.Context, countryId, index, size uint64) (heros []*entity.Hero, err error) {
	rows, err := d.db.QueryContext(ctx, "select id, name, hero_data from hero where country_id =? order by last_online_time desc limit ?, ?", &countryId, &index, &size)
	if err != nil {
		return nil, errors.Wrapf(err, "DB LoadNoGuildHeroListByName 失败")
	}

	array := make([]*entity.Hero, 0)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var heroData []byte
		err := rows.Scan(&id, &name, &heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadHeroListByNameAndCountry() 失败")
		}

		hero, err := d.parseHero(id, name, heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadHeroListByNameAndCountry() 失败")
		}

		array = append(array, hero)
	}

	return array, nil
}

func (d *SqlDbService) UpdateHeroName(ctx context.Context, id int64, originName, newName string) bool {

	_, err := d.db.ExecContext(ctx, "update hero set name=? where id=?", &newName, &id)
	if err != nil {
		logrus.WithError(err).Debugf("英雄改名，DB操作失败，名字重复?")
		return false
	}

	return true
}

// mail

func (d *SqlDbService) MaxMailId(ctx context.Context) (uint64, error) {
	var maxId uint64
	err := d.db.QueryRowContext(ctx, "select ifnull(max(id),0) from mail").Scan(&maxId)
	if err != nil {
		return 0, errors.Wrapf(err, "db.MaxMailId() 失败")
	}

	return maxId, nil
}

func mailSqlKey(hasMailId bool, keep, readed, has_report, has_prize, collected int32) int32 {

	key := int32(0)
	if hasMailId {
		key |= 1
	}

	key |= keep << 1
	key |= readed << 3
	key |= has_report << 5
	key |= has_prize << 7
	key |= collected << 9

	return key
}

func validKey(k int32) int32 {
	switch k {
	case 0, 1, 2:
		return k
	default:
		return 0
	}
}

func mailSql(hasMailId bool, keep, readed, has_report, has_prize, collected int32) string {
	b := bytes.Buffer{}
	b.WriteString("select data,keep,readed,collected from mail where ")

	if hasMailId {
		b.WriteString("id<? and ")
	}
	b.WriteString("receiver=?")
	b.WriteString(boolFilterString("keep", keep))
	b.WriteString(boolFilterString("readed", readed))
	b.WriteString(boolFilterString("has_report", has_report))
	b.WriteString(boolFilterString("has_prize", has_prize))
	if has_report == 2 {
		b.WriteString(" and report_tag=?")
	}
	b.WriteString(boolFilterString("collected", collected))

	b.WriteString(" order by id desc limit ?")
	return b.String()
}

func boolFilterString(fieldname string, filter int32) string {

	switch filter {
	case 0:
		return ""
	case 1:
		return " and " + fieldname + "=false"
	case 2:
		return " and " + fieldname + "=true"
	default:
		logrus.Panicf("boolFilterString invalid filter: %v", filter)
		return ""
	}
}

var mailSqlMap = buildMailSqlMap()

func buildMailSqlMap() map[int32]string {
	m := make(map[int32]string)

	for i := 0; i < 2; i++ {
		hasMailId := i == 0
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
							m[mailSqlKey(hasMailId, keep, readed, hasReport, hasPrize, collected)] = mailSql(hasMailId, keep, readed, hasReport, hasPrize, collected)
						}
					}
				}
			}
		}
	}

	return m
}

func (d *SqlDbService) LoadHeroMailList(ctx context.Context, heroId int64, minMailId uint64, keep, readed, has_report, has_prize, collected, report_tag int32, count uint64) ([]*shared_proto.MailProto, error) {

	count = u64.Min(u64.Max(count, d.datas.MiscConfig().MailMinBatchCount), d.datas.MiscConfig().MailMaxBatchCount)

	sqlStr := mailSqlMap[mailSqlKey(minMailId > 0, validKey(keep), validKey(readed), validKey(has_report), validKey(has_prize), validKey(collected))]

	var rows isql.Rows
	var err error
	var args []interface{}
	if minMailId > 0 {
		args = append(args, &minMailId)
	}
	args = append(args, &heroId)
	if has_report == 2 {
		args = append(args, &report_tag)
	}
	args = append(args, &count)

	rows, err = d.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "db.LoadHeroMailList() 失败")
	}

	datas := make([]*shared_proto.MailProto, 0)
	defer rows.Close()
	for rows.Next() {
		var data []byte
		var keep, readed, collected bool

		// data,keep,readed,collected

		if err := rows.Scan(&data, &keep, &readed, &collected); err != nil {
			return nil, errors.Wrapf(err, "db.LoadHeroMailList() 失败")
		}

		mailProto := &shared_proto.MailProto{}
		if err := mailProto.Unmarshal(data); err != nil {
			return nil, errors.Wrapf(err, "db.LoadHeroMailList() MailProto.Unmarshal失败")
		}

		mailProto.Keep = keep
		mailProto.Read = readed
		mailProto.HasPrize = mailProto.Prize != nil
		mailProto.Collected = collected

		datas = append(datas, mailProto)
	}

	return datas, nil
}

func (d *SqlDbService) LoadMailCountHasPrizeNotCollected(ctx context.Context, heroId int64) (count int, err error) {
	row := d.db.QueryRowContext(ctx, "select count(*) from mail where receiver=? and has_prize = true and collected = false", &heroId)

	err = row.Scan(&count)

	return
}

func (d *SqlDbService) LoadMailCountHasReportNotReaded(ctx context.Context, heroId int64, reportTag int32) (count int, err error) {
	row := d.db.QueryRowContext(ctx, "select count(*) from mail where receiver=? and has_report=true and report_tag=? and readed=false", &heroId, &reportTag)

	err = row.Scan(&count)

	return
}

func (d *SqlDbService) LoadMailCountNoReportNotReaded(ctx context.Context, heroId int64) (count int, err error) {
	row := d.db.QueryRowContext(ctx, "select count(*) from mail where receiver=? and has_report = false and readed = false", &heroId)

	err = row.Scan(&count)

	return
}

func (d *SqlDbService) IsCollectableMail(ctx context.Context, id uint64) (bool, error) {
	row := d.db.QueryRowContext(ctx, "select count(*) from mail where id=? and has_prize = true and collected = false", &id)
	var count int
	err := row.Scan(&count)
	return count > 0, err
}

func (d *SqlDbService) LoadCollectMailPrize(ctx context.Context, id uint64, heroId int64) (*resdata.Prize, error) {

	var data []byte
	err := d.db.QueryRowContext(ctx, "select data from mail where id=? and has_prize=true and collected=false", &id).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "db.LoadMailPrize() 失败")
	}

	if len(data) > 0 {
		proto := &shared_proto.MailProto{}
		if err := proto.Unmarshal(data); err != nil {
			return nil, errors.Wrapf(err, "db.LoadCollectMailPrize() MailProto.Unmarshal失败")
		}

		if proto.Prize != nil {
			return resdata.UnmarshalPrize(proto.Prize, d.datas), nil
		}
	}

	return nil, nil
}

func (d *SqlDbService) LoadMail(ctx context.Context, id uint64) ([]byte, error) {

	var data []byte
	err := d.db.QueryRowContext(ctx, "select data from mail where id=?", &id).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "db.LoadMailPrize() 失败")
	}

	return data, nil
}

func (d *SqlDbService) CreateMail(ctx context.Context, id uint64, receiver int64, data []byte, keep, has_report, has_prize bool, report_tag int32, time int64) error {
	//fmt.Println("创建一封新邮件（证明这函数有执行到）@AlbertFan")
	readed := false
	collected := false

	_, err := d.db.ExecContext(ctx, "insert into mail(id,receiver,data,keep,readed,has_report,report_tag,has_prize,collected,time)values(?,?,?,?,?,?,?,?,?,?)", &id, &receiver, &data, &keep, &readed, &has_report, &report_tag, &has_prize, &collected, &time)
	if err != nil {
		return errors.Wrapf(err, "db.CreateMail() 失败")
	}

	return nil
}

func (d *SqlDbService) DeleteMail(ctx context.Context, id uint64, heroId int64) error {

	_, err := d.db.ExecContext(ctx, "delete from mail where id=? and receiver=?", &id, &heroId)
	if err != nil {
		return errors.Wrapf(err, "db.DeleteMail() 失败")
	}

	return nil
}

func (d *SqlDbService) UpdateMailKeep(ctx context.Context, id uint64, heroId int64, keep bool) error {

	_, err := d.db.ExecContext(ctx, "update mail set keep=? where id=? and receiver=?", &keep, &id, &heroId)
	if err != nil {
		return errors.Wrapf(err, "db.UpdateMailKeep() 失败")
	}

	return nil
}

func (d *SqlDbService) UpdateMailRead(ctx context.Context, id uint64, heroId int64, read bool) error {

	_, err := d.db.ExecContext(ctx, "update mail set readed=? where id=? and receiver=?", &read, &id, &heroId)
	if err != nil {
		return errors.Wrapf(err, "db.UpdateMailRead() 失败")
	}

	return nil
}

func (d *SqlDbService) UpdateMailCollected(ctx context.Context, id uint64, heroId int64, collected bool) error {

	_, err := d.db.ExecContext(ctx, "update mail set collected=? where id=? and receiver=?", &collected, &id, &heroId)
	if err != nil {
		return errors.Wrapf(err, "db.UpdateMailCollected() 失败")
	}

	return nil
}

func (d *SqlDbService) ReadMultiMail(ctx context.Context, heroId int64, mailIds []uint64, hasReport bool) (*resdata.Prize, error) {

	if len(mailIds) <= 0 {
		return nil, nil
	}

	minMailId := mailIds[0]
	maxMailId := minMailId
	for _, id := range mailIds {
		maxMailId = u64.Max(maxMailId, id)
		minMailId = u64.Min(minMailId, id)
	}

	var prizeBuilder *resdata.PrizeBuilder
	if !hasReport {
		// 查找范围内，可领取的邮件列表
		rows, err := d.db.QueryContext(ctx, "select data from mail where id>=? and id<=? and receiver=? and has_prize=true and collected=false", &minMailId, &maxMailId, &heroId)
		if err != nil {
			return nil, errors.Wrapf(err, "db.ReadMultiMail() 失败")
		}

		defer rows.Close()
		for rows.Next() {
			var data []byte

			if err := rows.Scan(&data); err != nil {
				return nil, errors.Wrapf(err, "db.ReadMultiMail() 失败")
			}

			mailProto := &shared_proto.MailProto{}
			if err := mailProto.Unmarshal(data); err != nil {
				return nil, errors.Wrapf(err, "db.ReadMultiMail() MailProto.Unmarshal失败")
			}

			if mailProto.Prize != nil {
				if prizeBuilder == nil {
					prizeBuilder = resdata.NewPrizeBuilder()
				}

				prizeBuilder.Add(resdata.UnmarshalPrize(mailProto.Prize, d.datas))
			}
		}
	}

	// 将邮件的read和collected状态改成true
	_, err := d.db.ExecContext(ctx, "update mail set readed=true where id>=? and id<=? and receiver=?", &minMailId, &maxMailId, &heroId)
	if err != nil {
		return nil, errors.Wrapf(err, "db.ReadMultiMail() 失败")
	}

	var prize *resdata.Prize
	if prizeBuilder != nil {
		prize = prizeBuilder.Build()

		_, err = d.db.ExecContext(ctx, "update mail set collected=true where id>=? and id<=? and receiver=? and has_prize=true", &minMailId, &maxMailId, &heroId)
		if err != nil {
			return nil, errors.Wrapf(err, "db.ReadMultiMail() 失败")
		}
	}

	return prize, nil
}

func (d *SqlDbService) DeleteMultiMail(ctx context.Context, heroId int64, mailIds []uint64, hasReport bool) (err error) {
	if len(mailIds) <= 0 {
		return
	}

	minMailId := mailIds[0]
	maxMailId := minMailId
	for _, id := range mailIds {
		maxMailId = u64.Max(maxMailId, id)
		minMailId = u64.Min(minMailId, id)
	}

	if hasReport {
		_, err = d.db.ExecContext(ctx, "delete from mail where id>=? and id<=? and receiver=? and keep=false and readed=true and has_report=true", &minMailId, &maxMailId, &heroId)
	} else {
		_, err = d.db.ExecContext(ctx, "delete from mail where id>=? and id<=? and receiver=? and keep=false and readed=true and has_report=false and (has_prize=false or collected=true)", &minMailId, &maxMailId, &heroId)
	}
	return
}

// guild

func (d *SqlDbService) MaxGuildId(ctx context.Context) (int64, error) {
	var maxId int64
	err := d.db.QueryRowContext(ctx, "select ifnull(max(id),0) from guild").Scan(&maxId)
	if err != nil {
		return 0, errors.Wrapf(err, "db.MaxGuildId() 失败")
	}

	return maxId, nil
}

func (d *SqlDbService) CreateGuild(ctx context.Context, id int64, data []byte) error {

	sql := "insert into guild(id,data)value(?,?)"
	_, err := d.db.ExecContext(ctx, sql, &id, &data)
	if err != nil {
		return errors.Wrapf(err, "db.CreateGuild() 失败")
	}

	return nil
}

func (d *SqlDbService) LoadAllGuild(ctx context.Context) ([]*sharedguilddata.Guild, error) {

	rows, err := d.db.QueryContext(ctx, "select id,data from guild order by id")
	if err != nil {
		return nil, errors.Wrapf(err, "db.LoadAllGuild() 失败")
	}

	array := make([]*sharedguilddata.Guild, 0)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var guildData []byte
		err := rows.Scan(&id, &guildData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadAllSharedHeroData() 失败")
		}

		g, err := d.parseGuild(id, guildData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadAllSharedHeroData() 失败")
		}

		array = append(array, g)
	}

	return array, nil
}

func (d *SqlDbService) parseGuild(id int64, guildData []byte) (*sharedguilddata.Guild, error) {

	if len(guildData) <= 0 {
		return nil, errors.Errorf("DbSerbive.parseGuild(%d), guildData len = 0", id)
	}

	var proto *server_proto.GuildServerProto

	proto = &server_proto.GuildServerProto{}
	err := proto.Unmarshal(guildData)
	if err != nil {
		return nil, errors.Wrapf(err, "DbSerbive.parseGuild(%d), Unmarshal GuildServerProto fail", id)
	}

	return sharedguilddata.UnmarshalGuild(id, proto, d.datas, d.timeService.CurrentTime())
}

func (d *SqlDbService) LoadGuild(ctx context.Context, id int64) (*sharedguilddata.Guild, error) {
	//var name string &name, &flagName,
	//var flagName string
	var data []byte
	err := d.db.QueryRowContext(ctx, "select data from guild where id=?", &id).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "db.LoadGuild() 失败")
	}

	return d.parseGuild(id, data)
}

func (d *SqlDbService) SaveGuild(ctx context.Context, id int64, data []byte) error {

	sql := "update guild set data=? where id=?"
	_, err := d.db.ExecContext(ctx, sql, &data, &id)
	if err != nil {
		return errors.Wrapf(err, "db.SaveGuild() 失败")
	}

	return nil
}

func (d *SqlDbService) DeleteGuild(ctx context.Context, id int64) error {

	sql := "delete from guild where id=?"
	_, err := d.db.ExecContext(ctx, sql, &id)
	if err != nil {
		return errors.Wrapf(err, "db.DeleteGuild() 失败")
	}

	d.deleteGuildLogs(ctx, id)

	return nil
}

func (d *SqlDbService) deleteGuildLogs(ctx context.Context, id int64) error {

	sql := "delete from guild_logs where guild=?"
	_, err := d.db.ExecContext(ctx, sql, &id)
	if err != nil {
		return errors.Wrapf(err, "db.DeleteGuildLogs() 失败")
	}

	return nil
}

func (d *SqlDbService) LoadGuildLogs(ctx context.Context, guildId int64, logType shared_proto.GuildLogType, id int64, count uint64) ([]*shared_proto.GuildLogProto, error) {

	count = u64.Min(u64.Max(count, d.datas.MiscConfig().MailMinBatchCount), d.datas.MiscConfig().MailMaxBatchCount)

	lt := int(logType)

	var rows isql.Rows
	var err error
	if id == 0 {
		sql := "select id,data from guild_logs where guild=? and type=? order by id desc limit ?"
		rows, err = d.db.QueryContext(ctx, sql, &guildId, &lt, &count)
	} else {
		sql := "select id,data from guild_logs where guild=? and type=? and id<? order by id desc limit ?"
		rows, err = d.db.QueryContext(ctx, sql, &guildId, &lt, &id, &count)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "db.LoadGuildLogs() 失败")
	}

	datas := make([]*shared_proto.GuildLogProto, 0)
	defer rows.Close()
	for rows.Next() {

		var logId int64
		var data []byte
		if err := rows.Scan(&logId, &data); err != nil {
			return nil, errors.Wrapf(err, "db.LoadGuildLogs() 失败")
		}

		proto := &shared_proto.GuildLogProto{}
		if err := proto.Unmarshal(data); err != nil {
			return nil, errors.Wrapf(err, "db.LoadGuildLogs() GuildLogProto.Unmarshal失败")
		}

		proto.Id = i64.Int32(logId)

		datas = append(datas, proto)
	}

	return datas, nil
}

func (d *SqlDbService) InsertGuildLog(ctx context.Context, guildId int64, proto *shared_proto.GuildLogProto) error {
	lt := int(proto.Type)
	data := util.SafeMarshal(proto)

	sql := "insert into guild_logs(guild,type,data)value(?,?,?)"
	result, err := d.db.ExecContext(ctx, sql, &guildId, &lt, &data)
	if err != nil {
		return errors.Wrapf(err, "db.AddGuildLog() 失败")
	}

	id, _ := result.LastInsertId()
	proto.Id = i64.Int32(id)

	return nil
}

func (d *SqlDbService) LoadBaiZhanRecord(ctx context.Context, heroId int64, count uint64) (records isql.BytesArray, err error) {
	rows, err := d.db.QueryContext(ctx, "select data from bai_zhan where defender_id = ? OR attacker_id = ? order by time desc limit ?", &heroId, &heroId, &count)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "db.loadBaiZhanRecord() 失败")
	}

	records = make([][]byte, 0)
	defer rows.Close()
	for rows.Next() {
		var data []byte

		if err := rows.Scan(&data); err != nil {
			return nil, errors.Wrapf(err, "db.loadBaiZhanRecord() 失败")
		}

		records = append(records, data)
	}

	return records, nil
}

func (d *SqlDbService) InsertBaiZhanReplay(ctx context.Context, attackerId, defenderId int64, replay *shared_proto.BaiZhanReplayProto, isDefenderNpc bool, time int64) error {
	if isDefenderNpc {
		defenderId = 0
	}

	replayBytes := must.Marshal(replay)

	_, err := d.db.ExecContext(ctx, "insert into bai_zhan(attacker_id,defender_id,data,time)values(?,?,?,?)", &attackerId, &defenderId, &replayBytes, &time)
	if err != nil {
		return errors.Wrapf(err, "db.InsertBaiZhanReplay() 失败")
	}

	return nil
}

func (d *SqlDbService) AddChatMsg(ctx context.Context, senderId int64, room []byte, proto *shared_proto.ChatMsgProto) (int64, error) {
	chatMsgBytes := must.Marshal(proto)

	sql := "insert into chat_msg(sender,chat_type,room,time,msg)value(?,?,?,?,?)"
	result, err := d.db.ExecContext(ctx, sql, &senderId, &proto.ChatType, &room, &proto.SendTime, &chatMsgBytes)
	if err != nil {
		return 0, errors.Wrapf(err, "db.AddChatMsg() 失败")
	}

	id, _ := result.LastInsertId()
	proto.ChatId = i64.ToBytes(id)

	return id, nil
}

func (d *SqlDbService) UpdateChatMsg(ctx context.Context, id int64, proto *shared_proto.ChatMsgProto) bool {
	chatMsgBytes := must.Marshal(proto)

	sqlStr := "update chat_msg set msg=? where id=?"
	_, err := d.db.ExecContext(ctx, sqlStr, &id, &chatMsgBytes)
	if err != nil {
		logrus.WithError(err).Debugf("db.UpdateChatMsg() 失败")
		return false
	}
	return true
}

func (d *SqlDbService) RemoveChatMsg(ctx context.Context, senderId int64) error {

	sql := "delete from chat_msg where sender=?"
	_, err := d.db.ExecContext(ctx, sql, &senderId)
	if err != nil {
		return errors.Wrapf(err, "db.RemoveChatMsg() 失败")
	}

	return nil
}

func (d *SqlDbService) LoadChatMsg(ctx context.Context, id int64) (*shared_proto.ChatMsgProto, error) {
	sqlStr := "select msg from chat_msg where id=?"
	var data []byte
	err := d.db.QueryRowContext(ctx, sqlStr, &id).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	msg := &shared_proto.ChatMsgProto{}
	if err := msg.Unmarshal(data); err != nil {
		return nil, err
	}

	return msg, nil
}

func (d *SqlDbService) ListHeroChatMsg(ctx context.Context, room []byte, minChatId uint64) ([]*shared_proto.ChatMsgProto, error) {
	hasId := minChatId > 0

	var sqlStr string
	var param []interface{}
	if hasId {
		sqlStr = "select id,msg from chat_msg where id<? and room=? order by time desc limit ?"
		param = append(param, &minChatId)
	} else {
		sqlStr = "select id,msg from chat_msg where room=? order by time desc limit ?"
	}
	param = append(param, &room, &d.datas.MiscConfig().ChatBatchCount)

	rows, err := d.db.QueryContext(ctx, sqlStr, param...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "db.ListHeroChatMsg() 失败")
	}

	var msgs []*shared_proto.ChatMsgProto
	defer rows.Close()
	for rows.Next() {
		var id int64
		var data []byte
		if err := rows.Scan(&id, &data); err != nil {
			return nil, errors.Wrapf(err, "db.ListHeroChatMsg() 失败")
		}

		proto := &shared_proto.ChatMsgProto{}
		if err := proto.Unmarshal(data); err != nil {
			return nil, errors.Wrapf(err, "db.ListHeroChatMsg() unmarshal 失败")
		}

		proto.ChatId = i64.ToBytes(id)
		msgs = append(msgs, proto)
	}

	return msgs, nil
}

func (d *SqlDbService) UpdateChatWindow(ctx context.Context, heroId int64, room, targetSenderBytes []byte, addUnread bool, sendTime int32, updateSendTime bool) error {

	var sql string
	var params []interface{}
	if updateSendTime {
		if addUnread {
			sql = "update chat_window set unread_count=unread_count+1,target=?,time=? where hero_id=? and room=?"
		} else {
			sql = "update chat_window set target=?,time=? where hero_id=? and room=?"
		}
		params = []interface{}{&targetSenderBytes, &sendTime, &heroId, &room}
	} else {
		if addUnread {
			sql = "update chat_window set unread_count=unread_count+1,target=? where hero_id=? and room=?"
		} else {
			sql = "update chat_window set target=? where hero_id=? and room=?"
		}
		params = []interface{}{&targetSenderBytes, &heroId, &room}
	}

	result, err := d.db.ExecContext(ctx, sql, params...)
	if err != nil {
		return errors.Wrapf(err, "db.UpdateChatWindow() 失败")
	}

	if row, err := result.RowsAffected(); err != nil {
		return errors.Wrapf(err, "db.UpdateChatWindow() 失败")
	} else if row < 1 {
		// 执行insert操作
		unread := 0
		if addUnread {
			unread = 1
		}
		_, err := d.db.ExecContext(ctx, "insert ignore into chat_window(hero_id,room,time,unread_count,target)values(?,?,?,?,?)", &heroId, &room, &sendTime, &unread, &targetSenderBytes)
		if err != nil {
			return errors.Wrapf(err, "db.InsertChatWindow() 失败")
		}
	}
	return nil
}

func (d *SqlDbService) DeleteChatWindow(ctx context.Context, heroId int64, room []byte) error {

	sql := "delete from chat_window where hero_id=? and room=?"
	_, err := d.db.ExecContext(ctx, sql, &heroId, &room)
	if err != nil {
		return errors.Wrapf(err, "db.DeleteChatWindow() 失败")
	}
	return nil
}

func (d *SqlDbService) LoadUnreadChatCount(ctx context.Context, heroId int64) (uint64, error) {

	var result uint64

	sql := "select ifnull(sum(unread_count),0) from chat_window where hero_id=?"
	err := d.db.QueryRowContext(ctx, sql, &heroId).Scan(&result)
	if err != nil {
		return 0, errors.Wrapf(err, "db.LoadUnreadChatCount() 失败")
	}
	return result, nil
}

func (d *SqlDbService) ReadChat(ctx context.Context, heroId int64, room []byte) error {

	sql := "update chat_window set unread_count=0 where hero_id=? and room=?"
	_, err := d.db.ExecContext(ctx, sql, &heroId, &room)
	if err != nil {
		return errors.Wrapf(err, "db.UpdateChatWindow() 失败")
	}
	return nil
}

func (d *SqlDbService) ListHeroChatWindow(ctx context.Context, heroId int64) ([]uint64, isql.BytesArray, error) {
	sqlStr := "select unread_count,target from chat_window where hero_id=? order by time desc limit ?"
	rows, err := d.db.QueryContext(ctx, sqlStr, &heroId, &d.datas.MiscConfig().ChatWindowLimit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}

		return nil, nil, errors.Wrapf(err, "db.ListHeroChatWindow() 失败")
	}

	var unreadCount []uint64
	var sender [][]byte
	defer rows.Close()
	for rows.Next() {
		var unread uint64
		var data []byte
		if err := rows.Scan(&unread, &data); err != nil {
			return nil, nil, errors.Wrapf(err, "db.ListHeroChatWindow() 失败")
		}

		unreadCount = append(unreadCount, unread)
		sender = append(sender, data)
	}

	return unreadCount, sender, nil
}

func (d *SqlDbService) CreateFarmCube(ctx context.Context, cube *entity.FarmCube) error {
	sqlStr := "INSERT IGNORE INTO farm  (hero_id, cube, start_time, ripe_time, conflict_time, remove_time, res_id, steal_times) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	cube.PrepareSave()
	_, err := d.db.ExecContext(ctx, sqlStr,
		&cube.HeroId, &cube.Cube, &cube.StartTimeInt, &cube.RipeTimeInt,
		&cube.ConflictTimeInt, &cube.RemoveTimeInt, &cube.ResId, &cube.StealTimes)
	if err != nil {
		return errors.Wrapf(err, "db.CreateFarm() 失败")
	}

	return nil
}

func (d *SqlDbService) ResetFarmCubes(ctx context.Context, heroId int64) error {
	sqlStr := "DELETE FROM farm WHERE hero_id=? AND remove_time>0 "
	_, err := d.db.ExecContext(ctx, sqlStr, &heroId)
	if err != nil {
		return errors.Wrap(err, "db.ResetFarmCubes error")
	}

	sqlStr = "DELETE FROM farm WHERE hero_id=? AND conflict_time=0 "
	_, err = d.db.ExecContext(ctx, sqlStr, &heroId)
	if err != nil {
		return errors.Wrap(err, "db.ResetFarmCubes error")
	}
	return nil
}

func (d *SqlDbService) ResetConflictFarmCubes(ctx context.Context, heroId int64) error {
	sqlStr := "UPDATE farm SET start_time=0, ripe_time=0, res_id=0, steal_times=0 WHERE hero_id=? AND conflict_time>0 "
	_, err := d.db.ExecContext(ctx, sqlStr, &heroId)
	if err != nil {
		return errors.Wrapf(err, "db.ResetConflictFarmCubes() 失败")
	}

	return nil
}

func (d *SqlDbService) SaveFarmCube(ctx context.Context, cube *entity.FarmCube) error {
	sqlStr := "UPDATE farm SET start_time=?, ripe_time=?, conflict_time=?, remove_time=?, res_id=?, steal_times=? WHERE hero_id=? AND cube=? "
	cube.PrepareSave()
	_, err := d.db.ExecContext(ctx, sqlStr,
		&cube.StartTimeInt, &cube.RipeTimeInt, &cube.ConflictTimeInt, &cube.RemoveTimeInt, &cube.ResId,
		&cube.StealTimes, &cube.HeroId, &cube.Cube,
	)
	if err != nil {
		return errors.Wrapf(err, "db.SaveFarm() 失败")
	}

	return nil
}

func (d *SqlDbService) PlantFarmCube(ctx context.Context, heroId int64, cube cb.Cube, startTime int64, ripeTime int64, resId uint64) error {
	sqlStr := "UPDATE farm SET start_time=?, ripe_time=?, res_id=?, steal_times=0 WHERE hero_id=? AND cube=? "
	_, err := d.db.ExecContext(ctx, sqlStr, &startTime, &ripeTime, &resId, &heroId, &cube)
	if err != nil {
		return errors.Wrapf(err, "db.PlantFarmCube() 失败")
	}

	return nil
}

func (d *SqlDbService) UpdateFarmCubeState(ctx context.Context, heroId int64, cube cb.Cube, conflictedTime, removeTime int64) error {
	sqlStr := "UPDATE farm SET conflict_time=?, remove_time=? WHERE hero_id=? AND cube=? "
	_, err := d.db.ExecContext(ctx, sqlStr, &conflictedTime, &removeTime, &heroId, &cube)
	if err != nil {
		return errors.Wrapf(err, "db.UpdateFarmCubeState() 失败")
	}
	return nil
}

func (d *SqlDbService) GMFarmRipe(ctx context.Context, heroId int64, startTime, ripeTime int64) error {
	sqlStr := "UPDATE farm SET start_time=?, ripe_time=? WHERE hero_id=? "
	_, err := d.db.ExecContext(ctx, sqlStr, &startTime, &ripeTime, &heroId)
	if err != nil {
		return errors.Wrapf(err, "db.GMFarmRipe() 失败")
	}
	return nil
}

func (d *SqlDbService) SetFarmRipeTime(ctx context.Context, heroId int64, ripeTime int64) error {
	sqlStr := "UPDATE farm SET ripe_time=? WHERE hero_id=? "
	_, err := d.db.ExecContext(ctx, sqlStr, &ripeTime, &heroId)
	if err != nil {
		return errors.Wrapf(err, "db.SetFarmRipeTime() 失败")
	}
	return nil
}

func (d *SqlDbService) UpdateFarmCubeRipeTime(ctx context.Context, heroId int64, cube cb.Cube, startTime, ripeTime int64) error {
	sqlStr := "UPDATE farm SET start_time=?, ripe_time=? WHERE hero_id=? AND cube=? "
	_, err := d.db.ExecContext(ctx, sqlStr, &startTime, &ripeTime, &heroId, &cube)
	if err != nil {
		return errors.Wrapf(err, "db.UpdateFarmCubeRipeTime() 失败")
	}
	return nil
}

func (d *SqlDbService) LoadFarmCube(ctx context.Context, heroId int64, cb cb.Cube) (*entity.FarmCube, error) {
	sqlStr := "SELECT `start_time`, `ripe_time`, `conflict_time`, `remove_time`, `res_id`, `steal_times` FROM `farm` WHERE `hero_id`=? AND `cube`=? "

	cube := &entity.FarmCube{}
	cube.HeroId = heroId
	cube.Cube = cb
	err := d.db.QueryRowContext(ctx, sqlStr, &heroId, &cb).Scan(
		&cube.StartTimeInt, &cube.RipeTimeInt,
		&cube.ConflictTimeInt, &cube.RemoveTimeInt, &cube.ResId,
		&cube.StealTimes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "SqlDbService.LoadFarmCube(%v)", heroId)
	}

	cube.OnLoad(d.datas)
	return cube, nil
}

func rowsToArray(rows isql.Rows, datas *config.ConfigDatas) ([]*entity.FarmCube, error) {
	array := make([]*entity.FarmCube, 0)
	defer rows.Close()
	for rows.Next() {
		cube := &entity.FarmCube{}
		err := rows.Scan(
			&cube.HeroId, &cube.Cube, &cube.StartTimeInt, &cube.RipeTimeInt,
			&cube.ConflictTimeInt, &cube.RemoveTimeInt, &cube.ResId, &cube.StealTimes,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "db.rowsToArray() 失败")
		}

		cube.OnLoad(datas)
		array = append(array, cube)
	}

	return array, nil
}

var (
	sqlLoadStrBase = "SELECT `hero_id`, `cube`, `start_time`, `ripe_time`, `conflict_time`, `remove_time`, `res_id`, `steal_times` FROM `farm` WHERE `hero_id` = ? "
)

func (d *SqlDbService) LoadFarmCubes(ctx context.Context, heroId int64) ([]*entity.FarmCube, error) {
	sqlStr := sqlLoadStrBase

	rows, err := d.db.QueryContext(ctx, sqlStr, &heroId)
	if err != nil {
		return nil, errors.Wrapf(err, "db.LoadFarmCubes 失败")
	}

	return rowsToArray(rows, d.datas)
}

func (d *SqlDbService) LoadFarmHarvestCubes(ctx context.Context, heroId int64) ([]*entity.FarmCube, error) {
	sqlStr := sqlLoadStrBase + " AND `conflict_time`=0;"

	rows, err := d.db.QueryContext(ctx, sqlStr, &heroId)
	if err != nil {
		return nil, errors.Wrapf(err, "db.LoadFarmHarvestCubes 失败")
	}

	return rowsToArray(rows, d.datas)
}

func (d *SqlDbService) LoadFarmStealCubes(ctx context.Context, heroId int64, ripeTime int64, maxStealTimes uint64) ([]*entity.FarmCube, error) {
	sqlStr := sqlLoadStrBase + " AND `conflict_time`=0 AND `remove_time`=0 AND ripe_time<? AND steal_times<? "

	rows, err := d.db.QueryContext(ctx, sqlStr, &heroId, &ripeTime, &maxStealTimes)
	if err != nil {
		return nil, errors.Wrapf(err, "db.LoadFarmStealCubes 失败")
	}

	return rowsToArray(rows, d.datas)
}

func (d *SqlDbService) UpdateFarmStealTimes(ctx context.Context, heroId int64, cubes []cb.Cube) error {
	if len(cubes) <= 0 {
		return errors.Errorf("db.UpdateFarmStealTimes len(cubes) <= 0")
	}

	cubeIds := make([]string, 0)
	for _, c := range cubes {
		cubeIds = append(cubeIds, strconv.FormatUint(uint64(c), 10))
	}
	cubeCondition := " AND `cube` IN (" + strings.Join(cubeIds, ",") + ")"
	sqlStr := "UPDATE `farm` SET `steal_times`=`steal_times`+1 WHERE `hero_id`=? " + cubeCondition
	_, err := d.db.ExecContext(ctx, sqlStr, &heroId)
	if err != nil {
		return errors.Wrapf(err, "db.UpdateFarmStealTimes 失败")
	}

	return nil
}

func (d *SqlDbService) RemoveFarmCube(ctx context.Context, heroId int64, cube cb.Cube) error {
	sqlStr := "DELETE FROM `farm` WHERE `hero_id`=? AND `cube`=? AND `remove_time`>0 "
	_, err := d.db.ExecContext(ctx, sqlStr, &heroId, &cube)
	if err != nil {
		return errors.Wrapf(err, "db.RemoveFarmCube 失败")
	}

	return nil
}

func (d *SqlDbService) CreateFarmLog(ctx context.Context, logProto *shared_proto.FarmStealLogProto) error {
	sqlStr := "INSERT INTO `farm_log` (`hero_id`, `content`, `log_time`) VALUES (?,?,?);"
	heroIdInt, ok := idbytes.ToId(logProto.HeroId)
	if !ok {
		return errors.Errorf("db CreateFarmLog idbytes.ToId error")
	}

	content, err := logProto.Marshal()
	if err != nil {
		return errors.Wrapf(err, "db.CreateFarmLog() 失败")
	}

	_, err = d.db.ExecContext(ctx, sqlStr, &heroIdInt, &content, &logProto.LogTime)
	if err != nil {
		return errors.Wrapf(err, "db.CreateFarmLog() 失败")
	}

	return nil
}

func (d *SqlDbService) LoadFarmLog(ctx context.Context, heroId int64, size uint64) ([]*shared_proto.FarmStealLogProto, error) {
	sqlStr := "SELECT `content` FROM `farm_log` WHERE `hero_id`=? ORDER BY `log_time` DESC LIMIT ? ;"

	rows, err := d.db.QueryContext(ctx, sqlStr, &heroId, &size)
	if err != nil {
		return nil, errors.Wrapf(err, "db.LoadFarmLog 失败")
	}

	array := make([]*shared_proto.FarmStealLogProto, 0)
	defer rows.Close()
	for rows.Next() {
		var content []byte
		err := rows.Scan(&content)
		if err != nil {
			logrus.WithError(err).Debugf("db.LoadFarmLog rowsScan() 失败")
			continue
		}

		log := &shared_proto.FarmStealLogProto{}
		err = log.Unmarshal(content)
		if err != nil {
			logrus.WithError(err).Debugf("db.LoadFarmLog log.Unmarshal() 失败")
			continue
		}

		array = append(array, log)
	}

	return array, nil
}

func (d *SqlDbService) RemoveFarmLog(ctx context.Context, logTime int32) error {
	sqlStr := "DELETE FROM `farm_log` WHERE `log_time`<?;"
	_, err := d.db.ExecContext(ctx, sqlStr, &logTime)
	if err != nil {
		return errors.Wrapf(err, "db.RemoveFarmLog 失败")
	}

	return nil
}

func (d *SqlDbService) AddFarmSteal(ctx context.Context, heroId, thiefId int64, cube cb.Cube) error {
	sqlStr := "INSERT INTO `farm_steal` (`hero_id`, `cube`, `thief_id`) VALUES (?,?,?) "

	_, err := d.db.ExecContext(ctx, sqlStr, &heroId, &cube, &thiefId)
	if err != nil {
		return errors.Wrapf(err, "db.AddFarmSteal() 失败")
	}

	return nil
}

func (d *SqlDbService) LoadFarmStealCount(ctx context.Context, heroId, thiefId int64, cube cb.Cube) (count uint64, err error) {
	sqlStr := "SELECT COUNT(`cube`) `count` FROM `farm_steal` WHERE `hero_id`=? AND `thief_id`=? AND cube=? "
	err = d.db.QueryRowContext(ctx, sqlStr, &heroId, &thiefId, &cube).Scan(&count)
	if err != nil {
		err = errors.Wrapf(err, "db.LoadFarmStealCount() 失败")
		return
	}
	return
}

func (d *SqlDbService) LoadCanStealCube(ctx context.Context, heroId, thiefId, ripeTime int64, maxStealTimes uint64) (array []*entity.FarmCube, err error) {
	sqlStr := "select hero_id, cube, start_time, ripe_time, conflict_time, remove_time, res_id, steal_times from farm f where f.hero_id=? AND f.conflict_time=0 AND f.remove_time=0 AND f.ripe_time<? AND f.steal_times<? and not exists (select 1 from farm_steal fs where fs.hero_id=f.hero_id and thief_id=? and fs.cube=f.cube)"

	rows, err := d.db.QueryContext(ctx, sqlStr, &heroId, &ripeTime, &maxStealTimes, &thiefId)
	if err != nil {
		err = errors.Wrapf(err, "db.LoadFarmStealCount() 失败")
		return
	}
	return rowsToArray(rows, d.datas)
}

func (d *SqlDbService) LoadCanStealCount(ctx context.Context, heroId, thiefId int64, minExpireProtectTime int64, maxStealTime uint64) (count uint64, err error) {
	sqlStr := "select count(cube) from farm f where hero_id=?  AND conflict_time=0 AND remove_time=0 AND ripe_time<? AND steal_times<? and not exists(select 1 from farm_steal fs where fs.hero_id=f.hero_id and thief_id=? and fs.cube=f.cube)"

	err = d.db.QueryRowContext(ctx, sqlStr, &heroId, &minExpireProtectTime, &maxStealTime, &thiefId).Scan(&count)
	if err != nil {
		err = errors.Wrapf(err, "db.LoadCanStealCount() 失败")
		return
	}
	return
}

func (d *SqlDbService) RemoveFarmSteal(ctx context.Context, heroId int64, cubes []cb.Cube) error {
	if len(cubes) <= 0 {
		return errors.Errorf("db.RemoveFarmSteal len(cubes) <= 0")
	}

	cubeIds := make([]string, 0)
	for _, c := range cubes {
		cubeIds = append(cubeIds, fmt.Sprint(c))
	}

	cubeCondition := " AND `cube` IN (" + strings.Join(cubeIds, ",") + ")"
	sqlStr := "DELETE FROM `farm_steal` WHERE `hero_id`=? " + cubeCondition

	_, err := d.db.ExecContext(ctx, sqlStr, &heroId)
	if err != nil {
		return errors.Wrapf(err, "db.RemoveFarmSteal 失败")
	}

	return nil
}

// xuanyuan

func (d *SqlDbService) InsertXuanyRecord(ctx context.Context, heroId int64, record []byte) (int64, error) {

	query := "insert into xuany_record(hero_id,record)values(?,?)"
	result, err := d.db.ExecContext(ctx, query, &heroId, &record)
	if err != nil {
		return 0, errors.Wrapf(err, "db.InsertXuanyRecord 失败")
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrapf(err, "db.InsertXuanyRecord 获取 LastInsertId 失败")
	}

	return insertId, nil
}

func (d *SqlDbService) LoadXuanyRecord(ctx context.Context, heroId int64, id int64, up bool) ([]int64, isql.BytesArray, error) {

	var query string
	var params []interface{}
	if id == 0 {
		if up {
			query = "select id,record from xuany_record where hero_id=? order by id limit ?"
		} else {
			query = "select id,record from xuany_record where hero_id=? order by id desc limit ?"
		}

		params = []interface{}{&heroId, &d.datas.XuanyuanMiscData().RecordBatchCount}
	} else {
		if up {
			query = "select id,record from xuany_record where hero_id=? and id>? order by id limit ?"
		} else {
			query = "select id,record from xuany_record where hero_id=? and id<? order by id desc limit ?"
		}

		params = []interface{}{&heroId, &id, &d.datas.XuanyuanMiscData().RecordBatchCount}
	}

	if rows, err := d.db.QueryContext(ctx, query, params...); err != nil {
		return nil, nil, errors.Wrapf(err, "db.LoadXuanyRecord 失败")
	} else {
		defer rows.Close()

		var ids []int64
		var datas [][]byte
		for rows.Next() {
			var id int64
			var data []byte
			if err := rows.Scan(&id, &data); err != nil {
				return nil, nil, errors.Wrapf(err, "db.LoadXuanyRecord Scan 失败")
			}

			ids = append(ids, id)
			datas = append(datas, data)
		}
		return ids, datas, nil
	}
}

func (d *SqlDbService) LoadRecommendHeros(ctx context.Context, needLocation bool, location, minLevel, page, size uint64, excludeHeroId int64) (heros []*entity.Hero, err error) {
	var rows isql.Rows

	if needLocation {
		sqlStr := "select id, name, hero_data from hero where location = ? and level >= ? order by last_online_time desc limit ?, ?"
		rows, err = d.db.QueryContext(ctx, sqlStr, &location, &minLevel, &page, &size)
	} else {
		sqlStr := "select id, name, hero_data from hero where level >= ? order by last_online_time desc limit ?, ?"
		rows, err = d.db.QueryContext(ctx, sqlStr, &minLevel, &page, &size)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "DB LoadRecommendHeros QueryContext 失败")
	}

	array := make([]*entity.Hero, 0)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var heroData []byte
		err := rows.Scan(&id, &name, &heroData)
		if err != nil {
			logrus.WithError(err).Debugf("db.LoadRecommendHeros() Scan 失败")
			continue
		}

		hero, err := d.parseHero(id, name, heroData)
		if err != nil {
			logrus.WithError(err).Debugf("db.LoadRecommendHeros() parseHero 失败")
			continue
		}

		if hero.Id() == excludeHeroId {
			continue
		}

		array = append(array, hero)
	}

	return array, nil
}

func (d *SqlDbService) LoadHerosByName(ctx context.Context, text string, page, size uint64) (heros []*entity.Hero, err error) {
	text = text + "%"
	sqlStr := "select id, name, hero_data from hero where name like ? order by last_online_time desc limit ?, ?"
	rows, err := d.db.QueryContext(ctx, sqlStr, &text, &page, &size)

	if err != nil {
		return nil, errors.Wrapf(err, "DB LoadHerosByName QueryContext 失败")
	}

	array := make([]*entity.Hero, 0)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var heroData []byte
		err := rows.Scan(&id, &name, &heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadHerosByName() Scan 失败")
		}

		hero, err := d.parseHero(id, name, heroData)
		if err != nil {
			return nil, errors.Wrapf(err, "db.LoadHerosByName() parseHero 失败")
		}

		array = append(array, hero)
	}

	return array, nil
}

// hero int field

func (d *SqlDbService) loadHeroIdsByFieldValue(ctx context.Context, fieldName string, ids []int64, fieldValue int64) ([]int64, error) {
	if len(ids) <= 0 {
		return ids, nil
	}

	query := "select id from hero where id in (" + i64InCond(ids) + ") and " + fieldName + "=?"

	rows, err := d.db.QueryContext(ctx, query, &fieldValue)
	if err != nil {
		return nil, errors.Wrapf(err, "SqlDbService.LoadHeroIdsByFieldValue(%v)", ids)
	}
	defer rows.Close()

	result := make([]int64, 0, 1+len(ids)>>1)
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return nil, errors.Wrapf(err, "SqlDbService.LoadHeroIdsByFieldValue() 失败")
		}

		result = append(result, id)
	}
	return result, nil
}

func (d *SqlDbService) loadHeroField(ctx context.Context, fieldName string, id int64) (int64, error) {

	query := "select " + fieldName + " from hero where id=?"

	var result int64
	err := d.db.QueryRowContext(ctx, query, &id).Scan(&result)
	if err != nil {
		return 0, errors.Wrapf(err, "SqlDbService.LoadHeroField(%v)", id)
	}
	return result, nil
}

func (d *SqlDbService) updateHeroField(ctx context.Context, fieldName string, id int64, fieldValue int64) error {

	query := "update hero set " + fieldName + "=? where id=?"
	_, err := d.db.ExecContext(ctx, query, &fieldValue, &id)
	if err != nil {
		return errors.Wrapf(err, "SqlDbService.updateHeroField(%v)", id)
	}
	return nil
}

func (d *SqlDbService) updateHeroFieldIfExpected(ctx context.Context, fieldName string, id int64, expected, fieldValue int64) (bool, error) {

	query := "update hero set " + fieldName + "=? where id=? and " + fieldName + "=?"
	result, err := d.db.ExecContext(ctx, query, &fieldValue, &id, &expected)
	if err != nil {
		return false, errors.Wrapf(err, "SqlDbService.updateHeroFieldIfExpected(%v)", id)
	}

	if row, err := result.RowsAffected(); err != nil {
		return false, errors.Wrapf(err, "SqlDbService.updateHeroBoolFieldIfExpected(%v)", id)
	} else {
		return row > 0, nil
	}
}

// hero bool field

func boolCondClause(fieldName string, fieldIndex uint32, fieldValue bool) (string, int64) {
	if fieldValue {
		return fieldName + "&?<>0", 1 << fieldIndex
	} else {
		return fieldName + "&?=0", 1 << fieldIndex
	}
}

func boolSetClause(fieldName string, fieldIndex uint32, fieldValue bool) (string, int64) {
	if fieldValue {
		return fieldName + "=" + fieldName + "|?", 1 << fieldIndex
	} else {
		return fieldName + "=" + fieldName + "&?", ^(1 << fieldIndex)
	}
}

func (d *SqlDbService) loadHeroIdsByBoolFieldValue(ctx context.Context, fieldName string, ids []int64, fieldIndex uint32, fieldValue bool) ([]int64, error) {
	if len(ids) <= 0 {
		return ids, nil
	}

	condClause, mask := boolCondClause(fieldName, fieldIndex, fieldValue)
	query := "select id from hero where id in (" + i64InCond(ids) + ") and " + condClause
	rows, err := d.db.QueryContext(ctx, query, &mask)
	if err != nil {
		return nil, errors.Wrapf(err, "SqlDbService.LoadHeroIdsByFieldValue(%v)", ids)
	}
	defer rows.Close()

	result := make([]int64, 0, 1+len(ids)>>1)
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return nil, errors.Wrapf(err, "SqlDbService.LoadHeroIdsByFieldValue() 失败")
		}

		result = append(result, id)
	}
	return result, nil
}

func (d *SqlDbService) loadHeroBoolField(ctx context.Context, fieldName string, id int64, fieldIndex uint32) (bool, error) {

	field, err := d.loadHeroField(ctx, fieldName, id)
	if err != nil {
		return false, err
	}

	return uint64(field)&uint64(1<<fieldIndex) != 0, nil
}

func (d *SqlDbService) updateHeroBoolField(ctx context.Context, fieldName string, id int64, fieldIndex uint32, fieldValue bool) error {

	setClause, mask := boolSetClause(fieldName, fieldIndex, fieldValue)
	query := "update hero set " + setClause + " where id=?"
	_, err := d.db.ExecContext(ctx, query, &mask, &id)
	if err != nil {
		return errors.Wrapf(err, "SqlDbService.updateHeroField(%v)", id)
	}
	return nil
}

func (d *SqlDbService) updateHeroBoolFieldIfExpected(ctx context.Context, fieldName string, id int64, fieldIndex uint32, expected, fieldValue bool) (bool, error) {

	setClause, setMask := boolSetClause(fieldName, fieldIndex, fieldValue)
	condClause, condMask := boolCondClause(fieldName, fieldIndex, expected)

	query := "update hero set " + setClause + " where id=? and " + condClause

	result, err := d.db.ExecContext(ctx, query, &setMask, &id, &condMask)
	if err != nil {
		return false, errors.Wrapf(err, "SqlDbService.updateHeroBoolFieldIfExpected(%v)", id)
	}

	if row, err := result.RowsAffected(); err != nil {
		return false, errors.Wrapf(err, "SqlDbService.updateHeroBoolFieldIfExpected(%v)", id)
	} else {
		return row > 0, nil
	}
}

func i64InCond(ids []int64) string {
	b := &bytes.Buffer{}
	for i, id := range ids {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(strconv.FormatInt(id, 10))
	}
	return b.String()
}

const offlineBoolFieldName = "offline_bool"

func (d *SqlDbService) UpdateHeroOfflineBoolIfExpected(ctx context.Context, id int64, fieldIndex isql.OfflineBool, expected, fieldValue bool) (bool, error) {
	return d.updateHeroBoolFieldIfExpected(ctx, offlineBoolFieldName, id, uint32(fieldIndex), expected, fieldValue)
}

func (d *SqlDbService) AddMcWarRecord(ctx context.Context, mcWarId, mcId uint64, record *shared_proto.McWarFightRecordProto) error {
	sqlStr := "insert ignore into mc_war_record(war_id, mc_id, record) values (?, ?, ?)"
	data, err := record.Marshal()
	if err != nil {
		return errors.Wrapf(err, "db.AddMcWarRecord record.Marshal 失败")
	}
	_, err = d.db.ExecContext(ctx, sqlStr, &mcWarId, &mcId, &data)
	if err != nil {
		return errors.Wrapf(err, "db.AddMcWarRecord db 失败")
	}

	return nil
}

func (d *SqlDbService) AddMcWarHeroRecord(ctx context.Context, mcWarId, mcId uint64, heroId int64, record *shared_proto.McWarTroopAllRecordProto) error {
	sqlStr := "insert ignore into mc_war_hero_record(war_id, mc_id, hero_id, record) values (?, ?, ?, ?)"
	data, err := record.Marshal()
	if err != nil {
		return errors.Wrapf(err, "db.AddMcWarHeroRecord record.Marshal 失败")
	}
	_, err = d.db.ExecContext(ctx, sqlStr, &mcWarId, &mcId, &heroId, &data)
	if err != nil {
		return errors.Wrapf(err, "db.AddMcWarHeroRecord db 失败")
	}

	return nil
}

func (d *SqlDbService) AddMcWarGuildRecord(ctx context.Context, mcWarId, mcId uint64, guildId int64, record *shared_proto.McWarTroopsInfoProto) error {
	sqlStr := "insert ignore into mc_war_guild_record(war_id, mc_id, guild_id, record) values (?, ?, ?, ?)"
	data, err := record.Marshal()
	if err != nil {
		return errors.Wrapf(err, "db.AddMcWarGuildRecord record.Marshal 失败")
	}
	_, err = d.db.ExecContext(ctx, sqlStr, &mcWarId, &mcId, &guildId, &data)
	if err != nil {
		return errors.Wrapf(err, "db.AddMcWarGuildRecord db 失败")
	}

	return nil
}

func (d *SqlDbService) LoadMcWarRecord(ctx context.Context, mcWarId, mcId uint64) (record *shared_proto.McWarFightRecordProto, err error) {
	sqlStr := "select record from mc_war_record where war_id=? and mc_id=?"

	var data []byte
	err = d.db.QueryRowContext(ctx, sqlStr, &mcWarId, &mcId).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		err = errors.Wrapf(err, "SqlDbService.LoadMcWarRecord(war_id:%v, mc_id: %v)", mcWarId, mcId)
		return
	}

	record = &shared_proto.McWarFightRecordProto{}
	err = record.Unmarshal(data)
	if err != nil {
		err = errors.Wrapf(err, "SqlDbService.LoadMcWarRecord McWarFightRecordProto.Unmarshal")
		return
	}
	return
}

func (d *SqlDbService) LoadMcWarHeroRecord(ctx context.Context, mcWarId, mcId uint64, heroId int64) (record *shared_proto.McWarTroopAllRecordProto, err error) {
	sqlStr := "select record from mc_war_hero_record where war_id=? and mc_id=? and hero_id=?"
	var data []byte
	err = d.db.QueryRowContext(ctx, sqlStr, &mcWarId, &mcId, &heroId).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		err = errors.Wrapf(err, "SqlDbService.LoadMcWarHeroRecord(war_id:%v, mc_id: %v, hero_id:%v)", mcWarId, mcId, heroId)
		return
	}

	record = &shared_proto.McWarTroopAllRecordProto{}
	err = record.Unmarshal(data)
	if err != nil {
		err = errors.Wrapf(err, "SqlDbService.LoadMcWarHeroRecord McWarTroopAllRecordProto.Unmarshal")
		return
	}

	return
}

func (d *SqlDbService) LoadMcWarGuildRecord(ctx context.Context, mcWarId, mcId uint64, guildId int64) (troop *shared_proto.McWarTroopsInfoProto, err error) {
	sqlStr := "select record from mc_war_guild_record where war_id=? and mc_id=? and guild_id=?"
	var data []byte
	err = d.db.QueryRowContext(ctx, sqlStr, &mcWarId, &mcId, &guildId).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		err = errors.Wrapf(err, "SqlDbService.LoadMcWarGuildRecord(war_id:%v, mc_id: %v, guild_id:%v)", mcWarId, mcId, guildId)
		return
	}
	troop = &shared_proto.McWarTroopsInfoProto{}
	err = troop.Unmarshal(data)
	if err != nil {
		err = errors.Wrapf(err, "SqlDbService.LoadMcWarGuildRecord McWarTroopsInfoProto.Unmarshal")
		return
	}
	return
}

func (d *SqlDbService) LoadJoinedMcWarId(ctx context.Context, heroId int64) (warMcIdsObj *entity.JoinedMcWarIds, err error) {
	sqlStr := "select war_id, mc_id from mc_war_hero_record where hero_id=?"

	rows, err := d.db.QueryContext(ctx, sqlStr, &heroId)
	if err != nil {
		return nil, errors.Wrapf(err, "SqlDbService.LoadJoinedMcWarId(hero_id:%v) 失败", heroId)
	}
	defer rows.Close()

	warMcIdsObj = entity.NewJoinedMcWarIds()

	for rows.Next() {
		var warId, mcId int32
		if err = rows.Scan(&warId, &mcId); err != nil {
			return nil, errors.Wrapf(err, "SqlDbService.LoadJoinedMcWarId(hero_id:%v) 失败", heroId)
		}
		warMcIdsObj.WarMcIds[warId] = append(warMcIdsObj.WarMcIds[warId], mcId)
	}

	return
}

func (d *SqlDbService) DelMcWarHeroRecord(ctx context.Context, mcWarId int32) (err error) {
	sqlStr := "delete from mc_war_hero_record where war_id <?"
	_, err = d.db.ExecContext(ctx, sqlStr, &mcWarId)
	if err != nil {
		err = errors.Wrapf(err, "SqlDbService.DelMcWarHeroRecord(war_id:%v)", mcWarId)
		return
	}

	return
}

func (d *SqlDbService) DelMcWarHeroRecordWithHeroId(ctx context.Context, mcWarId int32, mcId uint64, heroId int64) (err error) {
	sqlStr := "delete from mc_war_hero_record where war_id <? and mc_id=? and hero_id=?"
	_, err = d.db.ExecContext(ctx, sqlStr, &mcWarId, &mcId, &heroId)
	if err != nil {
		err = errors.Wrapf(err, "SqlDbService.DelMcWarHeroRecordWithHeroId(war_id:%v, mc_id: %v, hero_id:%v)", mcWarId, mcId, heroId)
		return
	}

	return
}

func (d *SqlDbService) OrderExist(ctx context.Context, orderId string) (exist bool, err error) {
	sqlStr := "select count(orderid) from recharge where orderid=?"
	count := 0
	err = d.db.QueryRowContext(ctx, sqlStr, &orderId).Scan(&count)
	if err != nil {
		err = errors.Wrapf(err, "SqlDbService.OrderExist(orderId:%v)", orderId)
		return
	}
	return count > 0, nil
}

func (d *SqlDbService) CreateOrder(ctx context.Context, orderId string, orderAmount uint64, orderTime int64, pid, sid uint32, heroId int64, productId uint64, processTime int64) error {
	sqlStr := "insert into recharge(orderid,orderamount,ordertime,pid,sid,heroid,productid,processtime)values(?,?,?,?,?,?,?,?)"
	_, err := d.db.ExecContext(ctx, sqlStr, &orderId, &orderAmount, &orderTime, &pid, &sid, &heroId, &productId, &processTime)
	if err != nil {
		return errors.Wrapf(err, "SqlDbService.CreateOrder(orderId:%v)", orderId)
	}
	return nil
}
