package actors

import (
	"fmt"

	"github.com/ZaneH/keep-talking/internal/domain/entities"
)

func CreateModuleActor(bomb *entities.Bomb, module entities.Module) (ModuleActor, error) {
	switch module := module.(type) {
	case *entities.ClockModule:
		return NewStubModuleActor(module, 0), nil
	case *entities.SimpleWiresModule:
		return NewSimpleWiresModuleActor(module), nil
	case *entities.PasswordModule:
		return NewPasswordModuleActor(module), nil
	case *entities.BigButtonModule:
		return NewBigButtonModuleActor(module), nil
	case *entities.SimonSaysModule:
		return NewSimonSaysModuleActor(module), nil
	default:
		return nil, fmt.Errorf("unsupported module type: %v", module.GetType())
	}
}
