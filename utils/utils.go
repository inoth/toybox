package utils

import "encoding/json"

func Max[T int | int32 | int64 | float32 | float64 | uint | uint32 | uint64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T int | int32 | int64 | float32 | float64 | uint | uint32 | uint64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func FirstParam[T interface{}](defaultArg T, args []T) T {
	if len(args) > 0 {
		defaultArg = args[0]
	}
	return defaultArg
}

// 解析json
func JsonMarshal[T interface{}](str string) (T, error) {
	var res T
	err := json.Unmarshal([]byte(str), &res)
	return res, err
}

func JsonMarshalByte[T interface{}](str []byte) (T, error) {
	var res T
	err := json.Unmarshal([]byte(str), &res)
	return res, err
}

// 求交集
func Intersect[T string | int | float32 | float64](slice1, slice2 []T) []T {
	m := make(map[T]int)
	nn := make([]T, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 求差集 slice1-并集
func Difference[T string | int | float32 | float64](slice1, slice2 []T) []T {
	m := make(map[T]int)
	nn := make([]T, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

//Substr 字符串的截取
func Substr(str string, start int64, end int64) string {
	length := int64(len(str))
	if start < 0 || start > length {
		return ""
	}
	if end < 0 {
		return ""
	}
	if end > length {
		end = length
	}
	return string(str[start:end])
}

func getValueWithType(m map[string]interface{}, key string) (interface{}, bool) {
	val, ok := m[key]
	if !ok {
		return nil, false
	}
	return val, true
}

func GetIntValue(m map[string]interface{}, key string) (int, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return 0, false
	}
	if iVal, ok := val.(int); ok {
		return iVal, true
	}
	return 0, false
}

func GetFloatValue(m map[string]interface{}, key string) (float64, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return 0, false
	}
	if fVal, ok := val.(float64); ok {
		return fVal, true
	}
	return 0, false
}

func GetStringValue(m map[string]interface{}, key string) (string, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return "", false
	}
	if sVal, ok := val.(string); ok {
		return sVal, true
	}
	return "", false
}

func GetBoolValue(m map[string]interface{}, key string) (bool, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return false, false
	}
	if bVal, ok := val.(bool); ok {
		return bVal, true
	}
	return false, false
}
