package unicorm

import "fmt"

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
