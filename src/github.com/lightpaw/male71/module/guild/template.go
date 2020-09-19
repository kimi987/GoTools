package guild

import (
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/util/imath"
)

func newGuildTemplateArray(templates []*guild_data.NpcGuildTemplate, guilds []*sharedguilddata.Guild) []*guild_template {

	array := make([]*guild_template, 0, len(templates))
	for _, template := range templates {
		if template.RejectUserJoin {
			continue
		}

		array = append(array, newGuildTemplate(template, guilds))
	}

	return array
}

func newGuildTemplate(template *guild_data.NpcGuildTemplate, guilds []*sharedguilddata.Guild) *guild_template {
	t := &guild_template{}
	t.template = template

	names := template.GetCombineNames()
	flagNames := template.GetCombineFlagNames()

	n := imath.Min(len(names), len(flagNames))

out:
	for i := 0; i < n; i++ {
		name := names[i]
		flagName := flagNames[i]

		for _, g := range guilds {
			if name == g.Name() || flagName == g.FlagName() {
				continue out
			}
		}

		t.names = append(t.names, name)
		t.flagNames = append(t.flagNames, flagName)
	}

	return t
}

type guild_template struct {
	// 配置模板
	template *guild_data.NpcGuildTemplate

	// 允许使用的名字以及旗号
	names     []string
	flagNames []string
}

func (t *guild_template) Template() *guild_data.NpcGuildTemplate {
	return t.template
}

func (t *guild_template) HasName() bool {
	return len(t.names) > 0
}

func (t *guild_template) PopName() (guildName, flagName string) {

	if len(t.names) > 0 {
		guildName = t.names[0]
		copy(t.names, t.names[1:])
		t.names = t.names[:len(t.names)-1]
	}

	if len(t.flagNames) > 0 {
		flagName = t.flagNames[0]
		copy(t.flagNames, t.flagNames[1:])
		t.flagNames = t.flagNames[:len(t.flagNames)-1]
	}
	return
}
