package network

import (
	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type IpRangeToStaticIps map[string][]string

func NewIpRangeToStaticIps(input bftinput.Input) IpRangeToStaticIps {
	ipRangeToStaticIps := map[string][]string{}

	for _, network := range input.CloudConfig.Networks {
		for _, subnet := range network.Subnets {

			for _, job := range input.Jobs {
				for _, jobNetwork := range job.Networks {
					if jobNetwork.Name == network.Name {
						for _, ip := range jobNetwork.StaticIps {
							ipRangeToStaticIps[subnet.IpPool.IpRange] = append(ipRangeToStaticIps[subnet.IpPool.IpRange], ip)
						}
					}
				}
			}
		}
	}

	return ipRangeToStaticIps
}

func (i IpRangeToStaticIps) ReserveStaticIpsInPool(ipPool *bftinput.IpPool) {
	staticIps, ok := i[ipPool.IpRange]
	if !ok {
		return
	}

	for _, ip := range staticIps {
		ipPool.ReserveStaticIp(ip)
	}
}
