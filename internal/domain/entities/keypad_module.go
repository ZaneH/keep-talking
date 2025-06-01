package entities

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

const nStages = 4

type KeypadState struct {
	BaseModuleState
	// Symbols currently displayed on the keypad.
	DisplayedSymbols []valueobject.Symbol
	// Stores the activated symbols. Mapped to true when activated.
	ActivatedSymbols map[valueobject.Symbol]bool
	// Symbols that the user has to activate in order.
	solution []valueobject.Symbol
}

type KeypadModule struct {
	BaseModule
	State KeypadState
	rng   ports.RandomGenerator
}

func NewKeypadModule(rng ports.RandomGenerator) *KeypadModule {
	displayedSymbols := generateDisplayedSymbols(rng)
	solution := generateKeypadSolution(rng, displayedSymbols)

	return &KeypadModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		State: KeypadState{
			DisplayedSymbols: displayedSymbols,
			ActivatedSymbols: make(map[valueobject.Symbol]bool, nStages),
			solution:         solution,
		},
		rng: rng,
	}
}

func (m *KeypadModule) String() string {
	var result strings.Builder
	result.WriteString("\n")

	for i, sym := range m.State.DisplayedSymbols {
		if i > 0 {
			result.WriteString(", ")
		}
		if m.State.ActivatedSymbols[sym] {
			result.WriteString(fmt.Sprintf("[%s]", sym))
		} else {
			result.WriteString(fmt.Sprintf("%s", sym))
		}
	}

	return result.String()
}

func (m *KeypadModule) GetType() valueobject.ModuleType {
	return valueobject.Keypad
}

func (m *KeypadModule) SetState(state KeypadState) {
	m.State = state
}

func (m *KeypadModule) GetModuleState() ModuleState {
	return &m.State
}

func (m *KeypadModule) PressSymbol(sym valueobject.Symbol) (map[valueobject.Symbol]bool, bool, error) {
	if !slices.Contains(keypadSymbols[:], sym) {
		return m.State.ActivatedSymbols, false, fmt.Errorf("invalid symbol: %s", sym)
	}

	active := m.State.ActivatedSymbols[sym]
	if active {
		return m.State.ActivatedSymbols, false, fmt.Errorf("symbol %s is already activated", sym)
	}

	m.State.ActivatedSymbols[sym] = true
	idx := max(0, len(m.State.ActivatedSymbols)-1)

	// They hit the current symbol in the sequence
	if sym == m.State.solution[idx] {
		// They hit the last symbol in the sequence
		if len(m.State.ActivatedSymbols) == nStages {
			m.State.MarkAsSolved()
			return m.State.ActivatedSymbols, false, nil
		}

		// They hit a symbol in the sequence, but not the last one
		m.State.ActivatedSymbols[sym] = true
		return m.State.ActivatedSymbols, false, nil
	} else {
		// Incorrect input, reset state
		m.State.ActivatedSymbols = make(map[valueobject.Symbol]bool, nStages)
	}

	return m.State.ActivatedSymbols, true, nil
}

func generateKeypadSolution(rng ports.RandomGenerator, available []valueobject.Symbol) []valueobject.Symbol {
	solution := make([]valueobject.Symbol, 0, len(available))
	for len(solution) < nStages {
		symbol := available[rng.GetIntInRange(0, len(available)-1)]
		if !slices.Contains(solution, symbol) {
			solution = append(solution, symbol)
		}
	}
	return solution
}

var keypadSymbols = [...]valueobject.Symbol{
	valueobject.Copyright,
	valueobject.FilledStar,
	valueobject.HollowStar,
	valueobject.SmileyFace,
	valueobject.DoubleK,
	valueobject.Omega,
	valueobject.SquidKnife,
	valueobject.Pumpkin,
	valueobject.HookN,
	// valueobject.Teepee,
	valueobject.Six,
	valueobject.SquigglyN,
	valueobject.At,
	valueobject.Ae,
	valueobject.MeltedThree,
	valueobject.Euro,
	// valueobject.Circle,
	valueobject.NWithHat,
	valueobject.Dragon,
	valueobject.QuestionMark,
	valueobject.Paragraph,
	valueobject.RightC,
	valueobject.LeftC,
	valueobject.Pitchfork,
	// valueobject.Tripod,
	valueobject.Cursive,
	valueobject.Tracks,
	valueobject.Balloon,
	// valueobject.WeirdNose,
	valueobject.Upsidedowny,
	valueobject.Bt,
}

func generateDisplayedSymbols(rng ports.RandomGenerator) []valueobject.Symbol {
	displayed := make([]valueobject.Symbol, 0, nStages)
	for len(displayed) < nStages {
		symbol := keypadSymbols[rng.GetIntInRange(0, len(keypadSymbols)-1)]
		if !slices.Contains(displayed, symbol) {
			displayed = append(displayed, symbol)
		}
	}
	return displayed
}
