package actors

import "fmt"

type ActorError interface {
	Error() string
}

var (
	ErrInvalidModuleType    ActorError = fmt.Errorf("invalid module type")
	ErrInvalidModuleCommand ActorError = fmt.Errorf("invalid module command")
	ErrUnhandledMessageType ActorError = fmt.Errorf("unhandled message type")
)
