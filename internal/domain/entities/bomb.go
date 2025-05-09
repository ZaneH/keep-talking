package entities

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type Bomb struct {
	ID           uuid.UUID
	SerialNumber string
	StartingTime time.Duration
	StrikeCount  int
	MaxStrikes   int
	Faces        map[int]*BombFace
	Modules      map[uuid.UUID]Module
	Indicators   map[string]valueobject.Indicator
	Batteries    int
	Ports        []valueobject.Port
}

func NewBomb(config valueobject.BombConfig) *Bomb {
	return &Bomb{
		ID:           uuid.New(),
		SerialNumber: generateSerialNumber(),
		StartingTime: config.Timer,
		StrikeCount:  0,
		MaxStrikes:   config.MaxStrikes,
		Faces:        make(map[int]*BombFace),
		Modules:      make(map[uuid.UUID]Module),
		Indicators:   generateRandomIndicators(config.MaxIndicatorCount),
		Batteries:    generateRandomBatteryCount(config.MinBatteries, config.MaxBatteries),
		Ports:        generateRandomPorts(config.PortCount),
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

	log.Printf("adding module: %T", module)

	b.Modules[module.GetModuleID()] = module
	return nil
}

func (b *Bomb) GetTimeLeft() time.Duration {
	now := time.Now()
	return b.StartingTime - time.Since(now)
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

func generateRandomIndicators(count int) map[string]valueobject.Indicator {
	options := valueobject.AVAILABLE_INDICATOR_LABELS
	indicators := make(map[string]valueobject.Indicator, count)
	if count == 0 {
		return indicators
	}

	count = rand.Intn(count + 1)

	for range count {
		lit := rand.Intn(2) == 1
		randIdx := rand.Intn(len(options))
		label := options[randIdx]
		indicators[label] = valueobject.Indicator{
			Lit:   lit,
			Label: label,
		}
	}

	return indicators
}

func generateRandomBatteryCount(minBatteries int, maxBatteries int) int {
	if minBatteries == maxBatteries {
		return minBatteries
	}
	return rand.Intn(maxBatteries-minBatteries) + minBatteries
}

func generateRandomPorts(count int) []valueobject.Port {
	options := valueobject.AVAILABLE_PORTS
	ports := make([]valueobject.Port, 0, count)

	for range count {
		randIdx := rand.Intn(len(options))
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
