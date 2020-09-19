package sqldb

import (
	"database/sql"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/u64"
	"strings"
)

var upgrades = []func(*sql.DB, string, uint64){
	// 1号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getTableExist(db, dbName, "bai_zhan", "id") {
			sql := "CREATE TABLE bai_zhan (" +
				"`id` bigint NOT NULL AUTO_INCREMENT," +
				"`attacker_id` bigint," +
				"`defender_id` bigint," +
				"`data` mediumblob NOT NULL," +
				"`time` bigint," +
				"PRIMARY KEY (`id`)," +
				"INDEX `attacker_id` USING BTREE (attacker_id)," +
				"INDEX `defender_id` USING BTREE (defender_id)" +
				");"
			execSql(db, sql, version)
		}
	},
	// 2号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getTableExist(db, dbName, "user", "id") {
			sql := "CREATE TABLE user (" +
				"`id` bigint(20) NOT NULL," +
				"`misc` blob," +
				"PRIMARY KEY (`id`)" +
				") COMMENT='';"
			execSql(db, sql, version)
		}
	},
	// 3号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getIndexExist(db, dbName, "hero", "name") {
			sql := "ALTER TABLE `hero` ADD UNIQUE `name` (`name`);"
			execSql(db, sql, version)
		}
	},
	// 4号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getIndexExist(db, dbName, "guild_logs", "id") {
			sql := "CREATE TABLE `guild_logs` (" +
				"`id` bigint NOT NULL AUTO_INCREMENT," +
				"`guild` bigint(20) NOT NULL," +
				"`type` bigint(20) NOT NULL," +
				"`data` mediumblob NOT NULL," +
				"PRIMARY KEY (`id`)," +
				"INDEX `guild_type_id` USING BTREE (guild, type, id)" +
				");"
			execSql(db, sql, version)
		}
	},
	// 5号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getTableExist(db, dbName, "chat_msg", "id") {
			sql := "CREATE TABLE `chat_msg` (" +
				"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
				"`sender` bigint(20) NOT NULL," +
				"`room` blob NOT NULL," +
				"`time` bigint(20) NOT NULL," +
				"`msg` mediumblob NOT NULL," +
				"PRIMARY KEY (`id`)," +
				"INDEX `room_time` USING BTREE (room(255), time)" +
				");"
			execSql(db, sql, version)
		}
		if !getTableExist(db, dbName, "chat_window", "hero_id") {
			sql := "CREATE TABLE `chat_window` (" +
				"`hero_id` bigint(20) NOT NULL," +
				"`room` blob NOT NULL," +
				"`time` bigint(20) NOT NULL," +
				"`unread_count` int(11) NOT NULL," +
				"`target` mediumblob NOT NULL," +
				"PRIMARY KEY (`hero_id`, `room`(255))," +
				"INDEX `heroid_time` USING BTREE (hero_id, time)" +
				");"
			execSql(db, sql, version)
		}
	},
	// 6号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getIndexExist(db, dbName, "hero", "settings") {
			sql := "ALTER TABLE `hero` ADD COLUMN `settings` bigint(20) DEFAULT 0 AFTER `base_region`;"
			execSql(db, sql, version)
		}
	},
	// 7号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getTableExist(db, dbName, "farm", "hero_id") {
			sql1 := "CREATE TABLE `farm` (" +
				"`hero_id` bigint(20) NOT NULL," +
				"`cube` bigint(20) NOT NULL COMMENT '地块 id'," +
				"`start_time` bigint(20) DEFAULT NULL COMMENT '开始时间'," +
				"`ripe_time` bigint(20) DEFAULT NULL," +
				"`conflict_time` bigint(20) DEFAULT NULL," +
				"`remove_time` bigint(20) DEFAULT NULL," +
				"`res_id` int(11) DEFAULT NULL COMMENT '资源配置 id'," +
				"`steal_times` int(11) DEFAULT NULL," +
				"PRIMARY KEY (`hero_id`,`cube`)," +
				"INDEX `steal` (`hero_id`,`conflict_time`,`remove_time`,`ripe_time`) USING BTREE" +
				") ;"
			sql2 := "CREATE TABLE `farm_log` (" +
				"`hero_id` bigint(20) NOT NULL," +
				"`content` mediumblob NOT NULL," +
				"`log_time` bigint(20) DEFAULT NULL," +
				"INDEX `load` (`hero_id`,`log_time`)  USING BTREE " +
				");"
			sql3 := "CREATE TABLE `farm_steal` (" +
				"`hero_id` bigint(20) NOT NULL," +
				"`thief_id` bigint(20) NOT NULL," +
				"`cube` bigint(20) NOT NULL," +
				"PRIMARY KEY (`hero_id`,`thief_id`,`cube`)" +
				");"

			execSql(db, sql1, version)
			execSql(db, sql2, version)
			execSql(db, sql3, version)
		}
	},
	// 8号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getIndexExist(db, dbName, "hero", "guild_id") {
			sql1 := "ALTER TABLE `hero` ADD COLUMN `guild_id` bigint(20) DEFAULT 0 AFTER `settings`;"
			sql2 := "ALTER TABLE `hero` ADD INDEX `no_guild_hero` (`guild_id`, `name`) ;"
			execSql(db, sql1, version)
			execSql(db, sql2, version)
		}
	},
	// 9号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getIndexExist(db, dbName, "chat_msg", "chat_type") {
			sql1 := "ALTER TABLE `chat_msg` " +
				"ADD COLUMN `chat_type` int(11) NOT NULL AFTER `sender`, " +
				"ADD INDEX `type_time` USING BTREE (chat_type, time);"
			execSql(db, sql1, version)
		}
	},
	// 10号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getTableExist(db, dbName, "xuany_record", "id") {
			sql1 := "CREATE TABLE `xuany_record` (" +
				"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
				"`hero_id` bigint(20) NOT NULL," +
				"`record` mediumblob NOT NULL," +
				"PRIMARY KEY (`id`)," +
				"INDEX `hero_id` (hero_id)" +
				");"
			execSql(db, sql1, version)
		}
	},
	// 11号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getTableExist(db, dbName, "hero", "location") {
			sql1 := "ALTER TABLE `hero` ADD COLUMN `last_online_time` bigint(20) DEFAULT 0 AFTER `settings`;"
			sql2 := "ALTER TABLE `hero` ADD COLUMN `location` int(11) DEFAULT 0 AFTER `last_online_time`;"
			sql3 := "ALTER TABLE `hero` ADD COLUMN `level` int(11) DEFAULT 0 AFTER `location`;"
			sql4 := "ALTER TABLE `hero` ADD INDEX `list_by_location` (`location`, `last_online_time`, `level`) ;"
			sql5 := "ALTER TABLE `hero` ADD INDEX `last_online_time` (`last_online_time`, `level`) ;"
			execSql(db, sql1, version)
			execSql(db, sql2, version)
			execSql(db, sql3, version)
			execSql(db, sql4, version)
			execSql(db, sql5, version)
		}
	},
	// 12号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getTableExist(db, dbName, "hero", "offline_bool") {
			sql1 := "ALTER TABLE `hero` ADD COLUMN `offline_bool` bigint(20) NOT NULL DEFAULT '0' AFTER `guild_id`;"
			execSql(db, sql1, version)
		}
	},

	// 13号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getTableExist(db, dbName, "mc_war_record", "war_id") {
			sql1 := "CREATE TABLE `mc_war_record` (" +
				"`war_id` bigint(20) NOT NULL," +
				"`mc_id` bigint(20) NOT NULL," +
				"`record` mediumblob NOT NULL," +
				"PRIMARY KEY (`war_id`,`mc_id`)" +
				");"
			execSql(db, sql1, version)

			sql2 := "CREATE TABLE `mc_war_hero_record` (" +
				"`war_id` bigint(20) NOT NULL," +
				"`mc_id` bigint(20) NOT NULL," +
				"`hero_id` bigint(20) NOT NULL," +
				"`record` mediumblob NOT NULL," +
				"PRIMARY KEY (`war_id`,`mc_id`, `hero_id`)" +
				");"
			execSql(db, sql2, version)

		}
	},

	//// 14号升级
	//func(db *sql.DB, dbName string, version uint64) {
	//	if !getTableExist(db, dbName, "tlog_tmp", "id") {
	//		sql := "CREATE TABLE `tlog_tmp` (" +
	//			"`id` bigint NOT NULL AUTO_INCREMENT," +
	//			"`msg` blob NOT NULL," +
	//			"`create_time` bigint(20) NOT NULL," +
	//			"PRIMARY KEY (`id`)" +
	//			");"
	//		execSql(db, sql, version)
	//	}
	//},
	// 14号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getIndexExist(db, dbName, "mail", "report_tag") {
			sql := "ALTER TABLE `mail` ADD COLUMN `report_tag` int(11) DEFAULT 0 AFTER `has_report`;"
			execSql(db, sql, version)
		}
	},
	// 15号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getIndexExist(db, dbName, "mc_war_guild_record", "war_id") {
			sql := "CREATE TABLE `mc_war_guild_record` (" +
				"`war_id` bigint(20) NOT NULL," +
				"`mc_id` bigint(20) NOT NULL," +
				"`guild_id` bigint(20) NOT NULL," +
				"`record` mediumblob NOT NULL," +
				"PRIMARY KEY (`war_id`,`mc_id`, `guild_id`)" +
				");"
			execSql(db, sql, version)
		}
	},
	// 16号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getIndexExist(db, dbName, "recharge", "orderid") {
			sql := "CREATE TABLE `recharge` (" +
				"`orderid` varchar(200) NOT NULL," +
				"`orderamount` bigint(20) NOT NULL," +
				"`ordertime` bigint(20) NOT NULL," +
				"`pid` int(11) NOT NULL," +
				"`sid` int(11) NOT NULL," +
				"`heroid` bigint(20) NOT NULL," +
				"`productid` bigint(20) NOT NULL," +
				"`processtime` bigint(20) NOT NULL," +
				"PRIMARY KEY (`orderid`)" +
				");"
			execSql(db, sql, version)
		}
	},
	// 17号升级
	func(db *sql.DB, dbName string, version uint64) {
		if !getIndexExist(db, dbName, "hero", "guild_id") {
			sql1 := "ALTER TABLE `hero` ADD COLUMN `country_id` bigint(20) DEFAULT 0 AFTER `guild_id`;"
			sql2 := "ALTER TABLE `hero` ADD INDEX `country_hero` (`country_id`, `name`) ;"
			execSql(db, sql1, version)
			execSql(db, sql2, version)
		}
	},
}

