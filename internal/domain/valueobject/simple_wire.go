package valueobject

type SimpleWire struct {
	WireColor Color
	IsCut     bool
	Index     int
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
