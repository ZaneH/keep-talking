package valueobject

import "time"

type BombConfig struct {
	// Duration of the bomb timer
	Timer      time.Duration
	MaxStrikes int
	NumFaces   int
	// Minimum number of modules on the bomb
	MinModules        int
	MaxModulesPerFace int
	// Module type to probability mapping
	ModuleTypes map[ModuleType]float32
	// Minimum number of batteries on the bomb
	MinBatteries int
	// Maximum number of batteries on the bomb
	MaxBatteries      int
	MaxIndicatorCount int
	// Maximum number of ports on the bomb
	PortCount int
	// Max number of columns for modules
	Columns int
	// Max number of rows for modules
	Rows int
}

func NewDefaultBombConfig() BombConfig {
	return BombConfig{
		Timer:             5 * time.Minute,
		MaxStrikes:        3,
		NumFaces:          1,
		MinModules:        4,
		MaxModulesPerFace: 3,
		ModuleTypes: map[ModuleType]float32{
			Clock:       0.0, // One will be added manually
			Wires:       0.05,
			Password:    0.05,
			BigButton:   0.05,
			Keypad:      0.05,
			Simon:       0.05,
			WhosOnFirst: 0.05,
			Memory:      0.05,
		},
		MinBatteries:      1,
		MaxBatteries:      4,
		MaxIndicatorCount: 2,
		PortCount:         3,
		Columns:           3,
		Rows:              2,
	}
}
