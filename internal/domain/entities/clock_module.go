package entities

import (
	"fmt"
	"time"

	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type ClockModule struct {
	BaseModule
}

func NewClockModule() *ClockModule {
	return &ClockModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
	}
}

func (m *ClockModule) GetType() valueobject.ModuleType {
	return valueobject.Clock
}

func (m *ClockModule) GetModuleState() ModuleState {
	return nil
}

func (m *ClockModule) String() string {
	var result string

	duration := m.GetBomb().TimerDuration
	started_at := m.GetBomb().StartedAt
	if started_at != nil {
		remaining := duration - time.Since(*started_at)
		if remaining < 0 {
			remaining = 0
		}
		result = fmt.Sprintf("Clock Module: %s", remaining)
	} else {
		result = "Clock Module: Not started"
	}

	return result
}
