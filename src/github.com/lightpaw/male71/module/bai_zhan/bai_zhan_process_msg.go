package bai_zhan

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/bai_zhan"
)

//gogen:iface c2s_bai_zhan_challenge
func (m *BaiZhanModule) ProcessChallenge(hc iface.HeroController) {
	err := m.Challenge(hc)
	if err != nil {
		logrus.WithError(err).Debugln(err.Error())
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_collect_salary
func (m *BaiZhanModule) ProcessCollectSalary(hc iface.HeroController) {
	_, err := m.CollectSalary(hc)
	if err != nil {
		logrus.WithError(err).Debugln(err.Error())
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface
func (m *BaiZhanModule) ProcessCollectJunXianPrize(proto *bai_zhan.C2SCollectJunXianPrizeProto, hc iface.HeroController) {
	_, err := m.CollectJunXianPrize(proto.GetId(), hc)
	if err != nil {
		logrus.WithError(err).Debugln(err.Error())
		hc.Send(err.ErrMsg())
		return
	}
}

//gogen:iface c2s_self_record
func (m *BaiZhanModule) ProcessSelfRecord(proto *bai_zhan.C2SSelfRecordProto, hc iface.HeroController) {
	m.SelfRecord(proto.GetVersion(), hc)
}

//gogen:iface c2s_query_bai_zhan_info
func (m *BaiZhanModule) ProcessQueryBaiZhanInfo(hc iface.HeroController) {
	m.QueryBaiZhanInfo(hc)
}

//gogen:iface c2s_clear_last_jun_xian
func (m *BaiZhanModule) ProcesClearLastJunXian(hc iface.HeroController) {
	m.ClearLastJunXian(hc)
	hc.Send(bai_zhan.CLEAR_LAST_JUN_XIAN_S2C)
}

//gogen:iface c2s_request_self_rank
func (m *BaiZhanModule) ProcessRequestSelfRank(hc iface.HeroController) {
	m.RequestSelfRank(hc)
}

//gogen:iface c2s_request_rank
func (m *BaiZhanModule) ProcessRequestRank(proto *bai_zhan.C2SRequestRankProto, hc iface.HeroController) {
	m.RequestRank(proto.Self, uint64(proto.StartRank), hc)
}
