package valueobject

import "time"

type BombConfig struct {
	// Duration of the bomb timer
	Timer time.Duration
	// Number of acceptable strikes before the bomb explodes
	MaxStrikes int
	// Number of faces on the bomb that hold modules
	NumFaces int
	// Minimum number of modules on the bomb
	MinModules int
	// Max number of modules that can be on a single face
	MaxModulesPerFace int
	// Module type to probability mapping
	ModuleTypes map[ModuleType]float32
	// Minimum number of batteries on the bomb
	MinBatteries int
	// Maximum number of batteries on the bomb
	MaxBatteries int
	// Max number of lit indicators that can appear on edgework
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
			ClockModule:        0.0, // One will be added manually
			WiresModule:        0.05,
			PasswordModule:     0.05,
			BigButtonModule:    0.05,
			KeypadModule:       0.05,
			SimonModule:        0.05,
			WhosOnFirstModule:  0.05,
			MemoryModule:       0.05,
			MorseModule:        0.05,
			NeedyVentGasModule: 0.05,
			NeedyKnobModule:    0.05,
			MazeModule:         0.05,
		},
		MinBatteries:      1,
		MaxBatteries:      4,
		MaxIndicatorCount: 2,
		PortCount:         3,
		Columns:           3,
		Rows:              2,
	}
}
