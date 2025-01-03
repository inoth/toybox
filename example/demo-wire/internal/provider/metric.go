package provider

import "github.com/inoth/toybox/metric"

func NewMetric() *metric.Prometheus {
	return metric.New(
	// metric.WithNamespace("demo"),
	// metric.WithSubsystem("demo"),
	// metric.WithMetric(
	// 	metric.Metric{
	// 		Name: "http_request_total",
	// 		Desc: "Total number of HTTP requests",
	// 		Type: metric.CounterVec,
	// 		Args: []string{"code", "method", "handler", "host", "url"},
	// 	},
	// 	metric.Metric{
	// 		Name: "request_duration_seconds",
	// 		Desc: "The HTTP request latencies in seconds.",
	// 		Type: metric.HistogramVec,
	// 		Args: []string{"code", "method", "url"},
	// 	},
	// 	metric.Metric{
	// 		Name: "response_size_bytes",
	// 		Desc: "The HTTP response sizes in bytes.",
	// 		Type: metric.Summary,
	// 	},
	// 	metric.Metric{
	// 		Name: "request_size_bytes",
	// 		Desc: "The HTTP request sizes in bytes.",
	// 		Type: metric.Summary,
	// 	},
	// ),
	)
}
