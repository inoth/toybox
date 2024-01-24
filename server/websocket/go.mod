module github.com/inoth/toybox/server/websocket

go 1.21.6

replace github.com/inoth/toybox => ../../

replace github.com/inoth/toybox/component/logger => ../../component/logger

require (
	github.com/gorilla/websocket v1.5.1
	github.com/inoth/toybox v0.0.0
	github.com/inoth/toybox/component/logger v0.0.0
	github.com/pkg/errors v0.9.1
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/bytedance/sonic v1.10.2 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/go-resty/resty/v2 v2.11.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/arch v0.6.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
