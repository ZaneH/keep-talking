package entities

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/ZaneH/defuse.party-go/internal/application/helpers"
	"github.com/ZaneH/defuse.party-go/internal/domain/ports"
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
	"github.com/google/uuid"
)

const minWires = 3
const maxWires = 6

type WiresState struct {
	BaseModuleState
	Wires []valueobject.Wire
}

func NewWiresState(rng ports.RandomGenerator) WiresState {
	return WiresState{
		BaseModuleState: BaseModuleState{},
		Wires:           generateRandomWires(rng),
	}
}

type WiresModule struct {
	BaseModule
	State WiresState
	rng   ports.RandomGenerator
}

func NewWiresModule(rng ports.RandomGenerator) *WiresModule {
	return &WiresModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		State: NewWiresState(rng),
		rng:   rng,
	}
}

func (m *WiresModule) GetModuleState() ModuleState {
	return &m.State
}

func (m *WiresModule) GetType() valueobject.ModuleType {
	return valueobject.WiresModule
}

func (m *WiresModule) SetState(state WiresState) {
	m.State = state
}

func (m *WiresModule) String() string {
	var result strings.Builder
	result.WriteString("\n")
	for i, wire := range m.State.Wires {
		if wire.IsCut {
			fmt.Fprintf(&result, "Wire %d: %s (cut)\n", i, wire.WireColor)
		} else {
			fmt.Fprintf(&result, "Wire %d: %s\n", i, wire.WireColor)
		}
	}

	return result.String()
}

func generateRandomWires(rng ports.RandomGenerator) []valueobject.Wire {
	nWires := rng.GetIntInRange(minWires, maxWires)
	wires := make([]valueobject.Wire, nWires)

	maxPossibleGaps := maxWires - nWires
	extraPositions := 0
	if maxPossibleGaps > 0 {
		extraPositions = rng.GetIntInRange(0, maxPossibleGaps)
	}

	totalPositions := nWires + extraPositions
	possibleIndices := make([]int, totalPositions)
	for i := range totalPositions {
		possibleIndices[i] = i
	}

	rng.Shuffle(len(possibleIndices), func(i, j int) {
		possibleIndices[i], possibleIndices[j] = possibleIndices[j], possibleIndices[i]
	})

	selectedIndices := possibleIndices[:nWires]

	for i := range nWires {
		colorIndex := rng.GetIntInRange(0, len(wireColors)-1)
		color := wireColors[colorIndex]

		wires[i] = valueobject.Wire{
			WireColor: color,
			Position:  selectedIndices[i],
		}
	}

	return wires
}

func (m *WiresModule) cutSucceed() (strike bool, err error) {
	m.State.MarkAsSolved()
	return false, nil
}

func (m *WiresModule) CutWire(wirePos int) (strike bool, err error) {
	sorted := make([]valueobject.Wire, len(m.State.Wires))
	copy(sorted, m.State.Wires)

	// Create a sorted copy of the wires based on their positions
	slices.SortFunc(sorted, func(a, b valueobject.Wire) int {
		return a.Position - b.Position
	})

	// Get the index, accounting for any gaps created by extra positions
	wireIdx := slices.IndexFunc(sorted, func(w valueobject.Wire) bool {
		return w.Position == wirePos
	})

	// If the wire position is invalid
	if wireIdx == -1 {
		return false, errors.New("invalid wire position")
	}

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
		blueIdxs := colorIndices(sorted, valueobject.Blue)
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
		redIdxs := colorIndices(sorted, valueobject.Red)
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
		if len(colorIndices(sorted, valueobject.Blue)) == 1 {
			if wireIdx == 0 {
				return m.cutSucceed()
			} else {
				return true, nil
			}
		}

		// If there is more than one yellow wire, cut the last wire.
		yellowIdxs := colorIndices(sorted, valueobject.Yellow)
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
		if len(colorIndices(sorted, valueobject.Red)) == 1 && !hasColor(sorted, valueobject.Yellow) {
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
		if len(colorIndices(sorted, valueobject.Yellow)) == 1 && len(colorIndices(sorted, valueobject.White)) > 1 {
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

func hasColor(wires []valueobject.Wire, color valueobject.Color) bool {
	for _, wire := range wires {
		if wire.WireColor == color {
			return true
		}
	}
	return false
}

func colorIndices(wires []valueobject.Wire, color valueobject.Color) []int {
	indecies := []int{}
	for i, wire := range wires {
		if wire.WireColor == color {
			indecies = append(indecies, i)
		}
	}
	return indecies
}

var wireColors = [...]valueobject.Color{
	valueobject.Yellow,
	valueobject.Red,
	valueobject.Blue,
	valueobject.Black,
	valueobject.White,
}
