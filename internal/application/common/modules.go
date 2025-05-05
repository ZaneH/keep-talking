package common

type Module interface {
	IsSolved() bool
}

type ModuleState struct {
	IsSolved bool
}
