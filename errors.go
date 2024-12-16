package toybox

import "errors"

var (
	ErrNotConfig = errors.New("unable to load configuration")
	ErrRestart   = errors.New("restart")
)
