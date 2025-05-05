package entities

type WireModule struct {
	Id string
}

func (w *WireModule) CutWire(wireIndex int) (string, error) {
	// Logic to cut the wire
	return "Wire cut successfully", nil
}
