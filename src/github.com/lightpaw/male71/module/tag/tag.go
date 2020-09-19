package tag

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/tag"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/lock"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
)

func NewTagModule(configDatas iface.ConfigDatas, timeService iface.TimeService, heroDataService iface.HeroDataService, guildSnapshotService iface.GuildSnapshotService, worldService iface.WorldService) *TagModule {
	return &TagModule{
		configDatas:          configDatas,
		timeService:          timeService,
		heroDataService:      heroDataService,
		guildSnapshotService: guildSnapshotService,
		worldService:         worldService,
	}
}

//gogen:iface
type TagModule struct {
	configDatas          iface.ConfigDatas
	timeService          iface.TimeService
	heroDataService      iface.HeroDataService
	guildSnapshotService iface.GuildSnapshotService
	worldService         iface.WorldService
}

//gogen:iface
func (m *TagModule) ProcessAddOrUpdateTag(proto *tag.C2SAddOrUpdateTagProto, hc iface.HeroController) {
	if len(proto.GetTag()) <= 0 {
		logrus.WithField("tag", proto.GetTag()).Debugln("标签非法")
		hc.Send(tag.ERR_ADD_OR_UPDATE_TAG_FAIL_CONTENT_TOO_SHORT)
		return
	}

	if u64.FromInt(util.GetCharLen(proto.GetTag())) > m.configDatas.TagMiscData().MaxCharCount {
		logrus.WithField("tag", proto.GetTag()).Debugln("标签非法")
		hc.Send(tag.ERR_ADD_OR_UPDATE_TAG_FAIL_CONTENT_TOO_LONG)
		return
	}

	targetId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.WithField("id", proto.GetId()).Debugln("解析玩家id没解析出来")
		hc.Send(tag.ERR_ADD_OR_UPDATE_TAG_FAIL_TARGET_NOT_FOUND)
		return
	}

	var name string
	var guildId int64

	if hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		name = hero.Name()
		guildId = hero.GuildId()
		return
	}) {
		hc.Send(tag.ERR_ADD_OR_UPDATE_TAG_FAIL_TARGET_NOT_FOUND)
		return
	}

	var guildFlagName string
	if guildId != 0 {
		snapshot := m.guildSnapshotService.GetSnapshot(hc.Id())
		if snapshot != nil {
			guildFlagName = snapshot.FlagName
		}
	}

	var errMsg pbutil.Buffer

	var t *shared_proto.TagProto
	var record *shared_proto.TagRecordProto

	m.heroDataService.Func(targetId, func(hero *entity.Hero, err error) (heroChanged bool) {
		if err != nil {
			if err == lock.ErrEmpty {
				// 玩家不存在
				logrus.WithField("id", proto.GetId()).Debugln("玩家不存在")
				errMsg = tag.ERR_ADD_OR_UPDATE_TAG_FAIL_TARGET_NOT_FOUND
				return
			}

			logrus.WithField("id", proto.GetId()).Errorln("服务器错误")
			errMsg = tag.ERR_ADD_OR_UPDATE_TAG_FAIL_TARGET_NOT_FOUND
			return
		}

		heroTag := hero.Tag()

		if !heroTag.Exist(proto.GetTag()) {
			// 标签不存在，那就是新增咯
			if heroTag.TagCount() >= int(m.configDatas.TagMiscData().MaxCount) {
				logrus.WithField("tag count", heroTag.TagCount()).Debugln("标签数量已经满了")
				errMsg = tag.ERR_ADD_OR_UPDATE_TAG_FAIL_TARGET_TAG_FULL
				return
			}
		}

		// TODO 检验黑名单

		t, record = heroTag.AddTag(hc.IdBytes(), name, guildFlagName, proto.GetTag(), m.timeService.CurrentTime())

		return true
	})

	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	recordBytes := must.Marshal(record)
	tagBytes := must.Marshal(t)
	if targetId != hc.Id() {
		// 不是自己给自己加的
		m.worldService.Send(targetId, tag.NewS2cOtherTagMeMsg(recordBytes, tagBytes))
	}

	hc.Send(tag.NewS2cAddOrUpdateTagMsg(proto.GetId(), recordBytes, tagBytes))
}

var deleteEmptyTagsSuc = tag.NewS2cDeleteTagMsg([]string{}).Static()

//gogen:iface
func (m *TagModule) ProcessDeleteTag(proto *tag.C2SDeleteTagProto, hc iface.HeroController) {
	if len(proto.Tags) <= 0 {
		logrus.WithField("tag count", len(proto.Tags)).Debugln("标签起码要发一个")
		hc.Send(deleteEmptyTagsSuc)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroTag := hero.Tag()

		for _, t := range proto.Tags {
			heroTag.RemoveTag(t)
		}

		result.Add(tag.NewS2cDeleteTagMsg(proto.Tags))
		result.Changed()
		result.Ok()
	})
}
