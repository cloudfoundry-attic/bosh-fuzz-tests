package analyzer

import (
	"fmt"

	bftexpectation "github.com/cloudfoundry-incubator/bosh-fuzz-tests/expectation"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type nothingChangedComparator struct{}

func NewNothingChangedComparator() Comparator {
	return &nothingChangedComparator{}
}

func (n *nothingChangedComparator) Compare(previousInput bftinput.Input, currentInput bftinput.Input) []bftexpectation.Expectation {
	expectations := []bftexpectation.Expectation{}
	for _, job := range currentInput.Jobs {
		if n.nothingChanged(job, currentInput, previousInput) {
			expectations = append(expectations, bftexpectation.NewDebugLog(fmt.Sprintf("No instances to update for '%s'", job.Name)))
		}
	}

	return expectations
}

func (n *nothingChangedComparator) nothingChanged(job bftinput.Job, currentInput bftinput.Input, previousInput bftinput.Input) bool {
	prevJob, found := previousInput.FindJobByName(job.Name)
	if !found {
		return true
	}

	if !prevJob.IsEqual(job) {
		return false
	}

	for _, azName := range job.AvailabilityZones {
		currentAz, _ := currentInput.FindAzByName(azName)
		prevAz, _ := previousInput.FindAzByName(azName)
		if !currentAz.IsEqual(prevAz) {
			return false
		}
	}

	if job.PersistentDiskPool != "" {
		currentPersistentDiskPool, _ := currentInput.FindDiskPoolByName(job.PersistentDiskPool)
		prevPersistentDiskPool, _ := previousInput.FindDiskPoolByName(job.PersistentDiskPool)
		if !currentPersistentDiskPool.IsEqual(prevPersistentDiskPool) {
			return false
		}
	}

	if job.PersistentDiskType != "" {
		currentPersistentDiskType, _ := currentInput.FindDiskTypeByName(job.PersistentDiskType)
		prevPersistentDiskType, _ := previousInput.FindDiskTypeByName(job.PersistentDiskType)
		if !currentPersistentDiskType.IsEqual(prevPersistentDiskType) {
			return false
		}
	}

	for _, jobNetwork := range job.Networks {
		currentNetwork, _ := currentInput.FindNetworkByName(jobNetwork.Name)
		prevNetwork, _ := previousInput.FindNetworkByName(jobNetwork.Name)
		if !currentNetwork.IsEqual(prevNetwork) {
			return false
		}
	}

	return true
}
