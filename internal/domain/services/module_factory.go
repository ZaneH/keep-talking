package services

import (
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/ports"
)

type ModuleFactory struct {
	rng ports.RandomGenerator
}

func NewModuleFactory(rng ports.RandomGenerator) *ModuleFactory {
	return &ModuleFactory{rng: rng}
}

func (f *ModuleFactory) CreateClockModule() *entities.ClockModule {
	return entities.NewClockModule()
}

func (f *ModuleFactory) CreateWiresModule() *entities.WiresModule {
	return entities.NewWiresModule(f.rng)
}

func (f *ModuleFactory) CreatePasswordModule() *entities.PasswordModule {
	return entities.NewPasswordModule(f.rng, nil)
}

func (f *ModuleFactory) CreateSimonModule() *entities.SimonModule {
	return entities.NewSimonModule(f.rng, nil)
}

func (f *ModuleFactory) CreateBigButtonModule() *entities.BigButtonModule {
	return entities.NewBigButtonModule(f.rng)
}

func (f *ModuleFactory) CreateKeypadModule() *entities.KeypadModule {
	return entities.NewKeypadModule(f.rng)
}

func (f *ModuleFactory) CreateWhosOnFirstModule() *entities.WhosOnFirstModule {
	return entities.NewWhosOnFirstModule(f.rng)
}

func (f *ModuleFactory) CreateMemoryModule() *entities.MemoryModule {
	return entities.NewMemoryModule(f.rng)
}

func (f *ModuleFactory) CreateMorseModule() *entities.MorseModule {
	return entities.NewMorseModule(f.rng)
}

func (f *ModuleFactory) CreateNeedyVentGasModule() *entities.NeedyVentGasModule {
	return entities.NewNeedyVentGasModule(f.rng)
}

func (f *ModuleFactory) CreateNeedyKnobModule() *entities.NeedyKnobModule {
	return entities.NewNeedyKnobModule(f.rng)
}
