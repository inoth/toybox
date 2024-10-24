package udpsvr

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/inoth/toybox/util"
	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

const (
	name = "udp"
)

type UDPQuicServer struct {
	option

	m      sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	pool   sync.Pool

	listen *quic.Listener

	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client

	input  chan []byte
	output chan Message
}

func New(opts ...Option) *UDPQuicServer {
	o := option{
		WriteWait:       10 * time.Second,
		PongWait:        60 * time.Second,
		PingPeriod:      (60 * time.Second) * 9 / 10,
		MaxMessageSize:  1 << 10,
		ReadBufferSize:  1 << 10,
		WriteBufferSize: 1 << 10,
		ChannelSize:     100,
		CertFile:        "cert.pem",
		KeyFile:         "priv.key",
	}
	for _, opt := range opts {
		opt(&o)
	}
	if o.serverName == "" {
		o.serverName = name
	}
	uqs := &UDPQuicServer{
		option: o,
	}
	uqs.pool = sync.Pool{New: func() any {
		return &Context{svr: uqs}
	}}
	return uqs
}

func (uq *UDPQuicServer) Name() string {
	return uq.serverName
}

func (uq *UDPQuicServer) Start(ctx context.Context) error {
	uq.ctx, uq.cancel = context.WithCancel(ctx)

	uq.clients = make(map[string]*Client)
	uq.register = make(chan *Client, util.Max(1, int(uq.ChannelSize)))
	uq.unregister = make(chan *Client, util.Max(1, int(uq.ChannelSize)))
	uq.input = make(chan []byte, uq.ChannelSize)
	uq.output = make(chan Message, uq.ChannelSize)

	tlsConfig, err := generateTLSConfig(uq.CertFile, uq.KeyFile)
	if err != nil {
		return errors.Wrap(err, "load certificate failed")
	}
	listen, err := quic.ListenAddr(uq.Addr, tlsConfig, nil)
	if err != nil {
		return errors.Wrap(err, "listening failed")
	}
	fmt.Printf("Server listening on %s\n", uq.Addr)

	uq.listen = listen

	return uq.run()
}

func (uq *UDPQuicServer) Stop(ctx context.Context) error {
	close(uq.register)
	close(uq.unregister)
	close(uq.input)
	close(uq.output)

	return uq.listen.Close()
}

func (uq *UDPQuicServer) run() error {

	go uq.accept()

	for {
		select {
		case <-uq.ctx.Done():
			return uq.ctx.Err()
		case client := <-uq.register:
			uq.registerClient(client)
		case client := <-uq.unregister:
			uq.unregisterClient(client)
		case msg := <-uq.input:
			go func(msg []byte) {
				defer func() {
					if err := recover(); err != nil {
						log.Printf("%v\n", err)
					}
				}()
				uq.sendMessage(msg)
			}(msg)
		case msg := <-uq.output:
			if client, ok := uq.clients[msg.ID]; ok {
				client.send <- msg.Body
			}
		}
	}
}

func (uq *UDPQuicServer) accept() {
	for {
		select {
		case <-uq.ctx.Done():
			return
		default:
			conn, err := uq.listen.Accept(context.Background())
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			if err = NewClient(uq, conn); err != nil {
				fmt.Println("Error new client:", err)
				continue
			}
		}
	}
}

func (uq *UDPQuicServer) registerClient(client *Client) {
	uq.m.Lock()
	defer uq.m.Unlock()

	if val, ok := uq.clients[client.ID]; ok {
		val.Close()
	}
	uq.clients[client.ID] = client
}

func (uq *UDPQuicServer) unregisterClient(client *Client) {
	uq.m.Lock()
	defer uq.m.Unlock()

	if _, ok := uq.clients[client.ID]; ok {
		delete(uq.clients, client.ID)
		client.Close()
	}
}

func (uq *UDPQuicServer) SendMessage(msg []byte) {
	uq.input <- msg
}

func (uq *UDPQuicServer) sendMessage(msg []byte) {
	c := uq.pool.Get().(*Context)
	c.reset()

	c.send(msg)
	for _, handle := range uq.handles {
		if !c.state {
			break
		}
		handle(c)
	}

	uq.pool.Put(c)
}
