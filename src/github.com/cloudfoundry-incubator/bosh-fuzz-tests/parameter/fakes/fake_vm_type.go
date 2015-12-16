package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeVmType struct {
}

func NewFakeVmType() *FakeVmType {
	return &FakeVmType{}
}

func (s *FakeVmType) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	vmType := bftinput.VmTypeConfig{Name: "fake-vm-type"}

	input.CloudConfig.VmTypes = []bftinput.VmTypeConfig{
		vmType,
	}
	for j, _ := range input.Jobs {
		input.Jobs[j].VmType = vmType.Name
	}

	return input
}
