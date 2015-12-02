package deployment

import (
	"fmt"
	"math/rand"
	"sort"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
)

type IpPoolProvider interface {
	NewIpPool(numOfNeededIPs int) *bftinput.IpPool
	Reset()
}

type ipPoolProvider struct {
	called             int
	gatewayFourthOctet int
}

func NewIpPoolProvider() IpPoolProvider {
	return &ipPoolProvider{}
}

func (p *ipPoolProvider) NewIpPool(numOfNeededIPs int) *bftinput.IpPool {
	if numOfNeededIPs == 0 {
		numOfNeededIPs = rand.Intn(10)
	}

	if p.gatewayFourthOctet == 1 {
		p.gatewayFourthOctet = 254
	} else {
		p.gatewayFourthOctet = 1
	}

	prefix := fmt.Sprintf("192.168.%d", p.called)
	p.called += 1

	ipRange := fmt.Sprintf("%s.0/24", prefix)
	gateway := fmt.Sprintf("%s.%d", prefix, p.gatewayFourthOctet)

	numberOfReservedBorders := rand.Intn(6) // up to 6 borders of reserved ranges

	usedIps := []int{}
	reservedBorders := []int{}

	for _, i := range rand.Perm(254) {
		if i != 0 && i != p.gatewayFourthOctet {
			if len(usedIps) < numOfNeededIPs {
				usedIps = append(usedIps, i)
			} else if len(reservedBorders) < numberOfReservedBorders {
				reservedBorders = append(reservedBorders, i)
			} else {
				break
			}
		}
	}

	sort.Ints(usedIps)
	sort.Ints(reservedBorders)

	decider := NewRandomDecider()
	reservedRangeGenerator := NewReservedRangeGenerator(prefix, decider)
	reservedRanges := reservedRangeGenerator.Generate(usedIps, reservedBorders)

	availableIps := []string{}
	shuffledUsedIpsIndeces := rand.Perm(len(usedIps))
	for _, ipIndex := range shuffledUsedIpsIndeces {
		availableIps = append(availableIps, fmt.Sprintf("%s.%d", prefix, usedIps[ipIndex]))
	}

	return &bftinput.IpPool{
		IpRange:      ipRange,
		Gateway:      gateway,
		Reserved:     reservedRanges,
		AvailableIps: availableIps,
	}
}

func (p *ipPoolProvider) Reset() {
	p.called = 0
}
