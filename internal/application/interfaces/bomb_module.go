package interfaces

import "github.com/google/uuid"

type BombModule interface {
	Id() uuid.UUID
	IsSolved() bool
	ApplyInput(input any) error
}
