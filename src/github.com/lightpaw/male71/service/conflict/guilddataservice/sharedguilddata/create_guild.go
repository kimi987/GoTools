package sharedguilddata

import (
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/entity/npcid"
	"time"
	"github.com/lightpaw/male7/util/u64"
)

// Npc guild
func NewNpcGuild(template *guild_data.NpcGuildTemplate, newGuildId int64, guildName, flagName string, ctime time.Time, datas *config.ConfigDatas) *Guild {
	// 创建一个Npc帮派
	newGuild := NewGuild(newGuildId, guildName, flagName, datas, ctime)

	// 设置模板
	newGuild.SetNpcTemplate(template)
	newGuild.SetLevelData(template.Level)
	newGuild.SetCountry(template.Country)

	// 设置帮主
	leaderId := npcid.NewNpcMemberId(newGuildId, 0)
	newGuild.SetLeader(leaderId)

	member := newMember(leaderId, datas.GuildConfig().GetLeaderClassLevel(), ctime, nil, nil, template.Leader.EncodeSnapshot(leaderId))
	newGuild.AddMember(member)

	if memberCount := len(template.Contribution7Members()); memberCount > 0 {
		index := 0

		//设置帮派成员
	out:
		for _, classLevel := range datas.GuildConfig().GetNpcSetClassLevelArray() {
			n := template.Level.GetClassMemberCount(classLevel.Level)
			for i := uint64(0); i < n; i++ {

				// 设置官员
				memberId := npcid.NewNpcMemberId(newGuildId, uint64(index+1))
				npcData := template.GetNpc(u64.FromInt(index + 1))
				if npcData == nil {
					break out
				}

				member := newMember(memberId, classLevel, ctime, nil, nil,
					npcData.EncodeSnapshot(memberId))
				newGuild.AddMember(member)

				index++
				if index >= memberCount {
					break out
				}
			}
		}

		if index < memberCount {
			for i := index; i < memberCount; i++ {

				// 设置帮众
				memberId := npcid.NewNpcMemberId(newGuildId, uint64(i+1))
				npcData := template.GetNpc(u64.FromInt(i + 1))
				if npcData == nil {
					break
				}

				member := newMember(memberId, datas.GuildConfig().GetLowestClassLevel(), ctime, nil, nil,
					npcData.EncodeSnapshot(memberId))
				newGuild.AddMember(member)
			}
		}
	}

	// 设置联盟目标
	newGuild.TryUpdateTarget(datas.GuildConfig(), ctime, 0)

	return newGuild

}
