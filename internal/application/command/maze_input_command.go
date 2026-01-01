package command

import "github.com/ZaneH/keep-talking/internal/domain/valueobject"

type MazeCommand struct {
	BaseModuleInputCommand
	Direction valueobject.CardinalDirection
}

type MazeInputCommandResult struct {
	BaseModuleInputCommandResult
	Maze           valueobject.Maze
	PlayerPosition valueobject.Point2D
}
