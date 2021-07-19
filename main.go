package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"math/rand"
	"net/http"
)

type ClusterManager struct {
	//集群名称
	Namespace          string
	ProcessCounterDesc *prometheus.Desc
}

// Describe 指标描述
func (c *ClusterManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ProcessCounterDesc
}

// Collect 指标信息
func (c *ClusterManager) Collect(ch chan<- prometheus.Metric) {
	SystemProcessBycount := c.SystemState()
	for k, processCount := range SystemProcessBycount {
		ch <- prometheus.MustNewConstMetric(
			c.ProcessCounterDesc,
			prometheus.GaugeValue,
			processCount,
			k,
			k,
		)
	}
}

// SystemState 采集方法
func (c *ClusterManager) SystemState() (processCountByHost map[string]float64) {
	processCountByHost = map[string]float64{
		"193": float64(rand.Int31n(1000)),
		"194": float64(rand.Int31n(1000)),
	}
	return processCountByHost
}

// NewClusterManager 创建采集器struct
func NewClusterManager(namespace string) *ClusterManager {
	return &ClusterManager{
		Namespace: namespace,
		ProcessCounterDesc: prometheus.NewDesc(
			"clustermanager_process_total",
			"Number of restart process",
			[]string{"host", "ip"},
			prometheus.Labels{
				"namespace": namespace,
				"domain":    "www.baidu.com",
			}),
	}
}

func main() {
	workerA := NewClusterManager("test")
	//定期检查收集指标的合法性
	reg := prometheus.NewPedanticRegistry()
	//collector注册
	reg.MustRegister(workerA)
	//定义采集数据的收集器集合
	gatherers := prometheus.Gatherers{prometheus.DefaultGatherer, reg}
	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
	//http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8081", nil))
}
