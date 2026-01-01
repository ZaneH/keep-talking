package entities

import (
	"fmt"
	"slices"

	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

const nWhosOnFirstStages = 3
const nWhosOnFirstWords = 6

type WhosOnFirstState struct {
	BaseModuleState
	// Word on screen for the defuser to read aloud.
	ScreenWord string
	// 6 words displayed to the player.
	ButtonWords []string
	// Tracks the current stage, [1, 3].
	Stage int
}

func NewWhosOnFirstState(rng ports.RandomGenerator) WhosOnFirstState {
	return WhosOnFirstState{
		BaseModuleState: BaseModuleState{},
		ScreenWord:      generateScreenWord(rng),
		ButtonWords:     generateButtonWords(rng),
		Stage:           1,
	}
}

type WhosOnFirstModule struct {
	BaseModule
	State WhosOnFirstState
	rng   ports.RandomGenerator
}

func NewWhosOnFirstModule(rng ports.RandomGenerator) *WhosOnFirstModule {
	return &WhosOnFirstModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		State: NewWhosOnFirstState(rng),
		rng:   rng,
	}
}

func (m *WhosOnFirstModule) GetModuleState() ModuleState {
	return &m.State
}

func (m *WhosOnFirstModule) GetType() valueobject.ModuleType {
	return valueobject.WhosOnFirstModule
}

func (m *WhosOnFirstModule) SetState(state WhosOnFirstState) {
	m.State = state
}

func (m *WhosOnFirstModule) String() string {
	var result = "\n"
	result += fmt.Sprintf("Stage: %d\n", m.State.Stage)
	result += "Screen Word: " + m.State.ScreenWord + "\n"
	result += "Displayed Words: "
	for i, word := range m.State.ButtonWords {
		if i > 0 {
			result += ", "
		}
		result += word
	}

	return result
}

func (m *WhosOnFirstModule) PressWord(word string) (strike bool, err error) {
	lookSpotIdx, exists := whosOnFirstScreenWordsToLocation[m.State.ScreenWord]
	if !exists {
		return false, fmt.Errorf("invalid screen word: %s", m.State.ScreenWord)
	}

	lookWord := m.State.ButtonWords[lookSpotIdx]

	validWords, exists := whosOnFirstSolutions[lookWord]
	if !exists {
		return false, fmt.Errorf("no valid words exist for input word: %s", word)
	}

	earliestIdx := -1
	for _, word := range m.State.ButtonWords {
		idx := slices.Index(validWords, word)
		if idx != -1 {
			// If the word is found, check if it's the earliest valid word
			if earliestIdx == -1 || idx < earliestIdx {
				earliestIdx = idx
			}
		}
	}

	earliestWord := validWords[earliestIdx]
	if word != earliestWord {
		m.State.ButtonWords = generateButtonWords(m.rng)
		m.State.ScreenWord = generateScreenWord(m.rng)
		m.State.Stage = 1

		return true, nil
	}

	m.State.ButtonWords = generateButtonWords(m.rng)
	m.State.ScreenWord = generateScreenWord(m.rng)
	m.State.Stage++
	if m.State.Stage >= nWhosOnFirstStages+1 {
		m.State.MarkAsSolved()
		return false, nil
	}

	return false, nil
}

func generateButtonWords(rng ports.RandomGenerator) []string {
	buttonWords := make([]string, 0, nWhosOnFirstWords)
	for len(buttonWords) < nWhosOnFirstWords {
		word := whosOnFirstButtonWords[rng.GetIntInRange(0, len(whosOnFirstButtonWords)-1)]
		if slices.Contains(buttonWords, word) {
			continue
		}

		buttonWords = append(buttonWords, word)
	}
	return buttonWords
}

func generateScreenWord(rng ports.RandomGenerator) string {
	return whosOnFirstScreenWords[rng.GetIntInRange(0, len(whosOnFirstScreenWords)-1)]
}

