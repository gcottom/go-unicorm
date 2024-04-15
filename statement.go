package unicorm

import (
	"reflect"
	"strings"

	"github.com/lodmev/scanstruct"
)

type StatementBuilder struct{}
type _statement struct{ statement string }
type insert struct{ statement string }
type update struct{ statement string }
type set struct {
	statement string
	args      []any
}
type statementColumns struct{ statement string }
type into struct{ statement string }
type statementWhere struct {
	statement string
	args      []any
}
type values struct {
	statement string
	args      []any
}

func NewStatementBuilder() *_statement {
	return &_statement{}
}
func (s *_statement) Insert() *insert {
	s.statement = "INSERT"
	return &insert{s.statement}
}
func (s *_statement) Update(tableName string) *update {
	s.statement = "UPDATE " + tableName
	return &update{s.statement}
}
func (u *update) Set(isStruct bool, args ...any) *set {
	u.statement = u.statement + " SET"
	if isStruct {
		structArgs := []any{}
		for i := 0; i < reflect.TypeOf(args[0]).NumField(); i++ {
			u.statement = u.statement + " " + scanstruct.ToSnakeCase(reflect.ValueOf(args[0]).Type().Field(i).Name) + "=?,"
			structArgs = append(structArgs, reflect.ValueOf(args[0]).Field(i).Interface())
		}
		args = structArgs
	}
	u.statement = u.statement[:len(u.statement)-1]
	return &set{statement: u.statement, args: args}
}
func (s *set) Where(conditional Conditional) *statementWhere {
	s.statement = s.statement + " WHERE " + conditional.String()
	return &statementWhere{statement: s.statement, args: s.args}
}
func (w *statementWhere) Execute() (string, []any) {
	return w.statement + ";", w.args
}

func (i *insert) Into(tableName string) *into {
	i.statement = i.statement + " INTO " + tableName
	return &into{i.statement}
}
func (i *into) Values(args []any) *values {
	argC := len(args)
	i.statement = i.statement + " VALUES (" + strings.Repeat("?,", argC-1) + "?);"
	return &values{statement: i.statement, args: args}
}
func (v *values) Execute() (string, []any) {
	return v.statement, v.args
}
