package utils

import "encoding/json"

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
