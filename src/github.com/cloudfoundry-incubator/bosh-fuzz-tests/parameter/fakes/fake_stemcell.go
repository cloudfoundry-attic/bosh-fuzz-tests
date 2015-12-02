package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeStemcell struct {
}

func NewFakeStemcell() *FakeStemcell {
	return &FakeStemcell{}
}

func (s *FakeStemcell) Apply(input *bftinput.Input) *bftinput.Input {
	input.Stemcells = []bftinput.StemcellConfig{
		{Name: "fake-stemcell"},
	}

	return input
}
