package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	default_name = "metric"
)

type Option func(*Prometheus)

func defaultOption() Prometheus {
	return Prometheus{
		ready:      true,
		name:       default_name,
		Port:       ":8081",
		Subsystem:  "",
		Namespace:  "",
		collectors: make(map[string]prometheus.Collector),
	}
}

func WithMetrics(metrics ...Metric) Option {
	return func(pm *Prometheus) {
		pm.Metrics = append(pm.Metrics, metrics...)
	}
}

func WithMetric(name, desc, metricType string, args ...string) Option {
	return func(pm *Prometheus) {
		pm.Metrics = append(pm.Metrics, Metric{
			Name: name,
			Desc: desc,
			Type: metricType,
			Args: args,
		})
	}
}

func WithNamespace(namespace string) Option {
	return func(pm *Prometheus) {
		pm.Namespace = namespace
	}
}

func WithSubsystem(subsystem string) Option {
	return func(pm *Prometheus) {
		pm.Subsystem = subsystem
	}
}

func WithPort(port string) Option {
	return func(pm *Prometheus) {
		pm.Port = port
	}
}

func WithNewRegistry(cs ...prometheus.Collector) Option {
	if len(cs) <= 0 {
		cs = append(cs,
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
	}
	return func(pm *Prometheus) {
		if pm.reg != nil {
			return
		}
		pm.reg = prometheus.NewRegistry()
		pm.reg.MustRegister(cs...)
	}
}
