package gm

import (
	"github.com/lightpaw/male7/gen/iface"
	"strconv"
)

func (m *GmModule) newMingcGmGroup() *gm_group {
	return &gm_group{
		tab: "名城",
		handler: []*gm_handler{
			newIntHandler("营建记录", "1", m.mcBuildLog),
			newIntHandler("营建名城", "1", m.mcBuild),
			newIntHandler("名城推荐", "1", m.recommendMcBuild),
		},
	}
}

func (m *GmModule) mcBuildLog(mcId int64, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleMingc_c2s_mc_build_log(amount string, hc iface.HeroController)
	}); ok {
		im.handleMingc_c2s_mc_build_log(strconv.FormatInt(mcId, 10), hc)
	}

}

func (m *GmModule) mcBuild(mcId int64, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleMingc_c2s_mc_build(amount string, hc iface.HeroController)
	}); ok {
		im.handleMingc_c2s_mc_build(strconv.FormatInt(mcId, 10), hc)
	}

}

func (m *GmModule) recommendMcBuild(mcId int64, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleGuild_c2s_add_recommend_mc_build(amount string, hc iface.HeroController)
	}); ok {
		im.handleGuild_c2s_add_recommend_mc_build (strconv.FormatInt(mcId, 10), hc)
	}

}
