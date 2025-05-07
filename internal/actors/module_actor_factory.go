package actors

import (
	"fmt"

	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

func CreateModuleActor(moduleType valueobject.ModuleType) (ModuleActor, error) {
	switch moduleType {
	case valueobject.SimpleWires:
		return NewSimpleWiresModuleActor(), nil
	case valueobject.Password:
		return NewPasswordModuleActor(), nil
	default:
		return nil, fmt.Errorf("unsupported module type: %v", moduleType)
	}
}
