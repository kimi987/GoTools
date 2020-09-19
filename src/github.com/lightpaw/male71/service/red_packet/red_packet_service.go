package red_packet

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/red_packet"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	red_packet_msg "github.com/lightpaw/male7/gen/pb/red_packet"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/must"
	"crypto/rand"
	"math/big"
	"sync"
	"time"
)

func NewRedPacketService(datas iface.ConfigDatas, time iface.TimeService, heroSnapshot iface.HeroSnapshotService, db iface.DbService) *RedPacketService {
	s := &RedPacketService{}
	s.time = time
	s.heroSnapshot = heroSnapshot
	s.db = db
	s.datas = datas

	s.redPacketIdPrefixIncr = atomic.NewInt64(0)
	s.redPackets = make(map[int64]*entity.RedPacket)

	s.load()

	return s
}

//gogen:iface
type RedPacketService struct {
	time         iface.TimeService
	heroSnapshot iface.HeroSnapshotService
	db           iface.DbService
	datas        iface.ConfigDatas

	redPacketIdPrefixIncr *atomic.Int64

	redPackets map[int64]*entity.RedPacket

	mapLock sync.RWMutex
}

func (s *RedPacketService) load() {
	var bytes []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		bytes, err = s.db.LoadKey(ctx, server_proto.Key_RedPacket)
		return
	})
	if err != nil {
		logrus.WithError(err).Panicf("加载红包数据失败")
		return
	}

	if len(bytes) <= 0 {
		return
	}

	proto := &server_proto.AllRedPacketServerProto{}
	if err := proto.Unmarshal(bytes); err != nil {
		logrus.WithError(err).Panicf("解析红包数据失败")
		return
	}

	s.unmarshal(proto)

	logrus.Debugf("加载红包数据成功")
}

func (s *RedPacketService) unmarshal(proto *server_proto.AllRedPacketServerProto) {
	s.redPacketIdPrefixIncr = atomic.NewInt64(proto.CurrIdPrefix)

	ctime := s.time.CurrentTime()
	if proto.RedPackets != nil {
		for _, p := range proto.RedPackets {
			redPacket := entity.NewDefaultRedPacket()
			redPacket.Unmarshal(p, s.datas.RedPacketData())
			if ctime.After(redPacket.ExpiredTime().Add(s.datas.MiscConfig().RedPacketServerDelDuration)) {
				continue
			}
			s.redPackets[redPacket.Id()] = redPacket
		}
	}
}

func (s *RedPacketService) Close() {
	if err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		return s.db.SaveKey(ctx, server_proto.Key_RedPacket, must.Marshal(s.encodeServer()))
	}); err != nil {
		logrus.WithError(err).Error("保存红包数据失败")
		return
	}
	logrus.Debugf("保存红包数据成功")
}

func (s *RedPacketService) encodeServer() *server_proto.AllRedPacketServerProto {
	proto := &server_proto.AllRedPacketServerProto{}
	proto.CurrIdPrefix = s.redPacketIdPrefixIncr.Load()
	for _, r := range s.redPackets {
		proto.RedPackets = append(proto.RedPackets, r.EncodeServer())
	}
	return proto
}

func (s *RedPacketService) createRedPacketId() int64 {
	key, _ := rand.Int(rand.Reader, big.NewInt(0xFFFFFF))
	return s.redPacketIdPrefixIncr.Inc()<<24 + key.Int64()
}

