package toybox

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/inoth/toybox/common/signal"
	"github.com/inoth/toybox/component"
	"github.com/inoth/toybox/internal"
	"github.com/inoth/toybox/server"
)

var (
	Toml = "toml"
	Yaml = "yaml"
	Json = "json"
)

type ToyBox struct {
	cfgPath   string
	cfgType   string
	cfgSource string //local || http

	componentMap map[string]struct{}
	components   []component.Component
	servers      []server.Server
}

type Option func(tb *ToyBox)

func New(opts ...Option) *ToyBox {
	tb := &ToyBox{
		cfgPath:      "config/",
		cfgType:      Toml,
		componentMap: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(tb)
	}
	return tb
}

func (tb *ToyBox) init() error {
	cfgByte, err := tb.readFile(Toml)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}
	err = tb.resolveConfig(cfgByte)
	if err != nil {
		return fmt.Errorf("failed to resolve configuration: %v", err)
	}
	return nil
}

func (tb *ToyBox) Run() {
	err := tb.init()
	if err != nil {
		panic(err)
	}

	tb.initComponents()

	// 加载服务信息
	err = tb.loadServers()
	if err != nil {
		panic(err)
	}

	for _, svc := range tb.servers {
		go func(svc server.Server) {
			fmt.Printf("server %v is up and running\n", svc.Name())
			err := svc.Start()
			if err != nil {
				panic(fmt.Sprintf("%v service has an exception: %v\n", svc.Name(), err))
			}
		}(svc)
	}

	// 之后再使用监听信号阻塞主进程
	signal.ListenSignal()
}

func (tb *ToyBox) initComponents() {
	var wg sync.WaitGroup
	for _, comp := range tb.components {
		wg.Add(1)
		go func(comp component.Component) {
			defer wg.Done()

			err := comp.Init()
			if err != nil {
				panic(fmt.Sprintf("component initialization failed: %v", err))
			}
			tb.componentMap[comp.Name()] = struct{}{}
		}(comp)
	}
	wg.Wait()
}

func (tb *ToyBox) loadServers() error {
	if len(server.Servers) <= 0 {
		return fmt.Errorf("no available services found")
	}
	for _, svc := range server.Servers {
		server := svc()
		b := true
		for _, requir := range server.RequiredComponent() {
			if _, ok := tb.componentMap[requir]; !ok {
				b = false
			}
		}
		if !b {
			return fmt.Errorf("service %v missing essential components: %v", server.Name(), server.RequiredComponent())
		}
		tb.servers = append(tb.servers, server)
	}
	return nil
}

func (tb *ToyBox) readFile(fileType string) ([]byte, error) {
	dev := os.Getenv("GORUNEVN")
	if len(dev) > 0 {
		tb.cfgPath = tb.cfgPath + "/" + dev
	}
	files, err := internal.WalkPath(tb.cfgPath)
	if err != nil {
		return nil, err
	}
	var (
		cfgTmp  []byte
		cfgByte []byte
	)
	for _, file := range files {
		fileSlice := strings.Split(file, ".")
		cfgType := fileSlice[len(fileSlice)-1]
		if cfgType != Toml {
			continue
		}
		cfgTmp, err = os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败,Err:%s", err.Error())
		}
		cfgByte = append(cfgByte, []byte("\n")...)
		cfgByte = append(cfgByte, cfgTmp...)
	}
	return cfgByte, nil
}
