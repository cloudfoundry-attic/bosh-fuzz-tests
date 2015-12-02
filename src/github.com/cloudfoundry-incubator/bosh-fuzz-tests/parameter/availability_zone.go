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

func (a *availabilityZone) Apply(input *bftinput.Input) *bftinput.Input {
	azs := map[string]bool{}

	for j := range input.Jobs {
		input.Jobs[j].AvailabilityZones = a.azs[rand.Intn(len(a.azs))]

		for _, az := range input.Jobs[j].AvailabilityZones {
			if azs[az] != true {
				input.CloudConfig.AvailabilityZones = append(input.CloudConfig.AvailabilityZones, az)
			}
			azs[az] = true
		}
	}

	return input
}
