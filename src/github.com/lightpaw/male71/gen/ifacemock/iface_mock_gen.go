package ifacemock

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/coreos/etcd/clientv3"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/activitydata"
	"github.com/lightpaw/male7/config/bai_zhan_data"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/blockdata"
	"github.com/lightpaw/male7/config/body"
	"github.com/lightpaw/male7/config/buffer"
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/charge"
	"github.com/lightpaw/male7/config/combatdata"
	"github.com/lightpaw/male7/config/combine"
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/dianquan"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/config/dungeon"
	"github.com/lightpaw/male7/config/farm"
	"github.com/lightpaw/male7/config/fishing_data"
	"github.com/lightpaw/male7/config/function"
	"github.com/lightpaw/male7/config/gardendata"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/head"
	"github.com/lightpaw/male7/config/hebi"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/location"
	"github.com/lightpaw/male7/config/maildata"
	"github.com/lightpaw/male7/config/military_data"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/promdata"
	"github.com/lightpaw/male7/config/pushdata"
	"github.com/lightpaw/male7/config/pvetroop"
	"github.com/lightpaw/male7/config/question"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/config/random_event"
	"github.com/lightpaw/male7/config/rank_data"
	"github.com/lightpaw/male7/config/red_packet"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/config/season"
	"github.com/lightpaw/male7/config/settings"
	"github.com/lightpaw/male7/config/shop"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/config/spell"
	"github.com/lightpaw/male7/config/strategydata"
	"github.com/lightpaw/male7/config/strongerdata"
	"github.com/lightpaw/male7/config/survey"
	"github.com/lightpaw/male7/config/tag"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/male7/config/teach"
	"github.com/lightpaw/male7/config/towerdata"
	"github.com/lightpaw/male7/config/vip"
	"github.com/lightpaw/male7/config/xiongnu"
	"github.com/lightpaw/male7/config/xuanydata"
	"github.com/lightpaw/male7/config/zhanjiang"
	"github.com/lightpaw/male7/config/zhengwu"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/face"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/module/bai_zhan/bai_zhan_objs"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuface"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuinfo"
	"github.com/lightpaw/male7/pb/rpcpb/game2tss"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/db/isql"
	"github.com/lightpaw/male7/service/extratimesservice/extratimesface"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"github.com/lightpaw/male7/service/tss"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/i64/concurrent"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/rpc7"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"reflect"
	"sync"
	"time"
)

var mockFuncMapLocker = &sync.RWMutex{}
var mockFuncMap = make(map[interface{}]map[uintptr]interface{})

func Mock(obj, funcKey, funcValue interface{}) {
	mockFuncMapLocker.Lock()
	defer mockFuncMapLocker.Unlock()

	funcMap := mockFuncMap[obj]
	if funcMap == nil {
		funcMap = make(map[uintptr]interface{})
		mockFuncMap[obj] = funcMap
	}

	funcMap[getFunctionPointer(funcKey)] = funcValue
}

func getMockFunc(obj, funcKey interface{}) (funcValue interface{}) {
	mockFuncMapLocker.RLock()
	defer mockFuncMapLocker.RUnlock()

	funcMap := mockFuncMap[obj]
	if funcMap == nil {
		return nil
	}
	return funcMap[getFunctionPointer(funcKey)]
}

func getFunctionPointer(i interface{}) uintptr {
	return reflect.ValueOf(i).Pointer()
}

var ActivityModule = &MockActivityModule{}

type MockActivityModule struct{}

func (s *MockActivityModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockActivityModule) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockActivityModule.Close()")
		}
		f()
	}

}
func (s *MockActivityModule) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockActivityModule.OnHeroOnline()")
		}
		f(a0)
	}

}

var AwsService = &MockAwsService{}

type MockAwsService struct{}

func (s *MockAwsService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockAwsService) InitFirehoseEventLog() bool {
	fi := getMockFunc(s, s.InitFirehoseEventLog)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockAwsService.InitFirehoseEventLog()")
		}
		return f()
	}

	return false
}

var BaiZhanModule = &MockBaiZhanModule{}

type MockBaiZhanModule struct{}

func (s *MockBaiZhanModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockBaiZhanModule) Challenge(a0 iface.HeroController) msg.ErrMsg {
	fi := getMockFunc(s, s.Challenge)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController) msg.ErrMsg)
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.Challenge()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockBaiZhanModule) ClearLastJunXian(a0 iface.HeroController) {
	fi := getMockFunc(s, s.ClearLastJunXian)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.ClearLastJunXian()")
		}
		f(a0)
	}

}
func (s *MockBaiZhanModule) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.Close()")
		}
		f()
	}

}
func (s *MockBaiZhanModule) CollectJunXianPrize(a0 int32, a1 iface.HeroController) (bool, msg.ErrMsg) {
	fi := getMockFunc(s, s.CollectJunXianPrize)
	if fi != nil {
		f, ok := fi.(func(int32, iface.HeroController) (bool, msg.ErrMsg))
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.CollectJunXianPrize()")
		}
		return f(a0, a1)
	}

	return false, nil
}
func (s *MockBaiZhanModule) CollectSalary(a0 iface.HeroController) (bool, msg.ErrMsg) {
	fi := getMockFunc(s, s.CollectSalary)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController) (bool, msg.ErrMsg))
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.CollectSalary()")
		}
		return f(a0)
	}

	return false, nil
}
func (s *MockBaiZhanModule) GmResetChallengeTimes(a0 int64) {
	fi := getMockFunc(s, s.GmResetChallengeTimes)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.GmResetChallengeTimes()")
		}
		f(a0)
	}

}
func (s *MockBaiZhanModule) GmResetDaily() {
	fi := getMockFunc(s, s.GmResetDaily)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.GmResetDaily()")
		}
		f()
	}

}
func (s *MockBaiZhanModule) GmSetJunXian(a0 int64, a1 iface.HeroController) {
	fi := getMockFunc(s, s.GmSetJunXian)
	if fi != nil {
		f, ok := fi.(func(int64, iface.HeroController))
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.GmSetJunXian()")
		}
		f(a0, a1)
	}

}
func (s *MockBaiZhanModule) QueryBaiZhanInfo(a0 iface.HeroController) {
	fi := getMockFunc(s, s.QueryBaiZhanInfo)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.QueryBaiZhanInfo()")
		}
		f(a0)
	}

}
func (s *MockBaiZhanModule) RequestRank(a0 bool, a1 uint64, a2 iface.HeroController) {
	fi := getMockFunc(s, s.RequestRank)
	if fi != nil {
		f, ok := fi.(func(bool, uint64, iface.HeroController))
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.RequestRank()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockBaiZhanModule) RequestSelfRank(a0 iface.HeroController) {
	fi := getMockFunc(s, s.RequestSelfRank)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.RequestSelfRank()")
		}
		f(a0)
	}

}
func (s *MockBaiZhanModule) SelfRecord(a0 int32, a1 iface.HeroController) {
	fi := getMockFunc(s, s.SelfRecord)
	if fi != nil {
		f, ok := fi.(func(int32, iface.HeroController))
		if !ok {
			panic("invalid mock func, MockBaiZhanModule.SelfRecord()")
		}
		f(a0, a1)
	}

}

var BaiZhanService = &MockBaiZhanService{}

type MockBaiZhanService struct{}

func (s *MockBaiZhanService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockBaiZhanService) Func(a0 bai_zhan_objs.BaiZhanObjsFunc) {
	fi := getMockFunc(s, s.Func)
	if fi != nil {
		f, ok := fi.(func(bai_zhan_objs.BaiZhanObjsFunc))
		if !ok {
			panic("invalid mock func, MockBaiZhanService.Func()")
		}
		f(a0)
	}

}
func (s *MockBaiZhanService) GetBaiZhanObj(a0 int64) bai_zhan_objs.RHeroBaiZhanObj {
	fi := getMockFunc(s, s.GetBaiZhanObj)
	if fi != nil {
		f, ok := fi.(func(int64) bai_zhan_objs.RHeroBaiZhanObj)
		if !ok {
			panic("invalid mock func, MockBaiZhanService.GetBaiZhanObj()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockBaiZhanService) GetHistoryMaxJunXianLevel(a0 int64) uint64 {
	fi := getMockFunc(s, s.GetHistoryMaxJunXianLevel)
	if fi != nil {
		f, ok := fi.(func(int64) uint64)
		if !ok {
			panic("invalid mock func, MockBaiZhanService.GetHistoryMaxJunXianLevel()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockBaiZhanService) GetJunXianLevel(a0 int64) uint64 {
	fi := getMockFunc(s, s.GetJunXianLevel)
	if fi != nil {
		f, ok := fi.(func(int64) uint64)
		if !ok {
			panic("invalid mock func, MockBaiZhanService.GetJunXianLevel()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockBaiZhanService) GetPoint(a0 int64) uint64 {
	fi := getMockFunc(s, s.GetPoint)
	if fi != nil {
		f, ok := fi.(func(int64) uint64)
		if !ok {
			panic("invalid mock func, MockBaiZhanService.GetPoint()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockBaiZhanService) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockBaiZhanService.OnHeroOnline()")
		}
		f(a0)
	}

}
func (s *MockBaiZhanService) Stop(a0 bai_zhan_objs.BaiZhanObjsFunc) {
	fi := getMockFunc(s, s.Stop)
	if fi != nil {
		f, ok := fi.(func(bai_zhan_objs.BaiZhanObjsFunc))
		if !ok {
			panic("invalid mock func, MockBaiZhanService.Stop()")
		}
		f(a0)
	}

}
func (s *MockBaiZhanService) TimeOutFunc(a0 bai_zhan_objs.BaiZhanObjsFunc) bool {
	fi := getMockFunc(s, s.TimeOutFunc)
	if fi != nil {
		f, ok := fi.(func(bai_zhan_objs.BaiZhanObjsFunc) bool)
		if !ok {
			panic("invalid mock func, MockBaiZhanService.TimeOutFunc()")
		}
		return f(a0)
	}

	return false
}

var BroadcastService = &MockBroadcastService{}

type MockBroadcastService struct{}

func (s *MockBroadcastService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockBroadcastService) Broadcast(a0 string, a1 bool) {
	fi := getMockFunc(s, s.Broadcast)
	if fi != nil {
		f, ok := fi.(func(string, bool))
		if !ok {
			panic("invalid mock func, MockBroadcastService.Broadcast()")
		}
		f(a0, a1)
	}

}
func (s *MockBroadcastService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockBroadcastService.Close()")
		}
		f()
	}

}
func (s *MockBroadcastService) GetCaptainText(a0 *entity.Captain) string {
	fi := getMockFunc(s, s.GetCaptainText)
	if fi != nil {
		f, ok := fi.(func(*entity.Captain) string)
		if !ok {
			panic("invalid mock func, MockBroadcastService.GetCaptainText()")
		}
		return f(a0)
	}

	return ""
}
func (s *MockBroadcastService) GetEquipText(a0 *goods.EquipmentData) string {
	fi := getMockFunc(s, s.GetEquipText)
	if fi != nil {
		f, ok := fi.(func(*goods.EquipmentData) string)
		if !ok {
			panic("invalid mock func, MockBroadcastService.GetEquipText()")
		}
		return f(a0)
	}

	return ""
}

var BuffService = &MockBuffService{}

type MockBuffService struct{}

func (s *MockBuffService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

// 新增或替换buff
func (s *MockBuffService) AddBuffToSelf(a0 *data.BuffEffectData, a1 int64) bool {
	fi := getMockFunc(s, s.AddBuffToSelf)
	if fi != nil {
		f, ok := fi.(func(*data.BuffEffectData, int64) bool)
		if !ok {
			panic("invalid mock func, MockBuffService.AddBuffToSelf()")
		}
		return f(a0, a1)
	}

	return false
}

// 新增或替换buff
func (s *MockBuffService) AddBuffToTarget(a0 *data.BuffEffectData, a1 int64, a2 int64) bool {
	fi := getMockFunc(s, s.AddBuffToTarget)
	if fi != nil {
		f, ok := fi.(func(*data.BuffEffectData, int64, int64) bool)
		if !ok {
			panic("invalid mock func, MockBuffService.AddBuffToTarget()")
		}
		return f(a0, a1, a2)
	}

	return false
}
func (s *MockBuffService) Cancel(a0 int64, a1 []*entity.BuffInfo) bool {
	fi := getMockFunc(s, s.Cancel)
	if fi != nil {
		f, ok := fi.(func(int64, []*entity.BuffInfo) bool)
		if !ok {
			panic("invalid mock func, MockBuffService.Cancel()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockBuffService) CancelGroup(a0 int64, a1 uint64) bool {
	fi := getMockFunc(s, s.CancelGroup)
	if fi != nil {
		f, ok := fi.(func(int64, uint64) bool)
		if !ok {
			panic("invalid mock func, MockBuffService.CancelGroup()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockBuffService) UpdatePerSecond(a0 *entity.Hero, a1 herolock.LockResult, a2 *entity.BuffInfo) {
	fi := getMockFunc(s, s.UpdatePerSecond)
	if fi != nil {
		f, ok := fi.(func(*entity.Hero, herolock.LockResult, *entity.BuffInfo))
		if !ok {
			panic("invalid mock func, MockBuffService.UpdatePerSecond()")
		}
		f(a0, a1, a2)
	}

}

var ChatModule = &MockChatModule{}

type MockChatModule struct{}

func (s *MockChatModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var ChatService = &MockChatService{}

type MockChatService struct{}

func (s *MockChatService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockChatService) AddMsg(a0 *shared_proto.ChatMsgProto) {
	fi := getMockFunc(s, s.AddMsg)
	if fi != nil {
		f, ok := fi.(func(*shared_proto.ChatMsgProto))
		if !ok {
			panic("invalid mock func, MockChatService.AddMsg()")
		}
		f(a0)
	}

}
func (s *MockChatService) BroadcastSystemChat(a0 string) {
	fi := getMockFunc(s, s.BroadcastSystemChat)
	if fi != nil {
		f, ok := fi.(func(string))
		if !ok {
			panic("invalid mock func, MockChatService.BroadcastSystemChat()")
		}
		f(a0)
	}

}
func (s *MockChatService) GetCacheMsg() pbutil.Buffer {
	fi := getMockFunc(s, s.GetCacheMsg)
	if fi != nil {
		f, ok := fi.(func() pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockChatService.GetCacheMsg()")
		}
		return f()
	}

	return nil
}
func (s *MockChatService) GetChatSender(a0 int64) *shared_proto.ChatSenderProto {
	fi := getMockFunc(s, s.GetChatSender)
	if fi != nil {
		f, ok := fi.(func(int64) *shared_proto.ChatSenderProto)
		if !ok {
			panic("invalid mock func, MockChatService.GetChatSender()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockChatService) GetSystemChatRecord(a0 int64) pbutil.Buffer {
	fi := getMockFunc(s, s.GetSystemChatRecord)
	if fi != nil {
		f, ok := fi.(func(int64) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockChatService.GetSystemChatRecord()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockChatService) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockChatService.OnHeroOnline()")
		}
		f(a0)
	}

}
func (s *MockChatService) SaveDB(a0 int64, a1 int64, a2 shared_proto.ChatType, a3 []byte, a4 []byte, a5 *shared_proto.ChatMsgProto) int64 {
	fi := getMockFunc(s, s.SaveDB)
	if fi != nil {
		f, ok := fi.(func(int64, int64, shared_proto.ChatType, []byte, []byte, *shared_proto.ChatMsgProto) int64)
		if !ok {
			panic("invalid mock func, MockChatService.SaveDB()")
		}
		return f(a0, a1, a2, a3, a4, a5)
	}

	return 0
}

// 系统自动聊天，有DB操作
func (s *MockChatService) SysChat(a0 int64, a1 int64, a2 shared_proto.ChatType, a3 string, a4 shared_proto.ChatMsgType, a5 bool, a6 bool, a7 bool, a8 bool) {
	fi := getMockFunc(s, s.SysChat)
	if fi != nil {
		f, ok := fi.(func(int64, int64, shared_proto.ChatType, string, shared_proto.ChatMsgType, bool, bool, bool, bool))
		if !ok {
			panic("invalid mock func, MockChatService.SysChat()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8)
	}

}
func (s *MockChatService) SysChatFunc(a0 int64, a1 int64, a2 shared_proto.ChatType, a3 string, a4 shared_proto.ChatMsgType, a5 bool, a6 bool, a7 bool, a8 bool, a9 constants.ChatFunc, a10 constants.ChatFunc) int64 {
	fi := getMockFunc(s, s.SysChatFunc)
	if fi != nil {
		f, ok := fi.(func(int64, int64, shared_proto.ChatType, string, shared_proto.ChatMsgType, bool, bool, bool, bool, constants.ChatFunc, constants.ChatFunc) int64)
		if !ok {
			panic("invalid mock func, MockChatService.SysChatFunc()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
	}

	return 0
}
func (s *MockChatService) SysChatProtoFunc(a0 int64, a1 int64, a2 shared_proto.ChatType, a3 string, a4 shared_proto.ChatMsgType, a5 bool, a6 bool, a7 bool, a8 bool, a9 constants.ChatFunc) int64 {
	fi := getMockFunc(s, s.SysChatProtoFunc)
	if fi != nil {
		f, ok := fi.(func(int64, int64, shared_proto.ChatType, string, shared_proto.ChatMsgType, bool, bool, bool, bool, constants.ChatFunc) int64)
		if !ok {
			panic("invalid mock func, MockChatService.SysChatProtoFunc()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	}

	return 0
}
func (s *MockChatService) SysChatSendFunc(a0 int64, a1 int64, a2 shared_proto.ChatType, a3 string, a4 shared_proto.ChatMsgType, a5 bool, a6 bool, a7 bool, a8 bool, a9 constants.ChatFunc) int64 {
	fi := getMockFunc(s, s.SysChatSendFunc)
	if fi != nil {
		f, ok := fi.(func(int64, int64, shared_proto.ChatType, string, shared_proto.ChatMsgType, bool, bool, bool, bool, constants.ChatFunc) int64)
		if !ok {
			panic("invalid mock func, MockChatService.SysChatSendFunc()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	}

	return 0
}
func (s *MockChatService) UpdateDBRedPacket(a0 int64, a1 bool) bool {
	fi := getMockFunc(s, s.UpdateDBRedPacket)
	if fi != nil {
		f, ok := fi.(func(int64, bool) bool)
		if !ok {
			panic("invalid mock func, MockChatService.UpdateDBRedPacket()")
		}
		return f(a0, a1)
	}

	return false
}

var ClientConfigModule = &MockClientConfigModule{}

type MockClientConfigModule struct{}

func (s *MockClientConfigModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var ClusterService = &MockClusterService{}

type MockClusterService struct{}

func (s *MockClusterService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockClusterService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockClusterService.Close()")
		}
		f()
	}

}
func (s *MockClusterService) EtcdClient() *clientv3.Client {
	fi := getMockFunc(s, s.EtcdClient)
	if fi != nil {
		f, ok := fi.(func() *clientv3.Client)
		if !ok {
			panic("invalid mock func, MockClusterService.EtcdClient()")
		}
		return f()
	}

	return nil
}
func (s *MockClusterService) GMUpdateClientVersion(a0 string) {
	fi := getMockFunc(s, s.GMUpdateClientVersion)
	if fi != nil {
		f, ok := fi.(func(string))
		if !ok {
			panic("invalid mock func, MockClusterService.GMUpdateClientVersion()")
		}
		f(a0)
	}

}
func (s *MockClusterService) GetClientVersionMsg(a0 string, a1 string) pbutil.Buffer {
	fi := getMockFunc(s, s.GetClientVersionMsg)
	if fi != nil {
		f, ok := fi.(func(string, string) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockClusterService.GetClientVersionMsg()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockClusterService) GetConfig(a0 string) ([]byte, error) {
	fi := getMockFunc(s, s.GetConfig)
	if fi != nil {
		f, ok := fi.(func(string) ([]byte, error))
		if !ok {
			panic("invalid mock func, MockClusterService.GetConfig()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockClusterService) LoginClient() *rpc7.Client {
	fi := getMockFunc(s, s.LoginClient)
	if fi != nil {
		f, ok := fi.(func() *rpc7.Client)
		if !ok {
			panic("invalid mock func, MockClusterService.LoginClient()")
		}
		return f()
	}

	return nil
}
func (s *MockClusterService) RpcAddr() string {
	fi := getMockFunc(s, s.RpcAddr)
	if fi != nil {
		f, ok := fi.(func() string)
		if !ok {
			panic("invalid mock func, MockClusterService.RpcAddr()")
		}
		return f()
	}

	return ""
}

var ConfigDatas = &MockConfigDatas{}

type MockConfigDatas struct{}

func (s *MockConfigDatas) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockConfigDatas) AchieveTaskData() *config.AchieveTaskDataConfig {
	fi := getMockFunc(s, s.AchieveTaskData)
	if fi != nil {
		f, ok := fi.(func() *config.AchieveTaskDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.AchieveTaskData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) AchieveTaskStarPrizeData() *config.AchieveTaskStarPrizeDataConfig {
	fi := getMockFunc(s, s.AchieveTaskStarPrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.AchieveTaskStarPrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.AchieveTaskStarPrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ActiveDegreePrizeData() *config.ActiveDegreePrizeDataConfig {
	fi := getMockFunc(s, s.ActiveDegreePrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.ActiveDegreePrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ActiveDegreePrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ActiveDegreeTaskData() *config.ActiveDegreeTaskDataConfig {
	fi := getMockFunc(s, s.ActiveDegreeTaskData)
	if fi != nil {
		f, ok := fi.(func() *config.ActiveDegreeTaskDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ActiveDegreeTaskData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ActivityCollectionData() *config.ActivityCollectionDataConfig {
	fi := getMockFunc(s, s.ActivityCollectionData)
	if fi != nil {
		f, ok := fi.(func() *config.ActivityCollectionDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ActivityCollectionData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ActivityShowData() *config.ActivityShowDataConfig {
	fi := getMockFunc(s, s.ActivityShowData)
	if fi != nil {
		f, ok := fi.(func() *config.ActivityShowDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ActivityShowData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ActivityTaskData() *config.ActivityTaskDataConfig {
	fi := getMockFunc(s, s.ActivityTaskData)
	if fi != nil {
		f, ok := fi.(func() *config.ActivityTaskDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ActivityTaskData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ActivityTaskListModeData() *config.ActivityTaskListModeDataConfig {
	fi := getMockFunc(s, s.ActivityTaskListModeData)
	if fi != nil {
		f, ok := fi.(func() *config.ActivityTaskListModeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ActivityTaskListModeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) AmountShowSortData() *config.AmountShowSortDataConfig {
	fi := getMockFunc(s, s.AmountShowSortData)
	if fi != nil {
		f, ok := fi.(func() *config.AmountShowSortDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.AmountShowSortData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) AreaData() *config.AreaDataConfig {
	fi := getMockFunc(s, s.AreaData)
	if fi != nil {
		f, ok := fi.(func() *config.AreaDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.AreaData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) AssemblyData() *config.AssemblyDataConfig {
	fi := getMockFunc(s, s.AssemblyData)
	if fi != nil {
		f, ok := fi.(func() *config.AssemblyDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.AssemblyData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BaYeStageData() *config.BaYeStageDataConfig {
	fi := getMockFunc(s, s.BaYeStageData)
	if fi != nil {
		f, ok := fi.(func() *config.BaYeStageDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BaYeStageData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BaYeTaskData() *config.BaYeTaskDataConfig {
	fi := getMockFunc(s, s.BaYeTaskData)
	if fi != nil {
		f, ok := fi.(func() *config.BaYeTaskDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BaYeTaskData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BaiZhanMiscData() *bai_zhan_data.BaiZhanMiscData {
	fi := getMockFunc(s, s.BaiZhanMiscData)
	if fi != nil {
		f, ok := fi.(func() *bai_zhan_data.BaiZhanMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BaiZhanMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BaowuData() *config.BaowuDataConfig {
	fi := getMockFunc(s, s.BaowuData)
	if fi != nil {
		f, ok := fi.(func() *config.BaowuDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BaowuData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BaozNpcData() *config.BaozNpcDataConfig {
	fi := getMockFunc(s, s.BaozNpcData)
	if fi != nil {
		f, ok := fi.(func() *config.BaozNpcDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BaozNpcData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BaseLevelData() *config.BaseLevelDataConfig {
	fi := getMockFunc(s, s.BaseLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.BaseLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BaseLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BlackMarketData() *config.BlackMarketDataConfig {
	fi := getMockFunc(s, s.BlackMarketData)
	if fi != nil {
		f, ok := fi.(func() *config.BlackMarketDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BlackMarketData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BlackMarketGoodsData() *config.BlackMarketGoodsDataConfig {
	fi := getMockFunc(s, s.BlackMarketGoodsData)
	if fi != nil {
		f, ok := fi.(func() *config.BlackMarketGoodsDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BlackMarketGoodsData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BlackMarketGoodsGroupData() *config.BlackMarketGoodsGroupDataConfig {
	fi := getMockFunc(s, s.BlackMarketGoodsGroupData)
	if fi != nil {
		f, ok := fi.(func() *config.BlackMarketGoodsGroupDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BlackMarketGoodsGroupData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BlockData() *config.BlockDataConfig {
	fi := getMockFunc(s, s.BlockData)
	if fi != nil {
		f, ok := fi.(func() *config.BlockDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BlockData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BodyData() *config.BodyDataConfig {
	fi := getMockFunc(s, s.BodyData)
	if fi != nil {
		f, ok := fi.(func() *config.BodyDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BodyData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BranchTaskData() *config.BranchTaskDataConfig {
	fi := getMockFunc(s, s.BranchTaskData)
	if fi != nil {
		f, ok := fi.(func() *config.BranchTaskDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BranchTaskData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BroadcastData() *config.BroadcastDataConfig {
	fi := getMockFunc(s, s.BroadcastData)
	if fi != nil {
		f, ok := fi.(func() *config.BroadcastDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BroadcastData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BroadcastHelp() *data.BroadcastHelp {
	fi := getMockFunc(s, s.BroadcastHelp)
	if fi != nil {
		f, ok := fi.(func() *data.BroadcastHelp)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BroadcastHelp()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BuffEffectData() *config.BuffEffectDataConfig {
	fi := getMockFunc(s, s.BuffEffectData)
	if fi != nil {
		f, ok := fi.(func() *config.BuffEffectDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BuffEffectData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BufferData() *config.BufferDataConfig {
	fi := getMockFunc(s, s.BufferData)
	if fi != nil {
		f, ok := fi.(func() *config.BufferDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BufferData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BufferTypeData() *config.BufferTypeDataConfig {
	fi := getMockFunc(s, s.BufferTypeData)
	if fi != nil {
		f, ok := fi.(func() *config.BufferTypeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BufferTypeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BuildingData() *config.BuildingDataConfig {
	fi := getMockFunc(s, s.BuildingData)
	if fi != nil {
		f, ok := fi.(func() *config.BuildingDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BuildingData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BuildingEffectData() *config.BuildingEffectDataConfig {
	fi := getMockFunc(s, s.BuildingEffectData)
	if fi != nil {
		f, ok := fi.(func() *config.BuildingEffectDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BuildingEffectData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BuildingLayoutData() *config.BuildingLayoutDataConfig {
	fi := getMockFunc(s, s.BuildingLayoutData)
	if fi != nil {
		f, ok := fi.(func() *config.BuildingLayoutDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BuildingLayoutData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BuildingLayoutMiscData() *domestic_data.BuildingLayoutMiscData {
	fi := getMockFunc(s, s.BuildingLayoutMiscData)
	if fi != nil {
		f, ok := fi.(func() *domestic_data.BuildingLayoutMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BuildingLayoutMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BuildingUnlockData() *config.BuildingUnlockDataConfig {
	fi := getMockFunc(s, s.BuildingUnlockData)
	if fi != nil {
		f, ok := fi.(func() *config.BuildingUnlockDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BuildingUnlockData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BwzlPrizeData() *config.BwzlPrizeDataConfig {
	fi := getMockFunc(s, s.BwzlPrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.BwzlPrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BwzlPrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) BwzlTaskData() *config.BwzlTaskDataConfig {
	fi := getMockFunc(s, s.BwzlTaskData)
	if fi != nil {
		f, ok := fi.(func() *config.BwzlTaskDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.BwzlTaskData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CaptainAbilityData() *config.CaptainAbilityDataConfig {
	fi := getMockFunc(s, s.CaptainAbilityData)
	if fi != nil {
		f, ok := fi.(func() *config.CaptainAbilityDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CaptainAbilityData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CaptainData() *config.CaptainDataConfig {
	fi := getMockFunc(s, s.CaptainData)
	if fi != nil {
		f, ok := fi.(func() *config.CaptainDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CaptainData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CaptainFriendshipData() *config.CaptainFriendshipDataConfig {
	fi := getMockFunc(s, s.CaptainFriendshipData)
	if fi != nil {
		f, ok := fi.(func() *config.CaptainFriendshipDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CaptainFriendshipData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CaptainLevelData() *config.CaptainLevelDataConfig {
	fi := getMockFunc(s, s.CaptainLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.CaptainLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CaptainLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CaptainOfficialCountData() *config.CaptainOfficialCountDataConfig {
	fi := getMockFunc(s, s.CaptainOfficialCountData)
	if fi != nil {
		f, ok := fi.(func() *config.CaptainOfficialCountDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CaptainOfficialCountData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CaptainOfficialData() *config.CaptainOfficialDataConfig {
	fi := getMockFunc(s, s.CaptainOfficialData)
	if fi != nil {
		f, ok := fi.(func() *config.CaptainOfficialDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CaptainOfficialData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CaptainRarityData() *config.CaptainRarityDataConfig {
	fi := getMockFunc(s, s.CaptainRarityData)
	if fi != nil {
		f, ok := fi.(func() *config.CaptainRarityDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CaptainRarityData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CaptainRebirthLevelData() *config.CaptainRebirthLevelDataConfig {
	fi := getMockFunc(s, s.CaptainRebirthLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.CaptainRebirthLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CaptainRebirthLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CaptainStarData() *config.CaptainStarDataConfig {
	fi := getMockFunc(s, s.CaptainStarData)
	if fi != nil {
		f, ok := fi.(func() *config.CaptainStarDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CaptainStarData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ChargeObjData() *config.ChargeObjDataConfig {
	fi := getMockFunc(s, s.ChargeObjData)
	if fi != nil {
		f, ok := fi.(func() *config.ChargeObjDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ChargeObjData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ChargePrizeData() *config.ChargePrizeDataConfig {
	fi := getMockFunc(s, s.ChargePrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.ChargePrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ChargePrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CityEventData() *config.CityEventDataConfig {
	fi := getMockFunc(s, s.CityEventData)
	if fi != nil {
		f, ok := fi.(func() *config.CityEventDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CityEventData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CityEventLevelData() *config.CityEventLevelDataConfig {
	fi := getMockFunc(s, s.CityEventLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.CityEventLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CityEventLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CityEventMiscData() *domestic_data.CityEventMiscData {
	fi := getMockFunc(s, s.CityEventMiscData)
	if fi != nil {
		f, ok := fi.(func() *domestic_data.CityEventMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CityEventMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CollectionExchangeData() *config.CollectionExchangeDataConfig {
	fi := getMockFunc(s, s.CollectionExchangeData)
	if fi != nil {
		f, ok := fi.(func() *config.CollectionExchangeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CollectionExchangeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ColorData() *config.ColorDataConfig {
	fi := getMockFunc(s, s.ColorData)
	if fi != nil {
		f, ok := fi.(func() *config.ColorDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ColorData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CombatConfig() *combatdata.CombatConfig {
	fi := getMockFunc(s, s.CombatConfig)
	if fi != nil {
		f, ok := fi.(func() *combatdata.CombatConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CombatConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CombatMiscConfig() *combatdata.CombatMiscConfig {
	fi := getMockFunc(s, s.CombatMiscConfig)
	if fi != nil {
		f, ok := fi.(func() *combatdata.CombatMiscConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CombatMiscConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CombatScene() *config.CombatSceneConfig {
	fi := getMockFunc(s, s.CombatScene)
	if fi != nil {
		f, ok := fi.(func() *config.CombatSceneConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CombatScene()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CombineCost() *config.CombineCostConfig {
	fi := getMockFunc(s, s.CombineCost)
	if fi != nil {
		f, ok := fi.(func() *config.CombineCostConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CombineCost()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ConditionPlunder() *config.ConditionPlunderConfig {
	fi := getMockFunc(s, s.ConditionPlunder)
	if fi != nil {
		f, ok := fi.(func() *config.ConditionPlunderConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ConditionPlunder()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ConditionPlunderItem() *config.ConditionPlunderItemConfig {
	fi := getMockFunc(s, s.ConditionPlunderItem)
	if fi != nil {
		f, ok := fi.(func() *config.ConditionPlunderItemConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ConditionPlunderItem()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) Cost() *config.CostConfig {
	fi := getMockFunc(s, s.Cost)
	if fi != nil {
		f, ok := fi.(func() *config.CostConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.Cost()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CountdownPrizeData() *config.CountdownPrizeDataConfig {
	fi := getMockFunc(s, s.CountdownPrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.CountdownPrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CountdownPrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CountdownPrizeDescData() *config.CountdownPrizeDescDataConfig {
	fi := getMockFunc(s, s.CountdownPrizeDescData)
	if fi != nil {
		f, ok := fi.(func() *config.CountdownPrizeDescDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CountdownPrizeDescData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CountryData() *config.CountryDataConfig {
	fi := getMockFunc(s, s.CountryData)
	if fi != nil {
		f, ok := fi.(func() *config.CountryDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CountryData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CountryMiscData() *country.CountryMiscData {
	fi := getMockFunc(s, s.CountryMiscData)
	if fi != nil {
		f, ok := fi.(func() *country.CountryMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CountryMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CountryOfficialData() *config.CountryOfficialDataConfig {
	fi := getMockFunc(s, s.CountryOfficialData)
	if fi != nil {
		f, ok := fi.(func() *config.CountryOfficialDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CountryOfficialData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) CountryOfficialNpcData() *config.CountryOfficialNpcDataConfig {
	fi := getMockFunc(s, s.CountryOfficialNpcData)
	if fi != nil {
		f, ok := fi.(func() *config.CountryOfficialNpcDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.CountryOfficialNpcData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) DailyBargainData() *config.DailyBargainDataConfig {
	fi := getMockFunc(s, s.DailyBargainData)
	if fi != nil {
		f, ok := fi.(func() *config.DailyBargainDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.DailyBargainData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) DiscountColorData() *config.DiscountColorDataConfig {
	fi := getMockFunc(s, s.DiscountColorData)
	if fi != nil {
		f, ok := fi.(func() *config.DiscountColorDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.DiscountColorData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) DungeonChapterData() *config.DungeonChapterDataConfig {
	fi := getMockFunc(s, s.DungeonChapterData)
	if fi != nil {
		f, ok := fi.(func() *config.DungeonChapterDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.DungeonChapterData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) DungeonData() *config.DungeonDataConfig {
	fi := getMockFunc(s, s.DungeonData)
	if fi != nil {
		f, ok := fi.(func() *config.DungeonDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.DungeonData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) DungeonGuideTroopData() *config.DungeonGuideTroopDataConfig {
	fi := getMockFunc(s, s.DungeonGuideTroopData)
	if fi != nil {
		f, ok := fi.(func() *config.DungeonGuideTroopDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.DungeonGuideTroopData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) DungeonMiscData() *dungeon.DungeonMiscData {
	fi := getMockFunc(s, s.DungeonMiscData)
	if fi != nil {
		f, ok := fi.(func() *dungeon.DungeonMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.DungeonMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) DurationCardData() *config.DurationCardDataConfig {
	fi := getMockFunc(s, s.DurationCardData)
	if fi != nil {
		f, ok := fi.(func() *config.DurationCardDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.DurationCardData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EncodeClient() *shared_proto.Config {
	fi := getMockFunc(s, s.EncodeClient)
	if fi != nil {
		f, ok := fi.(func() *shared_proto.Config)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EncodeClient()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EquipCombineData() *config.EquipCombineDataConfig {
	fi := getMockFunc(s, s.EquipCombineData)
	if fi != nil {
		f, ok := fi.(func() *config.EquipCombineDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EquipCombineData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EquipCombineDatas() *combine.EquipCombineDatas {
	fi := getMockFunc(s, s.EquipCombineDatas)
	if fi != nil {
		f, ok := fi.(func() *combine.EquipCombineDatas)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EquipCombineDatas()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EquipmentData() *config.EquipmentDataConfig {
	fi := getMockFunc(s, s.EquipmentData)
	if fi != nil {
		f, ok := fi.(func() *config.EquipmentDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EquipmentData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EquipmentLevelData() *config.EquipmentLevelDataConfig {
	fi := getMockFunc(s, s.EquipmentLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.EquipmentLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EquipmentLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EquipmentQualityData() *config.EquipmentQualityDataConfig {
	fi := getMockFunc(s, s.EquipmentQualityData)
	if fi != nil {
		f, ok := fi.(func() *config.EquipmentQualityDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EquipmentQualityData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EquipmentRefinedData() *config.EquipmentRefinedDataConfig {
	fi := getMockFunc(s, s.EquipmentRefinedData)
	if fi != nil {
		f, ok := fi.(func() *config.EquipmentRefinedDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EquipmentRefinedData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EquipmentTaozConfig() *goods.EquipmentTaozConfig {
	fi := getMockFunc(s, s.EquipmentTaozConfig)
	if fi != nil {
		f, ok := fi.(func() *goods.EquipmentTaozConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EquipmentTaozConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EquipmentTaozData() *config.EquipmentTaozDataConfig {
	fi := getMockFunc(s, s.EquipmentTaozData)
	if fi != nil {
		f, ok := fi.(func() *config.EquipmentTaozDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EquipmentTaozData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EventLimitGiftConfig() *promdata.EventLimitGiftConfig {
	fi := getMockFunc(s, s.EventLimitGiftConfig)
	if fi != nil {
		f, ok := fi.(func() *promdata.EventLimitGiftConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EventLimitGiftConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EventLimitGiftData() *config.EventLimitGiftDataConfig {
	fi := getMockFunc(s, s.EventLimitGiftData)
	if fi != nil {
		f, ok := fi.(func() *config.EventLimitGiftDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EventLimitGiftData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EventOptionData() *config.EventOptionDataConfig {
	fi := getMockFunc(s, s.EventOptionData)
	if fi != nil {
		f, ok := fi.(func() *config.EventOptionDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EventOptionData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) EventPosition() *config.EventPositionConfig {
	fi := getMockFunc(s, s.EventPosition)
	if fi != nil {
		f, ok := fi.(func() *config.EventPositionConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.EventPosition()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ExchangeMiscData() *dianquan.ExchangeMiscData {
	fi := getMockFunc(s, s.ExchangeMiscData)
	if fi != nil {
		f, ok := fi.(func() *dianquan.ExchangeMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ExchangeMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FamilyName() *config.FamilyNameConfig {
	fi := getMockFunc(s, s.FamilyName)
	if fi != nil {
		f, ok := fi.(func() *config.FamilyNameConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FamilyName()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FamilyNameData() *config.FamilyNameDataConfig {
	fi := getMockFunc(s, s.FamilyNameData)
	if fi != nil {
		f, ok := fi.(func() *config.FamilyNameDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FamilyNameData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FarmMaxStealConfig() *config.FarmMaxStealConfigConfig {
	fi := getMockFunc(s, s.FarmMaxStealConfig)
	if fi != nil {
		f, ok := fi.(func() *config.FarmMaxStealConfigConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FarmMaxStealConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FarmMiscConfig() *farm.FarmMiscConfig {
	fi := getMockFunc(s, s.FarmMiscConfig)
	if fi != nil {
		f, ok := fi.(func() *farm.FarmMiscConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FarmMiscConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FarmOneKeyConfig() *config.FarmOneKeyConfigConfig {
	fi := getMockFunc(s, s.FarmOneKeyConfig)
	if fi != nil {
		f, ok := fi.(func() *config.FarmOneKeyConfigConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FarmOneKeyConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FarmResConfig() *config.FarmResConfigConfig {
	fi := getMockFunc(s, s.FarmResConfig)
	if fi != nil {
		f, ok := fi.(func() *config.FarmResConfigConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FarmResConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FemaleGivenName() *config.FemaleGivenNameConfig {
	fi := getMockFunc(s, s.FemaleGivenName)
	if fi != nil {
		f, ok := fi.(func() *config.FemaleGivenNameConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FemaleGivenName()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FishData() *config.FishDataConfig {
	fi := getMockFunc(s, s.FishData)
	if fi != nil {
		f, ok := fi.(func() *config.FishDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FishData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FishRandomer() *fishing_data.FishRandomer {
	fi := getMockFunc(s, s.FishRandomer)
	if fi != nil {
		f, ok := fi.(func() *fishing_data.FishRandomer)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FishRandomer()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FishingCaptainProbabilityData() *config.FishingCaptainProbabilityDataConfig {
	fi := getMockFunc(s, s.FishingCaptainProbabilityData)
	if fi != nil {
		f, ok := fi.(func() *config.FishingCaptainProbabilityDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FishingCaptainProbabilityData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FishingCostData() *config.FishingCostDataConfig {
	fi := getMockFunc(s, s.FishingCostData)
	if fi != nil {
		f, ok := fi.(func() *config.FishingCostDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FishingCostData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FishingShowData() *config.FishingShowDataConfig {
	fi := getMockFunc(s, s.FishingShowData)
	if fi != nil {
		f, ok := fi.(func() *config.FishingShowDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FishingShowData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FreeGiftData() *config.FreeGiftDataConfig {
	fi := getMockFunc(s, s.FreeGiftData)
	if fi != nil {
		f, ok := fi.(func() *config.FreeGiftDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FreeGiftData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) FunctionOpenData() *config.FunctionOpenDataConfig {
	fi := getMockFunc(s, s.FunctionOpenData)
	if fi != nil {
		f, ok := fi.(func() *config.FunctionOpenDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.FunctionOpenData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GardenConfig() *gardendata.GardenConfig {
	fi := getMockFunc(s, s.GardenConfig)
	if fi != nil {
		f, ok := fi.(func() *gardendata.GardenConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GardenConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GemData() *config.GemDataConfig {
	fi := getMockFunc(s, s.GemData)
	if fi != nil {
		f, ok := fi.(func() *config.GemDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GemData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GemDatas() *goods.GemDatas {
	fi := getMockFunc(s, s.GemDatas)
	if fi != nil {
		f, ok := fi.(func() *goods.GemDatas)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GemDatas()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetAchieveTaskData(a0 uint64) *taskdata.AchieveTaskData {
	fi := getMockFunc(s, s.GetAchieveTaskData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.AchieveTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetAchieveTaskData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetAchieveTaskDataArray() []*taskdata.AchieveTaskData {
	fi := getMockFunc(s, s.GetAchieveTaskDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.AchieveTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetAchieveTaskDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetAchieveTaskStarPrizeData(a0 uint64) *taskdata.AchieveTaskStarPrizeData {
	fi := getMockFunc(s, s.GetAchieveTaskStarPrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.AchieveTaskStarPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetAchieveTaskStarPrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetAchieveTaskStarPrizeDataArray() []*taskdata.AchieveTaskStarPrizeData {
	fi := getMockFunc(s, s.GetAchieveTaskStarPrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.AchieveTaskStarPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetAchieveTaskStarPrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetActiveDegreePrizeData(a0 uint64) *taskdata.ActiveDegreePrizeData {
	fi := getMockFunc(s, s.GetActiveDegreePrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.ActiveDegreePrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActiveDegreePrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetActiveDegreePrizeDataArray() []*taskdata.ActiveDegreePrizeData {
	fi := getMockFunc(s, s.GetActiveDegreePrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.ActiveDegreePrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActiveDegreePrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetActiveDegreeTaskData(a0 uint64) *taskdata.ActiveDegreeTaskData {
	fi := getMockFunc(s, s.GetActiveDegreeTaskData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.ActiveDegreeTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActiveDegreeTaskData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetActiveDegreeTaskDataArray() []*taskdata.ActiveDegreeTaskData {
	fi := getMockFunc(s, s.GetActiveDegreeTaskDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.ActiveDegreeTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActiveDegreeTaskDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetActivityCollectionData(a0 uint64) *activitydata.ActivityCollectionData {
	fi := getMockFunc(s, s.GetActivityCollectionData)
	if fi != nil {
		f, ok := fi.(func(uint64) *activitydata.ActivityCollectionData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActivityCollectionData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetActivityCollectionDataArray() []*activitydata.ActivityCollectionData {
	fi := getMockFunc(s, s.GetActivityCollectionDataArray)
	if fi != nil {
		f, ok := fi.(func() []*activitydata.ActivityCollectionData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActivityCollectionDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetActivityShowData(a0 uint64) *activitydata.ActivityShowData {
	fi := getMockFunc(s, s.GetActivityShowData)
	if fi != nil {
		f, ok := fi.(func(uint64) *activitydata.ActivityShowData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActivityShowData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetActivityShowDataArray() []*activitydata.ActivityShowData {
	fi := getMockFunc(s, s.GetActivityShowDataArray)
	if fi != nil {
		f, ok := fi.(func() []*activitydata.ActivityShowData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActivityShowDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetActivityTaskData(a0 uint64) *taskdata.ActivityTaskData {
	fi := getMockFunc(s, s.GetActivityTaskData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.ActivityTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActivityTaskData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetActivityTaskDataArray() []*taskdata.ActivityTaskData {
	fi := getMockFunc(s, s.GetActivityTaskDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.ActivityTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActivityTaskDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetActivityTaskListModeData(a0 uint64) *activitydata.ActivityTaskListModeData {
	fi := getMockFunc(s, s.GetActivityTaskListModeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *activitydata.ActivityTaskListModeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActivityTaskListModeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetActivityTaskListModeDataArray() []*activitydata.ActivityTaskListModeData {
	fi := getMockFunc(s, s.GetActivityTaskListModeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*activitydata.ActivityTaskListModeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetActivityTaskListModeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetAmountShowSortData(a0 uint64) *resdata.AmountShowSortData {
	fi := getMockFunc(s, s.GetAmountShowSortData)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.AmountShowSortData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetAmountShowSortData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetAmountShowSortDataArray() []*resdata.AmountShowSortData {
	fi := getMockFunc(s, s.GetAmountShowSortDataArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.AmountShowSortData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetAmountShowSortDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetAreaData(a0 uint64) *regdata.AreaData {
	fi := getMockFunc(s, s.GetAreaData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.AreaData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetAreaData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetAreaDataArray() []*regdata.AreaData {
	fi := getMockFunc(s, s.GetAreaDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.AreaData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetAreaDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetAssemblyData(a0 uint64) *regdata.AssemblyData {
	fi := getMockFunc(s, s.GetAssemblyData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.AssemblyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetAssemblyData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetAssemblyDataArray() []*regdata.AssemblyData {
	fi := getMockFunc(s, s.GetAssemblyDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.AssemblyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetAssemblyDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBaYeStageData(a0 uint64) *taskdata.BaYeStageData {
	fi := getMockFunc(s, s.GetBaYeStageData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.BaYeStageData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBaYeStageData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBaYeStageDataArray() []*taskdata.BaYeStageData {
	fi := getMockFunc(s, s.GetBaYeStageDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.BaYeStageData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBaYeStageDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBaYeTaskData(a0 uint64) *taskdata.BaYeTaskData {
	fi := getMockFunc(s, s.GetBaYeTaskData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.BaYeTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBaYeTaskData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBaYeTaskDataArray() []*taskdata.BaYeTaskData {
	fi := getMockFunc(s, s.GetBaYeTaskDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.BaYeTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBaYeTaskDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBaowuData(a0 uint64) *resdata.BaowuData {
	fi := getMockFunc(s, s.GetBaowuData)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.BaowuData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBaowuData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBaowuDataArray() []*resdata.BaowuData {
	fi := getMockFunc(s, s.GetBaowuDataArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.BaowuData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBaowuDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBaozNpcData(a0 uint64) *regdata.BaozNpcData {
	fi := getMockFunc(s, s.GetBaozNpcData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.BaozNpcData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBaozNpcData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBaozNpcDataArray() []*regdata.BaozNpcData {
	fi := getMockFunc(s, s.GetBaozNpcDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.BaozNpcData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBaozNpcDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBaseLevelData(a0 uint64) *domestic_data.BaseLevelData {
	fi := getMockFunc(s, s.GetBaseLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.BaseLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBaseLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBaseLevelDataArray() []*domestic_data.BaseLevelData {
	fi := getMockFunc(s, s.GetBaseLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.BaseLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBaseLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBlackMarketData(a0 uint64) *shop.BlackMarketData {
	fi := getMockFunc(s, s.GetBlackMarketData)
	if fi != nil {
		f, ok := fi.(func(uint64) *shop.BlackMarketData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBlackMarketData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBlackMarketDataArray() []*shop.BlackMarketData {
	fi := getMockFunc(s, s.GetBlackMarketDataArray)
	if fi != nil {
		f, ok := fi.(func() []*shop.BlackMarketData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBlackMarketDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBlackMarketGoodsData(a0 uint64) *shop.BlackMarketGoodsData {
	fi := getMockFunc(s, s.GetBlackMarketGoodsData)
	if fi != nil {
		f, ok := fi.(func(uint64) *shop.BlackMarketGoodsData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBlackMarketGoodsData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBlackMarketGoodsDataArray() []*shop.BlackMarketGoodsData {
	fi := getMockFunc(s, s.GetBlackMarketGoodsDataArray)
	if fi != nil {
		f, ok := fi.(func() []*shop.BlackMarketGoodsData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBlackMarketGoodsDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBlackMarketGoodsGroupData(a0 uint64) *shop.BlackMarketGoodsGroupData {
	fi := getMockFunc(s, s.GetBlackMarketGoodsGroupData)
	if fi != nil {
		f, ok := fi.(func(uint64) *shop.BlackMarketGoodsGroupData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBlackMarketGoodsGroupData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBlackMarketGoodsGroupDataArray() []*shop.BlackMarketGoodsGroupData {
	fi := getMockFunc(s, s.GetBlackMarketGoodsGroupDataArray)
	if fi != nil {
		f, ok := fi.(func() []*shop.BlackMarketGoodsGroupData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBlackMarketGoodsGroupDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBlockData(a0 uint64) *blockdata.BlockData {
	fi := getMockFunc(s, s.GetBlockData)
	if fi != nil {
		f, ok := fi.(func(uint64) *blockdata.BlockData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBlockData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBlockDataArray() []*blockdata.BlockData {
	fi := getMockFunc(s, s.GetBlockDataArray)
	if fi != nil {
		f, ok := fi.(func() []*blockdata.BlockData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBlockDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBodyData(a0 uint64) *body.BodyData {
	fi := getMockFunc(s, s.GetBodyData)
	if fi != nil {
		f, ok := fi.(func(uint64) *body.BodyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBodyData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBodyDataArray() []*body.BodyData {
	fi := getMockFunc(s, s.GetBodyDataArray)
	if fi != nil {
		f, ok := fi.(func() []*body.BodyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBodyDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBranchTaskData(a0 uint64) *taskdata.BranchTaskData {
	fi := getMockFunc(s, s.GetBranchTaskData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.BranchTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBranchTaskData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBranchTaskDataArray() []*taskdata.BranchTaskData {
	fi := getMockFunc(s, s.GetBranchTaskDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.BranchTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBranchTaskDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBroadcastData(a0 string) *data.BroadcastData {
	fi := getMockFunc(s, s.GetBroadcastData)
	if fi != nil {
		f, ok := fi.(func(string) *data.BroadcastData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBroadcastData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBroadcastDataArray() []*data.BroadcastData {
	fi := getMockFunc(s, s.GetBroadcastDataArray)
	if fi != nil {
		f, ok := fi.(func() []*data.BroadcastData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBroadcastDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBuffEffectData(a0 uint64) *data.BuffEffectData {
	fi := getMockFunc(s, s.GetBuffEffectData)
	if fi != nil {
		f, ok := fi.(func(uint64) *data.BuffEffectData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBuffEffectData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBuffEffectDataArray() []*data.BuffEffectData {
	fi := getMockFunc(s, s.GetBuffEffectDataArray)
	if fi != nil {
		f, ok := fi.(func() []*data.BuffEffectData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBuffEffectDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBufferData(a0 uint64) *buffer.BufferData {
	fi := getMockFunc(s, s.GetBufferData)
	if fi != nil {
		f, ok := fi.(func(uint64) *buffer.BufferData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBufferData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBufferDataArray() []*buffer.BufferData {
	fi := getMockFunc(s, s.GetBufferDataArray)
	if fi != nil {
		f, ok := fi.(func() []*buffer.BufferData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBufferDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBufferTypeData(a0 uint64) *buffer.BufferTypeData {
	fi := getMockFunc(s, s.GetBufferTypeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *buffer.BufferTypeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBufferTypeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBufferTypeDataArray() []*buffer.BufferTypeData {
	fi := getMockFunc(s, s.GetBufferTypeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*buffer.BufferTypeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBufferTypeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBuildingData(a0 uint64) *domestic_data.BuildingData {
	fi := getMockFunc(s, s.GetBuildingData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.BuildingData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBuildingData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBuildingDataArray() []*domestic_data.BuildingData {
	fi := getMockFunc(s, s.GetBuildingDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.BuildingData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBuildingDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBuildingEffectData(a0 int) *sub.BuildingEffectData {
	fi := getMockFunc(s, s.GetBuildingEffectData)
	if fi != nil {
		f, ok := fi.(func(int) *sub.BuildingEffectData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBuildingEffectData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBuildingEffectDataArray() []*sub.BuildingEffectData {
	fi := getMockFunc(s, s.GetBuildingEffectDataArray)
	if fi != nil {
		f, ok := fi.(func() []*sub.BuildingEffectData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBuildingEffectDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBuildingLayoutData(a0 uint64) *domestic_data.BuildingLayoutData {
	fi := getMockFunc(s, s.GetBuildingLayoutData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.BuildingLayoutData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBuildingLayoutData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBuildingLayoutDataArray() []*domestic_data.BuildingLayoutData {
	fi := getMockFunc(s, s.GetBuildingLayoutDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.BuildingLayoutData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBuildingLayoutDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBuildingUnlockData(a0 uint64) *domestic_data.BuildingUnlockData {
	fi := getMockFunc(s, s.GetBuildingUnlockData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.BuildingUnlockData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBuildingUnlockData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBuildingUnlockDataArray() []*domestic_data.BuildingUnlockData {
	fi := getMockFunc(s, s.GetBuildingUnlockDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.BuildingUnlockData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBuildingUnlockDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBwzlPrizeData(a0 uint64) *taskdata.BwzlPrizeData {
	fi := getMockFunc(s, s.GetBwzlPrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.BwzlPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBwzlPrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBwzlPrizeDataArray() []*taskdata.BwzlPrizeData {
	fi := getMockFunc(s, s.GetBwzlPrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.BwzlPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBwzlPrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetBwzlTaskData(a0 uint64) *taskdata.BwzlTaskData {
	fi := getMockFunc(s, s.GetBwzlTaskData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.BwzlTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBwzlTaskData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetBwzlTaskDataArray() []*taskdata.BwzlTaskData {
	fi := getMockFunc(s, s.GetBwzlTaskDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.BwzlTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetBwzlTaskDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainAbilityData(a0 uint64) *captain.CaptainAbilityData {
	fi := getMockFunc(s, s.GetCaptainAbilityData)
	if fi != nil {
		f, ok := fi.(func(uint64) *captain.CaptainAbilityData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainAbilityData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainAbilityDataArray() []*captain.CaptainAbilityData {
	fi := getMockFunc(s, s.GetCaptainAbilityDataArray)
	if fi != nil {
		f, ok := fi.(func() []*captain.CaptainAbilityData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainAbilityDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainData(a0 uint64) *captain.CaptainData {
	fi := getMockFunc(s, s.GetCaptainData)
	if fi != nil {
		f, ok := fi.(func(uint64) *captain.CaptainData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainDataArray() []*captain.CaptainData {
	fi := getMockFunc(s, s.GetCaptainDataArray)
	if fi != nil {
		f, ok := fi.(func() []*captain.CaptainData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainFriendshipData(a0 uint64) *captain.CaptainFriendshipData {
	fi := getMockFunc(s, s.GetCaptainFriendshipData)
	if fi != nil {
		f, ok := fi.(func(uint64) *captain.CaptainFriendshipData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainFriendshipData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainFriendshipDataArray() []*captain.CaptainFriendshipData {
	fi := getMockFunc(s, s.GetCaptainFriendshipDataArray)
	if fi != nil {
		f, ok := fi.(func() []*captain.CaptainFriendshipData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainFriendshipDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainLevelData(a0 uint64) *captain.CaptainLevelData {
	fi := getMockFunc(s, s.GetCaptainLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *captain.CaptainLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainLevelDataArray() []*captain.CaptainLevelData {
	fi := getMockFunc(s, s.GetCaptainLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*captain.CaptainLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainOfficialCountData(a0 uint64) *captain.CaptainOfficialCountData {
	fi := getMockFunc(s, s.GetCaptainOfficialCountData)
	if fi != nil {
		f, ok := fi.(func(uint64) *captain.CaptainOfficialCountData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainOfficialCountData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainOfficialCountDataArray() []*captain.CaptainOfficialCountData {
	fi := getMockFunc(s, s.GetCaptainOfficialCountDataArray)
	if fi != nil {
		f, ok := fi.(func() []*captain.CaptainOfficialCountData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainOfficialCountDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainOfficialData(a0 uint64) *captain.CaptainOfficialData {
	fi := getMockFunc(s, s.GetCaptainOfficialData)
	if fi != nil {
		f, ok := fi.(func(uint64) *captain.CaptainOfficialData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainOfficialData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainOfficialDataArray() []*captain.CaptainOfficialData {
	fi := getMockFunc(s, s.GetCaptainOfficialDataArray)
	if fi != nil {
		f, ok := fi.(func() []*captain.CaptainOfficialData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainOfficialDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainRarityData(a0 uint64) *captain.CaptainRarityData {
	fi := getMockFunc(s, s.GetCaptainRarityData)
	if fi != nil {
		f, ok := fi.(func(uint64) *captain.CaptainRarityData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainRarityData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainRarityDataArray() []*captain.CaptainRarityData {
	fi := getMockFunc(s, s.GetCaptainRarityDataArray)
	if fi != nil {
		f, ok := fi.(func() []*captain.CaptainRarityData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainRarityDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainRebirthLevelData(a0 uint64) *captain.CaptainRebirthLevelData {
	fi := getMockFunc(s, s.GetCaptainRebirthLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *captain.CaptainRebirthLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainRebirthLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainRebirthLevelDataArray() []*captain.CaptainRebirthLevelData {
	fi := getMockFunc(s, s.GetCaptainRebirthLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*captain.CaptainRebirthLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainRebirthLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainStarData(a0 uint64) *captain.CaptainStarData {
	fi := getMockFunc(s, s.GetCaptainStarData)
	if fi != nil {
		f, ok := fi.(func(uint64) *captain.CaptainStarData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainStarData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCaptainStarDataArray() []*captain.CaptainStarData {
	fi := getMockFunc(s, s.GetCaptainStarDataArray)
	if fi != nil {
		f, ok := fi.(func() []*captain.CaptainStarData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCaptainStarDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetChargeObjData(a0 uint64) *charge.ChargeObjData {
	fi := getMockFunc(s, s.GetChargeObjData)
	if fi != nil {
		f, ok := fi.(func(uint64) *charge.ChargeObjData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetChargeObjData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetChargeObjDataArray() []*charge.ChargeObjData {
	fi := getMockFunc(s, s.GetChargeObjDataArray)
	if fi != nil {
		f, ok := fi.(func() []*charge.ChargeObjData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetChargeObjDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetChargePrizeData(a0 uint64) *charge.ChargePrizeData {
	fi := getMockFunc(s, s.GetChargePrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *charge.ChargePrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetChargePrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetChargePrizeDataArray() []*charge.ChargePrizeData {
	fi := getMockFunc(s, s.GetChargePrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*charge.ChargePrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetChargePrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCityEventData(a0 uint64) *domestic_data.CityEventData {
	fi := getMockFunc(s, s.GetCityEventData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.CityEventData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCityEventData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCityEventDataArray() []*domestic_data.CityEventData {
	fi := getMockFunc(s, s.GetCityEventDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.CityEventData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCityEventDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCityEventLevelData(a0 uint64) *domestic_data.CityEventLevelData {
	fi := getMockFunc(s, s.GetCityEventLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.CityEventLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCityEventLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCityEventLevelDataArray() []*domestic_data.CityEventLevelData {
	fi := getMockFunc(s, s.GetCityEventLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.CityEventLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCityEventLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCollectionExchangeData(a0 uint64) *activitydata.CollectionExchangeData {
	fi := getMockFunc(s, s.GetCollectionExchangeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *activitydata.CollectionExchangeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCollectionExchangeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCollectionExchangeDataArray() []*activitydata.CollectionExchangeData {
	fi := getMockFunc(s, s.GetCollectionExchangeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*activitydata.CollectionExchangeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCollectionExchangeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetColorData(a0 uint64) *data.ColorData {
	fi := getMockFunc(s, s.GetColorData)
	if fi != nil {
		f, ok := fi.(func(uint64) *data.ColorData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetColorData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetColorDataArray() []*data.ColorData {
	fi := getMockFunc(s, s.GetColorDataArray)
	if fi != nil {
		f, ok := fi.(func() []*data.ColorData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetColorDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCombatScene(a0 string) *scene.CombatScene {
	fi := getMockFunc(s, s.GetCombatScene)
	if fi != nil {
		f, ok := fi.(func(string) *scene.CombatScene)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCombatScene()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCombatSceneArray() []*scene.CombatScene {
	fi := getMockFunc(s, s.GetCombatSceneArray)
	if fi != nil {
		f, ok := fi.(func() []*scene.CombatScene)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCombatSceneArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCombineCost(a0 int) *domestic_data.CombineCost {
	fi := getMockFunc(s, s.GetCombineCost)
	if fi != nil {
		f, ok := fi.(func(int) *domestic_data.CombineCost)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCombineCost()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCombineCostArray() []*domestic_data.CombineCost {
	fi := getMockFunc(s, s.GetCombineCostArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.CombineCost)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCombineCostArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetConditionPlunder(a0 uint64) *resdata.ConditionPlunder {
	fi := getMockFunc(s, s.GetConditionPlunder)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.ConditionPlunder)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetConditionPlunder()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetConditionPlunderArray() []*resdata.ConditionPlunder {
	fi := getMockFunc(s, s.GetConditionPlunderArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.ConditionPlunder)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetConditionPlunderArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetConditionPlunderItem(a0 uint64) *resdata.ConditionPlunderItem {
	fi := getMockFunc(s, s.GetConditionPlunderItem)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.ConditionPlunderItem)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetConditionPlunderItem()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetConditionPlunderItemArray() []*resdata.ConditionPlunderItem {
	fi := getMockFunc(s, s.GetConditionPlunderItemArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.ConditionPlunderItem)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetConditionPlunderItemArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCost(a0 int) *resdata.Cost {
	fi := getMockFunc(s, s.GetCost)
	if fi != nil {
		f, ok := fi.(func(int) *resdata.Cost)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCost()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCostArray() []*resdata.Cost {
	fi := getMockFunc(s, s.GetCostArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.Cost)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCostArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCountdownPrizeData(a0 uint64) *domestic_data.CountdownPrizeData {
	fi := getMockFunc(s, s.GetCountdownPrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.CountdownPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCountdownPrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCountdownPrizeDataArray() []*domestic_data.CountdownPrizeData {
	fi := getMockFunc(s, s.GetCountdownPrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.CountdownPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCountdownPrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCountdownPrizeDescData(a0 uint64) *domestic_data.CountdownPrizeDescData {
	fi := getMockFunc(s, s.GetCountdownPrizeDescData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.CountdownPrizeDescData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCountdownPrizeDescData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCountdownPrizeDescDataArray() []*domestic_data.CountdownPrizeDescData {
	fi := getMockFunc(s, s.GetCountdownPrizeDescDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.CountdownPrizeDescData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCountdownPrizeDescDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCountryData(a0 uint64) *country.CountryData {
	fi := getMockFunc(s, s.GetCountryData)
	if fi != nil {
		f, ok := fi.(func(uint64) *country.CountryData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCountryData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCountryDataArray() []*country.CountryData {
	fi := getMockFunc(s, s.GetCountryDataArray)
	if fi != nil {
		f, ok := fi.(func() []*country.CountryData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCountryDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCountryOfficialData(a0 int) *country.CountryOfficialData {
	fi := getMockFunc(s, s.GetCountryOfficialData)
	if fi != nil {
		f, ok := fi.(func(int) *country.CountryOfficialData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCountryOfficialData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCountryOfficialDataArray() []*country.CountryOfficialData {
	fi := getMockFunc(s, s.GetCountryOfficialDataArray)
	if fi != nil {
		f, ok := fi.(func() []*country.CountryOfficialData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCountryOfficialDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetCountryOfficialNpcData(a0 uint64) *country.CountryOfficialNpcData {
	fi := getMockFunc(s, s.GetCountryOfficialNpcData)
	if fi != nil {
		f, ok := fi.(func(uint64) *country.CountryOfficialNpcData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCountryOfficialNpcData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetCountryOfficialNpcDataArray() []*country.CountryOfficialNpcData {
	fi := getMockFunc(s, s.GetCountryOfficialNpcDataArray)
	if fi != nil {
		f, ok := fi.(func() []*country.CountryOfficialNpcData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetCountryOfficialNpcDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetDailyBargainData(a0 uint64) *promdata.DailyBargainData {
	fi := getMockFunc(s, s.GetDailyBargainData)
	if fi != nil {
		f, ok := fi.(func(uint64) *promdata.DailyBargainData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDailyBargainData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetDailyBargainDataArray() []*promdata.DailyBargainData {
	fi := getMockFunc(s, s.GetDailyBargainDataArray)
	if fi != nil {
		f, ok := fi.(func() []*promdata.DailyBargainData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDailyBargainDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetDiscountColorData(a0 uint64) *shop.DiscountColorData {
	fi := getMockFunc(s, s.GetDiscountColorData)
	if fi != nil {
		f, ok := fi.(func(uint64) *shop.DiscountColorData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDiscountColorData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetDiscountColorDataArray() []*shop.DiscountColorData {
	fi := getMockFunc(s, s.GetDiscountColorDataArray)
	if fi != nil {
		f, ok := fi.(func() []*shop.DiscountColorData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDiscountColorDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetDungeonChapterData(a0 uint64) *dungeon.DungeonChapterData {
	fi := getMockFunc(s, s.GetDungeonChapterData)
	if fi != nil {
		f, ok := fi.(func(uint64) *dungeon.DungeonChapterData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDungeonChapterData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetDungeonChapterDataArray() []*dungeon.DungeonChapterData {
	fi := getMockFunc(s, s.GetDungeonChapterDataArray)
	if fi != nil {
		f, ok := fi.(func() []*dungeon.DungeonChapterData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDungeonChapterDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetDungeonData(a0 uint64) *dungeon.DungeonData {
	fi := getMockFunc(s, s.GetDungeonData)
	if fi != nil {
		f, ok := fi.(func(uint64) *dungeon.DungeonData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDungeonData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetDungeonDataArray() []*dungeon.DungeonData {
	fi := getMockFunc(s, s.GetDungeonDataArray)
	if fi != nil {
		f, ok := fi.(func() []*dungeon.DungeonData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDungeonDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetDungeonGuideTroopData(a0 uint64) *dungeon.DungeonGuideTroopData {
	fi := getMockFunc(s, s.GetDungeonGuideTroopData)
	if fi != nil {
		f, ok := fi.(func(uint64) *dungeon.DungeonGuideTroopData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDungeonGuideTroopData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetDungeonGuideTroopDataArray() []*dungeon.DungeonGuideTroopData {
	fi := getMockFunc(s, s.GetDungeonGuideTroopDataArray)
	if fi != nil {
		f, ok := fi.(func() []*dungeon.DungeonGuideTroopData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDungeonGuideTroopDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetDurationCardData(a0 uint64) *promdata.DurationCardData {
	fi := getMockFunc(s, s.GetDurationCardData)
	if fi != nil {
		f, ok := fi.(func(uint64) *promdata.DurationCardData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDurationCardData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetDurationCardDataArray() []*promdata.DurationCardData {
	fi := getMockFunc(s, s.GetDurationCardDataArray)
	if fi != nil {
		f, ok := fi.(func() []*promdata.DurationCardData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetDurationCardDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetEquipCombineData(a0 uint64) *combine.EquipCombineData {
	fi := getMockFunc(s, s.GetEquipCombineData)
	if fi != nil {
		f, ok := fi.(func(uint64) *combine.EquipCombineData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipCombineData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetEquipCombineDataArray() []*combine.EquipCombineData {
	fi := getMockFunc(s, s.GetEquipCombineDataArray)
	if fi != nil {
		f, ok := fi.(func() []*combine.EquipCombineData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipCombineDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetEquipmentData(a0 uint64) *goods.EquipmentData {
	fi := getMockFunc(s, s.GetEquipmentData)
	if fi != nil {
		f, ok := fi.(func(uint64) *goods.EquipmentData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipmentData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetEquipmentDataArray() []*goods.EquipmentData {
	fi := getMockFunc(s, s.GetEquipmentDataArray)
	if fi != nil {
		f, ok := fi.(func() []*goods.EquipmentData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipmentDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetEquipmentLevelData(a0 uint64) *goods.EquipmentLevelData {
	fi := getMockFunc(s, s.GetEquipmentLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *goods.EquipmentLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipmentLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetEquipmentLevelDataArray() []*goods.EquipmentLevelData {
	fi := getMockFunc(s, s.GetEquipmentLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*goods.EquipmentLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipmentLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetEquipmentQualityData(a0 uint64) *goods.EquipmentQualityData {
	fi := getMockFunc(s, s.GetEquipmentQualityData)
	if fi != nil {
		f, ok := fi.(func(uint64) *goods.EquipmentQualityData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipmentQualityData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetEquipmentQualityDataArray() []*goods.EquipmentQualityData {
	fi := getMockFunc(s, s.GetEquipmentQualityDataArray)
	if fi != nil {
		f, ok := fi.(func() []*goods.EquipmentQualityData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipmentQualityDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetEquipmentRefinedData(a0 uint64) *goods.EquipmentRefinedData {
	fi := getMockFunc(s, s.GetEquipmentRefinedData)
	if fi != nil {
		f, ok := fi.(func(uint64) *goods.EquipmentRefinedData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipmentRefinedData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetEquipmentRefinedDataArray() []*goods.EquipmentRefinedData {
	fi := getMockFunc(s, s.GetEquipmentRefinedDataArray)
	if fi != nil {
		f, ok := fi.(func() []*goods.EquipmentRefinedData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipmentRefinedDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetEquipmentTaozData(a0 uint64) *goods.EquipmentTaozData {
	fi := getMockFunc(s, s.GetEquipmentTaozData)
	if fi != nil {
		f, ok := fi.(func(uint64) *goods.EquipmentTaozData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipmentTaozData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetEquipmentTaozDataArray() []*goods.EquipmentTaozData {
	fi := getMockFunc(s, s.GetEquipmentTaozDataArray)
	if fi != nil {
		f, ok := fi.(func() []*goods.EquipmentTaozData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEquipmentTaozDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetEventLimitGiftData(a0 uint64) *promdata.EventLimitGiftData {
	fi := getMockFunc(s, s.GetEventLimitGiftData)
	if fi != nil {
		f, ok := fi.(func(uint64) *promdata.EventLimitGiftData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEventLimitGiftData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetEventLimitGiftDataArray() []*promdata.EventLimitGiftData {
	fi := getMockFunc(s, s.GetEventLimitGiftDataArray)
	if fi != nil {
		f, ok := fi.(func() []*promdata.EventLimitGiftData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEventLimitGiftDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetEventOptionData(a0 uint64) *random_event.EventOptionData {
	fi := getMockFunc(s, s.GetEventOptionData)
	if fi != nil {
		f, ok := fi.(func(uint64) *random_event.EventOptionData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEventOptionData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetEventOptionDataArray() []*random_event.EventOptionData {
	fi := getMockFunc(s, s.GetEventOptionDataArray)
	if fi != nil {
		f, ok := fi.(func() []*random_event.EventOptionData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEventOptionDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetEventPosition(a0 uint64) *random_event.EventPosition {
	fi := getMockFunc(s, s.GetEventPosition)
	if fi != nil {
		f, ok := fi.(func(uint64) *random_event.EventPosition)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEventPosition()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetEventPositionArray() []*random_event.EventPosition {
	fi := getMockFunc(s, s.GetEventPositionArray)
	if fi != nil {
		f, ok := fi.(func() []*random_event.EventPosition)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetEventPositionArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFamilyName(a0 string) *data.FamilyName {
	fi := getMockFunc(s, s.GetFamilyName)
	if fi != nil {
		f, ok := fi.(func(string) *data.FamilyName)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFamilyName()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFamilyNameArray() []*data.FamilyName {
	fi := getMockFunc(s, s.GetFamilyNameArray)
	if fi != nil {
		f, ok := fi.(func() []*data.FamilyName)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFamilyNameArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFamilyNameData(a0 uint64) *country.FamilyNameData {
	fi := getMockFunc(s, s.GetFamilyNameData)
	if fi != nil {
		f, ok := fi.(func(uint64) *country.FamilyNameData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFamilyNameData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFamilyNameDataArray() []*country.FamilyNameData {
	fi := getMockFunc(s, s.GetFamilyNameDataArray)
	if fi != nil {
		f, ok := fi.(func() []*country.FamilyNameData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFamilyNameDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFarmMaxStealConfig(a0 uint64) *farm.FarmMaxStealConfig {
	fi := getMockFunc(s, s.GetFarmMaxStealConfig)
	if fi != nil {
		f, ok := fi.(func(uint64) *farm.FarmMaxStealConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFarmMaxStealConfig()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFarmMaxStealConfigArray() []*farm.FarmMaxStealConfig {
	fi := getMockFunc(s, s.GetFarmMaxStealConfigArray)
	if fi != nil {
		f, ok := fi.(func() []*farm.FarmMaxStealConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFarmMaxStealConfigArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFarmOneKeyConfig(a0 uint64) *farm.FarmOneKeyConfig {
	fi := getMockFunc(s, s.GetFarmOneKeyConfig)
	if fi != nil {
		f, ok := fi.(func(uint64) *farm.FarmOneKeyConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFarmOneKeyConfig()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFarmOneKeyConfigArray() []*farm.FarmOneKeyConfig {
	fi := getMockFunc(s, s.GetFarmOneKeyConfigArray)
	if fi != nil {
		f, ok := fi.(func() []*farm.FarmOneKeyConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFarmOneKeyConfigArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFarmResConfig(a0 uint64) *farm.FarmResConfig {
	fi := getMockFunc(s, s.GetFarmResConfig)
	if fi != nil {
		f, ok := fi.(func(uint64) *farm.FarmResConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFarmResConfig()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFarmResConfigArray() []*farm.FarmResConfig {
	fi := getMockFunc(s, s.GetFarmResConfigArray)
	if fi != nil {
		f, ok := fi.(func() []*farm.FarmResConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFarmResConfigArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFemaleGivenName(a0 string) *data.FemaleGivenName {
	fi := getMockFunc(s, s.GetFemaleGivenName)
	if fi != nil {
		f, ok := fi.(func(string) *data.FemaleGivenName)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFemaleGivenName()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFemaleGivenNameArray() []*data.FemaleGivenName {
	fi := getMockFunc(s, s.GetFemaleGivenNameArray)
	if fi != nil {
		f, ok := fi.(func() []*data.FemaleGivenName)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFemaleGivenNameArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFishData(a0 uint64) *fishing_data.FishData {
	fi := getMockFunc(s, s.GetFishData)
	if fi != nil {
		f, ok := fi.(func(uint64) *fishing_data.FishData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFishData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFishDataArray() []*fishing_data.FishData {
	fi := getMockFunc(s, s.GetFishDataArray)
	if fi != nil {
		f, ok := fi.(func() []*fishing_data.FishData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFishDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFishingCaptainProbabilityData(a0 uint64) *fishing_data.FishingCaptainProbabilityData {
	fi := getMockFunc(s, s.GetFishingCaptainProbabilityData)
	if fi != nil {
		f, ok := fi.(func(uint64) *fishing_data.FishingCaptainProbabilityData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFishingCaptainProbabilityData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFishingCaptainProbabilityDataArray() []*fishing_data.FishingCaptainProbabilityData {
	fi := getMockFunc(s, s.GetFishingCaptainProbabilityDataArray)
	if fi != nil {
		f, ok := fi.(func() []*fishing_data.FishingCaptainProbabilityData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFishingCaptainProbabilityDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFishingCostData(a0 uint64) *fishing_data.FishingCostData {
	fi := getMockFunc(s, s.GetFishingCostData)
	if fi != nil {
		f, ok := fi.(func(uint64) *fishing_data.FishingCostData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFishingCostData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFishingCostDataArray() []*fishing_data.FishingCostData {
	fi := getMockFunc(s, s.GetFishingCostDataArray)
	if fi != nil {
		f, ok := fi.(func() []*fishing_data.FishingCostData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFishingCostDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFishingShowData(a0 uint64) *fishing_data.FishingShowData {
	fi := getMockFunc(s, s.GetFishingShowData)
	if fi != nil {
		f, ok := fi.(func(uint64) *fishing_data.FishingShowData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFishingShowData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFishingShowDataArray() []*fishing_data.FishingShowData {
	fi := getMockFunc(s, s.GetFishingShowDataArray)
	if fi != nil {
		f, ok := fi.(func() []*fishing_data.FishingShowData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFishingShowDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFreeGiftData(a0 uint64) *promdata.FreeGiftData {
	fi := getMockFunc(s, s.GetFreeGiftData)
	if fi != nil {
		f, ok := fi.(func(uint64) *promdata.FreeGiftData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFreeGiftData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFreeGiftDataArray() []*promdata.FreeGiftData {
	fi := getMockFunc(s, s.GetFreeGiftDataArray)
	if fi != nil {
		f, ok := fi.(func() []*promdata.FreeGiftData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFreeGiftDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetFunctionOpenData(a0 uint64) *function.FunctionOpenData {
	fi := getMockFunc(s, s.GetFunctionOpenData)
	if fi != nil {
		f, ok := fi.(func(uint64) *function.FunctionOpenData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFunctionOpenData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetFunctionOpenDataArray() []*function.FunctionOpenData {
	fi := getMockFunc(s, s.GetFunctionOpenDataArray)
	if fi != nil {
		f, ok := fi.(func() []*function.FunctionOpenData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetFunctionOpenDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGemData(a0 uint64) *goods.GemData {
	fi := getMockFunc(s, s.GetGemData)
	if fi != nil {
		f, ok := fi.(func(uint64) *goods.GemData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGemData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGemDataArray() []*goods.GemData {
	fi := getMockFunc(s, s.GetGemDataArray)
	if fi != nil {
		f, ok := fi.(func() []*goods.GemData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGemDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGoodsCombineData(a0 uint64) *combine.GoodsCombineData {
	fi := getMockFunc(s, s.GetGoodsCombineData)
	if fi != nil {
		f, ok := fi.(func(uint64) *combine.GoodsCombineData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGoodsCombineData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGoodsCombineDataArray() []*combine.GoodsCombineData {
	fi := getMockFunc(s, s.GetGoodsCombineDataArray)
	if fi != nil {
		f, ok := fi.(func() []*combine.GoodsCombineData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGoodsCombineDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGoodsData(a0 uint64) *goods.GoodsData {
	fi := getMockFunc(s, s.GetGoodsData)
	if fi != nil {
		f, ok := fi.(func(uint64) *goods.GoodsData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGoodsData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGoodsDataArray() []*goods.GoodsData {
	fi := getMockFunc(s, s.GetGoodsDataArray)
	if fi != nil {
		f, ok := fi.(func() []*goods.GoodsData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGoodsDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGoodsQuality(a0 uint64) *goods.GoodsQuality {
	fi := getMockFunc(s, s.GetGoodsQuality)
	if fi != nil {
		f, ok := fi.(func(uint64) *goods.GoodsQuality)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGoodsQuality()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGoodsQualityArray() []*goods.GoodsQuality {
	fi := getMockFunc(s, s.GetGoodsQualityArray)
	if fi != nil {
		f, ok := fi.(func() []*goods.GoodsQuality)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGoodsQualityArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuanFuLevelData(a0 uint64) *domestic_data.GuanFuLevelData {
	fi := getMockFunc(s, s.GetGuanFuLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.GuanFuLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuanFuLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuanFuLevelDataArray() []*domestic_data.GuanFuLevelData {
	fi := getMockFunc(s, s.GetGuanFuLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.GuanFuLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuanFuLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildBigBoxData(a0 uint64) *guild_data.GuildBigBoxData {
	fi := getMockFunc(s, s.GetGuildBigBoxData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildBigBoxData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildBigBoxData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildBigBoxDataArray() []*guild_data.GuildBigBoxData {
	fi := getMockFunc(s, s.GetGuildBigBoxDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildBigBoxData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildBigBoxDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildClassLevelData(a0 uint64) *guild_data.GuildClassLevelData {
	fi := getMockFunc(s, s.GetGuildClassLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildClassLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildClassLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildClassLevelDataArray() []*guild_data.GuildClassLevelData {
	fi := getMockFunc(s, s.GetGuildClassLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildClassLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildClassLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildClassTitleData(a0 uint64) *guild_data.GuildClassTitleData {
	fi := getMockFunc(s, s.GetGuildClassTitleData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildClassTitleData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildClassTitleData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildClassTitleDataArray() []*guild_data.GuildClassTitleData {
	fi := getMockFunc(s, s.GetGuildClassTitleDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildClassTitleData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildClassTitleDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildDonateData(a0 uint64) *guild_data.GuildDonateData {
	fi := getMockFunc(s, s.GetGuildDonateData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildDonateData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildDonateData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildDonateDataArray() []*guild_data.GuildDonateData {
	fi := getMockFunc(s, s.GetGuildDonateDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildDonateData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildDonateDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildEventPrizeData(a0 uint64) *guild_data.GuildEventPrizeData {
	fi := getMockFunc(s, s.GetGuildEventPrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildEventPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildEventPrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildEventPrizeDataArray() []*guild_data.GuildEventPrizeData {
	fi := getMockFunc(s, s.GetGuildEventPrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildEventPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildEventPrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildLevelCdrData(a0 uint64) *guild_data.GuildLevelCdrData {
	fi := getMockFunc(s, s.GetGuildLevelCdrData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildLevelCdrData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildLevelCdrData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildLevelCdrDataArray() []*guild_data.GuildLevelCdrData {
	fi := getMockFunc(s, s.GetGuildLevelCdrDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildLevelCdrData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildLevelCdrDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildLevelData(a0 uint64) *guild_data.GuildLevelData {
	fi := getMockFunc(s, s.GetGuildLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildLevelDataArray() []*guild_data.GuildLevelData {
	fi := getMockFunc(s, s.GetGuildLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildLevelPrize(a0 uint64) *resdata.GuildLevelPrize {
	fi := getMockFunc(s, s.GetGuildLevelPrize)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.GuildLevelPrize)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildLevelPrize()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildLevelPrizeArray() []*resdata.GuildLevelPrize {
	fi := getMockFunc(s, s.GetGuildLevelPrizeArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.GuildLevelPrize)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildLevelPrizeArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildLogData(a0 string) *guild_data.GuildLogData {
	fi := getMockFunc(s, s.GetGuildLogData)
	if fi != nil {
		f, ok := fi.(func(string) *guild_data.GuildLogData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildLogData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildLogDataArray() []*guild_data.GuildLogData {
	fi := getMockFunc(s, s.GetGuildLogDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildLogData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildLogDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildPermissionShowData(a0 uint64) *guild_data.GuildPermissionShowData {
	fi := getMockFunc(s, s.GetGuildPermissionShowData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildPermissionShowData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildPermissionShowData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildPermissionShowDataArray() []*guild_data.GuildPermissionShowData {
	fi := getMockFunc(s, s.GetGuildPermissionShowDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildPermissionShowData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildPermissionShowDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildPrestigeEventData(a0 uint64) *guild_data.GuildPrestigeEventData {
	fi := getMockFunc(s, s.GetGuildPrestigeEventData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildPrestigeEventData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildPrestigeEventData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildPrestigeEventDataArray() []*guild_data.GuildPrestigeEventData {
	fi := getMockFunc(s, s.GetGuildPrestigeEventDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildPrestigeEventData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildPrestigeEventDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildPrestigePrizeData(a0 uint64) *guild_data.GuildPrestigePrizeData {
	fi := getMockFunc(s, s.GetGuildPrestigePrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildPrestigePrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildPrestigePrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildPrestigePrizeDataArray() []*guild_data.GuildPrestigePrizeData {
	fi := getMockFunc(s, s.GetGuildPrestigePrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildPrestigePrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildPrestigePrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildRankPrizeData(a0 uint64) *guild_data.GuildRankPrizeData {
	fi := getMockFunc(s, s.GetGuildRankPrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildRankPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildRankPrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildRankPrizeDataArray() []*guild_data.GuildRankPrizeData {
	fi := getMockFunc(s, s.GetGuildRankPrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildRankPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildRankPrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildTarget(a0 uint64) *guild_data.GuildTarget {
	fi := getMockFunc(s, s.GetGuildTarget)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildTarget)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildTarget()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildTargetArray() []*guild_data.GuildTarget {
	fi := getMockFunc(s, s.GetGuildTargetArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildTarget)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildTargetArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildTaskData(a0 uint64) *guild_data.GuildTaskData {
	fi := getMockFunc(s, s.GetGuildTaskData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildTaskData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildTaskDataArray() []*guild_data.GuildTaskData {
	fi := getMockFunc(s, s.GetGuildTaskDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildTaskDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildTaskEvaluateData(a0 uint64) *guild_data.GuildTaskEvaluateData {
	fi := getMockFunc(s, s.GetGuildTaskEvaluateData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildTaskEvaluateData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildTaskEvaluateData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildTaskEvaluateDataArray() []*guild_data.GuildTaskEvaluateData {
	fi := getMockFunc(s, s.GetGuildTaskEvaluateDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildTaskEvaluateData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildTaskEvaluateDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetGuildTechnologyData(a0 uint64) *guild_data.GuildTechnologyData {
	fi := getMockFunc(s, s.GetGuildTechnologyData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.GuildTechnologyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildTechnologyData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetGuildTechnologyDataArray() []*guild_data.GuildTechnologyData {
	fi := getMockFunc(s, s.GetGuildTechnologyDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.GuildTechnologyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetGuildTechnologyDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetHeadData(a0 string) *head.HeadData {
	fi := getMockFunc(s, s.GetHeadData)
	if fi != nil {
		f, ok := fi.(func(string) *head.HeadData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHeadData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetHeadDataArray() []*head.HeadData {
	fi := getMockFunc(s, s.GetHeadDataArray)
	if fi != nil {
		f, ok := fi.(func() []*head.HeadData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHeadDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetHebiPrizeData(a0 uint64) *hebi.HebiPrizeData {
	fi := getMockFunc(s, s.GetHebiPrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *hebi.HebiPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHebiPrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetHebiPrizeDataArray() []*hebi.HebiPrizeData {
	fi := getMockFunc(s, s.GetHebiPrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*hebi.HebiPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHebiPrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetHeroLevelData(a0 uint64) *herodata.HeroLevelData {
	fi := getMockFunc(s, s.GetHeroLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *herodata.HeroLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHeroLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetHeroLevelDataArray() []*herodata.HeroLevelData {
	fi := getMockFunc(s, s.GetHeroLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*herodata.HeroLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHeroLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetHeroLevelFundData(a0 uint64) *promdata.HeroLevelFundData {
	fi := getMockFunc(s, s.GetHeroLevelFundData)
	if fi != nil {
		f, ok := fi.(func(uint64) *promdata.HeroLevelFundData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHeroLevelFundData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetHeroLevelFundDataArray() []*promdata.HeroLevelFundData {
	fi := getMockFunc(s, s.GetHeroLevelFundDataArray)
	if fi != nil {
		f, ok := fi.(func() []*promdata.HeroLevelFundData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHeroLevelFundDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetHeroLevelSubData(a0 uint64) *data.HeroLevelSubData {
	fi := getMockFunc(s, s.GetHeroLevelSubData)
	if fi != nil {
		f, ok := fi.(func(uint64) *data.HeroLevelSubData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHeroLevelSubData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetHeroLevelSubDataArray() []*data.HeroLevelSubData {
	fi := getMockFunc(s, s.GetHeroLevelSubDataArray)
	if fi != nil {
		f, ok := fi.(func() []*data.HeroLevelSubData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHeroLevelSubDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetHomeNpcBaseData(a0 uint64) *basedata.HomeNpcBaseData {
	fi := getMockFunc(s, s.GetHomeNpcBaseData)
	if fi != nil {
		f, ok := fi.(func(uint64) *basedata.HomeNpcBaseData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHomeNpcBaseData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetHomeNpcBaseDataArray() []*basedata.HomeNpcBaseData {
	fi := getMockFunc(s, s.GetHomeNpcBaseDataArray)
	if fi != nil {
		f, ok := fi.(func() []*basedata.HomeNpcBaseData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetHomeNpcBaseDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetI18nData(a0 string) *i18n.I18nData {
	fi := getMockFunc(s, s.GetI18nData)
	if fi != nil {
		f, ok := fi.(func(string) *i18n.I18nData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetI18nData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetI18nDataArray() []*i18n.I18nData {
	fi := getMockFunc(s, s.GetI18nDataArray)
	if fi != nil {
		f, ok := fi.(func() []*i18n.I18nData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetI18nDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetIcon(a0 string) *icon.Icon {
	fi := getMockFunc(s, s.GetIcon)
	if fi != nil {
		f, ok := fi.(func(string) *icon.Icon)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetIcon()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetIconArray() []*icon.Icon {
	fi := getMockFunc(s, s.GetIconArray)
	if fi != nil {
		f, ok := fi.(func() []*icon.Icon)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetIconArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetJiuGuanData(a0 uint64) *military_data.JiuGuanData {
	fi := getMockFunc(s, s.GetJiuGuanData)
	if fi != nil {
		f, ok := fi.(func(uint64) *military_data.JiuGuanData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJiuGuanData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetJiuGuanDataArray() []*military_data.JiuGuanData {
	fi := getMockFunc(s, s.GetJiuGuanDataArray)
	if fi != nil {
		f, ok := fi.(func() []*military_data.JiuGuanData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJiuGuanDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetJunTuanNpcData(a0 uint64) *regdata.JunTuanNpcData {
	fi := getMockFunc(s, s.GetJunTuanNpcData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.JunTuanNpcData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJunTuanNpcData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetJunTuanNpcDataArray() []*regdata.JunTuanNpcData {
	fi := getMockFunc(s, s.GetJunTuanNpcDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.JunTuanNpcData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJunTuanNpcDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetJunTuanNpcPlaceData(a0 uint64) *regdata.JunTuanNpcPlaceData {
	fi := getMockFunc(s, s.GetJunTuanNpcPlaceData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.JunTuanNpcPlaceData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJunTuanNpcPlaceData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetJunTuanNpcPlaceDataArray() []*regdata.JunTuanNpcPlaceData {
	fi := getMockFunc(s, s.GetJunTuanNpcPlaceDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.JunTuanNpcPlaceData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJunTuanNpcPlaceDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetJunXianLevelData(a0 uint64) *bai_zhan_data.JunXianLevelData {
	fi := getMockFunc(s, s.GetJunXianLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *bai_zhan_data.JunXianLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJunXianLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetJunXianLevelDataArray() []*bai_zhan_data.JunXianLevelData {
	fi := getMockFunc(s, s.GetJunXianLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*bai_zhan_data.JunXianLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJunXianLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetJunXianPrizeData(a0 uint64) *bai_zhan_data.JunXianPrizeData {
	fi := getMockFunc(s, s.GetJunXianPrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *bai_zhan_data.JunXianPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJunXianPrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetJunXianPrizeDataArray() []*bai_zhan_data.JunXianPrizeData {
	fi := getMockFunc(s, s.GetJunXianPrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*bai_zhan_data.JunXianPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJunXianPrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetJunYingLevelData(a0 uint64) *military_data.JunYingLevelData {
	fi := getMockFunc(s, s.GetJunYingLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *military_data.JunYingLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJunYingLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetJunYingLevelDataArray() []*military_data.JunYingLevelData {
	fi := getMockFunc(s, s.GetJunYingLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*military_data.JunYingLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetJunYingLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetLocationData(a0 uint64) *location.LocationData {
	fi := getMockFunc(s, s.GetLocationData)
	if fi != nil {
		f, ok := fi.(func(uint64) *location.LocationData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetLocationData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetLocationDataArray() []*location.LocationData {
	fi := getMockFunc(s, s.GetLocationDataArray)
	if fi != nil {
		f, ok := fi.(func() []*location.LocationData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetLocationDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetLoginDayData(a0 uint64) *promdata.LoginDayData {
	fi := getMockFunc(s, s.GetLoginDayData)
	if fi != nil {
		f, ok := fi.(func(uint64) *promdata.LoginDayData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetLoginDayData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetLoginDayDataArray() []*promdata.LoginDayData {
	fi := getMockFunc(s, s.GetLoginDayDataArray)
	if fi != nil {
		f, ok := fi.(func() []*promdata.LoginDayData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetLoginDayDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMailData(a0 string) *maildata.MailData {
	fi := getMockFunc(s, s.GetMailData)
	if fi != nil {
		f, ok := fi.(func(string) *maildata.MailData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMailData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMailDataArray() []*maildata.MailData {
	fi := getMockFunc(s, s.GetMailDataArray)
	if fi != nil {
		f, ok := fi.(func() []*maildata.MailData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMailDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMainTaskData(a0 uint64) *taskdata.MainTaskData {
	fi := getMockFunc(s, s.GetMainTaskData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.MainTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMainTaskData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMainTaskDataArray() []*taskdata.MainTaskData {
	fi := getMockFunc(s, s.GetMainTaskDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.MainTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMainTaskDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMaleGivenName(a0 string) *data.MaleGivenName {
	fi := getMockFunc(s, s.GetMaleGivenName)
	if fi != nil {
		f, ok := fi.(func(string) *data.MaleGivenName)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMaleGivenName()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMaleGivenNameArray() []*data.MaleGivenName {
	fi := getMockFunc(s, s.GetMaleGivenNameArray)
	if fi != nil {
		f, ok := fi.(func() []*data.MaleGivenName)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMaleGivenNameArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMcBuildAddSupportData(a0 uint64) *mingcdata.McBuildAddSupportData {
	fi := getMockFunc(s, s.GetMcBuildAddSupportData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.McBuildAddSupportData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMcBuildAddSupportData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMcBuildAddSupportDataArray() []*mingcdata.McBuildAddSupportData {
	fi := getMockFunc(s, s.GetMcBuildAddSupportDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.McBuildAddSupportData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMcBuildAddSupportDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMcBuildGuildMemberPrizeData(a0 uint64) *mingcdata.McBuildGuildMemberPrizeData {
	fi := getMockFunc(s, s.GetMcBuildGuildMemberPrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.McBuildGuildMemberPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMcBuildGuildMemberPrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMcBuildGuildMemberPrizeDataArray() []*mingcdata.McBuildGuildMemberPrizeData {
	fi := getMockFunc(s, s.GetMcBuildGuildMemberPrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.McBuildGuildMemberPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMcBuildGuildMemberPrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMcBuildMcSupportData(a0 uint64) *mingcdata.McBuildMcSupportData {
	fi := getMockFunc(s, s.GetMcBuildMcSupportData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.McBuildMcSupportData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMcBuildMcSupportData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMcBuildMcSupportDataArray() []*mingcdata.McBuildMcSupportData {
	fi := getMockFunc(s, s.GetMcBuildMcSupportDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.McBuildMcSupportData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMcBuildMcSupportDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcBaseData(a0 uint64) *mingcdata.MingcBaseData {
	fi := getMockFunc(s, s.GetMingcBaseData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcBaseData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcBaseData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcBaseDataArray() []*mingcdata.MingcBaseData {
	fi := getMockFunc(s, s.GetMingcBaseDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcBaseData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcBaseDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcTimeData(a0 uint64) *mingcdata.MingcTimeData {
	fi := getMockFunc(s, s.GetMingcTimeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcTimeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcTimeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcTimeDataArray() []*mingcdata.MingcTimeData {
	fi := getMockFunc(s, s.GetMingcTimeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcTimeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcTimeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarBuildingData(a0 uint64) *mingcdata.MingcWarBuildingData {
	fi := getMockFunc(s, s.GetMingcWarBuildingData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcWarBuildingData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarBuildingData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarBuildingDataArray() []*mingcdata.MingcWarBuildingData {
	fi := getMockFunc(s, s.GetMingcWarBuildingDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcWarBuildingData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarBuildingDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarDrumStatData(a0 uint64) *mingcdata.MingcWarDrumStatData {
	fi := getMockFunc(s, s.GetMingcWarDrumStatData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcWarDrumStatData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarDrumStatData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarDrumStatDataArray() []*mingcdata.MingcWarDrumStatData {
	fi := getMockFunc(s, s.GetMingcWarDrumStatDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcWarDrumStatData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarDrumStatDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarMapData(a0 uint64) *mingcdata.MingcWarMapData {
	fi := getMockFunc(s, s.GetMingcWarMapData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcWarMapData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarMapData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarMapDataArray() []*mingcdata.MingcWarMapData {
	fi := getMockFunc(s, s.GetMingcWarMapDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcWarMapData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarMapDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarMultiKillData(a0 uint64) *mingcdata.MingcWarMultiKillData {
	fi := getMockFunc(s, s.GetMingcWarMultiKillData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcWarMultiKillData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarMultiKillData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarMultiKillDataArray() []*mingcdata.MingcWarMultiKillData {
	fi := getMockFunc(s, s.GetMingcWarMultiKillDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcWarMultiKillData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarMultiKillDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarNpcData(a0 uint64) *mingcdata.MingcWarNpcData {
	fi := getMockFunc(s, s.GetMingcWarNpcData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcWarNpcData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarNpcData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarNpcDataArray() []*mingcdata.MingcWarNpcData {
	fi := getMockFunc(s, s.GetMingcWarNpcDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcWarNpcData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarNpcDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarNpcGuildData(a0 uint64) *mingcdata.MingcWarNpcGuildData {
	fi := getMockFunc(s, s.GetMingcWarNpcGuildData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcWarNpcGuildData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarNpcGuildData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarNpcGuildDataArray() []*mingcdata.MingcWarNpcGuildData {
	fi := getMockFunc(s, s.GetMingcWarNpcGuildDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcWarNpcGuildData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarNpcGuildDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarSceneData(a0 uint64) *mingcdata.MingcWarSceneData {
	fi := getMockFunc(s, s.GetMingcWarSceneData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcWarSceneData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarSceneData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarSceneDataArray() []*mingcdata.MingcWarSceneData {
	fi := getMockFunc(s, s.GetMingcWarSceneDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcWarSceneData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarSceneDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarTouShiBuildingTargetData(a0 uint64) *mingcdata.MingcWarTouShiBuildingTargetData {
	fi := getMockFunc(s, s.GetMingcWarTouShiBuildingTargetData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcWarTouShiBuildingTargetData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarTouShiBuildingTargetData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarTouShiBuildingTargetDataArray() []*mingcdata.MingcWarTouShiBuildingTargetData {
	fi := getMockFunc(s, s.GetMingcWarTouShiBuildingTargetDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcWarTouShiBuildingTargetData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarTouShiBuildingTargetDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarTroopLastBeatWhenFailData(a0 uint64) *mingcdata.MingcWarTroopLastBeatWhenFailData {
	fi := getMockFunc(s, s.GetMingcWarTroopLastBeatWhenFailData)
	if fi != nil {
		f, ok := fi.(func(uint64) *mingcdata.MingcWarTroopLastBeatWhenFailData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarTroopLastBeatWhenFailData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMingcWarTroopLastBeatWhenFailDataArray() []*mingcdata.MingcWarTroopLastBeatWhenFailData {
	fi := getMockFunc(s, s.GetMingcWarTroopLastBeatWhenFailDataArray)
	if fi != nil {
		f, ok := fi.(func() []*mingcdata.MingcWarTroopLastBeatWhenFailData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMingcWarTroopLastBeatWhenFailDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMonsterCaptainData(a0 uint64) *monsterdata.MonsterCaptainData {
	fi := getMockFunc(s, s.GetMonsterCaptainData)
	if fi != nil {
		f, ok := fi.(func(uint64) *monsterdata.MonsterCaptainData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMonsterCaptainData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMonsterCaptainDataArray() []*monsterdata.MonsterCaptainData {
	fi := getMockFunc(s, s.GetMonsterCaptainDataArray)
	if fi != nil {
		f, ok := fi.(func() []*monsterdata.MonsterCaptainData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMonsterCaptainDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetMonsterMasterData(a0 uint64) *monsterdata.MonsterMasterData {
	fi := getMockFunc(s, s.GetMonsterMasterData)
	if fi != nil {
		f, ok := fi.(func(uint64) *monsterdata.MonsterMasterData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMonsterMasterData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetMonsterMasterDataArray() []*monsterdata.MonsterMasterData {
	fi := getMockFunc(s, s.GetMonsterMasterDataArray)
	if fi != nil {
		f, ok := fi.(func() []*monsterdata.MonsterMasterData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetMonsterMasterDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetNamelessCaptainData(a0 uint64) *captain.NamelessCaptainData {
	fi := getMockFunc(s, s.GetNamelessCaptainData)
	if fi != nil {
		f, ok := fi.(func(uint64) *captain.NamelessCaptainData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetNamelessCaptainData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetNamelessCaptainDataArray() []*captain.NamelessCaptainData {
	fi := getMockFunc(s, s.GetNamelessCaptainDataArray)
	if fi != nil {
		f, ok := fi.(func() []*captain.NamelessCaptainData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetNamelessCaptainDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetNormalShopGoods(a0 uint64) *shop.NormalShopGoods {
	fi := getMockFunc(s, s.GetNormalShopGoods)
	if fi != nil {
		f, ok := fi.(func(uint64) *shop.NormalShopGoods)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetNormalShopGoods()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetNormalShopGoodsArray() []*shop.NormalShopGoods {
	fi := getMockFunc(s, s.GetNormalShopGoodsArray)
	if fi != nil {
		f, ok := fi.(func() []*shop.NormalShopGoods)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetNormalShopGoodsArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetNpcBaseData(a0 uint64) *basedata.NpcBaseData {
	fi := getMockFunc(s, s.GetNpcBaseData)
	if fi != nil {
		f, ok := fi.(func(uint64) *basedata.NpcBaseData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetNpcBaseData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetNpcBaseDataArray() []*basedata.NpcBaseData {
	fi := getMockFunc(s, s.GetNpcBaseDataArray)
	if fi != nil {
		f, ok := fi.(func() []*basedata.NpcBaseData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetNpcBaseDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetNpcGuildTemplate(a0 uint64) *guild_data.NpcGuildTemplate {
	fi := getMockFunc(s, s.GetNpcGuildTemplate)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.NpcGuildTemplate)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetNpcGuildTemplate()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetNpcGuildTemplateArray() []*guild_data.NpcGuildTemplate {
	fi := getMockFunc(s, s.GetNpcGuildTemplateArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.NpcGuildTemplate)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetNpcGuildTemplateArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetNpcMemberData(a0 uint64) *guild_data.NpcMemberData {
	fi := getMockFunc(s, s.GetNpcMemberData)
	if fi != nil {
		f, ok := fi.(func(uint64) *guild_data.NpcMemberData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetNpcMemberData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetNpcMemberDataArray() []*guild_data.NpcMemberData {
	fi := getMockFunc(s, s.GetNpcMemberDataArray)
	if fi != nil {
		f, ok := fi.(func() []*guild_data.NpcMemberData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetNpcMemberDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetOptionPrize(a0 uint64) *random_event.OptionPrize {
	fi := getMockFunc(s, s.GetOptionPrize)
	if fi != nil {
		f, ok := fi.(func(uint64) *random_event.OptionPrize)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetOptionPrize()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetOptionPrizeArray() []*random_event.OptionPrize {
	fi := getMockFunc(s, s.GetOptionPrizeArray)
	if fi != nil {
		f, ok := fi.(func() []*random_event.OptionPrize)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetOptionPrizeArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetOuterCityBuildingData(a0 uint64) *domestic_data.OuterCityBuildingData {
	fi := getMockFunc(s, s.GetOuterCityBuildingData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.OuterCityBuildingData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetOuterCityBuildingData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetOuterCityBuildingDataArray() []*domestic_data.OuterCityBuildingData {
	fi := getMockFunc(s, s.GetOuterCityBuildingDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.OuterCityBuildingData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetOuterCityBuildingDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetOuterCityBuildingDescData(a0 uint64) *domestic_data.OuterCityBuildingDescData {
	fi := getMockFunc(s, s.GetOuterCityBuildingDescData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.OuterCityBuildingDescData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetOuterCityBuildingDescData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetOuterCityBuildingDescDataArray() []*domestic_data.OuterCityBuildingDescData {
	fi := getMockFunc(s, s.GetOuterCityBuildingDescDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.OuterCityBuildingDescData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetOuterCityBuildingDescDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetOuterCityData(a0 uint64) *domestic_data.OuterCityData {
	fi := getMockFunc(s, s.GetOuterCityData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.OuterCityData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetOuterCityData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetOuterCityDataArray() []*domestic_data.OuterCityData {
	fi := getMockFunc(s, s.GetOuterCityDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.OuterCityData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetOuterCityDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetOuterCityLayoutData(a0 uint64) *domestic_data.OuterCityLayoutData {
	fi := getMockFunc(s, s.GetOuterCityLayoutData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.OuterCityLayoutData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetOuterCityLayoutData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetOuterCityLayoutDataArray() []*domestic_data.OuterCityLayoutData {
	fi := getMockFunc(s, s.GetOuterCityLayoutDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.OuterCityLayoutData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetOuterCityLayoutDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetPassiveSpellData(a0 uint64) *spell.PassiveSpellData {
	fi := getMockFunc(s, s.GetPassiveSpellData)
	if fi != nil {
		f, ok := fi.(func(uint64) *spell.PassiveSpellData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPassiveSpellData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetPassiveSpellDataArray() []*spell.PassiveSpellData {
	fi := getMockFunc(s, s.GetPassiveSpellDataArray)
	if fi != nil {
		f, ok := fi.(func() []*spell.PassiveSpellData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPassiveSpellDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetPlunder(a0 uint64) *resdata.Plunder {
	fi := getMockFunc(s, s.GetPlunder)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.Plunder)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPlunder()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetPlunderArray() []*resdata.Plunder {
	fi := getMockFunc(s, s.GetPlunderArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.Plunder)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPlunderArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetPlunderGroup(a0 uint64) *resdata.PlunderGroup {
	fi := getMockFunc(s, s.GetPlunderGroup)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.PlunderGroup)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPlunderGroup()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetPlunderGroupArray() []*resdata.PlunderGroup {
	fi := getMockFunc(s, s.GetPlunderGroupArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.PlunderGroup)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPlunderGroupArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetPlunderItem(a0 uint64) *resdata.PlunderItem {
	fi := getMockFunc(s, s.GetPlunderItem)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.PlunderItem)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPlunderItem()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetPlunderItemArray() []*resdata.PlunderItem {
	fi := getMockFunc(s, s.GetPlunderItemArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.PlunderItem)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPlunderItemArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetPlunderPrize(a0 uint64) *resdata.PlunderPrize {
	fi := getMockFunc(s, s.GetPlunderPrize)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.PlunderPrize)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPlunderPrize()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetPlunderPrizeArray() []*resdata.PlunderPrize {
	fi := getMockFunc(s, s.GetPlunderPrizeArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.PlunderPrize)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPlunderPrizeArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetPrivacySettingData(a0 uint64) *settings.PrivacySettingData {
	fi := getMockFunc(s, s.GetPrivacySettingData)
	if fi != nil {
		f, ok := fi.(func(uint64) *settings.PrivacySettingData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPrivacySettingData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetPrivacySettingDataArray() []*settings.PrivacySettingData {
	fi := getMockFunc(s, s.GetPrivacySettingDataArray)
	if fi != nil {
		f, ok := fi.(func() []*settings.PrivacySettingData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPrivacySettingDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetPrize(a0 int) *resdata.Prize {
	fi := getMockFunc(s, s.GetPrize)
	if fi != nil {
		f, ok := fi.(func(int) *resdata.Prize)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPrize()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetPrizeArray() []*resdata.Prize {
	fi := getMockFunc(s, s.GetPrizeArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.Prize)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPrizeArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetProductData(a0 uint64) *charge.ProductData {
	fi := getMockFunc(s, s.GetProductData)
	if fi != nil {
		f, ok := fi.(func(uint64) *charge.ProductData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetProductData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetProductDataArray() []*charge.ProductData {
	fi := getMockFunc(s, s.GetProductDataArray)
	if fi != nil {
		f, ok := fi.(func() []*charge.ProductData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetProductDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetProsperityDamageBuffData(a0 uint64) *domestic_data.ProsperityDamageBuffData {
	fi := getMockFunc(s, s.GetProsperityDamageBuffData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.ProsperityDamageBuffData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetProsperityDamageBuffData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetProsperityDamageBuffDataArray() []*domestic_data.ProsperityDamageBuffData {
	fi := getMockFunc(s, s.GetProsperityDamageBuffDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.ProsperityDamageBuffData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetProsperityDamageBuffDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetPushData(a0 uint64) *pushdata.PushData {
	fi := getMockFunc(s, s.GetPushData)
	if fi != nil {
		f, ok := fi.(func(uint64) *pushdata.PushData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPushData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetPushDataArray() []*pushdata.PushData {
	fi := getMockFunc(s, s.GetPushDataArray)
	if fi != nil {
		f, ok := fi.(func() []*pushdata.PushData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPushDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetPveTroopData(a0 uint64) *pvetroop.PveTroopData {
	fi := getMockFunc(s, s.GetPveTroopData)
	if fi != nil {
		f, ok := fi.(func(uint64) *pvetroop.PveTroopData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPveTroopData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetPveTroopDataArray() []*pvetroop.PveTroopData {
	fi := getMockFunc(s, s.GetPveTroopDataArray)
	if fi != nil {
		f, ok := fi.(func() []*pvetroop.PveTroopData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetPveTroopDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetQuestionData(a0 uint64) *question.QuestionData {
	fi := getMockFunc(s, s.GetQuestionData)
	if fi != nil {
		f, ok := fi.(func(uint64) *question.QuestionData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetQuestionData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetQuestionDataArray() []*question.QuestionData {
	fi := getMockFunc(s, s.GetQuestionDataArray)
	if fi != nil {
		f, ok := fi.(func() []*question.QuestionData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetQuestionDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetQuestionPrizeData(a0 uint64) *question.QuestionPrizeData {
	fi := getMockFunc(s, s.GetQuestionPrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *question.QuestionPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetQuestionPrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetQuestionPrizeDataArray() []*question.QuestionPrizeData {
	fi := getMockFunc(s, s.GetQuestionPrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*question.QuestionPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetQuestionPrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetQuestionSayingData(a0 uint64) *question.QuestionSayingData {
	fi := getMockFunc(s, s.GetQuestionSayingData)
	if fi != nil {
		f, ok := fi.(func(uint64) *question.QuestionSayingData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetQuestionSayingData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetQuestionSayingDataArray() []*question.QuestionSayingData {
	fi := getMockFunc(s, s.GetQuestionSayingDataArray)
	if fi != nil {
		f, ok := fi.(func() []*question.QuestionSayingData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetQuestionSayingDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetRaceData(a0 int) *race.RaceData {
	fi := getMockFunc(s, s.GetRaceData)
	if fi != nil {
		f, ok := fi.(func(int) *race.RaceData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRaceData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetRaceDataArray() []*race.RaceData {
	fi := getMockFunc(s, s.GetRaceDataArray)
	if fi != nil {
		f, ok := fi.(func() []*race.RaceData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRaceDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetRandomEventData(a0 uint64) *random_event.RandomEventData {
	fi := getMockFunc(s, s.GetRandomEventData)
	if fi != nil {
		f, ok := fi.(func(uint64) *random_event.RandomEventData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRandomEventData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetRandomEventDataArray() []*random_event.RandomEventData {
	fi := getMockFunc(s, s.GetRandomEventDataArray)
	if fi != nil {
		f, ok := fi.(func() []*random_event.RandomEventData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRandomEventDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetRedPacketData(a0 uint64) *red_packet.RedPacketData {
	fi := getMockFunc(s, s.GetRedPacketData)
	if fi != nil {
		f, ok := fi.(func(uint64) *red_packet.RedPacketData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRedPacketData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetRedPacketDataArray() []*red_packet.RedPacketData {
	fi := getMockFunc(s, s.GetRedPacketDataArray)
	if fi != nil {
		f, ok := fi.(func() []*red_packet.RedPacketData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRedPacketDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetRegionAreaData(a0 uint64) *regdata.RegionAreaData {
	fi := getMockFunc(s, s.GetRegionAreaData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.RegionAreaData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionAreaData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetRegionAreaDataArray() []*regdata.RegionAreaData {
	fi := getMockFunc(s, s.GetRegionAreaDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.RegionAreaData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionAreaDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetRegionData(a0 uint64) *regdata.RegionData {
	fi := getMockFunc(s, s.GetRegionData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.RegionData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetRegionDataArray() []*regdata.RegionData {
	fi := getMockFunc(s, s.GetRegionDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.RegionData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetRegionMonsterData(a0 uint64) *regdata.RegionMonsterData {
	fi := getMockFunc(s, s.GetRegionMonsterData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.RegionMonsterData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionMonsterData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetRegionMonsterDataArray() []*regdata.RegionMonsterData {
	fi := getMockFunc(s, s.GetRegionMonsterDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.RegionMonsterData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionMonsterDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetRegionMultiLevelNpcData(a0 uint64) *regdata.RegionMultiLevelNpcData {
	fi := getMockFunc(s, s.GetRegionMultiLevelNpcData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.RegionMultiLevelNpcData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionMultiLevelNpcData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetRegionMultiLevelNpcDataArray() []*regdata.RegionMultiLevelNpcData {
	fi := getMockFunc(s, s.GetRegionMultiLevelNpcDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.RegionMultiLevelNpcData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionMultiLevelNpcDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetRegionMultiLevelNpcLevelData(a0 uint64) *regdata.RegionMultiLevelNpcLevelData {
	fi := getMockFunc(s, s.GetRegionMultiLevelNpcLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.RegionMultiLevelNpcLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionMultiLevelNpcLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetRegionMultiLevelNpcLevelDataArray() []*regdata.RegionMultiLevelNpcLevelData {
	fi := getMockFunc(s, s.GetRegionMultiLevelNpcLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.RegionMultiLevelNpcLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionMultiLevelNpcLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetRegionMultiLevelNpcTypeData(a0 int) *regdata.RegionMultiLevelNpcTypeData {
	fi := getMockFunc(s, s.GetRegionMultiLevelNpcTypeData)
	if fi != nil {
		f, ok := fi.(func(int) *regdata.RegionMultiLevelNpcTypeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionMultiLevelNpcTypeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetRegionMultiLevelNpcTypeDataArray() []*regdata.RegionMultiLevelNpcTypeData {
	fi := getMockFunc(s, s.GetRegionMultiLevelNpcTypeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.RegionMultiLevelNpcTypeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetRegionMultiLevelNpcTypeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetResCaptainData(a0 uint64) *resdata.ResCaptainData {
	fi := getMockFunc(s, s.GetResCaptainData)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.ResCaptainData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetResCaptainData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetResCaptainDataArray() []*resdata.ResCaptainData {
	fi := getMockFunc(s, s.GetResCaptainDataArray)
	if fi != nil {
		f, ok := fi.(func() []*resdata.ResCaptainData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetResCaptainDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetResistXiongNuData(a0 uint64) *xiongnu.ResistXiongNuData {
	fi := getMockFunc(s, s.GetResistXiongNuData)
	if fi != nil {
		f, ok := fi.(func(uint64) *xiongnu.ResistXiongNuData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetResistXiongNuData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetResistXiongNuDataArray() []*xiongnu.ResistXiongNuData {
	fi := getMockFunc(s, s.GetResistXiongNuDataArray)
	if fi != nil {
		f, ok := fi.(func() []*xiongnu.ResistXiongNuData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetResistXiongNuDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetResistXiongNuScoreData(a0 uint64) *xiongnu.ResistXiongNuScoreData {
	fi := getMockFunc(s, s.GetResistXiongNuScoreData)
	if fi != nil {
		f, ok := fi.(func(uint64) *xiongnu.ResistXiongNuScoreData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetResistXiongNuScoreData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetResistXiongNuScoreDataArray() []*xiongnu.ResistXiongNuScoreData {
	fi := getMockFunc(s, s.GetResistXiongNuScoreDataArray)
	if fi != nil {
		f, ok := fi.(func() []*xiongnu.ResistXiongNuScoreData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetResistXiongNuScoreDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetResistXiongNuWaveData(a0 uint64) *xiongnu.ResistXiongNuWaveData {
	fi := getMockFunc(s, s.GetResistXiongNuWaveData)
	if fi != nil {
		f, ok := fi.(func(uint64) *xiongnu.ResistXiongNuWaveData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetResistXiongNuWaveData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetResistXiongNuWaveDataArray() []*xiongnu.ResistXiongNuWaveData {
	fi := getMockFunc(s, s.GetResistXiongNuWaveDataArray)
	if fi != nil {
		f, ok := fi.(func() []*xiongnu.ResistXiongNuWaveData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetResistXiongNuWaveDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetSeasonData(a0 uint64) *season.SeasonData {
	fi := getMockFunc(s, s.GetSeasonData)
	if fi != nil {
		f, ok := fi.(func(uint64) *season.SeasonData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSeasonData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetSeasonDataArray() []*season.SeasonData {
	fi := getMockFunc(s, s.GetSeasonDataArray)
	if fi != nil {
		f, ok := fi.(func() []*season.SeasonData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSeasonDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetSecretTowerData(a0 uint64) *towerdata.SecretTowerData {
	fi := getMockFunc(s, s.GetSecretTowerData)
	if fi != nil {
		f, ok := fi.(func(uint64) *towerdata.SecretTowerData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSecretTowerData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetSecretTowerDataArray() []*towerdata.SecretTowerData {
	fi := getMockFunc(s, s.GetSecretTowerDataArray)
	if fi != nil {
		f, ok := fi.(func() []*towerdata.SecretTowerData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSecretTowerDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetSecretTowerWordsData(a0 uint64) *towerdata.SecretTowerWordsData {
	fi := getMockFunc(s, s.GetSecretTowerWordsData)
	if fi != nil {
		f, ok := fi.(func(uint64) *towerdata.SecretTowerWordsData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSecretTowerWordsData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetSecretTowerWordsDataArray() []*towerdata.SecretTowerWordsData {
	fi := getMockFunc(s, s.GetSecretTowerWordsDataArray)
	if fi != nil {
		f, ok := fi.(func() []*towerdata.SecretTowerWordsData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSecretTowerWordsDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetShop(a0 uint64) *shop.Shop {
	fi := getMockFunc(s, s.GetShop)
	if fi != nil {
		f, ok := fi.(func(uint64) *shop.Shop)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetShop()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetShopArray() []*shop.Shop {
	fi := getMockFunc(s, s.GetShopArray)
	if fi != nil {
		f, ok := fi.(func() []*shop.Shop)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetShopArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetSoldierLevelData(a0 uint64) *domestic_data.SoldierLevelData {
	fi := getMockFunc(s, s.GetSoldierLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.SoldierLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSoldierLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetSoldierLevelDataArray() []*domestic_data.SoldierLevelData {
	fi := getMockFunc(s, s.GetSoldierLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.SoldierLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSoldierLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetSpCollectionData(a0 uint64) *promdata.SpCollectionData {
	fi := getMockFunc(s, s.GetSpCollectionData)
	if fi != nil {
		f, ok := fi.(func(uint64) *promdata.SpCollectionData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSpCollectionData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetSpCollectionDataArray() []*promdata.SpCollectionData {
	fi := getMockFunc(s, s.GetSpCollectionDataArray)
	if fi != nil {
		f, ok := fi.(func() []*promdata.SpCollectionData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSpCollectionDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetSpell(a0 uint64) *spell.Spell {
	fi := getMockFunc(s, s.GetSpell)
	if fi != nil {
		f, ok := fi.(func(uint64) *spell.Spell)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSpell()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetSpellArray() []*spell.Spell {
	fi := getMockFunc(s, s.GetSpellArray)
	if fi != nil {
		f, ok := fi.(func() []*spell.Spell)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSpellArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetSpellData(a0 uint64) *spell.SpellData {
	fi := getMockFunc(s, s.GetSpellData)
	if fi != nil {
		f, ok := fi.(func(uint64) *spell.SpellData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSpellData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetSpellDataArray() []*spell.SpellData {
	fi := getMockFunc(s, s.GetSpellDataArray)
	if fi != nil {
		f, ok := fi.(func() []*spell.SpellData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSpellDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetSpellFacadeData(a0 uint64) *spell.SpellFacadeData {
	fi := getMockFunc(s, s.GetSpellFacadeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *spell.SpellFacadeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSpellFacadeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetSpellFacadeDataArray() []*spell.SpellFacadeData {
	fi := getMockFunc(s, s.GetSpellFacadeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*spell.SpellFacadeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSpellFacadeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetSpriteStat(a0 uint64) *data.SpriteStat {
	fi := getMockFunc(s, s.GetSpriteStat)
	if fi != nil {
		f, ok := fi.(func(uint64) *data.SpriteStat)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSpriteStat()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetSpriteStatArray() []*data.SpriteStat {
	fi := getMockFunc(s, s.GetSpriteStatArray)
	if fi != nil {
		f, ok := fi.(func() []*data.SpriteStat)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSpriteStatArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetStateData(a0 uint64) *spell.StateData {
	fi := getMockFunc(s, s.GetStateData)
	if fi != nil {
		f, ok := fi.(func(uint64) *spell.StateData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetStateData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetStateDataArray() []*spell.StateData {
	fi := getMockFunc(s, s.GetStateDataArray)
	if fi != nil {
		f, ok := fi.(func() []*spell.StateData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetStateDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetStrategyData(a0 uint64) *strategydata.StrategyData {
	fi := getMockFunc(s, s.GetStrategyData)
	if fi != nil {
		f, ok := fi.(func(uint64) *strategydata.StrategyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetStrategyData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetStrategyDataArray() []*strategydata.StrategyData {
	fi := getMockFunc(s, s.GetStrategyDataArray)
	if fi != nil {
		f, ok := fi.(func() []*strategydata.StrategyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetStrategyDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetStrategyEffectData(a0 uint64) *strategydata.StrategyEffectData {
	fi := getMockFunc(s, s.GetStrategyEffectData)
	if fi != nil {
		f, ok := fi.(func(uint64) *strategydata.StrategyEffectData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetStrategyEffectData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetStrategyEffectDataArray() []*strategydata.StrategyEffectData {
	fi := getMockFunc(s, s.GetStrategyEffectDataArray)
	if fi != nil {
		f, ok := fi.(func() []*strategydata.StrategyEffectData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetStrategyEffectDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetStrongerData(a0 uint64) *strongerdata.StrongerData {
	fi := getMockFunc(s, s.GetStrongerData)
	if fi != nil {
		f, ok := fi.(func(uint64) *strongerdata.StrongerData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetStrongerData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetStrongerDataArray() []*strongerdata.StrongerData {
	fi := getMockFunc(s, s.GetStrongerDataArray)
	if fi != nil {
		f, ok := fi.(func() []*strongerdata.StrongerData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetStrongerDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetSurveyData(a0 string) *survey.SurveyData {
	fi := getMockFunc(s, s.GetSurveyData)
	if fi != nil {
		f, ok := fi.(func(string) *survey.SurveyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSurveyData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetSurveyDataArray() []*survey.SurveyData {
	fi := getMockFunc(s, s.GetSurveyDataArray)
	if fi != nil {
		f, ok := fi.(func() []*survey.SurveyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetSurveyDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTaskBoxData(a0 uint64) *taskdata.TaskBoxData {
	fi := getMockFunc(s, s.GetTaskBoxData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.TaskBoxData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTaskBoxData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTaskBoxDataArray() []*taskdata.TaskBoxData {
	fi := getMockFunc(s, s.GetTaskBoxDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.TaskBoxData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTaskBoxDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTaskTargetData(a0 uint64) *taskdata.TaskTargetData {
	fi := getMockFunc(s, s.GetTaskTargetData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.TaskTargetData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTaskTargetData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTaskTargetDataArray() []*taskdata.TaskTargetData {
	fi := getMockFunc(s, s.GetTaskTargetDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.TaskTargetData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTaskTargetDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTeachChapterData(a0 uint64) *teach.TeachChapterData {
	fi := getMockFunc(s, s.GetTeachChapterData)
	if fi != nil {
		f, ok := fi.(func(uint64) *teach.TeachChapterData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTeachChapterData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTeachChapterDataArray() []*teach.TeachChapterData {
	fi := getMockFunc(s, s.GetTeachChapterDataArray)
	if fi != nil {
		f, ok := fi.(func() []*teach.TeachChapterData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTeachChapterDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTechnologyData(a0 uint64) *domestic_data.TechnologyData {
	fi := getMockFunc(s, s.GetTechnologyData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.TechnologyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTechnologyData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTechnologyDataArray() []*domestic_data.TechnologyData {
	fi := getMockFunc(s, s.GetTechnologyDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.TechnologyData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTechnologyDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetText(a0 string) *data.Text {
	fi := getMockFunc(s, s.GetText)
	if fi != nil {
		f, ok := fi.(func(string) *data.Text)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetText()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTextArray() []*data.Text {
	fi := getMockFunc(s, s.GetTextArray)
	if fi != nil {
		f, ok := fi.(func() []*data.Text)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTextArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTieJiangPuLevelData(a0 uint64) *domestic_data.TieJiangPuLevelData {
	fi := getMockFunc(s, s.GetTieJiangPuLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.TieJiangPuLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTieJiangPuLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTieJiangPuLevelDataArray() []*domestic_data.TieJiangPuLevelData {
	fi := getMockFunc(s, s.GetTieJiangPuLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.TieJiangPuLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTieJiangPuLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTimeLimitGiftData(a0 uint64) *promdata.TimeLimitGiftData {
	fi := getMockFunc(s, s.GetTimeLimitGiftData)
	if fi != nil {
		f, ok := fi.(func(uint64) *promdata.TimeLimitGiftData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTimeLimitGiftData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTimeLimitGiftDataArray() []*promdata.TimeLimitGiftData {
	fi := getMockFunc(s, s.GetTimeLimitGiftDataArray)
	if fi != nil {
		f, ok := fi.(func() []*promdata.TimeLimitGiftData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTimeLimitGiftDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTimeLimitGiftGroupData(a0 uint64) *promdata.TimeLimitGiftGroupData {
	fi := getMockFunc(s, s.GetTimeLimitGiftGroupData)
	if fi != nil {
		f, ok := fi.(func(uint64) *promdata.TimeLimitGiftGroupData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTimeLimitGiftGroupData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTimeLimitGiftGroupDataArray() []*promdata.TimeLimitGiftGroupData {
	fi := getMockFunc(s, s.GetTimeLimitGiftGroupDataArray)
	if fi != nil {
		f, ok := fi.(func() []*promdata.TimeLimitGiftGroupData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTimeLimitGiftGroupDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTimeRuleData(a0 uint64) *data.TimeRuleData {
	fi := getMockFunc(s, s.GetTimeRuleData)
	if fi != nil {
		f, ok := fi.(func(uint64) *data.TimeRuleData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTimeRuleData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTimeRuleDataArray() []*data.TimeRuleData {
	fi := getMockFunc(s, s.GetTimeRuleDataArray)
	if fi != nil {
		f, ok := fi.(func() []*data.TimeRuleData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTimeRuleDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTitleData(a0 uint64) *taskdata.TitleData {
	fi := getMockFunc(s, s.GetTitleData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.TitleData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTitleData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTitleDataArray() []*taskdata.TitleData {
	fi := getMockFunc(s, s.GetTitleDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.TitleData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTitleDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTitleTaskData(a0 uint64) *taskdata.TitleTaskData {
	fi := getMockFunc(s, s.GetTitleTaskData)
	if fi != nil {
		f, ok := fi.(func(uint64) *taskdata.TitleTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTitleTaskData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTitleTaskDataArray() []*taskdata.TitleTaskData {
	fi := getMockFunc(s, s.GetTitleTaskDataArray)
	if fi != nil {
		f, ok := fi.(func() []*taskdata.TitleTaskData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTitleTaskDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTowerData(a0 uint64) *towerdata.TowerData {
	fi := getMockFunc(s, s.GetTowerData)
	if fi != nil {
		f, ok := fi.(func(uint64) *towerdata.TowerData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTowerData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTowerDataArray() []*towerdata.TowerData {
	fi := getMockFunc(s, s.GetTowerDataArray)
	if fi != nil {
		f, ok := fi.(func() []*towerdata.TowerData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTowerDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTrainingLevelData(a0 uint64) *military_data.TrainingLevelData {
	fi := getMockFunc(s, s.GetTrainingLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *military_data.TrainingLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTrainingLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTrainingLevelDataArray() []*military_data.TrainingLevelData {
	fi := getMockFunc(s, s.GetTrainingLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*military_data.TrainingLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTrainingLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTreasuryTreeData(a0 uint64) *gardendata.TreasuryTreeData {
	fi := getMockFunc(s, s.GetTreasuryTreeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *gardendata.TreasuryTreeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTreasuryTreeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTreasuryTreeDataArray() []*gardendata.TreasuryTreeData {
	fi := getMockFunc(s, s.GetTreasuryTreeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*gardendata.TreasuryTreeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTreasuryTreeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTroopDialogueData(a0 uint64) *regdata.TroopDialogueData {
	fi := getMockFunc(s, s.GetTroopDialogueData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.TroopDialogueData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTroopDialogueData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTroopDialogueDataArray() []*regdata.TroopDialogueData {
	fi := getMockFunc(s, s.GetTroopDialogueDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.TroopDialogueData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTroopDialogueDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTroopDialogueTextData(a0 uint64) *regdata.TroopDialogueTextData {
	fi := getMockFunc(s, s.GetTroopDialogueTextData)
	if fi != nil {
		f, ok := fi.(func(uint64) *regdata.TroopDialogueTextData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTroopDialogueTextData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTroopDialogueTextDataArray() []*regdata.TroopDialogueTextData {
	fi := getMockFunc(s, s.GetTroopDialogueTextDataArray)
	if fi != nil {
		f, ok := fi.(func() []*regdata.TroopDialogueTextData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTroopDialogueTextDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetTutorData(a0 uint64) *military_data.TutorData {
	fi := getMockFunc(s, s.GetTutorData)
	if fi != nil {
		f, ok := fi.(func(uint64) *military_data.TutorData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTutorData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetTutorDataArray() []*military_data.TutorData {
	fi := getMockFunc(s, s.GetTutorDataArray)
	if fi != nil {
		f, ok := fi.(func() []*military_data.TutorData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetTutorDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetVipContinueDaysData(a0 uint64) *vip.VipContinueDaysData {
	fi := getMockFunc(s, s.GetVipContinueDaysData)
	if fi != nil {
		f, ok := fi.(func(uint64) *vip.VipContinueDaysData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetVipContinueDaysData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetVipContinueDaysDataArray() []*vip.VipContinueDaysData {
	fi := getMockFunc(s, s.GetVipContinueDaysDataArray)
	if fi != nil {
		f, ok := fi.(func() []*vip.VipContinueDaysData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetVipContinueDaysDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetVipLevelData(a0 uint64) *vip.VipLevelData {
	fi := getMockFunc(s, s.GetVipLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *vip.VipLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetVipLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetVipLevelDataArray() []*vip.VipLevelData {
	fi := getMockFunc(s, s.GetVipLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*vip.VipLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetVipLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetWorkshopDuration(a0 uint64) *domestic_data.WorkshopDuration {
	fi := getMockFunc(s, s.GetWorkshopDuration)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.WorkshopDuration)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetWorkshopDuration()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetWorkshopDurationArray() []*domestic_data.WorkshopDuration {
	fi := getMockFunc(s, s.GetWorkshopDurationArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.WorkshopDuration)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetWorkshopDurationArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetWorkshopLevelData(a0 uint64) *domestic_data.WorkshopLevelData {
	fi := getMockFunc(s, s.GetWorkshopLevelData)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.WorkshopLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetWorkshopLevelData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetWorkshopLevelDataArray() []*domestic_data.WorkshopLevelData {
	fi := getMockFunc(s, s.GetWorkshopLevelDataArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.WorkshopLevelData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetWorkshopLevelDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetWorkshopRefreshCost(a0 uint64) *domestic_data.WorkshopRefreshCost {
	fi := getMockFunc(s, s.GetWorkshopRefreshCost)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.WorkshopRefreshCost)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetWorkshopRefreshCost()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetWorkshopRefreshCostArray() []*domestic_data.WorkshopRefreshCost {
	fi := getMockFunc(s, s.GetWorkshopRefreshCostArray)
	if fi != nil {
		f, ok := fi.(func() []*domestic_data.WorkshopRefreshCost)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetWorkshopRefreshCostArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetXuanyuanRangeData(a0 uint64) *xuanydata.XuanyuanRangeData {
	fi := getMockFunc(s, s.GetXuanyuanRangeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *xuanydata.XuanyuanRangeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetXuanyuanRangeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetXuanyuanRangeDataArray() []*xuanydata.XuanyuanRangeData {
	fi := getMockFunc(s, s.GetXuanyuanRangeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*xuanydata.XuanyuanRangeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetXuanyuanRangeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetXuanyuanRankPrizeData(a0 uint64) *xuanydata.XuanyuanRankPrizeData {
	fi := getMockFunc(s, s.GetXuanyuanRankPrizeData)
	if fi != nil {
		f, ok := fi.(func(uint64) *xuanydata.XuanyuanRankPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetXuanyuanRankPrizeData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetXuanyuanRankPrizeDataArray() []*xuanydata.XuanyuanRankPrizeData {
	fi := getMockFunc(s, s.GetXuanyuanRankPrizeDataArray)
	if fi != nil {
		f, ok := fi.(func() []*xuanydata.XuanyuanRankPrizeData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetXuanyuanRankPrizeDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetZhanJiangChapterData(a0 uint64) *zhanjiang.ZhanJiangChapterData {
	fi := getMockFunc(s, s.GetZhanJiangChapterData)
	if fi != nil {
		f, ok := fi.(func(uint64) *zhanjiang.ZhanJiangChapterData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhanJiangChapterData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetZhanJiangChapterDataArray() []*zhanjiang.ZhanJiangChapterData {
	fi := getMockFunc(s, s.GetZhanJiangChapterDataArray)
	if fi != nil {
		f, ok := fi.(func() []*zhanjiang.ZhanJiangChapterData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhanJiangChapterDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetZhanJiangData(a0 uint64) *zhanjiang.ZhanJiangData {
	fi := getMockFunc(s, s.GetZhanJiangData)
	if fi != nil {
		f, ok := fi.(func(uint64) *zhanjiang.ZhanJiangData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhanJiangData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetZhanJiangDataArray() []*zhanjiang.ZhanJiangData {
	fi := getMockFunc(s, s.GetZhanJiangDataArray)
	if fi != nil {
		f, ok := fi.(func() []*zhanjiang.ZhanJiangData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhanJiangDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetZhanJiangGuanQiaData(a0 uint64) *zhanjiang.ZhanJiangGuanQiaData {
	fi := getMockFunc(s, s.GetZhanJiangGuanQiaData)
	if fi != nil {
		f, ok := fi.(func(uint64) *zhanjiang.ZhanJiangGuanQiaData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhanJiangGuanQiaData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetZhanJiangGuanQiaDataArray() []*zhanjiang.ZhanJiangGuanQiaData {
	fi := getMockFunc(s, s.GetZhanJiangGuanQiaDataArray)
	if fi != nil {
		f, ok := fi.(func() []*zhanjiang.ZhanJiangGuanQiaData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhanJiangGuanQiaDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetZhenBaoGeShopGoods(a0 uint64) *shop.ZhenBaoGeShopGoods {
	fi := getMockFunc(s, s.GetZhenBaoGeShopGoods)
	if fi != nil {
		f, ok := fi.(func(uint64) *shop.ZhenBaoGeShopGoods)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhenBaoGeShopGoods()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetZhenBaoGeShopGoodsArray() []*shop.ZhenBaoGeShopGoods {
	fi := getMockFunc(s, s.GetZhenBaoGeShopGoodsArray)
	if fi != nil {
		f, ok := fi.(func() []*shop.ZhenBaoGeShopGoods)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhenBaoGeShopGoodsArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetZhengWuCompleteData(a0 uint64) *zhengwu.ZhengWuCompleteData {
	fi := getMockFunc(s, s.GetZhengWuCompleteData)
	if fi != nil {
		f, ok := fi.(func(uint64) *zhengwu.ZhengWuCompleteData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhengWuCompleteData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetZhengWuCompleteDataArray() []*zhengwu.ZhengWuCompleteData {
	fi := getMockFunc(s, s.GetZhengWuCompleteDataArray)
	if fi != nil {
		f, ok := fi.(func() []*zhengwu.ZhengWuCompleteData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhengWuCompleteDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetZhengWuData(a0 uint64) *zhengwu.ZhengWuData {
	fi := getMockFunc(s, s.GetZhengWuData)
	if fi != nil {
		f, ok := fi.(func(uint64) *zhengwu.ZhengWuData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhengWuData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetZhengWuDataArray() []*zhengwu.ZhengWuData {
	fi := getMockFunc(s, s.GetZhengWuDataArray)
	if fi != nil {
		f, ok := fi.(func() []*zhengwu.ZhengWuData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhengWuDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GetZhengWuRefreshData(a0 uint64) *zhengwu.ZhengWuRefreshData {
	fi := getMockFunc(s, s.GetZhengWuRefreshData)
	if fi != nil {
		f, ok := fi.(func(uint64) *zhengwu.ZhengWuRefreshData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhengWuRefreshData()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockConfigDatas) GetZhengWuRefreshDataArray() []*zhengwu.ZhengWuRefreshData {
	fi := getMockFunc(s, s.GetZhengWuRefreshDataArray)
	if fi != nil {
		f, ok := fi.(func() []*zhengwu.ZhengWuRefreshData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GetZhengWuRefreshDataArray()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GoodsCheck() *goods.GoodsCheck {
	fi := getMockFunc(s, s.GoodsCheck)
	if fi != nil {
		f, ok := fi.(func() *goods.GoodsCheck)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GoodsCheck()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GoodsCombineData() *config.GoodsCombineDataConfig {
	fi := getMockFunc(s, s.GoodsCombineData)
	if fi != nil {
		f, ok := fi.(func() *config.GoodsCombineDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GoodsCombineData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GoodsConfig() *singleton.GoodsConfig {
	fi := getMockFunc(s, s.GoodsConfig)
	if fi != nil {
		f, ok := fi.(func() *singleton.GoodsConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GoodsConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GoodsData() *config.GoodsDataConfig {
	fi := getMockFunc(s, s.GoodsData)
	if fi != nil {
		f, ok := fi.(func() *config.GoodsDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GoodsData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GoodsQuality() *config.GoodsQualityConfig {
	fi := getMockFunc(s, s.GoodsQuality)
	if fi != nil {
		f, ok := fi.(func() *config.GoodsQualityConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GoodsQuality()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuanFuLevelData() *config.GuanFuLevelDataConfig {
	fi := getMockFunc(s, s.GuanFuLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.GuanFuLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuanFuLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildBigBoxData() *config.GuildBigBoxDataConfig {
	fi := getMockFunc(s, s.GuildBigBoxData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildBigBoxDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildBigBoxData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildClassLevelData() *config.GuildClassLevelDataConfig {
	fi := getMockFunc(s, s.GuildClassLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildClassLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildClassLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildClassTitleData() *config.GuildClassTitleDataConfig {
	fi := getMockFunc(s, s.GuildClassTitleData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildClassTitleDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildClassTitleData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildConfig() *singleton.GuildConfig {
	fi := getMockFunc(s, s.GuildConfig)
	if fi != nil {
		f, ok := fi.(func() *singleton.GuildConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildDonateData() *config.GuildDonateDataConfig {
	fi := getMockFunc(s, s.GuildDonateData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildDonateDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildDonateData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildEventPrizeData() *config.GuildEventPrizeDataConfig {
	fi := getMockFunc(s, s.GuildEventPrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildEventPrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildEventPrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildGenConfig() *singleton.GuildGenConfig {
	fi := getMockFunc(s, s.GuildGenConfig)
	if fi != nil {
		f, ok := fi.(func() *singleton.GuildGenConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildGenConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildLevelCdrData() *config.GuildLevelCdrDataConfig {
	fi := getMockFunc(s, s.GuildLevelCdrData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildLevelCdrDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildLevelCdrData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildLevelData() *config.GuildLevelDataConfig {
	fi := getMockFunc(s, s.GuildLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildLevelPrize() *config.GuildLevelPrizeConfig {
	fi := getMockFunc(s, s.GuildLevelPrize)
	if fi != nil {
		f, ok := fi.(func() *config.GuildLevelPrizeConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildLevelPrize()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildLogData() *config.GuildLogDataConfig {
	fi := getMockFunc(s, s.GuildLogData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildLogDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildLogData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildLogHelp() *guild_data.GuildLogHelp {
	fi := getMockFunc(s, s.GuildLogHelp)
	if fi != nil {
		f, ok := fi.(func() *guild_data.GuildLogHelp)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildLogHelp()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildPermissionShowData() *config.GuildPermissionShowDataConfig {
	fi := getMockFunc(s, s.GuildPermissionShowData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildPermissionShowDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildPermissionShowData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildPrestigeEventData() *config.GuildPrestigeEventDataConfig {
	fi := getMockFunc(s, s.GuildPrestigeEventData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildPrestigeEventDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildPrestigeEventData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildPrestigePrizeData() *config.GuildPrestigePrizeDataConfig {
	fi := getMockFunc(s, s.GuildPrestigePrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildPrestigePrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildPrestigePrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildRankPrizeData() *config.GuildRankPrizeDataConfig {
	fi := getMockFunc(s, s.GuildRankPrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildRankPrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildRankPrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildTarget() *config.GuildTargetConfig {
	fi := getMockFunc(s, s.GuildTarget)
	if fi != nil {
		f, ok := fi.(func() *config.GuildTargetConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildTarget()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildTaskData() *config.GuildTaskDataConfig {
	fi := getMockFunc(s, s.GuildTaskData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildTaskDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildTaskData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildTaskEvaluateData() *config.GuildTaskEvaluateDataConfig {
	fi := getMockFunc(s, s.GuildTaskEvaluateData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildTaskEvaluateDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildTaskEvaluateData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) GuildTechnologyData() *config.GuildTechnologyDataConfig {
	fi := getMockFunc(s, s.GuildTechnologyData)
	if fi != nil {
		f, ok := fi.(func() *config.GuildTechnologyDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.GuildTechnologyData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) HeadData() *config.HeadDataConfig {
	fi := getMockFunc(s, s.HeadData)
	if fi != nil {
		f, ok := fi.(func() *config.HeadDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.HeadData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) HebiMiscData() *hebi.HebiMiscData {
	fi := getMockFunc(s, s.HebiMiscData)
	if fi != nil {
		f, ok := fi.(func() *hebi.HebiMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.HebiMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) HebiPrizeData() *config.HebiPrizeDataConfig {
	fi := getMockFunc(s, s.HebiPrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.HebiPrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.HebiPrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) HeroCreateData() *heroinit.HeroCreateData {
	fi := getMockFunc(s, s.HeroCreateData)
	if fi != nil {
		f, ok := fi.(func() *heroinit.HeroCreateData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.HeroCreateData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) HeroInitData() *heroinit.HeroInitData {
	fi := getMockFunc(s, s.HeroInitData)
	if fi != nil {
		f, ok := fi.(func() *heroinit.HeroInitData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.HeroInitData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) HeroLevelData() *config.HeroLevelDataConfig {
	fi := getMockFunc(s, s.HeroLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.HeroLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.HeroLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) HeroLevelFundData() *config.HeroLevelFundDataConfig {
	fi := getMockFunc(s, s.HeroLevelFundData)
	if fi != nil {
		f, ok := fi.(func() *config.HeroLevelFundDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.HeroLevelFundData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) HeroLevelSubData() *config.HeroLevelSubDataConfig {
	fi := getMockFunc(s, s.HeroLevelSubData)
	if fi != nil {
		f, ok := fi.(func() *config.HeroLevelSubDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.HeroLevelSubData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) HomeNpcBaseData() *config.HomeNpcBaseDataConfig {
	fi := getMockFunc(s, s.HomeNpcBaseData)
	if fi != nil {
		f, ok := fi.(func() *config.HomeNpcBaseDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.HomeNpcBaseData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) I18nData() *config.I18nDataConfig {
	fi := getMockFunc(s, s.I18nData)
	if fi != nil {
		f, ok := fi.(func() *config.I18nDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.I18nData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) Icon() *config.IconConfig {
	fi := getMockFunc(s, s.Icon)
	if fi != nil {
		f, ok := fi.(func() *config.IconConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.Icon()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) JiuGuanData() *config.JiuGuanDataConfig {
	fi := getMockFunc(s, s.JiuGuanData)
	if fi != nil {
		f, ok := fi.(func() *config.JiuGuanDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.JiuGuanData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) JiuGuanMiscData() *military_data.JiuGuanMiscData {
	fi := getMockFunc(s, s.JiuGuanMiscData)
	if fi != nil {
		f, ok := fi.(func() *military_data.JiuGuanMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.JiuGuanMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) JunTuanNpcData() *config.JunTuanNpcDataConfig {
	fi := getMockFunc(s, s.JunTuanNpcData)
	if fi != nil {
		f, ok := fi.(func() *config.JunTuanNpcDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.JunTuanNpcData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) JunTuanNpcPlaceConfig() *regdata.JunTuanNpcPlaceConfig {
	fi := getMockFunc(s, s.JunTuanNpcPlaceConfig)
	if fi != nil {
		f, ok := fi.(func() *regdata.JunTuanNpcPlaceConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.JunTuanNpcPlaceConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) JunTuanNpcPlaceData() *config.JunTuanNpcPlaceDataConfig {
	fi := getMockFunc(s, s.JunTuanNpcPlaceData)
	if fi != nil {
		f, ok := fi.(func() *config.JunTuanNpcPlaceDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.JunTuanNpcPlaceData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) JunXianLevelData() *config.JunXianLevelDataConfig {
	fi := getMockFunc(s, s.JunXianLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.JunXianLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.JunXianLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) JunXianPrizeData() *config.JunXianPrizeDataConfig {
	fi := getMockFunc(s, s.JunXianPrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.JunXianPrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.JunXianPrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) JunYingLevelData() *config.JunYingLevelDataConfig {
	fi := getMockFunc(s, s.JunYingLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.JunYingLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.JunYingLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) JunYingMiscData() *military_data.JunYingMiscData {
	fi := getMockFunc(s, s.JunYingMiscData)
	if fi != nil {
		f, ok := fi.(func() *military_data.JunYingMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.JunYingMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) LocationData() *config.LocationDataConfig {
	fi := getMockFunc(s, s.LocationData)
	if fi != nil {
		f, ok := fi.(func() *config.LocationDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.LocationData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) LoginDayData() *config.LoginDayDataConfig {
	fi := getMockFunc(s, s.LoginDayData)
	if fi != nil {
		f, ok := fi.(func() *config.LoginDayDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.LoginDayData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MailData() *config.MailDataConfig {
	fi := getMockFunc(s, s.MailData)
	if fi != nil {
		f, ok := fi.(func() *config.MailDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MailData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MailHelp() *maildata.MailHelp {
	fi := getMockFunc(s, s.MailHelp)
	if fi != nil {
		f, ok := fi.(func() *maildata.MailHelp)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MailHelp()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MainCityMiscData() *domestic_data.MainCityMiscData {
	fi := getMockFunc(s, s.MainCityMiscData)
	if fi != nil {
		f, ok := fi.(func() *domestic_data.MainCityMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MainCityMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MainTaskData() *config.MainTaskDataConfig {
	fi := getMockFunc(s, s.MainTaskData)
	if fi != nil {
		f, ok := fi.(func() *config.MainTaskDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MainTaskData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MaleGivenName() *config.MaleGivenNameConfig {
	fi := getMockFunc(s, s.MaleGivenName)
	if fi != nil {
		f, ok := fi.(func() *config.MaleGivenNameConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MaleGivenName()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) McBuildAddSupportData() *config.McBuildAddSupportDataConfig {
	fi := getMockFunc(s, s.McBuildAddSupportData)
	if fi != nil {
		f, ok := fi.(func() *config.McBuildAddSupportDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.McBuildAddSupportData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) McBuildGuildMemberPrizeData() *config.McBuildGuildMemberPrizeDataConfig {
	fi := getMockFunc(s, s.McBuildGuildMemberPrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.McBuildGuildMemberPrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.McBuildGuildMemberPrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) McBuildMcSupportData() *config.McBuildMcSupportDataConfig {
	fi := getMockFunc(s, s.McBuildMcSupportData)
	if fi != nil {
		f, ok := fi.(func() *config.McBuildMcSupportDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.McBuildMcSupportData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) McBuildMiscData() *mingcdata.McBuildMiscData {
	fi := getMockFunc(s, s.McBuildMiscData)
	if fi != nil {
		f, ok := fi.(func() *mingcdata.McBuildMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.McBuildMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MilitaryConfig() *singleton.MilitaryConfig {
	fi := getMockFunc(s, s.MilitaryConfig)
	if fi != nil {
		f, ok := fi.(func() *singleton.MilitaryConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MilitaryConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcBaseData() *config.MingcBaseDataConfig {
	fi := getMockFunc(s, s.MingcBaseData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcBaseDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcBaseData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcMiscData() *mingcdata.MingcMiscData {
	fi := getMockFunc(s, s.MingcMiscData)
	if fi != nil {
		f, ok := fi.(func() *mingcdata.MingcMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcTimeData() *config.MingcTimeDataConfig {
	fi := getMockFunc(s, s.MingcTimeData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcTimeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcTimeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcWarBuildingData() *config.MingcWarBuildingDataConfig {
	fi := getMockFunc(s, s.MingcWarBuildingData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcWarBuildingDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcWarBuildingData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcWarDrumStatData() *config.MingcWarDrumStatDataConfig {
	fi := getMockFunc(s, s.MingcWarDrumStatData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcWarDrumStatDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcWarDrumStatData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcWarMapData() *config.MingcWarMapDataConfig {
	fi := getMockFunc(s, s.MingcWarMapData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcWarMapDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcWarMapData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcWarMultiKillData() *config.MingcWarMultiKillDataConfig {
	fi := getMockFunc(s, s.MingcWarMultiKillData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcWarMultiKillDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcWarMultiKillData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcWarNpcData() *config.MingcWarNpcDataConfig {
	fi := getMockFunc(s, s.MingcWarNpcData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcWarNpcDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcWarNpcData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcWarNpcGuildData() *config.MingcWarNpcGuildDataConfig {
	fi := getMockFunc(s, s.MingcWarNpcGuildData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcWarNpcGuildDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcWarNpcGuildData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcWarSceneData() *config.MingcWarSceneDataConfig {
	fi := getMockFunc(s, s.MingcWarSceneData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcWarSceneDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcWarSceneData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcWarTouShiBuildingTargetData() *config.MingcWarTouShiBuildingTargetDataConfig {
	fi := getMockFunc(s, s.MingcWarTouShiBuildingTargetData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcWarTouShiBuildingTargetDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcWarTouShiBuildingTargetData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MingcWarTroopLastBeatWhenFailData() *config.MingcWarTroopLastBeatWhenFailDataConfig {
	fi := getMockFunc(s, s.MingcWarTroopLastBeatWhenFailData)
	if fi != nil {
		f, ok := fi.(func() *config.MingcWarTroopLastBeatWhenFailDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MingcWarTroopLastBeatWhenFailData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MiscConfig() *singleton.MiscConfig {
	fi := getMockFunc(s, s.MiscConfig)
	if fi != nil {
		f, ok := fi.(func() *singleton.MiscConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MiscConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MiscGenConfig() *singleton.MiscGenConfig {
	fi := getMockFunc(s, s.MiscGenConfig)
	if fi != nil {
		f, ok := fi.(func() *singleton.MiscGenConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MiscGenConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MonsterCaptainData() *config.MonsterCaptainDataConfig {
	fi := getMockFunc(s, s.MonsterCaptainData)
	if fi != nil {
		f, ok := fi.(func() *config.MonsterCaptainDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MonsterCaptainData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) MonsterMasterData() *config.MonsterMasterDataConfig {
	fi := getMockFunc(s, s.MonsterMasterData)
	if fi != nil {
		f, ok := fi.(func() *config.MonsterMasterDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.MonsterMasterData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) NamelessCaptainData() *config.NamelessCaptainDataConfig {
	fi := getMockFunc(s, s.NamelessCaptainData)
	if fi != nil {
		f, ok := fi.(func() *config.NamelessCaptainDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.NamelessCaptainData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) NormalShopGoods() *config.NormalShopGoodsConfig {
	fi := getMockFunc(s, s.NormalShopGoods)
	if fi != nil {
		f, ok := fi.(func() *config.NormalShopGoodsConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.NormalShopGoods()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) NpcBaseData() *config.NpcBaseDataConfig {
	fi := getMockFunc(s, s.NpcBaseData)
	if fi != nil {
		f, ok := fi.(func() *config.NpcBaseDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.NpcBaseData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) NpcGuildSuffixName() *guild_data.NpcGuildSuffixName {
	fi := getMockFunc(s, s.NpcGuildSuffixName)
	if fi != nil {
		f, ok := fi.(func() *guild_data.NpcGuildSuffixName)
		if !ok {
			panic("invalid mock func, MockConfigDatas.NpcGuildSuffixName()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) NpcGuildTemplate() *config.NpcGuildTemplateConfig {
	fi := getMockFunc(s, s.NpcGuildTemplate)
	if fi != nil {
		f, ok := fi.(func() *config.NpcGuildTemplateConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.NpcGuildTemplate()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) NpcMemberData() *config.NpcMemberDataConfig {
	fi := getMockFunc(s, s.NpcMemberData)
	if fi != nil {
		f, ok := fi.(func() *config.NpcMemberDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.NpcMemberData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) OptionPrize() *config.OptionPrizeConfig {
	fi := getMockFunc(s, s.OptionPrize)
	if fi != nil {
		f, ok := fi.(func() *config.OptionPrizeConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.OptionPrize()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) OuterCityBuildingData() *config.OuterCityBuildingDataConfig {
	fi := getMockFunc(s, s.OuterCityBuildingData)
	if fi != nil {
		f, ok := fi.(func() *config.OuterCityBuildingDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.OuterCityBuildingData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) OuterCityBuildingDescData() *config.OuterCityBuildingDescDataConfig {
	fi := getMockFunc(s, s.OuterCityBuildingDescData)
	if fi != nil {
		f, ok := fi.(func() *config.OuterCityBuildingDescDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.OuterCityBuildingDescData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) OuterCityData() *config.OuterCityDataConfig {
	fi := getMockFunc(s, s.OuterCityData)
	if fi != nil {
		f, ok := fi.(func() *config.OuterCityDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.OuterCityData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) OuterCityLayoutData() *config.OuterCityLayoutDataConfig {
	fi := getMockFunc(s, s.OuterCityLayoutData)
	if fi != nil {
		f, ok := fi.(func() *config.OuterCityLayoutDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.OuterCityLayoutData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) PassiveSpellData() *config.PassiveSpellDataConfig {
	fi := getMockFunc(s, s.PassiveSpellData)
	if fi != nil {
		f, ok := fi.(func() *config.PassiveSpellDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.PassiveSpellData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) Plunder() *config.PlunderConfig {
	fi := getMockFunc(s, s.Plunder)
	if fi != nil {
		f, ok := fi.(func() *config.PlunderConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.Plunder()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) PlunderGroup() *config.PlunderGroupConfig {
	fi := getMockFunc(s, s.PlunderGroup)
	if fi != nil {
		f, ok := fi.(func() *config.PlunderGroupConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.PlunderGroup()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) PlunderItem() *config.PlunderItemConfig {
	fi := getMockFunc(s, s.PlunderItem)
	if fi != nil {
		f, ok := fi.(func() *config.PlunderItemConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.PlunderItem()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) PlunderPrize() *config.PlunderPrizeConfig {
	fi := getMockFunc(s, s.PlunderPrize)
	if fi != nil {
		f, ok := fi.(func() *config.PlunderPrizeConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.PlunderPrize()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) PrivacySettingData() *config.PrivacySettingDataConfig {
	fi := getMockFunc(s, s.PrivacySettingData)
	if fi != nil {
		f, ok := fi.(func() *config.PrivacySettingDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.PrivacySettingData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) Prize() *config.PrizeConfig {
	fi := getMockFunc(s, s.Prize)
	if fi != nil {
		f, ok := fi.(func() *config.PrizeConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.Prize()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ProductData() *config.ProductDataConfig {
	fi := getMockFunc(s, s.ProductData)
	if fi != nil {
		f, ok := fi.(func() *config.ProductDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ProductData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) PromotionMiscData() *promdata.PromotionMiscData {
	fi := getMockFunc(s, s.PromotionMiscData)
	if fi != nil {
		f, ok := fi.(func() *promdata.PromotionMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.PromotionMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ProsperityDamageBuffData() *config.ProsperityDamageBuffDataConfig {
	fi := getMockFunc(s, s.ProsperityDamageBuffData)
	if fi != nil {
		f, ok := fi.(func() *config.ProsperityDamageBuffDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ProsperityDamageBuffData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) PushData() *config.PushDataConfig {
	fi := getMockFunc(s, s.PushData)
	if fi != nil {
		f, ok := fi.(func() *config.PushDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.PushData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) PveTroopData() *config.PveTroopDataConfig {
	fi := getMockFunc(s, s.PveTroopData)
	if fi != nil {
		f, ok := fi.(func() *config.PveTroopDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.PveTroopData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) QuestionData() *config.QuestionDataConfig {
	fi := getMockFunc(s, s.QuestionData)
	if fi != nil {
		f, ok := fi.(func() *config.QuestionDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.QuestionData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) QuestionMiscData() *question.QuestionMiscData {
	fi := getMockFunc(s, s.QuestionMiscData)
	if fi != nil {
		f, ok := fi.(func() *question.QuestionMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.QuestionMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) QuestionPrizeData() *config.QuestionPrizeDataConfig {
	fi := getMockFunc(s, s.QuestionPrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.QuestionPrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.QuestionPrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) QuestionSayingData() *config.QuestionSayingDataConfig {
	fi := getMockFunc(s, s.QuestionSayingData)
	if fi != nil {
		f, ok := fi.(func() *config.QuestionSayingDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.QuestionSayingData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RaceConfig() *race.RaceConfig {
	fi := getMockFunc(s, s.RaceConfig)
	if fi != nil {
		f, ok := fi.(func() *race.RaceConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RaceConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RaceData() *config.RaceDataConfig {
	fi := getMockFunc(s, s.RaceData)
	if fi != nil {
		f, ok := fi.(func() *config.RaceDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RaceData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RandomEventData() *config.RandomEventDataConfig {
	fi := getMockFunc(s, s.RandomEventData)
	if fi != nil {
		f, ok := fi.(func() *config.RandomEventDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RandomEventData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RandomEventDataDictionary() *random_event.RandomEventDataDictionary {
	fi := getMockFunc(s, s.RandomEventDataDictionary)
	if fi != nil {
		f, ok := fi.(func() *random_event.RandomEventDataDictionary)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RandomEventDataDictionary()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RandomEventPositionDictionary() *random_event.RandomEventPositionDictionary {
	fi := getMockFunc(s, s.RandomEventPositionDictionary)
	if fi != nil {
		f, ok := fi.(func() *random_event.RandomEventPositionDictionary)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RandomEventPositionDictionary()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RankMiscData() *rank_data.RankMiscData {
	fi := getMockFunc(s, s.RankMiscData)
	if fi != nil {
		f, ok := fi.(func() *rank_data.RankMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RankMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RedPacketData() *config.RedPacketDataConfig {
	fi := getMockFunc(s, s.RedPacketData)
	if fi != nil {
		f, ok := fi.(func() *config.RedPacketDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RedPacketData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RegionAreaData() *config.RegionAreaDataConfig {
	fi := getMockFunc(s, s.RegionAreaData)
	if fi != nil {
		f, ok := fi.(func() *config.RegionAreaDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RegionAreaData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RegionConfig() *singleton.RegionConfig {
	fi := getMockFunc(s, s.RegionConfig)
	if fi != nil {
		f, ok := fi.(func() *singleton.RegionConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RegionConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RegionData() *config.RegionDataConfig {
	fi := getMockFunc(s, s.RegionData)
	if fi != nil {
		f, ok := fi.(func() *config.RegionDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RegionData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RegionGenConfig() *singleton.RegionGenConfig {
	fi := getMockFunc(s, s.RegionGenConfig)
	if fi != nil {
		f, ok := fi.(func() *singleton.RegionGenConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RegionGenConfig()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RegionMonsterData() *config.RegionMonsterDataConfig {
	fi := getMockFunc(s, s.RegionMonsterData)
	if fi != nil {
		f, ok := fi.(func() *config.RegionMonsterDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RegionMonsterData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RegionMultiLevelNpcData() *config.RegionMultiLevelNpcDataConfig {
	fi := getMockFunc(s, s.RegionMultiLevelNpcData)
	if fi != nil {
		f, ok := fi.(func() *config.RegionMultiLevelNpcDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RegionMultiLevelNpcData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RegionMultiLevelNpcLevelData() *config.RegionMultiLevelNpcLevelDataConfig {
	fi := getMockFunc(s, s.RegionMultiLevelNpcLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.RegionMultiLevelNpcLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RegionMultiLevelNpcLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) RegionMultiLevelNpcTypeData() *config.RegionMultiLevelNpcTypeDataConfig {
	fi := getMockFunc(s, s.RegionMultiLevelNpcTypeData)
	if fi != nil {
		f, ok := fi.(func() *config.RegionMultiLevelNpcTypeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.RegionMultiLevelNpcTypeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ResCaptainData() *config.ResCaptainDataConfig {
	fi := getMockFunc(s, s.ResCaptainData)
	if fi != nil {
		f, ok := fi.(func() *config.ResCaptainDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ResCaptainData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ResistXiongNuData() *config.ResistXiongNuDataConfig {
	fi := getMockFunc(s, s.ResistXiongNuData)
	if fi != nil {
		f, ok := fi.(func() *config.ResistXiongNuDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ResistXiongNuData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ResistXiongNuMisc() *xiongnu.ResistXiongNuMisc {
	fi := getMockFunc(s, s.ResistXiongNuMisc)
	if fi != nil {
		f, ok := fi.(func() *xiongnu.ResistXiongNuMisc)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ResistXiongNuMisc()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ResistXiongNuScoreData() *config.ResistXiongNuScoreDataConfig {
	fi := getMockFunc(s, s.ResistXiongNuScoreData)
	if fi != nil {
		f, ok := fi.(func() *config.ResistXiongNuScoreDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ResistXiongNuScoreData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ResistXiongNuWaveData() *config.ResistXiongNuWaveDataConfig {
	fi := getMockFunc(s, s.ResistXiongNuWaveData)
	if fi != nil {
		f, ok := fi.(func() *config.ResistXiongNuWaveDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ResistXiongNuWaveData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SeasonData() *config.SeasonDataConfig {
	fi := getMockFunc(s, s.SeasonData)
	if fi != nil {
		f, ok := fi.(func() *config.SeasonDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SeasonData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SeasonMiscData() *season.SeasonMiscData {
	fi := getMockFunc(s, s.SeasonMiscData)
	if fi != nil {
		f, ok := fi.(func() *season.SeasonMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SeasonMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SecretTowerData() *config.SecretTowerDataConfig {
	fi := getMockFunc(s, s.SecretTowerData)
	if fi != nil {
		f, ok := fi.(func() *config.SecretTowerDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SecretTowerData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SecretTowerMiscData() *towerdata.SecretTowerMiscData {
	fi := getMockFunc(s, s.SecretTowerMiscData)
	if fi != nil {
		f, ok := fi.(func() *towerdata.SecretTowerMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SecretTowerMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SecretTowerWordsData() *config.SecretTowerWordsDataConfig {
	fi := getMockFunc(s, s.SecretTowerWordsData)
	if fi != nil {
		f, ok := fi.(func() *config.SecretTowerWordsDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SecretTowerWordsData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SettingMiscData() *settings.SettingMiscData {
	fi := getMockFunc(s, s.SettingMiscData)
	if fi != nil {
		f, ok := fi.(func() *settings.SettingMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SettingMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) Shop() *config.ShopConfig {
	fi := getMockFunc(s, s.Shop)
	if fi != nil {
		f, ok := fi.(func() *config.ShopConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.Shop()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ShopMiscData() *shop.ShopMiscData {
	fi := getMockFunc(s, s.ShopMiscData)
	if fi != nil {
		f, ok := fi.(func() *shop.ShopMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ShopMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SoldierLevelData() *config.SoldierLevelDataConfig {
	fi := getMockFunc(s, s.SoldierLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.SoldierLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SoldierLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SpCollectionData() *config.SpCollectionDataConfig {
	fi := getMockFunc(s, s.SpCollectionData)
	if fi != nil {
		f, ok := fi.(func() *config.SpCollectionDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SpCollectionData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) Spell() *config.SpellConfig {
	fi := getMockFunc(s, s.Spell)
	if fi != nil {
		f, ok := fi.(func() *config.SpellConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.Spell()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SpellData() *config.SpellDataConfig {
	fi := getMockFunc(s, s.SpellData)
	if fi != nil {
		f, ok := fi.(func() *config.SpellDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SpellData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SpellFacadeData() *config.SpellFacadeDataConfig {
	fi := getMockFunc(s, s.SpellFacadeData)
	if fi != nil {
		f, ok := fi.(func() *config.SpellFacadeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SpellFacadeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SpriteStat() *config.SpriteStatConfig {
	fi := getMockFunc(s, s.SpriteStat)
	if fi != nil {
		f, ok := fi.(func() *config.SpriteStatConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SpriteStat()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) StateData() *config.StateDataConfig {
	fi := getMockFunc(s, s.StateData)
	if fi != nil {
		f, ok := fi.(func() *config.StateDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.StateData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) StrategyData() *config.StrategyDataConfig {
	fi := getMockFunc(s, s.StrategyData)
	if fi != nil {
		f, ok := fi.(func() *config.StrategyDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.StrategyData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) StrategyEffectData() *config.StrategyEffectDataConfig {
	fi := getMockFunc(s, s.StrategyEffectData)
	if fi != nil {
		f, ok := fi.(func() *config.StrategyEffectDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.StrategyEffectData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) StrongerData() *config.StrongerDataConfig {
	fi := getMockFunc(s, s.StrongerData)
	if fi != nil {
		f, ok := fi.(func() *config.StrongerDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.StrongerData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) SurveyData() *config.SurveyDataConfig {
	fi := getMockFunc(s, s.SurveyData)
	if fi != nil {
		f, ok := fi.(func() *config.SurveyDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.SurveyData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TagMiscData() *tag.TagMiscData {
	fi := getMockFunc(s, s.TagMiscData)
	if fi != nil {
		f, ok := fi.(func() *tag.TagMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TagMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TaskBoxData() *config.TaskBoxDataConfig {
	fi := getMockFunc(s, s.TaskBoxData)
	if fi != nil {
		f, ok := fi.(func() *config.TaskBoxDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TaskBoxData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TaskMiscData() *taskdata.TaskMiscData {
	fi := getMockFunc(s, s.TaskMiscData)
	if fi != nil {
		f, ok := fi.(func() *taskdata.TaskMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TaskMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TaskTargetData() *config.TaskTargetDataConfig {
	fi := getMockFunc(s, s.TaskTargetData)
	if fi != nil {
		f, ok := fi.(func() *config.TaskTargetDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TaskTargetData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TeachChapterData() *config.TeachChapterDataConfig {
	fi := getMockFunc(s, s.TeachChapterData)
	if fi != nil {
		f, ok := fi.(func() *config.TeachChapterDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TeachChapterData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TechnologyData() *config.TechnologyDataConfig {
	fi := getMockFunc(s, s.TechnologyData)
	if fi != nil {
		f, ok := fi.(func() *config.TechnologyDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TechnologyData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) Text() *config.TextConfig {
	fi := getMockFunc(s, s.Text)
	if fi != nil {
		f, ok := fi.(func() *config.TextConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.Text()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TextHelp() *data.TextHelp {
	fi := getMockFunc(s, s.TextHelp)
	if fi != nil {
		f, ok := fi.(func() *data.TextHelp)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TextHelp()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TieJiangPuLevelData() *config.TieJiangPuLevelDataConfig {
	fi := getMockFunc(s, s.TieJiangPuLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.TieJiangPuLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TieJiangPuLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TimeLimitGiftData() *config.TimeLimitGiftDataConfig {
	fi := getMockFunc(s, s.TimeLimitGiftData)
	if fi != nil {
		f, ok := fi.(func() *config.TimeLimitGiftDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TimeLimitGiftData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TimeLimitGiftGroupData() *config.TimeLimitGiftGroupDataConfig {
	fi := getMockFunc(s, s.TimeLimitGiftGroupData)
	if fi != nil {
		f, ok := fi.(func() *config.TimeLimitGiftGroupDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TimeLimitGiftGroupData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TimeRuleData() *config.TimeRuleDataConfig {
	fi := getMockFunc(s, s.TimeRuleData)
	if fi != nil {
		f, ok := fi.(func() *config.TimeRuleDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TimeRuleData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TitleData() *config.TitleDataConfig {
	fi := getMockFunc(s, s.TitleData)
	if fi != nil {
		f, ok := fi.(func() *config.TitleDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TitleData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TitleTaskData() *config.TitleTaskDataConfig {
	fi := getMockFunc(s, s.TitleTaskData)
	if fi != nil {
		f, ok := fi.(func() *config.TitleTaskDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TitleTaskData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TowerData() *config.TowerDataConfig {
	fi := getMockFunc(s, s.TowerData)
	if fi != nil {
		f, ok := fi.(func() *config.TowerDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TowerData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TrainingLevelData() *config.TrainingLevelDataConfig {
	fi := getMockFunc(s, s.TrainingLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.TrainingLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TrainingLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TreasuryTreeData() *config.TreasuryTreeDataConfig {
	fi := getMockFunc(s, s.TreasuryTreeData)
	if fi != nil {
		f, ok := fi.(func() *config.TreasuryTreeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TreasuryTreeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TroopDialogueData() *config.TroopDialogueDataConfig {
	fi := getMockFunc(s, s.TroopDialogueData)
	if fi != nil {
		f, ok := fi.(func() *config.TroopDialogueDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TroopDialogueData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TroopDialogueTextData() *config.TroopDialogueTextDataConfig {
	fi := getMockFunc(s, s.TroopDialogueTextData)
	if fi != nil {
		f, ok := fi.(func() *config.TroopDialogueTextDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TroopDialogueTextData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) TutorData() *config.TutorDataConfig {
	fi := getMockFunc(s, s.TutorData)
	if fi != nil {
		f, ok := fi.(func() *config.TutorDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.TutorData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) VipContinueDaysData() *config.VipContinueDaysDataConfig {
	fi := getMockFunc(s, s.VipContinueDaysData)
	if fi != nil {
		f, ok := fi.(func() *config.VipContinueDaysDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.VipContinueDaysData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) VipLevelData() *config.VipLevelDataConfig {
	fi := getMockFunc(s, s.VipLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.VipLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.VipLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) VipMiscData() *vip.VipMiscData {
	fi := getMockFunc(s, s.VipMiscData)
	if fi != nil {
		f, ok := fi.(func() *vip.VipMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.VipMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) WorkshopDuration() *config.WorkshopDurationConfig {
	fi := getMockFunc(s, s.WorkshopDuration)
	if fi != nil {
		f, ok := fi.(func() *config.WorkshopDurationConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.WorkshopDuration()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) WorkshopLevelData() *config.WorkshopLevelDataConfig {
	fi := getMockFunc(s, s.WorkshopLevelData)
	if fi != nil {
		f, ok := fi.(func() *config.WorkshopLevelDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.WorkshopLevelData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) WorkshopRefreshCost() *config.WorkshopRefreshCostConfig {
	fi := getMockFunc(s, s.WorkshopRefreshCost)
	if fi != nil {
		f, ok := fi.(func() *config.WorkshopRefreshCostConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.WorkshopRefreshCost()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) XuanyuanMiscData() *xuanydata.XuanyuanMiscData {
	fi := getMockFunc(s, s.XuanyuanMiscData)
	if fi != nil {
		f, ok := fi.(func() *xuanydata.XuanyuanMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.XuanyuanMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) XuanyuanRangeData() *config.XuanyuanRangeDataConfig {
	fi := getMockFunc(s, s.XuanyuanRangeData)
	if fi != nil {
		f, ok := fi.(func() *config.XuanyuanRangeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.XuanyuanRangeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) XuanyuanRankPrizeData() *config.XuanyuanRankPrizeDataConfig {
	fi := getMockFunc(s, s.XuanyuanRankPrizeData)
	if fi != nil {
		f, ok := fi.(func() *config.XuanyuanRankPrizeDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.XuanyuanRankPrizeData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ZhanJiangChapterData() *config.ZhanJiangChapterDataConfig {
	fi := getMockFunc(s, s.ZhanJiangChapterData)
	if fi != nil {
		f, ok := fi.(func() *config.ZhanJiangChapterDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ZhanJiangChapterData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ZhanJiangData() *config.ZhanJiangDataConfig {
	fi := getMockFunc(s, s.ZhanJiangData)
	if fi != nil {
		f, ok := fi.(func() *config.ZhanJiangDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ZhanJiangData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ZhanJiangGuanQiaData() *config.ZhanJiangGuanQiaDataConfig {
	fi := getMockFunc(s, s.ZhanJiangGuanQiaData)
	if fi != nil {
		f, ok := fi.(func() *config.ZhanJiangGuanQiaDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ZhanJiangGuanQiaData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ZhanJiangMiscData() *zhanjiang.ZhanJiangMiscData {
	fi := getMockFunc(s, s.ZhanJiangMiscData)
	if fi != nil {
		f, ok := fi.(func() *zhanjiang.ZhanJiangMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ZhanJiangMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ZhenBaoGeShopGoods() *config.ZhenBaoGeShopGoodsConfig {
	fi := getMockFunc(s, s.ZhenBaoGeShopGoods)
	if fi != nil {
		f, ok := fi.(func() *config.ZhenBaoGeShopGoodsConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ZhenBaoGeShopGoods()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ZhengWuCompleteData() *config.ZhengWuCompleteDataConfig {
	fi := getMockFunc(s, s.ZhengWuCompleteData)
	if fi != nil {
		f, ok := fi.(func() *config.ZhengWuCompleteDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ZhengWuCompleteData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ZhengWuData() *config.ZhengWuDataConfig {
	fi := getMockFunc(s, s.ZhengWuData)
	if fi != nil {
		f, ok := fi.(func() *config.ZhengWuDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ZhengWuData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ZhengWuMiscData() *zhengwu.ZhengWuMiscData {
	fi := getMockFunc(s, s.ZhengWuMiscData)
	if fi != nil {
		f, ok := fi.(func() *zhengwu.ZhengWuMiscData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ZhengWuMiscData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ZhengWuRandomData() *zhengwu.ZhengWuRandomData {
	fi := getMockFunc(s, s.ZhengWuRandomData)
	if fi != nil {
		f, ok := fi.(func() *zhengwu.ZhengWuRandomData)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ZhengWuRandomData()")
		}
		return f()
	}

	return nil
}
func (s *MockConfigDatas) ZhengWuRefreshData() *config.ZhengWuRefreshDataConfig {
	fi := getMockFunc(s, s.ZhengWuRefreshData)
	if fi != nil {
		f, ok := fi.(func() *config.ZhengWuRefreshDataConfig)
		if !ok {
			panic("invalid mock func, MockConfigDatas.ZhengWuRefreshData()")
		}
		return f()
	}

	return nil
}

// 已在线的用户

var ConnectedUser = &MockConnectedUser{}

type MockConnectedUser struct{}

func (s *MockConnectedUser) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockConnectedUser) Disconnect(a0 msg.ErrMsg) {
	fi := getMockFunc(s, s.Disconnect)
	if fi != nil {
		f, ok := fi.(func(msg.ErrMsg))
		if !ok {
			panic("invalid mock func, MockConnectedUser.Disconnect()")
		}
		f(a0)
	}

}
func (s *MockConnectedUser) DisconnectAndWait(a0 msg.ErrMsg) {
	fi := getMockFunc(s, s.DisconnectAndWait)
	if fi != nil {
		f, ok := fi.(func(msg.ErrMsg))
		if !ok {
			panic("invalid mock func, MockConnectedUser.DisconnectAndWait()")
		}
		f(a0)
	}

}
func (s *MockConnectedUser) GetHeroController() iface.HeroController {
	fi := getMockFunc(s, s.GetHeroController)
	if fi != nil {
		f, ok := fi.(func() iface.HeroController)
		if !ok {
			panic("invalid mock func, MockConnectedUser.GetHeroController()")
		}
		return f()
	}

	return nil
}
func (s *MockConnectedUser) Id() int64 {
	fi := getMockFunc(s, s.Id)
	if fi != nil {
		f, ok := fi.(func() int64)
		if !ok {
			panic("invalid mock func, MockConnectedUser.Id()")
		}
		return f()
	}

	return 0
}
func (s *MockConnectedUser) IsClosed() bool {
	fi := getMockFunc(s, s.IsClosed)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockConnectedUser.IsClosed()")
		}
		return f()
	}

	return false
}
func (s *MockConnectedUser) IsLoaded() bool {
	fi := getMockFunc(s, s.IsLoaded)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockConnectedUser.IsLoaded()")
		}
		return f()
	}

	return false
}
func (s *MockConnectedUser) LogoutType() uint64 {
	fi := getMockFunc(s, s.LogoutType)
	if fi != nil {
		f, ok := fi.(func() uint64)
		if !ok {
			panic("invalid mock func, MockConnectedUser.LogoutType()")
		}
		return f()
	}

	return 0
}

// 玩家杂项
func (s *MockConnectedUser) Misc() *server_proto.UserMiscProto {
	fi := getMockFunc(s, s.Misc)
	if fi != nil {
		f, ok := fi.(func() *server_proto.UserMiscProto)
		if !ok {
			panic("invalid mock func, MockConnectedUser.Misc()")
		}
		return f()
	}

	return nil
}

// 玩家杂项
func (s *MockConnectedUser) MiscNeedOfflineSave() bool {
	fi := getMockFunc(s, s.MiscNeedOfflineSave)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockConnectedUser.MiscNeedOfflineSave()")
		}
		return f()
	}

	return false
}

// 发送消息.
func (s *MockConnectedUser) Send(a0 pbutil.Buffer) {
	fi := getMockFunc(s, s.Send)
	if fi != nil {
		f, ok := fi.(func(pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockConnectedUser.Send()")
		}
		f(a0)
	}

}

// 发送在线路繁忙时可以被丢掉的消息
func (s *MockConnectedUser) SendAll(a0 []pbutil.Buffer) {
	fi := getMockFunc(s, s.SendAll)
	if fi != nil {
		f, ok := fi.(func([]pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockConnectedUser.SendAll()")
		}
		f(a0)
	}

}

// 发送在线路繁忙时可以被丢掉的消息
func (s *MockConnectedUser) SendIfFree(a0 pbutil.Buffer) {
	fi := getMockFunc(s, s.SendIfFree)
	if fi != nil {
		f, ok := fi.(func(pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockConnectedUser.SendIfFree()")
		}
		f(a0)
	}

}
func (s *MockConnectedUser) SetHeroController(a0 iface.HeroController) {
	fi := getMockFunc(s, s.SetHeroController)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockConnectedUser.SetHeroController()")
		}
		f(a0)
	}

}
func (s *MockConnectedUser) SetLoaded() {
	fi := getMockFunc(s, s.SetLoaded)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockConnectedUser.SetLoaded()")
		}
		f()
	}

}
func (s *MockConnectedUser) SetLogoutType(a0 uint64) {
	fi := getMockFunc(s, s.SetLogoutType)
	if fi != nil {
		f, ok := fi.(func(uint64))
		if !ok {
			panic("invalid mock func, MockConnectedUser.SetLogoutType()")
		}
		f(a0)
	}

}

// 玩家杂项
func (s *MockConnectedUser) SetMisc(a0 *server_proto.UserMiscProto) {
	fi := getMockFunc(s, s.SetMisc)
	if fi != nil {
		f, ok := fi.(func(*server_proto.UserMiscProto))
		if !ok {
			panic("invalid mock func, MockConnectedUser.SetMisc()")
		}
		f(a0)
	}

}
func (s *MockConnectedUser) Sid() uint32 {
	fi := getMockFunc(s, s.Sid)
	if fi != nil {
		f, ok := fi.(func() uint32)
		if !ok {
			panic("invalid mock func, MockConnectedUser.Sid()")
		}
		return f()
	}

	return 0
}
func (s *MockConnectedUser) TencentInfo() *shared_proto.TencentInfoProto {
	fi := getMockFunc(s, s.TencentInfo)
	if fi != nil {
		f, ok := fi.(func() *shared_proto.TencentInfoProto)
		if !ok {
			panic("invalid mock func, MockConnectedUser.TencentInfo()")
		}
		return f()
	}

	return nil
}

var CountryModule = &MockCountryModule{}

type MockCountryModule struct{}

func (s *MockCountryModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var CountryService = &MockCountryService{}

type MockCountryService struct{}

func (s *MockCountryService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockCountryService) AddPrestige(a0 uint64, a1 uint64) (uint64, bool) {
	fi := getMockFunc(s, s.AddPrestige)
	if fi != nil {
		f, ok := fi.(func(uint64, uint64) (uint64, bool))
		if !ok {
			panic("invalid mock func, MockCountryService.AddPrestige()")
		}
		return f(a0, a1)
	}

	return 0, false
}
func (s *MockCountryService) AfterChangeCountry(a0 int64, a1 uint64, a2 uint64, a3 int32, a4 int, a5 bool) {
	fi := getMockFunc(s, s.AfterChangeCountry)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, int32, int, bool))
		if !ok {
			panic("invalid mock func, MockCountryService.AfterChangeCountry()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockCountryService) AfterChangeNameVoteStart(a0 *entity.Country) {
	fi := getMockFunc(s, s.AfterChangeNameVoteStart)
	if fi != nil {
		f, ok := fi.(func(*entity.Country))
		if !ok {
			panic("invalid mock func, MockCountryService.AfterChangeNameVoteStart()")
		}
		f(a0)
	}

}
func (s *MockCountryService) AfterUpgradeTitle(a0 int64, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.AfterUpgradeTitle)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockCountryService.AfterUpgradeTitle()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockCountryService) BroadcastCountry(a0 pbutil.Buffer, a1 uint64) {
	fi := getMockFunc(s, s.BroadcastCountry)
	if fi != nil {
		f, ok := fi.(func(pbutil.Buffer, uint64))
		if !ok {
			panic("invalid mock func, MockCountryService.BroadcastCountry()")
		}
		f(a0, a1)
	}

}
func (s *MockCountryService) CancelCountryDestroy(a0 uint64) {
	fi := getMockFunc(s, s.CancelCountryDestroy)
	if fi != nil {
		f, ok := fi.(func(uint64))
		if !ok {
			panic("invalid mock func, MockCountryService.CancelCountryDestroy()")
		}
		f(a0)
	}

}
func (s *MockCountryService) ChangeCountryHost(a0 uint64, a1 int64) bool {
	fi := getMockFunc(s, s.ChangeCountryHost)
	if fi != nil {
		f, ok := fi.(func(uint64, int64) bool)
		if !ok {
			panic("invalid mock func, MockCountryService.ChangeCountryHost()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockCountryService) ChangeKing(a0 uint64, a1 int64) bool {
	fi := getMockFunc(s, s.ChangeKing)
	if fi != nil {
		f, ok := fi.(func(uint64, int64) bool)
		if !ok {
			panic("invalid mock func, MockCountryService.ChangeKing()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockCountryService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockCountryService.Close()")
		}
		f()
	}

}
func (s *MockCountryService) Countries() []*entity.Country {
	fi := getMockFunc(s, s.Countries)
	if fi != nil {
		f, ok := fi.(func() []*entity.Country)
		if !ok {
			panic("invalid mock func, MockCountryService.Countries()")
		}
		return f()
	}

	return nil
}
func (s *MockCountryService) CountriesMsg(a0 uint64) pbutil.Buffer {
	fi := getMockFunc(s, s.CountriesMsg)
	if fi != nil {
		f, ok := fi.(func(uint64) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockCountryService.CountriesMsg()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockCountryService) Country(a0 uint64) *entity.Country {
	fi := getMockFunc(s, s.Country)
	if fi != nil {
		f, ok := fi.(func(uint64) *entity.Country)
		if !ok {
			panic("invalid mock func, MockCountryService.Country()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockCountryService) CountryDestroy(a0 uint64) bool {
	fi := getMockFunc(s, s.CountryDestroy)
	if fi != nil {
		f, ok := fi.(func(uint64) bool)
		if !ok {
			panic("invalid mock func, MockCountryService.CountryDestroy()")
		}
		return f(a0)
	}

	return false
}
func (s *MockCountryService) CountryDetailMsg(a0 uint64) pbutil.Buffer {
	fi := getMockFunc(s, s.CountryDetailMsg)
	if fi != nil {
		f, ok := fi.(func(uint64) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockCountryService.CountryDetailMsg()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockCountryService) CountryFlagHeroName(a0 int64) string {
	fi := getMockFunc(s, s.CountryFlagHeroName)
	if fi != nil {
		f, ok := fi.(func(int64) string)
		if !ok {
			panic("invalid mock func, MockCountryService.CountryFlagHeroName()")
		}
		return f(a0)
	}

	return ""
}
func (s *MockCountryService) CountryName(a0 uint64) string {
	fi := getMockFunc(s, s.CountryName)
	if fi != nil {
		f, ok := fi.(func(uint64) string)
		if !ok {
			panic("invalid mock func, MockCountryService.CountryName()")
		}
		return f(a0)
	}

	return ""
}
func (s *MockCountryService) CountryPrestigeMsg(a0 uint64) pbutil.Buffer {
	fi := getMockFunc(s, s.CountryPrestigeMsg)
	if fi != nil {
		f, ok := fi.(func(uint64) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockCountryService.CountryPrestigeMsg()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockCountryService) DetailMsgCacheDisable(a0 uint64) {
	fi := getMockFunc(s, s.DetailMsgCacheDisable)
	if fi != nil {
		f, ok := fi.(func(uint64))
		if !ok {
			panic("invalid mock func, MockCountryService.DetailMsgCacheDisable()")
		}
		f(a0)
	}

}
func (s *MockCountryService) ForceOfficialAppoint(a0 uint64, a1 int64, a2 shared_proto.CountryOfficialType) bool {
	fi := getMockFunc(s, s.ForceOfficialAppoint)
	if fi != nil {
		f, ok := fi.(func(uint64, int64, shared_proto.CountryOfficialType) bool)
		if !ok {
			panic("invalid mock func, MockCountryService.ForceOfficialAppoint()")
		}
		return f(a0, a1, a2)
	}

	return false
}
func (s *MockCountryService) ForceOfficialDepose(a0 uint64, a1 int64) bool {
	fi := getMockFunc(s, s.ForceOfficialDepose)
	if fi != nil {
		f, ok := fi.(func(uint64, int64) bool)
		if !ok {
			panic("invalid mock func, MockCountryService.ForceOfficialDepose()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockCountryService) GmAppointKing(a0 uint64, a1 int64) {
	fi := getMockFunc(s, s.GmAppointKing)
	if fi != nil {
		f, ok := fi.(func(uint64, int64))
		if !ok {
			panic("invalid mock func, MockCountryService.GmAppointKing()")
		}
		f(a0, a1)
	}

}
func (s *MockCountryService) GmAppointOfficial(a0 uint64, a1 int64, a2 shared_proto.CountryOfficialType) {
	fi := getMockFunc(s, s.GmAppointOfficial)
	if fi != nil {
		f, ok := fi.(func(uint64, int64, shared_proto.CountryOfficialType))
		if !ok {
			panic("invalid mock func, MockCountryService.GmAppointOfficial()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockCountryService) GmDeposeKing(a0 uint64) {
	fi := getMockFunc(s, s.GmDeposeKing)
	if fi != nil {
		f, ok := fi.(func(uint64))
		if !ok {
			panic("invalid mock func, MockCountryService.GmDeposeKing()")
		}
		f(a0)
	}

}
func (s *MockCountryService) GmDestroy(a0 uint64) {
	fi := getMockFunc(s, s.GmDestroy)
	if fi != nil {
		f, ok := fi.(func(uint64))
		if !ok {
			panic("invalid mock func, MockCountryService.GmDestroy()")
		}
		f(a0)
	}

}
func (s *MockCountryService) GmOfficialDeposeAll(a0 uint64) {
	fi := getMockFunc(s, s.GmOfficialDeposeAll)
	if fi != nil {
		f, ok := fi.(func(uint64))
		if !ok {
			panic("invalid mock func, MockCountryService.GmOfficialDeposeAll()")
		}
		f(a0)
	}

}
func (s *MockCountryService) HeroCountry(a0 int64) uint64 {
	fi := getMockFunc(s, s.HeroCountry)
	if fi != nil {
		f, ok := fi.(func(int64) uint64)
		if !ok {
			panic("invalid mock func, MockCountryService.HeroCountry()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockCountryService) HeroOfficial(a0 uint64, a1 int64) shared_proto.CountryOfficialType {
	fi := getMockFunc(s, s.HeroOfficial)
	if fi != nil {
		f, ok := fi.(func(uint64, int64) shared_proto.CountryOfficialType)
		if !ok {
			panic("invalid mock func, MockCountryService.HeroOfficial()")
		}
		return f(a0, a1)
	}

	return 0
}
func (s *MockCountryService) IsCountryDestroyed(a0 uint64) bool {
	fi := getMockFunc(s, s.IsCountryDestroyed)
	if fi != nil {
		f, ok := fi.(func(uint64) bool)
		if !ok {
			panic("invalid mock func, MockCountryService.IsCountryDestroyed()")
		}
		return f(a0)
	}

	return false
}
func (s *MockCountryService) IsOnChangeNameVote(a0 uint64) bool {
	fi := getMockFunc(s, s.IsOnChangeNameVote)
	if fi != nil {
		f, ok := fi.(func(uint64) bool)
		if !ok {
			panic("invalid mock func, MockCountryService.IsOnChangeNameVote()")
		}
		return f(a0)
	}

	return false
}
func (s *MockCountryService) King(a0 uint64) int64 {
	fi := getMockFunc(s, s.King)
	if fi != nil {
		f, ok := fi.(func(uint64) int64)
		if !ok {
			panic("invalid mock func, MockCountryService.King()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockCountryService) LockHeroCapital(a0 int64) uint64 {
	fi := getMockFunc(s, s.LockHeroCapital)
	if fi != nil {
		f, ok := fi.(func(int64) uint64)
		if !ok {
			panic("invalid mock func, MockCountryService.LockHeroCapital()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockCountryService) LockHeroCountry(a0 int64) uint64 {
	fi := getMockFunc(s, s.LockHeroCountry)
	if fi != nil {
		f, ok := fi.(func(int64) uint64)
		if !ok {
			panic("invalid mock func, MockCountryService.LockHeroCountry()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockCountryService) MsgCacheDisable(a0 uint64) {
	fi := getMockFunc(s, s.MsgCacheDisable)
	if fi != nil {
		f, ok := fi.(func(uint64))
		if !ok {
			panic("invalid mock func, MockCountryService.MsgCacheDisable()")
		}
		f(a0)
	}

}
func (s *MockCountryService) OfficialAppoint(a0 int64, a1 int64, a2 uint64, a3 shared_proto.CountryOfficialType, a4 int32) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.OfficialAppoint)
	if fi != nil {
		f, ok := fi.(func(int64, int64, uint64, shared_proto.CountryOfficialType, int32) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockCountryService.OfficialAppoint()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return nil, nil
}
func (s *MockCountryService) OfficialDepose(a0 int64, a1 int64, a2 uint64) (shared_proto.CountryOfficialType, pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.OfficialDepose)
	if fi != nil {
		f, ok := fi.(func(int64, int64, uint64) (shared_proto.CountryOfficialType, pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockCountryService.OfficialDepose()")
		}
		return f(a0, a1, a2)
	}

	return 0, nil, nil
}
func (s *MockCountryService) OfficialLeave(a0 int64, a1 uint64) (shared_proto.CountryOfficialType, pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.OfficialLeave)
	if fi != nil {
		f, ok := fi.(func(int64, uint64) (shared_proto.CountryOfficialType, pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockCountryService.OfficialLeave()")
		}
		return f(a0, a1)
	}

	return 0, nil, nil
}
func (s *MockCountryService) OnHeroOnline(a0 iface.HeroController, a1 uint64) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, uint64))
		if !ok {
			panic("invalid mock func, MockCountryService.OnHeroOnline()")
		}
		f(a0, a1)
	}

}
func (s *MockCountryService) ReducePrestige(a0 uint64, a1 uint64) (uint64, bool) {
	fi := getMockFunc(s, s.ReducePrestige)
	if fi != nil {
		f, ok := fi.(func(uint64, uint64) (uint64, bool))
		if !ok {
			panic("invalid mock func, MockCountryService.ReducePrestige()")
		}
		return f(a0, a1)
	}

	return 0, false
}
func (s *MockCountryService) TutorialCountriesProto() *shared_proto.CountriesProto {
	fi := getMockFunc(s, s.TutorialCountriesProto)
	if fi != nil {
		f, ok := fi.(func() *shared_proto.CountriesProto)
		if !ok {
			panic("invalid mock func, MockCountryService.TutorialCountriesProto()")
		}
		return f()
	}

	return nil
}
func (s *MockCountryService) UpdateMcWarMsg(a0 *shared_proto.McWarProto) {
	fi := getMockFunc(s, s.UpdateMcWarMsg)
	if fi != nil {
		f, ok := fi.(func(*shared_proto.McWarProto))
		if !ok {
			panic("invalid mock func, MockCountryService.UpdateMcWarMsg()")
		}
		f(a0)
	}

}
func (s *MockCountryService) UpdateMingcsMsg(a0 *shared_proto.MingcsProto) {
	fi := getMockFunc(s, s.UpdateMingcsMsg)
	if fi != nil {
		f, ok := fi.(func(*shared_proto.MingcsProto))
		if !ok {
			panic("invalid mock func, MockCountryService.UpdateMingcsMsg()")
		}
		f(a0)
	}

}
func (s *MockCountryService) WalkCountryOnlineHero(a0 uint64, a1 iface.CountryHeroWalker) {
	fi := getMockFunc(s, s.WalkCountryOnlineHero)
	if fi != nil {
		f, ok := fi.(func(uint64, iface.CountryHeroWalker))
		if !ok {
			panic("invalid mock func, MockCountryService.WalkCountryOnlineHero()")
		}
		f(a0, a1)
	}

}

var DbService = &MockDbService{}

type MockDbService struct{}

func (s *MockDbService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockDbService) AddChatMsg(a0 context.Context, a1 int64, a2 []byte, a3 *shared_proto.ChatMsgProto) (int64, error) {
	fi := getMockFunc(s, s.AddChatMsg)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []byte, *shared_proto.ChatMsgProto) (int64, error))
		if !ok {
			panic("invalid mock func, MockDbService.AddChatMsg()")
		}
		return f(a0, a1, a2, a3)
	}

	return 0, nil
}
func (s *MockDbService) AddFarmSteal(a0 context.Context, a1 int64, a2 int64, a3 cb.Cube) error {
	fi := getMockFunc(s, s.AddFarmSteal)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int64, cb.Cube) error)
		if !ok {
			panic("invalid mock func, MockDbService.AddFarmSteal()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockDbService) AddMcWarGuildRecord(a0 context.Context, a1 uint64, a2 uint64, a3 int64, a4 *shared_proto.McWarTroopsInfoProto) error {
	fi := getMockFunc(s, s.AddMcWarGuildRecord)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, uint64, int64, *shared_proto.McWarTroopsInfoProto) error)
		if !ok {
			panic("invalid mock func, MockDbService.AddMcWarGuildRecord()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return nil
}
func (s *MockDbService) AddMcWarHeroRecord(a0 context.Context, a1 uint64, a2 uint64, a3 int64, a4 *shared_proto.McWarTroopAllRecordProto) error {
	fi := getMockFunc(s, s.AddMcWarHeroRecord)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, uint64, int64, *shared_proto.McWarTroopAllRecordProto) error)
		if !ok {
			panic("invalid mock func, MockDbService.AddMcWarHeroRecord()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return nil
}
func (s *MockDbService) AddMcWarRecord(a0 context.Context, a1 uint64, a2 uint64, a3 *shared_proto.McWarFightRecordProto) error {
	fi := getMockFunc(s, s.AddMcWarRecord)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, uint64, *shared_proto.McWarFightRecordProto) error)
		if !ok {
			panic("invalid mock func, MockDbService.AddMcWarRecord()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockDbService) CallingTimes() uint64 {
	fi := getMockFunc(s, s.CallingTimes)
	if fi != nil {
		f, ok := fi.(func() uint64)
		if !ok {
			panic("invalid mock func, MockDbService.CallingTimes()")
		}
		return f()
	}

	return 0
}
func (s *MockDbService) Close() error {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func() error)
		if !ok {
			panic("invalid mock func, MockDbService.Close()")
		}
		return f()
	}

	return nil
}
func (s *MockDbService) CreateFarmCube(a0 context.Context, a1 *entity.FarmCube) error {
	fi := getMockFunc(s, s.CreateFarmCube)
	if fi != nil {
		f, ok := fi.(func(context.Context, *entity.FarmCube) error)
		if !ok {
			panic("invalid mock func, MockDbService.CreateFarmCube()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockDbService) CreateFarmLog(a0 context.Context, a1 *shared_proto.FarmStealLogProto) error {
	fi := getMockFunc(s, s.CreateFarmLog)
	if fi != nil {
		f, ok := fi.(func(context.Context, *shared_proto.FarmStealLogProto) error)
		if !ok {
			panic("invalid mock func, MockDbService.CreateFarmLog()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockDbService) CreateGuild(a0 context.Context, a1 int64, a2 []byte) error {
	fi := getMockFunc(s, s.CreateGuild)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []byte) error)
		if !ok {
			panic("invalid mock func, MockDbService.CreateGuild()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) CreateHero(a0 context.Context, a1 *entity.Hero) (bool, error) {
	fi := getMockFunc(s, s.CreateHero)
	if fi != nil {
		f, ok := fi.(func(context.Context, *entity.Hero) (bool, error))
		if !ok {
			panic("invalid mock func, MockDbService.CreateHero()")
		}
		return f(a0, a1)
	}

	return false, nil
}
func (s *MockDbService) CreateMail(a0 context.Context, a1 uint64, a2 int64, a3 []byte, a4 bool, a5 bool, a6 bool, a7 int32, a8 int64) error {
	fi := getMockFunc(s, s.CreateMail)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, int64, []byte, bool, bool, bool, int32, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.CreateMail()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8)
	}

	return nil
}
func (s *MockDbService) CreateOrder(a0 context.Context, a1 string, a2 uint64, a3 int64, a4 uint32, a5 uint32, a6 int64, a7 uint64, a8 int64) error {
	fi := getMockFunc(s, s.CreateOrder)
	if fi != nil {
		f, ok := fi.(func(context.Context, string, uint64, int64, uint32, uint32, int64, uint64, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.CreateOrder()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8)
	}

	return nil
}
func (s *MockDbService) DelMcWarHeroRecord(a0 context.Context, a1 int32) error {
	fi := getMockFunc(s, s.DelMcWarHeroRecord)
	if fi != nil {
		f, ok := fi.(func(context.Context, int32) error)
		if !ok {
			panic("invalid mock func, MockDbService.DelMcWarHeroRecord()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockDbService) DelMcWarHeroRecordWithHeroId(a0 context.Context, a1 int32, a2 uint64, a3 int64) error {
	fi := getMockFunc(s, s.DelMcWarHeroRecordWithHeroId)
	if fi != nil {
		f, ok := fi.(func(context.Context, int32, uint64, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.DelMcWarHeroRecordWithHeroId()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockDbService) DeleteChatWindow(a0 context.Context, a1 int64, a2 []byte) error {
	fi := getMockFunc(s, s.DeleteChatWindow)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []byte) error)
		if !ok {
			panic("invalid mock func, MockDbService.DeleteChatWindow()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) DeleteGuild(a0 context.Context, a1 int64) error {
	fi := getMockFunc(s, s.DeleteGuild)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.DeleteGuild()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockDbService) DeleteMail(a0 context.Context, a1 uint64, a2 int64) error {
	fi := getMockFunc(s, s.DeleteMail)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.DeleteMail()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) DeleteMultiMail(a0 context.Context, a1 int64, a2 []uint64, a3 bool) error {
	fi := getMockFunc(s, s.DeleteMultiMail)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []uint64, bool) error)
		if !ok {
			panic("invalid mock func, MockDbService.DeleteMultiMail()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockDbService) FindSettingsOpen(a0 context.Context, a1 shared_proto.SettingType, a2 []int64) ([]int64, error) {
	fi := getMockFunc(s, s.FindSettingsOpen)
	if fi != nil {
		f, ok := fi.(func(context.Context, shared_proto.SettingType, []int64) ([]int64, error))
		if !ok {
			panic("invalid mock func, MockDbService.FindSettingsOpen()")
		}
		return f(a0, a1, a2)
	}

	return nil, nil
}
func (s *MockDbService) GMFarmRipe(a0 context.Context, a1 int64, a2 int64, a3 int64) error {
	fi := getMockFunc(s, s.GMFarmRipe)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int64, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.GMFarmRipe()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockDbService) HeroId(a0 context.Context, a1 string) (int64, error) {
	fi := getMockFunc(s, s.HeroId)
	if fi != nil {
		f, ok := fi.(func(context.Context, string) (int64, error))
		if !ok {
			panic("invalid mock func, MockDbService.HeroId()")
		}
		return f(a0, a1)
	}

	return 0, nil
}
func (s *MockDbService) HeroIdExist(a0 context.Context, a1 int64) (bool, error) {
	fi := getMockFunc(s, s.HeroIdExist)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) (bool, error))
		if !ok {
			panic("invalid mock func, MockDbService.HeroIdExist()")
		}
		return f(a0, a1)
	}

	return false, nil
}
func (s *MockDbService) HeroIds(a0 context.Context) ([]int64, error) {
	fi := getMockFunc(s, s.HeroIds)
	if fi != nil {
		f, ok := fi.(func(context.Context) ([]int64, error))
		if !ok {
			panic("invalid mock func, MockDbService.HeroIds()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockDbService) HeroNameExist(a0 context.Context, a1 string) (bool, error) {
	fi := getMockFunc(s, s.HeroNameExist)
	if fi != nil {
		f, ok := fi.(func(context.Context, string) (bool, error))
		if !ok {
			panic("invalid mock func, MockDbService.HeroNameExist()")
		}
		return f(a0, a1)
	}

	return false, nil
}
func (s *MockDbService) InsertBaiZhanReplay(a0 context.Context, a1 int64, a2 int64, a3 *shared_proto.BaiZhanReplayProto, a4 bool, a5 int64) error {
	fi := getMockFunc(s, s.InsertBaiZhanReplay)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int64, *shared_proto.BaiZhanReplayProto, bool, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.InsertBaiZhanReplay()")
		}
		return f(a0, a1, a2, a3, a4, a5)
	}

	return nil
}
func (s *MockDbService) InsertGuildLog(a0 context.Context, a1 int64, a2 *shared_proto.GuildLogProto) error {
	fi := getMockFunc(s, s.InsertGuildLog)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, *shared_proto.GuildLogProto) error)
		if !ok {
			panic("invalid mock func, MockDbService.InsertGuildLog()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) InsertXuanyRecord(a0 context.Context, a1 int64, a2 []byte) (int64, error) {
	fi := getMockFunc(s, s.InsertXuanyRecord)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []byte) (int64, error))
		if !ok {
			panic("invalid mock func, MockDbService.InsertXuanyRecord()")
		}
		return f(a0, a1, a2)
	}

	return 0, nil
}
func (s *MockDbService) IsCollectableMail(a0 context.Context, a1 uint64) (bool, error) {
	fi := getMockFunc(s, s.IsCollectableMail)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64) (bool, error))
		if !ok {
			panic("invalid mock func, MockDbService.IsCollectableMail()")
		}
		return f(a0, a1)
	}

	return false, nil
}
func (s *MockDbService) ListHeroChatMsg(a0 context.Context, a1 []byte, a2 uint64) ([]*shared_proto.ChatMsgProto, error) {
	fi := getMockFunc(s, s.ListHeroChatMsg)
	if fi != nil {
		f, ok := fi.(func(context.Context, []byte, uint64) ([]*shared_proto.ChatMsgProto, error))
		if !ok {
			panic("invalid mock func, MockDbService.ListHeroChatMsg()")
		}
		return f(a0, a1, a2)
	}

	return nil, nil
}
func (s *MockDbService) ListHeroChatWindow(a0 context.Context, a1 int64) ([]uint64, isql.BytesArray, error) {
	fi := getMockFunc(s, s.ListHeroChatWindow)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) ([]uint64, isql.BytesArray, error))
		if !ok {
			panic("invalid mock func, MockDbService.ListHeroChatWindow()")
		}
		return f(a0, a1)
	}

	return nil, nil, nil
}
func (s *MockDbService) LoadAllGuild(a0 context.Context) ([]*sharedguilddata.Guild, error) {
	fi := getMockFunc(s, s.LoadAllGuild)
	if fi != nil {
		f, ok := fi.(func(context.Context) ([]*sharedguilddata.Guild, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadAllGuild()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockDbService) LoadAllHeroData(a0 context.Context) ([]*entity.Hero, error) {
	fi := getMockFunc(s, s.LoadAllHeroData)
	if fi != nil {
		f, ok := fi.(func(context.Context) ([]*entity.Hero, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadAllHeroData()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockDbService) LoadAllRegionHero(a0 context.Context) ([]*entity.Hero, error) {
	fi := getMockFunc(s, s.LoadAllRegionHero)
	if fi != nil {
		f, ok := fi.(func(context.Context) ([]*entity.Hero, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadAllRegionHero()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockDbService) LoadBaiZhanRecord(a0 context.Context, a1 int64, a2 uint64) (isql.BytesArray, error) {
	fi := getMockFunc(s, s.LoadBaiZhanRecord)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, uint64) (isql.BytesArray, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadBaiZhanRecord()")
		}
		return f(a0, a1, a2)
	}

	return nil, nil
}
func (s *MockDbService) LoadCanStealCount(a0 context.Context, a1 int64, a2 int64, a3 int64, a4 uint64) (uint64, error) {
	fi := getMockFunc(s, s.LoadCanStealCount)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int64, int64, uint64) (uint64, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadCanStealCount()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return 0, nil
}
func (s *MockDbService) LoadCanStealCube(a0 context.Context, a1 int64, a2 int64, a3 int64, a4 uint64) ([]*entity.FarmCube, error) {
	fi := getMockFunc(s, s.LoadCanStealCube)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int64, int64, uint64) ([]*entity.FarmCube, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadCanStealCube()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return nil, nil
}
func (s *MockDbService) LoadChatMsg(a0 context.Context, a1 int64) (*shared_proto.ChatMsgProto, error) {
	fi := getMockFunc(s, s.LoadChatMsg)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) (*shared_proto.ChatMsgProto, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadChatMsg()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockDbService) LoadCollectMailPrize(a0 context.Context, a1 uint64, a2 int64) (*resdata.Prize, error) {
	fi := getMockFunc(s, s.LoadCollectMailPrize)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, int64) (*resdata.Prize, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadCollectMailPrize()")
		}
		return f(a0, a1, a2)
	}

	return nil, nil
}
func (s *MockDbService) LoadFarmCube(a0 context.Context, a1 int64, a2 cb.Cube) (*entity.FarmCube, error) {
	fi := getMockFunc(s, s.LoadFarmCube)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, cb.Cube) (*entity.FarmCube, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadFarmCube()")
		}
		return f(a0, a1, a2)
	}

	return nil, nil
}
func (s *MockDbService) LoadFarmCubes(a0 context.Context, a1 int64) ([]*entity.FarmCube, error) {
	fi := getMockFunc(s, s.LoadFarmCubes)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) ([]*entity.FarmCube, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadFarmCubes()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockDbService) LoadFarmHarvestCubes(a0 context.Context, a1 int64) ([]*entity.FarmCube, error) {
	fi := getMockFunc(s, s.LoadFarmHarvestCubes)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) ([]*entity.FarmCube, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadFarmHarvestCubes()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockDbService) LoadFarmLog(a0 context.Context, a1 int64, a2 uint64) ([]*shared_proto.FarmStealLogProto, error) {
	fi := getMockFunc(s, s.LoadFarmLog)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, uint64) ([]*shared_proto.FarmStealLogProto, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadFarmLog()")
		}
		return f(a0, a1, a2)
	}

	return nil, nil
}
func (s *MockDbService) LoadFarmStealCount(a0 context.Context, a1 int64, a2 int64, a3 cb.Cube) (uint64, error) {
	fi := getMockFunc(s, s.LoadFarmStealCount)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int64, cb.Cube) (uint64, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadFarmStealCount()")
		}
		return f(a0, a1, a2, a3)
	}

	return 0, nil
}
func (s *MockDbService) LoadFarmStealCubes(a0 context.Context, a1 int64, a2 int64, a3 uint64) ([]*entity.FarmCube, error) {
	fi := getMockFunc(s, s.LoadFarmStealCubes)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int64, uint64) ([]*entity.FarmCube, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadFarmStealCubes()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil, nil
}
func (s *MockDbService) LoadGuild(a0 context.Context, a1 int64) (*sharedguilddata.Guild, error) {
	fi := getMockFunc(s, s.LoadGuild)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) (*sharedguilddata.Guild, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadGuild()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockDbService) LoadGuildLogs(a0 context.Context, a1 int64, a2 shared_proto.GuildLogType, a3 int64, a4 uint64) ([]*shared_proto.GuildLogProto, error) {
	fi := getMockFunc(s, s.LoadGuildLogs)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, shared_proto.GuildLogType, int64, uint64) ([]*shared_proto.GuildLogProto, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadGuildLogs()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return nil, nil
}
func (s *MockDbService) LoadHero(a0 context.Context, a1 int64) (*entity.Hero, error) {
	fi := getMockFunc(s, s.LoadHero)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) (*entity.Hero, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadHero()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockDbService) LoadHeroCount(a0 context.Context) (uint64, error) {
	fi := getMockFunc(s, s.LoadHeroCount)
	if fi != nil {
		f, ok := fi.(func(context.Context) (uint64, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadHeroCount()")
		}
		return f(a0)
	}

	return 0, nil
}
func (s *MockDbService) LoadHeroListByCountry(a0 context.Context, a1 uint64, a2 uint64, a3 uint64) ([]*entity.Hero, error) {
	fi := getMockFunc(s, s.LoadHeroListByCountry)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, uint64, uint64) ([]*entity.Hero, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadHeroListByCountry()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil, nil
}
func (s *MockDbService) LoadHeroListByNameAndCountry(a0 context.Context, a1 string, a2 uint64, a3 uint64, a4 uint64) ([]*entity.Hero, error) {
	fi := getMockFunc(s, s.LoadHeroListByNameAndCountry)
	if fi != nil {
		f, ok := fi.(func(context.Context, string, uint64, uint64, uint64) ([]*entity.Hero, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadHeroListByNameAndCountry()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return nil, nil
}
func (s *MockDbService) LoadHeroMailList(a0 context.Context, a1 int64, a2 uint64, a3 int32, a4 int32, a5 int32, a6 int32, a7 int32, a8 int32, a9 uint64) ([]*shared_proto.MailProto, error) {
	fi := getMockFunc(s, s.LoadHeroMailList)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, uint64, int32, int32, int32, int32, int32, int32, uint64) ([]*shared_proto.MailProto, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadHeroMailList()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	}

	return nil, nil
}
func (s *MockDbService) LoadHerosByName(a0 context.Context, a1 string, a2 uint64, a3 uint64) ([]*entity.Hero, error) {
	fi := getMockFunc(s, s.LoadHerosByName)
	if fi != nil {
		f, ok := fi.(func(context.Context, string, uint64, uint64) ([]*entity.Hero, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadHerosByName()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil, nil
}
func (s *MockDbService) LoadJoinedMcWarId(a0 context.Context, a1 int64) (*entity.JoinedMcWarIds, error) {
	fi := getMockFunc(s, s.LoadJoinedMcWarId)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) (*entity.JoinedMcWarIds, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadJoinedMcWarId()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockDbService) LoadKey(a0 context.Context, a1 server_proto.Key) ([]byte, error) {
	fi := getMockFunc(s, s.LoadKey)
	if fi != nil {
		f, ok := fi.(func(context.Context, server_proto.Key) ([]byte, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadKey()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockDbService) LoadMail(a0 context.Context, a1 uint64) ([]byte, error) {
	fi := getMockFunc(s, s.LoadMail)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64) ([]byte, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadMail()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockDbService) LoadMailCountHasPrizeNotCollected(a0 context.Context, a1 int64) (int, error) {
	fi := getMockFunc(s, s.LoadMailCountHasPrizeNotCollected)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) (int, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadMailCountHasPrizeNotCollected()")
		}
		return f(a0, a1)
	}

	return 0, nil
}
func (s *MockDbService) LoadMailCountHasReportNotReaded(a0 context.Context, a1 int64, a2 int32) (int, error) {
	fi := getMockFunc(s, s.LoadMailCountHasReportNotReaded)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int32) (int, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadMailCountHasReportNotReaded()")
		}
		return f(a0, a1, a2)
	}

	return 0, nil
}
func (s *MockDbService) LoadMailCountNoReportNotReaded(a0 context.Context, a1 int64) (int, error) {
	fi := getMockFunc(s, s.LoadMailCountNoReportNotReaded)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) (int, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadMailCountNoReportNotReaded()")
		}
		return f(a0, a1)
	}

	return 0, nil
}
func (s *MockDbService) LoadMcWarGuildRecord(a0 context.Context, a1 uint64, a2 uint64, a3 int64) (*shared_proto.McWarTroopsInfoProto, error) {
	fi := getMockFunc(s, s.LoadMcWarGuildRecord)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, uint64, int64) (*shared_proto.McWarTroopsInfoProto, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadMcWarGuildRecord()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil, nil
}
func (s *MockDbService) LoadMcWarHeroRecord(a0 context.Context, a1 uint64, a2 uint64, a3 int64) (*shared_proto.McWarTroopAllRecordProto, error) {
	fi := getMockFunc(s, s.LoadMcWarHeroRecord)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, uint64, int64) (*shared_proto.McWarTroopAllRecordProto, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadMcWarHeroRecord()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil, nil
}
func (s *MockDbService) LoadMcWarRecord(a0 context.Context, a1 uint64, a2 uint64) (*shared_proto.McWarFightRecordProto, error) {
	fi := getMockFunc(s, s.LoadMcWarRecord)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, uint64) (*shared_proto.McWarFightRecordProto, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadMcWarRecord()")
		}
		return f(a0, a1, a2)
	}

	return nil, nil
}
func (s *MockDbService) LoadNoGuildHeroListByName(a0 context.Context, a1 string, a2 uint64, a3 uint64) ([]*entity.Hero, error) {
	fi := getMockFunc(s, s.LoadNoGuildHeroListByName)
	if fi != nil {
		f, ok := fi.(func(context.Context, string, uint64, uint64) ([]*entity.Hero, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadNoGuildHeroListByName()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil, nil
}
func (s *MockDbService) LoadRecommendHeros(a0 context.Context, a1 bool, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 int64) ([]*entity.Hero, error) {
	fi := getMockFunc(s, s.LoadRecommendHeros)
	if fi != nil {
		f, ok := fi.(func(context.Context, bool, uint64, uint64, uint64, uint64, int64) ([]*entity.Hero, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadRecommendHeros()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6)
	}

	return nil, nil
}
func (s *MockDbService) LoadUnreadChatCount(a0 context.Context, a1 int64) (uint64, error) {
	fi := getMockFunc(s, s.LoadUnreadChatCount)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) (uint64, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadUnreadChatCount()")
		}
		return f(a0, a1)
	}

	return 0, nil
}
func (s *MockDbService) LoadUserMisc(a0 context.Context, a1 int64) (*server_proto.UserMiscProto, error) {
	fi := getMockFunc(s, s.LoadUserMisc)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) (*server_proto.UserMiscProto, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadUserMisc()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockDbService) LoadXuanyRecord(a0 context.Context, a1 int64, a2 int64, a3 bool) ([]int64, isql.BytesArray, error) {
	fi := getMockFunc(s, s.LoadXuanyRecord)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int64, bool) ([]int64, isql.BytesArray, error))
		if !ok {
			panic("invalid mock func, MockDbService.LoadXuanyRecord()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil, nil, nil
}
func (s *MockDbService) MaxGuildId(a0 context.Context) (int64, error) {
	fi := getMockFunc(s, s.MaxGuildId)
	if fi != nil {
		f, ok := fi.(func(context.Context) (int64, error))
		if !ok {
			panic("invalid mock func, MockDbService.MaxGuildId()")
		}
		return f(a0)
	}

	return 0, nil
}
func (s *MockDbService) MaxMailId(a0 context.Context) (uint64, error) {
	fi := getMockFunc(s, s.MaxMailId)
	if fi != nil {
		f, ok := fi.(func(context.Context) (uint64, error))
		if !ok {
			panic("invalid mock func, MockDbService.MaxMailId()")
		}
		return f(a0)
	}

	return 0, nil
}
func (s *MockDbService) OrderExist(a0 context.Context, a1 string) (bool, error) {
	fi := getMockFunc(s, s.OrderExist)
	if fi != nil {
		f, ok := fi.(func(context.Context, string) (bool, error))
		if !ok {
			panic("invalid mock func, MockDbService.OrderExist()")
		}
		return f(a0, a1)
	}

	return false, nil
}
func (s *MockDbService) PlantFarmCube(a0 context.Context, a1 int64, a2 cb.Cube, a3 int64, a4 int64, a5 uint64) error {
	fi := getMockFunc(s, s.PlantFarmCube)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, cb.Cube, int64, int64, uint64) error)
		if !ok {
			panic("invalid mock func, MockDbService.PlantFarmCube()")
		}
		return f(a0, a1, a2, a3, a4, a5)
	}

	return nil
}
func (s *MockDbService) ReadChat(a0 context.Context, a1 int64, a2 []byte) error {
	fi := getMockFunc(s, s.ReadChat)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []byte) error)
		if !ok {
			panic("invalid mock func, MockDbService.ReadChat()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) ReadMultiMail(a0 context.Context, a1 int64, a2 []uint64, a3 bool) (*resdata.Prize, error) {
	fi := getMockFunc(s, s.ReadMultiMail)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []uint64, bool) (*resdata.Prize, error))
		if !ok {
			panic("invalid mock func, MockDbService.ReadMultiMail()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil, nil
}
func (s *MockDbService) RemoveChatMsg(a0 context.Context, a1 int64) error {
	fi := getMockFunc(s, s.RemoveChatMsg)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.RemoveChatMsg()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockDbService) RemoveFarmCube(a0 context.Context, a1 int64, a2 cb.Cube) error {
	fi := getMockFunc(s, s.RemoveFarmCube)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, cb.Cube) error)
		if !ok {
			panic("invalid mock func, MockDbService.RemoveFarmCube()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) RemoveFarmLog(a0 context.Context, a1 int32) error {
	fi := getMockFunc(s, s.RemoveFarmLog)
	if fi != nil {
		f, ok := fi.(func(context.Context, int32) error)
		if !ok {
			panic("invalid mock func, MockDbService.RemoveFarmLog()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockDbService) RemoveFarmSteal(a0 context.Context, a1 int64, a2 []cb.Cube) error {
	fi := getMockFunc(s, s.RemoveFarmSteal)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []cb.Cube) error)
		if !ok {
			panic("invalid mock func, MockDbService.RemoveFarmSteal()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) ResetConflictFarmCubes(a0 context.Context, a1 int64) error {
	fi := getMockFunc(s, s.ResetConflictFarmCubes)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.ResetConflictFarmCubes()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockDbService) ResetFarmCubes(a0 context.Context, a1 int64) error {
	fi := getMockFunc(s, s.ResetFarmCubes)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.ResetFarmCubes()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockDbService) SaveFarmCube(a0 context.Context, a1 *entity.FarmCube) error {
	fi := getMockFunc(s, s.SaveFarmCube)
	if fi != nil {
		f, ok := fi.(func(context.Context, *entity.FarmCube) error)
		if !ok {
			panic("invalid mock func, MockDbService.SaveFarmCube()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockDbService) SaveGuild(a0 context.Context, a1 int64, a2 []byte) error {
	fi := getMockFunc(s, s.SaveGuild)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []byte) error)
		if !ok {
			panic("invalid mock func, MockDbService.SaveGuild()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) SaveHero(a0 context.Context, a1 *entity.Hero) error {
	fi := getMockFunc(s, s.SaveHero)
	if fi != nil {
		f, ok := fi.(func(context.Context, *entity.Hero) error)
		if !ok {
			panic("invalid mock func, MockDbService.SaveHero()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockDbService) SaveKey(a0 context.Context, a1 server_proto.Key, a2 []byte) error {
	fi := getMockFunc(s, s.SaveKey)
	if fi != nil {
		f, ok := fi.(func(context.Context, server_proto.Key, []byte) error)
		if !ok {
			panic("invalid mock func, MockDbService.SaveKey()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) SetFarmRipeTime(a0 context.Context, a1 int64, a2 int64) error {
	fi := getMockFunc(s, s.SetFarmRipeTime)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.SetFarmRipeTime()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) UpdateChatMsg(a0 context.Context, a1 int64, a2 *shared_proto.ChatMsgProto) bool {
	fi := getMockFunc(s, s.UpdateChatMsg)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, *shared_proto.ChatMsgProto) bool)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateChatMsg()")
		}
		return f(a0, a1, a2)
	}

	return false
}
func (s *MockDbService) UpdateChatWindow(a0 context.Context, a1 int64, a2 []byte, a3 []byte, a4 bool, a5 int32, a6 bool) error {
	fi := getMockFunc(s, s.UpdateChatWindow)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []byte, []byte, bool, int32, bool) error)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateChatWindow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6)
	}

	return nil
}
func (s *MockDbService) UpdateFarmCubeRipeTime(a0 context.Context, a1 int64, a2 cb.Cube, a3 int64, a4 int64) error {
	fi := getMockFunc(s, s.UpdateFarmCubeRipeTime)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, cb.Cube, int64, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateFarmCubeRipeTime()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return nil
}
func (s *MockDbService) UpdateFarmCubeState(a0 context.Context, a1 int64, a2 cb.Cube, a3 int64, a4 int64) error {
	fi := getMockFunc(s, s.UpdateFarmCubeState)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, cb.Cube, int64, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateFarmCubeState()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return nil
}
func (s *MockDbService) UpdateFarmStealTimes(a0 context.Context, a1 int64, a2 []cb.Cube) error {
	fi := getMockFunc(s, s.UpdateFarmStealTimes)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, []cb.Cube) error)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateFarmStealTimes()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) UpdateHeroGuildId(a0 context.Context, a1 int64, a2 int64) error {
	fi := getMockFunc(s, s.UpdateHeroGuildId)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, int64) error)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateHeroGuildId()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) UpdateHeroName(a0 context.Context, a1 int64, a2 string, a3 string) bool {
	fi := getMockFunc(s, s.UpdateHeroName)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, string, string) bool)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateHeroName()")
		}
		return f(a0, a1, a2, a3)
	}

	return false
}
func (s *MockDbService) UpdateHeroOfflineBoolIfExpected(a0 context.Context, a1 int64, a2 isql.OfflineBool, a3 bool, a4 bool) (bool, error) {
	fi := getMockFunc(s, s.UpdateHeroOfflineBoolIfExpected)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, isql.OfflineBool, bool, bool) (bool, error))
		if !ok {
			panic("invalid mock func, MockDbService.UpdateHeroOfflineBoolIfExpected()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return false, nil
}
func (s *MockDbService) UpdateMailCollected(a0 context.Context, a1 uint64, a2 int64, a3 bool) error {
	fi := getMockFunc(s, s.UpdateMailCollected)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, int64, bool) error)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateMailCollected()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockDbService) UpdateMailKeep(a0 context.Context, a1 uint64, a2 int64, a3 bool) error {
	fi := getMockFunc(s, s.UpdateMailKeep)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, int64, bool) error)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateMailKeep()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockDbService) UpdateMailRead(a0 context.Context, a1 uint64, a2 int64, a3 bool) error {
	fi := getMockFunc(s, s.UpdateMailRead)
	if fi != nil {
		f, ok := fi.(func(context.Context, uint64, int64, bool) error)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateMailRead()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockDbService) UpdateSettings(a0 context.Context, a1 int64, a2 uint64) error {
	fi := getMockFunc(s, s.UpdateSettings)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, uint64) error)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateSettings()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockDbService) UpdateUserMisc(a0 context.Context, a1 int64, a2 *server_proto.UserMiscProto) error {
	fi := getMockFunc(s, s.UpdateUserMisc)
	if fi != nil {
		f, ok := fi.(func(context.Context, int64, *server_proto.UserMiscProto) error)
		if !ok {
			panic("invalid mock func, MockDbService.UpdateUserMisc()")
		}
		return f(a0, a1, a2)
	}

	return nil
}

var DepotModule = &MockDepotModule{}

type MockDepotModule struct{}

func (s *MockDepotModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var DianquanModule = &MockDianquanModule{}

type MockDianquanModule struct{}

func (s *MockDianquanModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var DomesticModule = &MockDomesticModule{}

type MockDomesticModule struct{}

func (s *MockDomesticModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockDomesticModule) UseBuffGoods(a0 *buffer.BufferData, a1 uint64, a2 iface.HeroController) {
	fi := getMockFunc(s, s.UseBuffGoods)
	if fi != nil {
		f, ok := fi.(func(*buffer.BufferData, uint64, iface.HeroController))
		if !ok {
			panic("invalid mock func, MockDomesticModule.UseBuffGoods()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockDomesticModule) UseMianGoods(a0 *buffer.BufferData, a1 uint64, a2 iface.HeroController) {
	fi := getMockFunc(s, s.UseMianGoods)
	if fi != nil {
		f, ok := fi.(func(*buffer.BufferData, uint64, iface.HeroController))
		if !ok {
			panic("invalid mock func, MockDomesticModule.UseMianGoods()")
		}
		f(a0, a1, a2)
	}

}

var DungeonModule = &MockDungeonModule{}

type MockDungeonModule struct{}

func (s *MockDungeonModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var EquipmentModule = &MockEquipmentModule{}

type MockEquipmentModule struct{}

func (s *MockEquipmentModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockEquipmentModule) OnAddGoodsEvent(a0 *entity.Hero, a1 herolock.LockResult, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.OnAddGoodsEvent)
	if fi != nil {
		f, ok := fi.(func(*entity.Hero, herolock.LockResult, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockEquipmentModule.OnAddGoodsEvent()")
		}
		f(a0, a1, a2, a3)
	}

}

var ExtraTimesService = &MockExtraTimesService{}

type MockExtraTimesService struct{}

func (s *MockExtraTimesService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockExtraTimesService) MultiLevelNpcMaxTimes() extratimesface.ExtraMaxTimes {
	fi := getMockFunc(s, s.MultiLevelNpcMaxTimes)
	if fi != nil {
		f, ok := fi.(func() extratimesface.ExtraMaxTimes)
		if !ok {
			panic("invalid mock func, MockExtraTimesService.MultiLevelNpcMaxTimes()")
		}
		return f()
	}

	return nil
}

var FarmModule = &MockFarmModule{}

type MockFarmModule struct{}

func (s *MockFarmModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var FarmService = &MockFarmService{}

type MockFarmService struct{}

func (s *MockFarmService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockFarmService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockFarmService.Close()")
		}
		f()
	}

}
func (s *MockFarmService) FuncNoWait(a0 entity.FarmFuncType) {
	fi := getMockFunc(s, s.FuncNoWait)
	if fi != nil {
		f, ok := fi.(func(entity.FarmFuncType))
		if !ok {
			panic("invalid mock func, MockFarmService.FuncNoWait()")
		}
		f(a0)
	}

}
func (s *MockFarmService) FuncWait(a0 string, a1 pbutil.Buffer, a2 iface.HeroController, a3 entity.FarmFuncType) {
	fi := getMockFunc(s, s.FuncWait)
	if fi != nil {
		f, ok := fi.(func(string, pbutil.Buffer, iface.HeroController, entity.FarmFuncType))
		if !ok {
			panic("invalid mock func, MockFarmService.FuncWait()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockFarmService) GMCanSteal(a0 int64) {
	fi := getMockFunc(s, s.GMCanSteal)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockFarmService.GMCanSteal()")
		}
		f(a0)
	}

}
func (s *MockFarmService) GMRipe(a0 int64) {
	fi := getMockFunc(s, s.GMRipe)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockFarmService.GMRipe()")
		}
		f(a0)
	}

}
func (s *MockFarmService) ReduceRipeTime(a0 int64, a1 time.Duration) {
	fi := getMockFunc(s, s.ReduceRipeTime)
	if fi != nil {
		f, ok := fi.(func(int64, time.Duration))
		if !ok {
			panic("invalid mock func, MockFarmService.ReduceRipeTime()")
		}
		f(a0, a1)
	}

}
func (s *MockFarmService) ReduceRipeTimePercent(a0 int64, a1 *entity.BuffInfo, a2 *entity.BuffInfo) {
	fi := getMockFunc(s, s.ReduceRipeTimePercent)
	if fi != nil {
		f, ok := fi.(func(int64, *entity.BuffInfo, *entity.BuffInfo))
		if !ok {
			panic("invalid mock func, MockFarmService.ReduceRipeTimePercent()")
		}
		f(a0, a1, a2)
	}

}

// 更新农场地块
func (s *MockFarmService) UpdateFarmCubeWithOffset(a0 int64, a1 uint64, a2 int, a3 int, a4 []cb.Cube, a5 bool, a6 time.Time) {
	fi := getMockFunc(s, s.UpdateFarmCubeWithOffset)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, int, int, []cb.Cube, bool, time.Time))
		if !ok {
			panic("invalid mock func, MockFarmService.UpdateFarmCubeWithOffset()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}

/*
 更新农场地块
 absCubes 所有本次变化的地块
 allConflictedBlocks 本次变化的地块中，有冲突的地块
 npcConflictOffsets 所有 npc 造成的冲突地块
*/
func (s *MockFarmService) UpdateFarmCubes(a0 int64, a1 uint64, a2 []cb.Cube, a3 []cb.Cube, a4 []cb.Cube, a5 int, a6 int, a7 time.Time) {
	fi := getMockFunc(s, s.UpdateFarmCubes)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, []cb.Cube, []cb.Cube, []cb.Cube, int, int, time.Time))
		if !ok {
			panic("invalid mock func, MockFarmService.UpdateFarmCubes()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7)
	}

}

var FightService = &MockFightService{}

type MockFightService struct{}

func (s *MockFightService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockFightService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockFightService.Close()")
		}
		f()
	}

}
func (s *MockFightService) SendFightRequest(a0 *entity.TlogFightContext, a1 *scene.CombatScene, a2 int64, a3 int64, a4 *shared_proto.CombatPlayerProto, a5 *shared_proto.CombatPlayerProto) *server_proto.CombatResponseServerProto {
	fi := getMockFunc(s, s.SendFightRequest)
	if fi != nil {
		f, ok := fi.(func(*entity.TlogFightContext, *scene.CombatScene, int64, int64, *shared_proto.CombatPlayerProto, *shared_proto.CombatPlayerProto) *server_proto.CombatResponseServerProto)
		if !ok {
			panic("invalid mock func, MockFightService.SendFightRequest()")
		}
		return f(a0, a1, a2, a3, a4, a5)
	}

	return nil
}
func (s *MockFightService) SendFightRequestReturnResult(a0 *entity.TlogFightContext, a1 *scene.CombatScene, a2 int64, a3 int64, a4 *shared_proto.CombatPlayerProto, a5 *shared_proto.CombatPlayerProto, a6 bool) *server_proto.CombatResponseServerProto {
	fi := getMockFunc(s, s.SendFightRequestReturnResult)
	if fi != nil {
		f, ok := fi.(func(*entity.TlogFightContext, *scene.CombatScene, int64, int64, *shared_proto.CombatPlayerProto, *shared_proto.CombatPlayerProto, bool) *server_proto.CombatResponseServerProto)
		if !ok {
			panic("invalid mock func, MockFightService.SendFightRequestReturnResult()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6)
	}

	return nil
}
func (s *MockFightService) SendMultiFightRequest(a0 *entity.TlogFightContext, a1 *scene.CombatScene, a2 []int64, a3 []int64, a4 []*shared_proto.CombatPlayerProto, a5 []*shared_proto.CombatPlayerProto, a6 int32, a7 int32, a8 int32) *server_proto.MultiCombatResponseServerProto {
	fi := getMockFunc(s, s.SendMultiFightRequest)
	if fi != nil {
		f, ok := fi.(func(*entity.TlogFightContext, *scene.CombatScene, []int64, []int64, []*shared_proto.CombatPlayerProto, []*shared_proto.CombatPlayerProto, int32, int32, int32) *server_proto.MultiCombatResponseServerProto)
		if !ok {
			panic("invalid mock func, MockFightService.SendMultiFightRequest()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8)
	}

	return nil
}
func (s *MockFightService) SendMultiFightRequestReturnResult(a0 *entity.TlogFightContext, a1 *scene.CombatScene, a2 []int64, a3 []int64, a4 []*shared_proto.CombatPlayerProto, a5 []*shared_proto.CombatPlayerProto, a6 int32, a7 int32, a8 int32, a9 bool) *server_proto.MultiCombatResponseServerProto {
	fi := getMockFunc(s, s.SendMultiFightRequestReturnResult)
	if fi != nil {
		f, ok := fi.(func(*entity.TlogFightContext, *scene.CombatScene, []int64, []int64, []*shared_proto.CombatPlayerProto, []*shared_proto.CombatPlayerProto, int32, int32, int32, bool) *server_proto.MultiCombatResponseServerProto)
		if !ok {
			panic("invalid mock func, MockFightService.SendMultiFightRequestReturnResult()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	}

	return nil
}

var FightXService = &MockFightXService{}

type MockFightXService struct{}

func (s *MockFightXService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockFightXService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockFightXService.Close()")
		}
		f()
	}

}

// 新版战斗
func (s *MockFightXService) SendFightRequest(a0 *entity.TlogFightContext, a1 *scene.CombatScene, a2 int64, a3 int64, a4 *shared_proto.CombatPlayerProto, a5 *shared_proto.CombatPlayerProto) *server_proto.CombatXResponseServerProto {
	fi := getMockFunc(s, s.SendFightRequest)
	if fi != nil {
		f, ok := fi.(func(*entity.TlogFightContext, *scene.CombatScene, int64, int64, *shared_proto.CombatPlayerProto, *shared_proto.CombatPlayerProto) *server_proto.CombatXResponseServerProto)
		if !ok {
			panic("invalid mock func, MockFightXService.SendFightRequest()")
		}
		return f(a0, a1, a2, a3, a4, a5)
	}

	return nil
}
func (s *MockFightXService) SendFightRequestReturnResult(a0 *entity.TlogFightContext, a1 *scene.CombatScene, a2 int64, a3 int64, a4 *shared_proto.CombatPlayerProto, a5 *shared_proto.CombatPlayerProto, a6 bool) *server_proto.CombatXResponseServerProto {
	fi := getMockFunc(s, s.SendFightRequestReturnResult)
	if fi != nil {
		f, ok := fi.(func(*entity.TlogFightContext, *scene.CombatScene, int64, int64, *shared_proto.CombatPlayerProto, *shared_proto.CombatPlayerProto, bool) *server_proto.CombatXResponseServerProto)
		if !ok {
			panic("invalid mock func, MockFightXService.SendFightRequestReturnResult()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6)
	}

	return nil
}

var FishingModule = &MockFishingModule{}

type MockFishingModule struct{}

func (s *MockFishingModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockFishingModule) GmFishingRate(a0 int64, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.GmFishingRate)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockFishingModule.GmFishingRate()")
		}
		f(a0, a1, a2)
	}

}

var GameExporter = &MockGameExporter{}

type MockGameExporter struct{}

func (s *MockGameExporter) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockGameExporter) GetMsgTimeCost() *metrics.MsgTimeCostSummary {
	fi := getMockFunc(s, s.GetMsgTimeCost)
	if fi != nil {
		f, ok := fi.(func() *metrics.MsgTimeCostSummary)
		if !ok {
			panic("invalid mock func, MockGameExporter.GetMsgTimeCost()")
		}
		return f()
	}

	return nil
}
func (s *MockGameExporter) GetRegisterCounter() *atomic.Uint64 {
	fi := getMockFunc(s, s.GetRegisterCounter)
	if fi != nil {
		f, ok := fi.(func() *atomic.Uint64)
		if !ok {
			panic("invalid mock func, MockGameExporter.GetRegisterCounter()")
		}
		return f()
	}

	return nil
}
func (s *MockGameExporter) Start() (http.Handler, error) {
	fi := getMockFunc(s, s.Start)
	if fi != nil {
		f, ok := fi.(func() (http.Handler, error))
		if !ok {
			panic("invalid mock func, MockGameExporter.Start()")
		}
		return f()
	}

	return nil, nil
}

var GameServer = &MockGameServer{}

type MockGameServer struct{}

func (s *MockGameServer) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockGameServer) GetRpcPort() uint32 {
	fi := getMockFunc(s, s.GetRpcPort)
	if fi != nil {
		f, ok := fi.(func() uint32)
		if !ok {
			panic("invalid mock func, MockGameServer.GetRpcPort()")
		}
		return f()
	}

	return 0
}
func (s *MockGameServer) GetTcpPort() uint32 {
	fi := getMockFunc(s, s.GetTcpPort)
	if fi != nil {
		f, ok := fi.(func() uint32)
		if !ok {
			panic("invalid mock func, MockGameServer.GetTcpPort()")
		}
		return f()
	}

	return 0
}
func (s *MockGameServer) Serve(a0 iface.ServeListener, a1 iface.ConnHandler) {
	fi := getMockFunc(s, s.Serve)
	if fi != nil {
		f, ok := fi.(func(iface.ServeListener, iface.ConnHandler))
		if !ok {
			panic("invalid mock func, MockGameServer.Serve()")
		}
		f(a0, a1)
	}

}

var GardenModule = &MockGardenModule{}

type MockGardenModule struct{}

func (s *MockGardenModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

// 宝石模块

var GemModule = &MockGemModule{}

type MockGemModule struct{}

func (s *MockGemModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var GmModule = &MockGmModule{}

type MockGmModule struct{}

func (s *MockGmModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var GuildModule = &MockGuildModule{}

type MockGuildModule struct{}

func (s *MockGuildModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockGuildModule) GmAddBigBoxEnergy(a0 int64, a1 uint64) {
	fi := getMockFunc(s, s.GmAddBigBoxEnergy)
	if fi != nil {
		f, ok := fi.(func(int64, uint64))
		if !ok {
			panic("invalid mock func, MockGuildModule.GmAddBigBoxEnergy()")
		}
		f(a0, a1)
	}

}
func (s *MockGuildModule) GmAddGuildBuildAmount(a0 int64, a1 uint64) {
	fi := getMockFunc(s, s.GmAddGuildBuildAmount)
	if fi != nil {
		f, ok := fi.(func(int64, uint64))
		if !ok {
			panic("invalid mock func, MockGuildModule.GmAddGuildBuildAmount()")
		}
		f(a0, a1)
	}

}
func (s *MockGuildModule) GmAddGuildYinliang(a0 int64, a1 int64) {
	fi := getMockFunc(s, s.GmAddGuildYinliang)
	if fi != nil {
		f, ok := fi.(func(int64, int64))
		if !ok {
			panic("invalid mock func, MockGuildModule.GmAddGuildYinliang()")
		}
		f(a0, a1)
	}

}
func (s *MockGuildModule) GmGiveGuildEventPrize(a0 *entity.Hero, a1 herolock.LockResult, a2 []*guild_data.GuildEventPrizeData) {
	fi := getMockFunc(s, s.GmGiveGuildEventPrize)
	if fi != nil {
		f, ok := fi.(func(*entity.Hero, herolock.LockResult, []*guild_data.GuildEventPrizeData))
		if !ok {
			panic("invalid mock func, MockGuildModule.GmGiveGuildEventPrize()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockGuildModule) GmMiaoGuildTechCd(a0 int64) {
	fi := getMockFunc(s, s.GmMiaoGuildTechCd)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockGuildModule.GmMiaoGuildTechCd()")
		}
		f(a0)
	}

}
func (s *MockGuildModule) GmOpenImpeachLeader(a0 int64) {
	fi := getMockFunc(s, s.GmOpenImpeachLeader)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockGuildModule.GmOpenImpeachLeader()")
		}
		f(a0)
	}

}
func (s *MockGuildModule) GmRemoveNpcGuild() {
	fi := getMockFunc(s, s.GmRemoveNpcGuild)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockGuildModule.GmRemoveNpcGuild()")
		}
		f()
	}

}
func (s *MockGuildModule) GmUpgradeGuildLevel(a0 int64) {
	fi := getMockFunc(s, s.GmUpgradeGuildLevel)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockGuildModule.GmUpgradeGuildLevel()")
		}
		f(a0)
	}

}
func (s *MockGuildModule) HandleGiveGuildEventPrize(a0 *entity.Hero, a1 herolock.LockResult, a2 int64, a3 []*guild_data.GuildEventPrizeData, a4 uint64) {
	fi := getMockFunc(s, s.HandleGiveGuildEventPrize)
	if fi != nil {
		f, ok := fi.(func(*entity.Hero, herolock.LockResult, int64, []*guild_data.GuildEventPrizeData, uint64))
		if !ok {
			panic("invalid mock func, MockGuildModule.HandleGiveGuildEventPrize()")
		}
		f(a0, a1, a2, a3, a4)
	}

}
func (s *MockGuildModule) OnHeroOnline(a0 iface.HeroController, a1 int64) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64))
		if !ok {
			panic("invalid mock func, MockGuildModule.OnHeroOnline()")
		}
		f(a0, a1)
	}

}

var GuildService = &MockGuildService{}

type MockGuildService struct{}

func (s *MockGuildService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

// 增加联盟任务进度
func (s *MockGuildService) AddGuildTaskProgress(a0 int64, a1 *guild_data.GuildTaskData, a2 uint64) {
	fi := getMockFunc(s, s.AddGuildTaskProgress)
	if fi != nil {
		f, ok := fi.(func(int64, *guild_data.GuildTaskData, uint64))
		if !ok {
			panic("invalid mock func, MockGuildService.AddGuildTaskProgress()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockGuildService) AddHufu(a0 uint64, a1 int64, a2 int64, a3 string, a4 string) bool {
	fi := getMockFunc(s, s.AddHufu)
	if fi != nil {
		f, ok := fi.(func(uint64, int64, int64, string, string) bool)
		if !ok {
			panic("invalid mock func, MockGuildService.AddHufu()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return false
}
func (s *MockGuildService) AddLog(a0 int64, a1 *shared_proto.GuildLogProto) {
	fi := getMockFunc(s, s.AddLog)
	if fi != nil {
		f, ok := fi.(func(int64, *shared_proto.GuildLogProto))
		if !ok {
			panic("invalid mock func, MockGuildService.AddLog()")
		}
		f(a0, a1)
	}

}
func (s *MockGuildService) AddLogWithMemberIds(a0 int64, a1 []int64, a2 *shared_proto.GuildLogProto) {
	fi := getMockFunc(s, s.AddLogWithMemberIds)
	if fi != nil {
		f, ok := fi.(func(int64, []int64, *shared_proto.GuildLogProto))
		if !ok {
			panic("invalid mock func, MockGuildService.AddLogWithMemberIds()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockGuildService) AddRecommendInviteHeros(a0 int64) {
	fi := getMockFunc(s, s.AddRecommendInviteHeros)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockGuildService.AddRecommendInviteHeros()")
		}
		f(a0)
	}

}
func (s *MockGuildService) Broadcast(a0 int64, a1 pbutil.Buffer) {
	fi := getMockFunc(s, s.Broadcast)
	if fi != nil {
		f, ok := fi.(func(int64, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockGuildService.Broadcast()")
		}
		f(a0, a1)
	}

}
func (s *MockGuildService) CheckAndAddRecommendInviteHeros(a0 int64) {
	fi := getMockFunc(s, s.CheckAndAddRecommendInviteHeros)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockGuildService.CheckAndAddRecommendInviteHeros()")
		}
		f(a0)
	}

}
func (s *MockGuildService) ClearSelfGuildMsgCache(a0 int64) {
	fi := getMockFunc(s, s.ClearSelfGuildMsgCache)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockGuildService.ClearSelfGuildMsgCache()")
		}
		f(a0)
	}

}
func (s *MockGuildService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockGuildService.Close()")
		}
		f()
	}

}
func (s *MockGuildService) Func(a0 sharedguilddata.Funcs) bool {
	fi := getMockFunc(s, s.Func)
	if fi != nil {
		f, ok := fi.(func(sharedguilddata.Funcs) bool)
		if !ok {
			panic("invalid mock func, MockGuildService.Func()")
		}
		return f(a0)
	}

	return false
}
func (s *MockGuildService) FuncGuild(a0 int64, a1 sharedguilddata.Func) {
	fi := getMockFunc(s, s.FuncGuild)
	if fi != nil {
		f, ok := fi.(func(int64, sharedguilddata.Func))
		if !ok {
			panic("invalid mock func, MockGuildService.FuncGuild()")
		}
		f(a0, a1)
	}

}
func (s *MockGuildService) GetGuildFlagName(a0 int64) string {
	fi := getMockFunc(s, s.GetGuildFlagName)
	if fi != nil {
		f, ok := fi.(func(int64) string)
		if !ok {
			panic("invalid mock func, MockGuildService.GetGuildFlagName()")
		}
		return f(a0)
	}

	return ""
}
func (s *MockGuildService) GetGuildIdByFlagName(a0 string) int64 {
	fi := getMockFunc(s, s.GetGuildIdByFlagName)
	if fi != nil {
		f, ok := fi.(func(string) int64)
		if !ok {
			panic("invalid mock func, MockGuildService.GetGuildIdByFlagName()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockGuildService) GetGuildIdByName(a0 string) int64 {
	fi := getMockFunc(s, s.GetGuildIdByName)
	if fi != nil {
		f, ok := fi.(func(string) int64)
		if !ok {
			panic("invalid mock func, MockGuildService.GetGuildIdByName()")
		}
		return f(a0)
	}

	return 0
}

// 根据国家id获取声望排名列表的消息
func (s *MockGuildService) GetGuildPrestigeRankMsg(a0 uint64, a1 time.Time) pbutil.Buffer {
	fi := getMockFunc(s, s.GetGuildPrestigeRankMsg)
	if fi != nil {
		f, ok := fi.(func(uint64, time.Time) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockGuildService.GetGuildPrestigeRankMsg()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockGuildService) GetSnapshot(a0 int64) *guildsnapshotdata.GuildSnapshot {
	fi := getMockFunc(s, s.GetSnapshot)
	if fi != nil {
		f, ok := fi.(func(int64) *guildsnapshotdata.GuildSnapshot)
		if !ok {
			panic("invalid mock func, MockGuildService.GetSnapshot()")
		}
		return f(a0)
	}

	return nil
}

// 筛选推荐联盟
func (s *MockGuildService) RecommendGuildList(a0 *snapshotdata.HeroSnapshot) []int64 {
	fi := getMockFunc(s, s.RecommendGuildList)
	if fi != nil {
		f, ok := fi.(func(*snapshotdata.HeroSnapshot) []int64)
		if !ok {
			panic("invalid mock func, MockGuildService.RecommendGuildList()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockGuildService) RecommendInviteHeroList(a0 uint64, a1 uint64, a2 []int64) []*snapshotdata.HeroSnapshot {
	fi := getMockFunc(s, s.RecommendInviteHeroList)
	if fi != nil {
		f, ok := fi.(func(uint64, uint64, []int64) []*snapshotdata.HeroSnapshot)
		if !ok {
			panic("invalid mock func, MockGuildService.RecommendInviteHeroList()")
		}
		return f(a0, a1, a2)
	}

	return nil
}
func (s *MockGuildService) RegisterCallback(a0 guildsnapshotdata.Callback) {
	fi := getMockFunc(s, s.RegisterCallback)
	if fi != nil {
		f, ok := fi.(func(guildsnapshotdata.Callback))
		if !ok {
			panic("invalid mock func, MockGuildService.RegisterCallback()")
		}
		f(a0)
	}

}
func (s *MockGuildService) RemoveRecommendInviteHero(a0 int64) {
	fi := getMockFunc(s, s.RemoveRecommendInviteHero)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockGuildService.RemoveRecommendInviteHero()")
		}
		f(a0)
	}

}
func (s *MockGuildService) RemoveSnapshot(a0 int64) {
	fi := getMockFunc(s, s.RemoveSnapshot)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockGuildService.RemoveSnapshot()")
		}
		f(a0)
	}

}
func (s *MockGuildService) SaveChangedGuild() {
	fi := getMockFunc(s, s.SaveChangedGuild)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockGuildService.SaveChangedGuild()")
		}
		f()
	}

}
func (s *MockGuildService) SelfGuildMsgCache() concurrent.I64BufferMap {
	fi := getMockFunc(s, s.SelfGuildMsgCache)
	if fi != nil {
		f, ok := fi.(func() concurrent.I64BufferMap)
		if !ok {
			panic("invalid mock func, MockGuildService.SelfGuildMsgCache()")
		}
		return f()
	}

	return nil
}
func (s *MockGuildService) SetGuildRankFunc(a0 sharedguilddata.GetGuildRankFunc, a1 sharedguilddata.GenerateRankMsgFunc) {
	fi := getMockFunc(s, s.SetGuildRankFunc)
	if fi != nil {
		f, ok := fi.(func(sharedguilddata.GetGuildRankFunc, sharedguilddata.GenerateRankMsgFunc))
		if !ok {
			panic("invalid mock func, MockGuildService.SetGuildRankFunc()")
		}
		f(a0, a1)
	}

}
func (s *MockGuildService) TimeoutFunc(a0 sharedguilddata.Funcs) bool {
	fi := getMockFunc(s, s.TimeoutFunc)
	if fi != nil {
		f, ok := fi.(func(sharedguilddata.Funcs) bool)
		if !ok {
			panic("invalid mock func, MockGuildService.TimeoutFunc()")
		}
		return f(a0)
	}

	return false
}
func (s *MockGuildService) UpdateGuildHeroSnapshot(a0 *sharedguilddata.Guild) {
	fi := getMockFunc(s, s.UpdateGuildHeroSnapshot)
	if fi != nil {
		f, ok := fi.(func(*sharedguilddata.Guild))
		if !ok {
			panic("invalid mock func, MockGuildService.UpdateGuildHeroSnapshot()")
		}
		f(a0)
	}

}
func (s *MockGuildService) UpdateSnapshot(a0 *sharedguilddata.Guild) *guildsnapshotdata.GuildSnapshot {
	fi := getMockFunc(s, s.UpdateSnapshot)
	if fi != nil {
		f, ok := fi.(func(*sharedguilddata.Guild) *guildsnapshotdata.GuildSnapshot)
		if !ok {
			panic("invalid mock func, MockGuildService.UpdateSnapshot()")
		}
		return f(a0)
	}

	return nil
}

var GuildSnapshotService = &MockGuildSnapshotService{}

type MockGuildSnapshotService struct{}

func (s *MockGuildSnapshotService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockGuildSnapshotService) GetGuildBasicProto(a0 int64) *shared_proto.GuildBasicProto {
	fi := getMockFunc(s, s.GetGuildBasicProto)
	if fi != nil {
		f, ok := fi.(func(int64) *shared_proto.GuildBasicProto)
		if !ok {
			panic("invalid mock func, MockGuildSnapshotService.GetGuildBasicProto()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockGuildSnapshotService) GetGuildLevel(a0 int64) uint64 {
	fi := getMockFunc(s, s.GetGuildLevel)
	if fi != nil {
		f, ok := fi.(func(int64) uint64)
		if !ok {
			panic("invalid mock func, MockGuildSnapshotService.GetGuildLevel()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockGuildSnapshotService) GetSnapshot(a0 int64) *guildsnapshotdata.GuildSnapshot {
	fi := getMockFunc(s, s.GetSnapshot)
	if fi != nil {
		f, ok := fi.(func(int64) *guildsnapshotdata.GuildSnapshot)
		if !ok {
			panic("invalid mock func, MockGuildSnapshotService.GetSnapshot()")
		}
		return f(a0)
	}

	return nil
}

// 注册监听snapshot变化callback
func (s *MockGuildSnapshotService) RegisterCallback(a0 guildsnapshotdata.Callback) {
	fi := getMockFunc(s, s.RegisterCallback)
	if fi != nil {
		f, ok := fi.(func(guildsnapshotdata.Callback))
		if !ok {
			panic("invalid mock func, MockGuildSnapshotService.RegisterCallback()")
		}
		f(a0)
	}

}
func (s *MockGuildSnapshotService) RemoveSnapshot(a0 int64) {
	fi := getMockFunc(s, s.RemoveSnapshot)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockGuildSnapshotService.RemoveSnapshot()")
		}
		f(a0)
	}

}
func (s *MockGuildSnapshotService) UpdateSnapshot(a0 *guildsnapshotdata.GuildSnapshot) {
	fi := getMockFunc(s, s.UpdateSnapshot)
	if fi != nil {
		f, ok := fi.(func(*guildsnapshotdata.GuildSnapshot))
		if !ok {
			panic("invalid mock func, MockGuildSnapshotService.UpdateSnapshot()")
		}
		f(a0)
	}

}

var HebiModule = &MockHebiModule{}

type MockHebiModule struct{}

func (s *MockHebiModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockHebiModule) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockHebiModule.Close()")
		}
		f()
	}

}
func (s *MockHebiModule) UpdateGuildInfo(a0 int64, a1 int64) {
	fi := getMockFunc(s, s.UpdateGuildInfo)
	if fi != nil {
		f, ok := fi.(func(int64, int64))
		if !ok {
			panic("invalid mock func, MockHebiModule.UpdateGuildInfo()")
		}
		f(a0, a1)
	}

}
func (s *MockHebiModule) UpdateGuildInfoBatch(a0 []int64, a1 int64) {
	fi := getMockFunc(s, s.UpdateGuildInfoBatch)
	if fi != nil {
		f, ok := fi.(func([]int64, int64))
		if !ok {
			panic("invalid mock func, MockHebiModule.UpdateGuildInfoBatch()")
		}
		f(a0, a1)
	}

}

var HeroController = &MockHeroController{}

type MockHeroController struct{}

func (s *MockHeroController) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockHeroController) AddTickFunc(a0 face.BFunc) {
	fi := getMockFunc(s, s.AddTickFunc)
	if fi != nil {
		f, ok := fi.(func(face.BFunc))
		if !ok {
			panic("invalid mock func, MockHeroController.AddTickFunc()")
		}
		f(a0)
	}

}
func (s *MockHeroController) Disconnect(a0 msg.ErrMsg) {
	fi := getMockFunc(s, s.Disconnect)
	if fi != nil {
		f, ok := fi.(func(msg.ErrMsg))
		if !ok {
			panic("invalid mock func, MockHeroController.Disconnect()")
		}
		f(a0)
	}

}
func (s *MockHeroController) Func(a0 herolock.Func) {
	fi := getMockFunc(s, s.Func)
	if fi != nil {
		f, ok := fi.(func(herolock.Func))
		if !ok {
			panic("invalid mock func, MockHeroController.Func()")
		}
		f(a0)
	}

}
func (s *MockHeroController) FuncNotError(a0 herolock.FuncNotError) bool {
	fi := getMockFunc(s, s.FuncNotError)
	if fi != nil {
		f, ok := fi.(func(herolock.FuncNotError) bool)
		if !ok {
			panic("invalid mock func, MockHeroController.FuncNotError()")
		}
		return f(a0)
	}

	return false
}
func (s *MockHeroController) FuncWithSend(a0 herolock.SendFunc) bool {
	fi := getMockFunc(s, s.FuncWithSend)
	if fi != nil {
		f, ok := fi.(func(herolock.SendFunc) bool)
		if !ok {
			panic("invalid mock func, MockHeroController.FuncWithSend()")
		}
		return f(a0)
	}

	return false
}

//GetBlockIndex 获取地块索引
func (s *MockHeroController) GetBlockIndex() interface{} {
	fi := getMockFunc(s, s.GetBlockIndex)
	if fi != nil {
		f, ok := fi.(func() interface{})
		if !ok {
			panic("invalid mock func, MockHeroController.GetBlockIndex()")
		}
		return f()
	}

	return 0
}
func (s *MockHeroController) GetCareCondition() *server_proto.MilitaryConditionProto {
	fi := getMockFunc(s, s.GetCareCondition)
	if fi != nil {
		f, ok := fi.(func() *server_proto.MilitaryConditionProto)
		if !ok {
			panic("invalid mock func, MockHeroController.GetCareCondition()")
		}
		return f()
	}

	return nil
}
func (s *MockHeroController) GetCareWaterTimesMap() map[int64]uint64 {
	fi := getMockFunc(s, s.GetCareWaterTimesMap)
	if fi != nil {
		f, ok := fi.(func() map[int64]uint64)
		if !ok {
			panic("invalid mock func, MockHeroController.GetCareWaterTimesMap()")
		}
		return f()
	}

	return nil
}
func (s *MockHeroController) GetClientIp() string {
	fi := getMockFunc(s, s.GetClientIp)
	if fi != nil {
		f, ok := fi.(func() string)
		if !ok {
			panic("invalid mock func, MockHeroController.GetClientIp()")
		}
		return f()
	}

	return ""
}
func (s *MockHeroController) GetClientIp32() uint32 {
	fi := getMockFunc(s, s.GetClientIp32)
	if fi != nil {
		f, ok := fi.(func() uint32)
		if !ok {
			panic("invalid mock func, MockHeroController.GetClientIp32()")
		}
		return f()
	}

	return 0
}
func (s *MockHeroController) GetIsInBackgroud() bool {
	fi := getMockFunc(s, s.GetIsInBackgroud)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockHeroController.GetIsInBackgroud()")
		}
		return f()
	}

	return false
}
func (s *MockHeroController) GetPf() uint32 {
	fi := getMockFunc(s, s.GetPf)
	if fi != nil {
		f, ok := fi.(func() uint32)
		if !ok {
			panic("invalid mock func, MockHeroController.GetPf()")
		}
		return f()
	}

	return 0
}
func (s *MockHeroController) GetViewArea() *realmface.ViewArea {
	fi := getMockFunc(s, s.GetViewArea)
	if fi != nil {
		f, ok := fi.(func() *realmface.ViewArea)
		if !ok {
			panic("invalid mock func, MockHeroController.GetViewArea()")
		}
		return f()
	}

	return nil
}

//GetWatchObjList 获取观察列表
func (s *MockHeroController) GetWatchObjList() map[interface{}]int {
	fi := getMockFunc(s, s.GetWatchObjList)
	if fi != nil {
		f, ok := fi.(func() map[interface{}]int)
		if !ok {
			panic("invalid mock func, MockHeroController.GetWatchObjList()")
		}
		return f()
	}

	return nil
}
func (s *MockHeroController) Id() int64 {
	fi := getMockFunc(s, s.Id)
	if fi != nil {
		f, ok := fi.(func() int64)
		if !ok {
			panic("invalid mock func, MockHeroController.Id()")
		}
		return f()
	}

	return 0
}
func (s *MockHeroController) IdBytes() []byte {
	fi := getMockFunc(s, s.IdBytes)
	if fi != nil {
		f, ok := fi.(func() []byte)
		if !ok {
			panic("invalid mock func, MockHeroController.IdBytes()")
		}
		return f()
	}

	return nil
}
func (s *MockHeroController) IsClosed() bool {
	fi := getMockFunc(s, s.IsClosed)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockHeroController.IsClosed()")
		}
		return f()
	}

	return false
}
func (s *MockHeroController) LastClickTime() time.Time {
	fi := getMockFunc(s, s.LastClickTime)
	if fi != nil {
		f, ok := fi.(func() time.Time)
		if !ok {
			panic("invalid mock func, MockHeroController.LastClickTime()")
		}
		return f()
	}

	return time.Time{}
}
func (s *MockHeroController) LockGetGuildId() (int64, bool) {
	fi := getMockFunc(s, s.LockGetGuildId)
	if fi != nil {
		f, ok := fi.(func() (int64, bool))
		if !ok {
			panic("invalid mock func, MockHeroController.LockGetGuildId()")
		}
		return f()
	}

	return 0, false
}
func (s *MockHeroController) LockHeroCountry() uint64 {
	fi := getMockFunc(s, s.LockHeroCountry)
	if fi != nil {
		f, ok := fi.(func() uint64)
		if !ok {
			panic("invalid mock func, MockHeroController.LockHeroCountry()")
		}
		return f()
	}

	return 0
}
func (s *MockHeroController) NextRefreshRecommendHeroTime() time.Time {
	fi := getMockFunc(s, s.NextRefreshRecommendHeroTime)
	if fi != nil {
		f, ok := fi.(func() time.Time)
		if !ok {
			panic("invalid mock func, MockHeroController.NextRefreshRecommendHeroTime()")
		}
		return f()
	}

	return time.Time{}
}
func (s *MockHeroController) NextSearchHeroTime() time.Time {
	fi := getMockFunc(s, s.NextSearchHeroTime)
	if fi != nil {
		f, ok := fi.(func() time.Time)
		if !ok {
			panic("invalid mock func, MockHeroController.NextSearchHeroTime()")
		}
		return f()
	}

	return time.Time{}
}
func (s *MockHeroController) NextSearchNoGuildHeros() time.Time {
	fi := getMockFunc(s, s.NextSearchNoGuildHeros)
	if fi != nil {
		f, ok := fi.(func() time.Time)
		if !ok {
			panic("invalid mock func, MockHeroController.NextSearchNoGuildHeros()")
		}
		return f()
	}

	return time.Time{}
}
func (s *MockHeroController) Pid() uint32 {
	fi := getMockFunc(s, s.Pid)
	if fi != nil {
		f, ok := fi.(func() uint32)
		if !ok {
			panic("invalid mock func, MockHeroController.Pid()")
		}
		return f()
	}

	return 0
}

//RemoveBlockIndex 退出时 移除出场景
func (s *MockHeroController) RemoveBlockIndex() {
	fi := getMockFunc(s, s.RemoveBlockIndex)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockHeroController.RemoveBlockIndex()")
		}
		f()
	}

}

// 发送消息.
func (s *MockHeroController) Send(a0 pbutil.Buffer) {
	fi := getMockFunc(s, s.Send)
	if fi != nil {
		f, ok := fi.(func(pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockHeroController.Send()")
		}
		f(a0)
	}

}

// 发送消息.
func (s *MockHeroController) SendAll(a0 []pbutil.Buffer) {
	fi := getMockFunc(s, s.SendAll)
	if fi != nil {
		f, ok := fi.(func([]pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockHeroController.SendAll()")
		}
		f(a0)
	}

}

// 发送在线路繁忙时可以被丢掉的消息
func (s *MockHeroController) SendIfFree(a0 pbutil.Buffer) {
	fi := getMockFunc(s, s.SendIfFree)
	if fi != nil {
		f, ok := fi.(func(pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockHeroController.SendIfFree()")
		}
		f(a0)
	}

}

//SetBlockIndex 设置地块索引(AOI)
func (s *MockHeroController) SetBlockIndex(a0 interface{}) interface{} {
	fi := getMockFunc(s, s.SetBlockIndex)
	if fi != nil {
		f, ok := fi.(func(interface{}) interface{})
		if !ok {
			panic("invalid mock func, MockHeroController.SetBlockIndex()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockHeroController) SetCareCondition(a0 *server_proto.MilitaryConditionProto) {
	fi := getMockFunc(s, s.SetCareCondition)
	if fi != nil {
		f, ok := fi.(func(*server_proto.MilitaryConditionProto))
		if !ok {
			panic("invalid mock func, MockHeroController.SetCareCondition()")
		}
		f(a0)
	}

}
func (s *MockHeroController) SetCareWaterTimesMap(a0 map[int64]uint64) {
	fi := getMockFunc(s, s.SetCareWaterTimesMap)
	if fi != nil {
		f, ok := fi.(func(map[int64]uint64))
		if !ok {
			panic("invalid mock func, MockHeroController.SetCareWaterTimesMap()")
		}
		f(a0)
	}

}
func (s *MockHeroController) SetIsInBackgroud(a0 time.Time, a1 bool) {
	fi := getMockFunc(s, s.SetIsInBackgroud)
	if fi != nil {
		f, ok := fi.(func(time.Time, bool))
		if !ok {
			panic("invalid mock func, MockHeroController.SetIsInBackgroud()")
		}
		f(a0, a1)
	}

}
func (s *MockHeroController) SetLastClickTime(a0 time.Time) {
	fi := getMockFunc(s, s.SetLastClickTime)
	if fi != nil {
		f, ok := fi.(func(time.Time))
		if !ok {
			panic("invalid mock func, MockHeroController.SetLastClickTime()")
		}
		f(a0)
	}

}
func (s *MockHeroController) SetNextSearchNoGuildHeros(a0 time.Time) {
	fi := getMockFunc(s, s.SetNextSearchNoGuildHeros)
	if fi != nil {
		f, ok := fi.(func(time.Time))
		if !ok {
			panic("invalid mock func, MockHeroController.SetNextSearchNoGuildHeros()")
		}
		f(a0)
	}

}
func (s *MockHeroController) SetViewArea(a0 *realmface.ViewArea) {
	fi := getMockFunc(s, s.SetViewArea)
	if fi != nil {
		f, ok := fi.(func(*realmface.ViewArea))
		if !ok {
			panic("invalid mock func, MockHeroController.SetViewArea()")
		}
		f(a0)
	}

}

//AddWatchObjList 设置观察对象列表 如果设置nil,表示清空
func (s *MockHeroController) SetWatchObjList(a0 map[interface{}]int) {
	fi := getMockFunc(s, s.SetWatchObjList)
	if fi != nil {
		f, ok := fi.(func(map[interface{}]int))
		if !ok {
			panic("invalid mock func, MockHeroController.SetWatchObjList()")
		}
		f(a0)
	}

}
func (s *MockHeroController) Sid() uint32 {
	fi := getMockFunc(s, s.Sid)
	if fi != nil {
		f, ok := fi.(func() uint32)
		if !ok {
			panic("invalid mock func, MockHeroController.Sid()")
		}
		return f()
	}

	return 0
}
func (s *MockHeroController) TickFunc() {
	fi := getMockFunc(s, s.TickFunc)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockHeroController.TickFunc()")
		}
		f()
	}

}
func (s *MockHeroController) TotalOnlineTime() time.Duration {
	fi := getMockFunc(s, s.TotalOnlineTime)
	if fi != nil {
		f, ok := fi.(func() time.Duration)
		if !ok {
			panic("invalid mock func, MockHeroController.TotalOnlineTime()")
		}
		return f()
	}

	return 0
}
func (s *MockHeroController) TryNextWriteOnlineLogTime(a0 time.Time, a1 time.Duration) bool {
	fi := getMockFunc(s, s.TryNextWriteOnlineLogTime)
	if fi != nil {
		f, ok := fi.(func(time.Time, time.Duration) bool)
		if !ok {
			panic("invalid mock func, MockHeroController.TryNextWriteOnlineLogTime()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockHeroController) UpdateIsInBackgroud(a0 time.Time) {
	fi := getMockFunc(s, s.UpdateIsInBackgroud)
	if fi != nil {
		f, ok := fi.(func(time.Time))
		if !ok {
			panic("invalid mock func, MockHeroController.UpdateIsInBackgroud()")
		}
		f(a0)
	}

}
func (s *MockHeroController) UpdateNextRefreshRecommendHeroTime(a0 time.Time) {
	fi := getMockFunc(s, s.UpdateNextRefreshRecommendHeroTime)
	if fi != nil {
		f, ok := fi.(func(time.Time))
		if !ok {
			panic("invalid mock func, MockHeroController.UpdateNextRefreshRecommendHeroTime()")
		}
		f(a0)
	}

}
func (s *MockHeroController) UpdateNextSearchHeroTime(a0 time.Time) {
	fi := getMockFunc(s, s.UpdateNextSearchHeroTime)
	if fi != nil {
		f, ok := fi.(func(time.Time))
		if !ok {
			panic("invalid mock func, MockHeroController.UpdateNextSearchHeroTime()")
		}
		f(a0)
	}

}

var HeroDataService = &MockHeroDataService{}

type MockHeroDataService struct{}

func (s *MockHeroDataService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockHeroDataService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockHeroDataService.Close()")
		}
		f()
	}

}
func (s *MockHeroDataService) Create(a0 *entity.Hero) error {
	fi := getMockFunc(s, s.Create)
	if fi != nil {
		f, ok := fi.(func(*entity.Hero) error)
		if !ok {
			panic("invalid mock func, MockHeroDataService.Create()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockHeroDataService) Exist(a0 int64) (bool, error) {
	fi := getMockFunc(s, s.Exist)
	if fi != nil {
		f, ok := fi.(func(int64) (bool, error))
		if !ok {
			panic("invalid mock func, MockHeroDataService.Exist()")
		}
		return f(a0)
	}

	return false, nil
}
func (s *MockHeroDataService) Func(a0 int64, a1 herolock.Func) {
	fi := getMockFunc(s, s.Func)
	if fi != nil {
		f, ok := fi.(func(int64, herolock.Func))
		if !ok {
			panic("invalid mock func, MockHeroDataService.Func()")
		}
		f(a0, a1)
	}

}
func (s *MockHeroDataService) FuncNotError(a0 int64, a1 herolock.FuncNotError) bool {
	fi := getMockFunc(s, s.FuncNotError)
	if fi != nil {
		f, ok := fi.(func(int64, herolock.FuncNotError) bool)
		if !ok {
			panic("invalid mock func, MockHeroDataService.FuncNotError()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockHeroDataService) FuncWithSend(a0 int64, a1 herolock.SendFunc) bool {
	fi := getMockFunc(s, s.FuncWithSend)
	if fi != nil {
		f, ok := fi.(func(int64, herolock.SendFunc) bool)
		if !ok {
			panic("invalid mock func, MockHeroDataService.FuncWithSend()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockHeroDataService) FuncWithSendError(a0 int64, a1 herolock.SendFunc) (bool, error) {
	fi := getMockFunc(s, s.FuncWithSendError)
	if fi != nil {
		f, ok := fi.(func(int64, herolock.SendFunc) (bool, error))
		if !ok {
			panic("invalid mock func, MockHeroDataService.FuncWithSendError()")
		}
		return f(a0, a1)
	}

	return false, nil
}
func (s *MockHeroDataService) NewHeroLocker(a0 int64) herolock.HeroLocker {
	fi := getMockFunc(s, s.NewHeroLocker)
	if fi != nil {
		f, ok := fi.(func(int64) herolock.HeroLocker)
		if !ok {
			panic("invalid mock func, MockHeroDataService.NewHeroLocker()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockHeroDataService) Put(a0 *entity.Hero) error {
	fi := getMockFunc(s, s.Put)
	if fi != nil {
		f, ok := fi.(func(*entity.Hero) error)
		if !ok {
			panic("invalid mock func, MockHeroDataService.Put()")
		}
		return f(a0)
	}

	return nil
}

var HeroSnapshotService = &MockHeroSnapshotService{}

type MockHeroSnapshotService struct{}

func (s *MockHeroSnapshotService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

// 获得英雄的snapshot, 并不保证是最新的. 尽量
// 就算英雄要从db中加载, 也不会触发callback.
// 返回nil也可能是数据库报错, 英雄未必不存在
func (s *MockHeroSnapshotService) Get(a0 int64) *snapshotdata.HeroSnapshot {
	fi := getMockFunc(s, s.Get)
	if fi != nil {
		f, ok := fi.(func(int64) *snapshotdata.HeroSnapshot)
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.Get()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockHeroSnapshotService) GetBasicProto(a0 int64) *shared_proto.HeroBasicProto {
	fi := getMockFunc(s, s.GetBasicProto)
	if fi != nil {
		f, ok := fi.(func(int64) *shared_proto.HeroBasicProto)
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.GetBasicProto()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockHeroSnapshotService) GetBasicSnapshotProto(a0 int64) *shared_proto.HeroBasicSnapshotProto {
	fi := getMockFunc(s, s.GetBasicSnapshotProto)
	if fi != nil {
		f, ok := fi.(func(int64) *shared_proto.HeroBasicSnapshotProto)
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.GetBasicSnapshotProto()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockHeroSnapshotService) GetFlagHeroName(a0 int64) string {
	fi := getMockFunc(s, s.GetFlagHeroName)
	if fi != nil {
		f, ok := fi.(func(int64) string)
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.GetFlagHeroName()")
		}
		return f(a0)
	}

	return ""
}

// 只从Cache中获取，不读取DB
func (s *MockHeroSnapshotService) GetFromCache(a0 int64) *snapshotdata.HeroSnapshot {
	fi := getMockFunc(s, s.GetFromCache)
	if fi != nil {
		f, ok := fi.(func(int64) *snapshotdata.HeroSnapshot)
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.GetFromCache()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockHeroSnapshotService) GetHeroName(a0 int64) string {
	fi := getMockFunc(s, s.GetHeroName)
	if fi != nil {
		f, ok := fi.(func(int64) string)
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.GetHeroName()")
		}
		return f(a0)
	}

	return ""
}
func (s *MockHeroSnapshotService) GetTlogHero(a0 int64) entity.TlogHero {
	fi := getMockFunc(s, s.GetTlogHero)
	if fi != nil {
		f, ok := fi.(func(int64) entity.TlogHero)
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.GetTlogHero()")
		}
		return f(a0)
	}

	return nil
}

// 创建个新的snapshot, 但是并没有保存. 等unlock后再调用Cache保存
// 必须是lock住Hero的情况下才能调用, 确保此时没有其他人能访问hero对象
func (s *MockHeroSnapshotService) NewSnapshot(a0 *entity.Hero) *snapshotdata.HeroSnapshot {
	fi := getMockFunc(s, s.NewSnapshot)
	if fi != nil {
		f, ok := fi.(func(*entity.Hero) *snapshotdata.HeroSnapshot)
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.NewSnapshot()")
		}
		return f(a0)
	}

	return nil
}

// 英雄下线时调用, 把snapshot移动到lru中
func (s *MockHeroSnapshotService) Offline(a0 int64) {
	fi := getMockFunc(s, s.Offline)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.Offline()")
		}
		f(a0)
	}

}

// 英雄上线时调用, 保存snapshot. 这个snapshot必须是没有变化的数据的, 只是上个线而已. 不会触发callback
// 如果英雄上线导致snapshot中缓存的数据有了变化, 必须再调用一次Update, 把这个变化告知其他系统
func (s *MockHeroSnapshotService) Online(a0 *snapshotdata.HeroSnapshot) {
	fi := getMockFunc(s, s.Online)
	if fi != nil {
		f, ok := fi.(func(*snapshotdata.HeroSnapshot))
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.Online()")
		}
		f(a0)
	}

}

// 注册监听snapshot变化callback
func (s *MockHeroSnapshotService) RegisterCallback(a0 snapshotdata.SnapshotCallback) {
	fi := getMockFunc(s, s.RegisterCallback)
	if fi != nil {
		f, ok := fi.(func(snapshotdata.SnapshotCallback))
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.RegisterCallback()")
		}
		f(a0)
	}

}

// 改变了英雄数据后, 缓存英雄snapshot, 必须是在unlock后调用
// 如果snapshot的版本号低于缓存中的版本号, 不会触发callback
func (s *MockHeroSnapshotService) Update(a0 *snapshotdata.HeroSnapshot) {
	fi := getMockFunc(s, s.Update)
	if fi != nil {
		f, ok := fi.(func(*snapshotdata.HeroSnapshot))
		if !ok {
			panic("invalid mock func, MockHeroSnapshotService.Update()")
		}
		f(a0)
	}

}

var IndividualServerConfig = &MockIndividualServerConfig{}

type MockIndividualServerConfig struct{}

func (s *MockIndividualServerConfig) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockIndividualServerConfig) GetDisablePush() bool {
	fi := getMockFunc(s, s.GetDisablePush)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetDisablePush()")
		}
		return f()
	}

	return false
}
func (s *MockIndividualServerConfig) GetDontEncrypt() bool {
	fi := getMockFunc(s, s.GetDontEncrypt)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetDontEncrypt()")
		}
		return f()
	}

	return false
}
func (s *MockIndividualServerConfig) GetGameAppID() string {
	fi := getMockFunc(s, s.GetGameAppID)
	if fi != nil {
		f, ok := fi.(func() string)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetGameAppID()")
		}
		return f()
	}

	return ""
}
func (s *MockIndividualServerConfig) GetHttpPort() int {
	fi := getMockFunc(s, s.GetHttpPort)
	if fi != nil {
		f, ok := fi.(func() int)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetHttpPort()")
		}
		return f()
	}

	return 0
}
func (s *MockIndividualServerConfig) GetIgnoreHeartBeat() bool {
	fi := getMockFunc(s, s.GetIgnoreHeartBeat)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetIgnoreHeartBeat()")
		}
		return f()
	}

	return false
}
func (s *MockIndividualServerConfig) GetIsAllowRobot() bool {
	fi := getMockFunc(s, s.GetIsAllowRobot)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetIsAllowRobot()")
		}
		return f()
	}

	return false
}
func (s *MockIndividualServerConfig) GetIsDebug() bool {
	fi := getMockFunc(s, s.GetIsDebug)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetIsDebug()")
		}
		return f()
	}

	return false
}
func (s *MockIndividualServerConfig) GetKafkaBrokerAddr() []string {
	fi := getMockFunc(s, s.GetKafkaBrokerAddr)
	if fi != nil {
		f, ok := fi.(func() []string)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetKafkaBrokerAddr()")
		}
		return f()
	}

	return nil
}
func (s *MockIndividualServerConfig) GetKafkaStart() bool {
	fi := getMockFunc(s, s.GetKafkaStart)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetKafkaStart()")
		}
		return f()
	}

	return false
}
func (s *MockIndividualServerConfig) GetLocalAddStr() string {
	fi := getMockFunc(s, s.GetLocalAddStr)
	if fi != nil {
		f, ok := fi.(func() string)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetLocalAddStr()")
		}
		return f()
	}

	return ""
}
func (s *MockIndividualServerConfig) GetPlatformID() int {
	fi := getMockFunc(s, s.GetPlatformID)
	if fi != nil {
		f, ok := fi.(func() int)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetPlatformID()")
		}
		return f()
	}

	return 0
}
func (s *MockIndividualServerConfig) GetPort() int {
	fi := getMockFunc(s, s.GetPort)
	if fi != nil {
		f, ok := fi.(func() int)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetPort()")
		}
		return f()
	}

	return 0
}
func (s *MockIndividualServerConfig) GetReplayPrefix() string {
	fi := getMockFunc(s, s.GetReplayPrefix)
	if fi != nil {
		f, ok := fi.(func() string)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetReplayPrefix()")
		}
		return f()
	}

	return ""
}
func (s *MockIndividualServerConfig) GetServerID() int {
	fi := getMockFunc(s, s.GetServerID)
	if fi != nil {
		f, ok := fi.(func() int)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetServerID()")
		}
		return f()
	}

	return 0
}
func (s *MockIndividualServerConfig) GetServerInfo() *shared_proto.HeroServerInfoProto {
	fi := getMockFunc(s, s.GetServerInfo)
	if fi != nil {
		f, ok := fi.(func() *shared_proto.HeroServerInfoProto)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetServerInfo()")
		}
		return f()
	}

	return nil
}
func (s *MockIndividualServerConfig) GetServerStartTime() time.Time {
	fi := getMockFunc(s, s.GetServerStartTime)
	if fi != nil {
		f, ok := fi.(func() time.Time)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetServerStartTime()")
		}
		return f()
	}

	return time.Time{}
}
func (s *MockIndividualServerConfig) GetSkipHeader() bool {
	fi := getMockFunc(s, s.GetSkipHeader)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetSkipHeader()")
		}
		return f()
	}

	return false
}
func (s *MockIndividualServerConfig) GetTlogStart() bool {
	fi := getMockFunc(s, s.GetTlogStart)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetTlogStart()")
		}
		return f()
	}

	return false
}
func (s *MockIndividualServerConfig) GetTlogTopic() string {
	fi := getMockFunc(s, s.GetTlogTopic)
	if fi != nil {
		f, ok := fi.(func() string)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetTlogTopic()")
		}
		return f()
	}

	return ""
}
func (s *MockIndividualServerConfig) GetZoneAreaID() int {
	fi := getMockFunc(s, s.GetZoneAreaID)
	if fi != nil {
		f, ok := fi.(func() int)
		if !ok {
			panic("invalid mock func, MockIndividualServerConfig.GetZoneAreaID()")
		}
		return f()
	}

	return 0
}

var KafkaService = &MockKafkaService{}

type MockKafkaService struct{}

func (s *MockKafkaService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockKafkaService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockKafkaService.Close()")
		}
		f()
	}

}
func (s *MockKafkaService) NewProducerMsg(a0 entity.Topic, a1 []byte) *sarama.ProducerMessage {
	fi := getMockFunc(s, s.NewProducerMsg)
	if fi != nil {
		f, ok := fi.(func(entity.Topic, []byte) *sarama.ProducerMessage)
		if !ok {
			panic("invalid mock func, MockKafkaService.NewProducerMsg()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockKafkaService) SendAsync(a0 *sarama.ProducerMessage) {
	fi := getMockFunc(s, s.SendAsync)
	if fi != nil {
		f, ok := fi.(func(*sarama.ProducerMessage))
		if !ok {
			panic("invalid mock func, MockKafkaService.SendAsync()")
		}
		f(a0)
	}

}
func (s *MockKafkaService) SendSync(a0 *sarama.ProducerMessage) error {
	fi := getMockFunc(s, s.SendSync)
	if fi != nil {
		f, ok := fi.(func(*sarama.ProducerMessage) error)
		if !ok {
			panic("invalid mock func, MockKafkaService.SendSync()")
		}
		return f(a0)
	}

	return nil
}

var LightpawHandler = &MockLightpawHandler{}

type MockLightpawHandler struct{}

func (s *MockLightpawHandler) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var LocationHeroCache = &MockLocationHeroCache{}

type MockLocationHeroCache struct{}

func (s *MockLocationHeroCache) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockLocationHeroCache) UpdateHero(a0 *snapshotdata.HeroSnapshot, a1 int64) {
	fi := getMockFunc(s, s.UpdateHero)
	if fi != nil {
		f, ok := fi.(func(*snapshotdata.HeroSnapshot, int64))
		if !ok {
			panic("invalid mock func, MockLocationHeroCache.UpdateHero()")
		}
		f(a0, a1)
	}

}
func (s *MockLocationHeroCache) UpdateLocation(a0 int64, a1 uint64) {
	fi := getMockFunc(s, s.UpdateLocation)
	if fi != nil {
		f, ok := fi.(func(int64, uint64))
		if !ok {
			panic("invalid mock func, MockLocationHeroCache.UpdateLocation()")
		}
		f(a0, a1)
	}

}

var MailModule = &MockMailModule{}

type MockMailModule struct{}

func (s *MockMailModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockMailModule) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockMailModule.OnHeroOnline()")
		}
		f(a0)
	}

}

// 过期函数@AlbertFan
// 发邮件
func (s *MockMailModule) SendMail(a0 int64, a1 uint64, a2 string, a3 string, a4 bool, a5 *shared_proto.FightReportProto, a6 *shared_proto.PrizeProto, a7 time.Time) bool {
	fi := getMockFunc(s, s.SendMail)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, string, string, bool, *shared_proto.FightReportProto, *shared_proto.PrizeProto, time.Time) bool)
		if !ok {
			panic("invalid mock func, MockMailModule.SendMail()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7)
	}

	return false
}
func (s *MockMailModule) SendProtoMail(a0 int64, a1 *shared_proto.MailProto, a2 time.Time) bool {
	fi := getMockFunc(s, s.SendProtoMail)
	if fi != nil {
		f, ok := fi.(func(int64, *shared_proto.MailProto, time.Time) bool)
		if !ok {
			panic("invalid mock func, MockMailModule.SendProtoMail()")
		}
		return f(a0, a1, a2)
	}

	return false
}
func (s *MockMailModule) SendReportMail(a0 int64, a1 *shared_proto.MailProto, a2 time.Time) bool {
	fi := getMockFunc(s, s.SendReportMail)
	if fi != nil {
		f, ok := fi.(func(int64, *shared_proto.MailProto, time.Time) bool)
		if !ok {
			panic("invalid mock func, MockMailModule.SendReportMail()")
		}
		return f(a0, a1, a2)
	}

	return false
}

var MetricsRegister = &MockMetricsRegister{}

type MockMetricsRegister struct{}

func (s *MockMetricsRegister) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockMetricsRegister) EnableDBMetrics() bool {
	fi := getMockFunc(s, s.EnableDBMetrics)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockMetricsRegister.EnableDBMetrics()")
		}
		return f()
	}

	return false
}
func (s *MockMetricsRegister) EnableMsgMetrics() bool {
	fi := getMockFunc(s, s.EnableMsgMetrics)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockMetricsRegister.EnableMsgMetrics()")
		}
		return f()
	}

	return false
}
func (s *MockMetricsRegister) EnableOnlineCountMetrics() bool {
	fi := getMockFunc(s, s.EnableOnlineCountMetrics)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockMetricsRegister.EnableOnlineCountMetrics()")
		}
		return f()
	}

	return false
}
func (s *MockMetricsRegister) EnablePanicMetrics() bool {
	fi := getMockFunc(s, s.EnablePanicMetrics)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockMetricsRegister.EnablePanicMetrics()")
		}
		return f()
	}

	return false
}
func (s *MockMetricsRegister) EnableRegisterCountMetrics() bool {
	fi := getMockFunc(s, s.EnableRegisterCountMetrics)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockMetricsRegister.EnableRegisterCountMetrics()")
		}
		return f()
	}

	return false
}
func (s *MockMetricsRegister) Register(a0 prometheus.Collector) {
	fi := getMockFunc(s, s.Register)
	if fi != nil {
		f, ok := fi.(func(prometheus.Collector))
		if !ok {
			panic("invalid mock func, MockMetricsRegister.Register()")
		}
		f(a0)
	}

}
func (s *MockMetricsRegister) RegisterFunc(a0 metrics.CollectFunc) {
	fi := getMockFunc(s, s.RegisterFunc)
	if fi != nil {
		f, ok := fi.(func(metrics.CollectFunc))
		if !ok {
			panic("invalid mock func, MockMetricsRegister.RegisterFunc()")
		}
		f(a0)
	}

}

var MilitaryModule = &MockMilitaryModule{}

type MockMilitaryModule struct{}

func (s *MockMilitaryModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockMilitaryModule) GmRate(a0 int64, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.GmRate)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockMilitaryModule.GmRate()")
		}
		f(a0, a1, a2)
	}

}

var MingcModule = &MockMingcModule{}

type MockMingcModule struct{}

func (s *MockMingcModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var MingcService = &MockMingcService{}

type MockMingcService struct{}

func (s *MockMingcService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockMingcService) AllInOneGuild() int64 {
	fi := getMockFunc(s, s.AllInOneGuild)
	if fi != nil {
		f, ok := fi.(func() int64)
		if !ok {
			panic("invalid mock func, MockMingcService.AllInOneGuild()")
		}
		return f()
	}

	return 0
}
func (s *MockMingcService) Build(a0 int64, a1 int64, a2 uint64, a3 uint64) bool {
	fi := getMockFunc(s, s.Build)
	if fi != nil {
		f, ok := fi.(func(int64, int64, uint64, uint64) bool)
		if !ok {
			panic("invalid mock func, MockMingcService.Build()")
		}
		return f(a0, a1, a2, a3)
	}

	return false
}
func (s *MockMingcService) CaptainHostGuild(a0 uint64) int64 {
	fi := getMockFunc(s, s.CaptainHostGuild)
	if fi != nil {
		f, ok := fi.(func(uint64) int64)
		if !ok {
			panic("invalid mock func, MockMingcService.CaptainHostGuild()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockMingcService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockMingcService.Close()")
		}
		f()
	}

}
func (s *MockMingcService) Country(a0 uint64) uint64 {
	fi := getMockFunc(s, s.Country)
	if fi != nil {
		f, ok := fi.(func(uint64) uint64)
		if !ok {
			panic("invalid mock func, MockMingcService.Country()")
		}
		return f(a0)
	}

	return 0
}

// 国家当前占领的本国初始名城
func (s *MockMingcService) CountryHoldInitMcs(a0 uint64) []*entity.Mingc {
	fi := getMockFunc(s, s.CountryHoldInitMcs)
	if fi != nil {
		f, ok := fi.(func(uint64) []*entity.Mingc)
		if !ok {
			panic("invalid mock func, MockMingcService.CountryHoldInitMcs()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockMingcService) DisableMcBuildLogCache(a0 uint64) {
	fi := getMockFunc(s, s.DisableMcBuildLogCache)
	if fi != nil {
		f, ok := fi.(func(uint64))
		if !ok {
			panic("invalid mock func, MockMingcService.DisableMcBuildLogCache()")
		}
		f(a0)
	}

}
func (s *MockMingcService) GetMcBuildGuildMemberPrize(a0 uint64) *resdata.Prize {
	fi := getMockFunc(s, s.GetMcBuildGuildMemberPrize)
	if fi != nil {
		f, ok := fi.(func(uint64) *resdata.Prize)
		if !ok {
			panic("invalid mock func, MockMingcService.GetMcBuildGuildMemberPrize()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockMingcService) GuildMingc(a0 int64) *entity.Mingc {
	fi := getMockFunc(s, s.GuildMingc)
	if fi != nil {
		f, ok := fi.(func(int64) *entity.Mingc)
		if !ok {
			panic("invalid mock func, MockMingcService.GuildMingc()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockMingcService) IsHoldCountryCapital(a0 int64, a1 uint64) bool {
	fi := getMockFunc(s, s.IsHoldCountryCapital)
	if fi != nil {
		f, ok := fi.(func(int64, uint64) bool)
		if !ok {
			panic("invalid mock func, MockMingcService.IsHoldCountryCapital()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockMingcService) McBuildLogMsg(a0 *entity.Mingc) pbutil.Buffer {
	fi := getMockFunc(s, s.McBuildLogMsg)
	if fi != nil {
		f, ok := fi.(func(*entity.Mingc) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockMingcService.McBuildLogMsg()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockMingcService) Mingc(a0 uint64) *entity.Mingc {
	fi := getMockFunc(s, s.Mingc)
	if fi != nil {
		f, ok := fi.(func(uint64) *entity.Mingc)
		if !ok {
			panic("invalid mock func, MockMingcService.Mingc()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockMingcService) MingcsMsg(a0 uint64) pbutil.Buffer {
	fi := getMockFunc(s, s.MingcsMsg)
	if fi != nil {
		f, ok := fi.(func(uint64) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockMingcService.MingcsMsg()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockMingcService) SetHostGuild(a0 uint64, a1 int64) bool {
	fi := getMockFunc(s, s.SetHostGuild)
	if fi != nil {
		f, ok := fi.(func(uint64, int64) bool)
		if !ok {
			panic("invalid mock func, MockMingcService.SetHostGuild()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockMingcService) UpdateMsg() pbutil.Buffer {
	fi := getMockFunc(s, s.UpdateMsg)
	if fi != nil {
		f, ok := fi.(func() pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockMingcService.UpdateMsg()")
		}
		return f()
	}

	return nil
}
func (s *MockMingcService) WalkMingcs(a0 entity.MingcFunc) {
	fi := getMockFunc(s, s.WalkMingcs)
	if fi != nil {
		f, ok := fi.(func(entity.MingcFunc))
		if !ok {
			panic("invalid mock func, MockMingcService.WalkMingcs()")
		}
		f(a0)
	}

}

var MingcWarModule = &MockMingcWarModule{}

type MockMingcWarModule struct{}

func (s *MockMingcWarModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var MingcWarService = &MockMingcWarService{}

type MockMingcWarService struct{}

func (s *MockMingcWarService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockMingcWarService) ApplyAst(a0 int64, a1 bool, a2 *mingcdata.MingcBaseData) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.ApplyAst)
	if fi != nil {
		f, ok := fi.(func(int64, bool, *mingcdata.MingcBaseData) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.ApplyAst()")
		}
		return f(a0, a1, a2)
	}

	return nil, nil
}
func (s *MockMingcWarService) ApplyAstNotice(a0 int64, a1 int64) bool {
	fi := getMockFunc(s, s.ApplyAstNotice)
	if fi != nil {
		f, ok := fi.(func(int64, int64) bool)
		if !ok {
			panic("invalid mock func, MockMingcWarService.ApplyAstNotice()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockMingcWarService) ApplyAtk(a0 int64, a1 *mingcdata.MingcBaseData, a2 uint64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.ApplyAtk)
	if fi != nil {
		f, ok := fi.(func(int64, *mingcdata.MingcBaseData, uint64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.ApplyAtk()")
		}
		return f(a0, a1, a2)
	}

	return nil, nil
}
func (s *MockMingcWarService) ApplyAtkNotice(a0 int64, a1 int64) bool {
	fi := getMockFunc(s, s.ApplyAtkNotice)
	if fi != nil {
		f, ok := fi.(func(int64, int64) bool)
		if !ok {
			panic("invalid mock func, MockMingcWarService.ApplyAtkNotice()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockMingcWarService) BuildFightStartMsg(a0 time.Time) (pbutil.Buffer, bool) {
	fi := getMockFunc(s, s.BuildFightStartMsg)
	if fi != nil {
		f, ok := fi.(func(time.Time) (pbutil.Buffer, bool))
		if !ok {
			panic("invalid mock func, MockMingcWarService.BuildFightStartMsg()")
		}
		return f(a0)
	}

	return nil, false
}
func (s *MockMingcWarService) CancelApplyAst(a0 int64, a1 *mingcdata.MingcBaseData) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.CancelApplyAst)
	if fi != nil {
		f, ok := fi.(func(int64, *mingcdata.MingcBaseData) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.CancelApplyAst()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockMingcWarService) CatchHistoryRecord(a0 int64, a1 int64) pbutil.Buffer {
	fi := getMockFunc(s, s.CatchHistoryRecord)
	if fi != nil {
		f, ok := fi.(func(int64, int64) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockMingcWarService.CatchHistoryRecord()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockMingcWarService) CatchTroopsRank(a0 int64, a1 uint64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.CatchTroopsRank)
	if fi != nil {
		f, ok := fi.(func(int64, uint64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.CatchTroopsRank()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockMingcWarService) CleanOnGuildRemoved(a0 int64) {
	fi := getMockFunc(s, s.CleanOnGuildRemoved)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockMingcWarService.CleanOnGuildRemoved()")
		}
		f(a0)
	}

}
func (s *MockMingcWarService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockMingcWarService.Close()")
		}
		f()
	}

}
func (s *MockMingcWarService) CurrMcWarStage() (int32, time.Time, time.Time) {
	fi := getMockFunc(s, s.CurrMcWarStage)
	if fi != nil {
		f, ok := fi.(func() (int32, time.Time, time.Time))
		if !ok {
			panic("invalid mock func, MockMingcWarService.CurrMcWarStage()")
		}
		return f()
	}

	return 0, time.Time{}, time.Time{}
}
func (s *MockMingcWarService) GmApplyAtkGuild(a0 uint64, a1 int64) {
	fi := getMockFunc(s, s.GmApplyAtkGuild)
	if fi != nil {
		f, ok := fi.(func(uint64, int64))
		if !ok {
			panic("invalid mock func, MockMingcWarService.GmApplyAtkGuild()")
		}
		f(a0, a1)
	}

}
func (s *MockMingcWarService) GmCampFail(a0 uint64, a1 bool) {
	fi := getMockFunc(s, s.GmCampFail)
	if fi != nil {
		f, ok := fi.(func(uint64, bool))
		if !ok {
			panic("invalid mock func, MockMingcWarService.GmCampFail()")
		}
		f(a0, a1)
	}

}
func (s *MockMingcWarService) GmChangeStage(a0 shared_proto.MingcWarState, a1 iface.HeroController, a2 time.Time) {
	fi := getMockFunc(s, s.GmChangeStage)
	if fi != nil {
		f, ok := fi.(func(shared_proto.MingcWarState, iface.HeroController, time.Time))
		if !ok {
			panic("invalid mock func, MockMingcWarService.GmChangeStage()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockMingcWarService) GmNewMingcWar() {
	fi := getMockFunc(s, s.GmNewMingcWar)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockMingcWarService.GmNewMingcWar()")
		}
		f()
	}

}
func (s *MockMingcWarService) GmSetAstAtkGuild(a0 uint64, a1 int64) {
	fi := getMockFunc(s, s.GmSetAstAtkGuild)
	if fi != nil {
		f, ok := fi.(func(uint64, int64))
		if !ok {
			panic("invalid mock func, MockMingcWarService.GmSetAstAtkGuild()")
		}
		f(a0, a1)
	}

}
func (s *MockMingcWarService) GmSetAstDefGuild(a0 uint64, a1 int64) {
	fi := getMockFunc(s, s.GmSetAstDefGuild)
	if fi != nil {
		f, ok := fi.(func(uint64, int64))
		if !ok {
			panic("invalid mock func, MockMingcWarService.GmSetAstDefGuild()")
		}
		f(a0, a1)
	}

}
func (s *MockMingcWarService) GmSetDefGuild(a0 uint64, a1 int64) {
	fi := getMockFunc(s, s.GmSetDefGuild)
	if fi != nil {
		f, ok := fi.(func(uint64, int64))
		if !ok {
			panic("invalid mock func, MockMingcWarService.GmSetDefGuild()")
		}
		f(a0, a1)
	}

}
func (s *MockMingcWarService) GuildMcWarType(a0 int64) (uint64, shared_proto.MingcWarGuildType) {
	fi := getMockFunc(s, s.GuildMcWarType)
	if fi != nil {
		f, ok := fi.(func(int64) (uint64, shared_proto.MingcWarGuildType))
		if !ok {
			panic("invalid mock func, MockMingcWarService.GuildMcWarType()")
		}
		return f(a0)
	}

	return 0, 0
}
func (s *MockMingcWarService) JoinFight(a0 iface.HeroController, a1 uint64, a2 []uint64, a3 []int32) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.JoinFight)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, uint64, []uint64, []int32) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.JoinFight()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil, nil
}
func (s *MockMingcWarService) JoiningFightMingc(a0 int64) (uint64, bool) {
	fi := getMockFunc(s, s.JoiningFightMingc)
	if fi != nil {
		f, ok := fi.(func(int64) (uint64, bool))
		if !ok {
			panic("invalid mock func, MockMingcWarService.JoiningFightMingc()")
		}
		return f(a0)
	}

	return 0, false
}
func (s *MockMingcWarService) McWarStartEndTime() (time.Time, time.Time) {
	fi := getMockFunc(s, s.McWarStartEndTime)
	if fi != nil {
		f, ok := fi.(func() (time.Time, time.Time))
		if !ok {
			panic("invalid mock func, MockMingcWarService.McWarStartEndTime()")
		}
		return f()
	}

	return time.Time{}, time.Time{}
}
func (s *MockMingcWarService) QuitFight(a0 int64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.QuitFight)
	if fi != nil {
		f, ok := fi.(func(int64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.QuitFight()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockMingcWarService) QuitWatch(a0 iface.HeroController, a1 uint64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.QuitWatch)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, uint64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.QuitWatch()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockMingcWarService) ReplyApplyAst(a0 int64, a1 int64, a2 *mingcdata.MingcBaseData, a3 bool) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.ReplyApplyAst)
	if fi != nil {
		f, ok := fi.(func(int64, int64, *mingcdata.MingcBaseData, bool) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.ReplyApplyAst()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil, nil
}
func (s *MockMingcWarService) SceneBack(a0 int64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.SceneBack)
	if fi != nil {
		f, ok := fi.(func(int64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.SceneBack()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockMingcWarService) SceneChangeMode(a0 int64, a1 shared_proto.MingcWarModeType) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.SceneChangeMode)
	if fi != nil {
		f, ok := fi.(func(int64, shared_proto.MingcWarModeType) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.SceneChangeMode()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockMingcWarService) SceneDrum(a0 int64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.SceneDrum)
	if fi != nil {
		f, ok := fi.(func(int64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.SceneDrum()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockMingcWarService) SceneMove(a0 int64, a1 cb.Cube) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.SceneMove)
	if fi != nil {
		f, ok := fi.(func(int64, cb.Cube) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.SceneMove()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockMingcWarService) SceneSpeedUp(a0 int64, a1 float64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.SceneSpeedUp)
	if fi != nil {
		f, ok := fi.(func(int64, float64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.SceneSpeedUp()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockMingcWarService) SceneTouShiBuildingFire(a0 int64, a1 cb.Cube) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.SceneTouShiBuildingFire)
	if fi != nil {
		f, ok := fi.(func(int64, cb.Cube) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.SceneTouShiBuildingFire()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockMingcWarService) SceneTouShiBuildingTurnTo(a0 int64, a1 cb.Cube, a2 bool) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.SceneTouShiBuildingTurnTo)
	if fi != nil {
		f, ok := fi.(func(int64, cb.Cube, bool) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.SceneTouShiBuildingTurnTo()")
		}
		return f(a0, a1, a2)
	}

	return nil, nil
}
func (s *MockMingcWarService) SceneTroopRelive(a0 int64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.SceneTroopRelive)
	if fi != nil {
		f, ok := fi.(func(int64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.SceneTroopRelive()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockMingcWarService) SendChat(a0 int64, a1 *shared_proto.ChatMsgProto) pbutil.Buffer {
	fi := getMockFunc(s, s.SendChat)
	if fi != nil {
		f, ok := fi.(func(int64, *shared_proto.ChatMsgProto) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockMingcWarService.SendChat()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockMingcWarService) UpdateMsg() {
	fi := getMockFunc(s, s.UpdateMsg)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockMingcWarService.UpdateMsg()")
		}
		f()
	}

}
func (s *MockMingcWarService) ViewMcWarMcMsg(a0 uint64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.ViewMcWarMcMsg)
	if fi != nil {
		f, ok := fi.(func(uint64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.ViewMcWarMcMsg()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockMingcWarService) ViewMcWarSceneMsg(a0 uint64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.ViewMcWarSceneMsg)
	if fi != nil {
		f, ok := fi.(func(uint64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.ViewMcWarSceneMsg()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockMingcWarService) ViewMsg(a0 uint64) pbutil.Buffer {
	fi := getMockFunc(s, s.ViewMsg)
	if fi != nil {
		f, ok := fi.(func(uint64) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockMingcWarService.ViewMsg()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockMingcWarService) ViewSceneTroopRecord(a0 int64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.ViewSceneTroopRecord)
	if fi != nil {
		f, ok := fi.(func(int64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.ViewSceneTroopRecord()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockMingcWarService) ViewSelfGuildProto(a0 int64) *shared_proto.McWarGuildProto {
	fi := getMockFunc(s, s.ViewSelfGuildProto)
	if fi != nil {
		f, ok := fi.(func(int64) *shared_proto.McWarGuildProto)
		if !ok {
			panic("invalid mock func, MockMingcWarService.ViewSelfGuildProto()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockMingcWarService) Watch(a0 iface.HeroController, a1 uint64) (pbutil.Buffer, pbutil.Buffer) {
	fi := getMockFunc(s, s.Watch)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, uint64) (pbutil.Buffer, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockMingcWarService.Watch()")
		}
		return f(a0, a1)
	}

	return nil, nil
}

var MiscModule = &MockMiscModule{}

type MockMiscModule struct{}

func (s *MockMiscModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockMiscModule) GmSetClientVersion(a0 string) {
	fi := getMockFunc(s, s.GmSetClientVersion)
	if fi != nil {
		f, ok := fi.(func(string))
		if !ok {
			panic("invalid mock func, MockMiscModule.GmSetClientVersion()")
		}
		f(a0)
	}

}
func (s *MockMiscModule) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockMiscModule.OnHeroOnline()")
		}
		f(a0)
	}

}
func (s *MockMiscModule) PrintClientLog(a0 int64, a1 string, a2 string, a3 string) {
	fi := getMockFunc(s, s.PrintClientLog)
	if fi != nil {
		f, ok := fi.(func(int64, string, string, string))
		if !ok {
			panic("invalid mock func, MockMiscModule.PrintClientLog()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockMiscModule) SendClientVersion(a0 sender.Sender) {
	fi := getMockFunc(s, s.SendClientVersion)
	if fi != nil {
		f, ok := fi.(func(sender.Sender))
		if !ok {
			panic("invalid mock func, MockMiscModule.SendClientVersion()")
		}
		f(a0)
	}

}
func (s *MockMiscModule) SendConfig(a0 *misc.C2SConfigProto, a1 sender.Sender) {
	fi := getMockFunc(s, s.SendConfig)
	if fi != nil {
		f, ok := fi.(func(*misc.C2SConfigProto, sender.Sender))
		if !ok {
			panic("invalid mock func, MockMiscModule.SendConfig()")
		}
		f(a0, a1)
	}

}
func (s *MockMiscModule) SendLuaConfig(a0 *misc.C2SConfigluaProto, a1 sender.Sender) {
	fi := getMockFunc(s, s.SendLuaConfig)
	if fi != nil {
		f, ok := fi.(func(*misc.C2SConfigluaProto, sender.Sender))
		if !ok {
			panic("invalid mock func, MockMiscModule.SendLuaConfig()")
		}
		f(a0, a1)
	}

}
func (s *MockMiscModule) SyncTime(a0 int32, a1 sender.Sender) {
	fi := getMockFunc(s, s.SyncTime)
	if fi != nil {
		f, ok := fi.(func(int32, sender.Sender))
		if !ok {
			panic("invalid mock func, MockMiscModule.SyncTime()")
		}
		f(a0, a1)
	}

}

var Modules = &MockModules{}

type MockModules struct{}

func (s *MockModules) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockModules) ActivityModule() iface.ActivityModule {
	fi := getMockFunc(s, s.ActivityModule)
	if fi != nil {
		f, ok := fi.(func() iface.ActivityModule)
		if !ok {
			panic("invalid mock func, MockModules.ActivityModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) BaiZhanModule() iface.BaiZhanModule {
	fi := getMockFunc(s, s.BaiZhanModule)
	if fi != nil {
		f, ok := fi.(func() iface.BaiZhanModule)
		if !ok {
			panic("invalid mock func, MockModules.BaiZhanModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) ChatModule() iface.ChatModule {
	fi := getMockFunc(s, s.ChatModule)
	if fi != nil {
		f, ok := fi.(func() iface.ChatModule)
		if !ok {
			panic("invalid mock func, MockModules.ChatModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) ClientConfigModule() iface.ClientConfigModule {
	fi := getMockFunc(s, s.ClientConfigModule)
	if fi != nil {
		f, ok := fi.(func() iface.ClientConfigModule)
		if !ok {
			panic("invalid mock func, MockModules.ClientConfigModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) CountryModule() iface.CountryModule {
	fi := getMockFunc(s, s.CountryModule)
	if fi != nil {
		f, ok := fi.(func() iface.CountryModule)
		if !ok {
			panic("invalid mock func, MockModules.CountryModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) DepotModule() iface.DepotModule {
	fi := getMockFunc(s, s.DepotModule)
	if fi != nil {
		f, ok := fi.(func() iface.DepotModule)
		if !ok {
			panic("invalid mock func, MockModules.DepotModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) DianquanModule() iface.DianquanModule {
	fi := getMockFunc(s, s.DianquanModule)
	if fi != nil {
		f, ok := fi.(func() iface.DianquanModule)
		if !ok {
			panic("invalid mock func, MockModules.DianquanModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) DomesticModule() iface.DomesticModule {
	fi := getMockFunc(s, s.DomesticModule)
	if fi != nil {
		f, ok := fi.(func() iface.DomesticModule)
		if !ok {
			panic("invalid mock func, MockModules.DomesticModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) DungeonModule() iface.DungeonModule {
	fi := getMockFunc(s, s.DungeonModule)
	if fi != nil {
		f, ok := fi.(func() iface.DungeonModule)
		if !ok {
			panic("invalid mock func, MockModules.DungeonModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) EquipmentModule() iface.EquipmentModule {
	fi := getMockFunc(s, s.EquipmentModule)
	if fi != nil {
		f, ok := fi.(func() iface.EquipmentModule)
		if !ok {
			panic("invalid mock func, MockModules.EquipmentModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) FarmModule() iface.FarmModule {
	fi := getMockFunc(s, s.FarmModule)
	if fi != nil {
		f, ok := fi.(func() iface.FarmModule)
		if !ok {
			panic("invalid mock func, MockModules.FarmModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) FishingModule() iface.FishingModule {
	fi := getMockFunc(s, s.FishingModule)
	if fi != nil {
		f, ok := fi.(func() iface.FishingModule)
		if !ok {
			panic("invalid mock func, MockModules.FishingModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) GardenModule() iface.GardenModule {
	fi := getMockFunc(s, s.GardenModule)
	if fi != nil {
		f, ok := fi.(func() iface.GardenModule)
		if !ok {
			panic("invalid mock func, MockModules.GardenModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) GemModule() iface.GemModule {
	fi := getMockFunc(s, s.GemModule)
	if fi != nil {
		f, ok := fi.(func() iface.GemModule)
		if !ok {
			panic("invalid mock func, MockModules.GemModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) GuildModule() iface.GuildModule {
	fi := getMockFunc(s, s.GuildModule)
	if fi != nil {
		f, ok := fi.(func() iface.GuildModule)
		if !ok {
			panic("invalid mock func, MockModules.GuildModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) HebiModule() iface.HebiModule {
	fi := getMockFunc(s, s.HebiModule)
	if fi != nil {
		f, ok := fi.(func() iface.HebiModule)
		if !ok {
			panic("invalid mock func, MockModules.HebiModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) MailModule() iface.MailModule {
	fi := getMockFunc(s, s.MailModule)
	if fi != nil {
		f, ok := fi.(func() iface.MailModule)
		if !ok {
			panic("invalid mock func, MockModules.MailModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) MilitaryModule() iface.MilitaryModule {
	fi := getMockFunc(s, s.MilitaryModule)
	if fi != nil {
		f, ok := fi.(func() iface.MilitaryModule)
		if !ok {
			panic("invalid mock func, MockModules.MilitaryModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) MingcModule() iface.MingcModule {
	fi := getMockFunc(s, s.MingcModule)
	if fi != nil {
		f, ok := fi.(func() iface.MingcModule)
		if !ok {
			panic("invalid mock func, MockModules.MingcModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) MingcWarModule() iface.MingcWarModule {
	fi := getMockFunc(s, s.MingcWarModule)
	if fi != nil {
		f, ok := fi.(func() iface.MingcWarModule)
		if !ok {
			panic("invalid mock func, MockModules.MingcWarModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) MiscModule() iface.MiscModule {
	fi := getMockFunc(s, s.MiscModule)
	if fi != nil {
		f, ok := fi.(func() iface.MiscModule)
		if !ok {
			panic("invalid mock func, MockModules.MiscModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) PromotionModule() iface.PromotionModule {
	fi := getMockFunc(s, s.PromotionModule)
	if fi != nil {
		f, ok := fi.(func() iface.PromotionModule)
		if !ok {
			panic("invalid mock func, MockModules.PromotionModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) QuestionModule() iface.QuestionModule {
	fi := getMockFunc(s, s.QuestionModule)
	if fi != nil {
		f, ok := fi.(func() iface.QuestionModule)
		if !ok {
			panic("invalid mock func, MockModules.QuestionModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) RandomEventModule() iface.RandomEventModule {
	fi := getMockFunc(s, s.RandomEventModule)
	if fi != nil {
		f, ok := fi.(func() iface.RandomEventModule)
		if !ok {
			panic("invalid mock func, MockModules.RandomEventModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) RankModule() iface.RankModule {
	fi := getMockFunc(s, s.RankModule)
	if fi != nil {
		f, ok := fi.(func() iface.RankModule)
		if !ok {
			panic("invalid mock func, MockModules.RankModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) RedPacketModule() iface.RedPacketModule {
	fi := getMockFunc(s, s.RedPacketModule)
	if fi != nil {
		f, ok := fi.(func() iface.RedPacketModule)
		if !ok {
			panic("invalid mock func, MockModules.RedPacketModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) RegionModule() iface.RegionModule {
	fi := getMockFunc(s, s.RegionModule)
	if fi != nil {
		f, ok := fi.(func() iface.RegionModule)
		if !ok {
			panic("invalid mock func, MockModules.RegionModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) RelationModule() iface.RelationModule {
	fi := getMockFunc(s, s.RelationModule)
	if fi != nil {
		f, ok := fi.(func() iface.RelationModule)
		if !ok {
			panic("invalid mock func, MockModules.RelationModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) SecretTowerModule() iface.SecretTowerModule {
	fi := getMockFunc(s, s.SecretTowerModule)
	if fi != nil {
		f, ok := fi.(func() iface.SecretTowerModule)
		if !ok {
			panic("invalid mock func, MockModules.SecretTowerModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) ShopModule() iface.ShopModule {
	fi := getMockFunc(s, s.ShopModule)
	if fi != nil {
		f, ok := fi.(func() iface.ShopModule)
		if !ok {
			panic("invalid mock func, MockModules.ShopModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) StrategyModule() iface.StrategyModule {
	fi := getMockFunc(s, s.StrategyModule)
	if fi != nil {
		f, ok := fi.(func() iface.StrategyModule)
		if !ok {
			panic("invalid mock func, MockModules.StrategyModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) StressModule() iface.StressModule {
	fi := getMockFunc(s, s.StressModule)
	if fi != nil {
		f, ok := fi.(func() iface.StressModule)
		if !ok {
			panic("invalid mock func, MockModules.StressModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) SurveyModule() iface.SurveyModule {
	fi := getMockFunc(s, s.SurveyModule)
	if fi != nil {
		f, ok := fi.(func() iface.SurveyModule)
		if !ok {
			panic("invalid mock func, MockModules.SurveyModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) TagModule() iface.TagModule {
	fi := getMockFunc(s, s.TagModule)
	if fi != nil {
		f, ok := fi.(func() iface.TagModule)
		if !ok {
			panic("invalid mock func, MockModules.TagModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) TaskModule() iface.TaskModule {
	fi := getMockFunc(s, s.TaskModule)
	if fi != nil {
		f, ok := fi.(func() iface.TaskModule)
		if !ok {
			panic("invalid mock func, MockModules.TaskModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) TeachModule() iface.TeachModule {
	fi := getMockFunc(s, s.TeachModule)
	if fi != nil {
		f, ok := fi.(func() iface.TeachModule)
		if !ok {
			panic("invalid mock func, MockModules.TeachModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) TowerModule() iface.TowerModule {
	fi := getMockFunc(s, s.TowerModule)
	if fi != nil {
		f, ok := fi.(func() iface.TowerModule)
		if !ok {
			panic("invalid mock func, MockModules.TowerModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) VipModule() iface.VipModule {
	fi := getMockFunc(s, s.VipModule)
	if fi != nil {
		f, ok := fi.(func() iface.VipModule)
		if !ok {
			panic("invalid mock func, MockModules.VipModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) XiongNuModule() iface.XiongNuModule {
	fi := getMockFunc(s, s.XiongNuModule)
	if fi != nil {
		f, ok := fi.(func() iface.XiongNuModule)
		if !ok {
			panic("invalid mock func, MockModules.XiongNuModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) XuanyuanModule() iface.XuanyuanModule {
	fi := getMockFunc(s, s.XuanyuanModule)
	if fi != nil {
		f, ok := fi.(func() iface.XuanyuanModule)
		if !ok {
			panic("invalid mock func, MockModules.XuanyuanModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) ZhanJiangModule() iface.ZhanJiangModule {
	fi := getMockFunc(s, s.ZhanJiangModule)
	if fi != nil {
		f, ok := fi.(func() iface.ZhanJiangModule)
		if !ok {
			panic("invalid mock func, MockModules.ZhanJiangModule()")
		}
		return f()
	}

	return nil
}
func (s *MockModules) ZhengWuModule() iface.ZhengWuModule {
	fi := getMockFunc(s, s.ZhengWuModule)
	if fi != nil {
		f, ok := fi.(func() iface.ZhengWuModule)
		if !ok {
			panic("invalid mock func, MockModules.ZhengWuModule()")
		}
		return f()
	}

	return nil
}

var ProductService = &MockProductService{}

type MockProductService struct{}

func (s *MockProductService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var PromotionModule = &MockPromotionModule{}

type MockPromotionModule struct{}

func (s *MockPromotionModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var PushService = &MockPushService{}

type MockPushService struct{}

func (s *MockPushService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockPushService) GmPush(a0 int64) {
	fi := getMockFunc(s, s.GmPush)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockPushService.GmPush()")
		}
		f(a0)
	}

}
func (s *MockPushService) MultiPush(a0 shared_proto.SettingType, a1 []int64, a2 int64) {
	fi := getMockFunc(s, s.MultiPush)
	if fi != nil {
		f, ok := fi.(func(shared_proto.SettingType, []int64, int64))
		if !ok {
			panic("invalid mock func, MockPushService.MultiPush()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockPushService) MultiPushFunc(a0 shared_proto.SettingType, a1 []int64, a2 int64, a3 pushdata.PushFunc) {
	fi := getMockFunc(s, s.MultiPushFunc)
	if fi != nil {
		f, ok := fi.(func(shared_proto.SettingType, []int64, int64, pushdata.PushFunc))
		if !ok {
			panic("invalid mock func, MockPushService.MultiPushFunc()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockPushService) MultiPushTitleContent(a0 shared_proto.SettingType, a1 string, a2 string, a3 []int64, a4 int64) {
	fi := getMockFunc(s, s.MultiPushTitleContent)
	if fi != nil {
		f, ok := fi.(func(shared_proto.SettingType, string, string, []int64, int64))
		if !ok {
			panic("invalid mock func, MockPushService.MultiPushTitleContent()")
		}
		f(a0, a1, a2, a3, a4)
	}

}
func (s *MockPushService) Push(a0 shared_proto.SettingType, a1 int64) {
	fi := getMockFunc(s, s.Push)
	if fi != nil {
		f, ok := fi.(func(shared_proto.SettingType, int64))
		if !ok {
			panic("invalid mock func, MockPushService.Push()")
		}
		f(a0, a1)
	}

}
func (s *MockPushService) PushFunc(a0 shared_proto.SettingType, a1 int64, a2 pushdata.PushFunc) {
	fi := getMockFunc(s, s.PushFunc)
	if fi != nil {
		f, ok := fi.(func(shared_proto.SettingType, int64, pushdata.PushFunc))
		if !ok {
			panic("invalid mock func, MockPushService.PushFunc()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockPushService) PushTitleContent(a0 shared_proto.SettingType, a1 string, a2 string, a3 int64) {
	fi := getMockFunc(s, s.PushTitleContent)
	if fi != nil {
		f, ok := fi.(func(shared_proto.SettingType, string, string, int64))
		if !ok {
			panic("invalid mock func, MockPushService.PushTitleContent()")
		}
		f(a0, a1, a2, a3)
	}

}

var QuestionModule = &MockQuestionModule{}

type MockQuestionModule struct{}

func (s *MockQuestionModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var RandomEventModule = &MockRandomEventModule{}

type MockRandomEventModule struct{}

func (s *MockRandomEventModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var RankModule = &MockRankModule{}

type MockRankModule struct{}

func (s *MockRankModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockRankModule) AddOrUpdateRankObj(a0 rankface.RankObj) {
	fi := getMockFunc(s, s.AddOrUpdateRankObj)
	if fi != nil {
		f, ok := fi.(func(rankface.RankObj))
		if !ok {
			panic("invalid mock func, MockRankModule.AddOrUpdateRankObj()")
		}
		f(a0)
	}

}
func (s *MockRankModule) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockRankModule.Close()")
		}
		f()
	}

}
func (s *MockRankModule) CountryOfficial(a0 int, a1 string, a2 uint64, a3 shared_proto.CountryOfficialType) []*shared_proto.HeroBasicSnapshotProto {
	fi := getMockFunc(s, s.CountryOfficial)
	if fi != nil {
		f, ok := fi.(func(int, string, uint64, shared_proto.CountryOfficialType) []*shared_proto.HeroBasicSnapshotProto)
		if !ok {
			panic("invalid mock func, MockRankModule.CountryOfficial()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockRankModule) RemoveRankObj(a0 shared_proto.RankType, a1 int64) {
	fi := getMockFunc(s, s.RemoveRankObj)
	if fi != nil {
		f, ok := fi.(func(shared_proto.RankType, int64))
		if !ok {
			panic("invalid mock func, MockRankModule.RemoveRankObj()")
		}
		f(a0, a1)
	}

}

// 百战千军类型的特殊排行榜不能从这里获取
func (s *MockRankModule) SingleRRankListFunc(a0 shared_proto.RankType, a1 rankface.RRankListFunc) bool {
	fi := getMockFunc(s, s.SingleRRankListFunc)
	if fi != nil {
		f, ok := fi.(func(shared_proto.RankType, rankface.RRankListFunc) bool)
		if !ok {
			panic("invalid mock func, MockRankModule.SingleRRankListFunc()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockRankModule) SubTypeRRankListFunc(a0 shared_proto.RankType, a1 uint64, a2 rankface.RRankListFunc) bool {
	fi := getMockFunc(s, s.SubTypeRRankListFunc)
	if fi != nil {
		f, ok := fi.(func(shared_proto.RankType, uint64, rankface.RRankListFunc) bool)
		if !ok {
			panic("invalid mock func, MockRankModule.SubTypeRRankListFunc()")
		}
		return f(a0, a1, a2)
	}

	return false
}
func (s *MockRankModule) UpdateBaiZhanRankList(a0 []rankface.RankObj) {
	fi := getMockFunc(s, s.UpdateBaiZhanRankList)
	if fi != nil {
		f, ok := fi.(func([]rankface.RankObj))
		if !ok {
			panic("invalid mock func, MockRankModule.UpdateBaiZhanRankList()")
		}
		f(a0)
	}

}
func (s *MockRankModule) UpdateXuanyRankList(a0 []rankface.RankObj) {
	fi := getMockFunc(s, s.UpdateXuanyRankList)
	if fi != nil {
		f, ok := fi.(func([]rankface.RankObj))
		if !ok {
			panic("invalid mock func, MockRankModule.UpdateXuanyRankList()")
		}
		f(a0)
	}

}

var Realm = &MockRealm{}

type MockRealm struct{}

func (s *MockRealm) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockRealm) AddAstDefendLog(a0 int64, a1 time.Time, a2 time.Time, a3 string, a4 uint64) bool {
	fi := getMockFunc(s, s.AddAstDefendLog)
	if fi != nil {
		f, ok := fi.(func(int64, time.Time, time.Time, string, uint64) bool)
		if !ok {
			panic("invalid mock func, MockRealm.AddAstDefendLog()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return false
}

// 把玩家的基地或行营加入进来, 可以是新英雄, 可以是随机迁城, 可以是高级快速迁城. 此时英雄的城必须不是流亡状态.
// 英雄的主城或行营当前必须不能已经属于其他地区管理.
// 加入前必须已预定坐标. processed为true的话, 就算err也不需要取消预定
// isHome表示是不是主城.
// return processed 是否已处理. err 表示是否处理有错. 根据err判断错误的类型
func (s *MockRealm) AddBase(a0 int64, a1 int, a2 int, a3 realmface.AddBaseType) (bool, error) {
	fi := getMockFunc(s, s.AddBase)
	if fi != nil {
		f, ok := fi.(func(int64, int, int, realmface.AddBaseType) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.AddBase()")
		}
		return f(a0, a1, a2, a3)
	}

	return false, nil
}
func (s *MockRealm) AddGuildWorkshop(a0 int64, a1 int, a2 int, a3 int32, a4 int32) bool {
	fi := getMockFunc(s, s.AddGuildWorkshop)
	if fi != nil {
		f, ok := fi.(func(int64, int, int, int32, int32) bool)
		if !ok {
			panic("invalid mock func, MockRealm.AddGuildWorkshop()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return false
}
func (s *MockRealm) AddHeroBaoZangMonster(a0 *regdata.BaozNpcData, a1 int, a2 int, a3 int64, a4 int32) (bool, bool) {
	fi := getMockFunc(s, s.AddHeroBaoZangMonster)
	if fi != nil {
		f, ok := fi.(func(*regdata.BaozNpcData, int, int, int64, int32) (bool, bool))
		if !ok {
			panic("invalid mock func, MockRealm.AddHeroBaoZangMonster()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return false, false
}
func (s *MockRealm) AddHomeNpc(a0 iface.HeroController, a1 []*basedata.HomeNpcBaseData) (bool, error) {
	fi := getMockFunc(s, s.AddHomeNpc)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, []*basedata.HomeNpcBaseData) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.AddHomeNpc()")
		}
		return f(a0, a1)
	}

	return false, nil
}
func (s *MockRealm) AddInvasionMonster(a0 int64, a1 shared_proto.MultiLevelNpcType, a2 uint64) {
	fi := getMockFunc(s, s.AddInvasionMonster)
	if fi != nil {
		f, ok := fi.(func(int64, shared_proto.MultiLevelNpcType, uint64))
		if !ok {
			panic("invalid mock func, MockRealm.AddInvasionMonster()")
		}
		f(a0, a1, a2)
	}

}

// 英雄操作导致繁荣度增加. 传入增加的量, 由这里执行具体增加的操作
// 升级建筑在英雄线程不要直接增加繁荣度, 要调这个方法来修改繁荣度.
func (s *MockRealm) AddProsperity(a0 int64, a1 uint64) (bool, error) {
	fi := getMockFunc(s, s.AddProsperity)
	if fi != nil {
		f, ok := fi.(func(int64, uint64) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.AddProsperity()")
		}
		return f(a0, a1)
	}

	return false, nil
}
func (s *MockRealm) AddXiongNuBase(a0 xiongnuface.RResistXiongNuInfo, a1 int, a2 int, a3 int, a4 int) (int64, int32, int32) {
	fi := getMockFunc(s, s.AddXiongNuBase)
	if fi != nil {
		f, ok := fi.(func(xiongnuface.RResistXiongNuInfo, int, int, int, int) (int64, int32, int32))
		if !ok {
			panic("invalid mock func, MockRealm.AddXiongNuBase()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return 0, 0, 0
}
func (s *MockRealm) AddXiongNuTroop(a0 int64, a1 []int64, a2 []*monsterdata.MonsterMasterData) bool {
	fi := getMockFunc(s, s.AddXiongNuTroop)
	if fi != nil {
		f, ok := fi.(func(int64, []int64, []*monsterdata.MonsterMasterData) bool)
		if !ok {
			panic("invalid mock func, MockRealm.AddXiongNuTroop()")
		}
		return f(a0, a1, a2)
	}

	return false
}
func (s *MockRealm) AroundBase(a0 int, a1 int, a2 int, a3 int) bool {
	fi := getMockFunc(s, s.AroundBase)
	if fi != nil {
		f, ok := fi.(func(int, int, int, int) bool)
		if !ok {
			panic("invalid mock func, MockRealm.AroundBase()")
		}
		return f(a0, a1, a2, a3)
	}

	return false
}

// 宝藏遣返
func (s *MockRealm) BaozRepatriate(a0 iface.HeroController, a1 int64, a2 int64) (bool, error) {
	fi := getMockFunc(s, s.BaozRepatriate)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64, int64) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.BaozRepatriate()")
		}
		return f(a0, a1, a2)
	}

	return false, nil
}
func (s *MockRealm) CalcMoveSpeed(a0 int64, a1 float64) float64 {
	fi := getMockFunc(s, s.CalcMoveSpeed)
	if fi != nil {
		f, ok := fi.(func(int64, float64) float64)
		if !ok {
			panic("invalid mock func, MockRealm.CalcMoveSpeed()")
		}
		return f(a0, a1)
	}

	return 0
}

// 班师回朝
func (s *MockRealm) CancelInvasion(a0 iface.HeroController, a1 int64) (bool, error) {
	fi := getMockFunc(s, s.CancelInvasion)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.CancelInvasion()")
		}
		return f(a0, a1)
	}

	return false, nil
}

// 取消预定的坐标. 由于某种原因, 哥来不了了.
func (s *MockRealm) CancelReservedPos(a0 int, a1 int) {
	fi := getMockFunc(s, s.CancelReservedPos)
	if fi != nil {
		f, ok := fi.(func(int, int))
		if !ok {
			panic("invalid mock func, MockRealm.CancelReservedPos()")
		}
		f(a0, a1)
	}

}

// 取消缓慢迁城
// 取消缓慢迁城，快速迁城，流亡，等等会自动取消缓慢迁移
func (s *MockRealm) CancelSlowMoveBase(a0 iface.HeroController) (bool, error) {
	fi := getMockFunc(s, s.CancelSlowMoveBase)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.CancelSlowMoveBase()")
		}
		return f(a0)
	}

	return false, nil
}

// 变更个人签名
func (s *MockRealm) ChangeSign(a0 int64, a1 string) {
	fi := getMockFunc(s, s.ChangeSign)
	if fi != nil {
		f, ok := fi.(func(int64, string))
		if !ok {
			panic("invalid mock func, MockRealm.ChangeSign()")
		}
		f(a0, a1)
	}

}

// 变更个人签名
func (s *MockRealm) ChangeTitle(a0 int64, a1 uint64) {
	fi := getMockFunc(s, s.ChangeTitle)
	if fi != nil {
		f, ok := fi.(func(int64, uint64))
		if !ok {
			panic("invalid mock func, MockRealm.ChangeTitle()")
		}
		f(a0, a1)
	}

}
func (s *MockRealm) CheckCanMoveBase(a0 int64, a1 int, a2 int, a3 bool) error {
	fi := getMockFunc(s, s.CheckCanMoveBase)
	if fi != nil {
		f, ok := fi.(func(int64, int, int, bool) error)
		if !ok {
			panic("invalid mock func, MockRealm.CheckCanMoveBase()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockRealm) CheckIsFucked(a0 int64) bool {
	fi := getMockFunc(s, s.CheckIsFucked)
	if fi != nil {
		f, ok := fi.(func(int64) bool)
		if !ok {
			panic("invalid mock func, MockRealm.CheckIsFucked()")
		}
		return f(a0)
	}

	return false
}
func (s *MockRealm) ClearAstDefendLog(a0 int64) {
	fi := getMockFunc(s, s.ClearAstDefendLog)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockRealm.ClearAstDefendLog()")
		}
		f(a0)
	}

}

// 创建集结
func (s *MockRealm) CreateAssembly(a0 iface.HeroController, a1 int64, a2 uint64, a3 uint64, a4 uint64, a5 time.Duration) (bool, error) {
	fi := getMockFunc(s, s.CreateAssembly)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64, uint64, uint64, uint64, time.Duration) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.CreateAssembly()")
		}
		return f(a0, a1, a2, a3, a4, a5)
	}

	return false, nil
}

// 自己驱逐自己城里的坏人
func (s *MockRealm) Expel(a0 iface.HeroController, a1 int64, a2 uint64) (bool, bool, string, error) {
	fi := getMockFunc(s, s.Expel)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64, uint64) (bool, bool, string, error))
		if !ok {
			panic("invalid mock func, MockRealm.Expel()")
		}
		return f(a0, a1, a2)
	}

	return false, false, "", nil
}
func (s *MockRealm) GetAstDefendHeros(a0 int64) []*shared_proto.HeroBasicProto {
	fi := getMockFunc(s, s.GetAstDefendHeros)
	if fi != nil {
		f, ok := fi.(func(int64) []*shared_proto.HeroBasicProto)
		if !ok {
			panic("invalid mock func, MockRealm.GetAstDefendHeros()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockRealm) GetAstDefendLogs() *server_proto.AllAstDefendLogProto {
	fi := getMockFunc(s, s.GetAstDefendLogs)
	if fi != nil {
		f, ok := fi.(func() *server_proto.AllAstDefendLogProto)
		if !ok {
			panic("invalid mock func, MockRealm.GetAstDefendLogs()")
		}
		return f()
	}

	return nil
}
func (s *MockRealm) GetAstDefendLogsByHero(a0 int64) []*shared_proto.AstDefendLogProto {
	fi := getMockFunc(s, s.GetAstDefendLogsByHero)
	if fi != nil {
		f, ok := fi.(func(int64) []*shared_proto.AstDefendLogProto)
		if !ok {
			panic("invalid mock func, MockRealm.GetAstDefendLogsByHero()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockRealm) GetAstDefendingTroopCount(a0 int64) uint64 {
	fi := getMockFunc(s, s.GetAstDefendingTroopCount)
	if fi != nil {
		f, ok := fi.(func(int64) uint64)
		if !ok {
			panic("invalid mock func, MockRealm.GetAstDefendingTroopCount()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockRealm) GetBaseLevel(a0 uint64) *domestic_data.BaseLevelData {
	fi := getMockFunc(s, s.GetBaseLevel)
	if fi != nil {
		f, ok := fi.(func(uint64) *domestic_data.BaseLevelData)
		if !ok {
			panic("invalid mock func, MockRealm.GetBaseLevel()")
		}
		return f(a0)
	}

	return nil
}

// 获得跟我土地有冲突的玩家id
func (s *MockRealm) GetConflictHeroIds(a0 iface.HeroController) (bool, []int64) {
	fi := getMockFunc(s, s.GetConflictHeroIds)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController) (bool, []int64))
		if !ok {
			panic("invalid mock func, MockRealm.GetConflictHeroIds()")
		}
		return f(a0)
	}

	return false, nil
}
func (s *MockRealm) GetDefendingTroopCount(a0 int64) uint64 {
	fi := getMockFunc(s, s.GetDefendingTroopCount)
	if fi != nil {
		f, ok := fi.(func(int64) uint64)
		if !ok {
			panic("invalid mock func, MockRealm.GetDefendingTroopCount()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockRealm) GetHeroBaozRoBase(a0 int64) *server_proto.RoBaseProto {
	fi := getMockFunc(s, s.GetHeroBaozRoBase)
	if fi != nil {
		f, ok := fi.(func(int64) *server_proto.RoBaseProto)
		if !ok {
			panic("invalid mock func, MockRealm.GetHeroBaozRoBase()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockRealm) GetMapData() *blockdata.StitchedBlocks {
	fi := getMockFunc(s, s.GetMapData)
	if fi != nil {
		f, ok := fi.(func() *blockdata.StitchedBlocks)
		if !ok {
			panic("invalid mock func, MockRealm.GetMapData()")
		}
		return f()
	}

	return nil
}
func (s *MockRealm) GetMaxXiongNuTroopFightingAmount(a0 int64, a1 int64) (bool, uint64) {
	fi := getMockFunc(s, s.GetMaxXiongNuTroopFightingAmount)
	if fi != nil {
		f, ok := fi.(func(int64, int64) (bool, uint64))
		if !ok {
			panic("invalid mock func, MockRealm.GetMaxXiongNuTroopFightingAmount()")
		}
		return f(a0, a1)
	}

	return false, 0
}
func (s *MockRealm) GetRadius() uint64 {
	fi := getMockFunc(s, s.GetRadius)
	if fi != nil {
		f, ok := fi.(func() uint64)
		if !ok {
			panic("invalid mock func, MockRealm.GetRadius()")
		}
		return f()
	}

	return 0
}
func (s *MockRealm) GetRoBase(a0 int64) *server_proto.RoBaseProto {
	fi := getMockFunc(s, s.GetRoBase)
	if fi != nil {
		f, ok := fi.(func(int64) *server_proto.RoBaseProto)
		if !ok {
			panic("invalid mock func, MockRealm.GetRoBase()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockRealm) GetRoBaseByPos(a0 int, a1 int) *server_proto.RoBaseProto {
	fi := getMockFunc(s, s.GetRoBaseByPos)
	if fi != nil {
		f, ok := fi.(func(int, int) *server_proto.RoBaseProto)
		if !ok {
			panic("invalid mock func, MockRealm.GetRoBaseByPos()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockRealm) GetRuinsBase(a0 int, a1 int) int64 {
	fi := getMockFunc(s, s.GetRuinsBase)
	if fi != nil {
		f, ok := fi.(func(int, int) int64)
		if !ok {
			panic("invalid mock func, MockRealm.GetRuinsBase()")
		}
		return f(a0, a1)
	}

	return 0
}
func (s *MockRealm) GetXiongNuInvateTargetCount(a0 int64) i64.GetU64 {
	fi := getMockFunc(s, s.GetXiongNuInvateTargetCount)
	if fi != nil {
		f, ok := fi.(func(int64) i64.GetU64)
		if !ok {
			panic("invalid mock func, MockRealm.GetXiongNuInvateTargetCount()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockRealm) GetXiongNuTroopInfo(a0 int64, a1 int64) *shared_proto.XiongNuBaseTroopProto {
	fi := getMockFunc(s, s.GetXiongNuTroopInfo)
	if fi != nil {
		f, ok := fi.(func(int64, int64) *shared_proto.XiongNuBaseTroopProto)
		if !ok {
			panic("invalid mock func, MockRealm.GetXiongNuTroopInfo()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockRealm) GmReduceProsperity(a0 int64, a1 uint64) bool {
	fi := getMockFunc(s, s.GmReduceProsperity)
	if fi != nil {
		f, ok := fi.(func(int64, uint64) bool)
		if !ok {
			panic("invalid mock func, MockRealm.GmReduceProsperity()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockRealm) GmRefreshBaoZangNpc() {
	fi := getMockFunc(s, s.GmRefreshBaoZangNpc)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockRealm.GmRefreshBaoZangNpc()")
		}
		f()
	}

}
func (s *MockRealm) GmSpeedUpFightMe(a0 int64) {
	fi := getMockFunc(s, s.GmSpeedUpFightMe)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockRealm.GmSpeedUpFightMe()")
		}
		f(a0)
	}

}

// 破坏联盟工坊
func (s *MockRealm) HurtGuildWorkshop(a0 iface.HeroController, a1 int64) (bool, bool) {
	fi := getMockFunc(s, s.HurtGuildWorkshop)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64) (bool, bool))
		if !ok {
			panic("invalid mock func, MockRealm.HurtGuildWorkshop()")
		}
		return f(a0, a1)
	}

	return false, false
}
func (s *MockRealm) Id() int64 {
	fi := getMockFunc(s, s.Id)
	if fi != nil {
		f, ok := fi.(func() int64)
		if !ok {
			panic("invalid mock func, MockRealm.Id()")
		}
		return f()
	}

	return 0
}

// 出发攻打/帮忙驱逐
// 没有err的话调用者还需要发送成功消息
func (s *MockRealm) Invasion(a0 iface.HeroController, a1 shared_proto.TroopOperate, a2 int64, a3 uint64, a4 uint64, a5 uint64) (bool, error) {
	fi := getMockFunc(s, s.Invasion)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, shared_proto.TroopOperate, int64, uint64, uint64, uint64) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.Invasion()")
		}
		return f(a0, a1, a2, a3, a4, a5)
	}

	return false, nil
}

// 出发侦察
func (s *MockRealm) InvasionInvestigate(a0 iface.HeroController, a1 int64) (bool, error) {
	fi := getMockFunc(s, s.InvasionInvestigate)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.InvasionInvestigate()")
		}
		return f(a0, a1)
	}

	return false, nil
}
func (s *MockRealm) IsEdgeNotHomePos(a0 int, a1 int) bool {
	fi := getMockFunc(s, s.IsEdgeNotHomePos)
	if fi != nil {
		f, ok := fi.(func(int, int) bool)
		if !ok {
			panic("invalid mock func, MockRealm.IsEdgeNotHomePos()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockRealm) IsPosOpened(a0 int, a1 int) bool {
	fi := getMockFunc(s, s.IsPosOpened)
	if fi != nil {
		f, ok := fi.(func(int, int) bool)
		if !ok {
			panic("invalid mock func, MockRealm.IsPosOpened()")
		}
		return f(a0, a1)
	}

	return false
}

// 加入集结
func (s *MockRealm) JoinAssembly(a0 iface.HeroController, a1 int64, a2 int64, a3 uint64) (bool, error) {
	fi := getMockFunc(s, s.JoinAssembly)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64, int64, uint64) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.JoinAssembly()")
		}
		return f(a0, a1, a2, a3)
	}

	return false, nil
}
func (s *MockRealm) Mian(a0 int64, a1 time.Time, a2 bool) (bool, error) {
	fi := getMockFunc(s, s.Mian)
	if fi != nil {
		f, ok := fi.(func(int64, time.Time, bool) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.Mian()")
		}
		return f(a0, a1, a2)
	}

	return false, nil
}

// 在同一个地图中移动基地
// 移动前必须已预定坐标. processed为true的话, 就算err也不需要取消预定
func (s *MockRealm) MoveBase(a0 iface.HeroController, a1 int, a2 int, a3 int, a4 int, a5 bool) (bool, error) {
	fi := getMockFunc(s, s.MoveBase)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int, int, int, int, bool) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.MoveBase()")
		}
		return f(a0, a1, a2, a3, a4, a5)
	}

	return false, nil
}
func (s *MockRealm) OnHeroLogin(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroLogin)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockRealm.OnHeroLogin()")
		}
		f(a0)
	}

}
func (s *MockRealm) QueryTroopUnit(a0 iface.HeroController, a1 int64, a2 int64) (bool, error) {
	fi := getMockFunc(s, s.QueryTroopUnit)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64, int64) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.QueryTroopUnit()")
		}
		return f(a0, a1, a2)
	}

	return false, nil
}
func (s *MockRealm) RandomAroundBase(a0 int, a1 int) (int, int, bool) {
	fi := getMockFunc(s, s.RandomAroundBase)
	if fi != nil {
		f, ok := fi.(func(int, int) (int, int, bool))
		if !ok {
			panic("invalid mock func, MockRealm.RandomAroundBase()")
		}
		return f(a0, a1)
	}

	return 0, 0, false
}
func (s *MockRealm) RandomBasePos() (int, int) {
	fi := getMockFunc(s, s.RandomBasePos)
	if fi != nil {
		f, ok := fi.(func() (int, int))
		if !ok {
			panic("invalid mock func, MockRealm.RandomBasePos()")
		}
		return f()
	}

	return 0, 0
}
func (s *MockRealm) ReduceProsperity(a0 int64, a1 uint64) bool {
	fi := getMockFunc(s, s.ReduceProsperity)
	if fi != nil {
		f, ok := fi.(func(int64, uint64) bool)
		if !ok {
			panic("invalid mock func, MockRealm.ReduceProsperity()")
		}
		return f(a0, a1)
	}

	return false
}

// 把玩家在这个地图上的基地或者行营移除. 玩家的部队必须都已不在外面. 而且不能是流亡状态且归这个地图管. (流亡状态的话, 都不归这里管)
func (s *MockRealm) RemoveBase(a0 int64, a1 bool, a2 *i18n.I18nRef, a3 *i18n.I18nRef) (bool, error, int, int) {
	fi := getMockFunc(s, s.RemoveBase)
	if fi != nil {
		f, ok := fi.(func(int64, bool, *i18n.I18nRef, *i18n.I18nRef) (bool, error, int, int))
		if !ok {
			panic("invalid mock func, MockRealm.RemoveBase()")
		}
		return f(a0, a1, a2, a3)
	}

	return false, nil, 0, 0
}

// 遣返
func (s *MockRealm) Repatriate(a0 iface.HeroController, a1 int64) (bool, error) {
	fi := getMockFunc(s, s.Repatriate)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.Repatriate()")
		}
		return f(a0, a1)
	}

	return false, nil
}

// 预定一个随机主城坐标, 返回的是个可以建主城的位置
func (s *MockRealm) ReserveNewHeroHomePos(a0 uint64) (bool, int, int) {
	fi := getMockFunc(s, s.ReserveNewHeroHomePos)
	if fi != nil {
		f, ok := fi.(func(uint64) (bool, int, int))
		if !ok {
			panic("invalid mock func, MockRealm.ReserveNewHeroHomePos()")
		}
		return f(a0)
	}

	return false, 0, 0
}

// 预定一个坐标
func (s *MockRealm) ReservePos(a0 int, a1 int) bool {
	fi := getMockFunc(s, s.ReservePos)
	if fi != nil {
		f, ok := fi.(func(int, int) bool)
		if !ok {
			panic("invalid mock func, MockRealm.ReservePos()")
		}
		return f(a0, a1)
	}

	return false
}

// 在同一个场景迁城，预定一个坐标
func (s *MockRealm) ReservePosForMoveBase(a0 int, a1 int, a2 int, a3 int) bool {
	fi := getMockFunc(s, s.ReservePosForMoveBase)
	if fi != nil {
		f, ok := fi.(func(int, int, int, int) bool)
		if !ok {
			panic("invalid mock func, MockRealm.ReservePosForMoveBase()")
		}
		return f(a0, a1, a2, a3)
	}

	return false
}

// 预定一个随机主城坐标, 返回的是个可以建主城的位置
func (s *MockRealm) ReserveRandomHomePos(a0 realmface.RandomPointType) (bool, int, int) {
	fi := getMockFunc(s, s.ReserveRandomHomePos)
	if fi != nil {
		f, ok := fi.(func(realmface.RandomPointType) (bool, int, int))
		if !ok {
			panic("invalid mock func, MockRealm.ReserveRandomHomePos()")
		}
		return f(a0)
	}

	return false, 0, 0
}

// 查看集结
func (s *MockRealm) ShowAssembly(a0 iface.HeroController, a1 int64, a2 int64, a3 int32) {
	fi := getMockFunc(s, s.ShowAssembly)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64, int64, int32))
		if !ok {
			panic("invalid mock func, MockRealm.ShowAssembly()")
		}
		f(a0, a1, a2, a3)
	}

}

// 加速
func (s *MockRealm) SpeedUp(a0 iface.HeroController, a1 int64, a2 int64, a3 float64, a4 uint64) (bool, error) {
	fi := getMockFunc(s, s.SpeedUp)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64, int64, float64, uint64) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.SpeedUp()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return false, nil
}
func (s *MockRealm) StartCareMilitary(a0 iface.HeroController) bool {
	fi := getMockFunc(s, s.StartCareMilitary)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController) bool)
		if !ok {
			panic("invalid mock func, MockRealm.StartCareMilitary()")
		}
		return f(a0)
	}

	return false
}

// 开始关心这个地图, 获得地图中所有主城的信息
func (s *MockRealm) StartCareRealm(a0 iface.HeroController, a1 int, a2 int, a3 int, a4 int) bool {
	fi := getMockFunc(s, s.StartCareRealm)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int, int, int, int) bool)
		if !ok {
			panic("invalid mock func, MockRealm.StartCareRealm()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return false
}
func (s *MockRealm) StopCareRealm(a0 iface.HeroController) bool {
	fi := getMockFunc(s, s.StopCareRealm)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController) bool)
		if !ok {
			panic("invalid mock func, MockRealm.StopCareRealm()")
		}
		return f(a0)
	}

	return false
}
func (s *MockRealm) TryRemoveBaseMian(a0 int64) bool {
	fi := getMockFunc(s, s.TryRemoveBaseMian)
	if fi != nil {
		f, ok := fi.(func(int64) bool)
		if !ok {
			panic("invalid mock func, MockRealm.TryRemoveBaseMian()")
		}
		return f(a0)
	}

	return false
}

// 改变英雄基础信息（含帮派）
func (s *MockRealm) UpdateHeroBasicInfoNoBlock(a0 int64) {
	fi := getMockFunc(s, s.UpdateHeroBasicInfoNoBlock)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockRealm.UpdateHeroBasicInfoNoBlock()")
		}
		f(a0)
	}

}
func (s *MockRealm) UpdateHeroRealmInfo(a0 int64, a1 bool, a2 bool, a3 bool) {
	fi := getMockFunc(s, s.UpdateHeroRealmInfo)
	if fi != nil {
		f, ok := fi.(func(int64, bool, bool, bool))
		if !ok {
			panic("invalid mock func, MockRealm.UpdateHeroRealmInfo()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockRealm) UpdateProsperity(a0 int64) bool {
	fi := getMockFunc(s, s.UpdateProsperity)
	if fi != nil {
		f, ok := fi.(func(int64) bool)
		if !ok {
			panic("invalid mock func, MockRealm.UpdateProsperity()")
		}
		return f(a0)
	}

	return false
}
func (s *MockRealm) UpdateProsperityBuff(a0 int64, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.UpdateProsperityBuff)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockRealm.UpdateProsperityBuff()")
		}
		f(a0, a1, a2, a3)
	}

}

// 手动升级老家等级
func (s *MockRealm) UpgradeBase(a0 iface.HeroController) (bool, error) {
	fi := getMockFunc(s, s.UpgradeBase)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController) (bool, error))
		if !ok {
			panic("invalid mock func, MockRealm.UpgradeBase()")
		}
		return f(a0)
	}

	return false, nil
}

var RealmService = &MockRealmService{}

type MockRealmService struct{}

func (s *MockRealmService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockRealmService) AddProsperityFunc(a0 int64, a1 int64, a2 uint64, a3 string) iface.Func {
	fi := getMockFunc(s, s.AddProsperityFunc)
	if fi != nil {
		f, ok := fi.(func(int64, int64, uint64, string) iface.Func)
		if !ok {
			panic("invalid mock func, MockRealmService.AddProsperityFunc()")
		}
		return f(a0, a1, a2, a3)
	}

	return nil
}
func (s *MockRealmService) CheckCanMoveBase(a0 int64, a1 int, a2 int, a3 bool) bool {
	fi := getMockFunc(s, s.CheckCanMoveBase)
	if fi != nil {
		f, ok := fi.(func(int64, int, int, bool) bool)
		if !ok {
			panic("invalid mock func, MockRealmService.CheckCanMoveBase()")
		}
		return f(a0, a1, a2, a3)
	}

	return false
}
func (s *MockRealmService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockRealmService.Close()")
		}
		f()
	}

}
func (s *MockRealmService) DoMoveBase(a0 shared_proto.GoodsMoveBaseType, a1 iface.Realm, a2 iface.HeroController, a3 int, a4 int, a5 int, a6 int, a7 bool) bool {
	fi := getMockFunc(s, s.DoMoveBase)
	if fi != nil {
		f, ok := fi.(func(shared_proto.GoodsMoveBaseType, iface.Realm, iface.HeroController, int, int, int, int, bool) bool)
		if !ok {
			panic("invalid mock func, MockRealmService.DoMoveBase()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7)
	}

	return false
}
func (s *MockRealmService) GetBigMap() iface.Realm {
	fi := getMockFunc(s, s.GetBigMap)
	if fi != nil {
		f, ok := fi.(func() iface.Realm)
		if !ok {
			panic("invalid mock func, MockRealmService.GetBigMap()")
		}
		return f()
	}

	return nil
}
func (s *MockRealmService) GetRealm(a0 int64) iface.Realm {
	fi := getMockFunc(s, s.GetRealm)
	if fi != nil {
		f, ok := fi.(func(int64) iface.Realm)
		if !ok {
			panic("invalid mock func, MockRealmService.GetRealm()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockRealmService) OnGuildSnapshotRemoved(a0 int64) {
	fi := getMockFunc(s, s.OnGuildSnapshotRemoved)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockRealmService.OnGuildSnapshotRemoved()")
		}
		f(a0)
	}

}

// callback
func (s *MockRealmService) OnGuildSnapshotUpdated(a0 *guildsnapshotdata.GuildSnapshot, a1 *guildsnapshotdata.GuildSnapshot) {
	fi := getMockFunc(s, s.OnGuildSnapshotUpdated)
	if fi != nil {
		f, ok := fi.(func(*guildsnapshotdata.GuildSnapshot, *guildsnapshotdata.GuildSnapshot))
		if !ok {
			panic("invalid mock func, MockRealmService.OnGuildSnapshotUpdated()")
		}
		f(a0, a1)
	}

}

// 坐标是已经占座占好了的, 直接调用realm.AddBase就可以了
func (s *MockRealmService) ReserveNewHeroHomePos(a0 uint64) (iface.Realm, int, int) {
	fi := getMockFunc(s, s.ReserveNewHeroHomePos)
	if fi != nil {
		f, ok := fi.(func(uint64) (iface.Realm, int, int))
		if !ok {
			panic("invalid mock func, MockRealmService.ReserveNewHeroHomePos()")
		}
		return f(a0)
	}

	return nil, 0, 0
}

// 坐标是已经占座占好了的, 直接调用realm.AddBase就可以了
func (s *MockRealmService) ReserveRandomHomePos(a0 realmface.RandomPointType) (iface.Realm, int, int) {
	fi := getMockFunc(s, s.ReserveRandomHomePos)
	if fi != nil {
		f, ok := fi.(func(realmface.RandomPointType) (iface.Realm, int, int))
		if !ok {
			panic("invalid mock func, MockRealmService.ReserveRandomHomePos()")
		}
		return f(a0)
	}

	return nil, 0, 0
}
func (s *MockRealmService) StartCareMilitary(a0 iface.HeroController) {
	fi := getMockFunc(s, s.StartCareMilitary)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockRealmService.StartCareMilitary()")
		}
		f(a0)
	}

}

var RedPacketModule = &MockRedPacketModule{}

type MockRedPacketModule struct{}

func (s *MockRedPacketModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var RedPacketService = &MockRedPacketService{}

type MockRedPacketService struct{}

func (s *MockRedPacketService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockRedPacketService) AllGrabbed(a0 int64) bool {
	fi := getMockFunc(s, s.AllGrabbed)
	if fi != nil {
		f, ok := fi.(func(int64) bool)
		if !ok {
			panic("invalid mock func, MockRedPacketService.AllGrabbed()")
		}
		return f(a0)
	}

	return false
}
func (s *MockRedPacketService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockRedPacketService.Close()")
		}
		f()
	}

}
func (s *MockRedPacketService) Create(a0 int64, a1 *red_packet.RedPacketData, a2 uint64, a3 string, a4 shared_proto.ChatType) (int64, string, msg.ErrMsg) {
	fi := getMockFunc(s, s.Create)
	if fi != nil {
		f, ok := fi.(func(int64, *red_packet.RedPacketData, uint64, string, shared_proto.ChatType) (int64, string, msg.ErrMsg))
		if !ok {
			panic("invalid mock func, MockRedPacketService.Create()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return 0, "", nil
}
func (s *MockRedPacketService) Exist(a0 int64) bool {
	fi := getMockFunc(s, s.Exist)
	if fi != nil {
		f, ok := fi.(func(int64) bool)
		if !ok {
			panic("invalid mock func, MockRedPacketService.Exist()")
		}
		return f(a0)
	}

	return false
}
func (s *MockRedPacketService) Expired(a0 int64, a1 time.Time) bool {
	fi := getMockFunc(s, s.Expired)
	if fi != nil {
		f, ok := fi.(func(int64, time.Time) bool)
		if !ok {
			panic("invalid mock func, MockRedPacketService.Expired()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockRedPacketService) Grab(a0 int64, a1 int64, a2 int64) (uint64, bool, *shared_proto.RedPacketProto, msg.ErrMsg) {
	fi := getMockFunc(s, s.Grab)
	if fi != nil {
		f, ok := fi.(func(int64, int64, int64) (uint64, bool, *shared_proto.RedPacketProto, msg.ErrMsg))
		if !ok {
			panic("invalid mock func, MockRedPacketService.Grab()")
		}
		return f(a0, a1, a2)
	}

	return 0, false, nil, nil
}
func (s *MockRedPacketService) Grabbed(a0 int64, a1 int64) bool {
	fi := getMockFunc(s, s.Grabbed)
	if fi != nil {
		f, ok := fi.(func(int64, int64) bool)
		if !ok {
			panic("invalid mock func, MockRedPacketService.Grabbed()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockRedPacketService) RedPacketChatId(a0 int64) int64 {
	fi := getMockFunc(s, s.RedPacketChatId)
	if fi != nil {
		f, ok := fi.(func(int64) int64)
		if !ok {
			panic("invalid mock func, MockRedPacketService.RedPacketChatId()")
		}
		return f(a0)
	}

	return 0
}
func (s *MockRedPacketService) SetRedPacketChatId(a0 int64, a1 int64) {
	fi := getMockFunc(s, s.SetRedPacketChatId)
	if fi != nil {
		f, ok := fi.(func(int64, int64))
		if !ok {
			panic("invalid mock func, MockRedPacketService.SetRedPacketChatId()")
		}
		f(a0, a1)
	}

}

var RegionModule = &MockRegionModule{}

type MockRegionModule struct{}

func (s *MockRegionModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockRegionModule) InitHeroBase(a0 iface.HeroController, a1 time.Time, a2 uint64, a3 realmface.AddBaseType) bool {
	fi := getMockFunc(s, s.InitHeroBase)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, time.Time, uint64, realmface.AddBaseType) bool)
		if !ok {
			panic("invalid mock func, MockRegionModule.InitHeroBase()")
		}
		return f(a0, a1, a2, a3)
	}

	return false
}
func (s *MockRegionModule) UseMianGoods(a0 uint64, a1 bool, a2 iface.HeroController) bool {
	fi := getMockFunc(s, s.UseMianGoods)
	if fi != nil {
		f, ok := fi.(func(uint64, bool, iface.HeroController) bool)
		if !ok {
			panic("invalid mock func, MockRegionModule.UseMianGoods()")
		}
		return f(a0, a1, a2)
	}

	return false
}

var RelationModule = &MockRelationModule{}

type MockRelationModule struct{}

func (s *MockRelationModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var ReminderService = &MockReminderService{}

type MockReminderService struct{}

func (s *MockReminderService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockReminderService) ChangeAttackOrRobCount(a0 int64, a1 int64, a2 int64, a3 int64, a4 bool) {
	fi := getMockFunc(s, s.ChangeAttackOrRobCount)
	if fi != nil {
		f, ok := fi.(func(int64, int64, int64, int64, bool))
		if !ok {
			panic("invalid mock func, MockReminderService.ChangeAttackOrRobCount()")
		}
		f(a0, a1, a2, a3, a4)
	}

}
func (s *MockReminderService) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockReminderService.OnHeroOnline()")
		}
		f(a0)
	}

}

var SeasonService = &MockSeasonService{}

type MockSeasonService struct{}

func (s *MockSeasonService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockSeasonService) GetSeasonTickTime() tickdata.TickTime {
	fi := getMockFunc(s, s.GetSeasonTickTime)
	if fi != nil {
		f, ok := fi.(func() tickdata.TickTime)
		if !ok {
			panic("invalid mock func, MockSeasonService.GetSeasonTickTime()")
		}
		return f()
	}

	return nil
}
func (s *MockSeasonService) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockSeasonService.OnHeroOnline()")
		}
		f(a0)
	}

}
func (s *MockSeasonService) Season() *season.SeasonData {
	fi := getMockFunc(s, s.Season)
	if fi != nil {
		f, ok := fi.(func() *season.SeasonData)
		if !ok {
			panic("invalid mock func, MockSeasonService.Season()")
		}
		return f()
	}

	return nil
}
func (s *MockSeasonService) SeasonByTime(a0 time.Time) *season.SeasonData {
	fi := getMockFunc(s, s.SeasonByTime)
	if fi != nil {
		f, ok := fi.(func(time.Time) *season.SeasonData)
		if !ok {
			panic("invalid mock func, MockSeasonService.SeasonByTime()")
		}
		return f(a0)
	}

	return nil
}

var SecretTowerModule = &MockSecretTowerModule{}

type MockSecretTowerModule struct{}

func (s *MockSecretTowerModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockSecretTowerModule) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockSecretTowerModule.Close()")
		}
		f()
	}

}
func (s *MockSecretTowerModule) OnHeroOffline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOffline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockSecretTowerModule.OnHeroOffline()")
		}
		f(a0)
	}

}
func (s *MockSecretTowerModule) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockSecretTowerModule.OnHeroOnline()")
		}
		f(a0)
	}

}

var ServerStartStopTimeService = &MockServerStartStopTimeService{}

type MockServerStartStopTimeService struct{}

func (s *MockServerStartStopTimeService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockServerStartStopTimeService) IsNormalStop() bool {
	fi := getMockFunc(s, s.IsNormalStop)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockServerStartStopTimeService.IsNormalStop()")
		}
		return f()
	}

	return false
}
func (s *MockServerStartStopTimeService) SaveStartTime() {
	fi := getMockFunc(s, s.SaveStartTime)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockServerStartStopTimeService.SaveStartTime()")
		}
		f()
	}

}
func (s *MockServerStartStopTimeService) SaveStopTime() {
	fi := getMockFunc(s, s.SaveStopTime)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockServerStartStopTimeService.SaveStopTime()")
		}
		f()
	}

}

var ServiceDep = &MockServiceDep{}

type MockServiceDep struct{}

func (s *MockServiceDep) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockServiceDep) Broadcast() iface.BroadcastService {
	fi := getMockFunc(s, s.Broadcast)
	if fi != nil {
		f, ok := fi.(func() iface.BroadcastService)
		if !ok {
			panic("invalid mock func, MockServiceDep.Broadcast()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Chat() iface.ChatService {
	fi := getMockFunc(s, s.Chat)
	if fi != nil {
		f, ok := fi.(func() iface.ChatService)
		if !ok {
			panic("invalid mock func, MockServiceDep.Chat()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Country() iface.CountryService {
	fi := getMockFunc(s, s.Country)
	if fi != nil {
		f, ok := fi.(func() iface.CountryService)
		if !ok {
			panic("invalid mock func, MockServiceDep.Country()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Datas() iface.ConfigDatas {
	fi := getMockFunc(s, s.Datas)
	if fi != nil {
		f, ok := fi.(func() iface.ConfigDatas)
		if !ok {
			panic("invalid mock func, MockServiceDep.Datas()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Db() iface.DbService {
	fi := getMockFunc(s, s.Db)
	if fi != nil {
		f, ok := fi.(func() iface.DbService)
		if !ok {
			panic("invalid mock func, MockServiceDep.Db()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Fight() iface.FightService {
	fi := getMockFunc(s, s.Fight)
	if fi != nil {
		f, ok := fi.(func() iface.FightService)
		if !ok {
			panic("invalid mock func, MockServiceDep.Fight()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) FightX() iface.FightXService {
	fi := getMockFunc(s, s.FightX)
	if fi != nil {
		f, ok := fi.(func() iface.FightXService)
		if !ok {
			panic("invalid mock func, MockServiceDep.FightX()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Guild() iface.GuildService {
	fi := getMockFunc(s, s.Guild)
	if fi != nil {
		f, ok := fi.(func() iface.GuildService)
		if !ok {
			panic("invalid mock func, MockServiceDep.Guild()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) GuildSnapshot() iface.GuildSnapshotService {
	fi := getMockFunc(s, s.GuildSnapshot)
	if fi != nil {
		f, ok := fi.(func() iface.GuildSnapshotService)
		if !ok {
			panic("invalid mock func, MockServiceDep.GuildSnapshot()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) HeroData() iface.HeroDataService {
	fi := getMockFunc(s, s.HeroData)
	if fi != nil {
		f, ok := fi.(func() iface.HeroDataService)
		if !ok {
			panic("invalid mock func, MockServiceDep.HeroData()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) HeroSnapshot() iface.HeroSnapshotService {
	fi := getMockFunc(s, s.HeroSnapshot)
	if fi != nil {
		f, ok := fi.(func() iface.HeroSnapshotService)
		if !ok {
			panic("invalid mock func, MockServiceDep.HeroSnapshot()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Mail() iface.MailModule {
	fi := getMockFunc(s, s.Mail)
	if fi != nil {
		f, ok := fi.(func() iface.MailModule)
		if !ok {
			panic("invalid mock func, MockServiceDep.Mail()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Mingc() iface.MingcService {
	fi := getMockFunc(s, s.Mingc)
	if fi != nil {
		f, ok := fi.(func() iface.MingcService)
		if !ok {
			panic("invalid mock func, MockServiceDep.Mingc()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Push() iface.PushService {
	fi := getMockFunc(s, s.Push)
	if fi != nil {
		f, ok := fi.(func() iface.PushService)
		if !ok {
			panic("invalid mock func, MockServiceDep.Push()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) SvrConf() iface.IndividualServerConfig {
	fi := getMockFunc(s, s.SvrConf)
	if fi != nil {
		f, ok := fi.(func() iface.IndividualServerConfig)
		if !ok {
			panic("invalid mock func, MockServiceDep.SvrConf()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Time() iface.TimeService {
	fi := getMockFunc(s, s.Time)
	if fi != nil {
		f, ok := fi.(func() iface.TimeService)
		if !ok {
			panic("invalid mock func, MockServiceDep.Time()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) Tlog() iface.TlogService {
	fi := getMockFunc(s, s.Tlog)
	if fi != nil {
		f, ok := fi.(func() iface.TlogService)
		if !ok {
			panic("invalid mock func, MockServiceDep.Tlog()")
		}
		return f()
	}

	return nil
}
func (s *MockServiceDep) World() iface.WorldService {
	fi := getMockFunc(s, s.World)
	if fi != nil {
		f, ok := fi.(func() iface.WorldService)
		if !ok {
			panic("invalid mock func, MockServiceDep.World()")
		}
		return f()
	}

	return nil
}

var ShopModule = &MockShopModule{}

type MockShopModule struct{}

func (s *MockShopModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var StrategyModule = &MockStrategyModule{}

type MockStrategyModule struct{}

func (s *MockStrategyModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockStrategyModule) GMStrategy(a0 uint64, a1 iface.HeroController) {
	fi := getMockFunc(s, s.GMStrategy)
	if fi != nil {
		f, ok := fi.(func(uint64, iface.HeroController))
		if !ok {
			panic("invalid mock func, MockStrategyModule.GMStrategy()")
		}
		f(a0, a1)
	}

}

var StressModule = &MockStressModule{}

type MockStressModule struct{}

func (s *MockStressModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var SurveyModule = &MockSurveyModule{}

type MockSurveyModule struct{}

func (s *MockSurveyModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockSurveyModule) GmGiveSurveyPrize(a0 int64, a1 string) {
	fi := getMockFunc(s, s.GmGiveSurveyPrize)
	if fi != nil {
		f, ok := fi.(func(int64, string))
		if !ok {
			panic("invalid mock func, MockSurveyModule.GmGiveSurveyPrize()")
		}
		f(a0, a1)
	}

}

var TagModule = &MockTagModule{}

type MockTagModule struct{}

func (s *MockTagModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var TaskModule = &MockTaskModule{}

type MockTaskModule struct{}

func (s *MockTaskModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockTaskModule) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockTaskModule.OnHeroOnline()")
		}
		f(a0)
	}

}

var TeachModule = &MockTeachModule{}

type MockTeachModule struct{}

func (s *MockTeachModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var TickerService = &MockTickerService{}

type MockTickerService struct{}

func (s *MockTickerService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockTickerService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockTickerService.Close()")
		}
		f()
	}

}
func (s *MockTickerService) GetDailyMcTickTime() tickdata.TickTime {
	fi := getMockFunc(s, s.GetDailyMcTickTime)
	if fi != nil {
		f, ok := fi.(func() tickdata.TickTime)
		if !ok {
			panic("invalid mock func, MockTickerService.GetDailyMcTickTime()")
		}
		return f()
	}

	return nil
}
func (s *MockTickerService) GetDailyTickTime() tickdata.TickTime {
	fi := getMockFunc(s, s.GetDailyTickTime)
	if fi != nil {
		f, ok := fi.(func() tickdata.TickTime)
		if !ok {
			panic("invalid mock func, MockTickerService.GetDailyTickTime()")
		}
		return f()
	}

	return nil
}
func (s *MockTickerService) GetDailyZeroTickTime() tickdata.TickTime {
	fi := getMockFunc(s, s.GetDailyZeroTickTime)
	if fi != nil {
		f, ok := fi.(func() tickdata.TickTime)
		if !ok {
			panic("invalid mock func, MockTickerService.GetDailyZeroTickTime()")
		}
		return f()
	}

	return nil
}
func (s *MockTickerService) GetPer10MinuteTickTime() tickdata.TickTime {
	fi := getMockFunc(s, s.GetPer10MinuteTickTime)
	if fi != nil {
		f, ok := fi.(func() tickdata.TickTime)
		if !ok {
			panic("invalid mock func, MockTickerService.GetPer10MinuteTickTime()")
		}
		return f()
	}

	return nil
}
func (s *MockTickerService) GetPer30MinuteTickTime() tickdata.TickTime {
	fi := getMockFunc(s, s.GetPer30MinuteTickTime)
	if fi != nil {
		f, ok := fi.(func() tickdata.TickTime)
		if !ok {
			panic("invalid mock func, MockTickerService.GetPer30MinuteTickTime()")
		}
		return f()
	}

	return nil
}
func (s *MockTickerService) GetPerHourTickTime() tickdata.TickTime {
	fi := getMockFunc(s, s.GetPerHourTickTime)
	if fi != nil {
		f, ok := fi.(func() tickdata.TickTime)
		if !ok {
			panic("invalid mock func, MockTickerService.GetPerHourTickTime()")
		}
		return f()
	}

	return nil
}
func (s *MockTickerService) GetPerMinuteTickTime() tickdata.TickTime {
	fi := getMockFunc(s, s.GetPerMinuteTickTime)
	if fi != nil {
		f, ok := fi.(func() tickdata.TickTime)
		if !ok {
			panic("invalid mock func, MockTickerService.GetPerMinuteTickTime()")
		}
		return f()
	}

	return nil
}
func (s *MockTickerService) GetWeeklyTickTime() tickdata.TickTime {
	fi := getMockFunc(s, s.GetWeeklyTickTime)
	if fi != nil {
		f, ok := fi.(func() tickdata.TickTime)
		if !ok {
			panic("invalid mock func, MockTickerService.GetWeeklyTickTime()")
		}
		return f()
	}

	return nil
}
func (s *MockTickerService) TickPer10Minute(a0 string, a1 iface.TickFunc) iface.Func {
	fi := getMockFunc(s, s.TickPer10Minute)
	if fi != nil {
		f, ok := fi.(func(string, iface.TickFunc) iface.Func)
		if !ok {
			panic("invalid mock func, MockTickerService.TickPer10Minute()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockTickerService) TickPer30Minute(a0 string, a1 iface.TickFunc) iface.Func {
	fi := getMockFunc(s, s.TickPer30Minute)
	if fi != nil {
		f, ok := fi.(func(string, iface.TickFunc) iface.Func)
		if !ok {
			panic("invalid mock func, MockTickerService.TickPer30Minute()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockTickerService) TickPerDay(a0 string, a1 iface.TickFunc) iface.Func {
	fi := getMockFunc(s, s.TickPerDay)
	if fi != nil {
		f, ok := fi.(func(string, iface.TickFunc) iface.Func)
		if !ok {
			panic("invalid mock func, MockTickerService.TickPerDay()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockTickerService) TickPerDayZero(a0 string, a1 iface.TickFunc) iface.Func {
	fi := getMockFunc(s, s.TickPerDayZero)
	if fi != nil {
		f, ok := fi.(func(string, iface.TickFunc) iface.Func)
		if !ok {
			panic("invalid mock func, MockTickerService.TickPerDayZero()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockTickerService) TickPerHour(a0 string, a1 iface.TickFunc) iface.Func {
	fi := getMockFunc(s, s.TickPerHour)
	if fi != nil {
		f, ok := fi.(func(string, iface.TickFunc) iface.Func)
		if !ok {
			panic("invalid mock func, MockTickerService.TickPerHour()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockTickerService) TickPerMinute(a0 string, a1 iface.TickFunc) iface.Func {
	fi := getMockFunc(s, s.TickPerMinute)
	if fi != nil {
		f, ok := fi.(func(string, iface.TickFunc) iface.Func)
		if !ok {
			panic("invalid mock func, MockTickerService.TickPerMinute()")
		}
		return f(a0, a1)
	}

	return nil
}
func (s *MockTickerService) TickTickPerWeek(a0 string, a1 iface.TickFunc) iface.Func {
	fi := getMockFunc(s, s.TickTickPerWeek)
	if fi != nil {
		f, ok := fi.(func(string, iface.TickFunc) iface.Func)
		if !ok {
			panic("invalid mock func, MockTickerService.TickTickPerWeek()")
		}
		return f(a0, a1)
	}

	return nil
}

var TimeLimitGiftService = &MockTimeLimitGiftService{}

type MockTimeLimitGiftService struct{}

func (s *MockTimeLimitGiftService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockTimeLimitGiftService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockTimeLimitGiftService.Close()")
		}
		f()
	}

}
func (s *MockTimeLimitGiftService) EncodeClient() []*shared_proto.TimeLimitGiftProto {
	fi := getMockFunc(s, s.EncodeClient)
	if fi != nil {
		f, ok := fi.(func() []*shared_proto.TimeLimitGiftProto)
		if !ok {
			panic("invalid mock func, MockTimeLimitGiftService.EncodeClient()")
		}
		return f()
	}

	return nil
}
func (s *MockTimeLimitGiftService) GetGiftEndTime(a0 uint64) (time.Time, bool) {
	fi := getMockFunc(s, s.GetGiftEndTime)
	if fi != nil {
		f, ok := fi.(func(uint64) (time.Time, bool))
		if !ok {
			panic("invalid mock func, MockTimeLimitGiftService.GetGiftEndTime()")
		}
		return f(a0)
	}

	return time.Time{}, false
}
func (s *MockTimeLimitGiftService) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockTimeLimitGiftService.OnHeroOnline()")
		}
		f(a0)
	}

}

var TimeService = &MockTimeService{}

type MockTimeService struct{}

func (s *MockTimeService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockTimeService) CurrentTime() time.Time {
	fi := getMockFunc(s, s.CurrentTime)
	if fi != nil {
		f, ok := fi.(func() time.Time)
		if !ok {
			panic("invalid mock func, MockTimeService.CurrentTime()")
		}
		return f()
	}

	return time.Time{}
}

var TlogBaseService = &MockTlogBaseService{}

type MockTlogBaseService struct{}

func (s *MockTlogBaseService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockTlogBaseService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockTlogBaseService.Close()")
		}
		f()
	}

}
func (s *MockTlogBaseService) WriteTlog(a0 string) bool {
	fi := getMockFunc(s, s.WriteTlog)
	if fi != nil {
		f, ok := fi.(func(string) bool)
		if !ok {
			panic("invalid mock func, MockTlogBaseService.WriteTlog()")
		}
		return f(a0)
	}

	return false
}

var TlogService = &MockTlogService{}

type MockTlogService struct{}

func (s *MockTlogService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockTlogService) BuildAccountRegister(a0 int64, a1 *shared_proto.TencentInfoProto) string {
	fi := getMockFunc(s, s.BuildAccountRegister)
	if fi != nil {
		f, ok := fi.(func(int64, *shared_proto.TencentInfoProto) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildAccountRegister()")
		}
		return f(a0, a1)
	}

	return ""
}
func (s *MockTlogService) BuildAdvanceSoulFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64) string {
	fi := getMockFunc(s, s.BuildAdvanceSoulFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildAdvanceSoulFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8)
	}

	return ""
}
func (s *MockTlogService) BuildAnswerFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 bool) string {
	fi := getMockFunc(s, s.BuildAnswerFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, bool) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildAnswerFlow()")
		}
		return f(a0, a1, a2, a3)
	}

	return ""
}
func (s *MockTlogService) BuildBaiZhanFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64) string {
	fi := getMockFunc(s, s.BuildBaiZhanFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildBaiZhanFlow()")
		}
		return f(a0, a1, a2, a3)
	}

	return ""
}
func (s *MockTlogService) BuildCareFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 string, a5 uint64, a6 uint64) string {
	fi := getMockFunc(s, s.BuildCareFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, string, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildCareFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6)
	}

	return ""
}
func (s *MockTlogService) BuildChangeCaptainFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64) string {
	fi := getMockFunc(s, s.BuildChangeCaptainFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildChangeCaptainFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	}

	return ""
}
func (s *MockTlogService) BuildChatFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64) string {
	fi := getMockFunc(s, s.BuildChatFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildChatFlow()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return ""
}
func (s *MockTlogService) BuildCityExpFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64) string {
	fi := getMockFunc(s, s.BuildCityExpFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildCityExpFlow()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return ""
}
func (s *MockTlogService) BuildEquipmentAddStarFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) string {
	fi := getMockFunc(s, s.BuildEquipmentAddStarFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildEquipmentAddStarFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6)
	}

	return ""
}
func (s *MockTlogService) BuildFarmFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 string, a4 uint64, a5 uint64, a6 uint64) string {
	fi := getMockFunc(s, s.BuildFarmFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, string, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildFarmFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6)
	}

	return ""
}
func (s *MockTlogService) BuildFishFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64) string {
	fi := getMockFunc(s, s.BuildFishFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildFishFlow()")
		}
		return f(a0, a1, a2, a3)
	}

	return ""
}
func (s *MockTlogService) BuildGUOGUANFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64) string {
	fi := getMockFunc(s, s.BuildGUOGUANFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildGUOGUANFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7)
	}

	return ""
}
func (s *MockTlogService) BuildGameSvrState() string {
	fi := getMockFunc(s, s.BuildGameSvrState)
	if fi != nil {
		f, ok := fi.(func() string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildGameSvrState()")
		}
		return f()
	}

	return ""
}
func (s *MockTlogService) BuildGameplayFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64) string {
	fi := getMockFunc(s, s.BuildGameplayFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildGameplayFlow()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return ""
}
func (s *MockTlogService) BuildGuideFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64) string {
	fi := getMockFunc(s, s.BuildGuideFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildGuideFlow()")
		}
		return f(a0, a1, a2, a3)
	}

	return ""
}
func (s *MockTlogService) BuildGuildFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64) string {
	fi := getMockFunc(s, s.BuildGuildFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildGuildFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5)
	}

	return ""
}
func (s *MockTlogService) BuildItemFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 int64) string {
	fi := getMockFunc(s, s.BuildItemFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, int64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildItemFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
	}

	return ""
}
func (s *MockTlogService) BuildKingExpFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64) string {
	fi := getMockFunc(s, s.BuildKingExpFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildKingExpFlow()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return ""
}
func (s *MockTlogService) BuildMailFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 string) string {
	fi := getMockFunc(s, s.BuildMailFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, string) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildMailFlow()")
		}
		return f(a0, a1, a2, a3)
	}

	return ""
}
func (s *MockTlogService) BuildMoneyFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) string {
	fi := getMockFunc(s, s.BuildMoneyFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildMoneyFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6)
	}

	return ""
}
func (s *MockTlogService) BuildMountRefreshFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) string {
	fi := getMockFunc(s, s.BuildMountRefreshFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildMountRefreshFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6)
	}

	return ""
}
func (s *MockTlogService) BuildMoveCitylFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 int64, a4 int64, a5 int64, a6 int64) string {
	fi := getMockFunc(s, s.BuildMoveCitylFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, int64, int64, int64, int64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildMoveCitylFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6)
	}

	return ""
}
func (s *MockTlogService) BuildNationalFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64) string {
	fi := getMockFunc(s, s.BuildNationalFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildNationalFlow()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return ""
}
func (s *MockTlogService) BuildPlayerCultivateFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64, a12 uint64) string {
	fi := getMockFunc(s, s.BuildPlayerCultivateFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildPlayerCultivateFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12)
	}

	return ""
}
func (s *MockTlogService) BuildPlayerEquipFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64) string {
	fi := getMockFunc(s, s.BuildPlayerEquipFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildPlayerEquipFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7)
	}

	return ""
}
func (s *MockTlogService) BuildPlayerExpDrugFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64) string {
	fi := getMockFunc(s, s.BuildPlayerExpDrugFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildPlayerExpDrugFlow()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return ""
}
func (s *MockTlogService) BuildPlayerHaunterFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64) string {
	fi := getMockFunc(s, s.BuildPlayerHaunterFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildPlayerHaunterFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7)
	}

	return ""
}
func (s *MockTlogService) BuildPlayerLogin(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64, a12 uint64, a13 uint64, a14 uint64, a15 uint64, a16 []uint64, a17 uint64, a18 []uint64, a19 uint64, a20 []uint64, a21 uint64) string {
	fi := getMockFunc(s, s.BuildPlayerLogin)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildPlayerLogin()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16, a17, a18, a19, a20, a21)
	}

	return ""
}
func (s *MockTlogService) BuildPlayerLogout(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64, a12 uint64, a13 uint64, a14 uint64, a15 uint64, a16 uint64, a17 uint64, a18 uint64, a19 []uint64, a20 uint64, a21 []uint64, a22 uint64, a23 []uint64, a24 uint64) string {
	fi := getMockFunc(s, s.BuildPlayerLogout)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildPlayerLogout()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16, a17, a18, a19, a20, a21, a22, a23, a24)
	}

	return ""
}
func (s *MockTlogService) BuildPlayerRegister(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto) string {
	fi := getMockFunc(s, s.BuildPlayerRegister)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildPlayerRegister()")
		}
		return f(a0, a1)
	}

	return ""
}
func (s *MockTlogService) BuildRefreshFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64) string {
	fi := getMockFunc(s, s.BuildRefreshFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildRefreshFlow()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return ""
}
func (s *MockTlogService) BuildResearchFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64) string {
	fi := getMockFunc(s, s.BuildResearchFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildResearchFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5)
	}

	return ""
}
func (s *MockTlogService) BuildResourceStockFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64) string {
	fi := getMockFunc(s, s.BuildResourceStockFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildResourceStockFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7)
	}

	return ""
}
func (s *MockTlogService) BuildRoundFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64, a12 uint64, a13 uint64, a14 uint64, a15 uint64, a16 uint64, a17 uint64, a18 uint64, a19 uint64, a20 uint64, a21 uint64, a22 uint64, a23 uint64, a24 uint64, a25 uint64, a26 uint64, a27 uint64, a28 uint64, a29 uint64, a30 uint64) string {
	fi := getMockFunc(s, s.BuildRoundFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildRoundFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16, a17, a18, a19, a20, a21, a22, a23, a24, a25, a26, a27, a28, a29, a30)
	}

	return ""
}
func (s *MockTlogService) BuildSnsFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64) string {
	fi := getMockFunc(s, s.BuildSnsFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildSnsFlow()")
		}
		return f(a0, a1, a2, a3)
	}

	return ""
}
func (s *MockTlogService) BuildSpeedUpFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64) string {
	fi := getMockFunc(s, s.BuildSpeedUpFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildSpeedUpFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7)
	}

	return ""
}
func (s *MockTlogService) BuildStrenghBuildingFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64) string {
	fi := getMockFunc(s, s.BuildStrenghBuildingFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildStrenghBuildingFlow()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return ""
}
func (s *MockTlogService) BuildStrenghEquipmentFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64) string {
	fi := getMockFunc(s, s.BuildStrenghEquipmentFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildStrenghEquipmentFlow()")
		}
		return f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11)
	}

	return ""
}
func (s *MockTlogService) BuildTaskFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64, a4 uint64) string {
	fi := getMockFunc(s, s.BuildTaskFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildTaskFlow()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return ""
}
func (s *MockTlogService) BuildVipLevelFlow(a0 entity.TlogHero, a1 *shared_proto.TencentInfoProto, a2 uint64, a3 uint64) string {
	fi := getMockFunc(s, s.BuildVipLevelFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64) string)
		if !ok {
			panic("invalid mock func, MockTlogService.BuildVipLevelFlow()")
		}
		return f(a0, a1, a2, a3)
	}

	return ""
}
func (s *MockTlogService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockTlogService.Close()")
		}
		f()
	}

}
func (s *MockTlogService) DontGenTlog() bool {
	fi := getMockFunc(s, s.DontGenTlog)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockTlogService.DontGenTlog()")
		}
		return f()
	}

	return false
}
func (s *MockTlogService) TlogAccountRegister(a0 int64) {
	fi := getMockFunc(s, s.TlogAccountRegister)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogAccountRegister()")
		}
		f(a0)
	}

}
func (s *MockTlogService) TlogAdvanceSoulFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64) {
	fi := getMockFunc(s, s.TlogAdvanceSoulFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogAdvanceSoulFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7)
	}

}
func (s *MockTlogService) TlogAdvanceSoulFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64) {
	fi := getMockFunc(s, s.TlogAdvanceSoulFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogAdvanceSoulFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7)
	}

}
func (s *MockTlogService) TlogAnswerFlow(a0 entity.TlogHero, a1 uint64, a2 bool) {
	fi := getMockFunc(s, s.TlogAnswerFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, bool))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogAnswerFlow()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogAnswerFlowById(a0 int64, a1 uint64, a2 bool) {
	fi := getMockFunc(s, s.TlogAnswerFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, bool))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogAnswerFlowById()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogBaiZhanFlow(a0 entity.TlogHero, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.TlogBaiZhanFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogBaiZhanFlow()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogBaiZhanFlowById(a0 int64, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.TlogBaiZhanFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogBaiZhanFlowById()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogCareFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 string, a4 uint64, a5 uint64) {
	fi := getMockFunc(s, s.TlogCareFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, string, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogCareFlow()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogCareFlowById(a0 int64, a1 uint64, a2 uint64, a3 string, a4 uint64, a5 uint64) {
	fi := getMockFunc(s, s.TlogCareFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, string, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogCareFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogChangeCaptainFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64) {
	fi := getMockFunc(s, s.TlogChangeCaptainFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogChangeCaptainFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8)
	}

}
func (s *MockTlogService) TlogChangeCaptainFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64) {
	fi := getMockFunc(s, s.TlogChangeCaptainFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogChangeCaptainFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8)
	}

}
func (s *MockTlogService) TlogChatFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogChatFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogChatFlow()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogChatFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogChatFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogChatFlowById()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogCityExpFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogCityExpFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogCityExpFlow()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogCityExpFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogCityExpFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogCityExpFlowById()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogEquipmentAddStarFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64) {
	fi := getMockFunc(s, s.TlogEquipmentAddStarFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogEquipmentAddStarFlow()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogEquipmentAddStarFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64) {
	fi := getMockFunc(s, s.TlogEquipmentAddStarFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogEquipmentAddStarFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogFarmFlow(a0 entity.TlogHero, a1 uint64, a2 string, a3 uint64, a4 uint64, a5 uint64) {
	fi := getMockFunc(s, s.TlogFarmFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, string, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogFarmFlow()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogFarmFlowById(a0 int64, a1 uint64, a2 string, a3 uint64, a4 uint64, a5 uint64) {
	fi := getMockFunc(s, s.TlogFarmFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, string, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogFarmFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogFishFlow(a0 entity.TlogHero, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.TlogFishFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogFishFlow()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogFishFlowById(a0 int64, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.TlogFishFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogFishFlowById()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogGUOGUANFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) {
	fi := getMockFunc(s, s.TlogGUOGUANFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogGUOGUANFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}
func (s *MockTlogService) TlogGUOGUANFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) {
	fi := getMockFunc(s, s.TlogGUOGUANFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogGUOGUANFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}
func (s *MockTlogService) TlogGameSvrState() {
	fi := getMockFunc(s, s.TlogGameSvrState)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockTlogService.TlogGameSvrState()")
		}
		f()
	}

}
func (s *MockTlogService) TlogGameplayFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogGameplayFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogGameplayFlow()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogGameplayFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogGameplayFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogGameplayFlowById()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogGuideFlow(a0 entity.TlogHero, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.TlogGuideFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogGuideFlow()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogGuideFlowById(a0 int64, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.TlogGuideFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogGuideFlowById()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogGuildFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64) {
	fi := getMockFunc(s, s.TlogGuildFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogGuildFlow()")
		}
		f(a0, a1, a2, a3, a4)
	}

}
func (s *MockTlogService) TlogGuildFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64) {
	fi := getMockFunc(s, s.TlogGuildFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogGuildFlowById()")
		}
		f(a0, a1, a2, a3, a4)
	}

}
func (s *MockTlogService) TlogItemFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 int64) {
	fi := getMockFunc(s, s.TlogItemFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, int64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogItemFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	}

}
func (s *MockTlogService) TlogItemFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 int64) {
	fi := getMockFunc(s, s.TlogItemFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, int64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogItemFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9)
	}

}
func (s *MockTlogService) TlogKingExpFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogKingExpFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogKingExpFlow()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogKingExpFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogKingExpFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogKingExpFlowById()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogMailFlow(a0 entity.TlogHero, a1 uint64, a2 string) {
	fi := getMockFunc(s, s.TlogMailFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, string))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogMailFlow()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogMailFlowById(a0 int64, a1 uint64, a2 string) {
	fi := getMockFunc(s, s.TlogMailFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, string))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogMailFlowById()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogMoneyFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64) {
	fi := getMockFunc(s, s.TlogMoneyFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogMoneyFlow()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogMoneyFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64) {
	fi := getMockFunc(s, s.TlogMoneyFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogMoneyFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogMountRefreshFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64) {
	fi := getMockFunc(s, s.TlogMountRefreshFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogMountRefreshFlow()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogMountRefreshFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64) {
	fi := getMockFunc(s, s.TlogMountRefreshFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogMountRefreshFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogMoveCitylFlow(a0 entity.TlogHero, a1 uint64, a2 int64, a3 int64, a4 int64, a5 int64) {
	fi := getMockFunc(s, s.TlogMoveCitylFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, int64, int64, int64, int64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogMoveCitylFlow()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogMoveCitylFlowById(a0 int64, a1 uint64, a2 int64, a3 int64, a4 int64, a5 int64) {
	fi := getMockFunc(s, s.TlogMoveCitylFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, int64, int64, int64, int64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogMoveCitylFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5)
	}

}
func (s *MockTlogService) TlogNationalFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogNationalFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogNationalFlow()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogNationalFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogNationalFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogNationalFlowById()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogPlayerCultivateFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64) {
	fi := getMockFunc(s, s.TlogPlayerCultivateFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerCultivateFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11)
	}

}
func (s *MockTlogService) TlogPlayerCultivateFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64) {
	fi := getMockFunc(s, s.TlogPlayerCultivateFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerCultivateFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11)
	}

}
func (s *MockTlogService) TlogPlayerEquipFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) {
	fi := getMockFunc(s, s.TlogPlayerEquipFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerEquipFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}
func (s *MockTlogService) TlogPlayerEquipFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) {
	fi := getMockFunc(s, s.TlogPlayerEquipFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerEquipFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}
func (s *MockTlogService) TlogPlayerExpDrugFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogPlayerExpDrugFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerExpDrugFlow()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogPlayerExpDrugFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogPlayerExpDrugFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerExpDrugFlowById()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogPlayerHaunterFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) {
	fi := getMockFunc(s, s.TlogPlayerHaunterFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerHaunterFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}
func (s *MockTlogService) TlogPlayerHaunterFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) {
	fi := getMockFunc(s, s.TlogPlayerHaunterFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerHaunterFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}
func (s *MockTlogService) TlogPlayerLogin(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64, a12 uint64, a13 uint64, a14 uint64, a15 []uint64, a16 uint64, a17 []uint64, a18 uint64, a19 []uint64, a20 uint64) {
	fi := getMockFunc(s, s.TlogPlayerLogin)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerLogin()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16, a17, a18, a19, a20)
	}

}
func (s *MockTlogService) TlogPlayerLoginById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64, a12 uint64, a13 uint64, a14 uint64, a15 []uint64, a16 uint64, a17 []uint64, a18 uint64, a19 []uint64, a20 uint64) {
	fi := getMockFunc(s, s.TlogPlayerLoginById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerLoginById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16, a17, a18, a19, a20)
	}

}
func (s *MockTlogService) TlogPlayerLogout(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64, a12 uint64, a13 uint64, a14 uint64, a15 uint64, a16 uint64, a17 uint64, a18 []uint64, a19 uint64, a20 []uint64, a21 uint64, a22 []uint64, a23 uint64) {
	fi := getMockFunc(s, s.TlogPlayerLogout)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerLogout()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16, a17, a18, a19, a20, a21, a22, a23)
	}

}
func (s *MockTlogService) TlogPlayerLogoutById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64, a12 uint64, a13 uint64, a14 uint64, a15 uint64, a16 uint64, a17 uint64, a18 []uint64, a19 uint64, a20 []uint64, a21 uint64, a22 []uint64, a23 uint64) {
	fi := getMockFunc(s, s.TlogPlayerLogoutById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerLogoutById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16, a17, a18, a19, a20, a21, a22, a23)
	}

}
func (s *MockTlogService) TlogPlayerRegister(a0 entity.TlogHero) {
	fi := getMockFunc(s, s.TlogPlayerRegister)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerRegister()")
		}
		f(a0)
	}

}
func (s *MockTlogService) TlogPlayerRegisterById(a0 int64) {
	fi := getMockFunc(s, s.TlogPlayerRegisterById)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogPlayerRegisterById()")
		}
		f(a0)
	}

}
func (s *MockTlogService) TlogRefreshFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogRefreshFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogRefreshFlow()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogRefreshFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogRefreshFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogRefreshFlowById()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogResearchFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64) {
	fi := getMockFunc(s, s.TlogResearchFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogResearchFlow()")
		}
		f(a0, a1, a2, a3, a4)
	}

}
func (s *MockTlogService) TlogResearchFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64) {
	fi := getMockFunc(s, s.TlogResearchFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogResearchFlowById()")
		}
		f(a0, a1, a2, a3, a4)
	}

}
func (s *MockTlogService) TlogResourceStockFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) {
	fi := getMockFunc(s, s.TlogResourceStockFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogResourceStockFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}
func (s *MockTlogService) TlogResourceStockFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) {
	fi := getMockFunc(s, s.TlogResourceStockFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogResourceStockFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}
func (s *MockTlogService) TlogRoundFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64, a12 uint64, a13 uint64, a14 uint64, a15 uint64, a16 uint64, a17 uint64, a18 uint64, a19 uint64, a20 uint64, a21 uint64, a22 uint64, a23 uint64, a24 uint64, a25 uint64, a26 uint64, a27 uint64, a28 uint64, a29 uint64) {
	fi := getMockFunc(s, s.TlogRoundFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogRoundFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16, a17, a18, a19, a20, a21, a22, a23, a24, a25, a26, a27, a28, a29)
	}

}
func (s *MockTlogService) TlogRoundFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64, a11 uint64, a12 uint64, a13 uint64, a14 uint64, a15 uint64, a16 uint64, a17 uint64, a18 uint64, a19 uint64, a20 uint64, a21 uint64, a22 uint64, a23 uint64, a24 uint64, a25 uint64, a26 uint64, a27 uint64, a28 uint64, a29 uint64) {
	fi := getMockFunc(s, s.TlogRoundFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogRoundFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15, a16, a17, a18, a19, a20, a21, a22, a23, a24, a25, a26, a27, a28, a29)
	}

}
func (s *MockTlogService) TlogSnsFlow(a0 entity.TlogHero, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.TlogSnsFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogSnsFlow()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogSnsFlowById(a0 int64, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.TlogSnsFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogSnsFlowById()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogSpeedUpFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) {
	fi := getMockFunc(s, s.TlogSpeedUpFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogSpeedUpFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}
func (s *MockTlogService) TlogSpeedUpFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64) {
	fi := getMockFunc(s, s.TlogSpeedUpFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogSpeedUpFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6)
	}

}
func (s *MockTlogService) TlogStrenghBuildingFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogStrenghBuildingFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogStrenghBuildingFlow()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogStrenghBuildingFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogStrenghBuildingFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogStrenghBuildingFlowById()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogStrenghEquipmentFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64) {
	fi := getMockFunc(s, s.TlogStrenghEquipmentFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogStrenghEquipmentFlow()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
	}

}
func (s *MockTlogService) TlogStrenghEquipmentFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 uint64, a7 uint64, a8 uint64, a9 uint64, a10 uint64) {
	fi := getMockFunc(s, s.TlogStrenghEquipmentFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogStrenghEquipmentFlowById()")
		}
		f(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10)
	}

}
func (s *MockTlogService) TlogTaskFlow(a0 entity.TlogHero, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogTaskFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogTaskFlow()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogTaskFlowById(a0 int64, a1 uint64, a2 uint64, a3 uint64) {
	fi := getMockFunc(s, s.TlogTaskFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogTaskFlowById()")
		}
		f(a0, a1, a2, a3)
	}

}
func (s *MockTlogService) TlogVipLevelFlow(a0 entity.TlogHero, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.TlogVipLevelFlow)
	if fi != nil {
		f, ok := fi.(func(entity.TlogHero, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogVipLevelFlow()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) TlogVipLevelFlowById(a0 int64, a1 uint64, a2 uint64) {
	fi := getMockFunc(s, s.TlogVipLevelFlowById)
	if fi != nil {
		f, ok := fi.(func(int64, uint64, uint64))
		if !ok {
			panic("invalid mock func, MockTlogService.TlogVipLevelFlowById()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockTlogService) WriteLog(a0 string) {
	fi := getMockFunc(s, s.WriteLog)
	if fi != nil {
		f, ok := fi.(func(string))
		if !ok {
			panic("invalid mock func, MockTlogService.WriteLog()")
		}
		f(a0)
	}

}

var TowerModule = &MockTowerModule{}

type MockTowerModule struct{}

func (s *MockTowerModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockTowerModule) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockTowerModule.Close()")
		}
		f()
	}

}

var TssClient = &MockTssClient{}

type MockTssClient struct{}

func (s *MockTssClient) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockTssClient) CallbackAddr() string {
	fi := getMockFunc(s, s.CallbackAddr)
	if fi != nil {
		f, ok := fi.(func() string)
		if !ok {
			panic("invalid mock func, MockTssClient.CallbackAddr()")
		}
		return f()
	}

	return ""
}
func (s *MockTssClient) CheckName(a0 string, a1 bool) (*game2tss.S2CUicJudgeUserInputNameV2Proto, error) {
	fi := getMockFunc(s, s.CheckName)
	if fi != nil {
		f, ok := fi.(func(string, bool) (*game2tss.S2CUicJudgeUserInputNameV2Proto, error))
		if !ok {
			panic("invalid mock func, MockTssClient.CheckName()")
		}
		return f(a0, a1)
	}

	return nil, nil
}
func (s *MockTssClient) Client() *rpc7.Client {
	fi := getMockFunc(s, s.Client)
	if fi != nil {
		f, ok := fi.(func() *rpc7.Client)
		if !ok {
			panic("invalid mock func, MockTssClient.Client()")
		}
		return f()
	}

	return nil
}
func (s *MockTssClient) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockTssClient.Close()")
		}
		f()
	}

}
func (s *MockTssClient) IsEnable() bool {
	fi := getMockFunc(s, s.IsEnable)
	if fi != nil {
		f, ok := fi.(func() bool)
		if !ok {
			panic("invalid mock func, MockTssClient.IsEnable()")
		}
		return f()
	}

	return false
}
func (s *MockTssClient) JudgeChat(a0 *game2tss.C2SUicJudgeUserInputChatV2Proto) (*game2tss.S2CUicJudgeUserInputChatV2Proto, error) {
	fi := getMockFunc(s, s.JudgeChat)
	if fi != nil {
		f, ok := fi.(func(*game2tss.C2SUicJudgeUserInputChatV2Proto) (*game2tss.S2CUicJudgeUserInputChatV2Proto, error))
		if !ok {
			panic("invalid mock func, MockTssClient.JudgeChat()")
		}
		return f(a0)
	}

	return nil, nil
}
func (s *MockTssClient) RegisterCallback(a0 tss.MsgCategory, a1 tss.Callback) {
	fi := getMockFunc(s, s.RegisterCallback)
	if fi != nil {
		f, ok := fi.(func(tss.MsgCategory, tss.Callback))
		if !ok {
			panic("invalid mock func, MockTssClient.RegisterCallback()")
		}
		f(a0, a1)
	}

}
func (s *MockTssClient) TryCheckName(a0 string, a1 sender.Sender, a2 string, a3 pbutil.Buffer, a4 pbutil.Buffer) bool {
	fi := getMockFunc(s, s.TryCheckName)
	if fi != nil {
		f, ok := fi.(func(string, sender.Sender, string, pbutil.Buffer, pbutil.Buffer) bool)
		if !ok {
			panic("invalid mock func, MockTssClient.TryCheckName()")
		}
		return f(a0, a1, a2, a3, a4)
	}

	return false
}

var VipModule = &MockVipModule{}

type MockVipModule struct{}

func (s *MockVipModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var WorldService = &MockWorldService{}

type MockWorldService struct{}

func (s *MockWorldService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

// 广播消息. 同一个[]byte会发送给每一个人. 调用了之后不可以再修改[]byte中的内容
// 消耗大, 会一个个shard锁住user map
func (s *MockWorldService) Broadcast(a0 pbutil.Buffer) {
	fi := getMockFunc(s, s.Broadcast)
	if fi != nil {
		f, ok := fi.(func(pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockWorldService.Broadcast()")
		}
		f(a0)
	}

}
func (s *MockWorldService) BroadcastIgnore(a0 pbutil.Buffer, a1 int64) {
	fi := getMockFunc(s, s.BroadcastIgnore)
	if fi != nil {
		f, ok := fi.(func(pbutil.Buffer, int64))
		if !ok {
			panic("invalid mock func, MockWorldService.BroadcastIgnore()")
		}
		f(a0, a1)
	}

}
func (s *MockWorldService) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockWorldService.Close()")
		}
		f()
	}

}
func (s *MockWorldService) FuncHero(a0 int64, a1 iface.HeroWalker) bool {
	fi := getMockFunc(s, s.FuncHero)
	if fi != nil {
		f, ok := fi.(func(int64, iface.HeroWalker) bool)
		if !ok {
			panic("invalid mock func, MockWorldService.FuncHero()")
		}
		return f(a0, a1)
	}

	return false
}
func (s *MockWorldService) GetTencentInfo(a0 int64) *shared_proto.TencentInfoProto {
	fi := getMockFunc(s, s.GetTencentInfo)
	if fi != nil {
		f, ok := fi.(func(int64) *shared_proto.TencentInfoProto)
		if !ok {
			panic("invalid mock func, MockWorldService.GetTencentInfo()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockWorldService) GetUserCloseSender(a0 int64) sender.ClosableSender {
	fi := getMockFunc(s, s.GetUserCloseSender)
	if fi != nil {
		f, ok := fi.(func(int64) sender.ClosableSender)
		if !ok {
			panic("invalid mock func, MockWorldService.GetUserCloseSender()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockWorldService) GetUserSender(a0 int64) sender.Sender {
	fi := getMockFunc(s, s.GetUserSender)
	if fi != nil {
		f, ok := fi.(func(int64) sender.Sender)
		if !ok {
			panic("invalid mock func, MockWorldService.GetUserSender()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockWorldService) IsDontPush(a0 int64) bool {
	fi := getMockFunc(s, s.IsDontPush)
	if fi != nil {
		f, ok := fi.(func(int64) bool)
		if !ok {
			panic("invalid mock func, MockWorldService.IsDontPush()")
		}
		return f(a0)
	}

	return false
}
func (s *MockWorldService) IsOnline(a0 int64) bool {
	fi := getMockFunc(s, s.IsOnline)
	if fi != nil {
		f, ok := fi.(func(int64) bool)
		if !ok {
			panic("invalid mock func, MockWorldService.IsOnline()")
		}
		return f(a0)
	}

	return false
}
func (s *MockWorldService) MultiSend(a0 []int64, a1 pbutil.Buffer) {
	fi := getMockFunc(s, s.MultiSend)
	if fi != nil {
		f, ok := fi.(func([]int64, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockWorldService.MultiSend()")
		}
		f(a0, a1)
	}

}
func (s *MockWorldService) MultiSendIgnore(a0 []int64, a1 pbutil.Buffer, a2 int64) {
	fi := getMockFunc(s, s.MultiSendIgnore)
	if fi != nil {
		f, ok := fi.(func([]int64, pbutil.Buffer, int64))
		if !ok {
			panic("invalid mock func, MockWorldService.MultiSendIgnore()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockWorldService) MultiSendMsgs(a0 []int64, a1 []pbutil.Buffer) {
	fi := getMockFunc(s, s.MultiSendMsgs)
	if fi != nil {
		f, ok := fi.(func([]int64, []pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockWorldService.MultiSendMsgs()")
		}
		f(a0, a1)
	}

}
func (s *MockWorldService) MultiSendMsgsIgnore(a0 []int64, a1 []pbutil.Buffer, a2 int64) {
	fi := getMockFunc(s, s.MultiSendMsgsIgnore)
	if fi != nil {
		f, ok := fi.(func([]int64, []pbutil.Buffer, int64))
		if !ok {
			panic("invalid mock func, MockWorldService.MultiSendMsgsIgnore()")
		}
		f(a0, a1, a2)
	}

}

// 尝试放入用户, 如果用户已在里面, 则返回旧的用户和false, 如果用户放入成功, 则返回nil, ok
func (s *MockWorldService) PutConnectedUserIfAbsent(a0 iface.ConnectedUser) (iface.ConnectedUser, bool) {
	fi := getMockFunc(s, s.PutConnectedUserIfAbsent)
	if fi != nil {
		f, ok := fi.(func(iface.ConnectedUser) (iface.ConnectedUser, bool))
		if !ok {
			panic("invalid mock func, MockWorldService.PutConnectedUserIfAbsent()")
		}
		return f(a0)
	}

	return nil, false
}

// 删除用户, 如果是同一个对象的话
func (s *MockWorldService) RemoveUserIfSame(a0 iface.ConnectedUser) bool {
	fi := getMockFunc(s, s.RemoveUserIfSame)
	if fi != nil {
		f, ok := fi.(func(iface.ConnectedUser) bool)
		if !ok {
			panic("invalid mock func, MockWorldService.RemoveUserIfSame()")
		}
		return f(a0)
	}

	return false
}
func (s *MockWorldService) Send(a0 int64, a1 pbutil.Buffer) {
	fi := getMockFunc(s, s.Send)
	if fi != nil {
		f, ok := fi.(func(int64, pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockWorldService.Send()")
		}
		f(a0, a1)
	}

}
func (s *MockWorldService) SendFunc(a0 int64, a1 iface.MsgFunc) {
	fi := getMockFunc(s, s.SendFunc)
	if fi != nil {
		f, ok := fi.(func(int64, iface.MsgFunc))
		if !ok {
			panic("invalid mock func, MockWorldService.SendFunc()")
		}
		f(a0, a1)
	}

}
func (s *MockWorldService) SendMsgs(a0 int64, a1 []pbutil.Buffer) {
	fi := getMockFunc(s, s.SendMsgs)
	if fi != nil {
		f, ok := fi.(func(int64, []pbutil.Buffer))
		if !ok {
			panic("invalid mock func, MockWorldService.SendMsgs()")
		}
		f(a0, a1)
	}

}
func (s *MockWorldService) WalkHero(a0 iface.HeroWalker) {
	fi := getMockFunc(s, s.WalkHero)
	if fi != nil {
		f, ok := fi.(func(iface.HeroWalker))
		if !ok {
			panic("invalid mock func, MockWorldService.WalkHero()")
		}
		f(a0)
	}

}
func (s *MockWorldService) WalkUser(a0 iface.UserWalker) {
	fi := getMockFunc(s, s.WalkUser)
	if fi != nil {
		f, ok := fi.(func(iface.UserWalker))
		if !ok {
			panic("invalid mock func, MockWorldService.WalkUser()")
		}
		f(a0)
	}

}

var XiongNuModule = &MockXiongNuModule{}

type MockXiongNuModule struct{}

func (s *MockXiongNuModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

// 启动一个定时任务，刷新
func (s *MockXiongNuModule) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockXiongNuModule.Close()")
		}
		f()
	}

}

// gm 命令开启
func (s *MockXiongNuModule) GmStart(a0 iface.HeroController, a1 int64) bool {
	fi := getMockFunc(s, s.GmStart)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController, int64) bool)
		if !ok {
			panic("invalid mock func, MockXiongNuModule.GmStart()")
		}
		return f(a0, a1)
	}

	return false
}

// 加入联盟后的处理
func (s *MockXiongNuModule) JoinGuild(a0 int64, a1 int64) {
	fi := getMockFunc(s, s.JoinGuild)
	if fi != nil {
		f, ok := fi.(func(int64, int64))
		if !ok {
			panic("invalid mock func, MockXiongNuModule.JoinGuild()")
		}
		f(a0, a1)
	}

}
func (s *MockXiongNuModule) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockXiongNuModule.OnHeroOnline()")
		}
		f(a0)
	}

}

var XiongNuService = &MockXiongNuService{}

type MockXiongNuService struct{}

func (s *MockXiongNuService) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockXiongNuService) AddInfo(a0 xiongnuface.ResistXiongNuInfo) {
	fi := getMockFunc(s, s.AddInfo)
	if fi != nil {
		f, ok := fi.(func(xiongnuface.ResistXiongNuInfo))
		if !ok {
			panic("invalid mock func, MockXiongNuService.AddInfo()")
		}
		f(a0)
	}

}
func (s *MockXiongNuService) GetInfo(a0 int64) xiongnuface.ResistXiongNuInfo {
	fi := getMockFunc(s, s.GetInfo)
	if fi != nil {
		f, ok := fi.(func(int64) xiongnuface.ResistXiongNuInfo)
		if !ok {
			panic("invalid mock func, MockXiongNuService.GetInfo()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockXiongNuService) GetRInfo(a0 int64) xiongnuface.RResistXiongNuInfo {
	fi := getMockFunc(s, s.GetRInfo)
	if fi != nil {
		f, ok := fi.(func(int64) xiongnuface.RResistXiongNuInfo)
		if !ok {
			panic("invalid mock func, MockXiongNuService.GetRInfo()")
		}
		return f(a0)
	}

	return nil
}
func (s *MockXiongNuService) IsStarted(a0 int64) bool {
	fi := getMockFunc(s, s.IsStarted)
	if fi != nil {
		f, ok := fi.(func(int64) bool)
		if !ok {
			panic("invalid mock func, MockXiongNuService.IsStarted()")
		}
		return f(a0)
	}

	return false
}
func (s *MockXiongNuService) IsTodayStarted(a0 int64) bool {
	fi := getMockFunc(s, s.IsTodayStarted)
	if fi != nil {
		f, ok := fi.(func(int64) bool)
		if !ok {
			panic("invalid mock func, MockXiongNuService.IsTodayStarted()")
		}
		return f(a0)
	}

	return false
}
func (s *MockXiongNuService) RemoveInfo(a0 xiongnuface.ResistXiongNuInfo) {
	fi := getMockFunc(s, s.RemoveInfo)
	if fi != nil {
		f, ok := fi.(func(xiongnuface.ResistXiongNuInfo))
		if !ok {
			panic("invalid mock func, MockXiongNuService.RemoveInfo()")
		}
		f(a0)
	}

}
func (s *MockXiongNuService) ResetDaily(a0 time.Time) {
	fi := getMockFunc(s, s.ResetDaily)
	if fi != nil {
		f, ok := fi.(func(time.Time))
		if !ok {
			panic("invalid mock func, MockXiongNuService.ResetDaily()")
		}
		f(a0)
	}

}
func (s *MockXiongNuService) Save() {
	fi := getMockFunc(s, s.Save)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockXiongNuService.Save()")
		}
		f()
	}

}
func (s *MockXiongNuService) SetTodayStarted(a0 int64) {
	fi := getMockFunc(s, s.SetTodayStarted)
	if fi != nil {
		f, ok := fi.(func(int64))
		if !ok {
			panic("invalid mock func, MockXiongNuService.SetTodayStarted()")
		}
		f(a0)
	}

}
func (s *MockXiongNuService) TodayJoinMap() *xiongnuinfo.TodayJoinMap {
	fi := getMockFunc(s, s.TodayJoinMap)
	if fi != nil {
		f, ok := fi.(func() *xiongnuinfo.TodayJoinMap)
		if !ok {
			panic("invalid mock func, MockXiongNuService.TodayJoinMap()")
		}
		return f()
	}

	return nil
}
func (s *MockXiongNuService) WalkInfo(a0 xiongnuface.WalkInfoFunc) {
	fi := getMockFunc(s, s.WalkInfo)
	if fi != nil {
		f, ok := fi.(func(xiongnuface.WalkInfoFunc))
		if !ok {
			panic("invalid mock func, MockXiongNuService.WalkInfo()")
		}
		f(a0)
	}

}
func (s *MockXiongNuService) XiongNuInfoMsg(a0 int64) pbutil.Buffer {
	fi := getMockFunc(s, s.XiongNuInfoMsg)
	if fi != nil {
		f, ok := fi.(func(int64) pbutil.Buffer)
		if !ok {
			panic("invalid mock func, MockXiongNuService.XiongNuInfoMsg()")
		}
		return f(a0)
	}

	return nil
}

var XuanyuanModule = &MockXuanyuanModule{}

type MockXuanyuanModule struct{}

func (s *MockXuanyuanModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

func (s *MockXuanyuanModule) AddChallenger(a0 int64, a1 *shared_proto.CombatPlayerProto, a2 uint64) {
	fi := getMockFunc(s, s.AddChallenger)
	if fi != nil {
		f, ok := fi.(func(int64, *shared_proto.CombatPlayerProto, uint64))
		if !ok {
			panic("invalid mock func, MockXuanyuanModule.AddChallenger()")
		}
		f(a0, a1, a2)
	}

}
func (s *MockXuanyuanModule) Close() {
	fi := getMockFunc(s, s.Close)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockXuanyuanModule.Close()")
		}
		f()
	}

}
func (s *MockXuanyuanModule) GetResetTickTime() tickdata.TickTime {
	fi := getMockFunc(s, s.GetResetTickTime)
	if fi != nil {
		f, ok := fi.(func() tickdata.TickTime)
		if !ok {
			panic("invalid mock func, MockXuanyuanModule.GetResetTickTime()")
		}
		return f()
	}

	return nil
}
func (s *MockXuanyuanModule) GmReset() {
	fi := getMockFunc(s, s.GmReset)
	if fi != nil {
		f, ok := fi.(func())
		if !ok {
			panic("invalid mock func, MockXuanyuanModule.GmReset()")
		}
		f()
	}

}
func (s *MockXuanyuanModule) OnHeroOnline(a0 iface.HeroController) {
	fi := getMockFunc(s, s.OnHeroOnline)
	if fi != nil {
		f, ok := fi.(func(iface.HeroController))
		if !ok {
			panic("invalid mock func, MockXuanyuanModule.OnHeroOnline()")
		}
		f(a0)
	}

}

var ZhanJiangModule = &MockZhanJiangModule{}

type MockZhanJiangModule struct{}

func (s *MockZhanJiangModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}

var ZhengWuModule = &MockZhengWuModule{}

type MockZhengWuModule struct{}

func (s *MockZhengWuModule) Mock(funcKey, funcValue interface{}) {
	Mock(s, funcKey, funcValue)
}
