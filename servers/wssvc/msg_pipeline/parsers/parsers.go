package parsers

import (
	"errors"

	"github.com/inoth/ino-toybox/servers/wssvc/accumulator"
)

type MsgParsers interface {
	Parser(body []byte, acc accumulator.Accumulator)
}

var Parsers = map[string]MsgParsers{}

func AddParsers(key string, parser MsgParsers) {
	Parsers[key] = parser
}

func GetParsers(key string) (MsgParsers, error) {
	if p, ok := Parsers[key]; ok {
		return p, nil
	}
	return nil, errors.New("not found parser.")
}
