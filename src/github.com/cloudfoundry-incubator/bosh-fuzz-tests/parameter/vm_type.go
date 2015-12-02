package parameter

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type vmType struct {
	definition    string
	nameGenerator bftnamegen.NameGenerator
}

func NewVmType(definition string, nameGenerator bftnamegen.NameGenerator) Parameter {
	return &vmType{
		definition:    definition,
		nameGenerator: nameGenerator,
	}
}

func (s *vmType) Apply(input *bftinput.Input) *bftinput.Input {
	for j, _ := range input.Jobs {
		if s.definition == "vm_type" {
			input.Jobs[j].VmType = s.nameGenerator.Generate(10)
			input.CloudConfig.VmTypes = append(
				input.CloudConfig.VmTypes,
				bftinput.VmTypeConfig{Name: input.Jobs[j].VmType},
			)
		} else if s.definition == "resource_pool" {
			input.Jobs[j].ResourcePool = s.nameGenerator.Generate(10)
			input.CloudConfig.ResourcePools = append(
				input.CloudConfig.ResourcePools,
				bftinput.ResourcePoolConfig{
					Name: input.Jobs[j].ResourcePool,
				},
			)
		}
	}

	return input
}
