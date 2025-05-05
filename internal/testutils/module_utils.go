package testutils

import (
	"testing"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/domain/factories"
)

func NewTestSimpleWireModule(t *testing.T) *actors.SimpleWireModuleActor {
	t.Helper()

	factory := factories.NewModuleFactory()
	module := factory.CreateSimpleWireModule(&factories.SimpleWireModuleConfig{
		WireCount:       3,
		SolutionIndices: []int{1},
	})

	return actors.NewSimpleWireModuleActor(module)
}

func NewTestPasswordModule(t *testing.T) *actors.PasswordModuleActor {
	t.Helper()

	factory := factories.NewModuleFactory()
	module := factory.CreatePasswordModule()
	module.SetSolution("test123")

	return actors.NewPasswordModuleActor(module)
}
