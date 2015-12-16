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
	properties := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	for i, _ := range input.CloudConfig.AvailabilityZones {
		input.CloudConfig.AvailabilityZones[i].CloudProperties = properties
	}

	for i, _ := range input.CloudConfig.VmTypes {
		input.CloudConfig.VmTypes[i].CloudProperties = properties
	}

	return input
}
