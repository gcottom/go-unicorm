package unicorm

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/lodmev/scanstruct"
)

func Parse[T any](functionName string, table *Table[T], val ...any) (string, string, []any, error) {
	var functionType string
	fieldList := getFieldList(table.RepoStruct)
	parts := splitParts(functionName)
	if len(parts) < 1 {
		return "", "", nil, errors.New("invalid function name")
	}
	_, ok := subjectMap[parts[0]]
	if !ok {
		return "", "", nil, errors.New("invalid function name: missing subject")
	}
	if slices.Contains(querySubjects, parts[0]) {
		functionType = QUERY
	} else if slices.Contains(statementSubjects, parts[0]) {
		functionType = STATEMENT
	} else {
		return "", "", nil, errors.New("invalid function name: can't determine if query or statement")
	}
	if functionType == QUERY {
		query, err := queryParse(parts[1:], fieldList, table.TableName)
		if err != nil {
			return "", "", nil, err
		}
		return functionType, query, nil, nil
	}
	if functionType == STATEMENT {
		statement, args, err := statementParse[T](parts, fieldList, table, val...)
		if err != nil {
			return "", "", nil, err
		}
		return functionType, statement, args, nil
	}
	return functionType, "", nil, nil
}
func splitParts(a string) []string {
	w := []string{}
	last := 0
	for i := 0; i < len(a); i++ {
		if unicode.IsUpper(rune(a[i])) {
			if last != 0 || i > 0 {
				w = append(w, a[last:i])
				last = i
			}
		}
	}
	w = append(w, a[last:])
	return w
}
func getFieldList[T any](strct T) []string {
	fieldList := make([]string, 0)
	for i := 0; i < reflect.ValueOf(strct).NumField(); i++ {
		fieldList = append(fieldList, scanstruct.ToSnakeCase(reflect.ValueOf(strct).Type().Field(i).Name))
	}
	return fieldList
}

func statementParse[T any](parts []string, fieldList []string, table *Table[T], val ...any) (string, []any, error) {
	if len(parts) < 2 {
		if parts[0] == SAVE {
			if len(val) == 1 {
				if w, ok := val[0].(T); ok {
					switch reflect.ValueOf(w).Field(0).Type().Kind() {
					case reflect.String:
						if reflect.ValueOf(w).Field(0).IsZero() {
							id := uuid.New().String()
							reflect.ValueOf(&w).Elem().Field(0).SetString(id)
							var args []any
							for i := 0; i < reflect.ValueOf(w).NumField(); i++ {
								args = append(args, reflect.ValueOf(w).Field(i).Interface())
							}
							fmt.Println(args)
							stmt, args := NewStatementBuilder().Insert().Into(table.TableName).Values(args).Execute()
							return stmt, args, nil

						}
						r, err := table.DBClient.Query(NewQueryBuilder().Select().Columns("COUNT(id)").From(table.TableName).Where(NewConditional("id", "?", "=")).Execute(), reflect.ValueOf(w).Field(0).Interface())
						if err != nil {
							panic(err)
						}
						var exists bool
						for r.Next() {
							r.Scan(&exists)
						}
						if exists {
							stmt, args := NewStatementBuilder().Update(table.TableName).Set(true, w).Where(NewConditional("id", "?", "=")).Execute()
							args = append(args, reflect.ValueOf(w).Field(0).Interface())
							return stmt, args, nil
						}
					case reflect.Int:
						if reflect.ValueOf(w).Field(0).IsZero() {
							r, err := table.DBClient.Query(NewQueryBuilder().Select().Columns("TOP 1 id").From(table.TableName + " ORDER BY id DESC").Execute())
							if err != nil {
								panic(err)
							}
							var tid int64
							for r.Next() {
								r.Scan(&tid)
							}
							tid++
							reflect.ValueOf(&w).Elem().Field(0).SetInt(tid)
							var args []any
							for i := 0; i < reflect.ValueOf(w).NumField(); i++ {
								args = append(args, reflect.ValueOf(w).Field(i).Interface())
							}
							stmt, args := NewStatementBuilder().Insert().Into(table.TableName).Values(args).Execute()
							return stmt, args, nil
						}
					}
				}
			}

		}
	}
	return "", nil, nil
}

func queryParse(parts []string, fieldList []string, tableName string) (string, error) {
	if len(parts) < 1 {
		return "", errors.New("query doens't have enough parts")
	}
	if !slices.Contains(queryJoiners, parts[0]) {
		if parts[0] == QUERY_SHORT_ALL {
			if len(parts) < 2 {
				//query is simple FindAll
				return NewQueryBuilder().Select().Columns("*").From(tableName).Execute(), nil
			}
			parts = parts[1:]
			if !slices.Contains(queryJoiners, parts[0]) {
				return "", errors.New("query missing query joiner")
			}
			if len(parts) > 1 {
				if parts[0] == QUERY_JOINER_BY {
					parts = parts[1:]
					if slices.Contains(fieldList, scanstruct.ToSnakeCase(strings.Join(parts, ""))) {
						return NewQueryBuilder().Select().Columns("*").From(tableName).Where(NewConditional(scanstruct.ToSnakeCase(parts[0]), "?", "=")).Execute(), nil
					}
					return "", errors.New("couldn't parse query")
				}
			}
		}
		return "", errors.New("query missing query joiner")
	}
	if parts[0] == QUERY_JOINER_BY {
		parts = parts[1:]
		if slices.Contains(fieldList, scanstruct.ToSnakeCase(strings.Join(parts, ""))) {
			return NewQueryBuilder().Select().Columns("*").From(tableName).Where(NewConditional(scanstruct.ToSnakeCase(strings.Join(parts, "")), "?", "=")).Execute(), nil
		}
		return "", errors.New("couldn't parse query")
	}
	return "", errors.New("couldn't parse query")
}

func getType(v any) reflect.Type {
	return reflect.ValueOf(v).Type()
}
