package module

import "errors"

type WireColor string

const (
	Red    WireColor = "red"
	Blue   WireColor = "blue"
	Green  WireColor = "green"
	Yellow WireColor = "yellow"
)

type WiresModule struct {
	Id      string
	Face    Face
	Wires   []WireColor
	Correct int // Index of the current wire to cut
	Solved  bool
}

func NewWiresModule(id string, face Face, wires []WireColor, correct int) *WiresModule {
	return &WiresModule{
		Id:      id,
		Face:    face,
		Wires:   wires,
		Correct: correct,
		Solved:  false,
	}
}

func (m *WiresModule) CutWire(index int) error {
	if index < 0 || index >= len(m.Wires) {
		return errors.New("invalid wire index")
	}

	if m.Solved {
		return nil
	}

	if index == m.Correct {
		m.Solved = true
		return nil
	}

	return errors.New("wrong wire")
}
