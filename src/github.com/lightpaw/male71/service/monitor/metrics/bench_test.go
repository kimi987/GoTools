package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

func BenchmarkTimer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer := prometheus.NewTimer(prometheus.ObserverFunc(empty))
		timer.ObserveDuration()
	}
}

func empty(s float64) {

}
