package websocket

type MessageHandleFunc func(ctx *Context)

func defaultMessageHandle(c *Context) {
	c.SendMessage(NewOutputMessage(c.input.Body, c.input.Id))
}

type OutputMessage struct {
	Targets []string
	Body    []byte
}

func NewOutputMessage(body []byte, targets ...string) OutputMessage {
	return OutputMessage{Body: body, Targets: targets}
}

type InputMessage struct {
	Id   string
	Body []byte
}

func NewMessage(id string, msg []byte) InputMessage {
	return InputMessage{Id: id, Body: msg}
}

func (ws *WebsocketServer) inputMessageHandle(msg InputMessage) {
	c := ws.pool.Get().(*Context)
	c.reset()
	c.input = msg

	ws.handleMessage(c)

	ws.output <- c.output
	ws.pool.Put(c)
}

func (ws *WebsocketServer) handleMessage(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	if len(ws.messageHandlers) <= 0 {
		defaultMessageHandle(c)
		return
	}
	for _, handler := range ws.messageHandlers {
		handler(c)
	}
}

func (ws *WebsocketServer) outputMessageHandle(msg OutputMessage) {
	for _, target := range msg.Targets {
		if client, ok := ws.clients[target]; ok {
			client.send <- msg.Body
		}
	}
}

func SendMessage(msg InputMessage) {
	hub.inputMessageHandle(msg)
}
