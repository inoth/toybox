package util

import (
	"reflect"
	"testing"
)

type A struct {
	Id   string
	Name string
	Age  int
}

type B struct {
	Id   string
	Name string
	Age  int
}

func TestConvertWithUnsafe(t *testing.T) {
	type testCase[I any, O any] struct {
		name string
		args A
		want *B
	}
	tests := []testCase[A, B]{
		// TODO: Add test cases.
		{name: "convert a->b", args: A{Id: "1", Name: "test", Age: 20}, want: &B{Id: "1", Name: "test", Age: 20}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ConvertWithUnsafe[A, B](&tt.args); err != nil || !reflect.DeepEqual(got, tt.want) {
				if err != nil {
					t.Errorf("ConvertWithUnsafe() err: %v", err)
				}
				t.Errorf("ConvertWithUnsafe() = %v, want %v", got, tt.want)
			}
		})
	}
}
