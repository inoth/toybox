package toml

import (
	"testing"

	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/config/file"
	"github.com/stretchr/testify/require"
)

type TestEntiy struct {
	name string `toml:"-"`
	Key  string `toml:"key"`
}

func (t *TestEntiy) Name() string {
	return t.name
}

func TestNewConfiguration(t *testing.T) {
	var (
		dir = "../"
	)
	tests := []struct {
		name   string
		dir    string
		data   TestEntiy
		expect string
	}{
		{
			name:   "local config",
			dir:    dir,
			data:   TestEntiy{name: "test1"},
			expect: "asdasdjaoisda1",
		},
		{
			name:   "local config",
			dir:    dir,
			data:   TestEntiy{name: "test2"},
			expect: "asdasdjaoisda2",
		},
		{
			name:   "local config",
			dir:    dir,
			data:   TestEntiy{name: "test3"},
			expect: "asdasdjaoisda3",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewConfiguration(
				config.WithSource(
					file.NewSource(test.dir),
				),
			)
			err := got.PrimitiveDecode(&test.data)
			if err != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.expect, test.data.Key, "NewConfiguration() = %v, want %v", test.data.Key, test.expect)
		})
	}
}
