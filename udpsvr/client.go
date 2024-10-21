package udpsvr

import (
	"context"

	"github.com/quic-go/quic-go"
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

	sess quic.Connection
	hub  *UDPQuicServer
}

func NewClient(hub *UDPQuicServer, sess quic.Connection)
