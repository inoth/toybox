package wssvr

import (
	"context"
	"log"
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
	pool     sync.Pool

	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client

	input  chan []byte
	output chan Message
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
	ws := &WebsocketServer{
		option: o,
	}
	ws.pool = sync.Pool{New: func() any {
		return &Context{ws: ws}
	}}
	return ws
}

func (w *WebsocketServer) Name() string {
	return name
}

func (w *WebsocketServer) Start(ctx context.Context) error {
	w.ctx, w.cancel = context.WithCancel(ctx)
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
	w.output = make(chan Message, w.ChannelSize)

	if len(w.handles) == 0 {
		w.handles = append(w.handles, defaultHandle())
	}

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
			go func(msg []byte) {
				defer func() {
					if err := recover(); err != nil {
						log.Printf("%v\n", err)
					}
				}()
				w.sendMessage(msg)
			}(msg)
		case msg := <-w.output:
			if client, ok := w.clients[msg.ID]; ok {
				client.send <- msg.Body
			}
		}
	}
}

func (w *WebsocketServer) sendMessage(msg []byte) {
	c := w.pool.Get().(*Context)
	c.reset()

	c.send(msg)
	for _, handle := range w.handles {
		if !c.state {
			break
		}
		handle(c)
	}

	w.pool.Put(c)
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

	if _, ok := w.clients[client.ID]; ok {
		delete(w.clients, client.ID)
		client.Close()
	}
}

func (w *WebsocketServer) SendMessage(msg []byte) {
	w.input <- msg
}
