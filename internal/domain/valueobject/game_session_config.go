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

func (c GameSessionConfig) Seed() string {
	return c.seed
}

func nonEmptySeed(seed string) string {
	if strings.TrimSpace(seed) == "" {
		return uuid.NewString()
	}
	return seed
}
