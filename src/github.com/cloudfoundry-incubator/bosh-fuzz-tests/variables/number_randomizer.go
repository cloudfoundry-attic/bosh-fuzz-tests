package variables

import (
	"math/rand"
)

type NumberRandomizer interface {
	Intn(n int) int
}

type DefaultNumberRandomizer struct {
}

func (d DefaultNumberRandomizer) Intn(n int) int {
	return rand.Intn(n)
}
