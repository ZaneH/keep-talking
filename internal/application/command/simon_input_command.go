package command

import "github.com/ZaneH/keep-talking/internal/domain/valueobject"

type SimonInputCommand struct {
	BaseModuleInputCommand
	Color valueobject.Color
}

type SimonInputCommandResult struct {
	BaseModuleInputCommandResult
	// Indicates if the current iteration of the sequence is finished.
	HasFinishedSeq  bool
	DisplaySequence []valueobject.Color
}
