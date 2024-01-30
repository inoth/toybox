package websocket

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/util"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var (
	hub *WebsocketServer
)

type WebsocketServer struct {
	name   string
	ready  bool
	lock   sync.RWMutex
	ctx    context.Context
	cancel func()

	pool            sync.Pool
	messageHandlers []MessageHandleFunc

	upgrader websocket.Upgrader

	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client

	input  chan InputMessage
	output chan OutputMessage

	ReadBufferSize    int64 `toml:"read_buffer_size"`
	WriteBufferSize   int64 `toml:"write_buffer_size"`
	MaxMsgChannelSize int64 `toml:"max_msg_channel_size"`

	WriteWait      time.Duration `toml:"write_wait"`
	PongWait       time.Duration `toml:"pong_wait"`
	PingPeriod     time.Duration `toml:"ping_period"`
	MaxMessageSize int64         `toml:"max_message_size"`
}

func NewWebsocketServer(opts ...Option) toybox.Option {
	ws := defaultOption()
	for _, opt := range opts {
		opt(&ws)
	}

	ws.pool = sync.Pool{
		New: func() any {
			return new(Context)
		},
	}

	hub = &ws
	return func(tb *toybox.ToyBox) {
		tb.AppendServer(&ws)
	}
}

func (ws *WebsocketServer) init(ctx context.Context) error {
	if !ws.ready {
		return fmt.Errorf("websocker server not ready")
	}
	ws.ctx, ws.cancel = context.WithCancel(ctx)
	ws.upgrader = websocket.Upgrader{
		ReadBufferSize:  int(ws.ReadBufferSize),
		WriteBufferSize: int(ws.WriteBufferSize),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws.register = make(chan *Client, util.Max(1, ws.MaxMsgChannelSize/10))
	ws.unregister = make(chan *Client, util.Max(1, ws.MaxMsgChannelSize/10))
	ws.input = make(chan InputMessage, ws.MaxMsgChannelSize)
	ws.output = make(chan OutputMessage, ws.MaxMsgChannelSize)
	return nil
}

func (ws *WebsocketServer) IsReady() {
	ws.ready = true
}

func (ws *WebsocketServer) Ready() bool {
	return ws.ready
}

func (ws *WebsocketServer) Name() string {
	return ws.name
}

func (ws *WebsocketServer) Run(ctx context.Context) error {
	defer ws.cancel()
	if err := ws.init(ctx); err != nil {
		return errors.Wrap(err, "ws.init failed")
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case client := <-ws.register:
			ws.registerClient(client)
		case client := <-ws.unregister:
			ws.unregisterClient(client)
		case msg, ok := <-ws.input:
			if !ok {
				ws.cancel()
				return fmt.Errorf("msg input channel closed")
			}
			ws.inputMessageHandle(msg)
		case msg, ok := <-ws.output:
			if !ok {
				ws.cancel()
				return fmt.Errorf("msg output channel closed")
			}
			ws.outputMessageHandle(msg)
		}
	}
	return nil
}

func (ws *WebsocketServer) registerClient(client *Client) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	if _, ok := ws.clients[client.id]; !ok {
		ws.clients[client.id] = client
	}
}

func (ws *WebsocketServer) unregisterClient(client *Client) {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	if _, ok := ws.clients[client.id]; ok {
		delete(ws.clients, client.id)
		client.stop()
	}
}
