package util

import (
	"fmt"
	"reflect"
)

func ToSlice(arr interface{}) []interface{} {
  v := reflect.ValueOf(arr)
  if v.Kind() != reflect.Slice {
    panic("toslice arr not slice")
  }
  l := v.Len()
  ret := make([]interface{}, l)
  for i := 0; i < l; i++ {
    ret[i] = v.Index(i).Interface()
  }
  return ret
}


func ToStr (i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func Substr(str string, start int, end int) string {
    rs := []rune(str)
    length := len(rs)

    if start < 0 || start > length {
	panic("start is wrong")
    }

    if end < 0 || end > length {
	panic("end is wrong")
    }

    return string(rs[start:end])
}
