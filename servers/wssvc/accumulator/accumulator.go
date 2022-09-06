package accumulator

import (
	"errors"

	"github.com/inoth/ino-toybox/servers/wssvc/models"
)

type Accumulator interface {
	Next(body models.MessageBody)
	Err(err error)
	ErrStr(msg string)
	Error() error
}

type accumulator struct {
	body chan<- models.MessageBody
	err  error
}

func (acc *accumulator) Next(body models.MessageBody) {
	acc.body <- body
}

func (acc *accumulator) Err(err error) {
	acc.err = err
}

func (acc *accumulator) ErrStr(msg string) {
	acc.err = errors.New(msg)
}

func (acc *accumulator) Error() error {
	return acc.err
}

func NewAccumulator(body chan<- models.MessageBody) Accumulator {
	return &accumulator{
		body: body,
	}
}
