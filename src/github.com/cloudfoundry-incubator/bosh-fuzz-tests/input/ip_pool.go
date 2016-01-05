package input

import (
	"fmt"
	"math/rand"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type IpPool struct {
	prefix            string
	IpRange           string
	Gateway           string
	Reserved          []string
	Static            []string
	staticIps         []int
	reservedStaticIps map[string]bool
}

func NewIpPool(
	prefix string,
	gatewayFourthOctet int,
	reserved []string,
) *IpPool {

	// shuffle static ips 200-253
	staticIps := []int{}
	shuffledIpIdxs := rand.Perm(54)
	for _, idx := range shuffledIpIdxs {
		staticIps = append(staticIps, 200+idx)
	}

	return &IpPool{
		prefix:            prefix,
		IpRange:           fmt.Sprintf("%s.0/24", prefix),
		Gateway:           fmt.Sprintf("%s.%d", prefix, gatewayFourthOctet),
		Reserved:          reserved,
		Static:            []string{fmt.Sprintf("%s.200-%s.253", prefix, prefix)},
		staticIps:         staticIps,
		reservedStaticIps: map[string]bool{},
	}
}

func (i *IpPool) NextStaticIp() (string, error) {
	var staticIp string
	var nextStaticIp int

	for {
		if len(i.staticIps) == 0 {
			return "", bosherr.Error("No more static IPs available")
		}

		nextStaticIp, i.staticIps = i.staticIps[0], i.staticIps[1:]
		staticIp = fmt.Sprintf("%s.%d", i.prefix, nextStaticIp)

		if _, ok := i.reservedStaticIps[staticIp]; !ok {
			break
		}
	}

	return staticIp, nil
}

func (i *IpPool) ReserveStaticIp(ip string) {
	i.reservedStaticIps[ip] = true
}

func (i *IpPool) Contains(ip string) bool {
	substring := ip[0:len(i.prefix)+1]
	return substring == (i.prefix + ".")
}
