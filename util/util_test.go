package util_test

import (
	"testing"

	"github.com/inoth/toybox/util"
)

func TestUUIDV7(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		ns   int
	}{
		{
			name: "default",
			ns:   16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.UUIDV7(tt.ns)
			// TODO: update the condition below to compare got with tt.want.
			if got == "" {
				t.Errorf("UUIDV7() = %v, want %v", got, tt.ns)
			}
			t.Logf("UUIDV7() = %v", got)
		})
	}
}
