package util

import (
	"strings"
	"sync"
)

var Sbp = sync.Pool{New: func() any {
	return &strings.Builder{}
}}

func GetString(s1, s2 string) string {
	builder := Sbp.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		Sbp.Put(builder)
	}()
	builder.WriteString(s1)
	builder.WriteString("-")
	builder.WriteString(s2)
	return builder.String()
}

func GetString3(s1, s2, s3 string) string {
	builder := Sbp.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		Sbp.Put(builder)
	}()
	builder.WriteString(s1)
	builder.WriteString("-")
	builder.WriteString(s2)
	builder.WriteString("-")
	builder.WriteString(s3)
	return builder.String()
}

func GetStrings(str ...string) string {
	builder := Sbp.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		Sbp.Put(builder)
	}()
	i := len(str)
	for j, s := range str {
		builder.WriteString(s)
		if j != (i - 1) {
			builder.WriteString("-")
		}
	}
	return builder.String()
}
