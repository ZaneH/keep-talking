package entities

import (
	"errors"
	"strings"

	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type BombFace struct {
	ModulesByPosition map[valueobject.ModulePosition]uuid.UUID
}

func NewBombFace() *BombFace {
	return &BombFace{
		ModulesByPosition: make(map[valueobject.ModulePosition]uuid.UUID),
	}
}

func (f *BombFace) AddModuleAt(module Module, position valueobject.ModulePosition) error {
	if _, exists := f.ModulesByPosition[position]; exists {
		return errors.New("position already occupied")
	}

	f.ModulesByPosition[position] = module.GetModuleID()
	return nil
}

func (f *BombFace) String() string {
	var result strings.Builder
	for position, moduleID := range f.ModulesByPosition {
		result.WriteString(position.String() + ": " + moduleID.String() + "\n")
	}
	return result.String()
}
