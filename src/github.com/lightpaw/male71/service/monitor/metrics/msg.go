package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

func NewMsgTimeCostSummary() *MsgTimeCostSummary {

	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  Namespace,
		Subsystem:  "msg",
		Name:       "time_cost",
		Help:       "msg handle time cost.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
		[]string{"msg"})

	return &MsgTimeCostSummary{
		summaryVec:  summaryVec,
		observerMap: make(map[int]prometheus.Observer),
	}
}

type MsgTimeCostSummary struct {
	// msg time cost
	summaryVec        *prometheus.SummaryVec
	observerMap       map[int]prometheus.Observer
	observerMapLocker sync.RWMutex
}

func (g *MsgTimeCostSummary) Collector() prometheus.Collector {
	return g.summaryVec
}

func (g *MsgTimeCostSummary) GetMsgTimeCostObserver(moduleId, sequenceId int) prometheus.Observer {
	key := moduleId<<16 | sequenceId

	o := g.getCacheMsgTimeCostObserver(key)
	if o != nil {
		return o
	}

	g.observerMapLocker.Lock()
	defer g.observerMapLocker.Unlock()

	o = g.summaryVec.WithLabelValues(fmt.Sprintf("%d-%d", moduleId, sequenceId))
	g.observerMap[key] = o

	return o
}

func (g *MsgTimeCostSummary) getCacheMsgTimeCostObserver(key int) prometheus.Observer {
	g.observerMapLocker.RLock()
	defer g.observerMapLocker.RUnlock()

	return g.observerMap[key]
}
