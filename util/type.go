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
