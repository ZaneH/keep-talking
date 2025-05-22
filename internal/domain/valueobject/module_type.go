package valueobject

type ModuleType int32

const (
	SimpleWires ModuleType = iota
	Password
	BigButton
	SimonSays
	Clock
)
