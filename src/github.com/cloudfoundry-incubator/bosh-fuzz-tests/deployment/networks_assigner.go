package deployment

import (
	"math/rand"
)

type NetworksAssigner interface {
	Assign(inputs []Input)
}

type networksAssigner struct {
	seed int64
}

func NewNetworksAssigner() NetworksAssigner {
	return &networksAssigner{}
}

func NewSeededNetworksAssigner(seed int64) NetworksAssigner {
	return &networksAssigner{seed: seed}
}

func (n *networksAssigner) Assign(inputs []Input) {
	if n.seed != 0 {
		rand.Seed(n.seed)
	}

	for i, _ := range inputs {
		for j, job := range inputs[i].Jobs {
			if job.AvailabilityZones == nil {
				inputs[i].Jobs[j].Network = "no-az"
			} else {
				inputs[i].Jobs[j].Network = "default"
			}
		}
		inputs[i].CloudConfig.Networks = []NetworkConfig{
			{
				Name:              "default",
				AvailabilityZones: inputs[i].CloudConfig.AvailabilityZones,
			},
			{
				Name: "no-az",
			},
		}
	}
}
