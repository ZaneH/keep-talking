package ports

type RandomGenerator interface {
	GetSeed() uint64
	GetIntInRange(min, max int) int
	Shuffle(length int, swapFunc func(i, j int))
	Float32(min, max float32) float32
}
