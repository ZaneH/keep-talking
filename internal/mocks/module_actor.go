package mocks

import (
	"context"

	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/google/uuid"
)

type MockModuleActor struct {
	ModuleID     uuid.UUID
	ProcessCount int
	Solved       bool
	LastCommand  interface{}
	ReturnValue  interface{}
	ReturnError  error
	Module       common.Module
}

func NewMockModuleActor() *MockModuleActor {
	return &MockModuleActor{
		ModuleID: uuid.New(),
	}
}

func (m *MockModuleActor) GetModuleID() uuid.UUID {
	return m.ModuleID
}

func (m *MockModuleActor) GetModule() common.Module {
	return m.Module
}

func (m *MockModuleActor) ProcessCommand(ctx context.Context, cmd interface{}) (interface{}, error) {
	m.ProcessCount++
	m.LastCommand = cmd
	return m.ReturnValue, m.ReturnError
}
