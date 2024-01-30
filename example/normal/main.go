package main

import (
	"demo1/router"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr"
	"github.com/inoth/toybox/server/ginsvr/httpgin"
	"github.com/inoth/toybox/server/ginsvr/validaton"
	"github.com/inoth/toybox/server/metric"
)

func main() {
	tb := toybox.New(
		toybox.WithLoadConf(),
		logger.WithLogger(),
		metric.NewPrometheus(
			metric.WithNamespace("test"),
			metric.WithSubsystem("test"),
			metric.WithMetrics(
				metric.Metric{
					Name: "requests_total",
					Desc: "How many HTTP requests processed, partitioned by status code and HTTP method.",
					Type: metric.CounterVec,
					Args: []string{"code", "method", "handler", "host", "url"},
				},
				metric.Metric{
					Name: "request_duration_seconds",
					Desc: "The HTTP request latencies in seconds.",
					Type: metric.HistogramVec,
					Args: []string{"code", "method", "url"},
				},
				metric.Metric{
					Name: "response_size_bytes",
					Desc: "The HTTP response sizes in bytes.",
					Type: metric.Summary,
				},
				metric.Metric{
					Name: "request_size_bytes",
					Desc: "The HTTP request sizes in bytes.",
					Type: metric.Summary,
				},
			)),
		httpgin.NewHttpGin(func(hgs *httpgin.HttpGinServer) {
			hgs.Use(ginsvr.Recovery())
		}, router.WithUserRouter(),
			httpgin.WithValidator(validaton.PhoneValidate, validaton.EmailValidate),
		),
	)
	if err := tb.Run(); err != nil {
		panic(err)
	}
}

// func(ut ut.Translator) error {
// 	return ut.Add("phoneValidate", "{0}不是一个合法的手机号", true)
// }, func(ut ut.Translator, fe validator.FieldError) string {
// 	t, _ := ut.T("phoneValidate", fe.Field())
// 	return t
// }
