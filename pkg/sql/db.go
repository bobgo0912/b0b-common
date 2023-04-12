package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/bobgo0912/b0b-common/pkg/config"
	"github.com/bobgo0912/b0b-common/pkg/constant"
	"github.com/bobgo0912/b0b-common/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lann/builder"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"strings"
)

const otelName = "b0b-common/mysql"

var (
	EduDb   *sqlx.DB
	TestDb  *sqlx.DB
	OrderDb *sqlx.DB
)

type BaseStore[T any] struct {
	Db        *sqlx.DB
	TableName string
}

func Db(dbname string, c *config.MysqlCfg) (*sqlx.DB, error) {
	if c == nil {
		c = config.Cfg.MysqlCfg[dbname]
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		c.Username, c.Password, c.Host, c.Port, dbname,
	)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Connect fail")
	}
	return db, nil
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
	span.SetAttributes(
		attribute.KeyValue{
			Key:   semconv.DBStatementKey,
			Value: attribute.StringValue(toSql),
		}, attribute.KeyValue{
			Key:   constant.DBParamKey,
			Value: attribute.StringValue(fmt.Sprint(param)),
		},
	)
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
func (s *BaseStore[T]) QueryByCondition(ctx context.Context, sb squirrel.SelectBuilder, columns ...string) (*T, error) {
	spanCtx, span := newOTELSpan(ctx, "DB.QueryByCondition")
	defer span.End()
	if len(columns) < 1 {
		columns = append(columns, "*")
	}

	toSql, param, err := sb.From(s.TableName).ToSql()
	if err != nil {
		log.Debug("QueryByCondition to sql fail err=", err.Error())
		return nil, errors.Wrap(err, "squirrel toSql fail")
	}
	span.SetAttributes(
		attribute.KeyValue{
			Key:   semconv.DBStatementKey,
			Value: attribute.StringValue(toSql),
		},
		attribute.KeyValue{
			Key:   constant.DBParamKey,
			Value: attribute.StringValue(fmt.Sprint(param)),
		},
	)
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
	sb = sb.From(s.TableName)
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
	span.SetAttributes(
		attribute.KeyValue{
			Key:   semconv.DBStatementKey,
			Value: attribute.StringValue(toSql),
		}, attribute.KeyValue{
			Key:   constant.DBParamKey,
			Value: attribute.StringValue(fmt.Sprint(param)),
		})
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
func (s *BaseStore[T]) QueryList(ctx context.Context, sb squirrel.SelectBuilder) ([]*T, error) {
	spanCtx, span := newOTELSpan(ctx, "DB.QueryList")
	defer span.End()
	from := sb.From(s.TableName)
	toSql, param, err := from.ToSql()
	span.SetAttributes(
		attribute.KeyValue{
			Key:   semconv.DBStatementKey,
			Value: attribute.StringValue(toSql),
		},
		attribute.KeyValue{
			Key:   constant.DBParamKey,
			Value: attribute.StringValue(fmt.Sprint(param)),
		})
	query, err := s.Db.QueryxContext(spanCtx, toSql, param...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
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
	return ts, nil
}

type ToSql struct {
	pred interface{}
	args []interface{}
}

func (t ToSql) ToSql() (string, []interface{}, error) {
	return t.pred.(string), nil, nil
}
func (s *BaseStore[T]) Count(ctx context.Context, sb squirrel.SelectBuilder) (uint64, error) {
	t := ToSql{pred: "count(*)"}

	set := builder.Set(sb, "Columns", []squirrel.Sqlizer{t})
	selectBuilder := set.(squirrel.SelectBuilder)
	toSql, i, err := selectBuilder.ToSql()
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

func (s *BaseStore[T]) Update() {

}

func (s *BaseStore[T]) Tx(ctx context.Context) (*sql.Tx, error) {
	return s.Db.BeginTx(ctx, nil)
}

func (s *BaseStore[T]) MultipleInsert(ctx context.Context, datas []*T) error {
	if len(datas) < 1 {
		return nil
	}
	_, span := newOTELSpan(ctx, "DB.MultipleInsert")
	defer span.End()

	insertBuilder := s.MultipleInsertBuilder(datas)
	toSql, param, err := insertBuilder.ToSql()
	if err != nil {
		log.Error("ToSql fail err=", err)
		return errors.Wrap(err, "ToSql fail")
	}
	_, err = s.Db.Exec(toSql, param...)
	if err != nil {
		log.Error("Exec fail err=", err)
		return errors.Wrap(err, "Exec fail")
	}
	return nil
}

func (s *BaseStore[T]) MultipleInsertBuilder(datas []*T) squirrel.InsertBuilder {
	r := reflect.ValueOf(*datas[0])
	t := r.Type()
	columns := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		get := t.Field(i).Tag.Get("db")
		if get == "" {
			continue
		}
		columns = append(columns, fmt.Sprintf("`%s`", get))
	}
	join := strings.Join(columns, ",")
	insertBuilder := squirrel.Insert(fmt.Sprintf("`%s`", s.TableName)).Columns(join)
	for _, data := range datas {
		of := reflect.ValueOf(*data)
		its := make([]interface{}, 0)
		for i := 0; i < of.NumField(); i++ {
			its = append(its, of.Field(i).Interface())
		}
		insertBuilder = insertBuilder.Values(its...)
	}
	return insertBuilder
}
func (s *BaseStore[T]) InsertBuilder(data *T) squirrel.InsertBuilder {
	r := reflect.ValueOf(*data)
	t := r.Type()
	columns := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		get := t.Field(i).Tag.Get("db")
		if get == "" {
			continue
		}
		columns = append(columns, fmt.Sprintf("`%s`", get))
	}
	join := strings.Join(columns, ",")
	insertBuilder := squirrel.Insert(fmt.Sprintf("`%s`", s.TableName)).Columns(join)
	of := reflect.ValueOf(*data)
	its := make([]interface{}, 0)
	for i := 0; i < of.NumField(); i++ {
		its = append(its, of.Field(i).Interface())
	}
	insertBuilder = insertBuilder.Values(its...)

	return insertBuilder
}
func (s *BaseStore[T]) Insert(ctx context.Context, data *T) error {
	if data == nil {
		return nil
	}
	_, span := newOTELSpan(ctx, "DB.Insert")
	defer span.End()

	insertBuilder := s.InsertBuilder(data)
	toSql, param, err := insertBuilder.ToSql()
	if err != nil {
		log.Error("ToSql fail err=", err)
		return errors.Wrap(err, "ToSql fail")
	}
	_, err = s.Db.Exec(toSql, param...)
	if err != nil {
		log.Error("Exec fail err=", err)
		return errors.Wrap(err, "Exec fail")
	}
	return nil
}
