package entities

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type BigButtonState struct {
	ModuleState
	buttonColor valueobject.Color
	// The label on the button
	label string
	// The final digit of the bomb's timer must match to solve. If it's nil,
	// tapping the button is the only valid solution.
	releaseDigit *int
}

func NewButtonState() BigButtonState {
	colorIdx := rand.Intn(len(bigButtonColors))
	labelIdx := rand.Intn(len(availableButtonWords))
	return BigButtonState{
		buttonColor:  bigButtonColors[colorIdx],
		label:        string(availableButtonWords[labelIdx]),
		releaseDigit: nil,
	}
}

type BigButtonModule struct {
	ModuleID uuid.UUID
	State    BigButtonState
	bomb     *Bomb
}

func NewBigButtonModule(bomb *Bomb) *BigButtonModule {
	return &BigButtonModule{
		ModuleID: uuid.New(),
		State:    NewButtonState(),
		bomb:     bomb,
	}
}

func (m *BigButtonModule) GetModuleID() uuid.UUID {
	return m.ModuleID
}

func (m *BigButtonModule) GetType() valueobject.ModuleType {
	return valueobject.SimpleWires
}

func (m *BigButtonModule) SetState(state BigButtonState) {
	m.State = state
}

func (m *BigButtonModule) String() string {
	var result string
	result += fmt.Sprintf("Button Color: %s\n", m.State.buttonColor)

	return result
}

func (m *BigButtonModule) IsSolved() bool {
	return false
}

func (m *BigButtonModule) PressButton(pressType valueobject.PressType) (*valueobject.Color, error) {
	handled := false
	var err error
	var color valueobject.Color

	switch pressType {
	case valueobject.PressTypeTap:
		handled, err = m.handleShortPress()
	case valueobject.PressTypeHold:
		handled, color = m.handleLongPress()
	case valueobject.PressTypeRelease:
		handled, err = m.handleLongPressRelease()
	default:
		return nil, errors.New("invalid press type")
	}

	if !handled {
		return nil, errors.New("press type not handled")
	}

	return &color, err
}

// Handles a short press (tap) of the button. There are certain conditions that must be met
// for the module to be marked as solved. If the conditions are not met, it will return an error.
func (m *BigButtonModule) handleShortPress() (handled bool, err error) {
	if m.bomb.Batteries > 1 && m.State.label == string(Detonate) {
		m.State.MarkSolved = true
		return true, nil
	}

	if m.bomb.Batteries > 2 && m.bomb.Indicators["FRK"].Lit {
		m.State.MarkSolved = true
		return true, nil
	}

	if m.State.buttonColor == valueobject.Red && m.State.label == string(Hold) {
		m.State.MarkSolved = true
		return true, nil
	}

	return false, nil
}

// Generates a random strip color regardless of the button color or label. There is no
// possibility of a strike from this action. If the button is to be pressed and held, this
// function will also set the release digit to either 4, 1, or 5.
func (m *BigButtonModule) handleLongPress() (handled bool, color valueobject.Color) {
	randIdx := rand.Intn(len(bigButtonStripColors))
	color = bigButtonStripColors[randIdx]

	return true, color
}

// Handles the release of a long press. If the release digit is nil, it will return an error.
// If the release digit matches the last digit of the bomb's timer, the module is marked as
// solved. Otherwise, it returns an error.
func (m *BigButtonModule) handleLongPressRelease() (handled bool, err error) {
	time := m.bomb.GetTimeLeft()
	if m.State.releaseDigit == nil {
		return true, errors.New("release digit is nil")
	} else {
		if int(time.Seconds())%10 == *m.State.releaseDigit {
			m.State.MarkSolved = true
		} else {
			return true, errors.New("release digit does not match")
		}
	}

	return false, nil
}

var bigButtonColors = [...]valueobject.Color{
	valueobject.Blue,
	valueobject.Red,
	valueobject.White,
	valueobject.Yellow,
	valueobject.Black,
}

var bigButtonStripColors = [...]valueobject.Color{
	valueobject.Blue,
	valueobject.Red,
	valueobject.White,
	valueobject.Yellow,
}

type bigButtonWords string

const (
	Abort    bigButtonWords = "Abort"
	Detonate bigButtonWords = "Detonate"
	Hold     bigButtonWords = "Hold"
	Press    bigButtonWords = "Press"
)

var availableButtonWords = [...]bigButtonWords{
	Abort,
	Detonate,
	Hold,
	Press,
}
