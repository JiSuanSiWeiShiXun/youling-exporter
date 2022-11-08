package collector

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	myMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "namespace_means_prefix",
		Name:      "DIY_Metric",
		Help:      "this word will be listed in comment section",
		ConstLabels: map[string]string{
			"key_name": "youling",
		},
	})
)

func recordMetrics(c context.Context) {
	// 这个collector只写了30条数据，一秒一条，30s后这个goroutine会被回收
	go func() {
		for {
			select {
			case <-c.Done():
				return
			default:
				myMetric.Inc()
				time.Sleep(time.Millisecond * 1000)
			}
		}
	}()
}

func TestCollector(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	prometheus.MustRegister(myMetric)
	recordMetrics(ctx)

	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("start[6789]\n")
	if err := http.ListenAndServe(":6789", nil); err != nil {
		panic(err)
	}
}
