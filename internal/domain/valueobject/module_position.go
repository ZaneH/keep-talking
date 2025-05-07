package valueobject

import "fmt"

type ModulePosition struct {
	Row    int
	Column int
	Face   int
}

func (p ModulePosition) String() string {
	return "(" + fmt.Sprint(p.Row) + ", " + fmt.Sprint(p.Column) + ", " + fmt.Sprint(p.Face) + ")"
}
