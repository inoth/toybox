package source

import (
	"github.com/inoth/ino-toybox/servers/wssvc/accumulator"
	"github.com/inoth/ino-toybox/servers/wssvc/models"
	"github.com/inoth/ino-toybox/servers/wssvc/msg_pipeline/process"
	"github.com/inoth/ino-toybox/util/auth"
)

//  解析token中数据, 获取发送者用户信息
type GetSourceMsgProcess struct{}

func (GetSourceMsgProcess) Process(msgbody models.MessageBody, acc accumulator.Accumulator) {
	switch body := msgbody.(type) {
	case *models.SendMessage:
		if len(body.Token) > 0 {
			user, err := auth.ParseToken(body.Token)
			if err != nil {
				acc.Err(err)
				return
			}
			body.Body.Source = models.UserInfo{
				Id:   user.Uid,
				Name: user.Name,
				Icon: user.Avater,
			}
			acc.Next(body)
			return
		} else if body.Body.Event == "system" {
			body.Body.Source = models.UserInfo{
				Id:   "system",
				Name: "system",
				Icon: "system.png",
			}
			acc.Next(body)
			return
		} else {
			acc.ErrStr("invalid message sources")
			return
		}
	default:
		acc.ErrStr("assertion message structure failed.")
		return
	}
}

func init() {
	process.AddProcess("source", &GetSourceMsgProcess{})
}
