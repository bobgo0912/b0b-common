package util

import (
	"fmt"
	"reflect"
	"strings"
)

func ConfigMask(value reflect.Value) string {
	var sb strings.Builder
	t := value.Type()
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		mask := t.Field(i).Tag.Get("mask")
		var ve interface{}
		if mask != "" {
			ve = "***"
		} else {
			ve = value.Field(i).Interface()
		}
		sb.WriteString(name)
		sb.WriteString(":")
		sb.WriteString(fmt.Sprint(ve))
		sb.WriteString(" ")
		fmt.Println(name, ve)
	}
	return sb.String()
}
