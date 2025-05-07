package valueobject

type BombFace int

const (
	Front BombFace = iota
	Back
)

type ModulePosition struct {
	Row    int
	Column int
	Face   BombFace
}
