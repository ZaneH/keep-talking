package services

import (
	"log"

	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

type BombFactoryImpl struct {
	moduleFactory *ModuleFactory
}

func NewBombFactory(moduleFactory *ModuleFactory) *BombFactoryImpl {
	return &BombFactoryImpl{
		moduleFactory: moduleFactory,
	}
}

func (f *BombFactoryImpl) CreateBomb(rng ports.RandomGenerator, config valueobject.BombConfig) *entities.Bomb {
	bomb := entities.NewBomb(rng, config)

	totalModules := config.NumFaces * config.MaxModulesPerFace
	totalModules = max(totalModules, config.MinModules)

	moduleTypes := make([]valueobject.ModuleType, 0)
	weights := make([]float32, 0)

	for moduleType, probability := range config.ModuleTypes {
		moduleTypes = append(moduleTypes, moduleType)
		weights = append(weights, probability)
	}

	modulesToAdd := make([]valueobject.ModuleType, totalModules)
	modulesToAdd[0] = valueobject.ClockModule

	for i := 1; i < totalModules; i++ {
		modulesToAdd[i] = selectWeightedModuleType(f.moduleFactory.rng, moduleTypes, weights)
	}

	f.placeModulesOnBomb(rng, bomb, modulesToAdd, config)

	return bomb
}

func (f *BombFactoryImpl) placeModulesOnBomb(rng ports.RandomGenerator, bomb *entities.Bomb, moduleTypes []valueobject.ModuleType, config valueobject.BombConfig) {
	availablePositions := make([]valueobject.ModulePosition, 0)

	for face := range config.NumFaces {
		for row := range config.Rows {
			for col := range config.Columns {
				position := valueobject.ModulePosition{Row: row, Column: col, Face: face}
				availablePositions = append(availablePositions, position)
			}
		}
	}

	rng.Shuffle(len(availablePositions), func(i, j int) {
		availablePositions[i], availablePositions[j] = availablePositions[j], availablePositions[i]
	})

	for i, moduleType := range moduleTypes {
		if i >= len(availablePositions) {
			break
		}

		position := availablePositions[i]
		module := f.createModule(bomb, moduleType, position)
		err := bomb.AddModule(module, position)
		if err != nil {
			log.Printf("error adding module to bomb: %v. skipping...", err)
			continue
		}
	}
}

func (f *BombFactoryImpl) createModule(bomb *entities.Bomb, moduleType valueobject.ModuleType, position valueobject.ModulePosition) entities.Module {
	var module entities.Module
	switch moduleType {
	case valueobject.ClockModule:
		module = f.moduleFactory.CreateClockModule()
	case valueobject.WiresModule:
		module = f.moduleFactory.CreateWiresModule()
	case valueobject.PasswordModule:
		module = f.moduleFactory.CreatePasswordModule()
	case valueobject.BigButtonModule:
		module = f.moduleFactory.CreateBigButtonModule()
	case valueobject.SimonModule:
		module = f.moduleFactory.CreateSimonModule()
	case valueobject.KeypadModule:
		module = f.moduleFactory.CreateKeypadModule()
	case valueobject.WhosOnFirstModule:
		module = f.moduleFactory.CreateWhosOnFirstModule()
	case valueobject.MemoryModule:
		module = f.moduleFactory.CreateMemoryModule()
	case valueobject.MorseModule:
		module = f.moduleFactory.CreateMorseModule()
	case valueobject.NeedyVentGasModule:
		module = f.moduleFactory.CreateNeedyVentGasModule()
	case valueobject.NeedyKnobModule:
		module = f.moduleFactory.CreateNeedyKnobModule()
	case valueobject.MazeModule:
		module = f.moduleFactory.CreateMazeModule()
	default:
		log.Printf("unknown module type %v, skipping...", moduleType)
		return nil
	}

	module.SetBomb(bomb)
	module.SetPosition(position)

	return module
}

func selectWeightedModuleType(rng ports.RandomGenerator, moduleTypes []valueobject.ModuleType, weights []float32) valueobject.ModuleType {
	// Simple weighted selection algorithm
	totalWeight := float32(0)
	for _, w := range weights {
		totalWeight += w
	}

	r := rng.Float32(0, totalWeight)
	for i, w := range weights {
		r -= w
		if r <= 0 {
			return moduleTypes[i]
		}
	}

	return moduleTypes[0]
}
