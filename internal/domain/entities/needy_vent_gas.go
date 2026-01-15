package entities

import (
	"strings"
	"time"

	"github.com/ZaneH/defuse.party-go/internal/domain/ports"
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
	"github.com/google/uuid"
)

type NeedyVentGasState struct {
	BaseModuleState
	// Question displayed to the player
	DisplayedQuestion string
	// Index of the displayed question
	questionIdx int8
	// Timestamp of when the current countdown started
	CountdownStartedAt int64
	// Duration in seconds of the countdown
	CountdownDuration int16
}

func NewNeedyVentGasState(rng ports.RandomGenerator) NeedyVentGasState {
	startQuestionIdx := rng.GetIntInRange(0, len(ventGasQuestions)-1)

	return NeedyVentGasState{
		BaseModuleState:    BaseModuleState{},
		DisplayedQuestion:  ventGasQuestions[startQuestionIdx],
		questionIdx:        int8(startQuestionIdx),
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  int16(30),
	}
}

type NeedyVentGasModule struct {
	BaseModule
	State NeedyVentGasState
	rng   ports.RandomGenerator
}

func NewNeedyVentGasModule(rng ports.RandomGenerator) *NeedyVentGasModule {
	return &NeedyVentGasModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		State: NewNeedyVentGasState(rng),
		rng:   rng,
	}
}

func (m *NeedyVentGasModule) String() string {
	var result strings.Builder
	result.WriteString("\n")
	// TODO

	return result.String()
}

func (m *NeedyVentGasModule) GetType() valueobject.ModuleType {
	return valueobject.NeedyVentGasModule
}

func (m *NeedyVentGasModule) SetState(state NeedyVentGasState) {
	m.State = state
}

func (m *NeedyVentGasModule) GetModuleState() ModuleState {
	return &m.State
}

func (m *NeedyVentGasModule) PressButton(input bool) (strike bool, err error) {
	// TODO: Factor in 2s delay

	a := ventGasAnswers[m.State.questionIdx]
	i := m.rng.GetIntInRange(0, int(len(ventGasQuestions)-1))
	m.State.DisplayedQuestion = ventGasQuestions[i]
	m.State.questionIdx = int8(i)
	m.State.CountdownStartedAt = time.Now().Unix()

	if input == a {
		return false, nil
	}

	return true, nil
}

func (m *NeedyVentGasModule) GetCurrentQuestion() string {
	return ventGasQuestions[m.State.questionIdx]
}

var ventGasQuestions = []string{
	"vent gas?",
	"detonate?",
}

var ventGasAnswers = []bool{
	true,
	false,
}
