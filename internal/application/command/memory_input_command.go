package command

type MemoryInputCommand struct {
	BaseModuleInputCommand
	ButtonIndex int
}

type MemoryInputCommandResult struct {
	BaseModuleInputCommandResult
	ScreenNumber     int
	DisplayedNumbers []int
	Stage            int
}
