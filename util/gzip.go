package util

import (
	"bytes"
	"compress/gzip"
	"io"
)

func DecompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var result bytes.Buffer
	if _, err := io.Copy(&result, reader); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

func CompressGzip(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	defer writer.Close()

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
