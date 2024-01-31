package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox"
	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/server/ginsvr"
	"github.com/inoth/toybox/server/ginsvr/httpgin"
	"github.com/inoth/toybox/server/metric"
	"github.com/inoth/toybox/util"
	"github.com/prometheus/client_golang/prometheus"
)

// http://localhost:8080
// http://localhost:8081/metrics

func main() {
	tb := toybox.New(
		toybox.WithLoadConf(),
		logger.WithLogger(),
		metric.NewPrometheus(
			metric.WithNewRegistry(),
			metric.WithNamespace("test_server"),
			metric.WithSubsystem("default"),
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
					// Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
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
			),
		),
		httpgin.NewHttpGin(func(hgs *httpgin.HttpGinServer) {
			hgs.Use(ginsvr.Recovery())
		}, func(hgs *httpgin.HttpGinServer) {
			hgs.GET("", Middleware(), func(c *gin.Context) {
				sleepTime := time.Duration(rand.Intn(21)) * time.Millisecond
				time.Sleep(sleepTime)
				c.String(200, util.RandStr(1000))
			})
			hgs.GET("/v1", Middleware(), func(c *gin.Context) {
				sleepTime := time.Duration(rand.Intn(21)) * time.Millisecond
				time.Sleep(sleepTime)
				c.String(200, util.RandStr(1000))
			})
			hgs.GET("/v2", Middleware(), func(c *gin.Context) {
				sleepTime := time.Duration(rand.Intn(21)) * time.Millisecond
				time.Sleep(sleepTime)
				c.String(200, util.RandStr(1000))
			})
		}),
	)
	if err := tb.Run(); err != nil {
		panic(err)
	}
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}
		start := time.Now()
		reqSz := computeApproximateRequestSize(c.Request)

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		resSz := float64(c.Writer.Size())

		metric.CallHistogramVec("request_duration_seconds", func(hv *prometheus.HistogramVec) {
			hv.WithLabelValues(status, c.Request.Method, c.Request.URL.Path).Observe(elapsed)
		})
		metric.CallCounterVec("requests_total", func(cv *prometheus.CounterVec) {
			cv.WithLabelValues(status, c.Request.Method, c.HandlerName(), c.Request.Method, c.Request.URL.Path).Inc()
		})
		metric.CallSummary("request_size_bytes", func(s prometheus.Summary) {
			s.Observe(float64(reqSz))
		})
		metric.CallSummary("response_size_bytes", func(s prometheus.Summary) {
			s.Observe(resSz)
		})
	}
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}
	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)
	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
