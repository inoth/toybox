package udpsvr

import "sync"

type HandlerFunc func(*Context)

type Context struct {
	body  []byte
	state bool

	Keys map[string]any
	m    sync.RWMutex
	uqs  *UDPQuicServer
}
