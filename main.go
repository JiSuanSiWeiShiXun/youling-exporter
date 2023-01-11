package main

import "C"
import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"nethogs-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	nethogsCollector   *collector.NethogsCollector
	monitorRecordMutex sync.Mutex
	once               sync.Once

	exPort string
	// procPorts string
	strParams  string
	scrapeChan chan *collector.NethogsMonitorRecord
)

func init() {
	// singleton
	once.Do(func() {
		nethogsCollector = collector.NewNethogsCollector()
		scrapeChan = make(chan *collector.NethogsMonitorRecord, 100) // 长度暂定100
	})
	flag.StringVar(&exPort, "web.listen-address", ":9356", "Address on which to expose metrics and web interface.")
	// flag.StringVar(&procPorts, "procports", "", "comma-separated list, port numbers to monitor.") // 为空就是所有进程
	flag.StringVar(&strParams, "pcap_filter", "", "pcap-filter passed to nethogs.")
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
}

// Scrape 从channel中不停读数据，
func Scrape() {
	log.Debug("starts scrape goroutine.")
	for {
		select {
		case record := <-scrapeChan:
			// 从channel中读出数据结构赋值给map

			// MonitorRecordMap数据的读写锁
			//monitorRecordMutex.Lock()
			//defer func() {
			//	monitorRecordMutex.Unlock()
			//}()
			key := fmt.Sprintf("%v | pid[%v]", record.Name, record.PID)
			if _, ok := nethogsCollector.RecordMap[key]; !ok {
				nethogsCollector.RecordMap[key] = new(collector.NethogsMonitorRecord)
			}
			nethogsCollector.RecordMap[key] = record

			log.Debugf("[tick]\n[key]\"%v\" [Name]%v [PID]%v [UID]%v [DEV]%v [sdbt]%v [rcbt]%v [sdkb]%v [rckb]%v\n[endtick]\n\n",
				key,
				nethogsCollector.RecordMap[key].Name,
				nethogsCollector.RecordMap[key].PID,
				nethogsCollector.RecordMap[key].UID,
				nethogsCollector.RecordMap[key].DeviceName,
				nethogsCollector.RecordMap[key].SentBytes,
				nethogsCollector.RecordMap[key].RecvBytes,
				nethogsCollector.RecordMap[key].SentKBs,
				nethogsCollector.RecordMap[key].RecvKBs,
			)
			time.Sleep(time.Millisecond*10) // 控制取消息频率上限为1s 100次，也就是可以1s处理100个进程的流量数据
		}
	}
}

func main() {
	// flag获取指定端口
	//webPort := 6789
	//targetPortList := []int{}
	flag.Parse()

	//context
	log.Infof("params: [%v]\n", strParams)
	go CallNethogs(strParams)
	go Scrape()
	//pprof
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatal(err)
		}
	}()

	prometheus.MustRegister(nethogsCollector)
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("start[%v]\n", exPort)
	if err := http.ListenAndServe(exPort, nil); err != nil {
		panic(err)
	}
}
