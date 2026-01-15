package command

import (
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
)

type PasswordLetterChangeCommand struct {
	BaseModuleInputCommand
	LetterIndex int
	Direction   valueobject.IncrementDecrement
}

type PasswordSubmitCommand struct {
	BaseModuleInputCommand
}

type PasswordCommandResult struct {
	BaseModuleInputCommandResult
	Letters string
}
