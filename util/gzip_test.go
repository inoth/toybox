package util

import "testing"

func TestGzip(t *testing.T) {
	str := "hello world!!!!"
	gzStr, err := CompressGzip([]byte(str))
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	buf, err := DecompressGzip(gzStr)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Log(string(buf))
}
