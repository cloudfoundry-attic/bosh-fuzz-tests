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

func (a *availabilityZone) Apply(input bftinput.Input) bftinput.Input {
	azs := map[string]bool{}
	input.CloudConfig.AvailabilityZones = nil

	for j := range input.Jobs {
		input.Jobs[j].AvailabilityZones = a.azs[rand.Intn(len(a.azs))]

		for _, name := range input.Jobs[j].AvailabilityZones {
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
