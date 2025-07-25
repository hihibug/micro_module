package utils

import (
	"bytes"
	"encoding/json"
	"reflect"
	"regexp"
	"unsafe"
)

// FmtLongInt 超过16位的数字转为字符串
func FmtLongInt(data interface{}) interface{} {
	j, _ := json.Marshal(data)
	reg := regexp.MustCompile(`:(\d{16,20})`)
	l := len(reg.FindAllString(string(j), -1)) //正则匹配16-20位的数字，如果找到了就开始正则替换并解析

	if l != 0 {
		var mapResult interface{}
		str := reg.ReplaceAllString(string(j), `:"${1}"`)
		_ = json.Unmarshal([]byte(str), &mapResult)
		data = &mapResult
	}

	return data
}

// IsLittleEndian 决定主机字节顺序是否为小端序
func IsLittleEndian() bool {
	n := 0x1234
	return *(*byte)(unsafe.Pointer(&n)) == 0x34
}

// IsInterfaceNil 判读接口是否为Nil
func IsInterfaceNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

// ToFormattedJSON traverses the golang value to formatted JSON string, such as a struct, map, slice, array and so on
func ToFormattedJSON(obj interface{}) (string, error) {
	bs, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = json.Indent(buf, bs, "", "\t")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// InArray 判断某一个值是否含在切片之中
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
