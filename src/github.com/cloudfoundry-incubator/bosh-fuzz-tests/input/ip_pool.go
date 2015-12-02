package input

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type IpPool struct {
	IpRange      string
	Gateway      string
	Reserved     []string
	Static       []string
	AvailableIps []string
}

func (i *IpPool) NextStaticIp() (string, error) {
	var ip string
	if len(i.AvailableIps) == 0 {
		return "", bosherr.Error("No more available")
	}
	ip, i.AvailableIps = i.AvailableIps[0], i.AvailableIps[1:]
	i.Static = append(i.Static, ip)
	return ip, nil
}
