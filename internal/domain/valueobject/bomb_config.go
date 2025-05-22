package valueobject

import "time"

type BombConfig struct {
	Timer             time.Duration
	MaxStrikes        int
	NumFaces          int
	MinModules        int
	MaxModulesPerFace int
	ModuleTypes       map[ModuleType]float32 // Module type to probability mapping
	MinBatteries      int
	MaxBatteries      int
	MaxIndicatorCount int
	PortCount         int
	Columns           int
	Rows              int
}

func NewDefaultBombConfig() BombConfig {
	return BombConfig{
		Timer:             5 * time.Minute,
		MaxStrikes:        3,
		NumFaces:          1,
		MinModules:        2,
		MaxModulesPerFace: 4,
		ModuleTypes: map[ModuleType]float32{
			SimpleWires: 0.1,
			// Password: 0.1,
			BigButton: 0.1,
			Clock:     0.0, // Will be added manually
		},
		MinBatteries:      1,
		MaxBatteries:      4,
		MaxIndicatorCount: 2,
		PortCount:         3,
		Columns:           3,
		Rows:              2,
	}
}
