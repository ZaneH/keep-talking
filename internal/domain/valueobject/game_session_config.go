package valueobject

import (
	"strings"

	"github.com/google/uuid"
)

type GameSessionConfig struct {
	seed        string
	BombConfigs []BombConfig
}

func NewEasyGameSessionConfig(seed string) GameSessionConfig {
	return GameSessionConfig{
		seed: nonEmptySeed(seed),
		BombConfigs: []BombConfig{
			NewDefaultBombConfig(),
		},
	}
}

// Creates config for level-based games
func NewGameSessionConfigFromLevel(seed string, level int) (GameSessionConfig, error) {
	builder := NewBombConfigBuilder()
	bombConfig, err := builder.FromLevel(level)
	if err != nil {
		return GameSessionConfig{}, err
	}

	if errs := ValidateBombConfig(bombConfig); errs.HasErrors() {
		return GameSessionConfig{}, errs
	}

	return GameSessionConfig{
		seed:        nonEmptySeed(seed),
		BombConfigs: []BombConfig{bombConfig},
	}, nil
}

// Creates config for preset missions
func NewGameSessionConfigFromMission(seed string, mission Mission) (GameSessionConfig, error) {
	builder := NewBombConfigBuilder()
	bombConfig, err := builder.FromMission(mission)
	if err != nil {
		return GameSessionConfig{}, err
	}

	return GameSessionConfig{
		seed:        nonEmptySeed(seed),
		BombConfigs: []BombConfig{bombConfig},
	}, nil
}

// Creates config from custom specification
func NewGameSessionConfigFromCustom(seed string, bombConfig BombConfig) (GameSessionConfig, error) {
	if errs := ValidateBombConfig(bombConfig); errs.HasErrors() {
		return GameSessionConfig{}, errs
	}

	return GameSessionConfig{
		seed:        nonEmptySeed(seed),
		BombConfigs: []BombConfig{bombConfig},
	}, nil
}

func (c GameSessionConfig) Seed() string {
	return c.seed
}

func nonEmptySeed(seed string) string {
	if strings.TrimSpace(seed) == "" {
		return uuid.NewString()
	}
	return seed
}
