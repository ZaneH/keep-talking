package entities

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

const WIRE_POSITIONS = 6

type SimpleWiresState struct {
	ModuleState
	SolutionIndices []int
	Wires           []valueobject.SimpleWire
}

func NewSimpleWiresState() SimpleWiresState {
	var wires []valueobject.SimpleWire
	var solution []int
	for {
		wires = generateRandomWires()
		solution = buildSolution(wires)
		if len(solution) > 0 {
			break
		}
	}

	return SimpleWiresState{
		Wires:           wires,
		SolutionIndices: solution,
	}
}

type SimpleWiresModule struct {
	ModuleID uuid.UUID
	State    SimpleWiresState
}

func NewSimpleWiresModule() *SimpleWiresModule {
	return &SimpleWiresModule{
		ModuleID: uuid.New(),
		State:    NewSimpleWiresState(),
	}
}

func (m *SimpleWiresModule) GetModuleID() uuid.UUID {
	return m.ModuleID
}

func (m *SimpleWiresModule) GetType() valueobject.ModuleType {
	return valueobject.SIMPLE_WIRES
}

func (m *SimpleWiresModule) SetState(state SimpleWiresState) {
	m.State = state
}

func (m *SimpleWiresModule) String() string {
	var result string = "\n"
	for i, wire := range m.State.Wires {
		if wire.IsCut {
			result += fmt.Sprintf("Wire %d: %s (cut)\n", i, wire.WireColor)
		} else {
			result += fmt.Sprintf("Wire %d: %s\n", i, wire.WireColor)
		}
	}

	return result
}

func buildSolution(wires []valueobject.SimpleWire) []int {
	solution := make([]int, 0)

	for i, wire := range wires {
		// TODO: Implement a guide to determine the solution
		// assume 'Red' is the solution for now
		if wire.WireColor == valueobject.SimpleWireColors[0] {
			solution = append(solution, i)
		}
	}

	return solution
}

func generateRandomWires() []valueobject.SimpleWire {
	wires := make([]valueobject.SimpleWire, WIRE_POSITIONS)

	for i := 0; i < WIRE_POSITIONS; i++ {
		if n := rand.Intn(100); n < 40 {
			continue
		}

		color := valueobject.SimpleWireColors[i%len(valueobject.SimpleWireColors)]
		wires[i] = valueobject.SimpleWire{
			WireColor: color,
		}
	}

	return wires
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

	isInSolution := false
	for _, index := range m.State.SolutionIndices {
		if index == wireIndex {
			isInSolution = true
			break
		}
	}

	if !isInSolution {
		return errors.New("wire was not meant to be cut")
	}

	if m.IsSolved() {
		m.State.MarkSolved = true
	}

	return nil
}
