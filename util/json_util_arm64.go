package util

import "encoding/json"

func JsonString(data interface{}) string {
	buf, err := json.Marshal(&data)
	if err != nil {
		return ""
	}
	return string(buf)
}

func JsonParse[T interface{}](data []byte) (T, error) {
	var res T
	err := json.Unmarshal(data, &res)
	return res, err
}
