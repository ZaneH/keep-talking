package command

import "github.com/ZaneH/keep-talking/internal/domain/valueobject"

type NeedyKnobCommand struct {
	BaseModuleInputCommand
}

type NeedyKnobCommandResult struct {
	BaseModuleInputCommandResult
	DisplayedPattern  [][]bool
	DialDirection     valueobject.CardinalDirection
	CoundownStartedAt int64
	CountdownDuration int16
}
