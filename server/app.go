package server

type Creator func() Server

var Servers map[string]Creator = make(map[string]Creator)

func Add(name string, creator Creator) {
	Servers[name] = creator
}
