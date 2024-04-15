package unicorm

import "fmt"

type conditionalBuilder struct {
	c Conditional
}
type conditionalColumn struct {
	c Conditional
}
type conditionalOperator struct {
	c Conditional
}
type conditionalValue struct {
	c Conditional
}

type Conditional struct {
	P          any
	Q          any
	Comparator string
}

func NewConditional(p any, q any, comparator string) Conditional {
	return Conditional{P: p, Q: q, Comparator: comparator}
}
func (c Conditional) String() string {
	return fmt.Sprintf("%s %s %s", c.P, c.Comparator, c.Q)
}

func NewConditionalBuilder() *conditionalBuilder {
	return &conditionalBuilder{}
}

func (c *conditionalBuilder) Column(columnName string) *conditionalColumn {
	c.c.P = columnName
	return &conditionalColumn{c.c}
}

func (cc *conditionalColumn) Operator(operator string) *conditionalOperator {
	cc.c.Comparator = operator
	return &conditionalOperator{cc.c}
}

func (co *conditionalOperator) Value(value any) *conditionalValue {
	co.c.Q = value
	return &conditionalValue{co.c}
}

func (c *conditionalValue) Execute() Conditional {
	return c.c
}
