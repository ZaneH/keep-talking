package entities

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"

	"github.com/ZaneH/keep-talking/internal/application/helpers"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

const minStages = 3
const maxStages = 5

type SimonSaysState struct {
	BaseModuleState
	// Stores each color to be displayed in order.
	DisplaySequence []valueobject.Color
	// When the user starts their input, the module needs to know what index of the sequence they
	// are at. This is used to determine if the user input is correct or not. 0 indicates that the
	// module is waiting for the user to input the first color in the sequence.
	inputCheckIdx int
}

type SimonSaysModule struct {
	BaseModule
	state SimonSaysState
}

func NewSimonSaysModule() *SimonSaysModule {
	return &SimonSaysModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		state: SimonSaysState{
			DisplaySequence: generateDisplaySequence(),
		},
	}
}

func (m *SimonSaysModule) String() string {
	var result strings.Builder
	result.WriteString("\n")
	result.WriteString(fmt.Sprintf("Strikes: %d\n", m.bomb.StrikeCount))
	result.WriteString(fmt.Sprintf("Serial Number: %s\n", m.bomb.SerialNumber))
	result.WriteString("Current sequence: ")
	for i, color := range m.state.DisplaySequence {
		if i == m.state.inputCheckIdx {
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

func (m *SimonSaysModule) PressColor(c valueobject.Color) (nextSeq []valueobject.Color, strike bool, err error) {
	nextSeqIdx := min(m.state.inputCheckIdx+2, len(m.state.DisplaySequence))
	nextSeq = m.state.DisplaySequence[:nextSeqIdx]
	if m.state.inputCheckIdx >= len(m.state.DisplaySequence) {
		return nextSeq, false, fmt.Errorf("input check index out of bounds")
	}

	if !slices.Contains(simonSaysColors[:], c) {
		return nextSeq, false, fmt.Errorf("invalid color: %s", c)
	}

	displayColor := m.state.DisplaySequence[m.state.inputCheckIdx]
	correctColor, err := m.translateColor(displayColor)
	if err != nil {
		return nextSeq, false, fmt.Errorf("error translating color: %w", err)
	}

	if c == correctColor {
		if m.state.inputCheckIdx == len(m.state.DisplaySequence)-1 {
			m.state.MarkSolved = true
		} else {
			m.state.inputCheckIdx++
		}

		return nextSeq, false, nil
	}

	return nextSeq, true, nil
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

func generateDisplaySequence() []valueobject.Color {
	sequence := make([]valueobject.Color, maxStages)
	for i := minStages; i < maxStages; i++ {
		sequence = append(sequence, simonSaysColors[rand.Intn(len(simonSaysColors))])
	}
	return sequence
}
