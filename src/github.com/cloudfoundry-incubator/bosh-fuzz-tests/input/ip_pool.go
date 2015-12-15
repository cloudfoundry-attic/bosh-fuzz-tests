package input

import "fmt"

type IpPool struct {
	prefix               string
	IpRange              string
	Gateway              string
	Reserved             []string
	Static               []string
	lastReservedStaticIp int
	reservedStaticIps    map[string]bool
}

func NewIpPool(
	prefix string,
	gatewayFourthOctet int,
	reserved []string,
) *IpPool {
	return &IpPool{
		prefix:               prefix,
		IpRange:              fmt.Sprintf("%s.0/24", prefix),
		Gateway:              fmt.Sprintf("%s.%d", prefix, gatewayFourthOctet),
		Reserved:             reserved,
		Static:               []string{fmt.Sprintf("%s.200-%s.253", prefix, prefix)},
		lastReservedStaticIp: 200,
		reservedStaticIps:    map[string]bool{},
	}
}

func (i *IpPool) NextStaticIp() (string, error) {
	var staticIp string

	for {
		staticIp = fmt.Sprintf("%s.%d", i.prefix, i.lastReservedStaticIp)
		i.lastReservedStaticIp += 1

		if _, ok := i.reservedStaticIps[staticIp]; !ok {
			break
		}
	}

	return staticIp, nil
}

func (i *IpPool) ReserveStaticIp(ip string) {
	i.reservedStaticIps[ip] = true
}
