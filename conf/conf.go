package conf

import (
	"context"
	"github/inoth/toybox"
	"sync"
)

var (
	conf = sync.Map{}
)

type Conf interface {
	Configuration(key string, cpt toybox.Component) error
}

// type JsonConfComponent struct{}

// type YamlConfComponent struct{}

type TomlConfComponent struct{}

func (tc *TomlConfComponent) Configuration(key string, cpt toybox.Component) error {
	return nil
}

func (tc *TomlConfComponent) Ready() bool {
	return false
}

func (tc *TomlConfComponent) Init(ctx context.Context) error {
	return nil
}
