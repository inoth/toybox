package models

import "encoding/json"

type MessageBody interface {
	GetJsonBody() []byte
}

type SendMessageBody struct {
	Ids  []string    `json:"ids"`
	Body MessageBody `json:"body"`
}

type TextMessage struct {
	Source  UserInfo `json:"source"`
	MsgType string   `json:"msgType"`
	Msg     string   `json:"msg"`
}

func (tm *TextMessage) GetJsonBody() []byte {
	buf, _ := json.Marshal(tm)
	return buf
}
