package entities

import (
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type BigButtonState struct {
	BaseModuleState
	ButtonColor valueobject.Color
	// The label on the button
	Label string
	// The final digit of the bomb's timer must match to solve. If it's nil, tapping the
	// button is the only valid solution.
	ReleaseDigit *int
}

func NewButtonState() BigButtonState {
	colorIdx := rand.Intn(len(bigButtonColors))
	labelIdx := rand.Intn(len(availableButtonWords))

	return BigButtonState{
		ButtonColor:     bigButtonColors[colorIdx],
		Label:           string(availableButtonWords[labelIdx]),
		ReleaseDigit:    nil,
		BaseModuleState: BaseModuleState{},
	}
}

type BigButtonModule struct {
	BaseModule
	State BigButtonState
}

func NewBigButtonModule(bomb *Bomb) *BigButtonModule {
	return &BigButtonModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
			bomb:     bomb,
		},
		State: NewButtonState(),
	}
}

func (m *BigButtonModule) GetType() valueobject.ModuleType {
	return valueobject.BigButton
}

func (m *BigButtonModule) GetModuleState() ModuleState {
	return &m.State
}

func (m *BigButtonModule) SetState(state BigButtonState) {
	m.State = state
}

func (m *BigButtonModule) String() string {
	var result string
	result += fmt.Sprintf("Button Color: %s\n", m.State.ButtonColor)

	return result
}

func (m *BigButtonModule) PressButton(pressType valueobject.PressType) (color *valueobject.Color, strike bool, err error) {
	handled := false

	switch pressType {
	case valueobject.PressTypeTap:
		handled, err = m.handleShortPress()
	case valueobject.PressTypeHold:
		handled, color = m.handleLongPress()
	case valueobject.PressTypeRelease:
		handled, err = m.handleLongPressRelease()
	default:
		return nil, true, errors.New("invalid press type")
	}

	// If not handled by a specific case, fallback to random color and release digit
	if !handled {
		if pressType == valueobject.PressTypeTap {
			return nil, true, errors.New("invalid short press")
		}

		if pressType == valueobject.PressTypeHold {
			m.State.ReleaseDigit = generateReleaseDigit()
			color, err = releaseDigitToStripColor(m.State.ReleaseDigit)
			if err != nil {
				log.Println("Error generating strip color for Big Button:", err)
				return nil, false, err
			}
		}
	}

	return color, false, err
}

// Handles a short press (tap) of the button. There are certain conditions that must be met
// for the module to be marked as solved. If the conditions are not met, it will return an error.
func (m *BigButtonModule) handleShortPress() (handled bool, err error) {
	if m.bomb.Batteries > 1 && m.State.Label == string(Detonate) {
		m.State.MarkAsSolved()
		return true, nil
	}

	if m.bomb.Batteries > 2 && m.bomb.Indicators["FRK"].Lit {
		m.State.MarkAsSolved()
		return true, nil
	}

	if m.State.ButtonColor == valueobject.Red && m.State.Label == string(Hold) {
		m.State.MarkAsSolved()
		return true, nil
	}

	return false, nil
}

// Generates a random strip color regardless of the button color or label. There is no
// possibility of a strike from this action. If the button is to be pressed and held, this
// function will also set the release digit to either 4, 1, or 5.
func (m *BigButtonModule) handleLongPress() (handled bool, color *valueobject.Color) {
	if m.State.ButtonColor == valueobject.Blue && m.State.Label == string(Abort) {
		m.State.ReleaseDigit = generateReleaseDigit()
		handled = true
	}

	if m.State.ButtonColor == valueobject.White && m.bomb.Indicators["CAR"].Lit {
		m.State.ReleaseDigit = generateReleaseDigit()
		handled = true
	}

	if m.State.ButtonColor == valueobject.Yellow {
		m.State.ReleaseDigit = generateReleaseDigit()
		handled = true
	}

	color, err := releaseDigitToStripColor(m.State.ReleaseDigit)
	if err != nil {
		log.Println("Error generating strip color for Big Button:", err)
		return false, nil
	}

	return handled, color
}

// Handles the release of a long press. If the release digit is nil, it will return an error.
// If the release digit matches the last digit of the bomb's timer, the module is marked as
// solved. Otherwise, it returns an error.
func (m *BigButtonModule) handleLongPressRelease() (handled bool, err error) {
	time := m.bomb.GetTimeLeft()
	if m.State.ReleaseDigit == nil {
		return true, errors.New("release digit is nil")
	} else {
		if int(time.Seconds())%10 == *m.State.ReleaseDigit {
			m.State.MarkAsSolved()
		} else {
			return true, errors.New("release digit does not match")
		}
	}

	return false, nil
}

var releaseDigits = [...]int{1, 4, 5}

func generateReleaseDigit() *int {
	digit := releaseDigits[rand.Intn(len(releaseDigits))]
	return &digit
}

func releaseDigitToStripColor(digit *int) (color *valueobject.Color, err error) {
	switch *digit {
	case 1:
		color = &bigButtonStripColors[rand.Intn(len(bigButtonStripColors))]
	case 4:
		value := valueobject.Blue
		color = &value
	case 5:
		value := valueobject.Red
		color = &value
	}

	if color != nil {
		return color, nil
	}

	return color, fmt.Errorf("invalid release digit: %d", *digit)
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
