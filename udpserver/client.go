package udpserver

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/inoth/toybox/util"
	"github.com/quic-go/quic-go"
)

const (
	IsDebug      = true
	lengthPrefix = 4
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	send   chan []byte
	closed atomic.Int32

	ID         string
	RemoteAddr string
	ctx        context.Context
	cancel     context.CancelFunc

	conn quic.Connection
	svr  *UDPQuicServer
}

func NewClient(svr *UDPQuicServer, conn quic.Connection) {
	if svr == nil || conn == nil {
		fmt.Println("NewClient failed")
		return
	}

	client := &Client{
		ID:         util.UUID(32),
		send:       make(chan []byte, svr.ChannelSize),
		conn:       conn,
		svr:        svr,
		RemoteAddr: conn.RemoteAddr().String(),
	}
	if IsDebug {
		client.ID = "testclient"
	}
	client.ctx, client.cancel = context.WithCancel(svr.ctx)

	stream, err := conn.AcceptStream(svr.ctx)
	if err != nil {
		fmt.Printf("Error client %s accepting stream: %v\n", client.ID, err)
	}

	go client.read(stream)
	go client.write(stream)

	svr.register <- client
}

func (c *Client) Close() {
	if c.closed.CompareAndSwap(0, 1) {
		_ = c.conn.CloseWithError(0, "connection closed")
		close(c.send)
		c.cancel()
	}
}

func (c *Client) read(stream quic.Stream) {
	defer func() {
		stream.Close()
		c.svr.unregister <- c
	}()
	buf := make([]byte, c.svr.MaxMessageSize)
	_ = stream.SetReadDeadline(time.Now().Add(c.svr.PongWait))
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
			if c.svr.Gzip {
				buf, err = util.DecompressGzip(buf)
				if err != nil {
					continue
				}
			}
			var idx uint32 = 0
			for int(idx) < n {
				msgLength := binary.BigEndian.Uint32(buf[idx : idx+lengthPrefix])
				msg := bytes.TrimSpace(bytes.Replace(buf[idx+lengthPrefix:idx+lengthPrefix+msgLength], newline, space, -1))
				idx = idx + lengthPrefix + msgLength
				c.svr.input <- msg
			}
		}
	}
}

func (c *Client) write(stream quic.Stream) {
	ticker := time.NewTicker(c.svr.PingPeriod)
	defer func() {
		ticker.Stop()
		stream.Close()
		c.svr.unregister <- c
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-stream.Context().Done():
			return
		case <-ticker.C:
			_ = stream.SetWriteDeadline(time.Now().Add(c.svr.WriteWait))
			if _, err := stream.Write([]byte{}); err != nil {
				return
			}
		case message, ok := <-c.send:
			_ = stream.SetWriteDeadline(time.Now().Add(c.svr.WriteWait))
			if !ok {
				_, _ = stream.Write([]byte{})
				return
			}
			msg := make([]byte, 4+len(message))
			binary.BigEndian.PutUint32(msg, uint32(len(message)))
			for i := 0; i < len(message); i++ {
				msg[i+4] = message[i]
			}
			if c.svr.Gzip {
				if compressed, err := util.CompressGzip(msg); err == nil {
					_, _ = stream.Write(compressed)
				}
			} else {
				_, _ = stream.Write(msg)
			}
		}
	}
}
