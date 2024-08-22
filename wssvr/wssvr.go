package wssvr

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/inoth/toybox/util"
)

const (
	name = "websocket"
)

type WebsocketServer struct {
	option

	m        sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
	upgrader websocket.Upgrader

	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client

	input  chan []byte
	output chan []byte
}

func New(opts ...Option) *WebsocketServer {
	o := option{
		WriteWait:       10 * time.Second,
		PongWait:        60 * time.Second,
		PingPeriod:      (60 * time.Second) * 9 / 10,
		MaxMessageSize:  1 << 10,
		ReadBufferSize:  1 << 10,
		WriteBufferSize: 1 << 10,
		ChannelSize:     100,
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &WebsocketServer{
		option: o,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *WebsocketServer) Name() string {
	return name
}

func (w *WebsocketServer) Start(ctx context.Context) error {
	w.upgrader = websocket.Upgrader{
		ReadBufferSize:  int(w.ReadBufferSize),
		WriteBufferSize: int(w.WriteBufferSize),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	w.clients = make(map[string]*Client)
	w.register = make(chan *Client, util.Max(1, int(w.ChannelSize)/10))
	w.unregister = make(chan *Client, util.Max(1, int(w.ChannelSize)/10))
	w.input = make(chan []byte, w.ChannelSize)
	w.output = make(chan []byte, w.ChannelSize)

	return w.run()
}

func (w *WebsocketServer) Stop(ctx context.Context) error {
	close(w.register)
	close(w.unregister)
	close(w.input)
	close(w.output)
	return nil
}

func (w *WebsocketServer) run() error {
	for {
		select {
		case <-w.ctx.Done():
			return w.ctx.Err()
		case client := <-w.register:
			w.registerClient(client)
		case client := <-w.unregister:
			w.unregisterClient(client)
		case msg := <-w.input:
			fmt.Printf("%v", string(msg))
		case msg := <-w.output:
			fmt.Printf("%v", string(msg))
		}
	}
}

func (w *WebsocketServer) registerClient(client *Client) {
	w.m.Lock()
	defer w.m.Unlock()

	if val, ok := w.clients[client.ID]; ok {
		val.Close()
	}
	w.clients[client.ID] = client
}

func (w *WebsocketServer) unregisterClient(client *Client) {
	w.m.Lock()
	defer w.m.Unlock()

	delete(w.clients, client.ID)
	client.Close()
}
