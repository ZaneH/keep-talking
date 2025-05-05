package command

import (
	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

type PasswordLetterChangeCommand struct {
	BaseModuleInputCommand
	LetterIndex int
	Direction   valueobject.IncrementDecrement
}

type PasswordSubmitCommand struct {
	BaseModuleInputCommand
}

type PasswordInputResult struct {
	Result *common.InputResult
}
