package main

import (
	"flag"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/pkg/errors"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"github.com/google/uuid"
	"github.com/lightpaw/male7combatx"
	"time"
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

	config := configDatas.CombatConfig()
	combatConfig, err := combatx.NewConfig(config.Encode())
	if err != nil {
		logrus.WithError(err).Panic("初始化战斗服配置失败")
	}

	proto := &server_proto.CombatXRequestServerProto{}

	proto.Seed = rand.Int63()
	proto.UploadFilePath = uuid.New().String()

	if defender.WallStat != nil {
		proto.MapRes = combatScene.WallMapRes
	} else {
		proto.MapRes = combatScene.MapRes
	}
	proto.MapXLen = u64.Int32(config.CombatXLen)
	proto.MapYLen = u64.Int32(config.CombatYLen)

	proto.AttackerId = int64(attacker.Id)
	proto.DefenserId = int64(defender.Id)

	proto.Attacker = attacker.GetPlayer()
	proto.Defenser = defender.GetPlayer()

	proto.ReturnResult = true

	proto.Debug = true

	uploader := combatx.NewLocalUploader("temp")
	response := combatx.Handle(combatConfig, uploader, proto)

	time.Sleep(time.Second)

	if response.ReturnCode != 0 {
		logrus.Debugln("战斗计算失败: %d, %s", response.ReturnCode, response.ReturnMsg)
		return
	}

	combatx.PrintResult(response.Result)

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