func upgrade(db *sql.DB, dbName string) {
	// 创建版本表...
	createDbVersionTableIfNotExist(db, dbName)

	// 检测版本号
	version := getDbVersion(db)
	gameVersion := uint64(len(upgrades))
	if version == gameVersion {
		logrus.Infof("数据库版本一致，%v", gameVersion)
		return
	}

	if version > gameVersion {
		if gameVersion == 0 {
			// 更新版本号到0
			logrus.Infof("重置数据库版本0")
			resetDbVersion(db)
			return
		}

		logrus.Panicf("数据库版本超前，服务器版本太低，数据库版本: %v, 服务器版本: %v", version, gameVersion)
	}

	// 开始升级
	for i := version; i < gameVersion; i++ {
		if upgrades[i] == nil {
			logrus.Panicf("数据库版本升级方法没找到，不知道怎么升，index: %v", i)
		}

		newVersion := i + 1
		logrus.Infof("数据库版本升级，版本号: %v", newVersion)
		func() {
			upgrades[i](db, dbName, newVersion)
			setDbVersion(db, newVersion)
		}()
	}
}

func createDbVersionTableIfNotExist(db *sql.DB, dbName string) {

	if getTableExist(db, dbName, "dbversion", "v") {
		return
	}

	for _, s := range strings.Split(initSql, ";") {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			execSql(db, s, 0)
		}
	}
}

