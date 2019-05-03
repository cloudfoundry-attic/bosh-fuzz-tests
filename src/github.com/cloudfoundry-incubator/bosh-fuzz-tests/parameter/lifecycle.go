package parameter

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type Lifecycle struct{}

func lifecycles(instanceGroup, previousInstanceGroup bftinput.InstanceGroup) []string {
	lifecycles := []string{"service"}

	if instanceGroup.PersistentDiskType == "" && previousInstanceGroup.PersistentDiskType == "" &&
		instanceGroup.PersistentDiskSize == 0 && previousInstanceGroup.PersistentDiskSize == 0 {
		lifecycles = append(lifecycles, "errand")
	}

	return lifecycles
}

func NewLifecycle() Parameter {
	return Lifecycle{}
}

func (l Lifecycle) Apply(input, previousInput bftinput.Input) bftinput.Input {
	newInput := bftinput.Input{
		DirectorUUID: input.DirectorUUID,
		Update:       input.Update,
		CloudConfig:  input.CloudConfig,
		Stemcells:    input.Stemcells,
		Variables:    input.Variables,
	}

	for _, group := range input.InstanceGroups {
		newInput.InstanceGroups = append(newInput.InstanceGroups, group)
	}

	for i, instanceGroup := range newInput.InstanceGroups {
		previousInstanceGroup := findInstanceGroupFromInput(instanceGroup, previousInput)

		cycles := lifecycles(instanceGroup, previousInstanceGroup)
		newInput.InstanceGroups[i].Lifecycle = cycles[rand.Intn(len(cycles))]
	}

	return newInput
}

func findInstanceGroupFromInput(desiredInstanceGroup bftinput.InstanceGroup, input bftinput.Input) bftinput.InstanceGroup {
	for _, instanceGroup := range input.InstanceGroups {
		if instanceGroup.Name == desiredInstanceGroup.Name {
			return instanceGroup
		} else {
			for _, migratedInstanceGroup := range desiredInstanceGroup.MigratedFrom {
				if instanceGroup.Name == migratedInstanceGroup.Name {
					return instanceGroup
				}
			}
		}
	}

	return bftinput.InstanceGroup{}
}
