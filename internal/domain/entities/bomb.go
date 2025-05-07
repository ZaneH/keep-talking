package entities

import (
	"math/rand"
	"strings"
	"time"

	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type Bomb struct {
	ID            uuid.UUID
	SerialNumber  string
	TimeRemaining time.Duration
	StartingTime  time.Duration
	StrikeCount   int
	MaxStrikes    int
	Faces         map[int]*BombFace
	Modules       map[uuid.UUID]Module
	Indicators    []valueobject.Indicator
	Batteries     int
	Ports         []valueobject.Port
}

func NewBomb(config valueobject.BombConfig) *Bomb {
	return &Bomb{
		ID:            uuid.New(),
		SerialNumber:  generateSerialNumber(),
		TimeRemaining: config.Timer,
		StartingTime:  config.Timer,
		StrikeCount:   0,
		MaxStrikes:    config.MaxStrikes,
		Modules:       make(map[uuid.UUID]Module),
		Indicators:    generateRandomIndicators(),
		Batteries:     generateRandomBatteryCount(config.MinBatteries, config.MaxBatteries),
		Ports:         generateRandomPorts(config.PortCount),
	}
}

func (b *Bomb) AddModule(module Module, faceIndex int, position valueobject.ModulePosition) error {
	face, exists := b.Faces[faceIndex]

	if !exists {
		face = NewBombFace()
		b.Faces[faceIndex] = face
	}

	if err := face.AddModule(module, position); err != nil {
		return err
	}

	b.Modules[module.GetModuleID()] = module
	return nil
}

func generateSerialNumber() string {
	var sb strings.Builder
	options := common.ALPHABET_UPPERCASE
	for range 6 {
		randIdx := rand.Intn(len(options))
		sb.WriteByte(options[randIdx])
	}

	return sb.String()
}

func generateRandomIndicators() []valueobject.Indicator {
	options := valueobject.AVAILABLE_INDICATOR_LABELS
	var indicators []valueobject.Indicator

	for range 5 {
		// Will result in 2.5 indicators on average
		if rand.Intn(2) == 1 {
			continue
		}

		lit := rand.Intn(2) == 1
		randIdx := rand.Intn(len(words))
		indicators = append(indicators, valueobject.Indicator{
			Lit:   lit,
			Label: options[randIdx],
		})
	}

	return indicators
}

func generateRandomBatteryCount(minBatteries int, maxBatteries int) int {
	return rand.Intn(maxBatteries-minBatteries) + minBatteries
}

func generateRandomPorts(count int) []valueobject.Port {
	options := valueobject.AVAILABLE_PORTS
	var ports []valueobject.Port

	for range count {
		randIdx := rand.Intn(len(options))
		ports = append(ports, options[randIdx])
	}

	return ports
}
