// bootstrap 展示了使用 bootstrap.Config 一键初始化注册中心和配置源。
// 需要一个已运行的 etcd 实例 (localhost:2379) 和本地配置文件。
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/bootstrap"
)

type httpServer struct {
	srv *http.Server
}

func (h *httpServer) Start(ctx context.Context) error {
	log.Printf("HTTP server listening on %s", h.srv.Addr)
	if err := h.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (h *httpServer) Stop(ctx context.Context) error {
	return h.srv.Shutdown(ctx)
}

func (h *httpServer) Endpoint() (string, error) {
	return fmt.Sprintf("http://%s", h.srv.Addr), nil
}

func main() {
	// bootstrap.Load() 会从以下来源合并配置（优先级从高到低）：
	// 1. 命令行参数 (--service-name, --registry-type 等)
	// 2. 环境变量 (TOYBOX_SERVICE_NAME, TOYBOX_REGISTRY_TYPE 等)
	// 3. 默认值
	//
	// 也可以像下面这样直接构造：
	cfg := &bootstrap.Config{
		ServiceName:    "bootstrap-service",
		ServiceVersion: "1.0.0",
		ConfigFile:     "config.yaml",
		ConfigFormat:   "yaml",
		// 如需服务注册，取消下方注释并确保 etcd 已运行：
		// RegistryType:      "etcd",
		// RegistryEndpoints: []string{"localhost:2379"},
		// RegistryTimeout:   bootstrap.Duration{Duration: 5 * time.Second},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from bootstrap example!")
	})

	srv := &httpServer{
		srv: &http.Server{Addr: ":8082", Handler: mux},
	}

	tb := toybox.New(
		toybox.WithBootstrap(cfg),
		toybox.WithServer(srv),
		toybox.WithStopTimeout(10*time.Second),
	)

	if err := tb.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
