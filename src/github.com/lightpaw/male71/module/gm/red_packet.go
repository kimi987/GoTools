package gm

import (
	"github.com/lightpaw/male7/gen/iface"
	red_packetpb "github.com/lightpaw/male7/gen/pb/red_packet"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/idbytes"
)

func (m *GmModule) newRedPacketGmGroup() *gm_group {
	group := &gm_group{
		tab: "red_packet",
		handler: []*gm_handler{
			newStringHandler("buy", "1", m.redPacketBuy),
			newStringHandler("create", "1", m.redPacketCreate),
			newIntHandler("grab", " ", m.redPacketGrab),
		},
	}
	return group
}

func (m *GmModule) redPacketBuy(dataId string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleRed_packet_c2s_buy(amount string, hc iface.HeroController)
	}); ok {
		im.handleRed_packet_c2s_buy(dataId, hc)
	}

}

func (m *GmModule) redPacketCreate(dataId string, hc iface.HeroController) {
	m.modules.RedPacketModule().(interface {
		ProcessCreate(*red_packetpb.C2SCreateProto, iface.HeroController)
	}).ProcessCreate(&red_packetpb.C2SCreateProto{

		DataId: parseInt32(dataId),

		Count: 2,

		ChatType:shared_proto.ChatType_ChatWorld,

		Text: "XXX",
	}, hc)

}

func (m *GmModule) redPacketGrab(id int64, hc iface.HeroController) {
	m.modules.RedPacketModule().(interface {
		ProcessGrab(*red_packetpb.C2SGrabProto, iface.HeroController)
	}).ProcessGrab(&red_packetpb.C2SGrabProto{
		Id: idbytes.ToBytes(id),
	}, hc)
}
