package valueobject

import "time"

type BombConfig struct {
	Timer             time.Duration
	MaxStrikes        int
	NumFaces          int
	MaxModulesPerFace int
	ModuleTypes       map[ModuleType]float32 // Module type to probability mapping
	MinBatteries      int
	MaxBatteries      int
	MaxIndicatorCount int
	PortCount         int
}

func NewDefaultBombConfig() BombConfig {
	return BombConfig{
		Timer:             5 * time.Minute,
		MaxStrikes:        3,
		NumFaces:          2,
		MaxModulesPerFace: 3,
		ModuleTypes: map[ModuleType]float32{
			SimpleWires: 0.24,
			Password:    0.1,
		},
		MinBatteries:      1,
		MaxBatteries:      4,
		MaxIndicatorCount: 2,
		PortCount:         3,
	}
}
