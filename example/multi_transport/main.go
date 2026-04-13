// multi_transport 展示了如何在 toybox 中同时运行多个 Transport：
// 例如同时运行 HTTP API 服务器和 内部 metrics 服务器。
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/inoth/toybox"
)

type httpTransport struct {
	name string
	srv  *http.Server
}

func newHTTPTransport(name, addr string, handler http.Handler) *httpTransport {
	return &httpTransport{
		name: name,
		srv:  &http.Server{Addr: addr, Handler: handler},
	}
}

func (h *httpTransport) Start(ctx context.Context) error {
	log.Printf("[%s] listening on %s", h.name, h.srv.Addr)
	if err := h.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (h *httpTransport) Stop(ctx context.Context) error {
	log.Printf("[%s] shutting down...", h.name)
	return h.srv.Shutdown(ctx)
}

func (h *httpTransport) Endpoint() (string, error) {
	return fmt.Sprintf("http://%s", h.srv.Addr), nil
}

func main() {
	// API 服务
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`)
	})
	apiMux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})
	apiServer := newHTTPTransport("api", ":8080", apiMux)

	// Metrics / 内部管理服务
	metricsMux := http.NewServeMux()
	metricsMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "uptime_seconds %d\n", 42)
	})
	metricsMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	metricsServer := newHTTPTransport("metrics", ":9090", metricsMux)

	tb := toybox.New(
		toybox.WithServiceInfo("multi-transport-service", "1.0.0"),
		toybox.WithServer(apiServer),
		toybox.WithServer(metricsServer),
		toybox.WithStopTimeout(15*time.Second),
		toybox.WithMetadata(map[string]string{
			"region": "cn-east-1",
			"env":    "staging",
		}),
	)

	if err := tb.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
