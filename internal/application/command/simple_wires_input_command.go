package command

type SimpleWiresInputCommand struct {
	BaseModuleInputCommand
	WireIndex int
}

type SimpleWiresInputCommandResult struct {
	Solved bool
	Strike bool
}
