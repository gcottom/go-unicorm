package unicorm

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"unicode"

	"github.com/lodmev/scanstruct"
)

func Parse[T any](functionName string, tableName string, strct T) (string, string, error) {
	var functionType string
	fieldList := getFieldList(strct)
	parts := splitParts(functionName)
	if len(parts) < 1 {
		return "", "", errors.New("invalid function name")
	}
	_, ok := subjectMap[parts[0]]
	if !ok {
		return "", "", errors.New("invalid function name: missing subject")
	}
	if slices.Contains(querySubjects, parts[0]) {
		functionType = "QUERY"
	} else if slices.Contains(statementSubjects, parts[0]) {
		functionType = "STATEMENT"
	} else {
		return "", "", errors.New("invalid function name: can't determine if query or statement")
	}
	if functionType == "QUERY" {
		query, err := queryParse(parts[1:], fieldList, tableName)
		if err != nil {
			return "", "", err
		}
		return functionType, query, nil
	}
	return functionType, "", nil
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
