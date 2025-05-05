package common

import "github.com/google/uuid"

type InputResult struct {
	ModuleId uuid.UUID
	Success  bool
}
