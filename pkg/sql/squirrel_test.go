package sql

import (
	"fmt"
	"reflect"
	"testing"
)
import sq "github.com/Masterminds/squirrel"

func TestFF(t *testing.T) {
	sql, i, err := sq.Insert("t_ts").
		Columns("s2", "c1").Values("asd", "sd").Values("dsd", "asa").ToSql()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sql)
	t.Log(i)

}

func TestSD(t *testing.T) {

	var sd = []*Stu{{
		Id:   213,
		Name: "123",
		Age:  1,
	}}
	values := sq.Insert("ssd").Values(1, "21", 12).Values(3, "232", 44)
	for _, stu := range sd {
		of := reflect.ValueOf(*stu)
		its := make([]interface{}, 0)
		for i := 0; i < of.NumField(); i++ {
			its = append(its, of.Field(i).Interface())
		}
		values = values.Values(its...)
	}
	sql, i, err := values.ToSql()
	fmt.Println(sql)
	fmt.Println(i)
	fmt.Println(err)
}
