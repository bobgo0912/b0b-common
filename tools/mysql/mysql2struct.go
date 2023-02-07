package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var pool *sqlx.DB
var config Config

func main() {
	flag.Parse()
	configPath := flag.Arg(0)
	file, e := os.Open(configPath)
	if e != nil {
		panic(fmt.Errorf("open fail fial,dirct:%s(%w)", configPath, e))
	}
	bytes, _ := ioutil.ReadAll(file)

	_ = json.Unmarshal(bytes, &config)
	var err error
	pool, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@%s(%s)/%s", config.Username, config.Password, config.Protocol, config.Address, config.Dbname))
	if err != nil {
		panic(err)
	}

	_ = os.MkdirAll(config.OutputDir+"/"+config.Dbname, os.ModePerm)
	tables := queryTables(config.TableName)

	builder := strings.Builder{}
	for _, table := range tables {
		//转换表名
		builder.Reset()
		packageTime := false
		//packageSql := false
		tableName := table.TableName
		//titleTable := CamelStr(tableName)
		titleTable := FirstUpper(tableName)
		builder.WriteString(fmt.Sprintf("//%s %s\ntype %s struct {\n", titleTable, table.TableComment, titleTable))
		//拼接字符串
		for _, column := range queryColumns(tableName) {
			//转换列名
			dataType := strings.ToUpper(column.DataType)
			value, ok := DataTypeMap[dataType]
			if ok {
				if column.IsNullable == "YES" {
					dataType = value[0]
					//	packageSql = true
				} else {
					dataType = value[0]
				}
				//是否需要 sql 包
				packageTime = dataType == "time.Time"

			} else {
				dataType = "string"
			}
			//拼接字符串
			camelStr := CamelStr(column.ColumnName)
			builder.WriteString(fmt.Sprintf("	%s %s `db:\"%s\" json:\"%s\"` //%s", camelStr, dataType, column.ColumnName, strings.ToLower(string(camelStr[0]))+camelStr[1:], column.ColumnComment))
			if column.ColumnKey != "" {
				builder.WriteString("(" + column.ColumnKey + ")")
			}
			builder.WriteString("\n")
		}
		builder.WriteString("}\n")
		fileStr := "package " + config.OutputPackage + "\nimport ("
		fileStr += "\"github.com/jmoiron/sqlx\"\n"
		fileStr += "\t\"github.com/bobgo0912/b0b-common/pkg/sql\"\n"

		//if packageSql {
		//	fileStr += "\"database/sql\"\n"
		//}
		if packageTime {
			fileStr += "\"time\"\n"
		}
		fileStr += ")\n\n"
		fileStr += "const " + titleTable + "TableName = \"" + tableName + "\"\n"

		titleDb := FirstUpper(config.Dbname)
		builder.WriteString("\ntype " + titleTable + "Store struct {\n\t*sql.BaseStore[" + titleTable + "]\n}\n\nfunc GetConnection() (*sqlx.DB, error) {\n\tif sql." + titleDb + "Db != nil {\n\t\treturn sql." + titleDb + "Db, nil\n\t}\n\tvar err error\n\tsql." + titleDb + "Db, err = sql.Db(\"" + config.Dbname + "\", nil)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn sql." + titleDb + "Db, nil\n}\n\nfunc Get" + titleTable + "Store() (*" + titleTable + "Store, error) {\n\tconnection, err := GetConnection()\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn &" + titleTable + "Store{&sql.BaseStore[" + titleTable + "]{Db: connection, TableName: " + titleTable + "TableName}}, nil\n}\n")
		fileStr += builder.String()

		join := path.Join(config.OutputDir, config.Dbname, tableName+".go")
		_ = os.WriteFile(join, []byte(fileStr), os.ModePerm)
	}
	_ = os.Chdir(config.OutputDir + "/" + config.Dbname)
	cmd := exec.Command("go", "fmt")
	out, e := cmd.CombinedOutput()
	if e != nil {
		panic(e)
	}
	fmt.Printf("格式化结果:\n%s\n", string(out))
}

// 查询所有的列
func queryColumns(tableName string) []Column {
	var results []Column
	e := pool.Select(&results, "select COLUMN_NAME,IS_NULLABLE,DATA_TYPE,COLUMN_KEY,COLUMN_COMMENT from information_schema.COLUMNS where TABLE_SCHEMA = ? and TABLE_NAME = ?", config.Dbname, tableName)
	if e != nil {
		panic(e)
	}
	return results
}

// 查询所有的表
func queryTables(tableName string) []Table {
	var tables []Table
	sql := "SELECT table_name ,table_comment FROM information_schema.TABLES WHERE table_schema = '" + config.Dbname + "'"
	if tableName != "" {
		sql += " and table_name = '" + tableName + "'"
	}
	sql += " ORDER BY table_name"
	e := pool.Select(&tables, sql)
	if e != nil {
		panic(e)
	}
	return tables
}

func CamelStr(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

type Column struct {
	ColumnName    string `db:"COLUMN_NAME"`
	IsNullable    string `db:"IS_NULLABLE"`
	DataType      string `db:"DATA_TYPE"`
	ColumnKey     string `db:"COLUMN_KEY"`
	ColumnComment string `db:"COLUMN_COMMENT"`
}

type Config struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Protocol      string `json:"protocol"`
	Address       string `json:"address"`
	Dbname        string `json:"dbname"`
	TableName     string `json:"tableName"`
	OutputDir     string `json:"outputDir"`
	OutputPackage string `json:"outputPackage"`
}

var DataTypeMap = map[string][]string{
	//整型
	"TINYINT":   {"int", "sql.NullInt64"},
	"SMALLINT":  {"int", "sql.NullInt64"},
	"MEDIUMINT": {"int", "sql.NullInt64"},
	"INT":       {"int", "sql.NullInt64"},
	"INTEGER":   {"int", "sql.NullInt64"},
	"BIGINT":    {"uint64", "sql.NullInt64"},
	//浮点数
	"FLOAT":   {"float64", "sql.NullFloat64"},
	"DOUBLE":  {"float64", "sql.NullFloat64"},
	"DECIMAL": {"float64", "sql.NullFloat64"},
	//时间
	"DATE":      {"time.Time", "sql.NullTime"},
	"TIME":      {"time.Time", "sql.NullTime"},
	"YEAR":      {"time.Time", "sql.NullTime"},
	"DATETIME":  {"time.Time", "sql.NullTime"},
	"TIMESTAMP": {"time.Time", "sql.NullTime"},
	//字符串
	"CHAR":       {"string", "sql.NullString"},
	"VARCHAR":    {"string", "sql.NullString"},
	"TINYBLOB":   {"string", "sql.NullString"},
	"TINYTEXT":   {"string", "sql.NullString"},
	"BLOB":       {"string", "sql.NullString"},
	"TEXT":       {"string", "sql.NullString"},
	"MEDIUMBLOB": {"string", "sql.NullString"},
	"MEDIUMTEXT": {"string", "sql.NullString"},
	"LONGBLOB":   {"string", "sql.NullString"},
	"LONGTEXT":   {"string", "sql.NullString"},
	"JSON":       {"string", "sql.NullString"},
}

type Table struct {
	TableName    string `db:"table_name"`
	TableComment string `db:"table_comment"`
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
