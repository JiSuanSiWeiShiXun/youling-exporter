package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type (
	NethogsMonitorRecord struct {
		RecordID   int    `json:"record_id" desc:"?"`
		Name       string `json:"name" desc:"？"`
		PID        int    `json:"pid" desc:"进程ID"`
		UID        uint32 `json:"uid" desc:"C uint32"`
		DeviceName string `json:"device_name" desc:"网卡名"`
		SentBytes  uint64 `json:"sent_bytes" desc:"统计到总共发送的数据量(byte)"`
		RecvBytes  uint64 `json:"recv_bytes" desc:"统计到总共的接收数据量(byte)"`
		SentKBs    uint64 `json:"sent_kbs" desc:"float64 delta时间内(1s) 发送的数据量，可以理解为瞬时速率(kb/s)"`
		RecvKBs    uint64 `json:"recv_kbs" desc:"float64 delta时间内(1s) 收到的数据量"`
	}

	NethogsCollector struct {
		RecordMap map[int]*NethogsMonitorRecord
		//isCollected bool //record中存的是临时数据 程序不做缓存，上报一次后就清除
		// 是否会存在一个prometheus pull频率和nethogs update频率差异 导致的数据丢失（中间还有一个collect()的调用频率）
		// 核心疑问还是collect()调用之后，写到只写chan的数据去哪儿了呢
		// desc *prometheus.Desc // desc并不一定要是collector的成员变量
	}
)

var (
	nethogsSentBytesDesc = prometheus.NewDesc( // 可以是个counter
		"sent_bytes_total_nethogs",
		"统计到总共发送的数据量(byte)",
		[]string{"name", "pid", "uid", "device_name", "port"},
		nil)

	nethogsRecvBytesDesc = prometheus.NewDesc(
		"recv_bytes_total_nethogs",
		"统计到总共的接收数据量(byte)",
		[]string{"name", "pid", "uid", "device_name", "port"},
		nil)

	nethogsSentKBsDesc = prometheus.NewDesc(
		"sent_KB_per_second_nethogs",
		"delta时间内(1s) 发送的数据量，可以理解为瞬时速率(kb/s)",
		[]string{"name", "pid", "uid", "device_name", "port"},
		nil)

	nethogsRecvKBsDesc = prometheus.NewDesc(
		"recv_KB_per_second_nethogs",
		"delta时间内(1s) 收到的数据量",
		[]string{"name", "pid", "uid", "device_name", "port"},
		nil)
)

func NewNethogsCollector() *NethogsCollector {
	return &NethogsCollector{
		RecordMap: make(map[int]*NethogsMonitorRecord),
	}
}

func (nc *NethogsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nethogsSentBytesDesc
	ch <- nethogsRecvBytesDesc
	ch <- nethogsSentKBsDesc
	ch <- nethogsRecvKBsDesc
}

func (nc *NethogsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, record := range nc.RecordMap {
		//if _, ok := pidSet[]; !ok {
		//	continue
		//}
		ch <- prometheus.MustNewConstMetric(
			nethogsSentBytesDesc,
			prometheus.CounterValue,
			float64(record.SentBytes),
			record.Name, // {"name", "pid", "uid", "device_name", "port"}
			strconv.Itoa(record.PID),
			strconv.Itoa(int(record.UID)),
			record.DeviceName,
			"[todo]port",
		)
		ch <- prometheus.MustNewConstMetric(
			nethogsRecvBytesDesc,
			prometheus.CounterValue,
			float64(record.RecvBytes),
			record.Name, // {"name", "pid", "uid", "device_name", "port"}
			strconv.Itoa(record.PID),
			strconv.Itoa(int(record.UID)),
			record.DeviceName,
			"[todo]port",
		)
		ch <- prometheus.MustNewConstMetric(
			nethogsSentKBsDesc,
			prometheus.GaugeValue,
			float64(record.SentKBs),
			record.Name, // {"name", "pid", "uid", "device_name", "port"}
			strconv.Itoa(record.PID),
			strconv.Itoa(int(record.UID)),
			record.DeviceName,
			"[todo]port",
		)
		ch <- prometheus.MustNewConstMetric(
			nethogsRecvKBsDesc,
			prometheus.GaugeValue,
			float64(record.RecvKBs),
			record.Name, // {"name", "pid", "uid", "device_name", "port"}
			strconv.Itoa(record.PID),
			strconv.Itoa(int(record.UID)),
			record.DeviceName,
			"[todo]port",
		)
	}
}
