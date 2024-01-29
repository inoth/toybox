package metric

import (
	"errors"

	"github.com/inoth/toybox/util"
	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	Id        string
	Name      string
	Desc      string
	Type      string
	Args      []string
	Collector prometheus.Collector
}

func NewMetric(subsystem, name, namespace, desc, metricType string, args ...string) *Metric {
	var collector prometheus.Collector
	switch metricType {
	case Counter:
		collector = prometheus.NewCounter(prometheus.CounterOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      name,
			Help:      desc,
		})
	case CounterVec:
		collector = prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      name,
			Help:      desc,
		}, args)
	case Gauge:
		collector = prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      name,
			Help:      desc,
		})
	case GaugeVec:
		collector = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      name,
			Help:      desc,
		}, args)
	case Histogram:
		collector = prometheus.NewHistogram(prometheus.HistogramOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      name,
			Help:      desc,
		})
	case HistogramVec:
		collector = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      name,
			Help:      desc,
		}, args)
	case Summary:
		collector = prometheus.NewSummary(prometheus.SummaryOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      name,
			Help:      desc,
		})
	case SummaryVec:
		collector = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Subsystem: subsystem,
			Namespace: namespace,
			Name:      name,
			Help:      desc,
		}, args)
	default:
		panic(errors.New("type of invalid indicator"))
	}

	return &Metric{
		Id:        util.UUID16(),
		Name:      name,
		Desc:      desc,
		Type:      metricType,
		Args:      args,
		Collector: collector,
	}
}
