package metric

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	id string
	// collector prometheus.Collector

	Name string   `toml:"name"`
	Desc string   `toml:"desc"`
	Type string   `toml:"type"`
	Args []string `toml:"args"`
}

type metrics struct {
	metricType string
	collector  prometheus.Collector
}

func (m *Metric) initMetric(subsystem, namespace string) *metrics {
	var collector prometheus.Collector
	switch m.Type {
	case Counter:
		collector = prometheus.NewCounter(prometheus.CounterOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      m.Name,
			Help:      m.Desc,
		})
	case CounterVec:
		collector = prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      m.Name,
			Help:      m.Desc,
		}, m.Args)
	case Gauge:
		collector = prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      m.Name,
			Help:      m.Desc,
		})
	case GaugeVec:
		collector = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      m.Name,
			Help:      m.Desc,
		}, m.Args)
	case Histogram:
		collector = prometheus.NewHistogram(prometheus.HistogramOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      m.Name,
			Help:      m.Desc,
		})
	case HistogramVec:
		collector = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      m.Name,
			Help:      m.Desc,
		}, m.Args)
	case Summary:
		collector = prometheus.NewSummary(prometheus.SummaryOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      m.Name,
			Help:      m.Desc,
		})
	case SummaryVec:
		collector = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      m.Name,
			Help:      m.Desc,
		}, m.Args)
	default:
		panic(errors.New("type of invalid indicator"))
	}
	return &metrics{
		metricType: m.Type,
		collector:  collector,
	}
}
