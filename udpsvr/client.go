package udpsvr

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/inoth/toybox/util"
	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// const addr = ":4242"
// func main() {
// 	// 设置监听地址
// 	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Server listening on %s\n", addr)
// 	for {
// 		// 接受客户端连接
// 		sess, err := listener.Accept(context.Background())
// 		if err != nil {
// 			log.Println("Error accepting connection:", err)
// 			continue
// 		}
// 		go handleClient(sess)
// 	}
// }
// func handleClient(sess quic.Connection) {
// 	fmt.Println("New connection accepted")
// 	stream, err := sess.AcceptStream(context.Background())
// 	if err != nil {
// 		log.Println("Error accepting stream:", err)
// 		return
// 	}
// 	// 简单地回显客户端发送的消息
// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := stream.Read(buf)
// 		if err != nil {
// 			log.Println("Error reading from stream:", err)
// 			return
// 		}

// 		fmt.Printf("Received: %s\n", string(buf[:n]))

//			_, err = stream.Write([]byte("Ack: " + string(buf[:n])))
//			if err != nil {
//				log.Println("Error writing to stream:", err)
//				return
//			}
//		}
//	}
//
//	func generateTLSConfig() *tls.Config {
//		// 创建简单的 TLS 配置，游戏可以使用自己签名的证书
//		cert, err := tls.LoadX509KeyPair("cert/cert.pem", "cert/priv.key")
//		if err != nil {
//			log.Fatal(err)
//		}
//		return &tls.Config{
//			Certificates: []tls.Certificate{cert},
//			NextProtos:   []string{"quic-echo-example"},
//		}
//	}
type Client struct {
	send chan []byte

	ID     string
	ctx    context.Context
	cancel context.CancelFunc

	stream quic.Stream
	svr    *UDPQuicServer
}

func NewClient(svr *UDPQuicServer, conn quic.Connection) error {
	if svr == nil {
		return fmt.Errorf("UDPQuicServer not init")
	}

	stream, err := conn.AcceptStream(svr.ctx)
	if err != nil {
		return errors.Wrap(err, "Error accepting stream")
	}

	client := &Client{
		ID:     util.UUID(32),
		send:   make(chan []byte, svr.ChannelSize),
		stream: stream,
		svr:    svr,
	}
	client.ctx, client.cancel = context.WithCancel(svr.ctx)

	go client.read()
	go client.write()

	svr.register <- client
	return nil
}

func (c *Client) Close() {
	close(c.send)
	c.cancel()
}

func (c *Client) read() {
	defer func() {
		c.svr.unregister <- c
		c.stream.Close()
	}()
	buf := make([]byte, c.svr.MaxMessageSize)
	c.stream.SetReadDeadline(time.Now().Add(c.svr.PongWait))
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			n, err := c.stream.Read(buf)
			if err != nil {
				return
			}
			msg := buf[:n]
			bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
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

func (c *Client) write() {
	ticker := time.NewTicker(c.svr.PingPeriod)
	defer func() {
		ticker.Stop()
		c.svr.unregister <- c
		c.stream.Close()
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.stream.SetWriteDeadline(time.Now().Add(c.svr.WriteWait))
			if _, err := c.stream.Write([]byte{}); err != nil {
				return
			}
		case message, ok := <-c.send:
			c.stream.SetWriteDeadline(time.Now().Add(c.svr.WriteWait))
			if !ok {
				c.stream.Write([]byte{})
				return
			}
			if c.svr.GZIP {
				if compressed, err := util.CompressGzip(message); err == nil {
					c.stream.Write(compressed)
				}
			} else {
				c.stream.Write(message)
			}
		}
	}
}
