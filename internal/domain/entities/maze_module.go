package entities

import (
	"fmt"
	"strings"

	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type MazeState struct {
	BaseModuleState
	// (X, Y) coordinates of goal from [0, 5]
	GoalPosition valueobject.Point2D
	// (X, Y) coordinates of player from [0, 5]
	PlayerPosition valueobject.Point2D
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
	result += fmt.Sprintf("Player: (%d, %d) | Goal: (%d, %d)\n",
		m.State.PlayerPosition.X, m.State.PlayerPosition.Y,
		m.State.GoalPosition.X, m.State.GoalPosition.Y)

	result += m.mazeToString()

	return result
}

func (m *MazeModule) PressButton(dir valueobject.CardinalDirection) (mazeMap valueobject.MazeMap, strike bool, err error) {
	return mazeA, true, nil
}

func generateRandomPosition(rng ports.RandomGenerator) valueobject.Point2D {
	x := rng.GetIntInRange(0, 5)
	y := rng.GetIntInRange(0, 5)
	return valueobject.Point2D{
		X: x,
		Y: y,
	}
}

var mazeA = valueobject.MazeMap{
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
		{}, {Bottom: true}, {Right: true}, {}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{Bottom: true}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true, Right: true},
	},
}

var mazeB = valueobject.MazeMap{
	{
		{Bottom: true}, {}, {Bottom: true, Right: true}, {}, {}, {Bottom: true, Right: true},
	},
	{
		{}, {Bottom: true, Right: true}, {}, {Bottom: true, Right: true}, {Bottom: true}, {Right: true},
	},
	{
		{Right: true}, {}, {Bottom: true, Right: true}, {}, {Bottom: true}, {Right: true},
	},
	{
		{}, {Bottom: true, Right: true}, {}, {Bottom: true, Right: true}, {Right: true}, {Right: true},
	},
	{
		{Right: true}, {Right: true}, {Right: true}, {}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{Bottom: true, Right: true}, {Bottom: true}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true},
		{Bottom: true, Right: true},
	},
}

var mazeC = valueobject.MazeMap{
	{
		{}, {Bottom: true}, {Right: true}, {Right: true}, {}, {Right: true},
	},
	{
		{Bottom: true, Right: true}, {Right: true}, {Right: true}, {Bottom: true}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{}, {Right: true}, {Right: true}, {}, {Right: true}, {Right: true},
	},
	{
		{Right: true}, {Right: true}, {Right: true}, {Right: true}, {Right: true}, {Right: true},
	},
	{
		{Right: true}, {Bottom: true}, {Bottom: true, Right: true}, {Right: true}, {Right: true}, {Right: true},
	},
	{
		{Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true, Right: true}, {Bottom: true},
		{Bottom: true, Right: true},
	},
}

var mazeD = valueobject.MazeMap{
	{
		{}, {Right: true}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {Right: true},
	},
	{
		{Right: true}, {Right: true}, {}, {Bottom: true}, {Bottom: true}, {Right: true},
	},
	{
		{Right: true}, {Bottom: true}, {Bottom: true, Right: true}, {}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{Right: true}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {Right: true},
	},
	{
		{}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {Right: true}, {Right: true},
	},
	{
		{Bottom: true}, {Bottom: true}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true, Right: true},
		{Bottom: true, Right: true},
	},
}

var mazeE = valueobject.MazeMap{
	{
		{Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {}, {Right: true},
	},
	{
		{}, {Bottom: true}, {Bottom: true}, {}, {Bottom: true, Right: true}, {Bottom: true, Right: true},
	},
	{
		{}, {Right: true}, {Bottom: true}, {Bottom: true, Right: true}, {}, {Right: true},
	},
	{
		{Right: true}, {Bottom: true}, {Bottom: true}, {Right: true}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{Right: true}, {}, {Bottom: true}, {Bottom: true}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{Bottom: true, Right: true}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true},
		{Bottom: true, Right: true},
	},
}

