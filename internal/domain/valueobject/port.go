package valueobject

type Port string

const (
	PortDVID   Port = "DVI-D"
	PortRCA    Port = "Stereo RCA"
	PortPS2    Port = "PS/2"
	PortRJ45   Port = "RJ-45"
	PortSerial Port = "Serial"
)

var AVAILABLE_PORTS = [...]Port{
	PortDVID,
	PortRCA,
	PortPS2,
	PortRJ45,
	PortSerial,
}
