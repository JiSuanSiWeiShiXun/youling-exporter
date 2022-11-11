package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
	"youling-exporter/collector"
)

var (
	nethogsCollector   *collector.NethogsCollector
	monitorRecordMutex sync.Mutex
	once               sync.Once
)

func init() {
	// singleton
	once.Do(func() {
		nethogsCollector = collector.NewNethogsCollector()
	})
}

func main() {
	// flag获取指定端口
	//webPort := 6789
	//targetPortList := []int{}

	//context
	go CallNethogs()

	prometheus.MustRegister(nethogsCollector)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("start[6789]\n")
	if err := http.ListenAndServe(":6789", nil); err != nil {
		panic(err)
	}
}
