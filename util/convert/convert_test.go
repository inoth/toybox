package convert

import (
	"reflect"
	"testing"
)

func TestConvertByUnsafe(t *testing.T) {
	type Foo struct {
		Id   string
		Name string
	}
	type Bar struct {
		Id   string
		Name string
	}

	tests := []struct {
		name       string
		input      *Foo
		wantOutput *Bar
		wantErr    bool
	}{
		{name: "test_1", input: &Foo{Id: "1111", Name: "1111"}, wantOutput: &Bar{Id: "1111", Name: "1111"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, err := ConvertByUnsafe[Foo, Bar](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertByUnsafe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("ConvertByUnsafe() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
