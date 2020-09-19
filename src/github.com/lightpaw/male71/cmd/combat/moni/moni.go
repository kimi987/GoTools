package main

import (
	"flag"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7combat/combat"
	"github.com/pkg/errors"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"github.com/google/uuid"
)

var (
	attackerId *uint64 = flag.Uint64("aid", 1, "进攻方id")
	defenderId *uint64 = flag.Uint64("did", 2, "防守方id")
	mapRes     *string = flag.String("map_res", "Battle_Field_1", "地图id")
	port       *int    = flag.Int("p", 7888, "端口")
)

// 模拟战斗
func main() {
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	check.PanicNotTrue(*attackerId != *defenderId, "进攻方id[%d]跟防守方id[%d]一样", attackerId, defenderId)

	ip, err := selfIp()
	if err != nil {
		logrus.WithError(err).Debugln("self ip get fail")
		return
	}

	configDatas := config.NewConfigDatas()
	attacker := configDatas.GetMonsterMasterData(*attackerId)
	if attacker == nil {
		logrus.Panicf("进攻方没有找到!%d", *attackerId)
		return
	}

	defender := configDatas.GetMonsterMasterData(*defenderId)
	if defender == nil {
		logrus.Panicf("防守方没有找到!%d", *defenderId)
		return
	}

	combatScene := configDatas.GetCombatScene(*mapRes)
	if combatScene == nil {
		logrus.Panicf("战斗场景没有找到!%s", *mapRes)
		return
	}

	raceDatas := make([]*shared_proto.RaceDataProto, 0, len(configDatas.GetRaceDataArray()))
	for _, rd := range configDatas.GetRaceDataArray() {
		raceDatas = append(raceDatas, rd.Encode())
	}

	uploadFilePath := fmt.Sprintf("http://%d.%d.%d.%d:%d/", ip[0], ip[1], ip[2], ip[3], *port)

	militaryConfig := configDatas.MilitaryConfig()
	proto := &server_proto.CombatRequestServerProto{
		Seed:           rand.Int63(),
		UploadFilePath: uuid.New().String(),
		AttackerId:     int64(attacker.Id),
		DefenserId:     int64(defender.Id),
		Attacker:       attacker.GetPlayer(),
		Defenser:       defender.GetPlayer(),

		MapXLen:  u64.Int32(militaryConfig.CombatXLen),
		MapYLen:  u64.Int32(militaryConfig.CombatYLen),
		MaxRound: u64.Int32(militaryConfig.CombatMaxRound),

		Coef:          i32.MultiF64(10000, militaryConfig.CombatCoef),
		CritRate:      i32.MultiF64(10000, militaryConfig.CombatCritRate),
		RestraintRate: i32.MultiF64(10000, militaryConfig.CombatRestraintRate),

		Races: raceDatas,

		ScorePercent:                militaryConfig.CombatScorePercent,
		MinWallAttackRound:          u64.Int32(militaryConfig.MinWallAttackRound),
		MaxWallAttachFixDamageRound: u64.Int32(militaryConfig.MaxWallAttachFixDamageRound),
		MaxWallBeenHurtPercent:      i32.MultiF64(10000, militaryConfig.MaxWallBeenHurtPercent),
	}
	if proto.Defenser.WallStat != nil {
		proto.MapRes = combatScene.WallMapRes
	} else {
		proto.MapRes = combatScene.MapRes
	}
	response := combat.LocalHandle(proto)

	if response.ReturnCode != 0 {
		logrus.Debugln("战斗计算失败: %d, %s", response.ReturnCode, response.ReturnMsg)
		return
	}

	logrus.Debugf("link: %#v", strings.Replace(response.Link, "{{local}}", uploadFilePath, -1))

	http.Handle("/", http.FileServer(http.Dir("temp")))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func selfIp() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []net.IP

	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && !ip.IsLoopback() {
				ips = append(ips, ip.To4())
			}
		}
	}

	// 找一个非 127.0.0.1, localhost 地址
	for _, ip := range ips {
		if strings.HasPrefix(ip.String(), "192.168.") {
			return ip, nil
		} else if strings.HasPrefix(ip.String(), "168.168.") {
			return ip, nil
		}
	}

	if len(ips) <= 0 {
		return nil, errors.New("没找到本机ip地址")
	}

	return ips[0], nil
}
