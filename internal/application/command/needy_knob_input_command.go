package command

import "github.com/ZaneH/defuse.party-go/internal/domain/valueobject"

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
