package metric

import (
	"context"
	"net/http"

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
)

type Prometheus struct {
	ready     bool
	Port      string `toml:"port"`
	Subsystem string `toml:"subsystem"`
	Namespace string `toml:"namespace"`

	Metrics    []Metric
	collectors map[string]*metrics
}

func (pm *Prometheus) IsReady() {
	pm.ready = true
}

func (pm *Prometheus) Ready() bool {
	return pm.ready
}

func (pm *Prometheus) Name() string {
	return pm.Namespace
}

func (pm *Prometheus) register() error {
	for _, item := range pm.Metrics {
		col := item.initMetric(pm.Subsystem, pm.Namespace)
		if _, ok := pm.collectors[item.Name]; !ok {
			pm.collectors[item.Name] = col

			if err := prometheus.Register(col.collector); err != nil {
				return errors.Wrap(err, "register initMetric err")
			}
		}
	}
	return nil
}

func (pm *Prometheus) Run(ctx context.Context) error {

	if err := pm.register(); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(pm.Port, nil); err != nil {
		return errors.Wrap(err, "run ListenAndServe err")
	}
	return nil
}
