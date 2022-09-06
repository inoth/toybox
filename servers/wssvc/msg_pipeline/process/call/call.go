package call

import (
	"regexp"
	"strings"

	"github.com/inoth/ino-toybox/servers/wssvc/accumulator"
	"github.com/inoth/ino-toybox/servers/wssvc/models"
	"github.com/inoth/ino-toybox/servers/wssvc/msg_pipeline/process"
)

type CallMsgProcess struct{}

func (CallMsgProcess) Process(msgbody models.MessageBody, acc accumulator.Accumulator) {
	switch body := msgbody.(type) {
	case *models.SendMessage:
		if body.Body.Event != models.EventCall {
			acc.Next(body)
			return
		}
		cmdpt := "^call\\s(expel|quit)\\s(\\w+)$"
		ok, _ := regexp.Match(cmdpt, []byte(body.Body.Msg))
		if !ok {
			acc.Next(body)
			return
		}
		if !handlerCmd(body, acc) {
			acc.Next(body)
			return
		}
	default:
		acc.Next(msgbody)
	}
}

func init() {
	process.AddProcess("deafult", &CallMsgProcess{})
}

func handlerCmd(body *models.SendMessage, acc accumulator.Accumulator) bool {
	params := strings.Split(body.Body.Msg, " ")
	if len(params) != 3 {
		return false
	}
	cmd := params[1]
	param := params[2]
	switch cmd {
	case "repeat": // 重复发送
		body.Body.Msg = param
		for i := 0; i < 2; i++ {
			acc.Next(body)
		}
		return true
	default:
		return false
	}
}
