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
	for _, instanceGroup := range currentInput.InstanceGroups {
		if n.nothingChanged(instanceGroup, currentInput, previousInputs) && n.isNotErrand(instanceGroup) {
			expectations = append(expectations, bftexpectation.NewDebugLog(fmt.Sprintf("No instances to update for '%s'", instanceGroup.Name)))
		}
	}

	return expectations
}

func (n *nothingChangedComparator) nothingChanged(instanceGroup bftinput.InstanceGroup, currentInput bftinput.Input, previousInputs []bftinput.Input) bool {
	mostRecentInput := previousInputs[len(previousInputs)-1]

	prevInstanceGroup, found := mostRecentInput.FindInstanceGroupByName(instanceGroup.Name)
	if !found {
		return false
	}

	if len(previousInputs) > 1 {
		inputBeforePrevious := previousInputs[len(previousInputs)-2]
		instanceGroupBeforePrevious, found := inputBeforePrevious.FindInstanceGroupByName(instanceGroup.Name)
		if found && instanceGroupBeforePrevious.HasPersistentDisk() && !prevInstanceGroup.HasPersistentDisk() {
			return false
		}

		for _, migratedFromConfig := range prevInstanceGroup.MigratedFrom {
			instanceGroupBeforePrevious, found := inputBeforePrevious.FindInstanceGroupByName(migratedFromConfig.Name)
			if found && instanceGroupBeforePrevious.HasPersistentDisk() && !prevInstanceGroup.HasPersistentDisk() {
				return false
			}
		}
	}

	if !prevInstanceGroup.IsEqual(instanceGroup) {
		return false
	}

	for _, azName := range instanceGroup.AvailabilityZones {
		currentAz, _ := currentInput.FindAzByName(azName)
		prevAz, _ := mostRecentInput.FindAzByName(azName)
		if !currentAz.IsEqual(prevAz) {
			return false
		}
	}

	if instanceGroup.PersistentDiskPool != "" {
		currentPersistentDiskPool, _ := currentInput.FindDiskPoolByName(instanceGroup.PersistentDiskPool)
		prevPersistentDiskPool, _ := mostRecentInput.FindDiskPoolByName(instanceGroup.PersistentDiskPool)
		if !currentPersistentDiskPool.IsEqual(prevPersistentDiskPool) {
			return false
		}
	}

	if instanceGroup.PersistentDiskType != "" {
		currentPersistentDiskType, _ := currentInput.FindDiskTypeByName(instanceGroup.PersistentDiskType)
		prevPersistentDiskType, _ := mostRecentInput.FindDiskTypeByName(instanceGroup.PersistentDiskType)
		if !currentPersistentDiskType.IsEqual(prevPersistentDiskType) {
			return false
		}
	}

	if instanceGroup.ResourcePool != "" {
		currentResourcePool, _ := currentInput.FindResourcePoolByName(instanceGroup.ResourcePool)
		prevResourcePool, _ := mostRecentInput.FindResourcePoolByName(instanceGroup.ResourcePool)
		if !currentResourcePool.IsEqual(prevResourcePool) {
			return false
		}
	}

	if instanceGroup.VmType != "" {
		currentVmType, _ := currentInput.FindVmTypeByName(instanceGroup.VmType)
		prevVmType, _ := mostRecentInput.FindVmTypeByName(instanceGroup.VmType)
		if !currentVmType.IsEqual(prevVmType) {
			return false
		}
	}

	if instanceGroup.Stemcell != "" {
		currentStemcell, _ := currentInput.FindStemcellByName(instanceGroup.Stemcell)
		prevStemcell, _ := mostRecentInput.FindStemcellByName(instanceGroup.Stemcell)
		if !currentStemcell.IsEqual(prevStemcell) {
			return false
		}
	}

	for _, instanceGroupNetwork := range instanceGroup.Networks {
		currentNetwork, _ := currentInput.FindNetworkByName(instanceGroupNetwork.Name)
		prevNetwork, _ := mostRecentInput.FindNetworkByName(instanceGroupNetwork.Name)
		if !currentNetwork.IsEqual(prevNetwork) {
			return false
		}
	}

	return true
}

func (n *nothingChangedComparator) isNotErrand(instanceGroup bftinput.InstanceGroup) bool {
	return instanceGroup.Lifecycle != "errand"
}
