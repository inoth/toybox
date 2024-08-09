package metric

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	name = "metric"
)

var (
	cm *CronsvcMetric
)

type CronsvcMetric struct {
	option

	svr *http.Server

	taskCount   prometheus.Counter
	currentTask prometheus.Gauge
	duration    *prometheus.HistogramVec
}

func New(opts ...Option) *CronsvcMetric {
	o := option{
		Port: ":9051",
	}
	for _, opt := range opts {
		opt(&o)
	}
	cm = &CronsvcMetric{
		option: o,
	}
	return cm
}

func (cc *CronsvcMetric) Name() string {
	return name
}

func (cc *CronsvcMetric) Start(ctx context.Context) error {
	cc.newMetrics()

	reg := prometheus.NewRegistry()

	reg.MustRegister(cc.taskCount)
	reg.MustRegister(cc.currentTask)
	reg.MustRegister(cc.duration)

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	cc.svr = &http.Server{Addr: cc.Port, Handler: mux}

	if err := cc.svr.ListenAndServe(); err != nil {
		return errors.Wrap(err, "start cronsvc metric err")
	}

	return nil
}

func (cc *CronsvcMetric) Stop(ctx context.Context) error {
	return cc.svr.Shutdown(ctx)
}

func (cc *CronsvcMetric) newMetrics() {
	cc.taskCount = prometheus.NewCounter(prometheus.CounterOpts{
		Subsystem: cc.Subsystem,
		Namespace: cc.Namespace,
		Name:      "task_run_total",
		Help:      "任务数量",
	})
	cc.currentTask = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: cc.Subsystem,
		Namespace: cc.Namespace,
		Name:      "current_task_total",
		Help:      "当前任务数量",
	})
	cc.duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: cc.Subsystem,
		Namespace: cc.Namespace,
		Name:      "task_run_duration_seconds",
		Help:      "任务运行耗时",
		Buckets:   prometheus.DefBuckets,
	}, []string{"task_id"})
}

func AddTaskCount(val float64) {
	if cm == nil {
		return
	}
	cm.taskCount.Add(val)
}

func SetCurrentTask(val float64) {
	if cm == nil {
		return
	}
	cm.currentTask.Set(val)
}

func SetDuration(taskId string, task_time float64) {
	if cm == nil {
		return
	}
	cm.duration.WithLabelValues(taskId).Observe(task_time)
}
