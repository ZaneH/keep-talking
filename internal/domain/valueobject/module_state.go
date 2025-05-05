package valueobject

import "github.com/google/uuid"

type ModuleType int32

const (
	SIMPLE_WIRES ModuleType = iota
	PASSWORD
)

type ModuleState struct {
	ModuleId       uuid.UUID
	ModuleType     ModuleType
	ModulePosition ModulePosition
	Solved         bool
}
