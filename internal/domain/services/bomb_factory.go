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

	moduleTypes := make([]valueobject.ModuleType, 0)
	weights := make([]float32, 0)

	for moduleType, probability := range config.ModuleTypes {
		moduleTypes = append(moduleTypes, moduleType)
		weights = append(weights, probability)
	}

	modulesToAdd := make([]valueobject.ModuleType, totalModules)
	for i := 0; i < totalModules; i++ {
		modulesToAdd[i] = selectWeightedModuleType(moduleTypes, weights)
	}

	f.placeModulesOnBomb(bomb, modulesToAdd, config)

	return bomb
}

func (f *BombFactoryImpl) placeModulesOnBomb(bomb *entities.Bomb, moduleTypes []valueobject.ModuleType, config valueobject.BombConfig) {
	availablePositions := make([]struct {
		face     int
		position valueobject.ModulePosition
	}, 0)

	for face := 0; face < config.NumFaces; face++ {
		for row := 0; row < config.MaxModulesPerFace; row++ {
			for col := 0; col < config.MaxModulesPerFace; col++ {
				position := valueobject.ModulePosition{Row: row, Column: col}
				availablePositions = append(availablePositions, struct {
					face     int
					position valueobject.ModulePosition
				}{face, position})
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

		location := availablePositions[i]
		module := f.createModule(moduleType)
		err := bomb.AddModule(module, location.face, location.position)
		if err != nil {
			log.Printf("error adding module to bomb: %v. skipping...", err)
			continue
		}
	}
}

func (f *BombFactoryImpl) createModule(moduleType valueobject.ModuleType) entities.Module {
	switch moduleType {
	case valueobject.SimpleWires:
		return entities.NewSimpleWiresModule()
	case valueobject.Password:
		return entities.NewPasswordModule(nil)
	default:
		return nil
	}
}

// Choose a module type randomly based on weights
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
