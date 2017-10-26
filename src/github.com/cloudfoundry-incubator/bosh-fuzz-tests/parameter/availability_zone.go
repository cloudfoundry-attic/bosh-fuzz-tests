package parameter

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type availabilityZone struct {
	azs [][]string
}

func NewAvailabilityZone(azs [][]string) Parameter {
	return &availabilityZone{
		azs: azs,
	}
}

func (a *availabilityZone) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	azs := map[string]bool{}
	input.CloudConfig.AvailabilityZones = nil

	azsList := [][]string{}
	if input.HasMigratedInstances() {
		for _, a := range a.azs {
			if len(a) == 0 {
				continue
			}
			azsList = append(azsList, a)
		}
	} else {
		azsList = a.azs
	}

	for j := range input.InstanceGroups {
		input.InstanceGroups[j].AvailabilityZones = azsList[rand.Intn(len(azsList))]

		for _, name := range input.InstanceGroups[j].AvailabilityZones {
			if azs[name] != true {
				input.CloudConfig.AvailabilityZones = append(
					input.CloudConfig.AvailabilityZones,
					bftinput.AvailabilityZone{
						Name: name,
					})
			}
			azs[name] = true
		}
	}

	return input
}
