package util

import (
	"reflect"
	"testing"
)

func TestHttpGetWith(t *testing.T) {
	type args struct {
		url     string
		params  map[string]string
		token   string
		headers map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "t1",
			args: args{
				url: "https://baidu.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HttpGetWith[map[string]any](tt.args.url, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpGetWith() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HttpGetWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
