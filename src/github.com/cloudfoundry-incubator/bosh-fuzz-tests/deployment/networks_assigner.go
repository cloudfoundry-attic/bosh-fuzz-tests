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
		networkPoolWithAzs := []NetworkConfig{}
		networkTypes := n.networks[rand.Intn(len(n.networks))]

		for _, networkType := range networkTypes {
			network := NetworkConfig{
				Name: n.nameGenerator.Generate(7),
				Type: networkType,
			}
			networkPoolWithAzs = append(networkPoolWithAzs, network)
		}

		networkPoolWithoutAzs := []NetworkConfig{}
		networkTypes = n.networks[rand.Intn(len(n.networks))]
		for _, networkType := range networkTypes {
			network := NetworkConfig{
				Name: n.nameGenerator.Generate(7),
				Type: networkType,
			}
			networkPoolWithoutAzs = append(networkPoolWithoutAzs, network)
		}

		// TODO: shuffle networks

		for k, network := range networkPoolWithAzs {
			if network.Type != "vip" {
				networkPoolWithAzs[k].Subnets = n.generateSubnets(inputs[i].CloudConfig.AvailabilityZones)

				for s, _ := range networkPoolWithAzs[k].Subnets {
					if network.Type == "manual" {
						ipPool := n.ipPoolProvider.NewIpPool()
						networkPoolWithAzs[k].Subnets[s].IpRange = ipPool.IpRange
						networkPoolWithAzs[k].Subnets[s].Gateway = ipPool.Gateway
						// subnet.Reserved = ipPool.Reserved
					}
				}
			}
		}

		for j, job := range inputs[i].Jobs {
			if job.AvailabilityZones == nil {
				inputs[i].Jobs[j].Networks = n.generateJobNetworks(networkPoolWithoutAzs, nil)

			} else {
				inputs[i].Jobs[j].Networks = n.generateJobNetworks(networkPoolWithAzs, job.AvailabilityZones)
			}
		}

		compilationNetworks := []NetworkConfig{}
		for _, network := range append(networkPoolWithAzs, networkPoolWithoutAzs...) {
			if len(network.Subnets) > 0 || network.Type == "vip" {
				inputs[i].CloudConfig.Networks = append(inputs[i].CloudConfig.Networks, network)

				if network.Type != "vip" {
					compilationNetworks = append(compilationNetworks, network)
				}
			}
		}

		compilationNetwork := compilationNetworks[rand.Intn(len(compilationNetworks))]
		inputs[i].CloudConfig.CompilationNetwork = compilationNetwork.Name
		azs := []string{}
		for _, s := range compilationNetwork.Subnets {
			azs = append(azs, s.AvailabilityZones...)
		}
		if len(azs) > 0 {
			inputs[i].CloudConfig.CompilationAvailabilityZone = azs[rand.Intn(len(azs))]
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
	}

	jobNetworks[rand.Intn(totalNumberOfJobNetworks)].DefaultDNSnGW = true

	return jobNetworks
}

func (n *networksAssigner) generateSubnets(azs []string) []SubnetConfig {
	subnets := []SubnetConfig{}

	placedAzs := NewPlacedAZs()
	for !placedAzs.AllPlaced(azs) {
		newAzs := n.randomCombination(azs)
		placedAzs.Place(newAzs)
		subnets = append(subnets, SubnetConfig{AvailabilityZones: newAzs})
	}

	return subnets
}

func (n *networksAssigner) randomCombination(items []string) []string {
	combination := []string{}
	totalNumberOfItems := rand.Intn(len(items)) + 1
	itemsToPick := rand.Perm(len(items))[:totalNumberOfItems]
	for _, i := range itemsToPick {
		combination = append(combination, items[i])
	}

	return combination
}

type PlacedAZs struct {
	azs map[string]bool
}

func NewPlacedAZs() *PlacedAZs {
	return &PlacedAZs{
		azs: map[string]bool{},
	}
}

func (a *PlacedAZs) Place(azs []string) {
	for _, az := range azs {
		a.azs[az] = true
	}
}

func (a *PlacedAZs) AllPlaced(azs []string) bool {
	for _, az := range azs {
		if !a.azs[az] {
			return false
		}
	}

	return true
}
