package network

import (
	"math/rand"

	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
)

type Assigner interface {
	Assign(input bftinput.Input) bftinput.Input
}

type assigner struct {
	networks       [][]string
	nameGenerator  bftnamegen.NameGenerator
	ipPoolProvider IpPoolProvider
	decider        bftdecider.Decider
}

func NewAssigner(
	networks [][]string,
	nameGenerator bftnamegen.NameGenerator,
	ipPoolProvider IpPoolProvider,
	decider bftdecider.Decider,
) Assigner {
	return &assigner{
		networks:       networks,
		nameGenerator:  nameGenerator,
		ipPoolProvider: ipPoolProvider,
		decider:        decider,
	}
}

func (n *assigner) Assign(input bftinput.Input) bftinput.Input {
	// 1. Generate Networks with/without AZs (network with types)
	// 2. Assign networks to each job (network with AZs) (job with network name)
	// 3. Generate IP specs for each network (network with IP specs)
	// 4. Aggregate instances to compute static IPs (network with static IPs) (job with static IP)

	n.ipPoolProvider.Reset()

	networkPoolWithAzs := []bftinput.NetworkConfig{}
	var networkTypes []string

	if len(input.CloudConfig.AvailabilityZones) > 0 {
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
				networkPoolWithAzs[k].Subnets = n.generateSubnets(input.CloudConfig.AvailabilityZones)
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

	for j, job := range input.Jobs {
		if job.AvailabilityZones == nil {
			input.Jobs[j].Networks = n.generateJobNetworks(networkPoolWithoutAzs)
		} else {
			input.Jobs[j].Networks = n.generateJobNetworks(networkPoolWithAzs)
		}
	}

	allNetworks := append(networkPoolWithAzs, networkPoolWithoutAzs...)
	n.assignStaticIps(allNetworks, input.Jobs)

	nonVipNetworks := []bftinput.NetworkConfig{}

	for _, network := range allNetworks {
		input.CloudConfig.Networks = append(input.CloudConfig.Networks, network)

		if network.Type != "vip" {
			nonVipNetworks = append(nonVipNetworks, network)
		}
	}

	compilationNetwork := nonVipNetworks[rand.Intn(len(nonVipNetworks))]
	input.CloudConfig.CompilationNetwork = compilationNetwork.Name
	azs := []string{}
	for _, s := range compilationNetwork.Subnets {
		azs = append(azs, s.AvailabilityZones...)
	}
	if len(azs) > 0 {
		input.CloudConfig.CompilationAvailabilityZone = azs[rand.Intn(len(azs))]
	}

	return input
}

func (n *assigner) generateJobNetworks(networkPool []bftinput.NetworkConfig) []bftinput.JobNetworkConfig {
	jobNetworks := []bftinput.JobNetworkConfig{}

	nonVipNetworks := []bftinput.NetworkConfig{}
	vipNetworks := []bftinput.NetworkConfig{}
	for _, network := range networkPool {
		if network.Type == "vip" {
			vipNetworks = append(vipNetworks, network)
		} else {
			nonVipNetworks = append(nonVipNetworks, network)
		}
	}

	totalNumberOfNonVipNetworks := rand.Intn(len(nonVipNetworks)) + 1 // can not be 0

	networksToPick := rand.Perm(len(nonVipNetworks))[:totalNumberOfNonVipNetworks]
	for _, k := range networksToPick {
		network := nonVipNetworks[k]
		jobNetworks = append(jobNetworks, bftinput.JobNetworkConfig{Name: network.Name})
	}

	jobNetworks[rand.Intn(len(jobNetworks))].DefaultDNSnGW = true

	if len(vipNetworks) != 0 {
		totalNumberOfVipNetworks := rand.Intn(len(vipNetworks)) // can be 0
		networksToPick = rand.Perm(len(vipNetworks))[:totalNumberOfVipNetworks]
		for _, k := range networksToPick {
			network := vipNetworks[k]
			jobNetworks = append(jobNetworks, bftinput.JobNetworkConfig{Name: network.Name})
		}
	}

	if len(jobNetworks) == 1 && !n.decider.IsYes() {
		// if we only have one network on job, we don't necessarily need to specify default DNS n GW
		jobNetworks[0].DefaultDNSnGW = false
	}

	return jobNetworks
}

func (n *assigner) generateSubnets(azs []bftinput.AvailabilityZone) []bftinput.SubnetConfig {
	subnets := []bftinput.SubnetConfig{}

	azNames := []string{}
	for _, az := range azs {
		azNames = append(azNames, az.Name)
	}

	placedAzs := NewPlacedAZs()
	for !placedAzs.AllPlaced(azNames) {
		newAzs := n.randomCombination(azNames)
		placedAzs.Place(newAzs)
		subnets = append(subnets, bftinput.SubnetConfig{AvailabilityZones: newAzs})
	}

	return subnets
}

func (n *assigner) generateSubnetsWithoutAzs() []bftinput.SubnetConfig {
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

func (n *assigner) aggregateNetworkJobs(jobs []bftinput.Job) map[string]JobsOnNetwork {
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

func (n *assigner) randomCombination(items []string) []string {
	combination := []string{}
	totalNumberOfItems := rand.Intn(len(items)) + 1
	itemsToPick := rand.Perm(len(items))[:totalNumberOfItems]
	for _, i := range itemsToPick {
		combination = append(combination, items[i])
	}

	return combination
}

func (n *assigner) assignStaticIps(networks []bftinput.NetworkConfig, jobs []bftinput.Job) {
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
				if !hasNetworkWithStaticIps && n.decider.IsYes() { // use static IPs
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

func (n *assigner) findIpPoolWithJobAZ(subnets []bftinput.SubnetConfig, azs []string) (*bftinput.IpPool, bool) {
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
