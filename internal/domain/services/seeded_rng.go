package services

import (
	"hash/fnv"
	"math/rand/v2"
)

type SeededRNG struct {
	seed uint64
	rng  *rand.Rand
}

func NewSeededRNG(seed uint64) *SeededRNG {
	return &SeededRNG{
		seed: seed,
		rng:  rand.New(rand.NewPCG(seed, seed)),
	}
}

func NewSeededRNGFromString(s string) *SeededRNG {
	seed := stringToSeed(s)
	return NewSeededRNG(seed)
}

func stringToSeed(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return uint64(h.Sum64())
}

func (s *SeededRNG) GetSeed() uint64 {
	return s.seed
}

func (s *SeededRNG) GetIntInRange(min, max int) int {
	return min + s.rng.IntN(max-min+1)
}

func (s *SeededRNG) Shuffle(length int, swapFunc func(i, j int)) {
	s.rng.Shuffle(length, swapFunc)
}

func (s *SeededRNG) Float32(min, max float32) float32 {
	return min + s.rng.Float32()*(max-min)
}
