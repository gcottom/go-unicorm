package unicorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/lodmev/scanstruct"
)

func InitTable[T, V any](superTable *T, entity V, db *sql.DB) {
	tableName := strings.ToLower(reflect.TypeOf(entity).Name())
	table := new(Table[V])
	table.TableName = tableName
	table.RepoStruct = entity
	table.DBClient = db
	table.structInfo = getFieldInfo(entity)
	reflect.ValueOf(superTable).Elem().FieldByName("Table").Set(reflect.ValueOf(table))
	setUpTable(table)
}

func getFieldInfo[T any](strct T) []structInfo {
	fieldList := make([]structInfo, 0)
	for i := 0; i < reflect.ValueOf(strct).NumField(); i++ {
		info := structInfo{
			FieldName: scanstruct.ToSnakeCase(reflect.ValueOf(strct).Type().Field(i).Name),
			FieldType: reflect.ValueOf(strct).Field(i).Type().Name(),
		}
		fieldList = append(fieldList, info)
	}
	return fieldList
}

func setUpTable[T any](table *Table[T]) {
	checkCreateDBTable(table)
	checkDBColumnsExist(table)
}
func checkCreateDBTable[T any](table *Table[T]) {
	vals := []string{}
	for i := 0; i < reflect.ValueOf(table.RepoStruct).NumField(); i++ {
		fn := scanstruct.ToSnakeCase(reflect.ValueOf(table.RepoStruct).Type().Field(i).Name)
		ft := ""
		switch reflect.ValueOf(table.RepoStruct).Field(i).Type().Name() {
		case "string":
			ft = "TEXT"
		case "int":
			ft = "INTEGER"
		case "bool":
			ft = "BOOLEAN"
		default:
			ft = "TEXT"
		}
		if fn == "id" {
			ft = ft + " NOT NULL PRIMARY KEY"
		}
		vals = append(vals, fmt.Sprintf("%s %s", fn, ft))

	}
	createStatement := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table.TableName, strings.Join(vals, ", "))
	table.DBClient.Exec(createStatement)
}

func checkDBColumnsExist[T any](table *Table[T]) {
	checkTableColExists := `SELECT COUNT(*) FROM pragma_table_info(?) WHERE name=?`
	var colExists bool
	var colsToAdd []string
	for i := 0; i < reflect.ValueOf(table.RepoStruct).NumField(); i++ {
		fn := scanstruct.ToSnakeCase(reflect.ValueOf(table.RepoStruct).Type().Field(i).Name)
		r, err := table.DBClient.Query(checkTableColExists, table.TableName, fn)
		if err != nil {
			panic(err)
		}
		for r.Next() {
			r.Scan(&colExists)
		}
		if !colExists {
			ft := ""
			switch reflect.ValueOf(table.RepoStruct).Field(i).Type().Name() {
			case "string":
				ft = "TEXT"
			case "int":
				ft = "INTEGER"
			case "bool":
				ft = "BOOLEAN"
			default:
				ft = "TEXT"
			}
			if fn == "id" {
				ft = ft + " NOT NULL PRIMARY KEY"
			}
			colsToAdd = append(colsToAdd, "ADD COLUMN "+fn+" "+ft)
		}
	}
	for _, cta := range colsToAdd {
		alterStatement := fmt.Sprintf("ALTER TABLE '%s' %s;", table.TableName, cta)
		table.DBClient.Exec(alterStatement)
	}
}

func (repo *Table[T]) AutoGenerate(val ...any) ([]T, error) {
	c, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(c).Name()
	sp := strings.Split(name, ".")
	function := sp[len(sp)-1]
	functiontype, q, args, err := Parse(function, repo, val...)
	if err != nil {
		panic(err)
	}
	fmt.Println(functiontype, q)
	if functiontype == QUERY {
		return ExecuteQuery(q, repo, val)
	}
	return ExecuteStatement(q, repo, args)
}

func ExecuteStatement[T any](statement string, r *Table[T], args []any) ([]T, error) {
	db := r.DBClient
	_, err := db.Exec(statement, args...)
	if err != nil {
		panic(err)
	}
	return nil, nil
}

func ExecuteQuery[T any](query string, r *Table[T], args []any) ([]T, error) {
	db := r.DBClient
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	results := []T{}
	var t T

	scanstruct.NameMapper = scanstruct.ToSnakeCase
	for rows.Next() {
		if err := scanstruct.Scan(&t, rows); err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, nil
}

func isCreateQuery(function string, val string) (bool, string) {
	v := make([]string, 2)
	if val == "query" {
		v[0] = "Get"
		v[1] = "get"
	} else {
		v[0] = "Save"
		v[1] = "save"
	}
	isQ := strings.HasPrefix(function, v[0]) || strings.HasPrefix(function, v[1])
	if isQ {
		a, f := strings.CutPrefix(function, v[0])
		if !f {
			a, _ = strings.CutPrefix(function, v[1])
		}
		return true, a
	}
	return false, function
}
