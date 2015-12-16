package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeCloudProperties struct {
}

func NewFakeCloudProperties() *FakeCloudProperties {
	return &FakeCloudProperties{}
}

func (s *FakeCloudProperties) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	for i, _ := range input.CloudConfig.AvailabilityZones {
		input.CloudConfig.AvailabilityZones[i].CloudProperties = map[string]string{
			"foo": "bar",
			"baz": "qux",
		}
	}

	return input
}
