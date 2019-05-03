package parameter

import (
	"math/rand"

	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type vmType struct {
	nameGenerator bftnamegen.NameGenerator
	reuseDecider  bftdecider.Decider
	logger        boshlog.Logger
}

func NewVmType(
	nameGenerator bftnamegen.NameGenerator,
	reuseDecider bftdecider.Decider,
	logger boshlog.Logger,
) Parameter {
	return &vmType{
		nameGenerator: nameGenerator,
		reuseDecider:  reuseDecider,
		logger:        logger,
	}
}

func (s *vmType) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	input.CloudConfig.VmTypes = nil

	usedVmTypes := map[string]bool{}

	for j, _ := range input.InstanceGroups {
		if input.IsV2() {
			reuseFromOtherInstanceGroup := s.reuseDecider.IsYes()
			if reuseFromOtherInstanceGroup && j > 0 {
				previousInstanceGroup := input.InstanceGroups[rand.Intn(j)]
				input.InstanceGroups[j].VmType = previousInstanceGroup.VmType

			} else {
				reuseFromPreviousDeploy := s.reuseDecider.IsYes()
				if !reuseFromPreviousDeploy || input.InstanceGroups[j].VmType == "" {
					input.InstanceGroups[j].VmType = s.nameGenerator.Generate(10)
				}
			}

			if usedVmTypes[input.InstanceGroups[j].VmType] != true {
				input.CloudConfig.VmTypes = append(
					input.CloudConfig.VmTypes,
					bftinput.VmTypeConfig{Name: input.InstanceGroups[j].VmType},
				)
			}
			usedVmTypes[input.InstanceGroups[j].VmType] = true
		}
	}

	return input
}
