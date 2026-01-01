package entities

import (
	"fmt"
	"strings"

	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type MazeModuleState struct {
	BaseModuleState
	// (X, Y) coordinates of goal from [0, 5]
	GoalPosition valueobject.Point2D
	// (X, Y) coordinates of player from [0, 5]
	PlayerPosition valueobject.Point2D
	// [0, 8] determines the maze layout
	Variant int
}

func NewMazeState(rng ports.RandomGenerator) MazeModuleState {
	return MazeModuleState{
		BaseModuleState: BaseModuleState{},
		GoalPosition:    generateRandomPosition(rng),
		PlayerPosition:  generateRandomPosition(rng),
		Variant:         rng.GetIntInRange(0, 8),
	}
}

type MazeModule struct {
	BaseModule
	State MazeModuleState
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
	return valueobject.MazeModule
}

func (m *MazeModule) SetState(state MazeModuleState) {
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

func (m *MazeModule) PressDirection(dir valueobject.CardinalDirection) (mazeMap valueobject.Maze, strike bool, err error) {
	mazeMap = variantToMaze(m.State.Variant)

	currentX := m.State.PlayerPosition.X
	currentY := m.State.PlayerPosition.Y

	newX := currentX
	newY := currentY

	switch dir {
	case valueobject.North:
		newY--
	case valueobject.South:
		newY++
	case valueobject.East:
		newX++
	case valueobject.West:
		newX--
	}

	if newX < 0 || newX > 5 || newY < 0 || newY > 5 {
		return mazeMap, true, nil
	}

	hasWall := false

	switch dir {
	case valueobject.North:
		if currentY > 0 {
			hasWall = mazeMap.Map[newY][currentX].Bottom
		} else {
			hasWall = true // Top edge is always a wall
		}

	case valueobject.South:
		if currentY < 5 {
			hasWall = mazeMap.Map[currentY][currentX].Bottom
		} else {
			hasWall = true // Bottom edge is always a wall
		}

	case valueobject.East:
		if currentX < 5 {
			hasWall = mazeMap.Map[currentY][currentX].Right
		} else {
			hasWall = true // Right edge is always a wall
		}

	case valueobject.West:
		if currentX > 0 {
			hasWall = mazeMap.Map[currentY][newX].Right
		} else {
			hasWall = true // Left edge is always a wall
		}
	}

	if hasWall {
		return mazeMap, true, nil
	}

	m.State.PlayerPosition.X = newX
	m.State.PlayerPosition.Y = newY

	return mazeMap, false, nil
}

func generateRandomPosition(rng ports.RandomGenerator) valueobject.Point2D {
	x := rng.GetIntInRange(0, 5)
	y := rng.GetIntInRange(0, 5)
	return valueobject.Point2D{
		X: x,
		Y: y,
	}
}

func variantToMaze(v int) valueobject.Maze {
	return mazes[v]
}

var mazes = []valueobject.Maze{mazeA, mazeB, mazeC, mazeD, mazeE, mazeF, mazeG, mazeH, mazeI}

var mazeA = valueobject.Maze{
	Marker1: valueobject.Point2D{X: 0, Y: 1},
	Marker2: valueobject.Point2D{X: 5, Y: 2},
	Map: [6][6]valueobject.MazeCell{
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
	},
}

var mazeB = valueobject.Maze{
	Marker1: valueobject.Point2D{X: 1, Y: 3},
	Marker2: valueobject.Point2D{X: 4, Y: 1},
	Map: [6][6]valueobject.MazeCell{
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
	},
}

var mazeC = valueobject.Maze{
	Marker1: valueobject.Point2D{X: 3, Y: 3},
	Marker2: valueobject.Point2D{X: 5, Y: 3},
	Map: [6][6]valueobject.MazeCell{
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
	},
}

var mazeD = valueobject.Maze{
	Marker1: valueobject.Point2D{X: 0, Y: 0},
	Marker2: valueobject.Point2D{X: 0, Y: 3},
	Map: [6][6]valueobject.MazeCell{
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
	},
}

var mazeE = valueobject.Maze{
	Marker1: valueobject.Point2D{X: 4, Y: 2},
	Marker2: valueobject.Point2D{X: 3, Y: 5},
	Map: [6][6]valueobject.MazeCell{
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
	},
}

var mazeF = valueobject.Maze{
	Marker1: valueobject.Point2D{X: 2, Y: 4},
	Marker2: valueobject.Point2D{X: 4, Y: 0},
	Map: [6][6]valueobject.MazeCell{
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
	},
}

var mazeG = valueobject.Maze{
	Marker1: valueobject.Point2D{X: 1, Y: 0},
	Marker2: valueobject.Point2D{X: 1, Y: 5},
	Map: [6][6]valueobject.MazeCell{
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
	},
}

var mazeH = valueobject.Maze{
	Marker1: valueobject.Point2D{X: 3, Y: 0},
	Marker2: valueobject.Point2D{X: 2, Y: 3},
	Map: [6][6]valueobject.MazeCell{
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
	},
}

var mazeI = valueobject.Maze{
	Marker1: valueobject.Point2D{X: 2, Y: 1},
	Marker2: valueobject.Point2D{X: 0, Y: 4},
	Map: [6][6]valueobject.MazeCell{
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
			if maze.Map[y][x].Right {
				sb.WriteString("|")
			} else {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")

		// Bottom walls
		sb.WriteString("+")
		for x := range 6 {
			if maze.Map[y][x].Bottom {
				sb.WriteString("---+")
			} else {
				sb.WriteString("   +")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
