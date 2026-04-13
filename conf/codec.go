package conf

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Format represents a config file format.
type Format string

const (
	FormatYAML Format = "yaml"
	FormatJSON Format = "json"
	FormatTOML Format = "toml"
)

// Codec handles marshaling and unmarshaling of config data in a specific format.
type Codec interface {
	// Unmarshal parses raw bytes into a generic map.
	Unmarshal(data []byte, v interface{}) error
	// Marshal encodes a value into bytes.
	Marshal(v interface{}) ([]byte, error)
}

// NewCodec returns a Codec for the given format.
func NewCodec(format Format) (Codec, error) {
	switch format {
	case FormatYAML:
		return &yamlCodec{}, nil
	case FormatJSON:
		return &jsonCodec{}, nil
	case FormatTOML:
		return &tomlCodec{}, nil
	default:
		return nil, fmt.Errorf("unsupported config format: %s", format)
	}
}

// yamlCodec implements Codec for YAML.
type yamlCodec struct{}

func (c *yamlCodec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (c *yamlCodec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// jsonCodec implements Codec for JSON.
type jsonCodec struct{}

func (c *jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (c *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// tomlCodec implements Codec for TOML.
type tomlCodec struct{}

func (c *tomlCodec) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

func (c *tomlCodec) Marshal(v interface{}) ([]byte, error) {
	buf, err := toml.Marshal(v)
	return buf, err
}
