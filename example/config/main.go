// config 展示了如何使用 toybox 的配置管理功能：
// 从本地 YAML 文件加载配置，支持热重载。
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/conf"
	"github.com/inoth/toybox/conf/source"
)

// configHTTPServer 实现 transport.Transport 和 conf.ConfigureMatcher 接口。
// toybox 在启动 transport 前，会通过 TransportName() 匹配配置 section，
// 然后调用 PrimitiveDecode 将该 section 反序列化到本结构体。
type configHTTPServer struct {
	Addr         string `yaml:"addr"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`

	srv *http.Server
}

// TransportName 实现 conf.ConfigureMatcher 接口，
// 返回值对应 config.yaml 中的顶层 key。
func (s *configHTTPServer) TransportName() string {
	return "http"
}

func newConfigHTTPServer(handler http.Handler) *configHTTPServer {
	return &configHTTPServer{
		srv: &http.Server{Handler: handler},
	}
}

func (s *configHTTPServer) Start(ctx context.Context) error {
	s.srv.Addr = s.Addr
	if d, err := time.ParseDuration(s.ReadTimeout); err == nil {
		s.srv.ReadTimeout = d
	}
	if d, err := time.ParseDuration(s.WriteTimeout); err == nil {
		s.srv.WriteTimeout = d
	}
	log.Printf("HTTP server listening on %s (read_timeout=%s, write_timeout=%s)", s.srv.Addr, s.ReadTimeout, s.WriteTimeout)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *configHTTPServer) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *configHTTPServer) Endpoint() (string, error) {
	return fmt.Sprintf("http://%s", s.srv.Addr), nil
}

func main() {
	configPath := "config.yaml"

	// 创建配置管理器：从本地文件加载 YAML，并开启文件变更监听
	fileSource := source.NewFile(configPath)
	fileWatcher := source.NewFileWatcher(configPath, 3*time.Second)

	mgr, err := conf.NewManager(fileSource,
		conf.WithFormat(conf.FormatYAML),
		conf.WithWatcher(fileWatcher),
	)
	if err != nil {
		log.Fatalf("init config manager: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from config example!")
	})

	srv := newConfigHTTPServer(mux)

	tb := toybox.New(
		toybox.WithServiceInfo("config-service", "1.0.0"),
		toybox.WithConfig(mgr),
		toybox.WithServer(srv),
	)

	if err := tb.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
