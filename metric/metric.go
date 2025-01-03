package metric

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	Counter      string = "counter"
	CounterVec   string = "counter_vec"
	Gauge        string = "gauge"
	GaugeVec     string = "gauge_vec"
	Histogram    string = "histogram"
	HistogramVec string = "histogram_vec"
	Summary      string = "summary"
	SummaryVec   string = "summary_vec"
)

type Metric struct {
	Name    string    `toml:"name"`
	Desc    string    `toml:"desc"`
	Type    string    `toml:"type"`
	Args    []string  `toml:"args"`
	Buckets []float64 `toml:"buckets"`
}

func (m *Metric) init(subsystem, namespace string) prometheus.Collector {
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
			Buckets:   m.Buckets,
		})
	case HistogramVec:
		collector = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      m.Name,
			Help:      m.Desc,
			Buckets:   m.Buckets,
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
	return collector
}
