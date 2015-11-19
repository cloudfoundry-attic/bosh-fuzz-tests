package deployment

import (
	"math/rand"
)

type NetworksAssigner interface {
	Assign(inputs []Input)
}

type networksAssigner struct {
	networks      [][]string
	nameGenerator NameGenerator
	seed          int64
}

func NewNetworksAssigner(networks [][]string, nameGenerator NameGenerator) NetworksAssigner {
	return &networksAssigner{networks: networks, nameGenerator: nameGenerator}
}

func NewSeededNetworksAssigner(networks [][]string, nameGenerator NameGenerator, seed int64) NetworksAssigner {
	return &networksAssigner{networks: networks, nameGenerator: nameGenerator, seed: seed}
}

func (n *networksAssigner) Assign(inputs []Input) {
	if n.seed != 0 {
		rand.Seed(n.seed)
	}

	for i, _ := range inputs {
		networkPool := []NetworkConfig{}
		networkTypes := n.networks[rand.Intn(len(n.networks))]
		for _, networkType := range networkTypes {
			networkName := n.nameGenerator.Generate(10)
			networkPool = append(networkPool, NetworkConfig{
				Name: networkName,
				Type: networkType,
			})
		}

		for j, job := range inputs[i].Jobs {
			totalNumberOfJobNetworks := rand.Intn(len(networkPool)) + 1
			networksToPick := rand.Perm(len(networkPool))[:totalNumberOfJobNetworks]
			for _, k := range networksToPick {
				network := networkPool[k]
				inputs[i].Jobs[j].Networks = append(inputs[i].Jobs[j].Networks, JobNetworkConfig{Name: network.Name})
				if job.AvailabilityZones != nil {
					subnet := SubnetConfig{AvailabilityZones: job.AvailabilityZones}
					networkPool[k].Subnets = append(networkPool[k].Subnets, subnet)
				}
				// TODO: handle nil azs
				// TODO: deduplicate to avoid overlapping subnets
			}

		}

		for _, network := range networkPool {
			inputs[i].CloudConfig.Networks = append(inputs[i].CloudConfig.Networks, network)
		}

		// Workaround: add default network to make compilation work
		defaultNetwork := NetworkConfig{
			Name: "default",
			Subnets: []SubnetConfig{
				SubnetConfig{inputs[i].CloudConfig.AvailabilityZones},
			},
		}
		inputs[i].CloudConfig.Networks = append(inputs[i].CloudConfig.Networks, defaultNetwork)
	}
}
