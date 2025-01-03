package ws

import (
	"fmt"

	"github.com/inoth/toybox/wsserver"
)

type body struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type MessageController struct{}

func NewMessageController() *MessageController {
	return &MessageController{}
}

func (m *MessageController) Handler() wsserver.HandlerFunc {
	return func(c *wsserver.Context) {
		fmt.Printf("%v\n", string(c.Body()))

		var data body
		err := c.BindJson(&data)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		c.String(data.ID, data.Body)
	}
}
