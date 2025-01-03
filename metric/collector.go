package metric

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	name = "metric"
)

type Prometheus struct {
	option

	svr        *http.Server
	collectors map[string]prometheus.Collector
}

func New(opts ...Option) *Prometheus {
	o := option{
		Port: ":9000",
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &Prometheus{
		option:     o,
		collectors: make(map[string]prometheus.Collector),
	}
}

func (p *Prometheus) Name() string {
	return name
}

func (p *Prometheus) Start(ctx context.Context) error {
	if len(p.Metrics) <= 0 {
		return fmt.Errorf("metrics is empty")
	}
	reg := prometheus.NewRegistry()

	for _, metric := range p.Metrics {
		col := metric.init(p.Subsystem, p.Namespace)
		if err := reg.Register(col); err != nil {
			continue
		}
		p.collectors[metric.Name] = col
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	p.svr = &http.Server{Addr: p.Port, Handler: mux}

	if err := p.svr.ListenAndServe(); err != nil && err != context.Canceled && err != http.ErrServerClosed {
		return errors.Wrap(err, "start cronsvc metric err")
	}
	return nil
}

func (p *Prometheus) Stop(ctx context.Context) error {
	return p.svr.Shutdown(ctx)
}

func (p *Prometheus) GetCounter(name string) prometheus.Counter {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(prometheus.Counter); ok {
			return col
		}
	}
	return nil
}

func (p *Prometheus) GetCounterVec(name string) *prometheus.CounterVec {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(*prometheus.CounterVec); ok {
			return col
		}
	}
	return nil
}

func (p *Prometheus) GetGauge(name string) prometheus.Gauge {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(prometheus.Gauge); ok {
			return col
		}
	}
	return nil
}

func (p *Prometheus) GetGaugeVec(name string) *prometheus.GaugeVec {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(*prometheus.GaugeVec); ok {
			return col
		}
	}
	return nil
}

func (p *Prometheus) GetHistogram(name string) prometheus.Histogram {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(prometheus.Histogram); ok {
			return col
		}
	}
	return nil
}

func (p *Prometheus) GetHistogramVec(name string) *prometheus.HistogramVec {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(*prometheus.HistogramVec); ok {
			return col
		}
	}
	return nil
}

func (p *Prometheus) GetSummary(name string) prometheus.Summary {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(prometheus.Summary); ok {
			return col
		}
	}
	return nil
}

func (p *Prometheus) GetSummaryVec(name string) *prometheus.SummaryVec {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(*prometheus.SummaryVec); ok {
			return col
		}
	}
	return nil
}

func (p *Prometheus) CallCounter(name string, fn func(prometheus.Counter)) {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(prometheus.Counter); ok {
			fn(col)
		}
	}
}

func (p *Prometheus) CallCounterVec(name string, fn func(*prometheus.CounterVec)) {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(*prometheus.CounterVec); ok {
			fn(col)
		}
	}
}

func (p *Prometheus) CallGauge(name string, fn func(prometheus.Gauge)) {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(prometheus.Gauge); ok {
			fn(col)
		}
	}
}

func (p *Prometheus) CallGaugeVec(name string, fn func(*prometheus.GaugeVec)) {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(*prometheus.GaugeVec); ok {
			fn(col)
		}
	}
}

func (p *Prometheus) CallHistogram(name string, fn func(prometheus.Histogram)) {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(prometheus.Histogram); ok {
			fn(col)
		}
	}
}

func (p *Prometheus) CallHistogramVec(name string, fn func(*prometheus.HistogramVec)) {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(*prometheus.HistogramVec); ok {
			fn(col)
		}
	}
}

func (p *Prometheus) CallSummary(name string, fn func(prometheus.Summary)) {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(prometheus.Summary); ok {
			fn(col)
		}
	}
}

func (p *Prometheus) CallSummaryVec(name string, fn func(*prometheus.SummaryVec)) {
	if val, ok := p.collectors[name]; ok {
		if col, ok := val.(*prometheus.SummaryVec); ok {
			fn(col)
		}
	}
}