var mazeF = valueobject.MazeMap{
	{
		{Right: true}, {}, {Right: true}, {Bottom: true}, {}, {Right: true},
	},
	{
		{Right: true}, {Right: true}, {Right: true}, {}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{}, {Bottom: true, Right: true}, {Bottom: true, Right: true}, {Right: true}, {}, {Bottom: true, Right: true},
	},
	{
		{Bottom: true}, {Right: true}, {}, {Right: true}, {Right: true}, {Right: true},
	},
	{
		{}, {Bottom: true, Right: true}, {Bottom: true, Right: true}, {Right: true}, {Bottom: true}, {Right: true},
	},
	{
		{Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true, Right: true}, {Bottom: true},
		{Bottom: true, Right: true},
	},
}

var mazeG = valueobject.MazeMap{
	{
		{}, {Bottom: true}, {Bottom: true}, {Right: true}, {}, {Right: true},
	},
	{
		{Right: true}, {}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{Bottom: true}, {Bottom: true, Right: true}, {}, {Bottom: true, Right: true}, {}, {Bottom: true, Right: true},
	},
	{
		{}, {Right: true}, {}, {Bottom: true}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{Right: true}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true}, {Right: true}, {Right: true},
	},
	{
		{Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true, Right: true},
	},
}

var mazeH = valueobject.MazeMap{
	{
		{Right: true}, {}, {Bottom: true}, {Right: true}, {}, {Right: true},
	},
	{
		{}, {Bottom: true}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{Right: true}, {}, {Bottom: true}, {Bottom: true}, {Right: true}, {Right: true},
	},
	{
		{Right: true}, {Bottom: true}, {Right: true}, {Bottom: true}, {Bottom: true}, {Bottom: true, Right: true},
	},
	{
		{Right: true}, {Right: true}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true, Right: true},
	},
	{
		{Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true}, {Bottom: true, Right: true},
	},
}

var mazeI = valueobject.MazeMap{
	{
		{Right: true}, {}, {Bottom: true}, {Bottom: true}, {}, {Right: true},
	},
	{
		{Right: true}, {Right: true}, {}, {Bottom: true, Right: true}, {Right: true}, {Right: true},
	},
	{
		{}, {Bottom: true}, {Bottom: true, Right: true}, {}, {Bottom: true, Right: true}, {Right: true},
	},
	{
		{Right: true}, {Right: true}, {}, {Bottom: true, Right: true}, {Bottom: true}, {Right: true},
	},
	{
		{Right: true}, {Right: true}, {Right: true}, {}, {Right: true}, {Bottom: true, Right: true},
	},
	{
		{Bottom: true}, {Bottom: true, Right: true}, {Bottom: true}, {Bottom: true, Right: true}, {Bottom: true},
		{Bottom: true, Right: true},
	},
}

func (m *MazeModule) mazeToString() string {
	var sb strings.Builder

	maze := mazeI

	// Top wall (always solid)
	sb.WriteString("+")
	for range 6 {
		sb.WriteString("---+")
	}
	sb.WriteString("\n")

	// For each row
	for y := range 6 {
		// Left wall (always solid) + cell contents
		sb.WriteString("|")
		for x := range 6 {
			// Cell content (player, goal, or empty)
			if m.State.PlayerPosition.X == x && m.State.PlayerPosition.Y == y {
				sb.WriteString(" P ")
			} else if m.State.GoalPosition.X == x && m.State.GoalPosition.Y == y {
				sb.WriteString(" G ")
			} else {
				sb.WriteString("   ")
			}

			// Right wall of this cell
			if maze[y][x].Right {
				sb.WriteString("|")
			} else {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")

		// Bottom walls
		sb.WriteString("+")
		for x := range 6 {
			if maze[y][x].Bottom {
				sb.WriteString("---+")
			} else {
				sb.WriteString("   +")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
