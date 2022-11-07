package collector

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	myMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "what's namespace?",
		Name:      "DIYMetric",
		Help:      "Is is comments?",
	})
)

func TestCollector(t *testing.T) {
	prometheus.MustRegister(myMetric)

}
