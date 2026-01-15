package command

import "github.com/ZaneH/defuse.party-go/internal/domain/valueobject"

type MazeCommand struct {
	BaseModuleInputCommand
	Direction valueobject.CardinalDirection
}

type MazeInputCommandResult struct {
	BaseModuleInputCommandResult
	PlayerPosition valueobject.Point2D
	GoalPosition   valueobject.Point2D
}
