//package main
//
//import (
//	"awesomeProject1/domain"
//	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//	"github.com/prometheus/common/log"
//	"net/http"
//	"sync"
//)
//
//// DomainCollector  采集器
//type DomainCollector struct {
//	Namespace       string
//	DomainCollector *prometheus.Desc
//	Domain          map[string]string
//}
//
//// Describe 指标描述
//func (c *DomainCollector) Describe(ch chan<- *prometheus.Desc) {
//	ch <- c.DomainCollector
//}
//
//// Collect 指标信息
//func (c *DomainCollector) Collect(ch chan<- prometheus.Metric) {
//	ExpiredTime := c.DomainDeadline()
//	for k, processCount := range ExpiredTime {
//		ch <- prometheus.MustNewConstMetric(
//			c.DomainCollector,
//			prometheus.GaugeValue,
//			processCount,
//			k,
//			k,
//		)
//	}
//}
//
//// DomainDeadline 采集方法
//func (c *DomainCollector) DomainDeadline() (DomainExpired map[string]float64) {
//	//域名列表
//	DomainList := map[string]string{
//		"kledu":   "baidu.com",
//		"kllive":  "17173.com",
//		"kllive2": "cctv.com"}
//	DomainExpired = make(map[string]float64)
//	mutex := sync.Mutex{}
//	for project, domainName := range DomainList {
//		go func(p string, d string) {
//			mutex.Lock()
//			DomainExpired[p] = domain.GetDomainExpired(d)
//			mutex.Unlock()
//		}(project, domainName)
//	}
//
//	return DomainExpired
//}
//
//// NewClusterManager 创建采集器struct
//func NewClusterManager(namespace string) *DomainCollector {
//
//	var c prometheus.Collector
//
//	return &DomainCollector{
//		Namespace: namespace,
//		DomainCollector: prometheus.NewDesc(
//			"domain_deadline_day",
//			" Domain name expired time",
//			[]string{"host", "project"},
//			prometheus.Labels{
//				"namespace": namespace,
//			}),
//	}
//}
//
//func main() {
//	workerA := NewClusterManager("pro")
//	//定期检查收集指标的合法性
//	reg := prometheus.NewPedanticRegistry()
//	//collector注册
//	reg.MustRegister(workerA)
//	//定义采集数据的收集器集合
//	gatherers := prometheus.Gatherers{prometheus.DefaultGatherer, reg}
//	h := promhttp.HandlerFor(gatherers,
//		promhttp.HandlerOpts{
//			ErrorLog:      log.NewErrorLogger(),
//			ErrorHandling: promhttp.ContinueOnError,
//		})
//	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
//		h.ServeHTTP(w, r)
//	})
//	//http.Handle("/metrics", promhttp.Handler())
//	log.Fatal(http.ListenAndServe(":8082", nil))
//}

package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr              = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	uniformDomain     = flag.Float64("uniform.domain", 0.0002, "The domain for the uniform distribution.")
	normDomain        = flag.Float64("normal.domain", 0.0002, "The domain for the normal distribution.")
	normMean          = flag.Float64("normal.mean", 0.00001, "The mean for the normal distribution.")
	oscillationPeriod = flag.Duration("oscillation-period", 10*time.Minute, "The duration of the rate oscillation period.")
)

var (
	// Create a summary to track fictional interservice RPC latencies for three
	// distinct services with different latency distributions. These services are
	// differentiated via a "service" label.
	expired = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "rpc_durations_seconds",
			Help:      "RPC latency distributions.",
			Namespace: "test",
		},
		[]string{"service"},
	)
	// The same as above, but now as a histogram, and only for the normal
	// distribution. The buckets are targeted to the parameters of the
	// normal distribution, with 20 buckets centered on the mean, each
	// half-sigma wide.
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(expired)
}

func main() {
	flag.Parse()

	// Periodically record some sample latencies for the three services.

	go func() {
		for {
			expired.WithLabelValues("day").Inc()
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
