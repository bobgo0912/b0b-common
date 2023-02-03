package sql

import "testing"
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
