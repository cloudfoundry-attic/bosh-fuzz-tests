package deployment

import (
	"math/rand"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type NetworksAssigner interface {
	Assign(inputs []bftinput.Input)
}

type networksAssigner struct {
	networks        [][]string
	nameGenerator   bftnamegen.NameGenerator
	ipPoolProvider  IpPoolProvider
	staticIpDecider Decider
}

func NewNetworksAssigner(networks [][]string, nameGenerator bftnamegen.NameGenerator, ipPoolProvider IpPoolProvider, staticIpDecider Decider) NetworksAssigner {
	return &networksAssigner{
		networks:        networks,
		nameGenerator:   nameGenerator,
		ipPoolProvider:  ipPoolProvider,
		staticIpDecider: staticIpDecider,
	}
}

func (n *networksAssigner) Assign(inputs []bftinput.Input) {
	// 1. Generate Networks with/without AZs (network with types)
	// 2. Assign networks to each job (network with AZs) (job with network name)
	// 3. Generate IP specs for each network (network with IP specs)
	// 4. Aggregate instances to compute static IPs (network with static IPs) (job with static I)

	for i, _ := range inputs {
		n.ipPoolProvider.Reset()

		networkPoolWithAzs := []bftinput.NetworkConfig{}
		var networkTypes []string

		if len(inputs[i].CloudConfig.AvailabilityZones) > 0 {
			networkTypes = n.networks[rand.Intn(len(n.networks))]

			for _, networkType := range networkTypes {
				network := bftinput.NetworkConfig{
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

		networkPoolWithoutAzs := []bftinput.NetworkConfig{}
		networkTypes = n.networks[rand.Intn(len(n.networks))]
		for _, networkType := range networkTypes {
			network := bftinput.NetworkConfig{
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

		nonVipNetworks := []bftinput.NetworkConfig{}

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

func (n *networksAssigner) generateJobNetworks(networkPool []bftinput.NetworkConfig, azs []string) []bftinput.JobNetworkConfig {
	jobNetworks := []bftinput.JobNetworkConfig{}

	totalNumberOfJobNetworks := rand.Intn(len(networkPool)) + 1
	networksToPick := rand.Perm(len(networkPool))[:totalNumberOfJobNetworks]
	for _, k := range networksToPick {
		network := networkPool[k]
		jobNetworks = append(jobNetworks, bftinput.JobNetworkConfig{Name: network.Name})
	}

	jobNetworks[rand.Intn(totalNumberOfJobNetworks)].DefaultDNSnGW = true

	return jobNetworks
}

func (n *networksAssigner) generateSubnets(azs []string) []bftinput.SubnetConfig {
	subnets := []bftinput.SubnetConfig{}

	placedAzs := NewPlacedAZs()
	for !placedAzs.AllPlaced(azs) {
		newAzs := n.randomCombination(azs)
		placedAzs.Place(newAzs)
		subnets = append(subnets, bftinput.SubnetConfig{AvailabilityZones: newAzs})
	}

	return subnets
}

func (n *networksAssigner) generateSubnetsWithoutAzs() []bftinput.SubnetConfig {
	subnets := []bftinput.SubnetConfig{}
	numberOfSubnets := rand.Intn(3) + 1 // up to 3

	for i := 0; i < numberOfSubnets; i++ {
		subnets = append(subnets, bftinput.SubnetConfig{})
	}

	return subnets
}

type JobsOnNetwork struct {
	Jobs           []bftinput.Job
	TotalInstances int
}

func (n *networksAssigner) aggregateNetworkJobs(jobs []bftinput.Job) map[string]JobsOnNetwork {
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

func (n *networksAssigner) assignStaticIps(networks []bftinput.NetworkConfig, jobs []bftinput.Job) {
	jobsOnNetworks := n.aggregateNetworkJobs(jobs)
	vipIpPool := n.ipPoolProvider.NewIpPool(254)

	// only use 1 network with static IPs because it is hard to generate multiple networks with
	// static IPs that can be distributed evenly across azs
	hasNetworkWithStaticIps := false

	for k, network := range networks {
		jobsOnNetwork := jobsOnNetworks[network.Name]

		if network.Type == "manual" {
			for s, _ := range network.Subnets {
				ipPool := n.ipPoolProvider.NewIpPool(jobsOnNetwork.TotalInstances)
				networks[k].Subnets[s].IpPool = ipPool
			}

			for _, job := range jobsOnNetwork.Jobs {
				if !hasNetworkWithStaticIps && n.staticIpDecider.IsYes() { // use static IPs
					hasNetworkWithStaticIps = true
					for ji := 0; ji < job.Instances; ji++ {
						subnetIpPool, found := n.findIpPoolWithJobAZ(networks[k].Subnets, job.AvailabilityZones)
						if found {
							staticIp, _ := subnetIpPool.NextStaticIp()
							for j, jobNetwork := range job.Networks {
								if jobNetwork.Name == network.Name {
									job.Networks[j].StaticIps = append(job.Networks[j].StaticIps, staticIp)
								}
							}
						}
					}
				}
			}
		} else if network.Type == "vip" {
			for _, job := range jobsOnNetwork.Jobs {
				for j, jobNetwork := range job.Networks {
					if jobNetwork.Name == network.Name {
						for ji := 0; ji < job.Instances; ji++ {
							staticIp, _ := vipIpPool.NextStaticIp()
							job.Networks[j].StaticIps = append(job.Networks[j].StaticIps, staticIp)
						}
					}
				}
			}
		}
	}
}

func (n *networksAssigner) findIpPoolWithJobAZ(subnets []bftinput.SubnetConfig, azs []string) (*bftinput.IpPool, bool) {
	shuffledSubnetIdxs := rand.Perm(len(subnets))
	shuffledSubnets := []bftinput.SubnetConfig{}
	for _, i := range shuffledSubnetIdxs {
		shuffledSubnets = append(shuffledSubnets, subnets[i])
	}

	for i, subnet := range shuffledSubnets {
		for _, subnetAz := range subnet.AvailabilityZones {
			for _, jobAz := range azs {
				if subnetAz == jobAz {
					return shuffledSubnets[i].IpPool, true
				}
			}
		}
	}

	return &bftinput.IpPool{}, false
}
