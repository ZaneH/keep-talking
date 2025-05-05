package entities

import (
	"errors"

	"github.com/google/uuid"
)

type PasswordModule struct {
	ModuleId uuid.UUID
	solution string
}

func NewPasswordModule(moduleId uuid.UUID, solution string) *PasswordModule {
	return &PasswordModule{
		ModuleId: moduleId,
		solution: solution,
	}
}

func (m *PasswordModule) CheckPassword(guess string) (bool, error) {
	if guess == m.solution {
		return true, nil
	}
	return false, errors.New("incorrect password")
}

func (m *PasswordModule) SetSolution(solution string) {
	m.solution = solution
}
