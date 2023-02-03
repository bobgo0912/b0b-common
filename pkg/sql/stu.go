package sql

import (
	"context"
	"github.com/Masterminds/squirrel"
)

const StuTableName = "stu"

type Stu struct {
	Id   uint64
	Name string
	Age  int
}

type StuStore struct {
	Store *BaseStore[Stu]
}

func GetStuStore() StuStore {
	return StuStore{Store: &BaseStore[Stu]{Db: Db("edu", nil), TableName: StuTableName}}
}
func (s *StuStore) QueryById(ctx context.Context, id uint64) (*Stu, error) {
	return s.Store.QueryById(ctx, id)
}
func (s *StuStore) QueryPage(ctx context.Context, sb squirrel.SelectBuilder, page, size uint64) (*Pagination[Stu], error) {
	return s.Store.QueryPage(ctx, sb, page, size)
}
