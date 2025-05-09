package valueobject

type Indicator struct {
	Label string
	Lit   bool
}

var AVAILABLE_INDICATOR_LABELS = []string{
	"OCR",
	"NSA",
	"CDN",
	"ZIP",
	"KEK",
	"KPR",
	"OSR",
	"BPM",
	"FRK",
	"JFK",
	"MTV",
	"RKG",
}
