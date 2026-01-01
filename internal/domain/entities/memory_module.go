package entities

import (
	"fmt"
	"log"

	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

// The memoryRound struct stores the label number and the button position
// to reference for future button presses.
type memoryRound struct {
	buttonNumber   int
	buttonPosition int
}

type MemoryState struct {
	BaseModuleState
	// Number on screen for the defuser to read aloud.
	ScreenNumber int
	// 4 numbers displayed to the player in display order.
	DisplayedNumbers []int
	// Past rounds tracked for the solution.
	pastRounds []memoryRound
	// Tracks the current stage, [1, 5].
	Stage int
}

func NewMemoryState(rng ports.RandomGenerator) MemoryState {
	return MemoryState{
		BaseModuleState:  BaseModuleState{},
		ScreenNumber:     rng.GetIntInRange(1, 4),
		DisplayedNumbers: generateMemoryDisplayedNumbers(rng),
		pastRounds:       make([]memoryRound, 0, 5),
		Stage:            1,
	}
}

type MemoryModule struct {
	BaseModule
	State MemoryState
	rng   ports.RandomGenerator
}

func NewMemoryModule(rng ports.RandomGenerator) *MemoryModule {
	return &MemoryModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		State: NewMemoryState(rng),
		rng:   rng,
	}
}

func (m *MemoryModule) GetModuleState() ModuleState {
	return &m.State
}

func (m *MemoryModule) GetType() valueobject.ModuleType {
	return valueobject.MemoryModule
}

func (m *MemoryModule) SetState(state MemoryState) {
	m.State = state
}

func (m *MemoryModule) String() string {
	var result = "\n"
	result += fmt.Sprintf("Stage: %d\n", m.State.Stage)
	result += fmt.Sprintf("Screen Number: %d\n", m.State.ScreenNumber)
	result += "Displayed Words: "
	for i, word := range m.State.DisplayedNumbers {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%d", word)
	}

	return result
}

