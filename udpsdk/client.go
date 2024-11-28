package udpsdk

import (
	"bytes"
	"context"
	"crypto/tls"
	"sync/atomic"
	"time"

	"github.com/inoth/toybox/util"
	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type UDPClient struct {
	send    chan []byte
	receive chan []byte

	closed atomic.Int32

	ID string

	ctx    context.Context
	cancel context.CancelFunc

	conn quic.Connection

	cfg *UDPClientConfig
}

type UDPClientConfig struct {
	Gzip           bool
	WriteWait      time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
	MaxMessageSize uint
	Addr           string
	TlsConf        *tls.Config
}

func NewClient(ctx context.Context, cfgs ...UDPClientConfig) (*UDPClient, error) {
	cfg := util.First(UDPClientConfig{
		Gzip:           false,
		WriteWait:      10 * time.Second,
		PongWait:       10 * time.Second,
		PingPeriod:     (10 * time.Second) * 9 / 10,
		MaxMessageSize: 1 << 10,
		Addr:           "localhost:4242",
		TlsConf: &tls.Config{
			InsecureSkipVerify: true, // 对于自签名证书
			NextProtos:         []string{"quic-echo-example"},
		},
	}, cfgs)

	conn, err := quic.DialAddr(context.Background(), cfg.Addr, cfg.TlsConf, nil)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	client := &UDPClient{
		ID:      util.UUID(32),
		conn:    conn,
		cfg:     &cfg,
		send:    make(chan []byte, 10),
		receive: make(chan []byte, 10),
	}
	client.ctx, client.cancel = context.WithCancel(ctx)

	stream, err := client.conn.OpenStreamSync(client.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "open stream failed")
	}

	go client.read(stream)
	go client.write(stream)

	return client, nil
}

func (c *UDPClient) Close() error {
	if c.closed.CompareAndSwap(0, 1) {
		c.conn.CloseWithError(0, "connection closed")
		close(c.send)
		close(c.receive)
		c.cancel()
	}
	return nil
}

func (c *UDPClient) read(stream quic.Stream) {
	defer func() {
		stream.Close()
		c.Close()
	}()
	buf := make([]byte, c.cfg.MaxMessageSize)
	stream.SetReadDeadline(time.Now().Add(time.Duration(c.cfg.PongWait)))
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-stream.Context().Done():
			return
		default:
			n, err := stream.Read(buf)
			if err != nil {
				return
			}
			msg := bytes.TrimSpace(bytes.Replace(buf[:n], newline, space, -1))
			if c.cfg.Gzip {
				if buf, err := util.DecompressGzip(msg); err == nil {
					c.receive <- buf
				}
			} else {
				c.receive <- msg
			}
		}
	}
}

func (c *UDPClient) write(stream quic.Stream) {
	ticker := time.NewTicker(c.cfg.PingPeriod)
	defer func() {
		ticker.Stop()
		stream.Close()
		c.Close()
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-stream.Context().Done():
			return
		case <-ticker.C:
			stream.SetWriteDeadline(time.Now().Add(c.cfg.WriteWait))
			if _, err := stream.Write([]byte{}); err != nil {
				return
			}
		case message, ok := <-c.send:
			stream.SetWriteDeadline(time.Now().Add(c.cfg.WriteWait))
			if !ok {
				stream.Write([]byte{})
				return
			}
			if c.cfg.Gzip {
				if compressed, err := util.CompressGzip(message); err == nil {
					stream.Write(compressed)
				}
			} else {
				stream.Write(message)
			}
			for i := 0; i < len(c.send); i++ {
				stream.Write(newline)
				stream.Write(<-c.send)
			}
		}
	}
}

func (c *UDPClient) SendMessage(msg []byte) {
	c.send <- msg
}

func (c *UDPClient) ReceiveMessage(msg chan<- []byte) {
	for v := range c.receive {
		msg <- v
	}
}
