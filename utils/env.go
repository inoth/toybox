package utils

import "os"

const (
	EnvDebug = "debug"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func GetRunEnv() string {
	e := os.Getenv("GORUNEVN")
	return e
}
