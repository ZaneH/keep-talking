package command

type NeedyVentGasCommand struct {
	BaseModuleInputCommand
	// True for Y, false for N
	Input bool
}

type NeedyVentGasCommandResult struct {
	BaseModuleInputCommandResult
	DisplayedQuestion  string
	CountdownStartedAt int64
	CountdownDuration  int16
}
