package valueobject

type SimpleWire struct {
	WireColor Color
	IsCut     bool
	Position  int
}

var SimpleWireColors = [...]Color{
	Yellow,
	Red,
	Blue,
	Black,
	White,
}
