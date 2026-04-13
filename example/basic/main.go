// basic 展示了 toybox 最基本的用法：
// 自定义 Transport 并通过 ToyBox 管理其生命周期。
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/inoth/toybox"
)

// httpServer 实现 transport.Transport 接口，包装 net/http.Server。
type httpServer struct {
	srv *http.Server
}

func newHTTPServer(addr string, handler http.Handler) *httpServer {
	return &httpServer{
		srv: &http.Server{Addr: addr, Handler: handler},
	}
}

func (h *httpServer) Start(ctx context.Context) error {
	log.Printf("HTTP server listening on %s", h.srv.Addr)
	if err := h.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (h *httpServer) Stop(ctx context.Context) error {
	log.Println("HTTP server shutting down...")
	return h.srv.Shutdown(ctx)
}

// Endpoint 实现 registry.Endpointer 接口, 便于服务注册时自动获取地址。
func (h *httpServer) Endpoint() (string, error) {
	return fmt.Sprintf("http://%s", h.srv.Addr), nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from toybox!")
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	srv := newHTTPServer(":8080", mux)

	tb := toybox.New(
		toybox.WithServiceInfo("basic-service", "1.0.0"),
		toybox.WithServer(srv),
	)

	if err := tb.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
