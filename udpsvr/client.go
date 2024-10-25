package udpsvr

import (
	"bytes"
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/inoth/toybox/util"
	"github.com/quic-go/quic-go"
)

const (
	IsDebug = false
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
		c.conn.CloseWithError(0, "connection closed")
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
	stream.SetReadDeadline(time.Now().Add(c.svr.PongWait))
	for {
		select {
		case <-c.ctx.Done():
		case <-stream.Context().Done():
			return
		default:
			n, err := stream.Read(buf)
			if err != nil {
				return
			}
			msg := buf[:n]
			msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
			if c.svr.GZIP {
				if buf, err := util.DecompressGzip(msg); err == nil {
					c.svr.input <- buf
				}
			} else {
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
		case <-stream.Context().Done():
			return
		case <-ticker.C:
			stream.SetWriteDeadline(time.Now().Add(c.svr.WriteWait))
			if _, err := stream.Write([]byte("PING")); err != nil {
				return
			}
		case message, ok := <-c.send:
			stream.SetWriteDeadline(time.Now().Add(c.svr.WriteWait))
			if !ok {
				stream.Write([]byte{})
				return
			}
			if c.svr.GZIP {
				if compressed, err := util.CompressGzip(message); err == nil {
					stream.Write(compressed)
				}
			} else {
				stream.Write(message)
			}
		}
	}
}
