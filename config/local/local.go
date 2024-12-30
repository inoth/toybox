package local

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/inoth/toybox/util/file"
	"github.com/pkg/errors"
)

type LocalConfig struct {
	Dir   string
	Paths []string
}

func NewSource(dir string) *LocalConfig {
	return &LocalConfig{
		Dir: dir,
	}
}

func (c *LocalConfig) Load() (str string, err error) {
	if c.Paths == nil {
		cfgEnv := os.Getenv("CONFIG_ENV")
		if cfgEnv != "" {
			c.Dir = filepath.Join(c.Dir, cfgEnv)
		}
		c.Paths, err = file.PathGlobPattern(c.Dir)
		if err != nil {
			return "", errors.Wrap(err, "no configuration available")
		}
	}
	var sb strings.Builder
	for _, path := range c.Paths {
		buf, err := file.ReadFile(path)
		if err != nil {
			fmt.Printf("%s read file err: %v", path, err)
			continue
		}
		sb.Write(buf)
		sb.WriteString("\n")
	}
	str = sb.String()
	return
}
