package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/inoth/toybox/udpclient"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := udpclient.NewClient(ctx)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer client.Close()

	msg := make(chan []byte)
	go client.ReceiveMessage(msg)

	go func() {
		for i := 0; i < 10; i++ {
			data := &body{
				ID:   "testclient",
				Body: fmt.Sprintf("Hello %d, Game Server!", i),
			}
			buf, _ := json.Marshal(data)

			client.SendMessage(buf)
			// time.Sleep(time.Millisecond * 300)
		}
		time.Sleep(time.Second * 3)
		cancel()
		close(msg)
	}()

	for buf := range msg {
		fmt.Printf("Received message: %s\n", string(buf))
	}

}

type body struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}
