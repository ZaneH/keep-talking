package command

import (
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
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
	DisplayedPattern     string
	DisplayedFrequency   float32
	SelectedFrequencyIdx int
}
