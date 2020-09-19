package military

import (
	"github.com/lightpaw/male7/config/confpath"
	"testing"
	"github.com/lightpaw/male7/config"
	. "github.com/onsi/gomega"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/google/uuid"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7combatx"
	"math/rand"
	"fmt"
)

func TestFightx(t *testing.T) {
	RegisterTestingT(t)

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	Ω(err).Should(Succeed())

	configProto := datas.CombatConfig().Encode()

	combatConfig, err := combatx.NewConfig(configProto)
	Ω(err).Should(Succeed())

	fmt.Println(len(datas.GetCaptainDataArray()))

	for i := 0; i < 1; i++ {
		doFightx(datas, combatConfig)
	}
}

func doFightx(datas *config.ConfigDatas, combatConfig *combatx.Config) {
	attackerId, attacker := genFightObj(datas.GetCaptainDataArray(), true)
	defenserId, defenser := genFightObj(datas.GetCaptainDataArray(), false)

	config := datas.CombatConfig()

	proto := &server_proto.CombatXRequestServerProto{}

	proto.Seed = rand.Int63()
	proto.UploadFilePath = uuid.New().String()
	proto.ReturnResult = true
	proto.MapRes = "map"

	proto.MapXLen = u64.Int32(config.CombatXLen)
	proto.MapYLen = u64.Int32(config.CombatYLen)

	proto.AttackerId = attackerId
	proto.DefenserId = defenserId

	proto.Attacker = attacker
	proto.Defenser = defenser

	// 设置race data
	raceDataMap := make(map[shared_proto.Race]*shared_proto.RaceDataProto)
	for _, t := range proto.Attacker.Troops {
		raceDataMap[t.Captain.Race] = datas.RaceConfig().GetProto(t.Captain.Race)
	}
	for _, t := range proto.Defenser.Troops {
		raceDataMap[t.Captain.Race] = datas.RaceConfig().GetProto(t.Captain.Race)
	}

	uploader := combatx.UploaderFunc(func(filename string, data []byte) (link, secondLink string, err error) {
		return
	})
	resp := combatx.Handle(combatConfig, uploader, proto)

	Ω(resp.ReturnCode).Should(BeEquivalentTo(0))
	combatx.PrintResult(resp.Result)
}
