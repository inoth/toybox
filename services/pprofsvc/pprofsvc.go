package pprofsvc

import (
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/inoth/toybox/utils"
)

type PprofSvc struct {
	Port string
}

func New(port ...string) *PprofSvc {
	return &PprofSvc{
		Port: utils.FirstParam(":9002", port),
	}
}

func (ps *PprofSvc) Start() error {
	env := os.Getenv("GORUNEVN")
	if env == "dev" || env == "debug" {
		return http.ListenAndServe(ps.Port, nil)
	} else {
		return nil
	}
}

func (gsvc *PprofSvc) Stop() {}
