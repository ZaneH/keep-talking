package command

import "github.com/ZaneH/keep-talking/internal/domain/valueobject"

type BigButtonInputCommand struct {
	BaseModuleInputCommand
	PressType        valueobject.PressType
	ReleaseTimestamp int64
}

type BigButtonInputCommandResult struct {
	BaseModuleInputCommandResult
	StripColor *valueobject.Color
}