var whosOnFirstScreenWords = [...]string{
	"YES", "FIRST", "DISPLAY", "OKAY",
	"SAYS", "NOTHING", "", "BLANK",
	"NO", "LED", "LEAD", "READ",
	"RED", "REED", "LEED", "HOLD ON",
	"YOU", "YOU ARE", "YOUR", "YOU’RE",
	"UR", "THERE", "THEY’RE", "THEIR",
	"THEY ARE", "SEE", "C", "CEE",
}

var whosOnFirstScreenWordsToLocation = map[string]int{
	"YES":      2,
	"FIRST":    1,
	"DISPLAY":  5,
	"OKAY":     1,
	"SAYS":     5,
	"NOTHING":  2,
	"":         4,
	"BLANK":    3,
	"NO":       5,
	"LED":      2,
	"LEAD":     5,
	"READ":     3,
	"RED":      3,
	"REED":     4,
	"LEED":     4,
	"HOLD ON":  5,
	"YOU":      4,
	"YOU ARE":  5,
	"YOUR":     3,
	"YOU’RE":   3,
	"UR":       0,
	"THERE":    5,
	"THEY’RE":  4,
	"THEIR":    3,
	"THEY ARE": 2,
	"SEE":      5,
	"C":        1,
	"CEE":      5,
}

var whosOnFirstButtonWords = [...]string{
	"READY", "FIRST", "NO", "BLANK",
	"NOTHING", "YES", "WHAT", "UHHH",
	"LEFT", "RIGHT", "MIDDLE", "OKAY",
	"WAIT", "PRESS", "YOU", "YOU ARE",
	"YOUR", "YOU’RE", "UR", "U", "UH HUH",
	"UH UH", "WHAT?", "DONE", "NEXT",
	"HOLD", "SURE", "LIKE",
}

