package deployment

import (
	"math/rand"
)

type NetworksAssigner interface {
	Assign(inputs []Input)
}

type networksAssigner struct {
	networks       [][]string
	nameGenerator  NameGenerator
	ipPoolProvider IpPoolProvider
	seed           int64
}

func NewNetworksAssigner(networks [][]string, nameGenerator NameGenerator) NetworksAssigner {
	return &networksAssigner{
		networks:       networks,
		nameGenerator:  nameGenerator,
		ipPoolProvider: NewIpPoolProvider(),
	}
}

func NewSeededNetworksAssigner(networks [][]string, nameGenerator NameGenerator, seed int64) NetworksAssigner {
	return &networksAssigner{
		networks:       networks,
		nameGenerator:  nameGenerator,
		ipPoolProvider: NewIpPoolProvider(),
		seed:           seed,
	}
}

func (n *networksAssigner) Assign(inputs []Input) {
	if n.seed != 0 {
		rand.Seed(n.seed)
	}

	for i, _ := range inputs {
		networkPool := []NetworkConfig{}
		networkTypes := n.networks[rand.Intn(len(n.networks))]
		for _, networkType := range networkTypes {
			networkName := n.nameGenerator.Generate(7)
			networkPool = append(networkPool, NetworkConfig{
				Name: networkName,
				Type: networkType,
			})
		}
		// TODO: shuffle networks

		for j, job := range inputs[i].Jobs {
			inputs[i].Jobs[j].Networks = n.generateJobNetworks(networkPool, job.AvailabilityZones)
		}

		for _, network := range networkPool {
			if len(network.Subnets) > 0 || network.Type == "vip" {
				inputs[i].CloudConfig.Networks = append(inputs[i].CloudConfig.Networks, network)
			}
		}
	}
}

func (n *networksAssigner) generateJobNetworks(networkPool []NetworkConfig, azs []string) []JobNetworkConfig {
	jobNetworks := []JobNetworkConfig{}

	totalNumberOfJobNetworks := rand.Intn(len(networkPool)) + 1
	networksToPick := rand.Perm(len(networkPool))[:totalNumberOfJobNetworks]
	for _, k := range networksToPick {
		network := networkPool[k]
		jobNetworks = append(jobNetworks, JobNetworkConfig{Name: network.Name})

		if network.Type != "vip" {
			subnet := SubnetConfig{AvailabilityZones: azs}
			ipPool := n.ipPoolProvider.NewIpPool()

			subnet.IpRange = ipPool.IpRange
			subnet.Gateway = ipPool.Gateway

			networkPool[k].Subnets = append(networkPool[k].Subnets, subnet)
		}
		// TODO: handle nil azs
		// TODO: reuse same subnet with all azs
	}

	jobNetworks[rand.Intn(totalNumberOfJobNetworks)].DefaultDNSnGW = true

	return jobNetworks
}
