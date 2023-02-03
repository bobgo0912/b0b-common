package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lann/builder"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const otelName = "b0b-common/mysql"

var (
	//stuDb = Db("edu", nil)
	stuDb *sql.DB
)

type BaseStore[T any] struct {
	Db        *sqlx.DB
	TableName string
}

func Db(dbname string, c *config.MysqlCfg) *sqlx.DB {
	if c == nil {
		c = &config.Cfg.MysqlCfg
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		c.UserName, c.Password, c.Host, c.Port, dbname,
	)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Panic(err)
	}
	return db
}

func (s *BaseStore[T]) QueryById(ctx context.Context, id uint64, columns ...string) (*T, error) {
	spanCtx, span := newOTELSpan(ctx, "DB.QueryById")
	defer span.End()
	if len(columns) < 1 {
		columns = append(columns, "*")
	}
	toSql, param, err := squirrel.Select(columns...).From(s.TableName).Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		log.Debug("QueryById to sql fail err=", err.Error())
		return nil, errors.Wrap(err, "squirrel toSql fail")
	}
	row := s.Db.QueryRowxContext(spanCtx, toSql, param...)
	if row.Err() != nil {
		log.Debug("QueryxContext fail err=", err.Error())
		return nil, errors.Wrap(err, "Query fail")
	}
	var d T
	err = row.StructScan(&d)
	if err != nil {
		log.Debug("StructScan fail err=", err.Error())
		return nil, errors.Wrap(err, "StructScan fail")
	}
	return &d, nil
}

func (s *BaseStore[T]) QueryPage(ctx context.Context, sb squirrel.SelectBuilder, page, size uint64) (*Pagination[T], error) {
	spanCtx, span := newOTELSpan(ctx, "DB.QueryPage")
	defer span.End()
	count, err := s.Count(spanCtx, sb)
	if err != nil {
		return nil, errors.Wrap(err, "Count fail")
	}
	p := Pagination[T]{Page: page, Size: size}
	if count == 0 {
		return &p, nil
	}
	p.Total = count
	toSql, param, err := s.page(sb, page, size).ToSql()
	if err != nil {
		log.Debug("QueryPage to toSql fail err=", err.Error())
		return nil, errors.Wrap(err, "QueryPage toSql fail")
	}
	query, err := s.Db.QueryxContext(spanCtx, toSql, param...)
	if err != nil {
		if err == sql.ErrNoRows {
			return &p, nil
		}
		log.Debug("QueryxContext fail err=", err.Error())
		return nil, errors.Wrap(err, "QueryxContext fail")
	}
	ts := make([]*T, 0)
	for query.Next() {
		var d T
		err = query.StructScan(&d)
		if err != nil {
			log.Debug("StructScan fail err=", err.Error())
			return nil, errors.Wrap(err, "StructScan fail")
		}
		ts = append(ts, &d)
	}
	p.Data = ts
	return &p, nil
}

type ToSql struct {
	Columns string
}

func (t *ToSql) ToSql() (sql string, args []interface{}, err error) {
	return t.Columns, nil, nil
}
func (s *BaseStore[T]) Count(ctx context.Context, sb squirrel.SelectBuilder) (uint64, error) {
	toSql, i, err := builder.Set(sb, "Columns", &ToSql{Columns: "count(*)"}).(squirrel.SelectBuilder).ToSql()
	if err != nil {
		log.Debug("Count toSql fail err=", err.Error())
		return 0, errors.Wrap(err, "Count toSql fail")
	}
	rowxContext := s.Db.QueryRowxContext(ctx, toSql, i...)
	if rowxContext.Err() != nil {
		log.Debug("QueryRowxContext fail err=", err.Error())
		return 0, errors.Wrap(err, "QueryRowxContext fail")
	}
	var count uint64
	err = rowxContext.Scan(&count)
	if err != nil {
		log.Debug("Scan fail err=", err.Error())
		return 0, errors.Wrap(err, "Scan fail")
	}
	return count, nil
}
func (s *BaseStore[T]) page(sb squirrel.SelectBuilder, page, size uint64) squirrel.SelectBuilder {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 50
	}
	offset := size * (page - 1)
	return sb.Limit(size).Offset(offset)
}
func newOTELSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	ctx, span := otel.Tracer(otelName).Start(ctx, name)
	span.SetAttributes(semconv.DBSystemMySQL)
	return ctx, span
}
