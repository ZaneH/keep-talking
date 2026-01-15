package valueobject

import "time"

// LevelConfigParams defines how each level scales
type LevelConfigParams struct {
	Timer             time.Duration
	MaxStrikes        int
	NumFaces          int
	MinModules        int
	MaxModulesPerFace int
	NeedyWeight       float32
	Rows              int
	Columns           int
}

// LevelConfigs maps level 1-10 to configuration parameters
var LevelConfigs = map[int]LevelConfigParams{
	1: {
		Timer:             5 * time.Minute,
		MaxStrikes:        3,
		NumFaces:          1,
		MinModules:        3,
		MaxModulesPerFace: 3,
		NeedyWeight:       0.0,
		Rows:              2,
		Columns:           3,
	},
	2: {
		Timer:             5 * time.Minute,
		MaxStrikes:        3,
		NumFaces:          1,
		MinModules:        4,
		MaxModulesPerFace: 4,
		NeedyWeight:       0.02,
		Rows:              2,
		Columns:           3,
	},
	3: {
		Timer:             4*time.Minute + 30*time.Second,
		MaxStrikes:        3,
		NumFaces:          1,
		MinModules:        5,
		MaxModulesPerFace: 5,
		NeedyWeight:       0.05,
		Rows:              2,
		Columns:           3,
	},
	4: {
		Timer:             4 * time.Minute,
		MaxStrikes:        3,
		NumFaces:          1,
		MinModules:        6,
		MaxModulesPerFace: 6,
		NeedyWeight:       0.08,
		Rows:              2,
		Columns:           3,
	},
	5: {
		Timer:             4 * time.Minute,
		MaxStrikes:        2,
		NumFaces:          2,
		MinModules:        8,
		MaxModulesPerFace: 5,
		NeedyWeight:       0.10,
		Rows:              2,
		Columns:           3,
	},
	6: {
		Timer:             3*time.Minute + 30*time.Second,
		MaxStrikes:        2,
		NumFaces:          2,
		MinModules:        10,
		MaxModulesPerFace: 5,
		NeedyWeight:       0.12,
		Rows:              2,
		Columns:           3,
	},
	7: {
		Timer:             3 * time.Minute,
		MaxStrikes:        2,
		NumFaces:          2,
		MinModules:        11,
		MaxModulesPerFace: 6,
		NeedyWeight:       0.15,
		Rows:              2,
		Columns:           3,
	},
	8: {
		Timer:             3 * time.Minute,
		MaxStrikes:        1,
		NumFaces:          3,
		MinModules:        14,
		MaxModulesPerFace: 5,
		NeedyWeight:       0.18,
		Rows:              2,
		Columns:           3,
	},
	9: {
		Timer:             2*time.Minute + 30*time.Second,
		MaxStrikes:        1,
		NumFaces:          3,
		MinModules:        16,
		MaxModulesPerFace: 6,
		NeedyWeight:       0.20,
		Rows:              2,
		Columns:           3,
	},
	10: {
		Timer:             2 * time.Minute,
		MaxStrikes:        1,
		NumFaces:          4,
		MinModules:        20,
		MaxModulesPerFace: 6,
		NeedyWeight:       0.25,
		Rows:              2,
		Columns:           3,
	},
}

func BombConfigFromLevel(level int) BombConfig {
	// Clamp level to valid range
	if level < 1 {
		level = 1
	}
	if level > 10 {
		level = 10
	}

	params := LevelConfigs[level]

	moduleTypes := map[ModuleType]float32{
		ClockModule:            0.0, // Always added manually
		WiresModule:            0.08,
		PasswordModule:         0.08,
		BigButtonModule:        0.08,
		KeypadModule:           0.08,
		SimonModule:            0.06,
		WhosOnFirstModule:      0.06,
		MemoryModule:           0.06,
		MorseModule:            0.05,
		MazeModule:             0.05,
		ComplicatedWiresModule: 0.06,
		WireSequenceModule:     0.05,
		NeedyVentGasModule:     params.NeedyWeight,
		NeedyKnobModule:        params.NeedyWeight,
	}

	return BombConfig{
		Timer:             params.Timer,
		MaxStrikes:        params.MaxStrikes,
		NumFaces:          params.NumFaces,
		MinModules:        params.MinModules,
		MaxModulesPerFace: params.MaxModulesPerFace,
		ModuleTypes:       moduleTypes,
		MinBatteries:      1,
		MaxBatteries:      2 + level/3,
		MaxIndicatorCount: 1 + level/4,
		PortCount:         2 + level/3,
		Columns:           params.Columns,
		Rows:              params.Rows,
	}
}
