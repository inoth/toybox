package metric

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/inoth/toybox"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// https://github.com/zsais/go-gin-prometheus/blob/master/middleware.go

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

var (
	prom *Prometheus
	once sync.Once
)

type Prometheus struct {
	ready      bool
	name       string
	reg        *prometheus.Registry
	collectors map[string]prometheus.Collector

	Port      string   `toml:"port"`
	Subsystem string   `toml:"subsystem"`
	Namespace string   `toml:"namespace"`
	Metrics   []Metric `toml:"metrics"`
}

func NewPrometheus(opts ...Option) toybox.Option {
	once.Do(func() {
		o := defaultOption()
		for _, opt := range opts {
			opt(&o)
		}
		prom = &o
	})
	return func(tb *toybox.ToyBox) {
		tb.AppendServer(prom)
	}
}

func (pm *Prometheus) IsReady() {
	pm.ready = true
}

func (pm *Prometheus) Ready() bool {
	return pm.ready
}

func (pm *Prometheus) Name() string {
	return pm.name
}

func (pm *Prometheus) register() error {
	for _, item := range pm.Metrics {
		col := item.init(pm.Subsystem, pm.Namespace)
		if _, ok := pm.collectors[item.Name]; !ok {
			if pm.reg != nil {
				if err := pm.reg.Register(col); err != nil {
					return errors.Wrap(err, "register metric err")
				}
			} else {
				if err := prometheus.Register(col); err != nil {
					return errors.Wrap(err, "register metric err")
				}
			}
			pm.collectors[item.Name] = col
		}
	}
	return nil
}

func (pm *Prometheus) Run(ctx context.Context) error {
	if !pm.ready {
		return fmt.Errorf("server %s not ready", pm.name)
	}

	if err := pm.register(); err != nil {
		return err
	}
	if pm.reg != nil {
		http.Handle("/metrics", promhttp.HandlerFor(pm.reg, promhttp.HandlerOpts{Registry: pm.reg}))
	} else {
		http.Handle("/metrics", promhttp.Handler())
	}
	if err := http.ListenAndServe(pm.Port, nil); err != nil {
		return errors.Wrap(err, "run metric server err")
	}
	return nil
}
