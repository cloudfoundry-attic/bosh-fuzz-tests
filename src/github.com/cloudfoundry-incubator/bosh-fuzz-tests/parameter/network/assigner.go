package network

import (
	"math/rand"

	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bftnamegen "github.com/cloudfoundry-incubator/bosh-fuzz-tests/name_generator"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Assigner interface {
	Assign(input bftinput.Input, previousInput bftinput.Input) bftinput.Input
}

type assigner struct {
	networks       [][]string
	nameGenerator  bftnamegen.NameGenerator
	ipPoolProvider IpPoolProvider
	decider        bftdecider.Decider
	logger         boshlog.Logger
}

func NewAssigner(
	networks [][]string,
	nameGenerator bftnamegen.NameGenerator,
	ipPoolProvider IpPoolProvider,
	decider bftdecider.Decider,
	logger boshlog.Logger,
) Assigner {
	return &assigner{
		networks:       networks,
		nameGenerator:  nameGenerator,
		ipPoolProvider: ipPoolProvider,
		decider:        decider,
		logger:         logger,
	}
}

func (n *assigner) Assign(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	// 1. Generate Networks with/without AZs (network with types)
	// 2. Assign networks to each instanceGroup (network with AZs) (instanceGroup with network name)
	// 3. Generate IP specs for each network (network with IP specs)
	// 4. Aggregate instances to compute static IPs (network with static IPs) (instanceGroup with static IP)

	n.ipPoolProvider.Reset()

	ipRangeToStaticIps := NewIpRangeToStaticIps(previousInput)

	networkPoolWithAzs := []bftinput.NetworkConfig{}
	var networkTypes []string

	networkReuser := NewReuser(
		input.CloudConfig.Networks,
		n.decider,
		n.nameGenerator,
	)

	if len(input.CloudConfig.AvailabilityZones) > 0 {
		networkTypes = n.networks[rand.Intn(len(n.networks))]

		for _, networkType := range networkTypes {
			network := networkReuser.CreateNetwork(networkType)
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
		network := networkReuser.CreateNetwork(networkType)
		networkPoolWithoutAzs = append(networkPoolWithoutAzs, network)
	}

	for k, network := range networkPoolWithoutAzs {
		if network.Type != "vip" {
			networkPoolWithoutAzs[k].Subnets = n.generateSubnetsWithoutAzs()
		}
	}

	for j, instanceGroup := range input.InstanceGroups {
		if instanceGroup.AvailabilityZones == nil {
			input.InstanceGroups[j].Networks = n.generateInstanceGroupNetworks(networkPoolWithoutAzs)
		} else {
			input.InstanceGroups[j].Networks = n.generateInstanceGroupNetworks(networkPoolWithAzs)
		}
	}

	allNetworks := append(networkPoolWithAzs, networkPoolWithoutAzs...)
	n.assignStaticIps(allNetworks, input.InstanceGroups, ipRangeToStaticIps, previousInput)

	nonVipNetworks := []bftinput.NetworkConfig{}
	input.CloudConfig.Networks = []bftinput.NetworkConfig{}

	for _, network := range allNetworks {
		input.CloudConfig.Networks = append(input.CloudConfig.Networks, network)

		if network.Type != "vip" {
			nonVipNetworks = append(nonVipNetworks, network)
		}
	}

	compilationNetwork := nonVipNetworks[rand.Intn(len(nonVipNetworks))]
	input.CloudConfig.Compilation.Network = compilationNetwork.Name
	azs := []string{}
	for _, s := range compilationNetwork.Subnets {
		azs = append(azs, s.AvailabilityZones...)
	}
	if len(azs) > 0 {
		input.CloudConfig.Compilation.AvailabilityZone = azs[rand.Intn(len(azs))]
	} else {
		input.CloudConfig.Compilation.AvailabilityZone = ""
	}

	return input
}

func (n *assigner) generateInstanceGroupNetworks(networkPool []bftinput.NetworkConfig) []bftinput.InstanceGroupNetworkConfig {
	instanceGroupNetworks := []bftinput.InstanceGroupNetworkConfig{}

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
		instanceGroupNetworks = append(instanceGroupNetworks, bftinput.InstanceGroupNetworkConfig{Name: network.Name})
	}

	instanceGroupNetworks[rand.Intn(len(instanceGroupNetworks))].DefaultDNSnGW = true

	if len(vipNetworks) != 0 {
		totalNumberOfVipNetworks := rand.Intn(len(vipNetworks)) // can be 0
		networksToPick = rand.Perm(len(vipNetworks))[:totalNumberOfVipNetworks]
		for _, k := range networksToPick {
			network := vipNetworks[k]
			instanceGroupNetworks = append(instanceGroupNetworks, bftinput.InstanceGroupNetworkConfig{Name: network.Name})
		}
	}

	if len(instanceGroupNetworks) == 1 && !n.decider.IsYes() {
		// if we only have one network on instanceGroup, we don't necessarily need to specify default DNS n GW
		instanceGroupNetworks[0].DefaultDNSnGW = false
	}

	return instanceGroupNetworks
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

type InstanceGroupsOnNetwork struct {
	InstanceGroups []bftinput.InstanceGroup
	TotalInstances int
}

func (n *assigner) aggregateNetworkInstanceGroups(instanceGroups []bftinput.InstanceGroup) map[string]InstanceGroupsOnNetwork {
	instanceGroupsOnNetworks := map[string]InstanceGroupsOnNetwork{}

	for _, instanceGroup := range instanceGroups {
		for _, instanceGroupNetwork := range instanceGroup.Networks {
			instanceGroupsOnNetworks[instanceGroupNetwork.Name] = InstanceGroupsOnNetwork{
				InstanceGroups: append(instanceGroupsOnNetworks[instanceGroupNetwork.Name].InstanceGroups, instanceGroup),
				TotalInstances: instanceGroupsOnNetworks[instanceGroupNetwork.Name].TotalInstances + instanceGroup.Instances,
			}
		}
	}

	return instanceGroupsOnNetworks
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

func (n *assigner) assignStaticIps(networks []bftinput.NetworkConfig, instanceGroups []bftinput.InstanceGroup, ipRangeToStaticIps IpRangeToStaticIps, previousInput bftinput.Input) {
	instanceGroupsOnNetworks := n.aggregateNetworkInstanceGroups(instanceGroups)
	vipIpPool := n.ipPoolProvider.NewIpPool(254)
	ipRangeToStaticIps.ReserveStaticIpsInPool(vipIpPool)

	for k, network := range networks {
		instanceGroupsOnNetwork := instanceGroupsOnNetworks[network.Name]

		if network.Type == "manual" {
			for s, _ := range network.Subnets {
				ipPool := n.ipPoolProvider.NewIpPool(instanceGroupsOnNetwork.TotalInstances)
				ipRangeToStaticIps.ReserveStaticIpsInPool(ipPool)
				networks[k].Subnets[s].IpPool = ipPool
			}

			for _, instanceGroup := range instanceGroupsOnNetwork.InstanceGroups {
				// only use 1 network with static IPs per instanceGroup because it is hard to generate multiple networks with
				// static IPs that can be distributed evenly across azs
				hasNetworkWithStaticIps := false

				for j, instanceGroupNetwork := range instanceGroup.Networks {
					instanceGroup.Networks[j].StaticIps = []string{}

					if !hasNetworkWithStaticIps && n.decider.IsYes() { // use static IPs
						hasNetworkWithStaticIps = true

						ipsToReuseFromPreviousDeploy := n.getInstanceGroupStaticIpsToReuse(previousInput, instanceGroup.Name, instanceGroupNetwork.Name)
						n.logger.Debug("networkAssigner", "Reusing IPs from previous deploy %#v for instance-group %s", ipsToReuseFromPreviousDeploy, instanceGroup.Name)

						for ji := 0; ji < instanceGroup.Instances; ji++ {
							subnetIpPool, found := n.findIpPoolWithInstanceGroupAZ(networks[k].Subnets, instanceGroup.AvailabilityZones)
							if found {
								var staticIp string
								if len(ipsToReuseFromPreviousDeploy) > 0 {
									var ipToReuse string
									ipToReuse, ipsToReuseFromPreviousDeploy = ipsToReuseFromPreviousDeploy[0], ipsToReuseFromPreviousDeploy[1:]
									if subnetIpPool.Contains(ipToReuse) {
										staticIp = ipToReuse
									}
								}

								if staticIp == "" {
									staticIp, _ = subnetIpPool.NextStaticIp()
								}

								if instanceGroupNetwork.Name == network.Name {
									instanceGroup.Networks[j].StaticIps = append(instanceGroup.Networks[j].StaticIps, staticIp)
								}
							}
						}
					}
				}
			}
		} else if network.Type == "vip" {
			for _, instanceGroup := range instanceGroupsOnNetwork.InstanceGroups {
				for j, instanceGroupNetwork := range instanceGroup.Networks {
					if instanceGroupNetwork.Name == network.Name {
						for ji := 0; ji < instanceGroup.Instances; ji++ {
							staticIp, _ := vipIpPool.NextStaticIp()
							instanceGroup.Networks[j].StaticIps = append(instanceGroup.Networks[j].StaticIps, staticIp)
						}
					}
				}
			}
		}
	}
}

func (n *assigner) getInstanceGroupStaticIpsToReuse(previousInput bftinput.Input, instanceGroupName string, networkName string) []string {
	staticIps := []string{}

	previousInstanceGroup, found := previousInput.FindInstanceGroupByName(instanceGroupName)
	if !found {
		return staticIps
	}

	for _, instanceGroupNetwork := range previousInstanceGroup.Networks {
		if instanceGroupNetwork.Name == networkName {
			for _, ip := range instanceGroupNetwork.StaticIps {
				staticIps = append(staticIps, ip)
			}
		}
	}

	if len(staticIps) == 0 {
		return staticIps
	}

	shuffledStaticIPsIdsx := rand.Perm(len(staticIps))
	ipsToReuse := rand.Intn(len(staticIps))

	shuffledStaticIps := []string{}

	for i := 0; i < ipsToReuse; i++ {
		shuffledStaticIps = append(shuffledStaticIps, staticIps[shuffledStaticIPsIdsx[i]])
	}

	return shuffledStaticIps
}

func (n *assigner) findIpPoolWithInstanceGroupAZ(subnets []bftinput.SubnetConfig, azs []string) (*bftinput.IpPool, bool) {
	shuffledSubnetIdxs := rand.Perm(len(subnets))
	shuffledSubnets := []bftinput.SubnetConfig{}
	for _, i := range shuffledSubnetIdxs {
		shuffledSubnets = append(shuffledSubnets, subnets[i])
	}

	for i, subnet := range shuffledSubnets {
		if len(subnet.AvailabilityZones) == 0 && len(azs) == 0 {
			return shuffledSubnets[i].IpPool, true
		}

		for _, subnetAz := range subnet.AvailabilityZones {
			for _, instanceGroupAz := range azs {
				if subnetAz == instanceGroupAz {
					return shuffledSubnets[i].IpPool, true
				}
			}
		}
	}

	return &bftinput.IpPool{}, false
}
