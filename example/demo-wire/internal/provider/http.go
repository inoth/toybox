package provider

import (
	"demo-wire/internal/controller"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/ginserver"
	"github.com/inoth/toybox/ginserver/middleware"
	"github.com/inoth/toybox/metric"
	"github.com/prometheus/client_golang/prometheus"
)

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

func RequestsTotal(p *metric.Prometheus) gin.HandlerFunc {
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

		p.CallHistogramVec("request_duration_seconds", func(hv *prometheus.HistogramVec) {
			hv.WithLabelValues(status, c.Request.Method, c.Request.URL.Path).Observe(elapsed)
		})
		p.CallCounterVec("requests_total", func(cv *prometheus.CounterVec) {
			cv.WithLabelValues(status, c.Request.Method, c.HandlerName(), c.Request.Method, c.Request.URL.Path).Inc()
		})
		p.CallSummary("request_size_bytes", func(s prometheus.Summary) {
			s.Observe(float64(reqSz))
		})
		p.CallSummary("response_size_bytes", func(s prometheus.Summary) {
			s.Observe(resSz)
		})
	}
}

func NewHttpServer(
	gc *controller.GreeterController,
	p *metric.Prometheus,
) *ginserver.GinHttpServer {
	return ginserver.NewHttp(
		ginserver.WithMiddleware(
			middleware.SetTraceId(),
			RequestsTotal(p),
		),
		ginserver.WithHandlers(gc),
	)
}

func NewHttp2Server(
	gc *controller.GreeterController,
) *ginserver.GinHttp2Server {
	return ginserver.NewHttp2(
		ginserver.WithMiddleware(middleware.SetTraceId()),
		ginserver.WithHandlers(gc),
	)
}

func NewHttp3Server(
	gc *controller.GreeterController,
) *ginserver.GinHttp3Server {
	return ginserver.NewHttp3(
		ginserver.WithMiddleware(middleware.SetTraceId()),
		ginserver.WithHandlers(gc),
	)
}
