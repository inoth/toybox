package license

import (
	"encoding/json"
	"fmt"
	"time"
)

type LicenseComponent struct{}

type Subscribe struct {
	AppName   string
	Name      string
	NodeLimit int
	State     bool
	Describe  string
	Expire    time.Time
}

type License struct {
	LegalMachine map[string]struct{}
	Subscribes   map[string]Subscribe
}

func (l *License) String() []byte {
	buf, err := json.Marshal(l)
	if err != nil {
		fmt.Printf("ERR: %v", err.Error())
		return []byte("")
	}
	return buf
}

func (l *License) AppendMachine(machineCode ...string) {
	for _, code := range machineCode {
		if _, ok := l.LegalMachine[code]; ok {
			fmt.Printf("WRAN: Duplicate Registration Machine %v", code)
			continue
		}
		l.LegalMachine[code] = struct{}{}
	}
}

func (l *License) AppendSubscribe(subs ...Subscribe) {
	for _, sub := range subs {
		if _, ok := l.Subscribes[sub.AppName]; ok {
			fmt.Printf("WRAN: Duplicate Registration Module %v", sub.AppName)
			continue
		}
		l.Subscribes[sub.AppName] = sub
	}
}