func (m *MemoryModule) PressButton(btnIdx int) (strike bool, err error) {
	if btnIdx < 0 || btnIdx >= len(m.State.DisplayedNumbers) {
		return false, fmt.Errorf("invalid button index: %d", btnIdx)
	}

	n := m.State.DisplayedNumbers[btnIdx]
	p := btnIdx + 1

	switch m.State.Stage {
	case 1:
		// If the display is 1, press the button in the second position.
		if m.State.ScreenNumber == 1 {
			if p == 2 {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 2, press the button in the second position.
		if m.State.ScreenNumber == 2 {
			if p == 2 {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 3, press the button in the third position.
		if m.State.ScreenNumber == 3 {
			if p == 3 {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 4, press the button in the fourth position.
		if m.State.ScreenNumber == 4 {
			if p == 4 {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
	case 2:
		// If the display is 1, press the button labeled "4".
		if m.State.ScreenNumber == 1 {
			if n == 4 {
				m.appendPastRound(n, findButtonLabeled(4, m.State.DisplayedNumbers))
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 2, press the button in the same position as you pressed in stage 1.
		if m.State.ScreenNumber == 2 {
			if p == m.State.pastRounds[0].buttonPosition {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 3, press the button in the first position.
		if m.State.ScreenNumber == 3 {
			if p == 1 {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 4, press the button in the same position as you pressed in stage 1.
		if m.State.ScreenNumber == 4 {
			if p == m.State.pastRounds[0].buttonPosition {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
	case 3:
		// If the display is 1, press the button with the same label you pressed in stage 2.
		if m.State.ScreenNumber == 1 {
			if n == m.State.pastRounds[1].buttonNumber {
				m.appendPastRound(n, findButtonLabeled(m.State.pastRounds[1].buttonNumber, m.State.DisplayedNumbers))
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 2, press the button with the same label you pressed in stage 1.
		if m.State.ScreenNumber == 2 {
			if n == m.State.pastRounds[0].buttonNumber {
				m.appendPastRound(n, findButtonLabeled(m.State.pastRounds[0].buttonNumber, m.State.DisplayedNumbers))
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 3, press the button in the third position.
		if m.State.ScreenNumber == 3 {
			if p == 3 {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 4, press the button labeled "4".
		if m.State.ScreenNumber == 4 {
			if n == 4 {
				m.appendPastRound(n, findButtonLabeled(4, m.State.DisplayedNumbers))
				m.incrementStage()
				return false, nil
			}
		}
	case 4:
		// If the display is 1, press the button in the same position as you pressed in stage 1.
		if m.State.ScreenNumber == 1 {
			if p == m.State.pastRounds[0].buttonPosition {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 2, press the button in the first position.
		if m.State.ScreenNumber == 2 {
			if p == 1 {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 3, press the button in the same position as you pressed in stage 2.
		if m.State.ScreenNumber == 3 {
			if p == m.State.pastRounds[1].buttonPosition {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 4, press the button in the same position as you pressed in stage 2.
		if m.State.ScreenNumber == 4 {
			if p == m.State.pastRounds[1].buttonPosition {
				m.appendPastRound(n, p)
				m.incrementStage()
				return false, nil
			}
		}
	case 5:
		// If the display is 1, press the button with the same label you pressed in stage 1.
		if m.State.ScreenNumber == 1 {
			if n == m.State.pastRounds[0].buttonNumber {
				m.appendPastRound(n, findButtonLabeled(m.State.pastRounds[0].buttonNumber, m.State.DisplayedNumbers))
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 2, press the button with the same label you pressed in stage 2.
		if m.State.ScreenNumber == 2 {
			if n == m.State.pastRounds[1].buttonNumber {
				m.appendPastRound(n, findButtonLabeled(m.State.pastRounds[1].buttonNumber, m.State.DisplayedNumbers))
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 3, press the button with the same label you pressed in stage 4.
		if m.State.ScreenNumber == 3 {
			if n == m.State.pastRounds[3].buttonNumber {
				m.appendPastRound(n, findButtonLabeled(m.State.pastRounds[3].buttonNumber, m.State.DisplayedNumbers))
				m.incrementStage()
				return false, nil
			}
		}
		// If the display is 4, press the button with the same label you pressed in stage 3.
		if m.State.ScreenNumber == 4 {
			if n == m.State.pastRounds[2].buttonNumber {
				m.appendPastRound(n, findButtonLabeled(m.State.pastRounds[2].buttonNumber, m.State.DisplayedNumbers))
				m.incrementStage()
				return false, nil
			}
		}
	default:
		return false, fmt.Errorf("invalid stage: %d", m.State.Stage)
	}

	// Incorrect button press, strike the bomb and reset the stage.
	m.State.pastRounds = make([]memoryRound, 0, 5)
	m.State.Stage = 1
	m.State.ScreenNumber = m.rng.GetIntInRange(1, 4)
	m.State.DisplayedNumbers = generateMemoryDisplayedNumbers(m.rng)
	return true, nil
}

func findButtonLabeled(number int, order []int) int {
	for i, n := range order {
		if n == number {
			return i + 1
		}
	}

	log.Fatalf("Button labeled %d not found in order %v", number, order)
	return -1
}

func generateMemoryDisplayedNumbers(rng ports.RandomGenerator) []int {
	order := []int{1, 2, 3, 4}
	rng.Shuffle(len(order), func(i int, j int) {
		order[i], order[j] = order[j], order[i]
	})

	return order
}

func (m *MemoryModule) appendPastRound(buttonNumber, buttonPosition int) {
	m.State.pastRounds = append(m.State.pastRounds, memoryRound{
		buttonNumber:   buttonNumber,
		buttonPosition: buttonPosition,
	})
}

func (m *MemoryModule) incrementStage() {
	m.State.Stage++
	m.State.ScreenNumber = m.rng.GetIntInRange(1, 4)
	m.State.DisplayedNumbers = generateMemoryDisplayedNumbers(m.rng)
	if m.State.Stage > 5 {
		m.State.MarkAsSolved()
	}
}
