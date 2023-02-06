package sql

import (
	"github.com/jmoiron/sqlx"
)

const StuTableName = "stu"

type Stu struct {
	Id   uint64 `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

type StuStore struct {
	*BaseStore[Stu]
}

func GetConnection() (*sqlx.DB, error) {
	if EduDb != nil {
		return EduDb, nil
	}
	var err error
	EduDb, err = Db("edu", nil)
	if err != nil {
		return nil, err
	}
	return EduDb, nil
}

func GetStuStore() (*StuStore, error) {
	connection, err := GetConnection()
	if err != nil {
		return nil, err
	}
	return &StuStore{&BaseStore[Stu]{Db: connection, TableName: StuTableName}}, nil
}

//func (s *StuStore) QueryById(ctx context.Context, id uint64) (*Stu, error) {
//	return s.Store.QueryById(ctx, id)
//}
//func (s *StuStore) QueryPage(ctx context.Context, sb squirrel.SelectBuilder, page, size uint64) (*Pagination[Stu], error) {
//	return s.Store.QueryPage(ctx, sb, page, size)
//}
