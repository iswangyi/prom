package main

import (
	"awesomeProject1/domain"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
)

const (
	domain1 = "www.baidu2.com"
	//domain2 = "www.17173.com"
)

// DomainCollector  采集器
type DomainCollector struct {
	Namespace       string
	DomainCollector *prometheus.Desc
	Domain          map[string]string
}

// Describe 指标描述
func (c *DomainCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.DomainCollector
}

// Collect 指标信息
func (c *DomainCollector) Collect(ch chan<- prometheus.Metric) {
	ExpiredTime := c.DomainDeadline()
	for k, processCount := range ExpiredTime {
		ch <- prometheus.MustNewConstMetric(
			c.DomainCollector,
			prometheus.GaugeValue,
			processCount,
			k,
		)
	}
}

// DomainDeadline 采集方法
func (c *DomainCollector) DomainDeadline() (processCountByHost map[string]float64) {

	processCountByHost = map[string]float64{
		"www": domain.GetDomainExpired(),
	}
	return processCountByHost
}

// NewClusterManager 创建采集器struct
func NewClusterManager(namespace string) *DomainCollector {
	return &DomainCollector{
		Namespace: namespace,
		DomainCollector: prometheus.NewDesc(
			"domain_deadline_hours",
			" Domain name expired time",
			[]string{"host"},
			prometheus.Labels{
				"namespace": namespace,
				"domain1":   domain1,
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
	log.Fatal(http.ListenAndServe(":8082", nil))
}
