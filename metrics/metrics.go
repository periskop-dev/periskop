package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// PrometheusNamespace namespace for Prometheus metrics
	PrometheusNamespace = "periskop"
)

// nolint
var (
	scrappedLabels = []string{"service_name"}
	// InstancesScrapped is a Prometheus gauge to track the number of instances scrapped
	InstancesScrapped = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: PrometheusNamespace,
			Name:      "instances_scrapped",
			Help:      "Number of instances scrapped.",
		},
		scrappedLabels,
	)
	// ErrorsScrapped is a Prometheus counter to track the total of errors scrapped
	ErrorsScrapped = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: PrometheusNamespace,
			Name:      "errors_scrapped_total",
			Help:      "Total number of errors scrapped.",
		},
		scrappedLabels,
	)
	// ServiceErrors is a Prometheus counter to track errors in the Periskop service
	ServiceErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: PrometheusNamespace,
			Name:      "application_errors_total",
			Help:      "Nu.",
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(InstancesScrapped)
	prometheus.MustRegister(ErrorsScrapped)
	prometheus.MustRegister(ServiceErrors)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}
