package decider

import (
	"math/rand"
)

type Decider interface {
	IsYes() bool
}

type randomDecider struct {
}

func NewRandomDecider() Decider {
	return &randomDecider{}
}

func (r *randomDecider) IsYes() bool {
	return rand.Intn(2) == 1
}
