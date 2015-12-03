package parameter

import (
	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type vmType struct {
	definition    string
	nameGenerator bftnamegen.NameGenerator
	reuseDecider  bftdecider.Decider
}

func NewVmType(definition string, nameGenerator bftnamegen.NameGenerator, reuseDecider bftdecider.Decider) Parameter {
	return &vmType{
		definition:    definition,
		nameGenerator: nameGenerator,
		reuseDecider:  reuseDecider,
	}
}

func (s *vmType) Apply(input bftinput.Input) bftinput.Input {
	input.CloudConfig.VmTypes = nil
	input.CloudConfig.ResourcePools = nil

	for j, _ := range input.Jobs {

		if s.definition == "vm_type" {
			input.Jobs[j].ResourcePool = ""

			if !s.reuseDecider.IsYes() || input.Jobs[j].VmType == "" {
				input.Jobs[j].VmType = s.nameGenerator.Generate(10)
			}

			input.CloudConfig.VmTypes = append(
				input.CloudConfig.VmTypes,
				bftinput.VmTypeConfig{Name: input.Jobs[j].VmType},
			)
		} else if s.definition == "resource_pool" {
			input.Jobs[j].VmType = ""

			if !s.reuseDecider.IsYes() || input.Jobs[j].ResourcePool == "" {
				input.Jobs[j].ResourcePool = s.nameGenerator.Generate(10)
			}

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
