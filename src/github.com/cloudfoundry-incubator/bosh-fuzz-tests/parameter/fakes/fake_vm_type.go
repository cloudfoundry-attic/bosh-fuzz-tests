package fakes

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type FakeVmType struct {
	definition string
}

func NewFakeVmType(definition string) *FakeVmType {
	return &FakeVmType{
		definition: definition,
	}
}

func (s *FakeVmType) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	if s.definition == "vm_type" {
		vmType := bftinput.VmTypeConfig{Name: "fake-vm-type"}

		input.CloudConfig.VmTypes = []bftinput.VmTypeConfig{
			vmType,
		}
		for j, _ := range input.Jobs {
			input.Jobs[j].VmType = vmType.Name
		}
	} else if s.definition == "resource_pool" {
		resourcePool := bftinput.ResourcePoolConfig{
			Name: "fake-resource-pool",
			Stemcell: bftinput.StemcellConfig{
				Name:    "foo-stemcell",
				Version: "1",
			},
		}

		input.CloudConfig.ResourcePools = []bftinput.ResourcePoolConfig{
			resourcePool,
		}
		for j, _ := range input.Jobs {
			input.Jobs[j].ResourcePool = resourcePool.Name
		}
	}

	return input
}
