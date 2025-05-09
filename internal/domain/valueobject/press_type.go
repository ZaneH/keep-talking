package valueobject

type PressType string

const (
	PressTypeTap     PressType = "Tap"
	PressTypeHold    PressType = "Hold"
	PressTypeRelease PressType = "Release"
)
