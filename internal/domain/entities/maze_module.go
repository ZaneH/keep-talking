package entities

import (
	"fmt"

	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type MazeState struct {
	BaseModuleState
	// (X, Y) coordinates of goal from [0, 5]
	GoalPosition [2]int
	// (X, Y) coordinates of player from [0, 5]
	PlayerPosition [2]int
	// [0, 8] determines the maze layout
	Variant int
}

func NewMazeState(rng ports.RandomGenerator) MazeState {
	return MazeState{
		BaseModuleState: BaseModuleState{},
		GoalPosition:    generateRandomPosition(rng),
		PlayerPosition:  generateRandomPosition(rng),
		Variant:         rng.GetIntInRange(0, 8),
	}
}

type MazeModule struct {
	BaseModule
	State MazeState
	rng   ports.RandomGenerator
}

func NewMazeModule(rng ports.RandomGenerator) *MazeModule {
	return &MazeModule{
		BaseModule: BaseModule{
			ModuleID: uuid.New(),
		},
		State: NewMazeState(rng),
		rng:   rng,
	}
}

func (m *MazeModule) GetModuleState() ModuleState {
	return &m.State
}

func (m *MazeModule) GetType() valueobject.ModuleType {
	return valueobject.Maze
}

func (m *MazeModule) SetState(state MazeState) {
	m.State = state
}

func (m *MazeModule) String() string {
	var result = "\n"

	result += fmt.Sprintf("Variant: %d\n", m.State.Variant)
	result += fmt.Sprintf("Player: (%d, %d) | Goal: (%d, %d)\n", m.State.PlayerPosition[0], m.State.PlayerPosition[1], m.State.GoalPosition[0], m.State.GoalPosition[1])

	return result
}

func (m *MazeModule) PressButton(dir valueobject.CardinalDirection) (mazeMap valueobject.MazeMap, strike bool, err error) {
	return maze1, true, nil
}

func generateRandomPosition(rng ports.RandomGenerator) [2]int {
	x := rng.GetIntInRange(0, 5)
	y := rng.GetIntInRange(0, 5)
	return [2]int{x, y}
}

var maze1 = valueobject.MazeMap{
	{
		{}, {Bottom: true}, {Right: true}, {}, {Bottom: true}, {Bottom: true, Right: true},
	},
	{
		{Right: true}, {}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true}, {Right: true},
	},
	{
		{Right: true}, {Bottom: true}, {Right: true}, {}, {Bottom: true}, {Right: true},
	},
	{
		{Right: true}, {Bottom: true}, {Bottom: true}, {Bottom: true, Right: true}, {Bottom: true}, {Right: true},
	},
	{
		{}, {Bottom: true}, {Right: true}, {}, {Bottom: true, Right: true}, {},
	},
	{
		{Bottom: true}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true, Right: true},
	},
}
