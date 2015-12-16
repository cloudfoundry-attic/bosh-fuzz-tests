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
	for _, job := range currentInput.Jobs {
		if s.jobStemcellChanged(job, currentInput, mostRecentInput) {
			expectations = append(expectations, bftexpectation.NewExistingInstanceDebugLog("stemcell_changed?", job.Name))
		}
	}

	return expectations
}

func (s *stemcellComparator) jobStemcellChanged(job bftinput.Job, currentInput bftinput.Input, mostRecentInput bftinput.Input) bool {
	prevJob, found := mostRecentInput.FindJobByName(job.Name)
	if !found {
		return false
	}

	var currentStemcell bftinput.StemcellConfig
	if job.Stemcell != "" {
		currentStemcell = s.findStemcellByAlias(job.Stemcell, currentInput)
	} else {
		currentStemcell = s.findResourcePoolStemcell(job.ResourcePool, currentInput)
	}

	if prevJob.Stemcell != "" {
		prevStemcell := s.findStemcellByAlias(prevJob.Stemcell, mostRecentInput)
		if prevStemcell.Version != currentStemcell.Version {
			s.logger.Debug("stemcell_comparator", "Stemcell versions don't match. Previous input: %#v, new input: %#v", mostRecentInput, currentInput)
			return true
		}
	} else {
		prevStemcell := s.findResourcePoolStemcell(prevJob.ResourcePool, mostRecentInput)
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
