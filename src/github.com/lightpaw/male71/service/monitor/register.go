package monitor

import (
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func NewMetricsRegister(serverConfig *kv.IndividualServerConfig) *MetricsRegister {
	return &MetricsRegister{
		serverConfig: serverConfig,
	}
}

//gogen:iface
type MetricsRegister struct {
	serverConfig *kv.IndividualServerConfig

	collectors []prometheus.Collector

	collectFuncs []metrics.CollectFunc
}

func (r *MetricsRegister) EnableDBMetrics() bool {
	return r.serverConfig.EnableDBMetrics
}

func (r *MetricsRegister) EnableOnlineCountMetrics() bool {
	return r.serverConfig.EnableOnlineCountMetrics
}

func (r *MetricsRegister) EnableRegisterCountMetrics() bool {
	return r.serverConfig.EnableRegisterMetrics
}

func (r *MetricsRegister) EnableMsgMetrics() bool {
	return r.serverConfig.EnableMsgMetrics
}

func (r *MetricsRegister) EnablePanicMetrics() bool {
	return r.serverConfig.EnablePanicMetrics
}

func (r *MetricsRegister) Register(collector prometheus.Collector) {
	r.collectors = append(r.collectors, collector)
}

func (r *MetricsRegister) RegisterFunc(collectFunc metrics.CollectFunc) {
	r.collectFuncs = append(r.collectFuncs, collectFunc)
}
