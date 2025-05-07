package entities

import (
	"errors"
	"math/rand"

	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type PasswordModuleState struct {
	Letters   [5][6]string // 5 letters, 6 options each
	Positions [5]int       // positions of the letters
	solution  string
	ModuleState
}

type PasswordModule struct {
	ModuleID uuid.UUID
	state    PasswordModuleState
}

func NewPasswordModule(providedSolution *string) *PasswordModule {
	var solution string
	if providedSolution == nil {
		solution = generateWord()
	} else {
		solution = *providedSolution
	}

	return &PasswordModule{
		ModuleID: uuid.New(),
		state: PasswordModuleState{
			Letters:   generateLetters(solution),
			Positions: [5]int{0, 0, 0, 0, 0},
			solution:  solution,
		},
	}
}

func (m *PasswordModule) GetCurrentGuess() string {
	guess := ""
	for i, pos := range m.state.Positions {
		guess += string(m.state.Letters[i][pos])
	}
	return guess
}

func (m *PasswordModule) IsSolved() bool {
	return m.state.solution == m.GetCurrentGuess()
}

func (m *PasswordModule) CheckPassword() error {
	if m.IsSolved() {
		m.state.MarkSolved = true
		return nil
	}

	return errors.New("incorrect password")
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
	var result string = "\n"
	for i := 0; i < len(m.state.Letters); i++ {
		result += "Letter " + string(rune('A'+i)) + ": "
		for j := 0; j < len(m.state.Letters[i]); j++ {
			if j == m.state.Positions[i] {
				result += "[" + m.state.Letters[i][j] + "] "
			} else {
				result += m.state.Letters[i][j] + " "
			}
		}
		result += "\n"
	}
	return result
}

func (m *PasswordModule) GetModuleID() uuid.UUID {
	return m.ModuleID
}

func (m *PasswordModule) GetType() valueobject.ModuleType {
	return valueobject.PASSWORD
}

var words = [...]string{
	"three",
	"apple",
	"world",
	"where",
	"water",
	"table",
	"chair",
}

func generateWord() string {
	randIdx := rand.Intn(len(words))
	return words[randIdx]
}

func generateLetters(solution string) [5][6]string {
	var letters [5][6]string

	for col := 0; col < len(solution) && col < 5; col++ {
		targetLetter := string(solution[col])
		targetPos := rand.Intn(6)

		// Track used letters in this column
		usedLetters := make(map[string]bool)
		usedLetters[targetLetter] = true

		for row := 0; row < 6; row++ {
			if row == targetPos {
				letters[col][row] = targetLetter
			} else {
				// Generate unique letter for this column
				var randomLetter string
				for {
					randomLetter = string(common.ALPHABET[rand.Intn(len(common.ALPHABET))])
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

		for row := 0; row < 6; row++ {
			var randomLetter string
			for {
				randomLetter = string(common.ALPHABET[rand.Intn(len(common.ALPHABET))])
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
