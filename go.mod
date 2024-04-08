module github.com/inoth/toybox

go 1.21

retract (
	[v1.0.0, v1.1.9]
	[v0.0.1, v0.9.9]
)

require (
	github.com/bytedance/sonic v1.10.2
	github.com/go-resty/resty/v2 v2.11.0
	github.com/google/uuid v1.6.0
	github.com/hpcloud/tail v1.0.0
	github.com/pkg/errors v0.9.1
)

require (
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/stretchr/testify v1.8.3 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)
