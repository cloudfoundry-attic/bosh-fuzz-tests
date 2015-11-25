package deployment

import (
	"math/rand"
)

type NetworksAssigner interface {
	Assign(inputs []Input)
}

type networksAssigner struct {
	networks        [][]string
	nameGenerator   NameGenerator
	ipPoolProvider  IpPoolProvider
	staticIpDecider Decider
	seed            int64
}

func NewNetworksAssigner(networks [][]string, nameGenerator NameGenerator, ipPoolProvider IpPoolProvider, staticIpDecider Decider) NetworksAssigner {
	return &networksAssigner{
		networks:        networks,
		nameGenerator:   nameGenerator,
		ipPoolProvider:  ipPoolProvider,
		staticIpDecider: staticIpDecider,
	}
}

func NewSeededNetworksAssigner(networks [][]string, nameGenerator NameGenerator, ipPoolProvider IpPoolProvider, staticIpDecider Decider, seed int64) NetworksAssigner {
	return &networksAssigner{
		networks:        networks,
		nameGenerator:   nameGenerator,
		ipPoolProvider:  ipPoolProvider,
		staticIpDecider: staticIpDecider,
		seed:            seed,
	}
}

func (n *networksAssigner) Assign(inputs []Input) {
	if n.seed != 0 {
		rand.Seed(n.seed)
	}

	// 1. Generate Networks with/without AZs (network with types)
	// 2. Assign networks to each job (network with AZs) (job with network name)
	// 3. Generate IP specs for each network (network with IP specs)
	// 4. Aggregate instances to compute static IPs (network with static IPs) (job with static I)

	for i, _ := range inputs {
		networkPoolWithAzs := []NetworkConfig{}
		var networkTypes []string

		if len(inputs[i].CloudConfig.AvailabilityZones) > 0 {
			networkTypes = n.networks[rand.Intn(len(n.networks))]

			for _, networkType := range networkTypes {
				network := NetworkConfig{
					Name: n.nameGenerator.Generate(7),
					Type: networkType,
				}
				networkPoolWithAzs = append(networkPoolWithAzs, network)
			}

			for k, network := range networkPoolWithAzs {
				if network.Type != "vip" {
					networkPoolWithAzs[k].Subnets = n.generateSubnets(inputs[i].CloudConfig.AvailabilityZones)
				}
			}
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

		for k, network := range networkPoolWithoutAzs {
			if network.Type != "vip" {
				networkPoolWithoutAzs[k].Subnets = n.generateSubnetsWithoutAzs()
			}
		}

		for j, job := range inputs[i].Jobs {
			if job.AvailabilityZones == nil {
				inputs[i].Jobs[j].Networks = n.generateJobNetworks(networkPoolWithoutAzs, nil)
			} else {
				inputs[i].Jobs[j].Networks = n.generateJobNetworks(networkPoolWithAzs, job.AvailabilityZones)
			}
		}

		allNetworks := append(networkPoolWithAzs, networkPoolWithoutAzs...)
		n.assignStaticIps(allNetworks, inputs[i].Jobs)

		nonVipNetworks := []NetworkConfig{}

		for _, network := range allNetworks {
			inputs[i].CloudConfig.Networks = append(inputs[i].CloudConfig.Networks, network)

			if network.Type != "vip" {
				nonVipNetworks = append(nonVipNetworks, network)
			}
		}

		compilationNetwork := nonVipNetworks[rand.Intn(len(nonVipNetworks))]
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

func (n *networksAssigner) generateSubnetsWithoutAzs() []SubnetConfig {
	subnets := []SubnetConfig{}
	numberOfSubnets := rand.Intn(3) + 1 // up to 3

	for i := 0; i < numberOfSubnets; i++ {
		subnets = append(subnets, SubnetConfig{})
	}

	return subnets
}

type JobsOnNetwork struct {
	Jobs           []Job
	TotalInstances int
}

func (n *networksAssigner) aggregateNetworkJobs(jobs []Job) map[string]JobsOnNetwork {
	jobsOnNetworks := map[string]JobsOnNetwork{}

	for _, job := range jobs {
		for _, jobNetwork := range job.Networks {
			jobsOnNetworks[jobNetwork.Name] = JobsOnNetwork{
				Jobs:           append(jobsOnNetworks[jobNetwork.Name].Jobs, job),
				TotalInstances: jobsOnNetworks[jobNetwork.Name].TotalInstances + job.Instances,
			}
		}
	}

	return jobsOnNetworks
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

func (n *networksAssigner) assignStaticIps(networks []NetworkConfig, jobs []Job) {
	jobsOnNetworks := n.aggregateNetworkJobs(jobs)

	vipIpPool := n.ipPoolProvider.NewIpPool(254)

	for k, network := range networks {
		jobsOnNetwork := jobsOnNetworks[network.Name]

		if network.Type == "manual" {
			for s, _ := range network.Subnets {
				ipPool := n.ipPoolProvider.NewIpPool(jobsOnNetwork.TotalInstances)
				networks[k].Subnets[s].IpPool = ipPool
			}

			for _, job := range jobsOnNetwork.Jobs {
				if n.staticIpDecider.IsYes() { // use static IPs
					jobNetwork, found := n.findJobNetworkWithName(network.Name, job.Networks)
					if found {
						for ji := 0; ji < job.Instances; ji++ {
							s := rand.Intn(len(network.Subnets))
							staticIp, _ := networks[k].Subnets[s].IpPool.NextStaticIp()
							jobNetwork.StaticIps = append(jobNetwork.StaticIps, staticIp)
						}
					}
				}
			}
		} else if network.Type == "vip" {
			for _, job := range jobsOnNetwork.Jobs {
				jobNetwork, found := n.findJobNetworkWithName(network.Name, job.Networks)
				if found {
					for ji := 0; ji < job.Instances; ji++ {
						staticIp, _ := vipIpPool.NextStaticIp()
						jobNetwork.StaticIps = append(jobNetwork.StaticIps, staticIp)
					}
				}
			}
		}
	}
}

func (n *networksAssigner) findJobNetworkWithName(networkName string, jobNetworks []JobNetworkConfig) (JobNetworkConfig, bool) {
	for j, jobNetwork := range jobNetworks {
		if jobNetwork.Name == networkName {
			return jobNetworks[j], true
		}
	}
	return JobNetworkConfig{}, false
}
