package entities

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"

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

func NewBigButtonModule() *BigButtonModule {
	return &BigButtonModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
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

func (m *BigButtonModule) PressButton(pressType valueobject.PressType, releaseTime int64) (stripColor *valueobject.Color, strike bool, err error) {
	handled := false

	switch pressType {
	case valueobject.PressTypeTap:
		handled = m.handleShortPress()
	case valueobject.PressTypeHold:
		handled, stripColor = m.handleLongPress()
	case valueobject.PressTypeRelease:
		handled, err = m.handleLongPressRelease(releaseTime)
	default:
		return nil, true, errors.New("invalid press type")
	}

	// If not handled by a specific case, fallback to random color and release digit
	if !handled {
		// Short tap was not handled, so it's a strike
		if pressType == valueobject.PressTypeTap {
			return nil, true, errors.New("invalid short press")
		}

		if pressType == valueobject.PressTypeHold {
			m.State.ReleaseDigit = generateReleaseDigit()
			stripColor, err = releaseDigitToStripColor(m.State.ReleaseDigit)
			if err != nil {
				log.Println("error generating strip color:", err)
				return nil, strike, err
			}
		}
	}

	return stripColor, strike, err
}

// Handles a short press (tap) of the button. There are certain conditions that must be met
// for the module to be marked as solved. If the conditions are not met, it will return an error.
func (m *BigButtonModule) handleShortPress() (handled bool) {
	if m.bomb.Batteries > 1 && m.State.Label == string(Detonate) {
		m.State.MarkAsSolved()
		return true
	}

	if m.bomb.Batteries > 2 && m.bomb.Indicators["FRK"].Lit {
		m.State.MarkAsSolved()
		return true
	}

	if m.State.ButtonColor == valueobject.Red && m.State.Label == string(Hold) {
		m.State.MarkAsSolved()
		return true
	}

	return false
}

// Generates a random strip color regardless of the button color or label. There is no
// possibility of a strike from this action. If the button is to be pressed and held, this
// function will also set the release digit to either 4, 1, or 5. If the button is not meant
// to be pressed and held, it will return nil for the strip color and no release digit will be set.
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
		// Unhandled, let the caller handle it
		return false, nil
	}

	return handled, color
}

// Handles the release of a long press. If the release digit is nil, it will return an error.
// If the release digit matches any digit of the bomb's timer, the module is marked as
// solved. Otherwise, it returns an error.
func (m *BigButtonModule) handleLongPressRelease(releaseTime int64) (handled bool, err error) {
	if m.State.ReleaseDigit == nil {
		return true, errors.New("release digit is nil")
	} else {
		// Check if any digit in MM:SS matches the release digit
		startedAt := m.bomb.StartedAt.Unix()
		duration := m.bomb.TimerDuration.Milliseconds() / 1000 // Convert to seconds

		// Calculate remaining time on the bomb (in seconds)
		elapsedTime := releaseTime - startedAt  // How many seconds have passed
		remainingTime := duration - elapsedTime // How many seconds remain

		minutes := remainingTime / 60
		seconds := remainingTime % 60
		releaseDigit := *m.State.ReleaseDigit

		// Check if MM:SS contains the release digit
		if strings.Contains(fmt.Sprintf("%02d:%02d", minutes, seconds), fmt.Sprintf("%d", releaseDigit)) {
			m.State.MarkAsSolved()
			return true, nil
		} else {
			return true, fmt.Errorf("release digit %d did not match any digit in MM:SS", releaseDigit)
		}
	}
}

var releaseDigits = [...]int{1, 4, 5}

func generateReleaseDigit() *int {
	digit := releaseDigits[rand.Intn(len(releaseDigits))]
	return &digit
}

func releaseDigitToStripColor(digit *int) (color *valueobject.Color, err error) {
	if digit == nil {
		return nil, errors.New("release digit is nil")
	}

	switch *digit {
	case 1:
		color = &bigButtonStripColors[rand.Intn(len(bigButtonStripColors))]
	case 4:
		value := valueobject.Blue
		color = &value
	case 5:
		value := valueobject.Yellow
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
