package valueobject

type SimpleWire struct {
	WireColor Color
	IsCut     bool
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
