package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeAvailabilityZone struct {
}

func NewFakeAvailabilityZone() *FakeAvailabilityZone {
	return &FakeAvailabilityZone{}
}

func (s *FakeAvailabilityZone) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	azs := []string{"z1"}

	input.CloudConfig.AvailabilityZones = []bftinput.AvailabilityZone{
		{
			Name: "z1",
		},
	}
	for j, _ := range input.Jobs {
		input.Jobs[j].AvailabilityZones = azs
	}

	return input
}
