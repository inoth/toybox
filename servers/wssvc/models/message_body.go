package models

import (
	"encoding/json"
	"time"
)

const (
	// 普通消息
	EventMsg = "msg"
	// 命令消息
	EventCall = "call"
	// 连接初始化
	EventInit = "init"
	// 系统消息
	EventSystem = "system"
)

type MessageBody interface {
	// 获取发送目标数组
	GetTargets() []string
	// 获取需要发送的body
	GetJsonBody() []byte
	// 实现拿到进行下一步需要的实体
	GenNextBody() MessageBody
}

type SourceMessage struct {
	Token string `json:"token"`
	// msg | auth | call
	Event     string    `json:"event"`
	EventBody EventBody `json:"eventBody"`
}
type EventBody struct {
	Targets   []string `json:"targets"`
	Msg       string   `json:"msg"`
	Timestamp int64    `json:"timestamp"`
}

func (sm SourceMessage) GetJsonBody() []byte {
	buf, _ := json.Marshal(sm)
	return buf
}

func (sm SourceMessage) GetTargets() []string {
	return sm.EventBody.Targets
}

func (sm SourceMessage) GenNextBody() MessageBody {
	return &SendMessage{
		Token: sm.Token,
		Ids:   sm.EventBody.Targets,
		Body: SendEventBody{
			Event:     sm.Event,
			Msg:       sm.EventBody.Msg,
			Timestamp: time.Now().Unix(),
		},
	}
}

type SendMessage struct {
	Ids   []string      `json:"-"`
	Token string        `json:"-"`
	Body  SendEventBody `json:"body"`
}

func (smb *SendMessage) GetJsonBody() []byte {
	return smb.Body.GetJsonBody()
}

func (smb *SendMessage) GetTargets() []string {
	return smb.Ids
}

func (SendMessage) GenNextBody() MessageBody {
	return nil
}

type SendEventBody struct {
	Source    UserInfo `json:"source"`
	Event     string   `json:"event"`
	Msg       string   `json:"msg"`
	Timestamp int64    `json:"timestamp"`
}

func (seb *SendEventBody) GetJsonBody() []byte {
	buf, _ := json.Marshal(seb)
	return buf
}
