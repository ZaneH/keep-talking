package valueobject

type SimpleWire struct {
	WireColor Color
	IsCut     bool
	Position  int
}

var SimpleWireColors = [...]Color{
	Red,
	Blue,
	White,
	Black,
	Yellow,
	Green,
	Orange,
	Pink,
}
