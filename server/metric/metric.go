package metric

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
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

func GetCounter(name string) prometheus.Counter {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(prometheus.Counter); ok {
			return col
		}
	}
	return nil
}

func GetCounterVec(name string) *prometheus.CounterVec {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(*prometheus.CounterVec); ok {
			return col
		}
	}
	return nil
}

func GetGauge(name string) prometheus.Gauge {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(prometheus.Gauge); ok {
			return col
		}
	}
	return nil
}

func GetGaugeVec(name string) *prometheus.GaugeVec {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(*prometheus.GaugeVec); ok {
			return col
		}
	}
	return nil
}

func GetHistogram(name string) prometheus.Histogram {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(prometheus.Histogram); ok {
			return col
		}
	}
	return nil
}

func GetHistogramVec(name string) *prometheus.HistogramVec {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(*prometheus.HistogramVec); ok {
			return col
		}
	}
	return nil
}

func GetSummary(name string) prometheus.Summary {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(prometheus.Summary); ok {
			return col
		}
	}
	return nil
}

func GetSummaryVec(name string) *prometheus.SummaryVec {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(*prometheus.SummaryVec); ok {
			return col
		}
	}
	return nil
}

func CallCounter(name string, fn func(prometheus.Counter)) {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(prometheus.Counter); ok {
			fn(col)
		}
	}
}

func CallCounterVec(name string, fn func(*prometheus.CounterVec)) {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(*prometheus.CounterVec); ok {
			fn(col)
		}
	}
}

func CallGauge(name string, fn func(prometheus.Gauge)) {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(prometheus.Gauge); ok {
			fn(col)
		}
	}
}

func CallGaugeVec(name string, fn func(*prometheus.GaugeVec)) {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(*prometheus.GaugeVec); ok {
			fn(col)
		}
	}
}

func CallHistogram(name string, fn func(prometheus.Histogram)) {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(prometheus.Histogram); ok {
			fn(col)
		}
	}
}

func CallHistogramVec(name string, fn func(*prometheus.HistogramVec)) {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(*prometheus.HistogramVec); ok {
			fn(col)
		}
	}
}

func CallSummary(name string, fn func(prometheus.Summary)) {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(prometheus.Summary); ok {
			fn(col)
		}
	}
}

func CallSummaryVec(name string, fn func(*prometheus.SummaryVec)) {
	if val, ok := prom.collectors[name]; ok {
		if col, ok := val.(*prometheus.SummaryVec); ok {
			fn(col)
		}
	}
}
