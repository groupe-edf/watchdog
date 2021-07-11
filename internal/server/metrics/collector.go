package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "watchdog_"
)

// Collector implements the prometheus.Collector interface
type Collector struct {
	Repositories *prometheus.Desc
}

func (collector Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.Repositories
}

func (collector Collector) Collect(ch chan<- prometheus.Metric) {

}

func NewCollector() Collector {
	return Collector{
		Repositories: prometheus.NewDesc(
			namespace+"repositories",
			"Number of Repositories",
			nil, nil,
		),
	}
}
