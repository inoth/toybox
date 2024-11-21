package uqs

import (
	"fmt"

	"github.com/inoth/toybox/udpsvr"
)

type MessageController struct {
}

func NewMessageController() *MessageController {
	return &MessageController{}
}

func (m *MessageController) Handler() udpsvr.HandlerFunc {
	return func(c *udpsvr.Context) {
		fmt.Printf("%v\n", string(c.Body()))

		var data body
		err := c.BindJson(&data)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		for i := 0; i < 3; i++ {
			c.String(data.ID, data.Body)
		}
	}
}

type body struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}
