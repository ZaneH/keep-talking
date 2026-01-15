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

	var modulesToAdd []valueobject.ModuleType

	// Check if we have explicit modules (from missions)
	if len(config.ExplicitModules) > 0 {
		modulesToAdd = f.expandExplicitModules(rng, config.ExplicitModules, config.MissionSection)
	} else {
		modulesToAdd = f.generateWeightedModules(rng, config)
	}

	// Always prepend clock
	modulesToAdd = append([]valueobject.ModuleType{valueobject.ClockModule}, modulesToAdd...)

	f.placeModulesOnBomb(rng, bomb, modulesToAdd, config)

	return bomb
}

func (f *BombFactoryImpl) expandExplicitModules(rng ports.RandomGenerator, specs []valueobject.ModuleSpec, section int) []valueobject.ModuleType {
	var modules []valueobject.ModuleType
	for _, spec := range specs {
		if len(spec.PossibleTypes) > 0 {
			idx := rng.GetIntInRange(0, len(spec.PossibleTypes)-1)
			chosenType := spec.PossibleTypes[idx]
			for i := 0; i < spec.Count; i++ {
				modules = append(modules, chosenType)
			}
		} else if spec.Type == valueobject.RandomModule {
			pool := valueobject.SectionModulePools[section]
			if len(pool) == 0 {
				pool = []valueobject.ModuleType{
					valueobject.WiresModule,
					valueobject.PasswordModule,
					valueobject.BigButtonModule,
					valueobject.KeypadModule,
					valueobject.SimonModule,
					valueobject.WhosOnFirstModule,
					valueobject.MemoryModule,
					valueobject.MorseModule,
					valueobject.MazeModule,
				}
			}
			for i := 0; i < spec.Count; i++ {
				idx := rng.GetIntInRange(0, len(pool)-1)
				randomType := pool[idx]
				modules = append(modules, randomType)
			}
		} else {
			for i := 0; i < spec.Count; i++ {
				modules = append(modules, spec.Type)
			}
		}
	}
	return modules
}

func (f *BombFactoryImpl) generateWeightedModules(rng ports.RandomGenerator, config valueobject.BombConfig) []valueobject.ModuleType {
	totalModules := config.NumFaces * config.MaxModulesPerFace
	totalModules = max(totalModules, config.MinModules)

	moduleTypes := make([]valueobject.ModuleType, 0)
	weights := make([]float32, 0)

	for moduleType, probability := range config.ModuleTypes {
		if moduleType != valueobject.ClockModule { // Skip clock, it's added separately
			moduleTypes = append(moduleTypes, moduleType)
			weights = append(weights, probability)
		}
	}

	modules := make([]valueobject.ModuleType, totalModules-1) // -1 for clock
	for i := range modules {
		modules[i] = selectWeightedModuleType(f.moduleFactory.rng, moduleTypes, weights)
	}

	return modules
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
