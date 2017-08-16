package parameter

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type fixedMigratedFrom struct {
}

func NewFixedMigratedFrom() Parameter {
	return &fixedMigratedFrom{}
}

func (f *fixedMigratedFrom) Apply(input bftinput.Input, previousInput bftinput.Input) bftinput.Input {
	for foundInstanceGroupIdx, instanceGroup := range input.InstanceGroups {
		previousInstanceGroup, found := previousInput.FindInstanceGroupByName(instanceGroup.Name)
		if found {
			if len(previousInstanceGroup.AvailabilityZones) == 0 && len(instanceGroup.AvailabilityZones) > 0 {
				staticIPs := f.sameStaticIps(instanceGroup, previousInstanceGroup, input)
				for _, ip := range staticIPs {
					f.assignMigratedFromBasedOnIp(ip, &input.InstanceGroups[foundInstanceGroupIdx])
				}
			}
		}
	}

	return input
}

func (f *fixedMigratedFrom) assignMigratedFromBasedOnIp(ip staticIPInfo, instanceGroupToUpdate *bftinput.InstanceGroup) {
	for _, subnet := range ip.Network.Subnets {
		if subnet.IpPool.Contains(ip.IP) {
			instanceGroupToUpdate.MigratedFrom = []bftinput.MigratedFromConfig{
				{
					Name:             instanceGroupToUpdate.Name,
					AvailabilityZone: subnet.AvailabilityZones[0],
				},
			}

			return
		}
	}
}

type staticIPInfo struct {
	IP      string
	Network bftinput.NetworkConfig
}

func (f *fixedMigratedFrom) sameStaticIps(instanceGroup bftinput.InstanceGroup, previousInstanceGroup bftinput.InstanceGroup, input bftinput.Input) []staticIPInfo {
	ips := []staticIPInfo{}
	for _, network := range instanceGroup.Networks {
		previousNetwork, networkFound := previousInstanceGroup.FindNetworkByName(network.Name)
		if networkFound {
			for _, currentIP := range network.StaticIps {
				for _, prevIP := range previousNetwork.StaticIps {
					if prevIP == currentIP {
						cloudNetwork, cloudNetworkFound := input.FindNetworkByName(network.Name)
						if cloudNetworkFound {
							ip := staticIPInfo{
								IP:      currentIP,
								Network: cloudNetwork,
							}
							ips = append(ips, ip)
						}
					}
				}
			}
		}
	}
	return ips
}
