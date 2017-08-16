package analyzer

import (
	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type stemcellComparator struct {
	logger boshlog.Logger
}

func NewStemcellComparator(logger boshlog.Logger) Comparator {
	return &stemcellComparator{
		logger: logger,
	}
}

func (s *stemcellComparator) Compare(previousInputs []bftinput.Input, currentInput bftinput.Input) []bftexpectation.Expectation {
	mostRecentInput := previousInputs[len(previousInputs)-1]
	expectations := []bftexpectation.Expectation{}
	for _, instanceGroup := range currentInput.InstanceGroups {
		if s.instanceGroupStemcellChanged(instanceGroup, currentInput, mostRecentInput) {
			expectations = append(expectations, bftexpectation.NewExistingInstanceDebugLog("stemcell_changed?", instanceGroup.Name))
		}
	}

	return expectations
}

func (s *stemcellComparator) instanceGroupStemcellChanged(instanceGroup bftinput.InstanceGroup, currentInput bftinput.Input, mostRecentInput bftinput.Input) bool {
	prevInstanceGroup, found := mostRecentInput.FindInstanceGroupByName(instanceGroup.Name)
	if !found {
		return false
	}

	var currentStemcell bftinput.StemcellConfig
	if instanceGroup.Stemcell != "" {
		currentStemcell = s.findStemcellByAlias(instanceGroup.Stemcell, currentInput)
	} else {
		currentStemcell = s.findResourcePoolStemcell(instanceGroup.ResourcePool, currentInput)
	}

	if prevInstanceGroup.Stemcell != "" {
		prevStemcell := s.findStemcellByAlias(prevInstanceGroup.Stemcell, mostRecentInput)
		if prevStemcell.Version != currentStemcell.Version {
			s.logger.Debug("stemcell_comparator", "Stemcell versions don't match. Previous input: %#v, new input: %#v", mostRecentInput, currentInput)
			return true
		}
	} else {
		prevStemcell := s.findResourcePoolStemcell(prevInstanceGroup.ResourcePool, mostRecentInput)
		if prevStemcell.Version != currentStemcell.Version {
			s.logger.Debug("stemcell_comparator", "Stemcell versions don't match. Previous input: %#v, new input: %#v", mostRecentInput, currentInput)
			return true
		}
	}

	return false
}

func (s *stemcellComparator) findStemcellByAlias(alias string, input bftinput.Input) bftinput.StemcellConfig {
	for _, stemcell := range input.Stemcells {
		if stemcell.Alias == alias {
			return stemcell
		}
	}

	return bftinput.StemcellConfig{}
}

func (s *stemcellComparator) findResourcePoolStemcell(resourcePoolName string, input bftinput.Input) bftinput.StemcellConfig {
	for _, resourcePool := range input.CloudConfig.ResourcePools {
		if resourcePool.Name == resourcePoolName {
			return resourcePool.Stemcell
		}
	}

	return bftinput.StemcellConfig{}
}
