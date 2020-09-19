package gamelogs

import "github.com/lightpaw/eventlog"

// LogName: create_hero_log
//
// 创建英雄日志，每个英雄只在第一次创建时候记录
//
// id (int64): 玩家id
// create_time (int64): 英雄创建时间，unix秒
// since_day (uint64): 开服第X天
func CreateHeroLog(pid, sid uint32, id int64, create_time int64, since_day uint64) {
	eventlog.Commit(eventlog.NewEvent("create_hero_log", pid, sid).
		With("id", id).
		With("create_time", create_time).
		With("since_day", since_day))
}

// LogName: hero_online_log
//
// 在线日志，以下3种情况都记录这个日志（state字段区分）
// 1、玩家上线
// 2、玩家离线
// 3、定时每小时的0分0秒，所有在线用户记录一次
//
// id (int64): 玩家id
// online (int64): 在线时长，单位秒
func HeroOnlineLog(pid, sid uint32, id int64, online int64) {
	eventlog.Commit(eventlog.NewEvent("hero_online_log", pid, sid).
		With("id", id).
		With("online", online))
}

// LogName: new_guide_log
//
// 新手引导日志，记录新手引导的过程
//
// id (int64): 玩家id
// progress (int32): 引导进度
// completed (bool): true表示已完成新手引导
func NewGuideLog(pid, sid uint32, id int64, progress int32, completed bool) {
	eventlog.Commit(eventlog.NewEvent("new_guide_log", pid, sid).
		With("id", id).
		With("progress", progress).
		With("completed", completed))
}

// LogName: device_login_log
//
// 设备登陆日志
//
// id (int64): 玩家id
// device_id (string): 设备唯一标识
// os_type (string): 客户端操作系统版本
func DeviceLoginLog(pid, sid uint32, id int64, device_id string, os_type string) {
	eventlog.Commit(eventlog.NewEvent("device_login_log", pid, sid).
		With("id", id).
		With("device_id", device_id).
		With("os_type", os_type))
}

// LogName: recharge_log
//
// 充值订单日志
//
// order_id (string): 订单id
// order_type (uint64): 订单类型（内购类型）
// order_amount (uint64): 订单金额（以平台为准，海外为美元，国内为人民币）
// id (int64): 玩家id
func RechargeLog(pid, sid uint32, order_id string, order_type uint64, order_amount uint64, id int64) {
	eventlog.Commit(eventlog.NewEvent("recharge_log", pid, sid).
		With("order_id", order_id).
		With("order_type", order_type).
		With("order_amount", order_amount).
		With("id", id))
}

// LogName: yuanbao_reduce_log
//
// 元宝进出流水日志
//
// id (int64): 玩家id
// t (uint64): 操作日志
// amount (uint64): 元宝数
func YuanbaoReduceLog(pid, sid uint32, id int64, t uint64, amount uint64) {
	eventlog.Commit(eventlog.NewEvent("yuanbao_reduce_log", pid, sid).
		With("id", id).
		With("t", t).
		With("amount", amount))
}

// LogName: dianquan_reduce_log
//
// 点券进出流水日志
//
// id (int64): 玩家id
// t (uint64): 操作日志
// amount (uint64): 元宝数
func DianquanReduceLog(pid, sid uint32, id int64, t uint64, amount uint64) {
	eventlog.Commit(eventlog.NewEvent("dianquan_reduce_log", pid, sid).
		With("id", id).
		With("t", t).
		With("amount", amount))
}

// LogName: upgrade_hero_level_log
//
// 君主升级日志
//
// id (int64): 玩家id
// level (uint64): 君主等级
func UpgradeHeroLevelLog(pid, sid uint32, id int64, level uint64) {
	eventlog.Commit(eventlog.NewEvent("upgrade_hero_level_log", pid, sid).
		With("id", id).
		With("level", level))
}

// LogName: update_base_level_log
//
// 主城等级日志
//
// id (int64): 玩家id
// level (uint64): 主城等级
func UpdateBaseLevelLog(pid, sid uint32, id int64, level uint64) {
	eventlog.Commit(eventlog.NewEvent("update_base_level_log", pid, sid).
		With("id", id).
		With("level", level))
}

// LogName: complete_task_log
//
// 任务完成日志
//
// id (int64): 玩家id
// t (uint64): 任务类型
// task_id (uint64): 任务id
func CompleteTaskLog(pid, sid uint32, id int64, t uint64, task_id uint64) {
	eventlog.Commit(eventlog.NewEvent("complete_task_log", pid, sid).
		With("id", id).
		With("t", t).
		With("task_id", task_id))
}

// LogName: unlock_captain_log
//
// 武将解锁日志
//
// id (int64): 玩家id
// captain_id (uint64): 武将id
// quality (uint64): 武将品质
func UnlockCaptainLog(pid, sid uint32, id int64, captain_id uint64, quality uint64) {
	eventlog.Commit(eventlog.NewEvent("unlock_captain_log", pid, sid).
		With("id", id).
		With("captain_id", captain_id).
		With("quality", quality))
}

// LogName: upgrade_building_log
//
// 武将解锁日志
//
// id (int64): 玩家id
// t (uint64): 建筑类型
// level (uint64): 建筑等级
func UpgradeBuildingLog(pid, sid uint32, id int64, t uint64, level uint64) {
	eventlog.Commit(eventlog.NewEvent("upgrade_building_log", pid, sid).
		With("id", id).
		With("t", t).
		With("level", level))
}

// LogName: upgrade_tech_log
//
// 武将解锁日志
//
// id (int64): 玩家id
// t (uint64): 科技类型
// level (uint64): 科技等级
func UpgradeTechLog(pid, sid uint32, id int64, t uint64, level uint64) {
	eventlog.Commit(eventlog.NewEvent("upgrade_tech_log", pid, sid).
		With("id", id).
		With("t", t).
		With("level", level))
}

// LogName: tower_challenge_log
//
// 武将解锁日志
//
// id (int64): 玩家id
// floor (uint64): 层数
// pass (bool): 是否通关
func TowerChallengeLog(pid, sid uint32, id int64, floor uint64, pass bool) {
	eventlog.Commit(eventlog.NewEvent("tower_challenge_log", pid, sid).
		With("id", id).
		With("floor", floor).
		With("pass", pass))
}

// LogName: dungeon_challenge_log
//
// 武将解锁日志
//
// id (int64): 玩家id
// chapter (uint64): 推图副本章节
// dungeon_id (uint64): 副本id
// difficult (uint64): 难度 1-普通 2-精英
// star (uint64): 星数
// pass (bool): 是否通关
func DungeonChallengeLog(pid, sid uint32, id int64, chapter uint64, dungeon_id uint64, difficult uint64, star uint64, pass bool) {
	eventlog.Commit(eventlog.NewEvent("dungeon_challenge_log", pid, sid).
		With("id", id).
		With("chapter", chapter).
		With("dungeon_id", dungeon_id).
		With("difficult", difficult).
		With("star", star).
		With("pass", pass))
}

// LogName: server_online_log
//
// 服务器在线日志
//
// count (uint64): 在线人数
func ServerOnlineLog(pid, sid uint32, count uint64) {
	eventlog.Commit(eventlog.NewEvent("server_online_log", pid, sid).
		With("count", count))
}
