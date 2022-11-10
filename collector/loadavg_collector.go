package collector

import "github.com/prometheus/client_golang/prometheus"

var nameSpace = "node"

type (
	NethogsCollector struct { // counter, gauge, histogram, summary
		metrics []NethogsDesc
	}

	NethogsDesc struct {
		desc      *prometheus.Desc     // 为什么要在collector中实现一个 desc？能否直接作为collector成员变量？
		valueType prometheus.ValueType // ?
	}
)

func NewNethogsCollector() *NethogsCollector {
	return &NethogsCollector{
		metrics: []NethogsDesc{
			{prometheus.NewDesc(nameSpace+"_load1", "1m load avg", nil, nil), prometheus.GaugeValue},
			{prometheus.NewDesc(nameSpace+"_load1", "1m load avg", nil, nil), prometheus.GaugeValue},
			{prometheus.NewDesc(nameSpace+"_load1", "1m load avg", nil, nil), prometheus.GaugeValue},
		},
	}
}

func (nc *NethogsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.metrics[0].desc
}

func (nc *NethogsCollector) Collect(ch chan<- prometheus.Metric) {
	//ch <- nc.metrics[0]
}
