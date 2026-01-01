package command

import "github.com/ZaneH/keep-talking/internal/domain/valueobject"

type MazeCommand struct {
	BaseModuleInputCommand
	Direction valueobject.CardinalDirection
}

type MazeCommandResult struct {
	BaseModuleInputCommandResult
	MazeMap valueobject.MazeMap
}
