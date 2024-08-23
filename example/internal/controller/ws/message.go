package ws

import (
	"fmt"

	"github.com/inoth/toybox/wssvr"
)

type MessageController struct {
}

func NewMessageController() *MessageController {
	return &MessageController{}
}

func (m *MessageController) Handler() wssvr.HandlerFunc {
	return func(c *wssvr.Context) {
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

type body struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}
