package command

import (
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

type KeypadInputCommand struct {
	BaseModuleInputCommand
	Symbol valueobject.Symbol
}

type KeypadInputCommandResult struct {
	BaseModuleInputCommandResult
	ActivatedSymbols map[valueobject.Symbol]bool
	DisplayedSymbols []valueobject.Symbol
}
