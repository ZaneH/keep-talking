package services

import (
	"log"
	rand "math/rand/v2"

	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

type BombFactoryImpl struct{}

func (f *BombFactoryImpl) CreateBomb(config valueobject.BombConfig) *entities.Bomb {
	bomb := entities.NewBomb(config)

	totalModules := config.NumFaces * config.MaxModulesPerFace
	if totalModules > 9 { // Arbitrary maximum
		totalModules = 9
	}

	if totalModules < config.MinModules {
		totalModules = config.MinModules
	}

	moduleTypes := make([]valueobject.ModuleType, 0)
	weights := make([]float32, 0)

	for moduleType, probability := range config.ModuleTypes {
		moduleTypes = append(moduleTypes, moduleType)
		weights = append(weights, probability)
	}

	modulesToAdd := make([]valueobject.ModuleType, totalModules)
	modulesToAdd[0] = valueobject.Clock

	for i := 1; i < totalModules; i++ {
		modulesToAdd[i] = selectWeightedModuleType(moduleTypes, weights)
	}

	f.placeModulesOnBomb(bomb, modulesToAdd, config)

	return bomb
}

func (f *BombFactoryImpl) placeModulesOnBomb(bomb *entities.Bomb, moduleTypes []valueobject.ModuleType, config valueobject.BombConfig) {
	availablePositions := make([]valueobject.ModulePosition, 0)

	for face := 0; face < config.NumFaces; face++ {
		for row := 0; row < config.Rows; row++ {
			for col := 0; col < config.Columns; col++ {
				position := valueobject.ModulePosition{Row: row, Column: col, Face: face}
				availablePositions = append(availablePositions, position)
			}
		}
	}

	rand.Shuffle(len(availablePositions), func(i, j int) {
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
	case valueobject.Clock:
		module = entities.NewClockModule()
	case valueobject.SimpleWires:
		module = entities.NewSimpleWiresModule()
	case valueobject.Password:
		module = entities.NewPasswordModule(nil)
	case valueobject.BigButton:
		module = entities.NewBigButtonModule()
	case valueobject.SimonSays:
		module = entities.NewSimonSaysModule(nil)
	default:
		return nil
	}

	module.SetBomb(bomb)
	module.SetPosition(position)

	return module
}

func selectWeightedModuleType(moduleTypes []valueobject.ModuleType, weights []float32) valueobject.ModuleType {
	// Simple weighted selection algorithm
	totalWeight := float32(0)
	for _, w := range weights {
		totalWeight += w
	}

	r := rand.Float32() * totalWeight
	for i, w := range weights {
		r -= w
		if r <= 0 {
			return moduleTypes[i]
		}
	}

	return moduleTypes[0]
}
