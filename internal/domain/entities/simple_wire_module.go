package entities

import (
	"errors"

	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type SimpleWireModule struct {
	ModuleID        uuid.UUID
	SolutionIndices []int
	Wires           []valueobject.SimpleWire
}

func NewSimpleWireModule(wires []valueobject.SimpleWire, solutionIndices []int) *SimpleWireModule {
	return &SimpleWireModule{
		SolutionIndices: solutionIndices,
		Wires:           wires,
	}
}

func (m *SimpleWireModule) IsSolved() bool {
	for _, index := range m.SolutionIndices {
		if index < 0 || index >= len(m.Wires) {
			return false
		}
		if !m.Wires[index].IsCut {
			return false
		}
	}
	return true
}

func (m *SimpleWireModule) CutWire(wireIndex int) (bool, error) {
	wire := &m.Wires[wireIndex]
	if wire.IsCut {
		return false, errors.New("wire already cut")
	}

	wire.IsCut = true
	return true, nil
}
