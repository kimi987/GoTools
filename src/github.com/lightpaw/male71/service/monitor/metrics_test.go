package monitor

import (
	"fmt"
	"github.com/lightpaw/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"testing"
	"time"
)

func TestHistogram(t *testing.T) {
	temps := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "pond_temperature_celsius",
		Help:    "The temperature of the frog pond.", // Sorry, we can't measure how badly it smells.
		Buckets: prometheus.LinearBuckets(20, 5, 5),  // 5 buckets, each 5 centigrade wide.
	})

	// Simulate some observations.
	//for i := 0; i < 1000; i++ {
	//	temps.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
	//}

	temps.Observe(41)

	// Just for demonstration, let's check the state of the histogram by
	// (ab)using its Write method (which is usually only used by Prometheus
	// internally).
	metric := &dto.Metric{}
	temps.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))

	metric = &dto.Metric{}
	temps.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))
}

func TestSummary(t *testing.T) {
	temps := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "pond_temperature_celsius",
		Help:       "The temperature of the frog pond.",
		Objectives: map[float64]float64{0.5: 1000, 0.9: 0.1, 0.99: 0.01},
	})

	// Simulate some observations.
	//for i := 0; i < 1000; i++ {
	//	temps.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
	//}

	for i := 0; i < 10; i++ {
		temps.Observe(float64(30 + i))
	}

	// Just for demonstration, let's check the state of the summary by
	// (ab)using its Write method (which is usually only used by Prometheus
	// internally).
	metric := &dto.Metric{}
	temps.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))

	// Just for demonstration, let's check the state of the summary by
	// (ab)using its Write method (which is usually only used by Prometheus
	// internally).
	metric = &dto.Metric{}
	temps.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))
}

func TestSummary0(t *testing.T) {
	temps := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "pond_temperature_celsius",
		Help:       "The temperature of the frog pond.",
		Objectives: map[float64]float64{},
	})

	// Simulate some observations.
	//for i := 0; i < 1000; i++ {
	//	temps.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
	//}

	for i := 0; i < 10; i++ {
		temps.Observe(float64(30 + i))
	}

	// Just for demonstration, let's check the state of the summary by
	// (ab)using its Write method (which is usually only used by Prometheus
	// internally).
	metric := &dto.Metric{}
	temps.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))
}

func TestCounter(t *testing.T) {
	pushCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "repository_pushes", // Note: No help string...
	})

	// 10 operations queued by the goroutine managing incoming requests.
	pushCounter.Add(10)
	// A worker goroutine has picked up a waiting operation.
	pushCounter.Inc()
	// And once more...
	pushCounter.Inc()

	metric := &dto.Metric{}
	pushCounter.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))

	metric = &dto.Metric{}
	pushCounter.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))
}

func TestGauge(t *testing.T) {
	opsQueued := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "our_company",
		Subsystem: "blob_storage",
		Name:      "ops_queued",
		Help:      "Number of blob storage operations waiting to be processed.",
	})

	// 10 operations queued by the goroutine managing incoming requests.
	opsQueued.Add(10)
	// A worker goroutine has picked up a waiting operation.
	opsQueued.Dec()
	// And once more...
	opsQueued.Dec()

	metric := &dto.Metric{}
	opsQueued.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))

	metric = &dto.Metric{}
	opsQueued.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))
}

func TestTimer(t *testing.T) {
	funcDuration := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "example_function_duration_seconds",
		Help: "Duration of the last call of an example function.",
	})

	// The Set method of the Gauge is used to observe the duration.
	timer := prometheus.NewTimer(prometheus.ObserverFunc(funcDuration.Set))

	time.Sleep(10 * time.Millisecond)

	timer.ObserveDuration()

	metric := &dto.Metric{}
	funcDuration.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))

	timer = prometheus.NewTimer(prometheus.ObserverFunc(funcDuration.Set))

	time.Sleep(10 * time.Millisecond)

	timer.ObserveDuration()

	metric = &dto.Metric{}
	funcDuration.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))
}
