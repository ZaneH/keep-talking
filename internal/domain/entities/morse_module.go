package entities

import (
	"fmt"
	"log"
	"strings"

	"github.com/ZaneH/defuse.party-go/internal/domain/ports"
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
	"github.com/google/uuid"
)

type MorseState struct {
	BaseModuleState
	// Index of the selected frequency chosen by the player
	SelectedFrequencyIdx int
	// Currently selected frequency as a float
	DisplayedFrequency float32
	// Pattern displayed to the player in dots and dashes
	DisplayedPattern string
	// Solution frequency that the player must match to solve the module
	solution float32
}

func NewMorseState(rng ports.RandomGenerator) MorseState {
	solutionWord := morseWords[rng.GetIntInRange(0, len(morseWords)-1)]
	solutionFreq := morseWordToFrequency(solutionWord)

	startIdx := rng.GetIntInRange(0, len(morseWords)-1)
	startFreq := morseWordToFrequency(morseWords[startIdx])

	return MorseState{
		BaseModuleState:      BaseModuleState{},
		SelectedFrequencyIdx: startIdx,
		DisplayedFrequency:   startFreq,
		DisplayedPattern:     morseTranslations[solutionWord],
		solution:             solutionFreq,
	}
}

type MorseModule struct {
	BaseModule
	State MorseState
	rng   ports.RandomGenerator
}

func NewMorseModule(rng ports.RandomGenerator) *MorseModule {
	return &MorseModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		State: NewMorseState(rng),
		rng:   rng,
	}
}

func (m *MorseModule) String() string {
	selectedWord := morseWords[m.State.SelectedFrequencyIdx]
	var result strings.Builder
	result.WriteString("\n")
	result.WriteString(fmt.Sprintf("Pattern: %s", m.State.DisplayedPattern))
	result.WriteString(fmt.Sprintf("%f MHz\n", morseWordToFrequency(selectedWord)))

	return result.String()
}

func (m *MorseModule) GetType() valueobject.ModuleType {
	return valueobject.MorseModule
}

func (m *MorseModule) SetState(state MorseState) {
	m.State = state
}

func (m *MorseModule) GetModuleState() ModuleState {
	return &m.State
}

func (m *MorseModule) PressTx() (strike bool, err error) {
	selectedWord := morseWords[m.State.SelectedFrequencyIdx]

	if m.State.solution != morseWordToFrequency(selectedWord) {
		return true, nil
	}

	m.State.MarkAsSolved()
	return false, nil
}

func (m *MorseModule) PressChangeFrequency(d valueobject.IncrementDecrement) float32 {
	if d == valueobject.Increment {
		m.State.SelectedFrequencyIdx++
	} else {
		m.State.SelectedFrequencyIdx--
	}

	if m.State.SelectedFrequencyIdx < 0 {
		m.State.SelectedFrequencyIdx = 0
	} else if m.State.SelectedFrequencyIdx >= len(morseWords) {
		m.State.SelectedFrequencyIdx = len(morseWords) - 1
	}

	selectedWord := morseWords[m.State.SelectedFrequencyIdx]
	return morseWordToFrequency(selectedWord)
}

func (m *MorseModule) GetCurrentFrequency() float32 {
	selectedWord := morseWords[m.State.SelectedFrequencyIdx]
	return morseWordToFrequency(selectedWord)
}

func morseWordToFrequency(word string) float32 {
	solution, exists := morseFrequencies[word]
	if !exists {
		log.Fatalf("morse word %s does not have a frequency mapping", word)
	}

	return solution
}

var morseWords = []string{
	"shell", "halls", "slick", "trick", "boxes",
	"leaks", "strobe", "bistro", "flick", "bombs",
	"break", "brick", "steak", "sting", "vector",
	"beats",
}

var morseTranslations = map[string]string{
	"shell":  "... .... . .-.. .-..",
	"halls":  ".... .- .-.. .-.. ...",
	"slick":  "... .-.. .. -.-. -.-",
	"trick":  "- .-. .. -.-. -.-",
	"boxes":  "-... --- -..- . ...",
	"leaks":  ".-.. . .- -.- ...",
	"strobe": "... - .-. --- -... .",
	"bistro": "-... .. ... - .-. ---",
	"flick":  "..-. .-.. .. -.-. -.-",
	"bombs":  "-... --- -- -... ...",
	"break":  "-... .-. . .- -.-",
	"brick":  "-... .-. .. -.-. -.-",
	"steak":  "... - . .- -.-",
	"sting":  "... - .. -. --.",
	"vector": "...- . -.-. - --- .-.",
	"beats":  "-... . .- - ...",
}

var morseFrequencies = map[string]float32{
	"shell":  3.505,
	"halls":  3.515,
	"slick":  3.522,
	"trick":  3.532,
	"boxes":  3.535,
	"leaks":  3.542,
	"strobe": 3.545,
	"bistro": 3.552,
	"flick":  3.555,
	"bombs":  3.565,
	"break":  3.572,
	"brick":  3.575,
	"steak":  3.582,
	"sting":  3.592,
	"vector": 3.595,
	"beats":  3.600,
}
