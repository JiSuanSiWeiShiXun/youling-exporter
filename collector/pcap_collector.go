package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	pcapNetIOSendTotalDesc = prometheus.NewDesc(
		"pcap_net_io_send_bytes_total",
		"总发送流量(bytes)",
		[]string{"groupname", "protocol"}, // groupname[ip:port], protocol[TCP\UDP]
		map[string]string{"source": "pcap"},
	)

	pcapNetIORecvTotalDesc = prometheus.NewDesc(
		"pcap_net_io_recv_bytes_total",
		"总接收流量(bytes)",
		[]string{"groupname", "protocol"},
		map[string]string{"source": "pcap"},
	)
)

type (
	// 控制 [prometheus scrape -> dataCache() update -> Collect() return]过程 的同步逻辑
	scrapeRequest struct {
		results chan<- prometheus.Metric
		done    chan struct{}
	}

	NetRecord struct {
		Name           string // <ip:port>
		Protocol       string
		SendBytesTotal uint64
		RecvBytesTotal uint64
	}

	PcapCollector struct {
		NetRecordCache map[string]NetRecord

		scrapeChan chan scrapeRequest
	}
)

// NewPcapCollector initialize a PcapCollector
func NewPcapCollector() *PcapCollector {
	pc := &PcapCollector{
		NetRecordCache: make(map[string]NetRecord),

		scrapeChan: make(chan scrapeRequest),
	}
	go pc.scrape()
	return pc
}

func (pc *PcapCollector) scrape() {
	for req := range pc.scrapeChan {
		// 每收到一次prometheus scrape数据的请求
		for _, r := range pc.NetRecordCache {
			req.results <- prometheus.MustNewConstMetric(
				pcapNetIOSendTotalDesc,
				prometheus.CounterValue,
				float64(r.SendBytesTotal),
				r.Name,
				r.Protocol,
			)
			req.results <- prometheus.MustNewConstMetric(
				pcapNetIORecvTotalDesc,
				prometheus.CounterValue,
				float64(r.RecvBytesTotal),
				r.Name,
				r.Protocol,
			)
		}
		req.done <- struct{}{}
	}
}

func (pc *PcapCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- pcapNetIOSendTotalDesc
	ch <- pcapNetIORecvTotalDesc
}

func (pc *PcapCollector) Collect(ch chan<- prometheus.Metric) {
	log.Trace("here comes a prometheus scrape")
	req := scrapeRequest{results: ch, done: make(chan struct{})}
	pc.scrapeChan <- req
	<-req.done
}
