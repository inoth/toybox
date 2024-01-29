package websocket

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/inoth/toybox/component/logger"
	"github.com/inoth/toybox/util"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	id     string
	cancel func()
	log    *logger.Logger
	hub    *WebsocketServer
	conn   *websocket.Conn
	send   chan []byte
}

func NewClient(w http.ResponseWriter, r *http.Request) (string, error) {
	if hub == nil {
		return "", fmt.Errorf("services are not yet ready")
	}
	conn, err := hub.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return "", errors.Wrap(err, "upgrade connect fail")
	}
	client := &Client{
		id:   util.RandStr(),
		hub:  hub,
		conn: conn,
		send: make(chan []byte, hub.MaxMessageSize),
	}
	client.log = logger.GetLogger(logger.LoggerConfig{ServerName: "wsid:" + client.id})

	var ctx context.Context
	ctx, client.cancel = context.WithCancel(hub.ctx)
	client.run(ctx)

	hub.register <- client
	return client.id, nil
}

func (c *Client) run(ctx context.Context) {
	go c.read(ctx)
	go c.write(ctx)
}

func (c *Client) stop() {
	c.cancel()
}

func (c *Client) read(ctx context.Context) {
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
		case <-ctx.Done():
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil || websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.log.Error(err.Error())
				return
			}
			msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
			c.hub.input <- NewMessage(c.id, msg)
		}
	}
}

func (c *Client) write(ctx context.Context) {
	ticker := time.NewTicker(c.hub.PingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.hub.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.log.Error(err.Error())
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
