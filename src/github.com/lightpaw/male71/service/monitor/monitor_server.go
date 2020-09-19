package monitor

import (
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/util/call"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/ctxfunc"
	"context"
	"github.com/lightpaw/logrus"
)

func NewGameExporter(register *MetricsRegister, dbService iface.DbService) *GameExporter {

	// 消息处理
	var msgTimeCost *metrics.MsgTimeCostSummary
	if register.EnableMsgMetrics() {
		msgTimeCost = metrics.NewMsgTimeCostSummary()

		// 注册
		register.Register(msgTimeCost.Collector())
	}

	registerCounter := atomic.NewUint64(0)
	ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		maxHeroCount, err := dbService.LoadHeroCount(ctx)
		if err != nil {
			logrus.WithError(err).Panic("获取注册玩家数失败")
		}
		registerCounter.Store(maxHeroCount)
		return
	})

	if register.EnableRegisterCountMetrics() {
		counter := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: metrics.Namespace,
			Subsystem: "register_count",
			Name:      "count",
			Help:      "Register user count.",
		})

		register.Register(counter)
		register.RegisterFunc(func() {
			counter.Set(float64(registerCounter.Load()))
		})
	}

	if register.EnablePanicMetrics() {
		panicCounter := prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: metrics.Namespace,
			Name:      "panic_count",
			Help:      "panic count.",
		})

		register.Register(panicCounter)
		register.RegisterFunc(func() {
			panicCounter.Add(float64(metrics.GetIncPanicCount()))
		})
	}

	g := &GameExporter{
		register:        register,
		msgTimeCost:     msgTimeCost,
		registerCounter: registerCounter,
	}

	return g
}

//gogen:iface
type GameExporter struct {
	register *MetricsRegister

	funcs []metrics.CollectFunc

	// msg time cost
	msgTimeCost *metrics.MsgTimeCostSummary

	registerCounter *atomic.Uint64
}

func (g *GameExporter) GetMsgTimeCost() *metrics.MsgTimeCostSummary {
	return g.msgTimeCost
}

func (g *GameExporter) GetRegisterCounter() *atomic.Uint64 {
	return g.registerCounter
}

func (g *GameExporter) loop() {

	// 定时收集
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:

			// 更新上报指标
			for _, f := range g.funcs {
				call.CatchPanic(f, "定时更新上报指标")
			}
		}
	}
}

func (g *GameExporter) Start() (http.Handler, error) {

	// 注册
	register := prometheus.NewRegistry()

	for _, c := range g.register.collectors {
		if err := register.Register(c); err != nil {
			return nil, errors.Wrapf(err, "注册上报指标失败")
		}
	}

	// 拷贝funcs
	funcs := make([]metrics.CollectFunc, len(g.register.collectFuncs))
	for i, f := range g.register.collectFuncs {
		funcs[i] = f
	}
	g.funcs = funcs
	if len(g.funcs) > 0 {
		go call.CatchLoopPanic(g.loop, "GameExporter")
	}

	return promhttp.HandlerFor(register, promhttp.HandlerOpts{}), nil
}
