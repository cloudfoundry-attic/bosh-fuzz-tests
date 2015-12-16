package parameter

import (
	"math/rand"

	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type vmType struct {
	definition    string
	nameGenerator bftnamegen.NameGenerator
	reuseDecider  bftdecider.Decider
	logger        boshlog.Logger
}

func NewVmType(
	definition string,
	nameGenerator bftnamegen.NameGenerator,
	reuseDecider bftdecider.Decider,
	logger boshlog.Logger,
) Parameter {
	return &vmType{
		definition:    definition,
		nameGenerator: nameGenerator,
		reuseDecider:  reuseDecider,
		logger:        logger,
	}
}

func (s *vmType) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	input.CloudConfig.VmTypes = nil
	input.CloudConfig.ResourcePools = nil

	usedVmTypes := map[string]bool{}

	for j, _ := range input.Jobs {
		if s.definition == "vm_type" {
			input.Jobs[j].ResourcePool = ""

			reuseFromOtherJob := s.reuseDecider.IsYes()
			if reuseFromOtherJob && j > 0 {
				previousJob := input.Jobs[rand.Intn(j)]
				input.Jobs[j].VmType = previousJob.VmType

			} else {
				reuseFromPreviousDeploy := s.reuseDecider.IsYes()
				if !reuseFromPreviousDeploy || input.Jobs[j].VmType == "" {
					input.Jobs[j].VmType = s.nameGenerator.Generate(10)
				}
			}

			if usedVmTypes[input.Jobs[j].VmType] != true {
				input.CloudConfig.VmTypes = append(
					input.CloudConfig.VmTypes,
					bftinput.VmTypeConfig{Name: input.Jobs[j].VmType},
				)
			}
			usedVmTypes[input.Jobs[j].VmType] = true

		} else if s.definition == "resource_pool" {
			input.Jobs[j].VmType = ""

			reuseFromOtherJob := s.reuseDecider.IsYes()
			if reuseFromOtherJob && j > 0 {
				previousJob := input.Jobs[rand.Intn(j)]
				input.Jobs[j].ResourcePool = previousJob.ResourcePool

			} else {
				reuseFromPreviousDeploy := s.reuseDecider.IsYes()
				if !reuseFromPreviousDeploy || input.Jobs[j].ResourcePool == "" {
					input.Jobs[j].ResourcePool = s.nameGenerator.Generate(10)
				}
			}

			if usedVmTypes[input.Jobs[j].ResourcePool] != true {
				input.CloudConfig.ResourcePools = append(
					input.CloudConfig.ResourcePools,
					bftinput.ResourcePoolConfig{
						Name: input.Jobs[j].ResourcePool,
					},
				)
			}
			usedVmTypes[input.Jobs[j].ResourcePool] = true
		}
	}

	return input
}
