package udpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"sync/atomic"
	"time"

	"github.com/inoth/toybox/util"
	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

const (
	lengthPrefix = 4
)

var (
	newline       = []byte{'\n'}
	space         = []byte{' '}
	defaultConfig = UDPClientConfig{
		WriteWait:      10 * time.Second,
		PongWait:       10 * time.Second,
		PingPeriod:     (10 * time.Second) * 9 / 10,
		MaxMessageSize: 1 << 10,
		Addr:           "localhost:4242",
		TlsConf: &tls.Config{
			InsecureSkipVerify: true, // 对于自签名证书
			NextProtos:         []string{"quic-echo-example"},
		},
	}
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
	WriteWait      time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
	MaxMessageSize uint
	Addr           string
	TlsConf        *tls.Config
}

func NewClient(ctx context.Context, cfgs ...UDPClientConfig) (*UDPClient, error) {
	cfg := util.First(defaultConfig, cfgs)

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
		_ = c.conn.CloseWithError(0, "connection closed")
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
	_ = stream.SetReadDeadline(time.Now().Add(time.Duration(c.cfg.PongWait)))
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
			if n < lengthPrefix {
				continue
			}
			var idx uint32 = 0
			for int(idx) < n {
				msgLength := binary.BigEndian.Uint32(buf[idx : idx+lengthPrefix])
				msg := bytes.TrimSpace(bytes.Replace(buf[idx+lengthPrefix:idx+lengthPrefix+msgLength], newline, space, -1))
				idx = idx + lengthPrefix + msgLength
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
			_ = stream.SetWriteDeadline(time.Now().Add(c.cfg.WriteWait))
			if _, err := stream.Write([]byte{}); err != nil {
				return
			}
		case message, ok := <-c.send:
			_ = stream.SetWriteDeadline(time.Now().Add(c.cfg.WriteWait))
			if !ok {
				_, _ = stream.Write([]byte{})
				return
			}
			msg := make([]byte, 4+len(message))
			binary.BigEndian.PutUint32(msg, uint32(len(message)))
			for i := 0; i < len(message); i++ {
				msg[i+4] = message[i]
			}
			_, _ = stream.Write(msg)
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
