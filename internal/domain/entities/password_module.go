package entities

import (
	"strings"

	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type PasswordState struct {
	BaseModuleState
	Letters   [5][6]string // 5 letters, 6 options each
	Positions [5]int       // positions of the letters
	solution  string
}

type PasswordModule struct {
	BaseModule
	state PasswordState
	rng   ports.RandomGenerator
}

func NewPasswordModule(rng ports.RandomGenerator, providedSolution *string) *PasswordModule {
	var solution string
	if providedSolution == nil {
		solution = generateWord(rng)
	} else {
		solution = *providedSolution
	}

	return &PasswordModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		state: PasswordState{
			Letters:   generateLetters(rng, solution),
			Positions: [5]int{0, 0, 0, 0, 0},
			solution:  solution,
		},
		rng: rng,
	}
}

func (m *PasswordModule) GetCurrentGuess() string {
	var guess strings.Builder
	for i, pos := range m.state.Positions {
		guess.WriteString(string(m.state.Letters[i][pos]))
	}
	return guess.String()
}

func (m *PasswordModule) CheckPassword() (strike bool) {
	if m.state.solution == m.GetCurrentGuess() {
		m.state.MarkAsSolved()
		return false
	}

	return true
}

func (m *PasswordModule) IncrementLetterOption(letterIndex int) {
	m.state.Positions[letterIndex]++
	if m.state.Positions[letterIndex] >= len(m.state.Letters[letterIndex]) {
		m.state.Positions[letterIndex] = 0
	}
}

func (m *PasswordModule) DecrementLetterOption(letterIndex int) {
	m.state.Positions[letterIndex]--
	if m.state.Positions[letterIndex] < 0 {
		m.state.Positions[letterIndex] = len(m.state.Letters[letterIndex]) - 1
	}
}

func (m *PasswordModule) String() string {
	var result strings.Builder
	result.WriteString("\n")
	for i := range len(m.state.Letters) {
		result.WriteString("Letter " + string(rune('A'+i)) + ": ")
		for j := range len(m.state.Letters[i]) {
			if j == m.state.Positions[i] {
				result.WriteString("[" + m.state.Letters[i][j] + "] ")
			} else {
				result.WriteString(m.state.Letters[i][j] + " ")
			}
		}
		result.WriteString("\n")
	}
	return result.String()
}

func (m *PasswordModule) GetType() valueobject.ModuleType {
	return valueobject.PasswordModule
}

func (m *PasswordModule) GetModuleState() ModuleState {
	return &m.state
}

var availablePasswordList = [...]string{
	"about", "after", "again", "below", "could",
	"every", "first", "found", "great", "house",
	"large", "learn", "never", "other", "place",
	"plant", "point", "right", "small", "sound",
	"spell", "still", "study", "their", "there",
	"these", "thing", "think", "three", "water",
	"where", "which", "world", "would", "write",
}

func generateWord(rng ports.RandomGenerator) string {
	randIdx := rng.GetIntInRange(0, len(availablePasswordList)-1)
	return availablePasswordList[randIdx]
}

func generateLetters(rng ports.RandomGenerator, solution string) [5][6]string {
	var letters [5][6]string

	for col := 0; col < len(solution) && col < 5; col++ {
		targetLetter := string(solution[col])
		targetPos := rng.GetIntInRange(0, 5)

		// Track used letters in this column
		usedLetters := make(map[string]bool)
		usedLetters[targetLetter] = true

		for row := range 6 {
			if row == targetPos {
				letters[col][row] = targetLetter
			} else {
				// Generate unique letter for this column
				var randomLetter string
				for {
					randomLetter = string(common.ALPHABET[rng.GetIntInRange(0, len(common.ALPHABET)-1)])
					if !usedLetters[randomLetter] {
						usedLetters[randomLetter] = true
						break
					}
				}
				letters[col][row] = randomLetter
			}
		}
	}

	// Fill remaining columns with unique random letters
	for col := len(solution); col < 5; col++ {
		usedLetters := make(map[string]bool)

		for row := range 6 {
			var randomLetter string
			for {
				randomLetter = string(common.ALPHABET[rng.GetIntInRange(0, len(common.ALPHABET)-1)])
				if !usedLetters[randomLetter] {
					usedLetters[randomLetter] = true
					break
				}
			}
			letters[col][row] = randomLetter
		}
	}

	return letters
}
