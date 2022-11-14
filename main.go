package main

import (
	"fmt"
    "flag"
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

    exPort string
    // procPorts string
    strParams string
)

func init() {
	// singleton
	once.Do(func() {
		nethogsCollector = collector.NewNethogsCollector()
	})
    flag.StringVar(&exPort, "web.listen-address", ":9356", "Address on which to expose metrics and web interface.")
    // flag.StringVar(&procPorts, "procports", "", "comma-separated list, port numbers to monitor.") // 为空就是所有进程
    flag.StringVar(&strParams, "pcap_filter", "", "pcap-filter passed to nethogs.")
}

func main() {
	// flag获取指定端口
	//webPort := 6789
	//targetPortList := []int{}
        flag.Parse()

	//context
        fmt.Printf("params: [%v]\n", strParams)
	go CallNethogs(strParams)

	prometheus.MustRegister(nethogsCollector)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("start[%v]\n", exPort)
	if err := http.ListenAndServe(exPort, nil); err != nil {
            panic(err)
	}
}
