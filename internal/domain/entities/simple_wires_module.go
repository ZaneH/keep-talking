package entities

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/ZaneH/keep-talking/internal/application/helpers"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

const minWires = 3
const maxWires = 6

type SimpleWiresState struct {
	ModuleState
	Wires []valueobject.SimpleWire
}

func NewSimpleWiresState() SimpleWiresState {
	var wires = generateRandomWires()

	return SimpleWiresState{
		Wires: wires,
	}
}

type SimpleWiresModule struct {
	ModuleID uuid.UUID
	State    SimpleWiresState
	bomb     *Bomb
}

func NewSimpleWiresModule(bomb *Bomb) *SimpleWiresModule {
	return &SimpleWiresModule{
		ModuleID: uuid.New(),
		State:    NewSimpleWiresState(),
		bomb:     bomb,
	}
}

func (m *SimpleWiresModule) GetModuleID() uuid.UUID {
	return m.ModuleID
}

func (m *SimpleWiresModule) GetModuleState() ModuleState {
	return m.State.ModuleState
}

func (m *SimpleWiresModule) GetType() valueobject.ModuleType {
	return valueobject.SimpleWires
}

func (m *SimpleWiresModule) SetState(state SimpleWiresState) {
	m.State = state
}

func (m *SimpleWiresModule) GetBomb() *Bomb {
	return m.bomb
}

func (m *SimpleWiresModule) String() string {
	var result = "\n"
	for i, wire := range m.State.Wires {
		if wire.IsCut {
			result += fmt.Sprintf("Wire %d: %s (cut)\n", i, wire.WireColor)
		} else {
			result += fmt.Sprintf("Wire %d: %s\n", i, wire.WireColor)
		}
	}

	return result
}

func generateRandomWires() []valueobject.SimpleWire {
	nWires := rand.Intn(maxWires-minWires+1) + minWires
	wires := make([]valueobject.SimpleWire, nWires)

	for i := 0; i < nWires; i++ {
		color := valueobject.SimpleWireColors[i%len(valueobject.SimpleWireColors)]
		wires[i] = valueobject.SimpleWire{
			WireColor: color,
		}
	}

	return wires
}

func (m *SimpleWiresModule) cutSucceed() (bool, error) {
	m.State.MarkSolved = true
	return false, nil
}

func (m *SimpleWiresModule) CutWire(wireIndex int) (strike bool, err error) {
	if wireIndex < 0 || wireIndex >= len(m.State.Wires) {
		return false, errors.New("invalid wire index")
	}

	wire := &m.State.Wires[wireIndex]
	if wire.IsCut {
		return false, errors.New("wire already cut")
	}

	wire.IsCut = true

	if len(m.State.Wires) == 3 {
		// If there are no red wires, cut the second wire.
		if !hasColor(m.State.Wires, valueobject.Red) {
			if wireIndex == 1 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If the last wire is white, cut the last wire.
		if m.State.Wires[len(m.State.Wires)-1].WireColor == valueobject.White {
			if wireIndex == len(m.State.Wires)-1 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is more than one blue wire, cut the last blue wire.
		blueIdxs := colorIndecies(m.State.Wires, valueobject.Blue)
		if len(blueIdxs) > 1 {
			if wireIndex == blueIdxs[len(blueIdxs)-1] {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}
	} else if len(m.State.Wires) == 4 {
		// If there is more than one red wire and the last digit of the serial number is odd,
		// cut the last red wire.
		redIdxs := colorIndecies(m.State.Wires, valueobject.Red)
		if len(redIdxs) > 1 && helpers.SerialNumbersEndsWithOddDigit(m.bomb.SerialNumber) {
			if wireIndex == redIdxs[len(redIdxs)-1] {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If the last wire is yellow and there are no red wires, cut the first wire.
		if m.State.Wires[len(m.State.Wires)-1].WireColor == valueobject.Yellow && !hasColor(m.State.Wires, valueobject.Red) {
			if wireIndex == 0 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is exactly one blue wire, cut the first wire.
		if len(colorIndecies(m.State.Wires, valueobject.Blue)) == 1 {
			if wireIndex == 0 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is more than one yellow wire, cut the last wire.
		yellowIdxs := colorIndecies(m.State.Wires, valueobject.Yellow)
		if len(yellowIdxs) > 1 {
			if wireIndex == yellowIdxs[len(yellowIdxs)-1] {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// Otherwise, cut the second wire.
		if wireIndex == 1 {
			return m.cutSucceed()
		}
	} else if len(m.State.Wires) == 5 {
		// If the last wire is black and the last digit of the serial number is odd,
		// cut the fourth wire.
		if m.State.Wires[len(m.State.Wires)-1].WireColor == valueobject.Black && helpers.SerialNumbersEndsWithOddDigit(m.bomb.SerialNumber) {
			if wireIndex == 3 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is exactly one red wire and there are no yellow wires, cut the first wire.
		if len(colorIndecies(m.State.Wires, valueobject.Red)) == 1 && !hasColor(m.State.Wires, valueobject.Yellow) {
			if wireIndex == 0 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there are no black wires, cut the second wire.
		if !hasColor(m.State.Wires, valueobject.Black) {
			if wireIndex == 1 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// Otherwise, cut the first wire
		if wireIndex == 0 {
			return m.cutSucceed()
		}
	} else if len(m.State.Wires) == 6 {
		// If there are no yellow wires and the last digit of the serial number is odd,
		// cut the third wire.
		if !hasColor(m.State.Wires, valueobject.Yellow) && helpers.SerialNumbersEndsWithOddDigit(m.bomb.SerialNumber) {
			if wireIndex == 2 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is exactly one yellow wire and there is more than one white wire,
		// cut the fourth wire.
		if len(colorIndecies(m.State.Wires, valueobject.Yellow)) == 1 && len(colorIndecies(m.State.Wires, valueobject.White)) > 1 {
			if wireIndex == 3 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there are no red wires, cut the last wire.
		if !hasColor(m.State.Wires, valueobject.Red) {
			if wireIndex == len(m.State.Wires)-1 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// Otherwise, cut the fourth wire.
		if wireIndex == 3 {
			return m.cutSucceed()
		}
	}

	return true, nil
}

func hasColor(wires []valueobject.SimpleWire, color valueobject.Color) bool {
	for _, wire := range wires {
		if wire.WireColor == color {
			return true
		}
	}
	return false
}

func colorIndecies(wires []valueobject.SimpleWire, color valueobject.Color) []int {
	indecies := []int{}
	for i, wire := range wires {
		if wire.WireColor == color {
			indecies = append(indecies, i)
		}
	}
	return indecies
}
