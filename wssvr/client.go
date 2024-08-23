package wssvr

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/inoth/toybox/util"
	"github.com/pkg/errors"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	send chan []byte

	ID string

	ctx    context.Context
	cancel context.CancelFunc

	conn *websocket.Conn
	hub  *WebsocketServer
}

func (c *Client) Close() {
	close(c.send)
	c.cancel()
}

func NewClient(hub *WebsocketServer, w http.ResponseWriter, r *http.Request) (string, error) {
	if hub == nil {
		panic(fmt.Errorf("WebsocketServer not init"))
	}
	conn, err := hub.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return "", errors.Wrap(err, "init upgrader failed")
	}
	client := &Client{
		ID:   util.UUID(32),
		send: make(chan []byte, hub.ChannelSize),
		conn: conn,
		hub:  hub,
	}
	client.ctx, client.cancel = context.WithCancel(hub.ctx)

	go client.read()
	go client.write()

	hub.register <- client
	return client.ID, nil
}

func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(c.hub.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.hub.PongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.hub.PongWait))
		return nil
	})
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil || websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
			}
			msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
			c.hub.input <- msg
		}
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(c.hub.PingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.hub.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(c.hub.WriteWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			for i := 0; i < len(c.send); i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
