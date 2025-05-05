package command

import (
	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type SubmitPasswordCommand struct {
	SessionId      uuid.UUID
	ModulePosition valueobject.ModulePosition
	Password       string
}

type SubmitPasswordResult struct {
	Result *common.InputResult
}
