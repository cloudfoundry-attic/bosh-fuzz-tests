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
	for foundJobIdx, job := range input.Jobs {
		previousJob, found := previousInput.FindJobByName(job.Name)
		if found {
			if len(previousJob.AvailabilityZones) == 0 && len(job.AvailabilityZones) > 0 {
				staticIPs := f.sameStaticIps(job, previousJob, input)
				for _, ip := range staticIPs {
					f.assignMigratedFromBasedOnIp(ip, &input.Jobs[foundJobIdx])
				}
			}
		}
	}

	return input
}

func (f *fixedMigratedFrom) assignMigratedFromBasedOnIp(ip staticIPInfo, jobToUpdate *bftinput.Job) {
	for _, subnet := range ip.Network.Subnets {
		if subnet.IpPool.Contains(ip.IP) {
			jobToUpdate.MigratedFrom = []bftinput.MigratedFromConfig{
				{
					Name:             jobToUpdate.Name,
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

func (f *fixedMigratedFrom) sameStaticIps(job bftinput.Job, previousJob bftinput.Job, input bftinput.Input) []staticIPInfo {
	ips := []staticIPInfo{}
	for _, network := range job.Networks {
		previousNetwork, networkFound := previousJob.FindNetworkByName(network.Name)
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
