package deafult

import (
	"github.com/inoth/ino-toybox/servers/wssvc/accumulator"
	"github.com/inoth/ino-toybox/servers/wssvc/models"
	"github.com/inoth/ino-toybox/servers/wssvc/msg_pipeline/process"
)

type DeafultMsgProcess struct{}

func (DeafultMsgProcess) Process(msgbody models.MessageBody, acc accumulator.Accumulator) {
	acc.Next(msgbody)
}

func init() {
	process.AddProcess("deafult", &DeafultMsgProcess{})
}
