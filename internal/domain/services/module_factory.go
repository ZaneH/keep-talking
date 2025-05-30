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

func (f *ModuleFactory) CreateSimpleWiresModule() *entities.SimpleWiresModule {
	return entities.NewSimpleWiresModule(f.rng)
}

func (f *ModuleFactory) CreatePasswordModule() *entities.PasswordModule {
	return entities.NewPasswordModule(f.rng, nil)
}

func (f *ModuleFactory) CreateSimonSaysModule() *entities.SimonSaysModule {
	return entities.NewSimonSaysModule(f.rng, nil)
}

func (f *ModuleFactory) CreateBigButtonModule() *entities.BigButtonModule {
	return entities.NewBigButtonModule(f.rng)
}
