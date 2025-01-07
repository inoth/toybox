package udpserver

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestRunUdpServer(t *testing.T) {
	us := New(
		WithTLSFile("../test_resources/cert/cert.pem", "../test_resources/cert/priv.key"),
		WithHandler(
			func(c *Context) {
				fmt.Println(string(c.Body()))
				c.String("testclient", "hello testclient")
			},
		))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := us.Start(ctx)
	if err != nil {
		t.Fatalf("%v", err)
		return
	}
	t.Log("ok")
}
