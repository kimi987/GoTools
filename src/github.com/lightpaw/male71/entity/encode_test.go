package entity

import (
	"fmt"
	"testing"
)

type test_proto struct {
	Id      int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Head    string `protobuf:"bytes,10,opt,name=head,proto3" json:"head,omitempty"`
	BigHead string `protobuf:"bytes,11,opt,name=big_head,json=bigHead,proto3" json:"big_head,omitempty"`
	Exp     uint64 `protobuf:"varint,13,opt,name=exp,proto3" json:"exp,omitempty"`
	Male    bool   `protobuf:"varint,14,opt,name=male,proto3" json:"male,omitempty"`
	// 改名
	NextChangeNameTime  int64  `protobuf:"varint,15,opt,name=next_change_name_time,json=nextChangeNameTime,proto3" json:"next_change_name_time,omitempty"`
	ChangeHeroNameTimes uint64 `protobuf:"varint,18,opt,name=change_hero_name_times,json=changeHeroNameTimes,proto3" json:"change_hero_name_times,omitempty"`
	DailyResetTime      int64  `protobuf:"varint,16,opt,name=daily_reset_time,json=dailyResetTime,proto3" json:"daily_reset_time,omitempty"`
	// 帮派
	GuildId       int64  `protobuf:"varint,31,opt,name=guild_id,json=guildId,proto3" json:"guild_id,omitempty"`
	GuildName     string `protobuf:"bytes,32,opt,name=guild_name,json=guildName,proto3" json:"guild_name,omitempty"`
	GuildFlagName string `protobuf:"bytes,33,opt,name=guild_flag_name,json=guildFlagName,proto3" json:"guild_flag_name,omitempty"`
	// 元宝
	Yuanbao uint64 `protobuf:"varint,41,opt,name=yuanbao,proto3" json:"yuanbao,omitempty"`
	// 资源
	Gold  uint64 `protobuf:"varint,42,opt,name=gold,proto3" json:"gold,omitempty"`
	Food  uint64 `protobuf:"varint,43,opt,name=food,proto3" json:"food,omitempty"`
	Wood  uint64 `protobuf:"varint,44,opt,name=wood,proto3" json:"wood,omitempty"`
	Stone uint64 `protobuf:"varint,45,opt,name=stone,proto3" json:"stone,omitempty"`
}

func Benchmark_Plain_Set(b *testing.B) {
	var proto *test_proto
	for i := 0; i < b.N; i++ {
		proto = &test_proto{}
		proto.Id = 0
		proto.Name = ""
		proto.Head = ""
		proto.Exp = 0
		proto.Male = true

		proto.NextChangeNameTime = 0
		proto.ChangeHeroNameTimes = 0
		proto.DailyResetTime = 0

		proto.GuildId = 0
		proto.GuildName = ""
		proto.GuildFlagName = ""

		// 元宝
		proto.Yuanbao = 0

		proto.Gold = 0
		proto.Food = 0
		proto.Wood = 0
		proto.Stone = 0
	}
	if proto.Id == 1 {
		fmt.Println(proto)
	}
}

func Benchmark_Constructor_Set(b *testing.B) {
	var proto *test_proto
	for i := 0; i < b.N; i++ {
		proto = &test_proto{

			Id:   0,
			Name: "",
			Head: "",
			Exp:  0,
			Male: true,

			NextChangeNameTime:  0,
			ChangeHeroNameTimes: 0,
			DailyResetTime:      0,

			GuildId:       0,
			GuildName:     "",
			GuildFlagName: "",

			// 元宝
			Yuanbao: 0,

			Gold:  0,
			Food:  0,
			Wood:  0,
			Stone: 0,
		}
	}

	if proto.Id == 1 {
		fmt.Println(proto)
	}
}
