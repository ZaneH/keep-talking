package command

import (
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
	"github.com/google/uuid"
)

type ConfigType int

const (
	ConfigTypeDefault ConfigType = iota
	ConfigTypeLevel
	ConfigTypeMission
	ConfigTypeCustom
)

type CreateGameCommand struct {
	Seed       string
	ConfigType ConfigType

	// Level-based config (1-10)
	Level int

	// Mission-based config
	Mission valueobject.Mission

	// Custom config
	CustomConfig *valueobject.BombConfig
}

type CreateGameCommandResult struct {
	SessionID uuid.UUID
	Config    valueobject.BombConfig
}
