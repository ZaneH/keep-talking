package entities

import (
	"errors"

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

func (f *BombFace) AddModule(module Module, position valueobject.ModulePosition) error {
	if _, exists := f.ModulesByPosition[position]; exists {
		return errors.New("position already occupied")
	}

	f.ModulesByPosition[position] = module.GetModuleID()
	return nil
}