var whosOnFirstSolutions = map[string][]string{
	"READY":   {"YES", "OKAY", "WHAT", "MIDDLE", "LEFT", "PRESS", "RIGHT", "BLANK", "READY", "NO", "FIRST", "UHHH", "NOTHING", "WAIT"},
	"FIRST":   {"LEFT", "OKAY", "YES", "MIDDLE", "NO", "RIGHT", "NOTHING", "UHHH", "WAIT", "READY", "BLANK", "WHAT", "PRESS", "FIRST"},
	"NO":      {"BLANK", "UHHH", "WAIT", "FIRST", "WHAT", "READY", "RIGHT", "YES", "NOTHING", "LEFT", "PRESS", "OKAY", "NO", "MIDDLE"},
	"BLANK":   {"WAIT", "RIGHT", "OKAY", "MIDDLE", "BLANK", "PRESS", "READY", "NOTHING", "NO", "WHAT", "LEFT", "UHHH", "YES", "FIRST"},
	"NOTHING": {"UHHH", "RIGHT", "OKAY", "MIDDLE", "YES", "BLANK", "NO", "PRESS", "LEFT", "WHAT", "WAIT", "FIRST", "NOTHING", "READY"},
	"YES":     {"OKAY", "RIGHT", "UHHH", "MIDDLE", "FIRST", "WHAT", "PRESS", "READY", "NOTHING", "YES", "LEFT", "BLANK", "NO", "WAIT"},
	"WHAT":    {"UHHH", "WHAT", "LEFT", "NOTHING", "READY", "BLANK", "MIDDLE", "NO", "OKAY", "FIRST", "WAIT", "YES", "PRESS", "RIGHT"},
	"UHHH":    {"READY", "NOTHING", "LEFT", "WHAT", "OKAY", "YES", "RIGHT", "NO", "PRESS", "BLANK", "UHHH", "MIDDLE", "WAIT", "FIRST"},
	"LEFT":    {"RIGHT", "LEFT", "FIRST", "NO", "MIDDLE", "YES", "BLANK", "WHAT", "UHHH", "WAIT", "PRESS", "READY", "OKAY", "NOTHING"},
	"RIGHT":   {"YES", "NOTHING", "READY", "PRESS", "NO", "WAIT", "WHAT", "RIGHT", "MIDDLE", "LEFT", "UHHH", "BLANK", "OKAY", "FIRST"},
	"MIDDLE":  {"BLANK", "READY", "OKAY", "WHAT", "NOTHING", "PRESS", "NO", "WAIT", "LEFT", "MIDDLE", "RIGHT", "FIRST", "UHHH", "YES"},
	"OKAY":    {"MIDDLE", "NO", "FIRST", "YES", "UHHH", "NOTHING", "WAIT", "OKAY", "LEFT", "READY", "BLANK", "PRESS", "WHAT", "RIGHT"},
	"WAIT":    {"UHHH", "NO", "BLANK", "OKAY", "YES", "LEFT", "FIRST", "PRESS", "WHAT", "WAIT", "NOTHING", "READY", "RIGHT", "MIDDLE"},
	"PRESS":   {"RIGHT", "MIDDLE", "YES", "READY", "PRESS", "OKAY", "NOTHING", "UHHH", "BLANK", "LEFT", "FIRST", "WHAT", "NO", "WAIT"},
	"YOU":     {"SURE", "YOU ARE", "YOUR", "YOU’RE", "NEXT", "UH HUH", "UR", "HOLD", "WHAT?", "YOU", "UH UH", "LIKE", "DONE", "U"},
	"YOU ARE": {"YOUR", "NEXT", "LIKE", "UH HUH", "WHAT?", "DONE", "UH UH", "HOLD", "YOU", "U", "YOU’RE", "SURE", "UR", "YOU ARE"},
	"YOUR":    {"UH UH", "YOU ARE", "UH HUH", "YOUR", "NEXT", "UR", "SURE", "U", "YOU’RE", "YOU", "WHAT?", "HOLD", "LIKE", "DONE"},
	"YOU’RE":  {"YOU", "YOU’RE", "UR", "NEXT", "UH UH", "YOU ARE", "U", "YOUR", "WHAT?", "UH HUH", "SURE", "DONE", "LIKE", "HOLD"},
	"UR":      {"ONE", "U", "UR", "UH HUH", "WHAT?", "SURE", "YOUR", "HOLD", "YOU’RE", "LIKE", "NEXT", "UH UH", "YOU ARE", "YOU"},
	"U":       {"UH HUH", "SURE", "NEXT", "WHAT?", "YOU’RE", "UR", "UH UH", "DONE", "U", "YOU", "LIKE", "HOLD", "YOU ARE", "YOUR"},
	"UH HUH":  {"UH HUH", "YOUR", "YOU ARE", "YOU", "DONE", "HOLD", "UH UH", "NEXT", "SURE", "LIKE", "YOU’RE", "UR", "U", "WHAT?"},
	"UH UH":   {"UR", "U", "YOU ARE", "YOU’RE", "NEXT", "UH UH", "DONE", "YOU", "UH HUH", "LIKE", "YOUR", "SURE", "HOLD", "WHAT?"},
	"WHAT?":   {"YOU", "HOLD", "YOU’RE", "YOUR", "U", "DONE", "UH UH", "LIKE", "YOU ARE", "UH HUH", "UR", "NEXT", "WHAT?", "SURE"},
	"DONE":    {"SURE", "UH HUH", "NEXT", "WHAT?", "YOUR", "UR", "YOU’RE", "HOLD", "LIKE", "YOU", "U", "YOU ARE", "UH UH", "DONE"},
	"NEXT":    {"WHAT?", "UH HUH", "UH UH", "YOUR", "HOLD", "SURE", "NEXT", "LIKE", "DONE", "YOU ARE", "UR", "YOU’RE", "U", "YOU"},
	"HOLD":    {"YOU ARE", "U", "DONE", "UH UH", "YOU", "UR", "SURE", "WHAT?", "YOU’RE", "NEXT", "HOLD", "UH HUH", "YOUR", "LIKE"},
	"SURE":    {"YOU ARE", "DONE", "LIKE", "YOU’RE", "YOU", "HOLD", "UH HUH", "UR", "SURE", "U", "WHAT?", "NEXT", "YOUR", "UH UH"},
	"LIKE":    {"YOU’RE", "NEXT", "U", "UR", "HOLD", "DONE", "UH UH", "WHAT?", "UH HUH", "YOU", "LIKE", "SURE", "YOU ARE", "YOUR"},
}
