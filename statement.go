package unicorm

import "strings"

type StatementBuilder struct{}
type _statement struct{ statement string }
type insert struct{ statement string }
type update struct{ statement string }
type statementColumns struct{ statement string }
type into struct{ statement string }
type statementWhere struct{ statement string }
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
func (s *_statement) Update() *update {
	return &update{s.statement}
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
