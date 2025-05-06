package entities

import (
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type SimpleWiresState struct {
	common.ModuleState
	SolutionIndices []int
	Wires           []valueobject.SimpleWire
}

type SimpleWiresModule struct {
	ModuleID uuid.UUID
	State    SimpleWiresState
}

func NewSimpleWiresModule(wires []valueobject.SimpleWire, solutionIndices []int) *SimpleWiresModule {
	return &SimpleWiresModule{
		ModuleID: uuid.New(),
		State: SimpleWiresState{
			Wires:           wires,
			SolutionIndices: solutionIndices,
		},
	}
}

func (m *SimpleWiresModule) IsSolved() bool {
	for _, index := range m.State.SolutionIndices {
		if index < 0 || index >= len(m.State.Wires) {
			return false
		}
		if !m.State.Wires[index].IsCut {
			return false
		}
	}
	return true
}

func (m *SimpleWiresModule) CutWire(wireIndex int) error {
	wire := &m.State.Wires[wireIndex]
	if wire.IsCut {
		return errors.New("wire already cut")
	}

	wire.IsCut = true
	return nil
}
