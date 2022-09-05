package json

import (
	"encoding/json"

	"github.com/inoth/ino-toybox/components/logger"
	"github.com/inoth/ino-toybox/servers/wssvc/accumulator"
	"github.com/inoth/ino-toybox/servers/wssvc/models"
	"github.com/inoth/ino-toybox/servers/wssvc/msg_pipline/parsers"
)

type JsonParser struct{}

func (JsonParser) Parser(msgbody []byte, acc accumulator.Accumulator) {
	logger.Log.Info(string(msgbody))

	var body models.TextMessage
	err := json.Unmarshal(msgbody, &body)
	if err != nil {
		logger.Log.Error(err.Error())
		acc.Err(err)
		return
	}
	acc.Next(&body)
}

func init() {
	parsers.AddParsers("json", &JsonParser{})
}
