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

func (n *nothingChangedComparator) Compare(previousInputs []bftinput.Input, currentInput bftinput.Input) []bftexpectation.Expectation {
	expectations := []bftexpectation.Expectation{}
	for _, job := range currentInput.Jobs {
		if n.nothingChanged(job, currentInput, previousInputs) {
			expectations = append(expectations, bftexpectation.NewDebugLog(fmt.Sprintf("No instances to update for '%s'", job.Name)))
		}
	}

	return expectations
}

func (n *nothingChangedComparator) nothingChanged(job bftinput.Job, currentInput bftinput.Input, previousInputs []bftinput.Input) bool {
	mostRecentInput := previousInputs[len(previousInputs)-1]

	prevJob, found := mostRecentInput.FindJobByName(job.Name)
	if !found {
		return false
	}

	if len(previousInputs) > 1 {
		inputBeforePrevious := previousInputs[len(previousInputs)-2]
		jobBeforePrevious, found := inputBeforePrevious.FindJobByName(job.Name)
		if found && jobBeforePrevious.HasPersistentDisk() && !prevJob.HasPersistentDisk() {
			return false
		}

		for _, migratedFromConfig := range prevJob.MigratedFrom {
			jobBeforePrevious, found := inputBeforePrevious.FindJobByName(migratedFromConfig.Name)
			if found && jobBeforePrevious.HasPersistentDisk() && !prevJob.HasPersistentDisk() {
				return false
			}
		}
	}

	if !prevJob.IsEqual(job) {
		return false
	}

	for _, azName := range job.AvailabilityZones {
		currentAz, _ := currentInput.FindAzByName(azName)
		prevAz, _ := mostRecentInput.FindAzByName(azName)
		if !currentAz.IsEqual(prevAz) {
			return false
		}
	}

	if job.PersistentDiskPool != "" {
		currentPersistentDiskPool, _ := currentInput.FindDiskPoolByName(job.PersistentDiskPool)
		prevPersistentDiskPool, _ := mostRecentInput.FindDiskPoolByName(job.PersistentDiskPool)
		if !currentPersistentDiskPool.IsEqual(prevPersistentDiskPool) {
			return false
		}
	}

	if job.PersistentDiskType != "" {
		currentPersistentDiskType, _ := currentInput.FindDiskTypeByName(job.PersistentDiskType)
		prevPersistentDiskType, _ := mostRecentInput.FindDiskTypeByName(job.PersistentDiskType)
		if !currentPersistentDiskType.IsEqual(prevPersistentDiskType) {
			return false
		}
	}

	if job.ResourcePool != "" {
		currentResourcePool, _ := currentInput.FindResourcePoolByName(job.ResourcePool)
		prevResourcePool, _ := mostRecentInput.FindResourcePoolByName(job.ResourcePool)
		if !currentResourcePool.IsEqual(prevResourcePool) {
			return false
		}
	}

	if job.VmType != "" {
		currentVmType, _ := currentInput.FindVmTypeByName(job.VmType)
		prevVmType, _ := mostRecentInput.FindVmTypeByName(job.VmType)
		if !currentVmType.IsEqual(prevVmType) {
			return false
		}
	}

	if job.Stemcell != "" {
		currentStemcell, _ := currentInput.FindStemcellByName(job.Stemcell)
		prevStemcell, _ := mostRecentInput.FindStemcellByName(job.Stemcell)
		if !currentStemcell.IsEqual(prevStemcell) {
			return false
		}
	}

	for _, jobNetwork := range job.Networks {
		currentNetwork, _ := currentInput.FindNetworkByName(jobNetwork.Name)
		prevNetwork, _ := mostRecentInput.FindNetworkByName(jobNetwork.Name)
		if !currentNetwork.IsEqual(prevNetwork) {
			return false
		}
	}

	return true
}