func (s *RedPacketService) Create(heroId int64, data *red_packet.RedPacketData, count uint64, text string, chatType shared_proto.ChatType) (id int64, jsonStr string, errMsg msg.ErrMsg) {
	if !entity.CheckCanCreateRedPacket(data.Money, data.MinPartMoney, count) {
		logrus.Debugf("创建红包，数量错误1")
		errMsg = red_packet_msg.ErrCreateFailCountErr
		return
	}

	ctime := s.time.CurrentTime()
	id = s.createRedPacketId()
	if s.get(id) != nil {
		return
	}

	packet, succ := entity.CreateRedPacket(id, data, count, ctime, text, s.heroSnapshot.GetBasicProto(heroId), chatType)
	if !succ {
		logrus.Debugf("创建红包，数量错误2")
		errMsg = red_packet_msg.ErrCreateFailCountErr
		return
	}

	if succ = s.put(packet); !succ {
		logrus.Debugf("创建红包，红包已经存在:%v", packet.Id())
		errMsg = red_packet_msg.ErrCreateFailServerErr
		return
	}


	if str, err := packet.BuildChatJson(); err != nil {
		logrus.WithError(err).Errorf("创建红包，marshal json 异常")
	} else {
		jsonStr = str
	}

	logrus.Debugf("创建红包:%+v\n", packet)

	return
}

func (s *RedPacketService) put(redPacket *entity.RedPacket) (succ bool) {
	s.mapLock.Lock()
	defer s.mapLock.Unlock()

	if _, ok := s.redPackets[redPacket.Id()]; ok {
		logrus.Errorf("RedPacketService 创建红包时，id 重复.%v", redPacket.Id())
		return
	}

	s.redPackets[redPacket.Id()] = redPacket
	succ = true
	return
}

func (s *RedPacketService) get(id int64) (packet *entity.RedPacket) {
	s.mapLock.RLock()
	defer s.mapLock.RUnlock()

	return s.redPackets[id]
}

func (s *RedPacketService) RedPacketChatId(redPacketId int64) (chatId int64) {
	packet := s.get(redPacketId)
	if packet == nil {
		logrus.Debugf("红包聊天id，红包 id:%v 不存在", redPacketId)
		return
	}
	chatId = packet.ChatId()
	return
}

func (s *RedPacketService) SetRedPacketChatId(redPacketId, chatId int64) {
	packet := s.get(redPacketId)
	if packet == nil {
		logrus.Debugf("设置红包聊天id，红包 id:%v 不存在", redPacketId)
		return
	}
	packet.SetChatId(chatId)
}

func (s *RedPacketService) Grab(redPacketId, heroId, gid int64) (grabMoney uint64, allGrabbed bool, p *shared_proto.RedPacketProto, errMsg msg.ErrMsg) {
	packet := s.get(redPacketId)
	if packet == nil {
		logrus.Debugf("抢红包，红包 id:%v 不存在", redPacketId)
		errMsg = red_packet_msg.ErrGrabFailInvalidId
		return
	}

	if packet.ChatType() == shared_proto.ChatType_ChatGuild {
		if packet.GuildId() != gid {
			errMsg = red_packet_msg.ErrGrabFailNotSameGuild
			return
		}
	}

	if s.Grabbed(redPacketId, heroId) {
		logrus.Debugf("抢红包，红包 id:%v hero:%v 抢过了", redPacketId, heroId)
		p = packet.Encode()
		return
	}

	ctime := s.time.CurrentTime()
	hero := s.heroSnapshot.GetBasicProto(heroId)

	grabMoney, allGrabbed = packet.Grab(heroId, hero, ctime)
	p = packet.Encode()

	return
}

func (s *RedPacketService) Grabbed(redPacketId, heroId int64) (grabbed bool) {
	packet := s.get(redPacketId)
	if packet == nil {
		return
	}
	return packet.Grabbed(heroId)
}

func (s *RedPacketService) AllGrabbed(redPacketId int64) (grabbed bool) {
	packet := s.get(redPacketId)
	if packet == nil {
		return
	}
	return packet.AllGrabbed()
}

func (s *RedPacketService) Expired(redPacketId int64, ctime time.Time) (expired bool) {
	packet := s.get(redPacketId)
	if packet == nil {
		return
	}
	return ctime.After(packet.ExpiredTime())
}

func (s *RedPacketService) Exist(redPacketId int64) bool {
	return s.get(redPacketId) != nil
}
