package command

type WhosOnFirstInputCommand struct {
	BaseModuleInputCommand
	Word string
}

type WhosOnFirstInputCommandResult struct {
	BaseModuleInputCommandResult
	ScreenWord  string
	ButtonWords []string
	Stage       int
}
