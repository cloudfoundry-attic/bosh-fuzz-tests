package network

import (
	"fmt"
	"math/rand"
	"sort"

	bftdecider "github.com/cloudfoundry-incubator/bosh-fuzz-tests/decider"
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

	numberOfReservedBorders := rand.Intn(6) // up to 6 borders of reserved ranges

	usedIps := []int{}
	reservedBorders := []int{}

	firstStaticIp := 200

	for _, i := range rand.Perm(firstStaticIp) {
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

	decider := bftdecider.NewRandomDecider()
	reservedRangeGenerator := NewReservedRangeGenerator(prefix, decider)
	reservedRanges := reservedRangeGenerator.Generate(usedIps, reservedBorders)

	return bftinput.NewIpPool(
		prefix,
		p.gatewayFourthOctet,
		reservedRanges,
	)
}

func (p *ipPoolProvider) Reset() {
	p.called = 0
}
