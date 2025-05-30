package entities

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ZaneH/keep-talking/internal/application/helpers"
	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

const minStages = 3
const maxStages = 5

type SimonSaysState struct {
	BaseModuleState
	// Stores the current sequence to be displayed to the user. Grows as the user progresses.
	DisplaySequence []valueobject.Color
	// When the user starts their input, the module needs to know what index of the sequence they
	// are at. This is used to determine if the user input is correct or not. 0 indicates that the
	// module is waiting for the user to input the first color in the sequence.
	InputCheckIdx int
	// The number of sequences that the user will need to complete to solve the module.
	nStages int
}

type SimonSaysModule struct {
	BaseModule
	state SimonSaysState
	rng   ports.RandomGenerator
}

func NewSimonSaysModule(rng ports.RandomGenerator, nStages *int) *SimonSaysModule {
	n := 0
	if nStages == nil {
		n = minStages + rng.GetIntInRange(minStages, maxStages)
	}

	return &SimonSaysModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		state: SimonSaysState{
			DisplaySequence: []valueobject.Color{generateRandomSimonSaysColor(rng)},
			InputCheckIdx:   0,
			nStages:         n,
		},
		rng: rng,
	}
}

func (m *SimonSaysModule) String() string {
	var result strings.Builder
	result.WriteString("\n")
	result.WriteString(fmt.Sprintf("Strikes: %d\n", m.bomb.StrikeCount))
	result.WriteString(fmt.Sprintf("Serial Number: %s\n", m.bomb.SerialNumber))
	result.WriteString("Current sequence: ")
	for i, color := range m.state.DisplaySequence {
		if i == m.state.InputCheckIdx {
			result.WriteString(fmt.Sprintf("[%s] ", color))
		} else {
			result.WriteString(string(color) + " ")
		}
	}
	return result.String()
}

func (m *SimonSaysModule) GetType() valueobject.ModuleType {
	return valueobject.SimonSays
}

func (m *SimonSaysModule) SetState(state SimonSaysState) {
	m.state = state
}

func (m *SimonSaysModule) GetModuleState() ModuleState {
	return &m.state
}

func (m *SimonSaysModule) PressColor(c valueobject.Color) (finishedSeq bool, nextSeq []valueobject.Color, strike bool, err error) {
	if m.state.InputCheckIdx >= len(m.state.DisplaySequence) {
		return false, nextSeq, false, fmt.Errorf("input check index out of bounds")
	}

	if !slices.Contains(simonSaysColors[:], c) {
		return false, nextSeq, false, fmt.Errorf("invalid color: %s", c)
	}

	displayColor := m.state.DisplaySequence[m.state.InputCheckIdx]
	correctColor, err := m.translateColor(displayColor)
	if err != nil {
		return false, nextSeq, false, fmt.Errorf("error translating color: %w", err)
	}

	if c == correctColor {
		// They hit the last color in the current sequence
		// solved if they have completed all stages
		// otherwise, add a new color to the sequence
		if m.state.InputCheckIdx == len(m.state.DisplaySequence)-1 {
			if m.state.InputCheckIdx+1 == m.state.nStages {
				m.state.MarkAsSolved()
				return true, nextSeq, false, nil
			}

			newColor := generateRandomSimonSaysColor(m.rng)
			m.state.InputCheckIdx = 0
			m.state.DisplaySequence = append(m.state.DisplaySequence, newColor)
			return true, m.state.DisplaySequence, false, nil
		} else {
			// They hit a color in the sequence, but not the last one
			// Just increment the input check index
			m.state.InputCheckIdx++
			return false, nextSeq, false, nil
		}
	} else {
		// Incorrect input, reset state
		m.state.InputCheckIdx = 0
		m.state.DisplaySequence = []valueobject.Color{generateRandomSimonSaysColor(m.rng)}
		return false, m.state.DisplaySequence, true, nil
	}
}

func (m *SimonSaysModule) translateColor(c valueobject.Color) (translated valueobject.Color, err error) {
	if helpers.SerialNumberContainsVowel(m.bomb.SerialNumber) {
		if m.bomb.StrikeCount == 0 {
			switch c {
			case valueobject.Red:
				translated = valueobject.Blue
			case valueobject.Green:
				translated = valueobject.Yellow
			case valueobject.Blue:
				translated = valueobject.Red
			case valueobject.Yellow:
				translated = valueobject.Green
			}
		} else if m.bomb.StrikeCount == 1 {
			switch c {
			case valueobject.Red:
				translated = valueobject.Yellow
			case valueobject.Green:
				translated = valueobject.Blue
			case valueobject.Blue:
				translated = valueobject.Green
			case valueobject.Yellow:
				translated = valueobject.Red
			}
		} else if m.bomb.StrikeCount >= 2 {
			switch c {
			case valueobject.Red:
				translated = valueobject.Green
			case valueobject.Green:
				translated = valueobject.Yellow
			case valueobject.Blue:
				translated = valueobject.Red
			case valueobject.Yellow:
				translated = valueobject.Blue
			}
		}
	} else {
		if m.bomb.StrikeCount == 0 {
			switch c {
			case valueobject.Red:
				translated = valueobject.Blue
			case valueobject.Green:
				translated = valueobject.Green
			case valueobject.Blue:
				translated = valueobject.Yellow
			case valueobject.Yellow:
				translated = valueobject.Red
			}
		} else if m.bomb.StrikeCount == 1 {
			switch c {
			case valueobject.Red:
				translated = valueobject.Red
			case valueobject.Green:
				translated = valueobject.Yellow
			case valueobject.Blue:
				translated = valueobject.Blue
			case valueobject.Yellow:
				translated = valueobject.Green
			}
		} else if m.bomb.StrikeCount >= 2 {
			switch c {
			case valueobject.Red:
				translated = valueobject.Yellow
			case valueobject.Green:
				translated = valueobject.Blue
			case valueobject.Blue:
				translated = valueobject.Green
			case valueobject.Yellow:
				translated = valueobject.Red
			}
		}
	}

	return translated, nil
}

var simonSaysColors = [...]valueobject.Color{
	valueobject.Red,
	valueobject.Green,
	valueobject.Blue,
	valueobject.Yellow,
}

func generateRandomSimonSaysColor(rng ports.RandomGenerator) valueobject.Color {
	return simonSaysColors[rng.GetIntInRange(0, len(simonSaysColors)-1)]
}
