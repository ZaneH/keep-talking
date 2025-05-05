package module

type Module interface {
	Id() string
	IsSolved() bool
	ApplyInput(input any) error
}

func (m *Module) Id() string {
	return m.Id
}
