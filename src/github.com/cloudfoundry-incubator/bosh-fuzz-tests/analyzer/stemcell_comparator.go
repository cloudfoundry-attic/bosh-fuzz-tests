package analyzer

import (
	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type stemcellComparator struct {
	expectationFactory bftexpectation.Factory
}

func NewStemcellComparator(expectationFactory bftexpectation.Factory) Comparator {
	return &stemcellComparator{
		expectationFactory: expectationFactory,
	}
}

func (s *stemcellComparator) Compare(previousInput bftinput.Input, currentInput bftinput.Input) []bftexpectation.Expectation {
	expectations := []bftexpectation.Expectation{}
	for _, job := range currentInput.Jobs {
		if s.jobStemcellChanged(job, currentInput, previousInput) {
			expectations = append(expectations, s.expectationFactory.CreateDebugLog("stemcell_changed?"))
		}
	}

	return expectations
}

func (s *stemcellComparator) jobStemcellChanged(job bftinput.Job, previousInput bftinput.Input, currentInput bftinput.Input) bool {
	prevJob, found := s.findJobByName(job.Name, previousInput)
	if !found {
		return false
	}

	var currentStemcell bftinput.StemcellConfig
	if job.Stemcell != "" {
		currentStemcell = s.findVmTypeStemcellByAlias(job.Stemcell, currentInput)
	} else {
		currentStemcell = s.findResourcePoolStemcell(job.ResourcePool, currentInput)
	}

	if prevJob.Stemcell != "" {
		prevStemcell := s.findVmTypeStemcellByAlias(prevJob.Stemcell, previousInput)
		if prevStemcell.Version != currentStemcell.Version {
			return true
		}
	} else {
		prevStemcell := s.findResourcePoolStemcell(prevJob.ResourcePool, previousInput)
		if prevStemcell.Version != currentStemcell.Version {
			return true
		}
	}

	return false
}

func (s *stemcellComparator) findVmTypeStemcellByAlias(alias string, input bftinput.Input) bftinput.StemcellConfig {
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

func (s *stemcellComparator) findJobByName(jobName string, input bftinput.Input) (bftinput.Job, bool) {
	for _, job := range input.Jobs {
		if job.Name == jobName {
			return job, true
		}
	}
	return bftinput.Job{}, false
}
