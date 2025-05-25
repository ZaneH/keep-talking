package command

import "github.com/ZaneH/keep-talking/internal/domain/valueobject"

type SimonSaysInputCommand struct {
	BaseModuleInputCommand
	Color valueobject.Color
}

type SimonSaysInputCommandResult struct {
	BaseModuleInputCommandResult
	// Indicates if the current iteration of the sequence is finished.
	HasFinishedSeq  bool
	DisplaySequence []valueobject.Color
}
