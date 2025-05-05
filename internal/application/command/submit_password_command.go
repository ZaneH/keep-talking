package command

import (
	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/google/uuid"
)

type SubmitPasswordCommand struct {
	ModuleId uuid.UUID
	Password string
}

type SubmitPasswordResult struct {
	Result *common.InputResult
}
