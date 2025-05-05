package valueobjects

type Face int

const (
	Front Face = iota
	Back
)

type ModulePosition struct {
	Row    int
	Column int
	Face   Face
}
