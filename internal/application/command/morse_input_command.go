package command

import (
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

type MorseChangeFrequencyCommand struct {
	BaseModuleInputCommand
	Direction valueobject.IncrementDecrement
}

type MorseTxCommand struct {
	BaseModuleInputCommand
}

type MorseCommandResult struct {
	BaseModuleInputCommandResult
	DisplayedFrequency float32
}
