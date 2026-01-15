package entities

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ZaneH/defuse.party-go/internal/application/common"
	"github.com/ZaneH/defuse.party-go/internal/domain/ports"
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
	"github.com/google/uuid"
)

type Bomb struct {
	ID            uuid.UUID
	SerialNumber  string
	TimerDuration time.Duration
	StartedAt     *time.Time
	StrikeCount   int
	MaxStrikes    int
	Faces         map[int]*BombFace
	Modules       map[uuid.UUID]Module
	Indicators    map[string]valueobject.Indicator
	Batteries     int
	Ports         []valueobject.Port
}

func NewBomb(rng ports.RandomGenerator, config valueobject.BombConfig) *Bomb {
	return &Bomb{
		ID:            uuid.New(),
		SerialNumber:  generateSerialNumber(rng),
		TimerDuration: config.Timer,
		StartedAt:     nil,
		StrikeCount:   0,
		MaxStrikes:    config.MaxStrikes,
		Faces:         make(map[int]*BombFace),
		Modules:       make(map[uuid.UUID]Module),
		Indicators:    generateRandomIndicators(rng, config.MaxIndicatorCount),
		Batteries:     generateRandomBatteryCount(rng, config.MinBatteries, config.MaxBatteries),
		Ports:         generateRandomPorts(rng, config.PortCount),
	}
}

func (b *Bomb) AddModule(module Module, position valueobject.ModulePosition) error {
	face, exists := b.Faces[position.Face]

	if !exists {
		face = NewBombFace()
		b.Faces[position.Face] = face
	}

	if err := face.AddModuleAt(module, position); err != nil {
		return err
	}

	log.Printf("adding module: %T", module)

	b.Modules[module.GetModuleID()] = module
	return nil
}

func (b *Bomb) AddStrike() {
	b.StrikeCount++
}

func (b *Bomb) GetTimeLeft() time.Duration {
	now := time.Now()
	return b.TimerDuration - time.Since(now)
}

func (b *Bomb) StartTimer() {
	if b.StartedAt != nil {
		return
	}
	now := time.Now()
	b.StartedAt = &now
}

func generateSerialNumber(rng ports.RandomGenerator) string {
	var sb strings.Builder
	options := common.ALPHABET_UPPERCASE
	for range 6 {
		randIdx := rng.GetIntInRange(0, len(options)-1)
		sb.WriteByte(options[randIdx])
	}

	return sb.String()
}

func generateRandomIndicators(rng ports.RandomGenerator, count int) map[string]valueobject.Indicator {
	options := valueobject.AVAILABLE_INDICATOR_LABELS
	indicators := make(map[string]valueobject.Indicator, count)
	if count == 0 {
		return indicators
	}

	count = rng.GetIntInRange(0, count)

	for range count {
		lit := rng.GetIntInRange(0, 1) == 1
		randIdx := rng.GetIntInRange(0, len(options)-1)
		label := options[randIdx]
		indicators[label] = valueobject.Indicator{
			Lit:   lit,
			Label: label,
		}
	}

	return indicators
}

func generateRandomBatteryCount(rng ports.RandomGenerator, minBatteries int, maxBatteries int) int {
	if minBatteries >= maxBatteries {
		return minBatteries
	}
	return rng.GetIntInRange(minBatteries, maxBatteries)
}

func generateRandomPorts(rng ports.RandomGenerator, count int) []valueobject.Port {
	options := valueobject.AVAILABLE_PORTS
	ports := make([]valueobject.Port, 0, count)

	for range count {
		randIdx := rng.GetIntInRange(0, len(options)-1)
		ports = append(ports, options[randIdx])
	}

	return ports
}

func (b *Bomb) String() string {
	var sb strings.Builder
	sb.WriteString("Bomb ID: " + b.ID.String() + "\n")
	sb.WriteString("Serial Number: " + b.SerialNumber + "\n")
	sb.WriteString("Time Remaining: " + b.GetTimeLeft().String() + "\n")
	sb.WriteString("Strike Count: " + fmt.Sprint(b.StrikeCount) + "\n")
	sb.WriteString("Max Strikes: " + fmt.Sprint(b.MaxStrikes) + "\n")
	sb.WriteString("Batteries: " + fmt.Sprint(b.Batteries) + "\n")
	sb.WriteString("Ports: " + fmt.Sprintf("%+v", b.Ports) + "\n")

	for faceIndex, face := range b.Faces {
		sb.WriteString("Face " + fmt.Sprint(faceIndex) + ":\n" + face.String() + "\n")
	}

	return sb.String()
}
