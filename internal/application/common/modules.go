package common

type Module interface {
	IsSolved() bool
}

type ModuleState struct {
	MarkSolved bool
}
