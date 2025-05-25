package entities

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"

	"github.com/ZaneH/keep-talking/internal/application/helpers"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

const minWires = 3
const maxWires = 6

type SimpleWiresState struct {
	BaseModuleState
	Wires []valueobject.SimpleWire
}

func NewSimpleWiresState() SimpleWiresState {
	return SimpleWiresState{
		BaseModuleState: BaseModuleState{},
		Wires:           generateRandomWires(),
	}
}

type SimpleWiresModule struct {
	BaseModule
	State SimpleWiresState
}

func NewSimpleWiresModule() *SimpleWiresModule {
	return &SimpleWiresModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		State: NewSimpleWiresState(),
	}
}

func (m *SimpleWiresModule) GetModuleState() ModuleState {
	return &m.State
}

func (m *SimpleWiresModule) GetType() valueobject.ModuleType {
	return valueobject.SimpleWires
}

func (m *SimpleWiresModule) SetState(state SimpleWiresState) {
	m.State = state
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

	maxPossibleGaps := maxWires - nWires
	extraPositions := 0
	if maxPossibleGaps > 0 {
		extraPositions = rand.Intn(maxPossibleGaps + 1)
	}

	totalPositions := nWires + extraPositions
	possibleIndices := make([]int, totalPositions)
	for i := 0; i < totalPositions; i++ {
		possibleIndices[i] = i
	}

	rand.Shuffle(len(possibleIndices), func(i, j int) {
		possibleIndices[i], possibleIndices[j] = possibleIndices[j], possibleIndices[i]
	})

	selectedIndices := possibleIndices[:nWires]

	for i := 0; i < nWires; i++ {
		colorIndex := rand.Intn(len(valueobject.SimpleWireColors))
		color := valueobject.SimpleWireColors[colorIndex]

		wires[i] = valueobject.SimpleWire{
			WireColor: color,
			Position:  selectedIndices[i],
		}
	}

	return wires
}

func (m *SimpleWiresModule) cutSucceed() (strike bool, err error) {
	m.State.MarkAsSolved()
	return false, nil
}

func (m *SimpleWiresModule) CutWire(wirePos int) (strike bool, err error) {
	sorted := make([]valueobject.SimpleWire, len(m.State.Wires))
	copy(sorted, m.State.Wires)

	// Create a sorted copy of the wires based on their positions
	slices.SortFunc(sorted, func(a, b valueobject.SimpleWire) int {
		return a.Position - b.Position
	})

	// Get the index, accounting for any gaps created by extra positions
	wireIdx := slices.IndexFunc(sorted, func(w valueobject.SimpleWire) bool {
		return w.Position == wirePos
	})

	wire := sorted[wireIdx]

	if wire.IsCut {
		return false, errors.New("wire already cut")
	}

	wire.IsCut = true

	if len(sorted) == 3 {
		// If there are no red wires, cut the second wire.
		if !hasColor(sorted, valueobject.Red) {
			if wireIdx == 1 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If the last wire is white, cut the last wire.
		if sorted[len(sorted)-1].WireColor == valueobject.White {
			if wireIdx == len(sorted)-1 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is more than one blue wire, cut the last blue wire.
		blueIdxs := colorIndecies(sorted, valueobject.Blue)
		if len(blueIdxs) > 1 {
			if wireIdx == blueIdxs[len(blueIdxs)-1] {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// Otherwise cut the last wire.
		if wireIdx == len(sorted)-1 {
			return m.cutSucceed()
		} else {
			return true, nil
		}
	} else if len(sorted) == 4 {
		// If there is more than one red wire and the last digit of the serial number is odd,
		// cut the last red wire.
		redIdxs := colorIndecies(sorted, valueobject.Red)
		if len(redIdxs) > 1 && helpers.SerialNumbersEndsWithOddDigit(m.bomb.SerialNumber) {
			if wireIdx == redIdxs[len(redIdxs)-1] {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If the last wire is yellow and there are no red wires, cut the first wire.
		if sorted[len(sorted)-1].WireColor == valueobject.Yellow && !hasColor(sorted, valueobject.Red) {
			if wireIdx == 0 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is exactly one blue wire, cut the first wire.
		if len(colorIndecies(sorted, valueobject.Blue)) == 1 {
			if wireIdx == 0 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is more than one yellow wire, cut the last wire.
		yellowIdxs := colorIndecies(sorted, valueobject.Yellow)
		if len(yellowIdxs) > 1 {
			if wireIdx == yellowIdxs[len(yellowIdxs)-1] {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// Otherwise, cut the second wire.
		if wireIdx == 1 {
			return m.cutSucceed()
		}
	} else if len(sorted) == 5 {
		// If the last wire is black and the last digit of the serial number is odd,
		// cut the fourth wire.
		if sorted[len(sorted)-1].WireColor == valueobject.Black && helpers.SerialNumbersEndsWithOddDigit(m.bomb.SerialNumber) {
			if wireIdx == 3 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is exactly one red wire and there are no yellow wires, cut the first wire.
		if len(colorIndecies(sorted, valueobject.Red)) == 1 && !hasColor(sorted, valueobject.Yellow) {
			if wireIdx == 0 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there are no black wires, cut the second wire.
		if !hasColor(sorted, valueobject.Black) {
			if wireIdx == 1 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// Otherwise, cut the first wire
		if wireIdx == 0 {
			return m.cutSucceed()
		}
	} else if len(sorted) == 6 {
		// If there are no yellow wires and the last digit of the serial number is odd,
		// cut the third wire.
		if !hasColor(sorted, valueobject.Yellow) && helpers.SerialNumbersEndsWithOddDigit(m.bomb.SerialNumber) {
			if wireIdx == 2 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is exactly one yellow wire and there is more than one white wire,
		// cut the fourth wire.
		if len(colorIndecies(sorted, valueobject.Yellow)) == 1 && len(colorIndecies(sorted, valueobject.White)) > 1 {
			if wireIdx == 3 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there are no red wires, cut the last wire.
		if !hasColor(sorted, valueobject.Red) {
			if wireIdx == len(sorted)-1 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// Otherwise, cut the fourth wire.
		if wireIdx == 3 {
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
