package websocket

import (
	"sync"
	"time"
)

const (
	default_name = "websocket"
)

type Option func(*WebsocketServer)

func defaultOption() WebsocketServer {
	return WebsocketServer{
		name:              default_name,
		ready:             true,
		lock:              sync.RWMutex{},
		clients:           make(map[string]*Client),
		WriteWait:         10 * time.Second,
		PongWait:          60 * time.Second,
		PingPeriod:        (60 * time.Second) * 9 / 10,
		MaxMessageSize:    1 << uint(10),
		ReadBufferSize:    1 << uint(10),
		WriteBufferSize:   1 << uint(10),
		MaxMsgChannelSize: 100,
	}
}

func WithName(name string) Option {
	return func(ws *WebsocketServer) {
		ws.name = name
	}
}

func WithMessageHandler(handleFunc ...MessageHandleFunc) Option {
	return func(ws *WebsocketServer) {
		ws.messageHandlers = append(ws.messageHandlers, handleFunc...)
	}
}
