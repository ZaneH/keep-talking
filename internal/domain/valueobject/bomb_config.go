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
			Clock:        0.0, // One will be added manually
			Wires:        0.05,
			Password:     0.05,
			BigButton:    0.05,
			Keypad:       0.05,
			Simon:        0.05,
			WhosOnFirst:  0.05,
			Memory:       0.05,
			Morse:        0.05,
			NeedyVentGas: 0.05,
			NeedyKnob:    0.05,
			Maze:         0.05,
		},
		MinBatteries:      1,
		MaxBatteries:      4,
		MaxIndicatorCount: 2,
		PortCount:         3,
		Columns:           3,
		Rows:              2,
	}
}
