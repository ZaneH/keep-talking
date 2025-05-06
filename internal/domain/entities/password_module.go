package entities

import (
	"errors"
	"math/rand"

	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/google/uuid"
)

type PasswordModuleState struct {
	Letters   [5][6]string // 5 letters, 6 options each
	Positions [5]int       // positions of the letters
	solution  string
	common.ModuleState
}

type PasswordModule struct {
	ModuleID uuid.UUID
	state    PasswordModuleState
}

func NewPasswordModule(solution string) *PasswordModule {
	return &PasswordModule{
		ModuleID: uuid.New(),
		state: PasswordModuleState{
			Letters:   generateLetters(),
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

func (m *PasswordModule) CheckPassword() (bool, error) {
	if m.state.solution == "" {
		return false, errors.New("solution not set")
	}

	if m.GetCurrentGuess() == m.state.solution {
		return true, nil
	}

	return false, errors.New("incorrect password")
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

func (m *PasswordModule) SetSolution(solution string) {
	m.state.solution = solution
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

var commonLetters = "abcdefghijklmnopqrstuvwxyz"

func generateLetters() [5][6]string {
	randIdx := rand.Intn(len(words))
	word := words[randIdx]

	var letters [5][6]string

	for col := 0; col < len(word) && col < 5; col++ {
		targetLetter := string(word[col])
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
					randomLetter = string(commonLetters[rand.Intn(len(commonLetters))])
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
	for col := len(word); col < 5; col++ {
		usedLetters := make(map[string]bool)

		for row := 0; row < 6; row++ {
			var randomLetter string
			for {
				randomLetter = string(commonLetters[rand.Intn(len(commonLetters))])
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
