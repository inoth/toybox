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
