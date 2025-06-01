package actors

import (
	"fmt"

	"github.com/ZaneH/keep-talking/internal/domain/entities"
)

func CreateModuleActor(bomb *entities.Bomb, module entities.Module) (ModuleActor, error) {
	switch module := module.(type) {
	case *entities.ClockModule:
		return NewStubModuleActor(module, 0), nil
	case *entities.WiresModule:
		return NewWiresModuleActor(module), nil
	case *entities.PasswordModule:
		return NewPasswordModuleActor(module), nil
	case *entities.BigButtonModule:
		return NewBigButtonModuleActor(module), nil
	case *entities.SimonModule:
		return NewSimonModuleActor(module), nil
	case *entities.KeypadModule:
		return NewKeypadModuleActor(module), nil
	case *entities.WhosOnFirstModule:
		return NewWhosOnFirstModuleActor(module), nil
	default:
		return nil, fmt.Errorf("unsupported module type: %v", module.GetType())
	}
}
