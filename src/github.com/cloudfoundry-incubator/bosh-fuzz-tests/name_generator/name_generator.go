package name_generator

import (
	"math/rand"
)

type NameGenerator interface {
	Generate(length int) string
}

type nameGenerator struct {
	absolute bool
}

func NewNameGenerator() NameGenerator {
	return &nameGenerator{}
}

func NewAbsoluteNameGenerator() NameGenerator {
	return &nameGenerator{
		absolute: true,
	}
}

func (n *nameGenerator) Generate(length int) string {
	b := make([]rune, length)
	b[0] = firstCharRunes[rand.Intn(len(firstCharRunes))]

	for i := 1; i < len(b); i++ {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	if n.absolute {
		return "/" + string(b)
	}
	return string(b)
}

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var firstCharRunes = []rune("abcdefghijklmnopqrstuvwxyz")