func execSql(db *sql.DB, sql string, version uint64) {
	_, err := db.Exec(sql)
	if err != nil {
		logrus.WithField("version", version).WithField("sql", sql).WithError(err).Panicf("数据库升级，执行Sql 失败")
	}
}

func setDbVersion(db *sql.DB, toSet uint64) {
	result, err := db.Exec(`update dbversion set v = ? where v = ?`, toSet, u64.Sub(toSet, 1))
	if err != nil {
		logrus.WithError(err).Panicf("数据库升级，更新DbVersion 失败")
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Panicf("数据库升级，更新DbVersion 失败")
	}

	if rowCount != 1 {
		logrus.WithError(err).Panicf("数据库升级，更新DbVersion 失败，受影响的行数 != 1, row: %v", rowCount)
	}
}

func resetDbVersion(db *sql.DB) {
	result, err := db.Exec(`update dbversion set v = 0`)
	if err != nil {
		logrus.WithError(err).Panicf("数据库升级，更新DbVersion 失败")
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Panicf("数据库升级，更新DbVersion 失败")
	}

	if rowCount != 1 {
		logrus.WithError(err).Panicf("数据库升级，更新DbVersion 失败，受影响的行数 != 1, row: %v", rowCount)
	}
}

func getDbVersion(db *sql.DB) uint64 {
	return getUint64(db, `select v from dbversion`)
}

func getTableExist(db *sql.DB, dbName, tableName, columnName string) bool {
	sql0 := "select count(*) from `INFORMATION_SCHEMA`.`COLUMNS` " +
		"where TABLE_SCHEMA = '%s' and TABLE_NAME = '%s' and COLUMN_NAME = '%s'"
	sql := fmt.Sprintf(sql0, dbName, tableName, columnName)
	return getUint64(db, sql) > 0
}

func getIndexExist(db *sql.DB, dbName, tableName, indexName string) bool {
	sql0 := "select count(*) from `INFORMATION_SCHEMA`.`STATISTICS` " +
		"where TABLE_SCHEMA = '%s' and TABLE_NAME = '%s' and INDEX_NAME = '%s'"

	sql := fmt.Sprintf(sql0, dbName, tableName, indexName)
	return getUint64(db, sql) > 0
}

func getUint64(db *sql.DB, sql string) uint64 {
	row := db.QueryRow(sql)

	var dbVersion uint64
	err := row.Scan(&dbVersion)
	if err != nil {
		logrus.WithField("sql", sql).WithError(err).Panicf("数据库升级，getUint64 error")
	}

	return dbVersion
}
