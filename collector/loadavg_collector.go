package collector

import "github.com/prometheus/client_golang/prometheus"

var nameSpace = "node"

type (
	LoadavgCollector struct { // counter, gauge, histogram, summary
		metrics []LoadavgDesc
	}

	LoadavgDesc struct {
		desc      *prometheus.Desc     // 为什么要在collector中实现一个 desc？能否直接作为collector成员变量？
		valueType prometheus.ValueType // ?
	}
)

func NewLoadavgCollector() *LoadavgCollector {
	return &LoadavgCollector{
		metrics: []LoadavgDesc{
			{prometheus.NewDesc(nameSpace+"_load1", "1m load avg", nil, nil), prometheus.GaugeValue},
			{prometheus.NewDesc(nameSpace+"_load1", "1m load avg", nil, nil), prometheus.GaugeValue},
			{prometheus.NewDesc(nameSpace+"_load1", "1m load avg", nil, nil), prometheus.GaugeValue},
		},
	}
}

func (nc *LoadavgCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.metrics[0].desc
}

func (nc *LoadavgCollector) Collect(ch chan<- prometheus.Metric) {
	//ch <- nc.metrics[0]
}
