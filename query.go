package unicorm

import "strings"

type QueryBuilder struct{}
type query struct{ statement string }
type _select struct{ statement string }
type columns struct{ statement string }
type from struct{ statement string }
type where struct{ statement string }

func NewQueryBuilder() *query {
	return &query{}
}
func (q *QueryBuilder) New() *query {
	return &query{}
}
func (q *query) Select() *_select {
	q.statement = "SELECT"
	return &_select{q.statement}
}
func (s *_select) Columns(fieldNames ...string) *columns {
	s.statement = s.statement + " " + strings.Join(fieldNames, ", ")
	return &columns{s.statement}
}
func (c *columns) From(tableNames ...string) *from {
	c.statement = c.statement + " FROM " + strings.Join(tableNames, ", ")
	return &from{c.statement}
}
func (f *from) Execute() string {
	return f.statement + ";"
}
func (f *from) Where(conditional Conditional) *where {
	f.statement = f.statement + " WHERE " + conditional.String()
	return &where{f.statement}
}
func (w *where) Execute() string {
	return w.statement + ";"
}
func (w *where) AndWhere(conditional Conditional) *where {
	w.statement = w.statement + " AND " + conditional.String()
	return &where{w.statement}
}
func (w *where) OrWhere(conditional Conditional) *where {
	w.statement = w.statement + " OR " + conditional.String()
	return &where{w.statement}
}
