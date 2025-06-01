package entities

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

const nKeypadStages = 4

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
	displayed, col := generateDisplayedSymbols(rng)
	solution := generateKeypadSolution(displayed, col)

	return &KeypadModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		State: KeypadState{
			DisplayedSymbols: displayed,
			ActivatedSymbols: make(map[valueobject.Symbol]bool, nKeypadStages),
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
	active := m.State.ActivatedSymbols[sym]
	if active {
		return m.State.ActivatedSymbols, false, fmt.Errorf("symbol %s is already activated", sym)
	}

	m.State.ActivatedSymbols[sym] = true
	idx := max(0, len(m.State.ActivatedSymbols)-1)

	// They hit the current symbol in the sequence
	if sym == m.State.solution[idx] {
		// They hit the last symbol in the sequence
		if len(m.State.ActivatedSymbols) == nKeypadStages {
			m.State.MarkAsSolved()
			return m.State.ActivatedSymbols, false, nil
		}

		// They hit a symbol in the sequence, but not the last one
		m.State.ActivatedSymbols[sym] = true
		return m.State.ActivatedSymbols, false, nil
	} else {
		// Incorrect input, reset state
		m.State.ActivatedSymbols = make(map[valueobject.Symbol]bool, nKeypadStages)
	}

	return m.State.ActivatedSymbols, true, nil
}

var keypadCol1 = [...]valueobject.Symbol{
	valueobject.Balloon,
	valueobject.At,
	valueobject.UpsideDownY,
	valueobject.SquigglyN,
	valueobject.SquidKnife,
	valueobject.HookN,
	valueobject.LeftC,
}

var keypadCol2 = [...]valueobject.Symbol{
	valueobject.Euro,
	valueobject.Balloon,
	valueobject.LeftC,
	valueobject.Cursive,
	valueobject.HollowStar,
	valueobject.HookN,
	valueobject.QuestionMark,
}

var keypadCol3 = [...]valueobject.Symbol{
	valueobject.Copyright,
	valueobject.Pumpkin,
	valueobject.Cursive,
	valueobject.DoubleK,
	valueobject.MeltedThree,
	valueobject.UpsideDownY,
	valueobject.HollowStar,
}

var keypadCol4 = [...]valueobject.Symbol{
	valueobject.Six,
	valueobject.Paragraph,
	valueobject.Bt,
	valueobject.SquidKnife,
	valueobject.DoubleK,
	valueobject.QuestionMark,
	valueobject.SmileyFace,
}

var keypadCol5 = [...]valueobject.Symbol{
	valueobject.Pitchfork,
	valueobject.SmileyFace,
	valueobject.Bt,
	valueobject.RightC,
	valueobject.Paragraph,
	valueobject.Dragon,
	valueobject.FilledStar,
}

var keypadCol6 = [...]valueobject.Symbol{
	valueobject.Six,
	valueobject.Euro,
	valueobject.Tracks,
	valueobject.Ae,
	valueobject.Pitchfork,
	valueobject.NWithHat,
	valueobject.Omega,
}

var columns = [][]valueobject.Symbol{
	keypadCol1[:],
	keypadCol2[:],
	keypadCol3[:],
	keypadCol4[:],
	keypadCol5[:],
	keypadCol6[:],
}

func generateDisplayedSymbols(rng ports.RandomGenerator) (displayed []valueobject.Symbol, column []valueobject.Symbol) {
	displayed = make([]valueobject.Symbol, 0, nKeypadStages)
	randomColumn := rng.GetIntInRange(0, len(columns)-1)
	keypadSymbols := columns[randomColumn]

	for len(displayed) < nKeypadStages {
		symbol := keypadSymbols[rng.GetIntInRange(0, len(keypadSymbols)-1)]
		if !slices.Contains(displayed, symbol) {
			displayed = append(displayed, symbol)
		}
	}

	return displayed, keypadSymbols
}

func generateKeypadSolution(displayed []valueobject.Symbol, column []valueobject.Symbol) []valueobject.Symbol {
	solution := make([]valueobject.Symbol, 0, len(displayed))
	for c := range column {
		for _, sym := range displayed {
			if column[c] == sym {
				solution = append(solution, sym)
				break
			}
		}
	}

	return solution
}
