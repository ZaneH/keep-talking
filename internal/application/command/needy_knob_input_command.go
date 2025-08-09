package command

type NeedyKnobCommand struct {
	BaseModuleInputCommand
}

type NeedyKnobCommandResult struct {
	BaseModuleInputCommandResult
	DisplayedPattern  [][]bool
	CoundownStartedAt int64
	CountdownDuration int16
}
