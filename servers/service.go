package servers

/*
	常用服务, 比如 http、websocket 监听等
*/

type Server interface {
	Start() error
}
