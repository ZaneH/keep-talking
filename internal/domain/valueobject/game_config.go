package valueobject

import "fmt"

type GameConfig struct {
	ProbabilitySimpleWires float32
	ProbabilityPassword    float32

	ModulesPerBomb    int
	MaxModulesPerFace int
	NumFaces          int
}

func NewDefaultGameConfig() GameConfig {
	return GameConfig{
		ProbabilitySimpleWires: 0.24,
		ProbabilityPassword:    0.1,
		ModulesPerBomb:         3,
		MaxModulesPerFace:      3,
		NumFaces:               2,
	}
}

func (c *GameConfig) Validate() error {
	if c.ProbabilitySimpleWires < 0 || c.ProbabilitySimpleWires > 1 {
		return fmt.Errorf("probabilitySimpleWires must be between 0 and 1")
	}

	if c.ProbabilityPassword < 0 || c.ProbabilityPassword > 1 {
		return fmt.Errorf("probabilityPassword must be between 0 and 1")
	}

	return nil
}
