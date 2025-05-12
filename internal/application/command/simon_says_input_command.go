package command

import "github.com/ZaneH/keep-talking/internal/domain/valueobject"

type SimonSaysInputCommand struct {
	BaseModuleInputCommand
	Color valueobject.Color
}

type SimonSaysInputCommandResult struct {
	BaseModuleInputCommandResult
	NextColorSequence []valueobject.Color
}
