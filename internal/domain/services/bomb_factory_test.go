package services_test

import (
	"testing"
	"time"

	"github.com/ZaneH/keep-talking/internal/domain/services"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestBombFactory_CreateBomb(t *testing.T) {
	// Arrange
	f := services.BombFactoryImpl{}
	c := valueobject.BombConfig{
		Timer:             1 * time.Minute,
		MaxStrikes:        3,
		NumFaces:          1,
		MaxModulesPerFace: 1,
		ModuleTypes: map[valueobject.ModuleType]float32{
			valueobject.SimpleWires: 0.5,
			valueobject.Password:    0.5,
		},
		MinBatteries:      0,
		MaxBatteries:      0,
		MaxIndicatorCount: 0,
		PortCount:         0,
	}

	// Act
	bomb := f.CreateBomb(c)

	// Assert
	if bomb == nil {
		t.Errorf("Expected bomb to be created, but got nil")
		return
	}

	assert.Equal(t, bomb.StartingTime, c.Timer, "Expected bomb timer to be %v, but got %v", c.Timer, bomb.StartingTime)
	assert.Equal(t, bomb.MaxStrikes, c.MaxStrikes, "Expected bomb max strikes to be %d, but got %d", c.MaxStrikes, bomb.MaxStrikes)
	assert.Equal(t, bomb.Batteries, c.MinBatteries, "Expected bomb batteries to be %d, but got %d", 0, bomb.Batteries)
	assert.Equal(t, bomb.Ports, []valueobject.Port{}, "Expected bomb ports to be empty, but got %v", bomb.Ports)
	assert.Equal(t, bomb.Indicators, map[string]valueobject.Indicator{}, "Expected bomb indicators to be empty, but got %v", bomb.Indicators)
	assert.Equal(t, len(bomb.Faces), 1)
}
