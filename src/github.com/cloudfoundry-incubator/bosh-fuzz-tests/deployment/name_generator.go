package deployment

import (
	"math/rand"
)

type NameGenerator interface {
	Generate(length int) string
}

type nameGenerator struct {
	seed int64
}

func NewNameGenerator() NameGenerator {
	return &nameGenerator{}
}

func NewSeededNameGenerator(seed int64) NameGenerator {
	return &nameGenerator{seed: seed}
}

func (n *nameGenerator) Generate(length int) string {
	if n.seed != 0 {
		rand.Seed(n.seed)
	}

	b := make([]rune, length)
	b[0] = firstCharRunes[rand.Intn(len(firstCharRunes))]

	for i := 1; i < len(b); i++ {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var firstCharRunes = []rune("abcdefghijklmnopqrstuvwxyz")
