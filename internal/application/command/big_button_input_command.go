package command

import "github.com/ZaneH/keep-talking/internal/domain/valueobject"

type BigButtonInputCommand struct {
	BaseModuleInputCommand
	Action valueobject.PressType
}

type BigButtonInputCommandResult struct {
	BaseModuleInputCommandResult
	StripColor valueobject.Color
}
