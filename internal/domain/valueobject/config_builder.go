package valueobject

// Handles conversion from various config sources to BombConfig
type BombConfigBuilder struct{}

func NewBombConfigBuilder() *BombConfigBuilder {
	return &BombConfigBuilder{}
}

// Creates a BombConfig from a preset mission
func (b *BombConfigBuilder) FromMission(mission Mission) (BombConfig, error) {
	if err := ValidateMission(mission); err != nil {
		return BombConfig{}, err
	}

	def := MissionDefinitions[mission]

	// Calculate total modules for grid sizing
	totalModules := 0
	for _, spec := range def.Modules {
		totalModules += spec.Count
	}
	totalModules++ // +1 for clock

	// Set all module type weights to 0 (will use explicit list)
	moduleTypes := make(map[ModuleType]float32)
	moduleTypes[ClockModule] = 0.0

	config := BombConfig{
		Timer:             def.Timer,
		MaxStrikes:        def.MaxStrikes,
		NumFaces:          def.NumFaces,
		MinModules:        totalModules,
		MaxModulesPerFace: (totalModules + def.NumFaces - 1) / def.NumFaces,
		ModuleTypes:       moduleTypes,
		MinBatteries:      1,
		MaxBatteries:      3,
		MaxIndicatorCount: 2,
		PortCount:         3,
		Columns:           def.Columns,
		Rows:              def.Rows,
		ExplicitModules:   def.Modules,
		MissionSection:    def.Section,
	}

	return config, nil
}

// Creates a BombConfig from a difficulty level (1-10)
func (b *BombConfigBuilder) FromLevel(level int) (BombConfig, error) {
	if err := ValidateLevel(level); err != nil {
		return BombConfig{}, err
	}
	return BombConfigFromLevel(level), nil
}

func DefaultModuleWeights() map[ModuleType]float32 {
	return map[ModuleType]float32{
		ClockModule:            0.0,
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
		NeedyVentGasModule:     0.05,
		NeedyKnobModule:        0.05,
	}
}
