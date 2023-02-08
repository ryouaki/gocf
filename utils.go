package gocf

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func InterfaceToString(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func IndexOfStringArray(arr []string, el string) int {
	return HasInArray(arr, el, func(a string, b string) bool {
		return a == b
	})
}

func HasInArray[T any](arr []T, el T, cp func(a T, b T) bool) int {
	if cp == nil {
		return -1
	}
	for i, v := range arr {
		if cp(v, el) {
			return i
		}
	}

	return -1
}

func GoCFLog(args ...any) {
	goArgs := make([]any, 1, 4)
	goArgs[0] = "[GoCF]:"
	for _, v := range args {
		goArgs = append(goArgs, v)
	}
	fmt.Println(goArgs...)
}
