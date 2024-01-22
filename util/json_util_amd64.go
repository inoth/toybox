package util

import "github.com/bytedance/sonic" // sonic 不支持 arm 环境

func JsonString(data interface{}) string {
	buf, err := sonic.Marshal(&data)
	if err != nil {
		return ""
	}
	return string(buf)
}

func JsonParse[T interface{}](data []byte) (T, error) {
	var res T
	err := sonic.Unmarshal(data, &res)
	return res, err
}
