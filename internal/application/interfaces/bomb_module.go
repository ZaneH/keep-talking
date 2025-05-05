package module

type BombModule interface {
	ID() string
	IsSolved() bool
	ApplyInput(input any) error
}
