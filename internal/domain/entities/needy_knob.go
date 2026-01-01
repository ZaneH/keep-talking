package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type NeedyKnobState struct {
	BaseModuleState
	// Light rows displayed to the player
	DisplayedPattern [][]bool
	// Current direction of the dial
	DialDirection valueobject.CardinalDirection
	// Timestamp of when the current countdown started
	CountdownStartedAt int64
	// Duration in seconds of the countdown
	CountdownDuration int16
}

func NewNeedyKnobState(rng ports.RandomGenerator) NeedyKnobState {
	displayPatternIdx := rng.GetIntInRange(0, len(knobLightStates)-1)

	return NeedyKnobState{
		BaseModuleState: BaseModuleState{},
		DisplayedPattern: [][]bool{
			knobLightStates[displayPatternIdx].firstRow,
			knobLightStates[displayPatternIdx].secondRow,
		},
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  int16(30),
	}
}

type NeedyKnobModule struct {
	BaseModule
	State NeedyKnobState
	rng   ports.RandomGenerator
}

func NewNeedyKnobModule(rng ports.RandomGenerator) *NeedyKnobModule {
	return &NeedyKnobModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		State: NewNeedyKnobState(rng),
		rng:   rng,
	}
}

func (m *NeedyKnobModule) String() string {
	var result strings.Builder
	result.WriteString("\n")

	for _, row := range m.State.DisplayedPattern {
		for _, light := range row {
			if light {
				result.WriteString("X ")
			} else {
				result.WriteString("  ")
			}
		}
		result.WriteString("\n")
	}

	result.WriteString(fmt.Sprintf("Dial Direction: %v\n", m.State.DialDirection))

	return result.String()
}

func (m *NeedyKnobModule) GetType() valueobject.ModuleType {
	return valueobject.NeedyKnobModule
}

func (m *NeedyKnobModule) SetState(state NeedyKnobState) {
	m.State = state
}

func (m *NeedyKnobModule) GetModuleState() ModuleState {
	return &m.State
}

func (m *NeedyKnobModule) RotateDial() (err error) {
	switch m.State.DialDirection {
	case valueobject.North:
		m.State.DialDirection = valueobject.East
	case valueobject.East:
		m.State.DialDirection = valueobject.South
	case valueobject.South:
		m.State.DialDirection = valueobject.West
	case valueobject.West:
		m.State.DialDirection = valueobject.North
	default:
		return fmt.Errorf("unknown direction %v", m.State.DialDirection)
	}

	return nil
}

type KnobLightState struct {
	solution  valueobject.CardinalDirection
	firstRow  []bool
	secondRow []bool
}

var knobLightStates = []KnobLightState{
	// Up
	{
		solution: valueobject.North,
		firstRow: []bool{
			false, false, true, false, true, true,
		},
		secondRow: []bool{
			true, true, true, true, false, true,
		},
	},
	{
		solution: valueobject.North,
		firstRow: []bool{
			true, false, true, false, true, false,
		},
		secondRow: []bool{
			false, true, true, false, true, true,
		},
	},

	// Down
	{
		solution: valueobject.South,
		firstRow: []bool{
			false, true, true, false, false, true,
		},
		secondRow: []bool{
			true, true, true, true, false, true,
		},
	},
	{
		solution: valueobject.South,
		firstRow: []bool{
			true, false, true, false, true, false,
		},
		secondRow: []bool{
			false, true, false, false, false, true,
		},
	},

	// Left
	{
		solution: valueobject.West,
		firstRow: []bool{
			false, false, false, false, true, false,
		},
		secondRow: []bool{
			true, false, false, true, true, true,
		},
	},
	{
		solution: valueobject.West,
		firstRow: []bool{
			false, false, false, false, true, false,
		},
		secondRow: []bool{
			false, false, false, true, true, false,
		},
	},

	// Right
	{
		solution: valueobject.East,
		firstRow: []bool{
			true, false, true, true, true, true,
		},
		secondRow: []bool{
			true, true, true, false, true, false,
		},
	},
	{
		solution: valueobject.East,
		firstRow: []bool{
			true, false, true, true, false, false,
		},
		secondRow: []bool{
			true, true, true, false, true, false,
		},
	},
}
