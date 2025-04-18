package wsserver

import "sync"

var clientPool = sync.Pool{
	New: func() any {
		return &Client{}
	},
}

func ClientGet(sendSize ...int) *Client {
	client := clientPool.Get().(*Client)
	client.Reset(sendSize...)
	return client
}

func ClientPut(client *Client) {
	if client == nil {
		return
	}
	clientPool.Put(client)
}
