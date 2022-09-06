package normal

import (
	"github.com/inoth/ino-toybox/servers/wssvc/accumulator"
	"github.com/inoth/ino-toybox/servers/wssvc/models"
	"github.com/inoth/ino-toybox/servers/wssvc/msg_pipeline/process"
)

type NormalMsgProcess struct{}

func (NormalMsgProcess) Process(msgbody models.MessageBody, acc accumulator.Accumulator) {
	acc.Next(msgbody)
}

func init() {
	process.AddProcess("normal", &NormalMsgProcess{})
}
