package utils

import (
	"net"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
)

var LocalIP = net.ParseIP("127.0.0.1")
