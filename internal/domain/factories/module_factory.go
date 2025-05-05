package factories

import (
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type ModuleFactory struct {
}

func NewModuleFactory() *ModuleFactory {
	return &ModuleFactory{}
}

func (mf *ModuleFactory) CreateSimpleWireModule(config *SimpleWireModuleConfig) *entities.SimpleWireModule {
	if config == nil {
		config = &SimpleWireModuleConfig{
			WireCount: 4,
			// TODO: Implement a better way to determine solution indices
			SolutionIndices: []int{2},
		}
	}

	wires := make([]valueobject.SimpleWire, config.WireCount)

	for i := 0; i < config.WireCount; i++ {
		color := valueobject.SimpleWireColors[i%len(valueobject.SimpleWireColors)]
		wires[i] = valueobject.SimpleWire{
			WireColor: color,
		}
	}

	return entities.NewSimpleWireModule(wires, config.SolutionIndices)
}

func (f *ModuleFactory) CreatePasswordModule() *entities.PasswordModule {
	moduleId := uuid.New()
	// TODO: Implement a real password logic
	solution := "password123"

	return entities.NewPasswordModule(moduleId, solution)
}

type SimpleWireModuleConfig struct {
	WireCount       int
	SolutionIndices []int
}
