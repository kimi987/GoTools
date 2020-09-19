package stress

import (
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/stress"
)

//gogen:iface
type StressModule struct {
	isDebug bool
}

func NewStressModule(config *kv.IndividualServerConfig) *StressModule {
	return &StressModule{isDebug: config.IsDebug}
}

//gogen:iface
func (s *StressModule) Ping(proto *stress.C2SRobotPingProto, hc iface.HeroController) {
	if !s.isDebug {
		hc.Disconnect(nil)
		return
	}
	hc.Send(stress.NewS2cRobotPingMsg(proto.Time))
}
