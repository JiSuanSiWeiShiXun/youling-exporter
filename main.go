package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"github.com/JiSuanSiWeiShiXun/pcap_exporter/collector"
	"github.com/JiSuanSiWeiShiXun/pcap_exporter/pcap"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	pc   *collector.PcapCollector
	once sync.Once

	exPort string
	filter string //= "udp portrange 4800-4900"
)

func init() {
	// singleton
	once.Do(func() {
		pc = collector.NewPcapCollector()
	})
	flag.StringVar(&exPort, "web.listen-address", ":9356", "Address on which to expose metrics and web interface.")
	flag.StringVar(&filter, "pcap_filter", "", "pcap-filter passed to nethogs.")
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
}

func main() {
	flag.Parse()
	log.Infof("params: [%v]\n", filter)
	go pcap.PacketCapture(pc, filter)
	////pprof
	//go func() {
	//	if err := http.ListenAndServe(":6060", nil); err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	prometheus.MustRegister(pc)
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("start[%v]\n", exPort)
	if err := http.ListenAndServe(exPort, nil); err != nil {
		log.Panic(err)
	}
}
