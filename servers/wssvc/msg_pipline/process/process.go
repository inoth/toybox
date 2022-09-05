package process

import (
	"errors"

	"github.com/inoth/ino-toybox/servers/wssvc/accumulator"
	"github.com/inoth/ino-toybox/servers/wssvc/models"
)

type MsgProcess interface {
	Process(msgbody models.MessageBody, acc accumulator.Accumulator)
}

var Process = map[string]MsgProcess{}

func AddProcess(key string, decorator MsgProcess) {
	Process[key] = decorator
}

func GetParsers(key string) (MsgProcess, error) {
	if p, ok := Process[key]; ok {
		return p, nil
	}
	return nil, errors.New("not found process.")
}
